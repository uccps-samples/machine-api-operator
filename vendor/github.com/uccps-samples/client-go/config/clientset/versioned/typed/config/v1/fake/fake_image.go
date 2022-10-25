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

// FakeImages implements ImageInterface
type FakeImages struct {
	Fake *FakeConfigV1
}

var imagesResource = schema.GroupVersionResource{Group: "config.uccp.io", Version: "v1", Resource: "images"}

var imagesKind = schema.GroupVersionKind{Group: "config.uccp.io", Version: "v1", Kind: "Image"}

// Get takes name of the image, and returns the corresponding image object, and an error if there is any.
func (c *FakeImages) Get(ctx context.Context, name string, options v1.GetOptions) (result *configv1.Image, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(imagesResource, name), &configv1.Image{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Image), err
}

// List takes label and field selectors, and returns the list of Images that match those selectors.
func (c *FakeImages) List(ctx context.Context, opts v1.ListOptions) (result *configv1.ImageList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(imagesResource, imagesKind, opts), &configv1.ImageList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &configv1.ImageList{ListMeta: obj.(*configv1.ImageList).ListMeta}
	for _, item := range obj.(*configv1.ImageList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested images.
func (c *FakeImages) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(imagesResource, opts))
}

// Create takes the representation of a image and creates it.  Returns the server's representation of the image, and an error, if there is any.
func (c *FakeImages) Create(ctx context.Context, image *configv1.Image, opts v1.CreateOptions) (result *configv1.Image, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(imagesResource, image), &configv1.Image{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Image), err
}

// Update takes the representation of a image and updates it. Returns the server's representation of the image, and an error, if there is any.
func (c *FakeImages) Update(ctx context.Context, image *configv1.Image, opts v1.UpdateOptions) (result *configv1.Image, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(imagesResource, image), &configv1.Image{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Image), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeImages) UpdateStatus(ctx context.Context, image *configv1.Image, opts v1.UpdateOptions) (*configv1.Image, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(imagesResource, "status", image), &configv1.Image{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Image), err
}

// Delete takes name of the image and deletes it. Returns an error if one occurs.
func (c *FakeImages) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(imagesResource, name), &configv1.Image{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeImages) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(imagesResource, listOpts)

	_, err := c.Fake.Invokes(action, &configv1.ImageList{})
	return err
}

// Patch applies the patch and returns the patched image.
func (c *FakeImages) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configv1.Image, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(imagesResource, name, pt, data, subresources...), &configv1.Image{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configv1.Image), err
}
