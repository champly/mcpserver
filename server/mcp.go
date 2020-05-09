package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/champly/mcpserver/resource"
	"github.com/champly/mcpserver/resource/mock"
	"github.com/champly/mcpserver/types"
	"google.golang.org/grpc"
	mcp "istio.io/api/mcp/v1alpha1"
	"istio.io/istio/pkg/mcp/rate"
	"istio.io/istio/pkg/mcp/server"
	"istio.io/istio/pkg/mcp/source"
	"istio.io/istio/pkg/mcp/testing/monitoring"
	"k8s.io/klog"
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

	serv := grpc.NewServer()
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
	once.Do(func() {
		h.PushFunc = pushResponse
	})

	snap, ok := resource.FactorySnap[req.Collection]
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

	resp, err := snap.All(req)
	if err != nil {
		klog.Fatalf("get all %s resource failed: %s", req.Collection, err)
	}
	if resp != nil {
		h.PushFunc(&source.WatchResponse{
			Collection: req.Collection,
			Version:    resp.Version,
			Resources:  resp.Resources,
			Request:    req,
		})
	}
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
