package types

import (
	mcp "istio.io/api/mcp/v1alpha1"
	"istio.io/istio/pkg/mcp/source"
)

// Istio CRD string
const (
	IstioCRDVirtualService  = "istio/networking/v1alpha3/virtualservices"
	IstioCRDDestinationRule = "istio/networking/v1alpha3/destinationrules"
	IstioCRDEnvoyFilter     = "istio/networking/v1alpha3/envoyfilters"
	IstioCRDServiceEntry    = "istio/networking/v1alpha3/serviceentries"
)

type ResourceSnap struct {
	Version   string
	Resources []*mcp.Resource
}

type Snap interface {
	All(req *source.Request) (*ResourceSnap, error)
}

type Source interface {
	Push(collection string, snap *ResourceSnap)
}

type LogicServer interface {
	Start(stop <-chan struct{})
}
