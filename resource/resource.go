package resource

import (
	"time"

	"github.com/champly/mcpserver/types"
	"k8s.io/klog/v2"
)

var FactorySnap map[string]types.Snap

func init() {
	FactorySnap = make(map[string]types.Snap)

}

// Registry registry need care resource
func Registry(ele string, snap types.Snap) {
	if _, ok := FactorySnap[ele]; ok {
		klog.Errorf("duplicate registry resource:%s", ele)
		return
	}
	FactorySnap[ele] = snap
}

// GetAllResource get all pilot watch resource
func GetAllResource() []string {
	// // istio 1.5+
	// for _, col := range collections.Pilot.All() {
	// 	cols = append(cols, col.Name().String())
	// }

	// istio version < 1.5
	cols := []string{"istio/authentication/v1alpha1/meshpolicies", "istio/authentication/v1alpha1/policies", "istio/config/v1alpha2/adapters", "istio/config/v1alpha2/httpapispecbindings", "istio/config/v1alpha2/httpapispecs", "istio/config/v1alpha2/legacy/apikeys", "istio/config/v1alpha2/legacy/authorizations", "istio/config/v1alpha2/legacy/bypasses", "istio/config/v1alpha2/legacy/checknothings", "istio/config/v1alpha2/legacy/circonuses", "istio/config/v1alpha2/legacy/cloudwatches", "istio/config/v1alpha2/legacy/deniers", "istio/config/v1alpha2/legacy/dogstatsds", "istio/config/v1alpha2/legacy/edges", "istio/config/v1alpha2/legacy/fluentds", "istio/config/v1alpha2/legacy/kubernetesenvs", "istio/config/v1alpha2/legacy/kuberneteses", "istio/config/v1alpha2/legacy/listcheckers", "istio/config/v1alpha2/legacy/listentries", "istio/config/v1alpha2/legacy/logentries", "istio/config/v1alpha2/legacy/memquotas", "istio/config/v1alpha2/legacy/metrics", "istio/config/v1alpha2/legacy/noops", "istio/config/v1alpha2/legacy/opas", "istio/config/v1alpha2/legacy/prometheuses", "istio/config/v1alpha2/legacy/quotas", "istio/config/v1alpha2/legacy/rbacs", "istio/config/v1alpha2/legacy/redisquotas", "istio/config/v1alpha2/legacy/reportnothings", "istio/config/v1alpha2/legacy/signalfxs", "istio/config/v1alpha2/legacy/solarwindses", "istio/config/v1alpha2/legacy/stackdrivers", "istio/config/v1alpha2/legacy/statsds", "istio/config/v1alpha2/legacy/stdios", "istio/config/v1alpha2/legacy/tracespans", "istio/config/v1alpha2/legacy/zipkins", "istio/config/v1alpha2/templates", "istio/mesh/v1alpha1/MeshConfig", "istio/mixer/v1/config/client/quotaspecbindings", "istio/mixer/v1/config/client/quotaspecs", "istio/networking/v1alpha3/destinationrules", "istio/networking/v1alpha3/envoyfilters", "istio/networking/v1alpha3/gateways", "istio/networking/v1alpha3/serviceentries", "istio/networking/v1alpha3/sidecars", "istio/networking/v1alpha3/synthetic/serviceentries", "istio/networking/v1alpha3/virtualservices", "istio/policy/v1beta1/attributemanifests", "istio/policy/v1beta1/handlers", "istio/policy/v1beta1/instances", "istio/policy/v1beta1/rules", "istio/rbac/v1alpha1/clusterrbacconfigs", "istio/rbac/v1alpha1/rbacconfigs", "istio/rbac/v1alpha1/servicerolebindings", "istio/rbac/v1alpha1/serviceroles", "istio/security/v1beta1/authorizationpolicies", "k8s/apps/v1/deployments", "k8s/core/v1/namespaces", "k8s/core/v1/pods", "k8s/core/v1/services", "istio/security/v1beta1/requestauthentications", "istio/security/v1beta1/peerauthentications", "istio/networking/v1alpha3/workloadentries"}

	return cols
}

// BuildVersion build resource snap version
func BuildVersion() string {
	return time.Now().Format("20060102150405")
}
