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
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	testcore "k8s.io/client-go/testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
)

// _dsGVR is a GroupVersionResource for apps/v1 DaemonSets.
var _dsGVR = v1.SchemeGroupVersion.WithResource("daemonsets")

func TestAwaitDaemonSetCreation(t *testing.T) {
	tests := []struct {
		name   string
		given  *unstructured.Unstructured
		events func(*clockwork.FakeClock, *unstructured.Unstructured) <-chan watch.Event

		want    v1.DaemonSetStatus
		wantErr string
	}{
		{
			name:    "RollingUpdate timed out",
			given:   dsWithRollingUpdate(),
			want:    dsStatusProgressing(),
			events:  dsCreateEventsWithTimeout,
			wantErr: "timed out waiting for the condition",
		},
		{
			name:   "RollingUpdate successful",
			given:  dsWithRollingUpdate(),
			want:   dsStatusRunning(),
			events: dsCreateEventsWithoutTimeout,
		},
		{
			name:    "OnDelete timed out",
			given:   dsWithOnDelete(),
			want:    dsStatusProgressing(),
			events:  dsCreateEventsWithTimeout,
			wantErr: "timed out waiting for the condition",
		},
		{
			name:   "OnDelete successful",
			given:  dsWithOnDelete(),
			want:   dsStatusRunning(),
			events: dsCreateEventsWithoutTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pconfig, clientset, clock := fakeProviderConfig(context.Background(), t)
			config := CreateConfig{
				ProviderConfig: pconfig,
				Inputs:         tt.given,
			}

			// Register a watcher to inject our test events.
			w := watch.NewRaceFreeFake()
			clientset.PrependWatchReactor("daemonsets", testcore.DefaultWatchReactor(w, nil))
			go func() {
				_ = clock.BlockUntilContext(context.Background(), 1) // Timeout sleeper
				for e := range tt.events(clock, tt.given) {
					w.Action(e.Type, e.Object)
					// The Fake's ObjectStore doesn't stay in sync with watch
					// events. We manually update it so subsequent Get
					// operations will correctly reflect the "server's" state.
					_ = clientset.Tracker().Update(_dsGVR, e.Object, "default")
				}
			}()

			ds, err := Creation(config)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assertDaemonSetStatus(t, tt.want, ds)
		})
	}
}

func TestAwaitDaemonSetUpdate(t *testing.T) {
	tests := []struct {
		name   string
		given  *unstructured.Unstructured
		events func(*clockwork.FakeClock, *unstructured.Unstructured) <-chan watch.Event

		want    v1.DaemonSetStatus
		wantErr string
	}{
		{
			name:    "OnDelete timed out",
			given:   dsWithOnDelete(),
			want:    dsStatusProgressing(),
			events:  dsUpdateEventsWithTimeout,
			wantErr: "timed out waiting for the condition",
		},
		{
			name:   "OnDelete successful",
			given:  dsWithOnDelete(),
			events: dsUpdateEventsWithoutTimeout,
			want:   dsStatusRunning(),
		},
		{
			name:    "RollingUpdate timed out",
			given:   dsWithRollingUpdate(),
			want:    dsStatusProgressing(),
			events:  dsUpdateEventsWithTimeout,
			wantErr: "timed out waiting for the condition",
		},
		{
			name:   "RollingUpdate successful",
			given:  dsWithRollingUpdate(),
			want:   dsStatusRunning(),
			events: dsUpdateEventsWithoutTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pconfig, clientset, clock := fakeProviderConfig(context.Background(), t)
			config := UpdateConfig{
				ProviderConfig: pconfig,
				OldInputs:      tt.given,
				OldOutputs:     tt.given,
				Inputs:         tt.given,
			}

			w := watch.NewRaceFreeFake()
			clientset.PrependWatchReactor("daemonsets", testcore.DefaultWatchReactor(w, nil))
			go func() {
				_ = clock.BlockUntilContext(context.Background(), 1) // Timeout sleeper
				for e := range tt.events(clock, tt.given) {
					w.Action(e.Type, e.Object)
					_ = clientset.Tracker().Update(_dsGVR, e.Object, "default")
				}
			}()

			// Update expects the resource to already exist.
			err := clientset.Tracker().Create(_dsGVR, tt.given, "default")
			require.NoError(t, err)

			ds, err := Update(config)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assertDaemonSetStatus(t, tt.want, ds)
		})
	}
}

