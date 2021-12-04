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

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	runinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1alpha1/run"
	runreconciler "github.com/tektoncd/pipeline/pkg/client/injection/reconciler/pipeline/v1alpha1/run"
	pipelinecontroller "github.com/tektoncd/pipeline/pkg/controller"
	"github.com/vincent-pli/manual-approve-tekton/pkg/apis/approverequests/v1alpha1"
	scheme "github.com/vincent-pli/manual-approve-tekton/pkg/client/clientset/versioned/scheme"
	approverequestclient "github.com/vincent-pli/manual-approve-tekton/pkg/client/injection/client"
	approverequestinformer "github.com/vincent-pli/manual-approve-tekton/pkg/client/injection/informers/approverequests/v1alpha1/approverequest"
)

// Reconciler implements addressableservicereconciler.Interface for
// AddressableService resources.
type Enqueue struct {
	impl   controller.Impl
	logger *zap.SugaredLogger
}

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)
	pipelineclientset := pipelineclient.Get(ctx)
	runInformer := runinformer.Get(ctx)
	approveRequestinformer := approverequestinformer.Get(ctx)
	approverequestclientset := approverequestclient.Get(ctx)

	logger.Info("Creating reconciler.")
	r := &Reconciler{
		pipelineClientSet:       pipelineclientset,
		approverequestClientSet: approverequestclientset,
		runLister:               runInformer.Lister(),
	}
	impl := runreconciler.NewImpl(ctx, r)
	r.Tracker = impl.Tracker

	logger.Info("Setting up event handlers.")
	enqueue := &Enqueue{
		impl:   *impl,
		logger: logger,
	}

	approveRequestinformer.Informer().AddEventHandler(controller.HandleAll(enqueue.EnqueueReferenceRun))

	runInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: pipelinecontroller.FilterRunRef(v1alpha1.SchemeGroupVersion.String(), "ApproveRequest"),
		Handler:    controller.HandleAll(impl.Enqueue),
	})

	return impl
}

func (e *Enqueue) EnqueueReferenceRun(obj interface{}) {
	ar, ok := obj.(*v1alpha1.ApproveRequest)
	if !ok {
		e.logger.Error("Not a ApproveRequest")
		return
	}

	// If we can determine the controller ref of this object, then
	// add that object to our workqueue.
	// TODO, problem here, will cause useless reconcile, need enhancement
	for _, request := range ar.Status.Requests {
		fmt.Println("xxxxxxxxxxxxxxxxxx")
		if request.RequestName != "" && request.Approved {
			e.impl.EnqueueKey(types.NamespacedName{Namespace: ar.GetNamespace(), Name: request.RequestName})
		}
	}
}

func addTypeInformationToObject(obj runtime.Object) error {
	gvks, _, err := scheme.Scheme.ObjectKinds(obj)
	if err != nil {
		return fmt.Errorf("missing apiVersion or kind and cannot assign it; %w", err)
	}

	for _, gvk := range gvks {
		if len(gvk.Kind) == 0 {
			continue
		}
		if len(gvk.Version) == 0 || gvk.Version == runtime.APIVersionInternal {
			continue
		}
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		break
	}

	return nil
}
