// Copyright 2016-2018, Pulumi Corporation.
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
	"fmt"
	"log"
	"sort"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
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

func watchAddedEvent(obj runtime.Object) watch.Event {
	return watch.Event{
		Type:   watch.Added,
		Object: obj,
	}
}

// nolint
func stringifyEvents(events []v1.Event) string {
	var output string
	for _, e := range events {
		output += fmt.Sprintf("\n   * %s (%s): %s: %s",
			e.InvolvedObject.Name, e.InvolvedObject.Kind,
			e.Reason, e.Message)
	}
	return output
}

// nolint
func getLastWarningsForObject(
	clientForEvents dynamic.ResourceInterface, namespace, name, kind string, limit int,
) ([]v1.Event, error) {
	m := map[string]string{
		"involvedObject.name": name,
		"involvedObject.kind": kind,
	}
	if namespace != "" {
		m["involvedObject.namespace"] = namespace
	}

	fs := fields.Set(m).String()
	glog.V(9).Infof("Looking up events via this selector: %q", fs)
	out, err := clientForEvents.List(metav1.ListOptions{
		FieldSelector: fs,
	})
	if err != nil {
		return nil, err
	}

	items := out.Items
	var events []v1.Event
	for _, item := range items {
		// Round trip conversion from `Unstructured` to `v1.Event`. There doesn't seem to be a good way
		// to do this conversion in client-go, and this is not a performance-critical section. When we
		// upgrade to client-go v7, we can replace it with `runtime.DefaultUnstructuredConverter`.
		var warning v1.Event
		str, err := item.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}

		err = warning.Unmarshal(str)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, warning)
	}

	// Bring latest events to the top, for easy access
	sort.Slice(events, func(i, j int) bool {
		return events[i].LastTimestamp.After(events[j].LastTimestamp.Time)
	})

	glog.V(9).Infof("Received '%d' events for %s/%s (%s)",
		len(events), namespace, name, kind)

	// It would be better to sort & filter on the server-side
	// but API doesn't seem to support it
	var warnings []v1.Event
	warnCount := 0
	uniqueWarnings := make(map[string]v1.Event)
	for _, e := range events {
		if warnCount >= limit {
			break
		}

		if e.Type == v1.EventTypeWarning {
			_, found := uniqueWarnings[e.Message]
			if found {
				continue
			}
			warnings = append(warnings, e)
			uniqueWarnings[e.Message] = e
			warnCount++
		}
	}

	return warnings, nil
}

// --------------------------------------------------------------------------

// Version helpers.

// --------------------------------------------------------------------------

// ServerVersion attempts to retrieve the server version from k8s.
// Returns the configured default version in case this fails.
func ServerVersion(cdi discovery.CachedDiscoveryInterface) serverVersion {
	var version serverVersion
	if sv, err := cdi.ServerVersion(); err == nil {
		if v, err := parseVersion(sv); err == nil {
			version = v
		} else {
			version = defaultVersion()
		}
	} else {
		version = defaultVersion()
	}

	return version
}

// --------------------------------------------------------------------------

// Response helpers.

// --------------------------------------------------------------------------

func is404(err error) bool {
	statusErr, ok := err.(*errors.StatusError)
	return ok && statusErr.ErrStatus.Code == 404
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
	}

	owners := obj.GetOwnerReferences()
	for _, owner := range owners {
		var ownerAPIVersion string
		switch owner.Kind {
		case "Deployment":
			ownerAPIVersion = canonicalizeDeploymentAPIVersion(owner.APIVersion)
		case "StatefulSet":
			ownerAPIVersion = canonicalizeStatefulSetAPIVersion(owner.APIVersion)
		}

		if ownerAPIVersion == possibleOwnerAPIVersion &&
			possibleOwner.GetKind() == owner.Kind && possibleOwner.GetName() == owner.Name {
			return true
		}
	}

	return false
}
