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

// FakeIngresses implements IngressInterface
type FakeIngresses struct {
	Fake *FakeConfigV1
}

var ingressesResource = schema.GroupVersionResource{Group: "config.uccp.io", Version: "v1", Resource: "ingresses"}

var ingressesKind = schema.GroupVersionKind{Group: "config.uccp.io", Version: "v1", Kind: "Ingress"}

// Get takes name of the ingress, and returns the corresponding ingress object, and an error if there is any.
func (c *FakeIngresses) Get(ctx context.Context, name string, options v1.GetOptions) (result *configv1.Ingress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(ingressesResource, name), &configv1.Ingress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Ingress), err
}

// List takes label and field selectors, and returns the list of Ingresses that match those selectors.
func (c *FakeIngresses) List(ctx context.Context, opts v1.ListOptions) (result *configv1.IngressList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(ingressesResource, ingressesKind, opts), &configv1.IngressList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &configv1.IngressList{ListMeta: obj.(*configv1.IngressList).ListMeta}
	for _, item := range obj.(*configv1.IngressList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ingresses.
func (c *FakeIngresses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(ingressesResource, opts))
}

// Create takes the representation of a ingress and creates it.  Returns the server's representation of the ingress, and an error, if there is any.
func (c *FakeIngresses) Create(ctx context.Context, ingress *configv1.Ingress, opts v1.CreateOptions) (result *configv1.Ingress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(ingressesResource, ingress), &configv1.Ingress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Ingress), err
}

// Update takes the representation of a ingress and updates it. Returns the server's representation of the ingress, and an error, if there is any.
func (c *FakeIngresses) Update(ctx context.Context, ingress *configv1.Ingress, opts v1.UpdateOptions) (result *configv1.Ingress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(ingressesResource, ingress), &configv1.Ingress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Ingress), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeIngresses) UpdateStatus(ctx context.Context, ingress *configv1.Ingress, opts v1.UpdateOptions) (*configv1.Ingress, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(ingressesResource, "status", ingress), &configv1.Ingress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Ingress), err
}

// Delete takes name of the ingress and deletes it. Returns an error if one occurs.
func (c *FakeIngresses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(ingressesResource, name), &configv1.Ingress{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeIngresses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(ingressesResource, listOpts)

	_, err := c.Fake.Invokes(action, &configv1.IngressList{})
	return err
}

// Patch applies the patch and returns the patched ingress.
func (c *FakeIngresses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configv1.Ingress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(ingressesResource, name, pt, data, subresources...), &configv1.Ingress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Ingress), err
}
