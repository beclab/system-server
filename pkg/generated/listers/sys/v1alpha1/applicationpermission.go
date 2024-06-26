// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ApplicationPermissionLister helps list ApplicationPermissions.
// All objects returned here must be treated as read-only.
type ApplicationPermissionLister interface {
	// List lists all ApplicationPermissions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ApplicationPermission, err error)
	// ApplicationPermissions returns an object that can list and get ApplicationPermissions.
	ApplicationPermissions(namespace string) ApplicationPermissionNamespaceLister
	ApplicationPermissionListerExpansion
}

// applicationPermissionLister implements the ApplicationPermissionLister interface.
type applicationPermissionLister struct {
	indexer cache.Indexer
}

// NewApplicationPermissionLister returns a new ApplicationPermissionLister.
func NewApplicationPermissionLister(indexer cache.Indexer) ApplicationPermissionLister {
	return &applicationPermissionLister{indexer: indexer}
}

// List lists all ApplicationPermissions in the indexer.
func (s *applicationPermissionLister) List(selector labels.Selector) (ret []*v1alpha1.ApplicationPermission, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ApplicationPermission))
	})
	return ret, err
}

// ApplicationPermissions returns an object that can list and get ApplicationPermissions.
func (s *applicationPermissionLister) ApplicationPermissions(namespace string) ApplicationPermissionNamespaceLister {
	return applicationPermissionNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ApplicationPermissionNamespaceLister helps list and get ApplicationPermissions.
// All objects returned here must be treated as read-only.
type ApplicationPermissionNamespaceLister interface {
	// List lists all ApplicationPermissions in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ApplicationPermission, err error)
	// Get retrieves the ApplicationPermission from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ApplicationPermission, error)
	ApplicationPermissionNamespaceListerExpansion
}

// applicationPermissionNamespaceLister implements the ApplicationPermissionNamespaceLister
// interface.
type applicationPermissionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ApplicationPermissions in the indexer for a given namespace.
func (s applicationPermissionNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ApplicationPermission, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ApplicationPermission))
	})
	return ret, err
}

// Get retrieves the ApplicationPermission from the indexer for a given namespace and name.
func (s applicationPermissionNamespaceLister) Get(name string) (*v1alpha1.ApplicationPermission, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("applicationpermission"), name)
	}
	return obj.(*v1alpha1.ApplicationPermission), nil
}