func TestAwaitDaemonSetRead(t *testing.T) {
	tests := []struct {
		name   string
		given  *unstructured.Unstructured
		events func(*clockwork.FakeClock, *unstructured.Unstructured) <-chan watch.Event

		want    v1.DaemonSetStatus
		wantErr string
	}{
		{
			name:    "OnDelete still progressing",
			given:   dsWithStatus(dsWithOnDelete(), dsStatusProgressing()),
			want:    dsStatusProgressing(),
			events:  dsUpdateEventsWithoutTimeout,
			wantErr: "Resource 'foo' was created but failed to initialize",
		},
		{
			name:   "OnDelete successful",
			given:  dsWithStatus(dsWithOnDelete(), dsStatusRunning()),
			events: dsUpdateEventsWithoutTimeout,
			want:   dsStatusRunning(),
		},
		{
			name:    "RollingUpdate still progressing",
			given:   dsWithRollingUpdate(),
			want:    dsStatusProgressing(),
			events:  dsUpdateEventsWithoutTimeout,
			wantErr: "Resource 'foo' was created but failed to initialize",
		},
		{
			name:   "RollingUpdate successful",
			given:  dsWithStatus(dsWithRollingUpdate(), dsStatusRunning()),
			want:   dsStatusRunning(),
			events: dsUpdateEventsWithoutTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pconfig, clientset, _ := fakeProviderConfig(context.Background(), t)
			config := ReadConfig{
				ProviderConfig: pconfig,
				Inputs:         tt.given,
				Name:           tt.given.GetName(),
			}

			err := clientset.Tracker().Create(_dsGVR, tt.given, tt.given.GetNamespace())
			require.NoError(t, err)

			ds, err := Read(config)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assertDaemonSetStatus(t, tt.want, ds)
		})
	}
}

