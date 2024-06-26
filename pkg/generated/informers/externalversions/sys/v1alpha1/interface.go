// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "bytetrade.io/web3os/system-server/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// ApplicationPermissions returns a ApplicationPermissionInformer.
	ApplicationPermissions() ApplicationPermissionInformer
	// ProviderRegistries returns a ProviderRegistryInformer.
	ProviderRegistries() ProviderRegistryInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// ApplicationPermissions returns a ApplicationPermissionInformer.
func (v *version) ApplicationPermissions() ApplicationPermissionInformer {
	return &applicationPermissionInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ProviderRegistries returns a ProviderRegistryInformer.
func (v *version) ProviderRegistries() ProviderRegistryInformer {
	return &providerRegistryInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
