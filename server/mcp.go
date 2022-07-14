package server

import (
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/champly/mcpserver/resource"
	"github.com/champly/mcpserver/resource/mock"
	"github.com/champly/mcpserver/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	mcp "istio.io/api/mcp/v1alpha1"
	"istio.io/istio/pkg/mcp/rate"
	"istio.io/istio/pkg/mcp/server"
	"istio.io/istio/pkg/mcp/source"
	"istio.io/istio/pkg/mcp/testing/monitoring"
	"k8s.io/klog/v2"
)

var (
	once           = &sync.Once{}
	defaultVersion = "v1"
)

// Server sink mcp server
type sinkServer struct {
	grpc.ServerStream
	*sourceHarness

	opt         *option
	logicServer types.LogicServer
}

// NewMCPServer build sink mcp server
func newMCPServer(opt *option) *sinkServer {
	serv := &sinkServer{
		opt:           opt,
		sourceHarness: newSourceHarness(),
	}

	serv.logicServer = mock.New(serv)
	return serv
}

func (s *sinkServer) Start(stop <-chan struct{}) {
	options := &source.Options{
		Watcher:            s,
		CollectionsOptions: source.CollectionOptionsFromSlice(resource.GetAllResource()),
		Reporter:           monitoring.NewInMemoryStatsContext(),
		ConnRateLimiter:    rate.NewRateLimiter(s.opt.Freq, s.opt.BurstSize),
	}
	serverOptions := &source.ServerOptions{
		AuthChecker: &server.AllowAllChecker{},
		RateLimiter: rate.NewRateLimiter(s.opt.Freq, s.opt.BurstSize).Create(),
	}

	srv := source.NewServer(options, serverOptions)

	addr := fmt.Sprintf("%s:%d", s.opt.Address, s.opt.GRPCPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		klog.Fatalf("listen %s failed:%s", err)
	}
	klog.Infoln("listen", addr)

	serv := grpc.NewServer(getServerGrpcOptions()...)
	mcp.RegisterResourceSourceServer(serv, srv)

	go func() {
		if err := serv.Serve(l); err != nil {
			klog.Errorf("grcp Serve errr:%s", err)
		}
	}()

	go func() {
		s.logicServer.Start(stop)
	}()

	<-stop
	serv.Stop()
}

type sourceHarness struct {
	PushFunc source.PushResponseFunc
}

func newSourceHarness() *sourceHarness {
	return &sourceHarness{}
}

func (h *sourceHarness) Watch(req *source.Request, pushResponse source.PushResponseFunc, peerAddr string) source.CancelWatchFunc {
	h.PushFunc = pushResponse

	snapHandler, ok := resource.FactorySnap[req.Collection]
	if !ok {
		if req.VersionInfo == defaultVersion {
			klog.Infof("needless resource ack:%+v", req)
			return nil
		}
		h.PushFunc(&source.WatchResponse{
			Collection: req.Collection,
			Version:    defaultVersion,
			Request:    req,
		})
		return nil
	}

	snap, err := snapHandler.All()
	if err != nil {
		klog.Fatalf("get all %s resource failed: %s", req.Collection, err)
	}

	if snap == nil {
		return nil
	}

	if snap.Version == req.VersionInfo {
		return nil
	}

	h.PushFunc(&source.WatchResponse{
		Collection: req.Collection,
		Version:    snap.Version,
		Resources:  snap.Resources,
		Request:    req,
	})
	return nil
}

func (h *sourceHarness) Push(collection string, snap *types.ResourceSnap) {
	if h.PushFunc == nil {
		return
	}
	h.PushFunc(&source.WatchResponse{
		Collection: collection,
		Version:    snap.Version,
		Resources:  snap.Resources,
	})
}

func getServerGrpcOptions() []grpc.ServerOption {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(uint32(1024)),
		grpc.MaxRecvMsgSize(int(4*1024*1024)),
		grpc.InitialWindowSize(int32(1024*1024)),
		grpc.InitialConnWindowSize(int32(1024*1024)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:                  30 * time.Second,
			Timeout:               10 * time.Second,
			MaxConnectionIdle:     time.Duration(math.MaxInt64),
			MaxConnectionAgeGrace: 10 * time.Second,
		}),
		// Relax keepalive enforcement policy requirements to avoid dropping connections due to too many pings.
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	return grpcOptions
}
