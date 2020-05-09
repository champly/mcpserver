package mock

import (
	"net/http"
	"sync"

	"github.com/champly/mcpserver/resource"
	"github.com/champly/mcpserver/types"
	"github.com/gin-gonic/gin"
	ptypes "github.com/gogo/protobuf/types"
	mcp "istio.io/api/mcp/v1alpha1"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pkg/mcp/source"
	"k8s.io/klog"
)

type mockserviceentry struct {
	l      sync.Mutex
	snap   *types.ResourceSnap
	source types.Source
}

func (se *mockserviceentry) All(req *source.Request) (*types.ResourceSnap, error) {
	se.l.Lock()
	defer se.l.Unlock()

	if se.snap != nil {
		if req.VersionInfo == se.snap.Version {
			klog.Infof("mock resource %s confirm", req.Collection)
			return nil, nil
		}
		return se.snap, nil
	}

	se.createNew()

	return se.snap, nil
}

func (se *mockserviceentry) createNew() {
	se.snap = &types.ResourceSnap{
		Version:   resource.BuildVersion(),
		Resources: []*mcp.Resource{},
	}

	data := &networking.ServiceEntry{
		Hosts: []string{"dubbo-mosn.io.dubbo.DemoService-sayHello"},
		Ports: []*networking.Port{
			{
				Number:   20882,
				Name:     "aabb-server",
				Protocol: "TCP",
			},
		},
		Location:   networking.ServiceEntry_MESH_INTERNAL,
		Resolution: networking.ServiceEntry_STATIC,
		Endpoints: []*networking.WorkloadEntry{
			{
				Address: "10.13.160.40",
			},
			{
				Address: "10.13.160.93",
			},
			{
				Address: "10.13.160.66",
			},
		},
	}
	b, _ := ptypes.MarshalAny(data)

	se.snap.Resources = append(se.snap.Resources, &mcp.Resource{
		Metadata: &mcp.Metadata{
			Name:    "aabb-server",
			Version: resource.BuildVersion(),
		},
		Body: b,
	})
}

func (se *mockserviceentry) Update(c *gin.Context) {
	se.l.Lock()
	defer se.l.Unlock()

	se.snap = &types.ResourceSnap{
		Version:   resource.BuildVersion(),
		Resources: []*mcp.Resource{},
	}

	data := &networking.ServiceEntry{
		Hosts: []string{"dubbo-mosn.io.dubbo.DemoService-sayHello"},
		Ports: []*networking.Port{
			{
				Number:   20882,
				Name:     "aabb-server",
				Protocol: "TCP",
			},
		},
		Location:   networking.ServiceEntry_MESH_INTERNAL,
		Resolution: networking.ServiceEntry_STATIC,
		Endpoints: []*networking.WorkloadEntry{
			{
				Address: "10.13.160.40",
			},
			{
				Address: "10.13.160.93",
			},
			{
				Address: "10.13.160.16",
			},
		},
	}
	b, _ := ptypes.MarshalAny(data)

	se.snap.Resources = append(se.snap.Resources, &mcp.Resource{
		Metadata: &mcp.Metadata{
			Name:    "aabb-server",
			Version: resource.BuildVersion(),
		},
		Body: b,
	})

	se.source.Push(types.IstioCRDServiceEntry, se.snap)

	c.JSON(http.StatusOK, "ok")
}
