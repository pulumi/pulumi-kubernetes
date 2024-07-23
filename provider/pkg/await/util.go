// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package await

import (
	"errors"
	"net/http"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

const trueStatus = "True"

// --------------------------------------------------------------------------

// Event helpers.
//
// Unlike the vast majority of our client (which does not use concrete types at all), we take a
// concrete dependency on `v1.Event` because it is a fundamental type used to communicate important
// status updates.

// --------------------------------------------------------------------------

// chanWatcher is a `watch.Interface` implementation meant to make it easy to mock the Kubernetes
// client-go's resource watchers. This is useful any place we'd want to provide a series of updates
// to resources from a source other than the Kubernetes API server. For example, for testing we
// might want to mock the API server, providing synthetic updates to some resource.
type chanWatcher struct {
	results chan watch.Event
}

var _ watch.Interface = (*chanWatcher)(nil)

func (mw *chanWatcher) Stop() {}

func (mw *chanWatcher) ResultChan() <-chan watch.Event {
	return mw.results
}

// watchAddedEvent returns a watch.Event of type Added.
func watchAddedEvent(obj runtime.Object) watch.Event {
	return watch.Event{
		Type:   watch.Added,
		Object: obj,
	}
}

// watchAddedEvent returns a watch.Event of type Modified.
func watchModifiedEvent(obj runtime.Object) watch.Event {
	return watch.Event{
		Type:   watch.Modified,
		Object: obj,
	}
}

// watchAddedEvent returns a watch.Event of type Deleted.
func watchDeletedEvent(obj runtime.Object) watch.Event {
	return watch.Event{
		Type:   watch.Deleted,
		Object: obj,
	}
}

// --------------------------------------------------------------------------

// Response helpers.

// --------------------------------------------------------------------------

// is404 returns true if any error in err's tree is a Kubernetes StatusError
// with code 404.
func is404(err error) bool {
	var statusErr *k8serrors.StatusError
	if errors.As(err, &statusErr) {
		return statusErr.ErrStatus.Code == http.StatusNotFound
	}
	return false
}

// --------------------------------------------------------------------------

// Ownership helpers.

// --------------------------------------------------------------------------

// TODO: Remove in favor of PodAggregator.
func isOwnedBy(obj, possibleOwner *unstructured.Unstructured) bool {
	if possibleOwner == nil {
		return false
	}

	var possibleOwnerAPIVersion string

	// Canonicalize apiVersion.
	switch possibleOwner.GetKind() {
	case "Deployment":
		possibleOwnerAPIVersion = canonicalizeDeploymentAPIVersion(possibleOwner.GetAPIVersion())
	case "StatefulSet":
		possibleOwnerAPIVersion = canonicalizeStatefulSetAPIVersion(possibleOwner.GetAPIVersion())
	default:
		possibleOwnerAPIVersion = possibleOwner.GetAPIVersion()
	}

	owners := obj.GetOwnerReferences()
	for _, owner := range owners {
		var ownerAPIVersion string
		switch owner.Kind {
		case "Deployment":
			ownerAPIVersion = canonicalizeDeploymentAPIVersion(owner.APIVersion)
		case "StatefulSet":
			ownerAPIVersion = canonicalizeStatefulSetAPIVersion(owner.APIVersion)
		default:
			ownerAPIVersion = owner.APIVersion
		}

		if ownerAPIVersion == possibleOwnerAPIVersion &&
			possibleOwner.GetKind() == owner.Kind &&
			possibleOwner.GetName() == owner.Name {
			return true
		}
	}

	return false
}
