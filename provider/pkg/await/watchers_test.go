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

	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
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
		name   string
		fields map[string]any
		want   string
	}{
		{
			name: "unknown occurrence time",
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: unknown, count: 1)",
		},
		{
			name: "legacy repeated warning",
			fields: map[string]any{
				"firstTimestamp": "2026-01-01T08:00:00Z",
				"lastTimestamp":  "2026-01-01T09:00:00Z",
				"count":          int64(4),
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: 2026-01-01T09:00:00Z, count: 4)",
		},
		{
			name: "series overrides legacy occurrence",
			fields: map[string]any{
				"lastTimestamp": "2026-01-01T09:00:00Z",
				"count":         int64(4),
				"series": map[string]any{
					"lastObservedTime": "2026-01-01T10:00:00.123456Z",
					"count":            int64(7),
				},
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: 2026-01-01T10:00:00.123456Z, count: 7)",
		},
		{
			name: "series without time uses legacy last timestamp",
			fields: map[string]any{
				"lastTimestamp": "2026-01-01T09:00:00Z",
				"count":         int64(4),
				"series": map[string]any{
					"count": int64(7),
				},
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: 2026-01-01T09:00:00Z, count: 7)",
		},
		{
			name: "modern singleton",
			fields: map[string]any{
				"eventTime": "2026-01-01T10:00:00.123456Z",
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: 2026-01-01T10:00:00.123456Z, count: 1)",
		},
		{
			name: "legacy singleton first timestamp",
			fields: map[string]any{
				"firstTimestamp": "2026-01-01T08:00:00Z",
				"count":          int64(1),
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: 2026-01-01T08:00:00Z, count: 1)",
		},
		{
			name: "repeated warning without latest time",
			fields: map[string]any{
				"eventTime":      "2026-01-01T08:00:00.000000Z",
				"firstTimestamp": "2026-01-01T08:00:00Z",
				"metadata": map[string]any{
					"creationTimestamp": "2026-01-01T08:00:00Z",
				},
				"count": int64(4),
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: unknown, count: 4)",
		},
		{
			name: "series without latest time",
			fields: map[string]any{
				"eventTime":      "2026-01-01T08:00:00.000000Z",
				"firstTimestamp": "2026-01-01T08:00:00Z",
				"count":          int64(1),
				"series": map[string]any{
					"count": int64(7),
				},
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: unknown, count: 7)",
		},
		{
			name: "UTC timestamp",
			fields: map[string]any{
				"lastTimestamp": "2026-01-01T12:00:00+02:00",
				"count":         int64(2),
			},
			want: "warning StackUpdateFailure: Failed to update Stack (last seen: 2026-01-01T10:00:00Z, count: 2)",
		},
		{
			name: "normal",
			fields: map[string]any{
				"type":          "Normal",
				"reason":        "Killing",
				"message":       "frog blast the vent core",
				"lastTimestamp": "2026-01-01T09:00:00Z",
				"count":         int64(4),
			},
			want: "debug Killing: frog blast the vent core",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &strings.Builder{}
			obs := NewEventAggregator(context.Background(), nil, log{buf}, owner)
			obj := &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Event",
				"reason":     "StackUpdateFailure",
				"type":       "Warning",
				"message":    "Failed to update Stack",
			}}
			for field, value := range tt.fields {
				obj.Object[field] = value
			}

			err := obs.Observe(watch.Event{Type: watch.Added, Object: obj})

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestEventAggregatorDeduplicatesOccurrences(t *testing.T) {
	ctx := context.Background()
	engine := fakehost.NewEngineServer(t)
	logger := logging.NewLogger(ctx, &fakehost.HostClient{Engine: engine}, "")
	obs := NewEventAggregator(ctx, nil, logger, &unstructured.Unstructured{})
	obj := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion":    "v1",
		"kind":          "Event",
		"reason":        "FailedUpdate",
		"type":          "Warning",
		"message":       "Could not update volume claim",
		"lastTimestamp": "2026-01-01T09:00:00Z",
		"count":         int64(4),
	}}

	assert.NoError(t, obs.Observe(watch.Event{Type: watch.Added, Object: obj}))
	assert.NoError(t, obs.Observe(watch.Event{Type: watch.Modified, Object: obj}))
	obj.Object["lastTimestamp"] = "2026-01-01T10:00:00Z"
	obj.Object["count"] = int64(5)
	assert.NoError(t, obs.Observe(watch.Event{Type: watch.Modified, Object: obj}))

	var messages []string
	for _, entry := range engine.Logs() {
		messages = append(messages, entry.GetMessage())
	}
	assert.Equal(t, []string{
		"FailedUpdate: Could not update volume claim (last seen: 2026-01-01T09:00:00Z, count: 4)",
		"FailedUpdate: Could not update volume claim (last seen: 2026-01-01T10:00:00Z, count: 5)",
	}, messages)
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
