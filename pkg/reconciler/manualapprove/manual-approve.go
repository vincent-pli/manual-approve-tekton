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

package manualapprove

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	clientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	runreconciler "github.com/tektoncd/pipeline/pkg/client/injection/reconciler/pipeline/v1alpha1/run"
	listersalpha "github.com/tektoncd/pipeline/pkg/client/listers/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/reconciler/events"
	approverequestsv1alpha1 "github.com/vincent-pli/manual-approve-tekton/pkg/apis/approverequests/v1alpha1"
	approverequestclientset "github.com/vincent-pli/manual-approve-tekton/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"
)

// Reconciler implements addressableservicereconciler.Interface for
// AddressableService resources.
type Reconciler struct {
	// Tracker builds an index of what resources are watching other resources
	// so that we can immediately react to changes tracked resources.
	Tracker tracker.Interface

	pipelineClientSet clientset.Interface
	// Listers index properties about resources
	runLister listersalpha.RunLister

	approverequestClientSet approverequestclientset.Interface
}

// Check that our Reconciler implements Interface
var _ runreconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, run *v1alpha1.Run) reconciler.Event {
	var merr error
	logger := logging.FromContext(ctx)
	logger.Infof("Reconciling Run %s/%s at %v", run.Namespace, run.Name, time.Now())

	// Check that the Run references a APPROVEREQUEST CRD.  The logic is controller.go should ensure that only this type of Run
	// is reconciled this controller but it never hurts to do some bullet-proofing.
	if run.Spec.Ref == nil ||
		run.Spec.Ref.APIVersion != approverequestsv1alpha1.SchemeGroupVersion.String() ||
		run.Spec.Ref.Kind != "ApproveRequest" {
		logger.Errorf("Received control for a Run %s/%s that does not reference a ApproveRequest CRD", run.Namespace, run.Name)
		return nil
	}

	// If the Run has not started, initialize the Condition and set the start time.
	if !run.HasStarted() {
		logger.Infof("Starting new Run %s/%s", run.Namespace, run.Name)
		run.Status.InitializeConditions()
		// In case node time was not synchronized, when controller has been scheduled to other nodes.
		if run.Status.StartTime.Sub(run.CreationTimestamp.Time) < 0 {
			logger.Warnf("Run %s createTimestamp %s is after the Run started %s", run.Name, run.CreationTimestamp, run.Status.StartTime)
			run.Status.StartTime = &run.CreationTimestamp
		}

		// Emit events. During the first reconcile the status of the Run may change twice
		// from not Started to Started and then to Running, so we need to sent the event here
		// and at the end of 'Reconcile' again.
		// We also want to send the "Started" event as soon as possible for anyone who may be waiting
		// on the event to perform user facing initialisations, such has reset a CI check status
		afterCondition := run.Status.GetCondition(apis.ConditionSucceeded)
		events.Emit(ctx, nil, afterCondition, run)
	}

	if run.IsDone() {
		logger.Infof("Run %s/%s is done", run.Namespace, run.Name)
		return nil
	}

	if run.IsCancelled() {
		logger.Infof("Run %s/%s is cancelled", run.Namespace, run.Name)
		return nil
	}

	// Store the condition before reconcile
	beforeCondition := run.Status.GetCondition(apis.ConditionSucceeded)

	// Reconcile the Run
	if err := r.reconcile(ctx, run); err != nil {
		logger.Errorf("Reconcile error: %v", err.Error())
		merr = multierror.Append(merr, err)
	}

	afterCondition := run.Status.GetCondition(apis.ConditionSucceeded)
	events.Emit(ctx, beforeCondition, afterCondition, run)

	// Only transient errors that should retry the reconcile are returned.
	return merr
}

func (r *Reconciler) reconcile(ctx context.Context, run *v1alpha1.Run) error {
	// logger := logging.FromContext(ctx)

	// Get the ApproveRequest referenced by the Run
	ar, err := r.getApproveRequest(ctx, run)
	if err != nil {
		return nil
	}

	arCopy := ar.DeepCopy()
	//if approveRequest.requestName is not nil
	if arCopy.Spec.RequestName == "" {
		arCopy.Spec.RequestName = run.Name
		arCopy.Status.Approved = false

		run.Status.MarkRunRunning(approverequestsv1alpha1.ApproveRequestRunReasonRunning.String(),
			"There is no taskrun in original pr mark as failed, wait: %s", time.Now().String())

		_, err = r.approverequestClientSet.CustomV1alpha1().ApproveRequests(run.Namespace).Update(ctx, arCopy, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Update ApproveRequest: %s failed: %w", fmt.Sprintf("%s/%s", arCopy.Namespace, arCopy.Name), err)
		}
		return nil
	} else if arCopy.Status.Approved {
		run.Status.MarkRunSucceeded(approverequestsv1alpha1.ApproveRequestRunReasonSucceeded.String(), "The approve request is approved: %s/%s", ar.Name, ar.Namespace)
	}

	return nil
}

func (r *Reconciler) getApproveRequest(ctx context.Context, run *v1alpha1.Run) (*approverequestsv1alpha1.ApproveRequest, error) {
	var approverequest *approverequestsv1alpha1.ApproveRequest

	if run.Spec.Ref != nil && run.Spec.Ref.Name != "" {
		// Use the k8 client to get the TaskLoop rather than the lister.  This avoids a timing issue where
		// the TaskLoop is not yet in the lister cache if it is created at nearly the same time as the Run.
		// See https://github.com/tektoncd/pipeline/issues/2740 for discussion on this issue.
		ar, err := r.approverequestClientSet.CustomV1alpha1().ApproveRequests(run.Namespace).Get(ctx, run.Spec.Ref.Name, metav1.GetOptions{})
		if err != nil {
			run.Status.MarkRunFailed(approverequestsv1alpha1.ApproveRequestRunReasonCouldntGet.String(),
				"Error retrieving TaskLoop for Run %s/%s: %s",
				run.Namespace, run.Name, err)
			return nil, fmt.Errorf("Error retrieving ApproveRequeset for Run %s: %w", fmt.Sprintf("%s/%s", run.Namespace, run.Name), err)
		}

		approverequest = ar
	} else {
		// Run does not require name but here it does.
		run.Status.MarkRunFailed(approverequestsv1alpha1.ApproveRequestRunReasonCouldntGet.String(),
			"Missing spec.ref.name for Run %s/%s",
			run.Namespace, run.Name)
		return nil, fmt.Errorf("Missing spec.ref.name for Run %s", fmt.Sprintf("%s/%s", run.Namespace, run.Name))
	}

	return approverequest, nil
}
