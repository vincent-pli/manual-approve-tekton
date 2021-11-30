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

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	runinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1alpha1/run"
	runreconciler "github.com/tektoncd/pipeline/pkg/client/injection/reconciler/pipeline/v1alpha1/run"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)
	pipelineclientset := pipelineclient.Get(ctx)
	runInformer := runinformer.Get(ctx)

	logger.Info("Creating reconciler.")
	r := &Reconciler{
		pipelineClientSet: pipelineclientset,
		runLister:         runInformer.Lister(),
	}
	impl := runreconciler.NewImpl(ctx, r)
	r.Tracker = impl.Tracker

	logger.Info("Setting up event handlers.")
	runInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
