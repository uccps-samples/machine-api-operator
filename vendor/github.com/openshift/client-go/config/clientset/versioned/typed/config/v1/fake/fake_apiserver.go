// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	configv1 "github.com/uccps-samples/api/config/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAPIServers implements APIServerInterface
type FakeAPIServers struct {
	Fake *FakeConfigV1
}

var apiserversResource = schema.GroupVersionResource{Group: "config.openshift.io", Version: "v1", Resource: "apiservers"}

var apiserversKind = schema.GroupVersionKind{Group: "config.openshift.io", Version: "v1", Kind: "APIServer"}

// Get takes name of the aPIServer, and returns the corresponding aPIServer object, and an error if there is any.
func (c *FakeAPIServers) Get(ctx context.Context, name string, options v1.GetOptions) (result *configv1.APIServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(apiserversResource, name), &configv1.APIServer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.APIServer), err
}

// List takes label and field selectors, and returns the list of APIServers that match those selectors.
func (c *FakeAPIServers) List(ctx context.Context, opts v1.ListOptions) (result *configv1.APIServerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(apiserversResource, apiserversKind, opts), &configv1.APIServerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &configv1.APIServerList{ListMeta: obj.(*configv1.APIServerList).ListMeta}
	for _, item := range obj.(*configv1.APIServerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aPIServers.
func (c *FakeAPIServers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(apiserversResource, opts))
}

// Create takes the representation of a aPIServer and creates it.  Returns the server's representation of the aPIServer, and an error, if there is any.
func (c *FakeAPIServers) Create(ctx context.Context, aPIServer *configv1.APIServer, opts v1.CreateOptions) (result *configv1.APIServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(apiserversResource, aPIServer), &configv1.APIServer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.APIServer), err
}

// Update takes the representation of a aPIServer and updates it. Returns the server's representation of the aPIServer, and an error, if there is any.
func (c *FakeAPIServers) Update(ctx context.Context, aPIServer *configv1.APIServer, opts v1.UpdateOptions) (result *configv1.APIServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(apiserversResource, aPIServer), &configv1.APIServer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.APIServer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAPIServers) UpdateStatus(ctx context.Context, aPIServer *configv1.APIServer, opts v1.UpdateOptions) (*configv1.APIServer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(apiserversResource, "status", aPIServer), &configv1.APIServer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.APIServer), err
}

// Delete takes name of the aPIServer and deletes it. Returns an error if one occurs.
func (c *FakeAPIServers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(apiserversResource, name), &configv1.APIServer{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAPIServers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(apiserversResource, listOpts)

	_, err := c.Fake.Invokes(action, &configv1.APIServerList{})
	return err
}

// Patch applies the patch and returns the patched aPIServer.
func (c *FakeAPIServers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configv1.APIServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(apiserversResource, name, pt, data, subresources...), &configv1.APIServer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.APIServer), err
}
