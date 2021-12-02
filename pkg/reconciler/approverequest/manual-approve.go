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

	cache "github.com/patrickmn/go-cache"
	approverequestsv1alpha1 "github.com/vincent-pli/manual-approve-tekton/pkg/apis/approverequests/v1alpha1"
	approverequestclientset "github.com/vincent-pli/manual-approve-tekton/pkg/client/clientset/versioned"
	requestreconciler "github.com/vincent-pli/manual-approve-tekton/pkg/client/injection/reconciler/approverequests/v1alpha1/approverequest"
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

	approverequestClientSet approverequestclientset.Interface
	cache                   *cache.Cache
}

// Check that our Reconciler implements Interface
var _ requestreconciler.Interface = (*Reconciler)(nil)
var _ requestreconciler.Finalizer = (*Reconciler)(nil)

func (r *Reconciler) FinalizeKind(ctx context.Context, ar *approverequestsv1alpha1.ApproveRequest) reconciler.Event {
	logger := logging.FromContext(ctx)

	logger.Infof("start remove the approverequest %s/%s", ar.Namespace, ar.Name)
	key := ar.Namespace + "/" + ar.Name
	_, found := r.cache.Get(key)
	if !found {
		return nil
	}

	r.cache.Delete(key)

	return nil
}

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, ar *approverequestsv1alpha1.ApproveRequest) reconciler.Event {
	logger := logging.FromContext(ctx)

	logger.Infof("start cache the approverequest %s/%s", ar.Namespace, ar.Name)
	key := ar.Namespace + "/" + ar.Name

	r.cache.Set(key, ar.Status.Requests, cache.NoExpiration)
	return nil
}