func TestAwaitDaemonSetDelete(t *testing.T) {
	ensureExists := func(clientset *fake.SimpleDynamicClient, ds *unstructured.Unstructured) {
		err := clientset.Tracker().Create(_dsGVR, ds, ds.GetNamespace())
		require.NoError(t, err)
	}

	dontDeleteImmediately := func(clientset *fake.SimpleDynamicClient, ds *unstructured.Unstructured) {
		// Intercept and ignore delete events -- this way we remove objects
		// from our fake's tracker only after we receive a deletion event from
		// the server, not when the client sends the initial delete request.
		clientset.PrependReactor(
			"delete",
			"daemonsets",
			func(_ testcore.Action) (handled bool, ret runtime.Object, err error) {
				return true, ds, nil
			},
		)
	}

	tests := []struct {
		name    string
		given   *unstructured.Unstructured
		setup   []func(*fake.SimpleDynamicClient, *unstructured.Unstructured)
		events  func(*unstructured.Unstructured) <-chan watch.Event
		timeout time.Duration

		wantErr string
	}{
		{
			name:   "Deletion succeeds immediately",
			given:  dsWithRollingUpdate(),
			setup:  []func(*fake.SimpleDynamicClient, *unstructured.Unstructured){ensureExists},
			events: nil,
		},
		{
			name:  "Deletion succeeds eventually",
			given: dsWithRollingUpdate(),
			setup: []func(*fake.SimpleDynamicClient, *unstructured.Unstructured){
				ensureExists,
				dontDeleteImmediately,
			},
			events: func(ds *unstructured.Unstructured) <-chan watch.Event {
				c := make(chan watch.Event, 1)
				c <- watchDeletedEvent(ds)
				return c
			},
		},
		{
			name:  "Deletion timed out",
			given: dsWithRollingUpdate(),
			setup: []func(*fake.SimpleDynamicClient, *unstructured.Unstructured){
				ensureExists,
				dontDeleteImmediately,
			},
			events: func(_ *unstructured.Unstructured) <-chan watch.Event {
				c := make(chan watch.Event, 1)
				return c
			},
			timeout: 1 * time.Second,
			wantErr: "timed out waiting for the condition",
		},
		{
			name:   "DaemonSet was already deleted",
			given:  dsWithRollingUpdate(),
			setup:  nil, // Initial GET will 404
			events: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pconfig, clientset, _ := fakeProviderConfig(context.Background(), t)
			config := DeleteConfig{
				ProviderConfig: pconfig,
				Inputs:         tt.given,
				Outputs:        tt.given,
				Name:           tt.given.GetName(),
				Timeout:        tt.timeout.Seconds(),
			}

			w := watch.NewRaceFreeFake()
			clientset.PrependWatchReactor("daemonsets", testcore.DefaultWatchReactor(w, nil))

			go func() {
				if tt.events == nil {
					return
				}
				for e := range tt.events(tt.given) {
					w.Action(e.Type, e.Object)
				}
			}()

			for _, s := range tt.setup {
				s(clientset, tt.given)
			}

			err := Deletion(config)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

// fakeProviderConfig returns helpers appropriate for testing await logic.
func fakeProviderConfig(
	ctx context.Context,
	t *testing.T,
) (ProviderConfig, *fake.SimpleDynamicClient, *clockwork.FakeClock) {
	clock := clockwork.NewFakeClock()

	host := &fakehost.HostClient{}

	client, disco, _, clientset := fake.NewSimpleDynamicClient()

	resources, err := openapi.GetResourceSchemasForClient(disco)
	require.NoError(t, err)

	urn := resource.NewURN(
		tokens.QName("teststack"),
		tokens.PackageName("testproj"),
		tokens.Type(""),
		tokens.Type("kubernetes:apps/v1:DaemonSet"),
		"testresource",
	)

	config := ProviderConfig{
		Context:           ctx,
		Host:              host,
		URN:               urn,
		InitialAPIVersion: corev1.SchemeGroupVersion.String(),
		FieldManager:      "test",
		ClusterVersion:    &cluster.ServerVersion{Major: 1, Minor: 29},
		ClientSet:         client,
		DedupLogger:       logging.NewLogger(ctx, host, urn),
		Resources:         resources,
		clock:             clock,
		Factories:         informers.NewFactories(t.Context()),
	}

	return config, clientset, clock
}

// assertDaemonSetStatus asserts a DaemonSet has the expected status.
func assertDaemonSetStatus(
	t *testing.T,
	status v1.DaemonSetStatus,
	uns *unstructured.Unstructured,
) {
	meta, err := clients.FromUnstructured(uns)
	require.NoError(t, err)
	ds, ok := meta.(*v1.DaemonSet)
	require.True(t, ok)

	assert.Equal(t, status, ds.Status)
}

// dsCreateEventsWithoutTimeout simulates a brand new DaemonSet which is
// created without pods and then rolled out within 2 minutes.
func dsCreateEventsWithoutTimeout(
	clock *clockwork.FakeClock,
	ds *unstructured.Unstructured,
) <-chan watch.Event {
	events := make(chan watch.Event, 1)
	go func() {
		events <- watchModifiedEvent(dsWithStatus(ds, dsStatusCreating()))
		clock.Advance(1 * time.Minute)
		for e := range dsUpdateEventsWithoutTimeout(clock, ds) {
			events <- e
		}
	}()
	return events
}

// dsCreateEventsWithTimeout simulates a brand new DaemonSet which is
// created without pods and then takes an hour to roll out.
func dsCreateEventsWithTimeout(
	clock *clockwork.FakeClock,
	ds *unstructured.Unstructured,
) <-chan watch.Event {
	events := make(chan watch.Event, 1)
	go func() {
		events <- watchModifiedEvent(dsWithStatus(ds, dsStatusCreating()))
		clock.Advance(1 * time.Minute)
		for e := range dsUpdateEventsWithTimeout(clock, ds) {
			events <- e
		}
	}()
	return events
}

// dsUpdateEventsWithoutTimeout simulates a DaemonSet rollout which takes one
// minute to complete.
func dsUpdateEventsWithoutTimeout(
	clock *clockwork.FakeClock,
	ds *unstructured.Unstructured,
) <-chan watch.Event {
	events := make(chan watch.Event, 1)
	go func() {
		events <- watchModifiedEvent(dsWithStatus(ds, dsStatusPending()))
		clock.Advance(1 * time.Minute)
		events <- watchModifiedEvent(dsWithStatus(ds, dsStatusProgressing()))
		events <- watchAddedEvent(dsReadyPod(ds.GetNamespace(), "ready-pod", ds.GetName()))
		clock.Advance(1 * time.Minute)
		events <- watchModifiedEvent(dsWithStatus(ds, dsStatusRunning()))
	}()
	return events
}

// dsUpdateEventsWithTimeout simulates a DaemonSet rollout which takes one hour
// to complete.
func dsUpdateEventsWithTimeout(
	clock *clockwork.FakeClock,
	ds *unstructured.Unstructured,
) <-chan watch.Event {
	events := make(chan watch.Event, 1)
	go func() {
		events <- watchModifiedEvent(dsWithStatus(ds, dsStatusProgressing()))
		events <- watchModifiedEvent(dsFailedPod(ds.GetNamespace(), "failed", ds.GetName()))
		clock.Advance(2 * _defaultDaemonSetTimeout)
	}()
	return events
}

// dsWithStatus returns a copy of the DaemonSet with the given status.
func dsWithStatus(
	uns *unstructured.Unstructured,
	status v1.DaemonSetStatus,
) *unstructured.Unstructured {
	obj, _ := clients.FromUnstructured(uns)
	ds := obj.(*v1.DaemonSet)
	ds.Status = status
	res, _ := clients.ToUnstructured(ds)
	return res
}

// dsWithRollingUpdate returns a DaemonSet with the default RollingUpdate strategy.
func dsWithRollingUpdate() *unstructured.Unstructured {
	return mustDecodeUnstructured(`
	{
		"apiVersion": "apps/v1",
		"kind":       "DaemonSet",
		"metadata": {
			"name": "foo",
			"namespace": "default",
			"generation": 1
		},
		"spec": {
			"minReadySeconds":      300,
			"revisionHistoryLimit": 1,
			"updateStrategy": {
				"type": "RollingUpdate",
				"rollingUpdate": {
					"maxUnavailable": 1
				}
			},
			"template": {
				"spec": {
					"containers": [
						{
							"name":  "foo",
							"image": "nginx"
						}
					]
				}
			}
		},
		"status": {
			"currentNumberScheduled": 0,
			"desiredNumberScheduled": 1,
			"numberMisscheduled":     0,
			"numberReady":            0,
			"observedGeneration":     1
		}
	}
	`)
}

// dsWithOnDelete returns a DaemonSet with OnDelete rollout strategy. This
// requires the user to manually delete old pods in order to finish the
// rollout.
func dsWithOnDelete() *unstructured.Unstructured {
	return mustDecodeUnstructured(`
	{
		"apiVersion": "apps/v1",
		"kind":       "DaemonSet",
		"metadata": {
			"name": "foo",
			"namespace": "default",
			"generation": 1
		},
		"spec": {
			"minReadySeconds":      300,
			"revisionHistoryLimit": 1,
			"updateStrategy": {
				"type": "OnDelete"
			},
			"template": {
				"spec": {
					"containers": [
						{
							"name":  "foo",
							"image": "nginx"
						}
					]
				}
			}
		},
		"status": {
			"currentNumberScheduled": 0,
			"desiredNumberScheduled": 1,
			"numberMisscheduled":     0,
			"numberReady":            0,
			"observedGeneration":     1
		}
	}
	`)
}

// dsStatusCreating represents a DaemonSet which has just been created and
// which has no pods currently running.
func dsStatusCreating() v1.DaemonSetStatus {
	return v1.DaemonSetStatus{
		DesiredNumberScheduled: 2,
		CurrentNumberScheduled: 0,
		UpdatedNumberScheduled: 0,
		NumberAvailable:        0,
		NumberMisscheduled:     0,
		ObservedGeneration:     1,
	}
}

// dsStatusPending represents a DaemonSet which is already deployed but which
// has an update pending, so older pods are considered misscheduled.
func dsStatusPending() v1.DaemonSetStatus {
	return v1.DaemonSetStatus{
		DesiredNumberScheduled: 2,
		CurrentNumberScheduled: 2,
		UpdatedNumberScheduled: 0,
		NumberReady:            2,
		NumberAvailable:        2,
		NumberMisscheduled:     2,
		ObservedGeneration:     1,
	}
}

// dsStatusProgressing represents a DaemonSet which is partway through an
// update.
func dsStatusProgressing() v1.DaemonSetStatus {
	return v1.DaemonSetStatus{
		DesiredNumberScheduled: 2,
		CurrentNumberScheduled: 2,
		UpdatedNumberScheduled: 1,
		NumberReady:            1,
		NumberAvailable:        1,
		NumberMisscheduled:     1,
		ObservedGeneration:     1,
	}
}

// dsStatusRunning represents a DaemonSet which has been fully rolled out.
func dsStatusRunning() v1.DaemonSetStatus {
	return v1.DaemonSetStatus{
		DesiredNumberScheduled: 2,
		CurrentNumberScheduled: 2,
		UpdatedNumberScheduled: 2,
		NumberReady:            2,
		NumberAvailable:        2,
		NumberMisscheduled:     0,
		ObservedGeneration:     1,
	}
}

// dsReadyPod returns a running Pod owned by the given DaemonSet.
func dsReadyPod(namespace, name, dsName string) *unstructured.Unstructured {
	return mustDecodeUnstructured(fmt.Sprintf(`
	{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
			"creationTimestamp": "2024-04-12T18:04:44Z",
			"generateName": "%s-",
			"labels": {
				"controller-revision-hash": "675b9b5d7f",
				"pod-template-generation": "1"
			},
			"name": "%s",
			"namespace": "%s",
			"ownerReferences": [
				{
					"apiVersion": "apps/v1",
					"blockOwnerDeletion": true,
					"controller": true,
					"kind": "DaemonSet",
					"name": "%s",
					"uid": "d21151ed-7cf4-4da3-86b3-095c40d40c33"
				}
			],
			"resourceVersion": "50154",
			"uid": "ca2a65f7-ddc3-480e-b6ca-ed9da544b85e"
		},
		"spec": {
			"affinity": {
				"nodeAffinity": {
					"requiredDuringSchedulingIgnoredDuringExecution": {
						"nodeSelectorTerms": [
							{
								"matchFields": [
									{
										"key": "metadata.name",
										"operator": "In",
										"values": [
											"orbstack"
										]
									}
								]
							}
						]
					}
				}
			},
			"containers": [
				{
					"image": "nginx:stable-alpine3.17-slim",
					"imagePullPolicy": "IfNotPresent",
					"name": "nginx",
					"resources": {},
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"volumeMounts": [
						{
							"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
							"name": "kube-api-access-6hrqm",
							"readOnly": true
						}
					]
				}
			],
			"dnsPolicy": "ClusterFirst",
			"enableServiceLinks": true,
			"nodeName": "orbstack",
			"preemptionPolicy": "PreemptLowerPriority",
			"priority": 0,
			"restartPolicy": "Always",
			"schedulerName": "default-scheduler",
			"securityContext": {},
			"serviceAccount": "default",
			"serviceAccountName": "default",
			"terminationGracePeriodSeconds": 30,
			"tolerations": [
				{
					"effect": "NoExecute",
					"key": "node.kubernetes.io/not-ready",
					"operator": "Exists"
				},
				{
					"effect": "NoExecute",
					"key": "node.kubernetes.io/unreachable",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/disk-pressure",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/memory-pressure",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/pid-pressure",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/unschedulable",
					"operator": "Exists"
				}
			],
			"volumes": [
				{
					"name": "kube-api-access-6hrqm",
					"projected": {
						"defaultMode": 420,
						"sources": [
							{
								"serviceAccountToken": {
									"expirationSeconds": 3607,
									"path": "token"
								}
							},
							{
								"configMap": {
									"items": [
										{
											"key": "ca.crt",
											"path": "ca.crt"
										}
									],
									"name": "kube-root-ca.crt"
								}
							},
							{
								"downwardAPI": {
									"items": [
										{
											"fieldRef": {
												"apiVersion": "v1",
												"fieldPath": "metadata.namespace"
											},
											"path": "namespace"
										}
									]
								}
							}
						]
					}
				}
			]
		},
		"status": {
			"conditions": [
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:04:44Z",
					"status": "True",
					"type": "Initialized"
				},
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:04:46Z",
					"status": "True",
					"type": "Ready"
				},
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:04:46Z",
					"status": "True",
					"type": "ContainersReady"
				},
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:04:44Z",
					"status": "True",
					"type": "PodScheduled"
				}
			],
			"containerStatuses": [
				{
					"containerID": "docker://cd48fdb9811c00ef37c7b6bcc0d54cbe01f188987cafe121395f64661f3bf5f1",
					"image": "nginx:stable-alpine-slim",
					"imageID": "docker-pullable://nginx@sha256:0a8c5686d40beca3cf231e223668cf77c91344d731e7d6d34984e91a938e10f6",
					"lastState": {},
					"name": "nginx",
					"ready": true,
					"restartCount": 0,
					"started": true,
					"state": {
						"running": {
							"startedAt": "2024-04-12T18:04:45Z"
						}
					}
				}
			],
			"hostIP": "198.19.249.2",
			"phase": "Running",
			"podIP": "192.168.194.33",
			"podIPs": [
				{
					"ip": "192.168.194.33"
				},
				{
					"ip": "fd07:b51a:cc66:a::21"
				}
			],
			"qosClass": "BestEffort",
			"startTime": "2024-04-12T18:04:44Z"
		}
	}
	`, dsName, name, namespace, dsName))
}

// dsFailedPod returns a Pod created by a DaemonSet that has failed.
func dsFailedPod(namespace, name, dsName string) *unstructured.Unstructured {
	//nolint:lll // Test data.
	return mustDecodeUnstructured(fmt.Sprintf(`
	{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
			"creationTimestamp": "2024-04-12T18:05:18Z",
			"generateName": "%s-",
			"labels": {
				"controller-revision-hash": "6cc847697b",
				"foo": "bar",
				"pod-template-generation": "2"
			},
			"name": "%s",
			"namespace": "%s",
			"ownerReferences": [
				{
					"apiVersion": "apps/v1",
					"blockOwnerDeletion": true,
					"controller": true,
					"kind": "DaemonSet",
					"name": "%s",
					"uid": "d21151ed-7cf4-4da3-86b3-095c40d40c33"
				}
			],
			"resourceVersion": "50187",
			"uid": "50c27d01-b26e-48e6-8924-2699aff914cb"
		},
		"spec": {
			"affinity": {
				"nodeAffinity": {
					"requiredDuringSchedulingIgnoredDuringExecution": {
						"nodeSelectorTerms": [
							{
								"matchFields": [
									{
										"key": "metadata.name",
										"operator": "In",
										"values": [
											"orbstack"
										]
									}
								]
							}
						]
					}
				}
			},
			"containers": [
				{
					"image": "nginx:busted",
					"imagePullPolicy": "IfNotPresent",
					"name": "nginx",
					"resources": {},
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"volumeMounts": [
						{
							"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
							"name": "kube-api-access-qm628",
							"readOnly": true
						}
					]
				}
			],
			"dnsPolicy": "ClusterFirst",
			"enableServiceLinks": true,
			"nodeName": "orbstack",
			"preemptionPolicy": "PreemptLowerPriority",
			"priority": 0,
			"restartPolicy": "Always",
			"schedulerName": "default-scheduler",
			"securityContext": {},
			"serviceAccount": "default",
			"serviceAccountName": "default",
			"terminationGracePeriodSeconds": 30,
			"tolerations": [
				{
					"effect": "NoExecute",
					"key": "node.kubernetes.io/not-ready",
					"operator": "Exists"
				},
				{
					"effect": "NoExecute",
					"key": "node.kubernetes.io/unreachable",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/disk-pressure",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/memory-pressure",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/pid-pressure",
					"operator": "Exists"
				},
				{
					"effect": "NoSchedule",
					"key": "node.kubernetes.io/unschedulable",
					"operator": "Exists"
				}
			],
			"volumes": [
				{
					"name": "kube-api-access-qm628",
					"projected": {
						"defaultMode": 420,
						"sources": [
							{
								"serviceAccountToken": {
									"expirationSeconds": 3607,
									"path": "token"
								}
							},
							{
								"configMap": {
									"items": [
										{
											"key": "ca.crt",
											"path": "ca.crt"
										}
									],
									"name": "kube-root-ca.crt"
								}
							},
							{
								"downwardAPI": {
									"items": [
										{
											"fieldRef": {
												"apiVersion": "v1",
												"fieldPath": "metadata.namespace"
											},
											"path": "namespace"
										}
									]
								}
							}
						]
					}
				}
			]
		},
		"status": {
			"conditions": [
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:05:18Z",
					"status": "True",
					"type": "Initialized"
				},
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:05:18Z",
					"message": "containers with unready status: [nginx]",
					"reason": "ContainersNotReady",
					"status": "False",
					"type": "Ready"
				},
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:05:18Z",
					"message": "containers with unready status: [nginx]",
					"reason": "ContainersNotReady",
					"status": "False",
					"type": "ContainersReady"
				},
				{
					"lastProbeTime": null,
					"lastTransitionTime": "2024-04-12T18:05:18Z",
					"status": "True",
					"type": "PodScheduled"
				}
			],
			"containerStatuses": [
				{
					"image": "nginx:busted",
					"imageID": "",
					"lastState": {},
					"name": "nginx",
					"ready": false,
					"restartCount": 0,
					"started": false,
					"state": {
						"waiting": {
							"message": "rpc error: code = Unknown desc = Error response from daemon: manifest for nginx:busted not found: manifest unknown: manifest unknown",
							"reason": "ErrImagePull"
						}
					}
				}
			],
			"hostIP": "198.19.249.2",
			"phase": "Pending",
			"podIP": "192.168.194.34",
			"podIPs": [
				{
					"ip": "192.168.194.34"
				},
				{
					"ip": "fd07:b51a:cc66:a::22"
				}
			],
			"qosClass": "BestEffort",
			"startTime": "2024-04-12T18:05:18Z"
		}
	}
	`, dsName, name, namespace, dsName))
}
