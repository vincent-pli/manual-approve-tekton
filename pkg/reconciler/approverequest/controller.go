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

package approverequest

import (
	"context"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	cache "github.com/patrickmn/go-cache"
	approverequestclient "github.com/vincent-pli/manual-approve-tekton/pkg/client/injection/client"
	approverequestinformer "github.com/vincent-pli/manual-approve-tekton/pkg/client/injection/informers/approverequests/v1alpha1/approverequest"
	requestreconciler "github.com/vincent-pli/manual-approve-tekton/pkg/client/injection/reconciler/approverequests/v1alpha1/approverequest"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)
	approveRequestinformer := approverequestinformer.Get(ctx)
	approverequestclientset := approverequestclient.Get(ctx)

	c := cache.New(0, 0)
	logger.Info("Creating reconciler.")
	r := &Reconciler{
		cache:                   c,
		approverequestClientSet: approverequestclientset,
	}
	impl := requestreconciler.NewImpl(ctx, r, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			SkipStatusUpdates: true,
		}
	})
	r.Tracker = impl.Tracker

	logger.Info("Setting up event handlers.")

	approveRequestinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	//start web server for rest
	webServer := WebServer{
		cache:                   c,
		approverequestClientSet: approverequestclientset,
	}
	go webServer.Start(context.TODO())

	return impl
}
