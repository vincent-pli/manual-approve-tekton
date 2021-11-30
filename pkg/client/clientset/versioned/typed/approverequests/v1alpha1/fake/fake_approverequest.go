/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/vincent-pli/manual-approve-tekton/pkg/apis/approverequests/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeApproveRequests implements ApproveRequestInterface
type FakeApproveRequests struct {
	Fake *FakeCustomV1alpha1
	ns   string
}

var approverequestsResource = schema.GroupVersionResource{Group: "custom.tektoncd.dev", Version: "v1alpha1", Resource: "approverequests"}

var approverequestsKind = schema.GroupVersionKind{Group: "custom.tektoncd.dev", Version: "v1alpha1", Kind: "ApproveRequest"}

// Get takes name of the approveRequest, and returns the corresponding approveRequest object, and an error if there is any.
func (c *FakeApproveRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ApproveRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(approverequestsResource, c.ns, name), &v1alpha1.ApproveRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApproveRequest), err
}

// List takes label and field selectors, and returns the list of ApproveRequests that match those selectors.
func (c *FakeApproveRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ApproveRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(approverequestsResource, approverequestsKind, c.ns, opts), &v1alpha1.ApproveRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ApproveRequestList{ListMeta: obj.(*v1alpha1.ApproveRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.ApproveRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested approveRequests.
func (c *FakeApproveRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(approverequestsResource, c.ns, opts))

}

// Create takes the representation of a approveRequest and creates it.  Returns the server's representation of the approveRequest, and an error, if there is any.
func (c *FakeApproveRequests) Create(ctx context.Context, approveRequest *v1alpha1.ApproveRequest, opts v1.CreateOptions) (result *v1alpha1.ApproveRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(approverequestsResource, c.ns, approveRequest), &v1alpha1.ApproveRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApproveRequest), err
}

// Update takes the representation of a approveRequest and updates it. Returns the server's representation of the approveRequest, and an error, if there is any.
func (c *FakeApproveRequests) Update(ctx context.Context, approveRequest *v1alpha1.ApproveRequest, opts v1.UpdateOptions) (result *v1alpha1.ApproveRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(approverequestsResource, c.ns, approveRequest), &v1alpha1.ApproveRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApproveRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeApproveRequests) UpdateStatus(ctx context.Context, approveRequest *v1alpha1.ApproveRequest, opts v1.UpdateOptions) (*v1alpha1.ApproveRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(approverequestsResource, "status", c.ns, approveRequest), &v1alpha1.ApproveRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApproveRequest), err
}

// Delete takes name of the approveRequest and deletes it. Returns an error if one occurs.
func (c *FakeApproveRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(approverequestsResource, c.ns, name), &v1alpha1.ApproveRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApproveRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(approverequestsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ApproveRequestList{})
	return err
}

// Patch applies the patch and returns the patched approveRequest.
func (c *FakeApproveRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ApproveRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(approverequestsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ApproveRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApproveRequest), err
}
