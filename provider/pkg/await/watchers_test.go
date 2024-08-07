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
	"time"

	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
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
			name: "related warning",
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
			want: "warning [stack/my-stack-28073241] StackUpdateFailure: Failed to update Stack",
		},
		{
			name: "related info",
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
			want: "",
		},
		{
			name: "unrelated warning",
			event: watch.Event{Type: watch.Added, Object: &unstructured.Unstructured{Object: map[string]any{
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
			}}},
			want: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			source := condition.Static(make(chan watch.Event, 1))

			buf := &strings.Builder{}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			obs := NewEventAggregator(context.Background(), source, log{buf}, owner)

			seen := make(chan struct{})
			go obs.Range(func(e watch.Event) bool {
				_ = obs.Observe(e)
				seen <- struct{}{}
				return true
			})

			source <- tt.event
			select {
			case <-seen:
			case <-ctx.Done():
			}
			assert.Equal(t, buf.String(), tt.want)
		})
	}
}

type log struct{ io.Writer }

func (l log) LogMessage(m checkerlog.Message) {
	fmt.Fprintf(l, "%s %s", m.Severity, m.String())
}
