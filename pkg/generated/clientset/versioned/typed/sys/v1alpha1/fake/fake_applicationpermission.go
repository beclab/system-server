// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "bytetrade.io/web3os/system-server/pkg/apis/sys/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeApplicationPermissions implements ApplicationPermissionInterface
type FakeApplicationPermissions struct {
	Fake *FakeSysV1alpha1
	ns   string
}

var applicationpermissionsResource = v1alpha1.SchemeGroupVersion.WithResource("applicationpermissions")

var applicationpermissionsKind = v1alpha1.SchemeGroupVersion.WithKind("ApplicationPermission")

// Get takes name of the applicationPermission, and returns the corresponding applicationPermission object, and an error if there is any.
func (c *FakeApplicationPermissions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ApplicationPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(applicationpermissionsResource, c.ns, name), &v1alpha1.ApplicationPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationPermission), err
}

// List takes label and field selectors, and returns the list of ApplicationPermissions that match those selectors.
func (c *FakeApplicationPermissions) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ApplicationPermissionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(applicationpermissionsResource, applicationpermissionsKind, c.ns, opts), &v1alpha1.ApplicationPermissionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ApplicationPermissionList{ListMeta: obj.(*v1alpha1.ApplicationPermissionList).ListMeta}
	for _, item := range obj.(*v1alpha1.ApplicationPermissionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested applicationPermissions.
func (c *FakeApplicationPermissions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(applicationpermissionsResource, c.ns, opts))

}

// Create takes the representation of a applicationPermission and creates it.  Returns the server's representation of the applicationPermission, and an error, if there is any.
func (c *FakeApplicationPermissions) Create(ctx context.Context, applicationPermission *v1alpha1.ApplicationPermission, opts v1.CreateOptions) (result *v1alpha1.ApplicationPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(applicationpermissionsResource, c.ns, applicationPermission), &v1alpha1.ApplicationPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationPermission), err
}

// Update takes the representation of a applicationPermission and updates it. Returns the server's representation of the applicationPermission, and an error, if there is any.
func (c *FakeApplicationPermissions) Update(ctx context.Context, applicationPermission *v1alpha1.ApplicationPermission, opts v1.UpdateOptions) (result *v1alpha1.ApplicationPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(applicationpermissionsResource, c.ns, applicationPermission), &v1alpha1.ApplicationPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationPermission), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeApplicationPermissions) UpdateStatus(ctx context.Context, applicationPermission *v1alpha1.ApplicationPermission, opts v1.UpdateOptions) (*v1alpha1.ApplicationPermission, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(applicationpermissionsResource, "status", c.ns, applicationPermission), &v1alpha1.ApplicationPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationPermission), err
}

// Delete takes name of the applicationPermission and deletes it. Returns an error if one occurs.
func (c *FakeApplicationPermissions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(applicationpermissionsResource, c.ns, name, opts), &v1alpha1.ApplicationPermission{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApplicationPermissions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(applicationpermissionsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ApplicationPermissionList{})
	return err
}

// Patch applies the patch and returns the patched applicationPermission.
func (c *FakeApplicationPermissions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ApplicationPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(applicationpermissionsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ApplicationPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationPermission), err
}
