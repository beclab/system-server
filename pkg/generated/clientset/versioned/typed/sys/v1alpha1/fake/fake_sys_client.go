// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "bytetrade.io/web3os/system-server/pkg/generated/clientset/versioned/typed/sys/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSysV1alpha1 struct {
	*testing.Fake
}

func (c *FakeSysV1alpha1) ApplicationPermissions(namespace string) v1alpha1.ApplicationPermissionInterface {
	return &FakeApplicationPermissions{c, namespace}
}

func (c *FakeSysV1alpha1) ProviderRegistries(namespace string) v1alpha1.ProviderRegistryInterface {
	return &FakeProviderRegistries{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSysV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
