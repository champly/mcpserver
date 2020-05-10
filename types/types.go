package types

import (
	mcp "istio.io/api/mcp/v1alpha1"
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
	All() (*ResourceSnap, error)
}

type Source interface {
	Push(collection string, snap *ResourceSnap)
}

type LogicServer interface {
	Start(stop <-chan struct{})
}
