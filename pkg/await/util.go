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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/dynamic"
)

// --------------------------------------------------------------------------

// Event helpers.
//
// Unlike the vast majority of our client (which does not use concrete types at all), we take a
// concrete dependency on `v1.Event` because it is a fundamental type used to communicate important
// status updates.

// --------------------------------------------------------------------------

func stringifyEvents(events []v1.Event) string {
	var output string
	for _, e := range events {
		output += fmt.Sprintf("\n   * %s (%s): %s: %s",
			e.InvolvedObject.Name, e.InvolvedObject.Kind,
			e.Reason, e.Message)
	}
	return output
}

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

	items := out.(*unstructured.UnstructuredList).Items
	events := []v1.Event{}
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

// Path helpers.

// --------------------------------------------------------------------------

func resourceListEquals(x, y v1.ResourceList) bool {
	for k, v := range x {
		yValue, ok := y[k]
		if !ok {
			return false
		}
		if v.Cmp(yValue) != 0 {
			return false
		}
	}
	for k, v := range y {
		xValue, ok := x[k]
		if !ok {
			return false
		}
		if v.Cmp(xValue) != 0 {
			return false
		}
	}
	return true
}

// --------------------------------------------------------------------------

// Response helpers.

// --------------------------------------------------------------------------

func is404(err error) bool {
	statusErr, ok := err.(*errors.StatusError)
	return ok && statusErr.ErrStatus.Code == 404
}
