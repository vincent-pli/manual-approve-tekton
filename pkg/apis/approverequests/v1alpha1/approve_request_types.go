/*
Copyright 2019 The Knative Authors

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// JobRunReason represents a reason for the Run "Succeeded" condition
type RunReason string

func (e RunReason) String() string {
	return string(e)
}

// const (
// 	// ApproveRequestConditionReady is set when the revision is starting to materialize
// 	// runtime resources, and becomes true when those resources are ready.
// 	// ApproveRequestConditionReady = apis.ConditionReady
// 	// ExceptionRunReasonInternalError indicates that the Exception failed due to an internal error in the reconciler
// 	ExceptionRunReasonInternalError RunReason = "ExceptionInternalError"
// 	// ExceptionRunReasonCouldntCancel indicates that the Exception failed due to an internal error in the reconciler
// 	ExceptionRunReasonCouldntCancel RunReason = "CouldntCancel"
// 	// ExceptionRunReasonCouldntGet indicates that the associated Exception couldn't be retrieved
// 	// ApproveRequestRunReasonCouldntGet RunReason = "CouldntGet"
// 	// ExceptionRunReasonCouldntGetOriginalPipelinerun indicates that the associated Exception couldn't be retrieved
// 	ExceptionRunReasonCouldntGetOriginalPipelinerun RunReason = "CouldntGetOriginalPipelinerun"
// 	// ExceptionRunReasonNoError indicates that the original Pipelinerun has no error
// 	ExceptionRunReasonNoError RunReason = "NoError"
// 	// ExceptionRunReasonCoundntCreate indicates that could not create new pipelinerun
// 	ExceptionRunReasonCoundntCreate RunReason = "CoundntCreate"
// )

// ApproveRequest is a Knative abstraction that encapsulates the interface by which Knative
// components express a desire to have a particular image cached.
//
// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApproveRequest struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the ApproveRequest (from the client).
	// +optional
	Spec ApproveRequestSpec `json:"spec,omitempty"`

	// Status communicates the observed state of the ApproveRequest (from the controller).
	// +optional
	Status ApproveRequestStatus `json:"status,omitempty"`
}

var (
	// Check that ApproveRequest can be validated and defaulted.
	_ apis.Validatable   = (*ApproveRequest)(nil)
	_ apis.Defaultable   = (*ApproveRequest)(nil)
	_ kmeta.OwnerRefable = (*ApproveRequest)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*ApproveRequest)(nil)
)

// ApproveRequestSpec holds the desired state of the ApproveRequest (from the client).
type ApproveRequestSpec struct {
	// RequstName holds the name of the Kubernetes Service to expose as an "addressable".
	Approver string `json:"approver,omitempty"`
}

const (
	// ApproveRequestConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	ApproveRequestConditionReady = apis.ConditionReady
	// ApproveRequestRunReasonCouldntGet indicates that the associated Exception couldn't be retrieved
	ApproveRequestRunReasonCouldntGet RunReason = "CouldntGet"
	// ApproveRequestRunReasonRunning indicates that the new created pipelinerun is running
	ApproveRequestRunReasonRunning RunReason = "Running"
	// ApproveRequestRunReasonSucceeded indicates that created Pipelinerun success or no error in original Pipelinerun
	ApproveRequestRunReasonSucceeded RunReason = "Succeeded"
)

type Request struct {
	// RequstName holds the name of the Kubernetes Service to expose as an "addressable".
	RequestName      string      `json:"requstName,omitempty"`
	RequestTimestamp metav1.Time `json:"requestTimestamp,omitempty" protobuf:"bytes,8,opt,name=requestTimestamp"`
	// Approved shows if the request has been approved or not
	Approved         bool        `json:"approved,omitempty"`
	ApproveTimestamp metav1.Time `json:"approveTimestamp,omitempty" protobuf:"bytes,8,opt,name=approveTimestamp"`
}

// ApproveRequestStatus communicates the observed state of the ApproveRequest (from the controller).
type ApproveRequestStatus struct {
	duckv1.Status `json:",inline,omitempty"`
	// RequstName holds the name of the Kubernetes Service to expose as an "addressable".
	Requests []Request `json:"requests"`
}

// ApproveRequestList is a list of ApproveRequest resources
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApproveRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ApproveRequest `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (as *ApproveRequest) GetStatus() *duckv1.Status {
	return &as.Status.Status
}
