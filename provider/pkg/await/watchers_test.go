// Copyright 2024, Pulumi Corporation.
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
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

func TestEventAggregator(t *testing.T) {
	owner := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "pulumi.com/v1",
		"kind":       "Stack",
		"metadata": map[string]any{
			"name":      "my-stack-28073241",
			"namespace": "operator",
			"uid":       "9bd08f1a-fa5b-40a5-ba41-bf69899a4416",
		},
	}}

	tests := []struct {
		name  string
		event watch.Event
		want  string
	}{
		{
			name: "warning",
			event: watch.Event{Type: watch.Added, Object: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Event",
				"reason":     "StackUpdateFailure",
				"type":       "Warning",
				"message":    "Failed to update Stack",
				"involvedObject": map[string]any{
					"kind": "Stack",
					"name": "my-stack-28073241",
					"uid":  "9bd08f1a-fa5b-40a5-ba41-bf69899a4416",
				},
			}}},
			want: "warning StackUpdateFailure: Failed to update Stack",
		},
		{
			name: "normal",
			event: watch.Event{Type: watch.Added, Object: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Event",
				"reason":     "Killing",
				"type":       "Normal",
				"message":    "frog blast the vent core",
				"involvedObject": map[string]any{
					"kind": "Pod",
					"name": "mypod-7854ff8877-p9ksk",
					"uid":  "9bd08f1a-fa5b-40a5-ba41-bf69899a4416",
				},
			}}},
			want: "debug Killing: frog blast the vent core",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &strings.Builder{}
			obs := NewEventAggregator(context.Background(), nil, log{buf}, owner)

			err := obs.Observe(tt.event)

			assert.NoError(t, err)
			assert.Equal(t, buf.String(), tt.want)
		})
	}
}

func TestRelatedEvent(t *testing.T) {
	owner := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "pulumi.com/v1",
		"kind":       "Stack",
		"metadata": map[string]any{
			"name":      "my-stack-28073241",
			"namespace": "operator",
			"uid":       "9bd08f1a-fa5b-40a5-ba41-bf69899a4416",
		},
	}}

	tests := []struct {
		name  string
		owner *unstructured.Unstructured
		event *unstructured.Unstructured
		want  bool
	}{
		{
			name:  "related",
			owner: owner,
			event: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Event",
				"reason":     "StackUpdateFailure",
				"type":       "Warning",
				"message":    "Failed to update Stack",
				"involvedObject": map[string]any{
					"kind": "Stack",
					"name": "my-stack-28073241",
					"uid":  "9bd08f1a-fa5b-40a5-ba41-bf69899a4416",
				},
			}},
			want: true,
		},
		{
			name:  "unrelated",
			owner: owner,
			event: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Event",
				"reason":     "Killing",
				"type":       "Warning",
				"message":    "Stopping container nginx",
				"involvedObject": map[string]any{
					"kind": "Pod",
					"name": "some-other-name",
					"uid":  "some-other-uid",
				},
			}},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := relatedEvents(tt.owner)(tt.event)
			assert.Equal(t, tt.want, actual)
		})
	}
}

type log struct{ io.Writer }

func (l log) Log(sev diag.Severity, msg string) {
	fmt.Fprintf(l, "%s %s", sev, msg)
}

func (l log) LogStatus(sev diag.Severity, msg string) {
	l.Log(sev, msg)
}
