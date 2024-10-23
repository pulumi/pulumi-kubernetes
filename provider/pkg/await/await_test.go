// Copyright 2021-2024, Pulumi Corporation.
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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/internal"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/watcher"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/managedfields"
	"k8s.io/apimachinery/pkg/util/managedfields/managedfieldstest"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	kfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	kubetesting "k8s.io/client-go/testing"
	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/yaml"
)

var (
	testServerVersion = &cluster.ServerVersion{Major: 1, Minor: 29}

	validPodUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name":      "foo",
				"namespace": "default",
			},
			"spec": map[string]any{
				"containers": []any{
					map[string]any{
						"name":  "foo",
						"image": "nginx",
					},
				},
			},
		},
	}

	unreadyNode = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Node",
			"metadata": map[string]any{
				"name": "foo",
			},
			"status": map[string]any{
				"conditions": []any{
					map[string]any{
						"type":   "Ready",
						"status": "False",
					},
				},
			},
		},
	}

	validClusterRoleUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]any{
				"name": "foo",
			},
			"rules": []any{},
		},
	}

	deprecatedClusterRoleUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "rbac.authorization.k8s.io/v1beta1",
			"kind":       "ClusterRole",
			"metadata": map[string]any{
				"name": "foo",
			},
			"rules": []any{},
		},
	}
)

var (
	podGVR                = corev1.SchemeGroupVersion.WithResource("pods")
	serviceUnavailableErr = apierrors.NewServiceUnavailable("test")
	removedAPIErr         = &kinds.RemovedAPIError{}
)

func TestCreation(t *testing.T) {
	type testCtx struct {
		host   *fakehost.HostClient
		config *CreateConfig
		disco  *fake.SimpleDiscovery
		mapper *fake.SimpleRESTMapper
		client *fake.SimpleDynamicClient
	}
	type args struct {
		preview         bool
		serverSideApply bool
		resType         tokens.Type
		inputs          *unstructured.Unstructured
	}
	type client struct {
		RESTMapperF    func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper
		GenericClientF func(client dynamic.Interface) dynamic.Interface
	}

	type expectF func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error)

	// awaiters

	touch := func(t *testing.T, ctx testCtx) awaiter {
		return func(cac awaitConfig) error {
			require.False(t, metadata.SkipAwaitLogic(cac.currentOutputs), "Await logic should not execute when SkipWait is set")

			// get the live object from the fake API Server
			require.Equal(t, cac.currentOutputs.GetNamespace(), cac.currentOutputs.GetNamespace(), "Live object should have a namespace")
			require.Equal(t, cac.currentOutputs.GetName(), cac.currentOutputs.GetName(), "Live object should have a name")
			gvr, err := clients.GVRForGVK(cac.clientSet.RESTMapper, cac.currentOutputs.GroupVersionKind())
			require.NoError(t, err)
			live, err := ctx.client.Tracker().Get(gvr, cac.currentOutputs.GetNamespace(), cac.currentOutputs.GetName())
			require.NoError(t, err, "Live object should exist in the API Server")
			pod, ok := live.(*unstructured.Unstructured)
			require.True(t, ok)

			// mutate the live object to simulate a observable status update.
			err = unstructured.SetNestedField(pod.Object, "Running", "status", "phase")
			require.NoError(t, err)
			err = ctx.client.Tracker().Update(gvr, live, cac.currentOutputs.GetNamespace())
			require.NoError(t, err)
			return nil
		}
	}

	awaitError := func(t *testing.T, ctx testCtx) awaiter {
		return func(cac awaitConfig) error {
			return serviceUnavailableErr
		}
	}

	awaitUnexpected := func(t *testing.T, ctx testCtx) awaiter {
		return func(cac awaitConfig) error {
			require.Fail(t, "Unexpected call to awaiter")
			return nil
		}
	}

	// expectations

	failed := func(target error) expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.ErrorAs(t, err, &target)
		}
	}
	previewed := func(ns, name string) expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Equal(t, ns, actual.GetNamespace(), "Object should have the expected namespace")
			require.Equal(t, name, actual.GetName(), "Object should have the expected name")
		}
	}
	created := func(ns, name string) expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Equal(t, ns, actual.GetNamespace(), "Object should have the expected namespace")
			require.Equal(t, name, actual.GetName(), "Object should have the expected name")

			gvr, err := clients.GVRForGVK(ctx.mapper, ctx.config.Inputs.GroupVersionKind())
			require.NoError(t, err)
			_, err = ctx.client.Tracker().Get(gvr, ns, name)
			require.NoError(t, err, "Live object should exist in the API Server")
		}
	}
	touched := func() expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.NoError(t, err)
			actualPhase, ok, err := unstructured.NestedString(actual.Object, "status", "phase")
			require.NoError(t, err)
			require.True(t, ok, "Object should have a status.phase field")
			require.Equal(t, actualPhase, "Running", "Object should have status.phase of 'Running'")
		}
	}

	tests := []struct {
		name           string
		client         client
		args           args
		expect         []expectF
		genericEnabled bool
		awaiter        func(t *testing.T, ctx testCtx) awaiter
	}{
		{
			name: "AwaitError",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: awaitError,
			expect:  []expectF{failed(serviceUnavailableErr)},
		},
		{
			name: "NoMatchError",
			client: client{
				RESTMapperF: func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper {
					// return a mapper that returns a NoMatchError until it is reset
					return FlakyRESTMapper(mapper, &meta.NoResourceMatchError{PartialResource: podGVR})
				},
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: touch,
			expect:  []expectF{created("default", "foo" /* after retry */)},
		},
		{
			name: "ServiceUnavailable",
			client: client{
				RESTMapperF: func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper {
					return FailedRESTMapper(mapper, serviceUnavailableErr)
				},
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			expect: []expectF{failed(serviceUnavailableErr)},
		},
		{
			name: "RemovedAPI",
			args: args{
				resType: tokens.Type("kubernetes:rbac.authorization.k8s.io/v1beta1:ClusterRole"),
				inputs:  deprecatedClusterRoleUnstructured,
			},
			expect: []expectF{failed(removedAPIErr)},
		},
		{
			name: "Namespaced",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: touch,
			expect:  []expectF{created("default", "foo"), touched()},
		},
		{
			name: "NonNamespaced",
			args: args{
				resType: tokens.Type("kubernetes:rbac.authorization.k8s.io/v1:ClusterRole"),
				inputs:  validClusterRoleUnstructured,
			},
			awaiter: touch,
			expect:  []expectF{created("", "foo"), touched()},
		},
		{
			name: "GenerateName",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  withGenerateName(validPodUnstructured),
			},
			awaiter: touch,
			expect:  []expectF{created("default", "foo-generated"), touched()},
		},
		{
			name: "SkipAwait",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  withSkipAwait(validPodUnstructured),
			},
			awaiter: awaitUnexpected,
			expect:  []expectF{created("default", "foo")},
		},
		{
			name: "NoAwaiter",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Node"),
				inputs:  unreadyNode,
			},
			awaiter: nil,
			expect:  []expectF{created("", "foo")},
		},
		{
			name: "GenericAwaiter",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Node"),
				inputs:  withReadyCondition(unreadyNode),
			},
			genericEnabled: true,
			expect:         []expectF{created("", "foo")},
		},
		{
			name: "AwaitError",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: awaitError,
			expect:  []expectF{failed(serviceUnavailableErr)},
		},
		{
			name:   "Preview",
			client: client{
				// FUTURE: return a client that requires dry-run mode
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
				preview: true,
			},
			awaiter: awaitUnexpected,
			expect:  []expectF{previewed("default", "foo")},
		},
		// FUTURE: test server-side apply (depends on https://github.com/kubernetes/kubernetes/issues/115598)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host := &fakehost.HostClient{}
			client, disco, mapper, clientset := fake.NewSimpleDynamicClient()
			resources, err := openapi.GetResourceSchemasForClient(disco)
			require.NoError(t, err)

			if tt.client.GenericClientF != nil {
				client.GenericClient = tt.client.GenericClientF(client.GenericClient)
			}
			if tt.client.RESTMapperF != nil {
				client.RESTMapper = tt.client.RESTMapperF(client.RESTMapper)
			}
			if tt.genericEnabled {
				t.Setenv("PULUMI_K8S_AWAIT_ALL", "true")
			}

			urn := resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"), tokens.Type(""), tt.args.resType, "testresource")
			config := CreateConfig{
				ProviderConfig: ProviderConfig{
					Context:           context.Background(),
					Host:              host,
					URN:               urn,
					InitialAPIVersion: corev1.SchemeGroupVersion.String(),
					FieldManager:      "test",
					ClusterVersion:    testServerVersion,
					ClientSet:         client,
					DedupLogger:       logging.NewLogger(context.Background(), host, urn),
					Resources:         resources,
					ServerSideApply:   tt.args.serverSideApply,
					awaiters:          map[string]awaitSpec{},
				},
				Inputs:  tt.args.inputs,
				Preview: tt.args.preview,
			}
			testCtx := testCtx{
				host:   host,
				config: &config,
				disco:  disco,
				mapper: mapper,
				client: clientset,
			}
			if tt.awaiter != nil {
				id := fmt.Sprintf("%s/%s", tt.args.inputs.GetAPIVersion(), tt.args.inputs.GetKind())
				config.awaiters[id] = awaitSpec{
					await: wrap(tt.awaiter(t, testCtx)),
				}
			}
			actual, err := Creation(config)
			for _, e := range tt.expect {
				e(t, testCtx, actual, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type testCtx struct {
		host   *fakehost.HostClient
		config *UpdateConfig
		disco  *fake.SimpleDiscovery
		mapper *fake.SimpleRESTMapper
		client *fake.SimpleDynamicClient
	}
	type args struct {
		preview         bool
		serverSideApply bool
		resType         tokens.Type
		inputs          *unstructured.Unstructured
		oldInputs       *unstructured.Unstructured
		oldOutputs      *unstructured.Unstructured
	}
	type client struct {
		RESTMapperF    func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper
		GenericClientF func(client dynamic.Interface) dynamic.Interface
	}

	type expectF func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error)

	// awaiters

	touch := func(t *testing.T, ctx testCtx) awaiter {
		return func(cac awaitConfig) error {
			require.False(t, metadata.SkipAwaitLogic(cac.currentOutputs), "Await logic should not execute when SkipWait is set")

			// get the live object from the fake API Server
			require.Equal(t, cac.currentOutputs.GetNamespace(), cac.currentOutputs.GetNamespace(), "Live object should have a namespace")
			require.Equal(t, cac.currentOutputs.GetName(), cac.currentOutputs.GetName(), "Live object should have a name")
			gvr, err := clients.GVRForGVK(cac.clientSet.RESTMapper, cac.currentOutputs.GroupVersionKind())
			require.NoError(t, err)
			live, err := ctx.client.Tracker().Get(gvr, cac.currentOutputs.GetNamespace(), cac.currentOutputs.GetName())
			require.NoError(t, err, "Live object should exist in the API Server")
			pod, ok := live.(*unstructured.Unstructured)
			require.True(t, ok)

			// mutate the live object to simulate a observable status update.
			err = unstructured.SetNestedField(pod.Object, "Running", "status", "phase")
			require.NoError(t, err)
			err = ctx.client.Tracker().Update(gvr, live, cac.currentOutputs.GetNamespace())
			require.NoError(t, err)
			return nil
		}
	}

	awaitError := func(t *testing.T, ctx testCtx) awaiter {
		return func(cac awaitConfig) error {
			return serviceUnavailableErr
		}
	}

	awaitUnexpected := func(t *testing.T, ctx testCtx) awaiter {
		return func(cac awaitConfig) error {
			require.Fail(t, "Unexpected call to awaiter")
			return nil
		}
	}

	// expectations

	failed := func(target error) expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.ErrorAs(t, err, &target)
		}
	}
	previewed := func(ns, name string) expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Equal(t, ns, actual.GetNamespace(), "Object should have the expected namespace")
			require.Equal(t, name, actual.GetName(), "Object should have the expected name")
		}
	}
	updated := func(ns, name string) expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Equal(t, ns, actual.GetNamespace(), "Object should have the expected namespace")
			require.Equal(t, name, actual.GetName(), "Object should have the expected name")

			gvr, err := clients.GVRForGVK(ctx.mapper, ctx.config.Inputs.GroupVersionKind())
			require.NoError(t, err)
			_, err = ctx.client.Tracker().Get(gvr, ns, name)
			require.NoError(t, err, "Live object should exist in the API Server")
		}
	}
	touched := func() expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			require.NoError(t, err)
			actualPhase, ok, err := unstructured.NestedString(actual.Object, "status", "phase")
			require.NoError(t, err)
			require.True(t, ok, "Object should have a status.phase field")
			require.Equal(t, actualPhase, "Running", "Object should have status.phase of 'Running'")
		}
	}
	logged := func() expectF {
		return func(t *testing.T, ctx testCtx, actual *unstructured.Unstructured, err error) {
			// FUTURE: assert that a log message was emitted to the fake host
		}
	}

	tests := []struct {
		name           string
		client         client
		args           args
		expect         []expectF
		genericEnabled bool
		awaiter        func(t *testing.T, ctx testCtx) awaiter
	}{
		{
			name: "NoMatchError",
			client: client{
				RESTMapperF: func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper {
					// return a mapper that returns a NoMatchError until it is reset
					return FlakyRESTMapper(mapper, &meta.NoResourceMatchError{PartialResource: podGVR})
				},
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: touch,
			expect:  []expectF{updated("default", "foo" /* after retry */)},
		},
		{
			name: "ServiceUnavailable",
			client: client{
				RESTMapperF: func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper {
					return FailedRESTMapper(mapper, serviceUnavailableErr)
				},
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			expect: []expectF{failed(serviceUnavailableErr)},
		},
		{
			name: "RemovedAPI",
			args: args{
				resType: tokens.Type("kubernetes:rbac.authorization.k8s.io/v1beta1:ClusterRole"),
				inputs:  deprecatedClusterRoleUnstructured,
			},
			expect: []expectF{failed(removedAPIErr)},
		},
		{
			name: "Namespaced",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: touch,
			expect:  []expectF{updated("default", "foo"), touched()},
		},
		{
			name: "NonNamespaced",
			args: args{
				resType: tokens.Type("kubernetes:rbac.authorization.k8s.io/v1:ClusterRole"),
				inputs:  validClusterRoleUnstructured,
			},
			awaiter: touch,
			expect:  []expectF{updated("", "foo"), touched()},
		},
		{
			name: "SkipAwait",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  withSkipAwait(validPodUnstructured),
			},
			awaiter: awaitUnexpected,
			expect:  []expectF{updated("default", "foo")},
		},
		// TODO: Handle this for create as well
		{
			name: "NoAwaiter",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Node"),
				inputs:  unreadyNode,
			},
			awaiter: nil,
			expect:  []expectF{updated("", "foo"), logged()},
		},
		{
			name: "GenericAwaiter",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Node"),
				inputs:  withReadyCondition(unreadyNode),
			},
			genericEnabled: true,
			expect:         []expectF{updated("", "foo"), logged()},
		},
		{
			name: "AwaitError",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
			},
			awaiter: awaitError,
			expect:  []expectF{failed(serviceUnavailableErr)},
		},
		{
			name:   "Preview",
			client: client{
				// FUTURE: return a client that requires dry-run mode
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				inputs:  validPodUnstructured,
				preview: true,
			},
			awaiter: awaitUnexpected,
			expect:  []expectF{previewed("default", "foo")},
		},
		// FUTURE: test server-side apply (depends on https://github.com/kubernetes/kubernetes/issues/115598)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host := &fakehost.HostClient{}
			oldInputs, oldOutputs := tt.args.inputs, tt.args.inputs
			if tt.args.oldInputs != nil {
				oldInputs = tt.args.oldInputs
			}
			if tt.args.oldOutputs != nil {
				oldOutputs = tt.args.oldOutputs
				oldOutputs = oldOutputs.DeepCopy()
				oldOutputs.SetName("old-outputs-name")
			}
			if tt.genericEnabled {
				t.Setenv("PULUMI_K8S_AWAIT_ALL", "true")
			}
			client, disco, mapper, clientset := fake.NewSimpleDynamicClient(fake.WithObjects(oldOutputs))
			resources, err := openapi.GetResourceSchemasForClient(disco)
			require.NoError(t, err)

			if tt.client.GenericClientF != nil {
				client.GenericClient = tt.client.GenericClientF(client.GenericClient)
			}
			if tt.client.RESTMapperF != nil {
				client.RESTMapper = tt.client.RESTMapperF(client.RESTMapper)
			}

			urn := resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"), tokens.Type(""), tt.args.resType, "testresource")
			config := UpdateConfig{
				ProviderConfig: ProviderConfig{
					Context:           context.Background(),
					Host:              host,
					URN:               urn,
					InitialAPIVersion: corev1.SchemeGroupVersion.String(),
					FieldManager:      "test",
					ClusterVersion:    testServerVersion,
					ClientSet:         client,
					DedupLogger:       logging.NewLogger(context.Background(), host, urn),
					Resources:         resources,
					ServerSideApply:   tt.args.serverSideApply,
					awaiters:          map[string]awaitSpec{},
				},
				OldInputs:  oldInputs,
				OldOutputs: oldOutputs,
				Inputs:     tt.args.inputs,
				Preview:    tt.args.preview,
			}
			testCtx := testCtx{
				host:   host,
				config: &config,
				disco:  disco,
				mapper: mapper,
				client: clientset,
			}
			if tt.awaiter != nil {
				id := fmt.Sprintf("%s/%s", tt.args.inputs.GetAPIVersion(), tt.args.inputs.GetKind())
				config.awaiters[id] = awaitSpec{
					await: wrap(tt.awaiter(t, testCtx)),
				}
			}
			actual, err := Update(config)
			for _, e := range tt.expect {
				e(t, testCtx, actual, err)
			}
		})
	}
}

func TestDeletion(t *testing.T) {
	type testCtx struct {
		ctx    context.Context
		cancel context.CancelFunc
		host   *fakehost.HostClient
		config *DeleteConfig
		disco  *fake.SimpleDiscovery
		mapper *fake.SimpleRESTMapper
		client *fake.SimpleDynamicClient
	}
	type args struct {
		serverSideApply bool
		resType         tokens.Type
		objects         []runtime.Object
		inputs          *unstructured.Unstructured
		outputs         *unstructured.Unstructured
		name            string
	}
	type client struct {
		RESTMapperF    func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper
		GenericClientF func(client dynamic.Interface) dynamic.Interface
	}

	type reactionF func(t *testing.T, ctx testCtx, action kubetesting.Action) (handled bool, ret runtime.Object, err error)

	type expectF func(t *testing.T, ctx testCtx, err error)

	// reactions

	suppressDeletion := func(t *testing.T, ctx testCtx, action kubetesting.Action) (bool, runtime.Object, error) {
		return true, ctx.config.Outputs, nil
	}

	cancelAwait := func(t *testing.T, ctx testCtx, action kubetesting.Action) (bool, runtime.Object, error) {
		ctx.cancel()
		return false, nil, nil
	}

	// awaiters

	awaitNoop := func(t *testing.T, ctx testCtx) condition.Satisfier {
		return condition.NewImmediate(ctx.config.DedupLogger, ctx.config.Inputs)
	}

	awaitError := func(t *testing.T, ctx testCtx) condition.Satisfier {
		return condition.NewFailure(serviceUnavailableErr)
	}

	awaitUnexpected := func(t *testing.T, ctx testCtx) condition.Satisfier {
		return condition.NewFailure(fmt.Errorf("unexpected call to await"))
	}

	// expectations

	failed := func(target error) expectF {
		return func(t *testing.T, ctx testCtx, err error) {
			require.ErrorAs(t, err, &target)
		}
	}

	succeeded := func() expectF {
		return func(t *testing.T, ctx testCtx, err error) {
			require.NoError(t, err)
		}
	}
	deleted := func(ns, name string) expectF {
		return func(t *testing.T, ctx testCtx, _ error) {
			gvr, err := clients.GVRForGVK(ctx.mapper, ctx.config.Inputs.GroupVersionKind())
			require.NoError(t, err)
			_, err = ctx.client.Tracker().Get(gvr, ns, name)
			require.Truef(t, apierrors.IsNotFound(err), "expected notfound, got an object")
		}
	}

	tests := []struct {
		name      string
		client    client
		args      args
		expect    []expectF
		reaction  []reactionF
		watcher   *watch.RaceFreeFakeWatcher
		condition func(*testing.T, testCtx) condition.Satisfier
	}{
		{
			name: "ServiceUnavailable",
			client: client{
				RESTMapperF: func(mapper meta.ResettableRESTMapper) meta.ResettableRESTMapper {
					return FailedRESTMapper(mapper, serviceUnavailableErr)
				},
			},
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			expect: []expectF{failed(serviceUnavailableErr)},
		},
		{
			name: "Namespaced",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			condition: awaitNoop,
			expect:    []expectF{succeeded(), deleted("default", "foo")},
		},
		{
			name: "NonNamespaced",
			args: args{
				resType: tokens.Type("kubernetes:rbac.authorization.k8s.io/v1:ClusterRole"),
				name:    "foo",
				objects: []runtime.Object{validClusterRoleUnstructured},
				inputs:  validClusterRoleUnstructured,
				outputs: validClusterRoleUnstructured,
			},
			condition: awaitNoop,
			expect:    []expectF{succeeded(), deleted("default", "foo")},
		},
		{
			name: "Gone",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{ /* empty */ },
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			condition: awaitUnexpected,
			expect:    []expectF{succeeded()},
		},
		{
			name: "SkipAwait",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  withSkipAwait(validPodUnstructured),
				outputs: validPodUnstructured,
			},
			reaction: []reactionF{suppressDeletion}, // suppress deletion to safeguard that the built-in watcher is not used.
			expect:   []expectF{succeeded()},
		},
		{
			name: "AwaitError",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			condition: awaitError,
			expect:    []expectF{failed(serviceUnavailableErr)},
		},
		{
			name: "Deleted",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			condition: awaitNoop,
			expect:    []expectF{succeeded(), deleted("default", "foo")},
		},
		{
			name: "Cancel",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			reaction: []reactionF{cancelAwait, suppressDeletion},
			expect:   []expectF{failed(&cancellationError{})},
		},
		{
			name: "CancelWithRecovery",
			args: args{
				resType: tokens.Type("kubernetes:core/v1:Pod"),
				name:    "foo",
				objects: []runtime.Object{validPodUnstructured},
				inputs:  validPodUnstructured,
				outputs: validPodUnstructured,
			},
			reaction:  []reactionF{cancelAwait},
			condition: awaitNoop,
			expect:    []expectF{succeeded()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			require.NotEmpty(t, tt.name, "Test case must have a name")

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			host := &fakehost.HostClient{}
			client, disco, mapper, clientset := fake.NewSimpleDynamicClient(fake.WithObjects(tt.args.objects...))
			resources, err := openapi.GetResourceSchemasForClient(disco)
			require.NoError(t, err)

			if tt.client.GenericClientF != nil {
				client.GenericClient = tt.client.GenericClientF(client.GenericClient)
			}
			if tt.client.RESTMapperF != nil {
				client.RESTMapper = tt.client.RESTMapperF(client.RESTMapper)
			}

			urn := resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"), tokens.Type(""), tt.args.resType, "testresource")
			config := DeleteConfig{
				ProviderConfig: ProviderConfig{
					Context:           ctx,
					Host:              host,
					URN:               urn,
					InitialAPIVersion: corev1.SchemeGroupVersion.String(),
					FieldManager:      "test",
					ClusterVersion:    testServerVersion,
					ClientSet:         client,
					DedupLogger:       logging.NewLogger(context.Background(), host, urn),
					Resources:         resources,
					ServerSideApply:   tt.args.serverSideApply,
					awaiters:          map[string]awaitSpec{},
				},
				Inputs:  tt.args.inputs,
				Outputs: tt.args.outputs,
				Name:    tt.args.name,
			}
			testCtx := testCtx{
				ctx:    ctx,
				cancel: cancel,
				host:   host,
				config: &config,
				disco:  disco,
				mapper: mapper,
				client: clientset,
			}
			for i := len(tt.reaction) - 1; i >= 0; i-- {
				reaction := tt.reaction[i]
				clientset.PrependReactor("*", "*", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
					return reaction(t, testCtx, action)
				})
			}
			if tt.watcher != nil {
				clientset.PrependWatchReactor("*", func(action kubetesting.Action) (handled bool, ret watch.Interface, err error) {
					return true, tt.watcher, nil
				})
			}
			if tt.condition != nil {
				config.condition = tt.condition(t, testCtx)
			}
			err = Deletion(config)
			for _, e := range tt.expect {
				e(t, testCtx, err)
			}
		})
	}
}

func TestAwaitSSAConflict(t *testing.T) {
	client, _, _, clientset := fake.NewSimpleDynamicClient()

	pod := validPodUnstructured.DeepCopy()
	pod.SetNamespace("default")

	urn := resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"), tokens.Type(""), tokens.Type(""), "testresource")
	pconfig := ProviderConfig{
		Context:         context.Background(),
		Host:            &fakehost.HostClient{},
		URN:             urn,
		FieldManager:    "test",
		ClientSet:       client,
		ServerSideApply: true,
	}
	config := CreateConfig{
		ProviderConfig: pconfig,
		Inputs:         pod,
	}

	// Intercept the SSA and respond with a conflict error.
	clientset.PrependReactor("patch", "pods", func(_ kubetesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewApplyConflict(nil, "conflict")
	})

	wantErr := "\nThe resource managed by field manager \"test\" had an apply conflict:"

	// Attempt to create the pod with SSA.
	_, err := Creation(config)
	assert.ErrorContains(t, err, wantErr)

	// We need a valid pod in our Tracker to avoid the Fake's 404.
	err = clientset.Tracker().Add(pod)
	require.NoError(t, err)

	// Attempt to update the pod with SSA.
	_, err = Update(UpdateConfig{
		ProviderConfig: pconfig,
		OldOutputs:     pod,
		Inputs:         pod,
	})
	assert.ErrorContains(t, err, wantErr)
}

func Test_Watcher_Interface_Cancel(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// Cancel should occur before `WatchUntil` because predicate always returns false.
	err := watcher.ForObject(cancelCtx, &mockResourceInterface{}, "").
		WatchUntil(func(_ *unstructured.Unstructured) bool { return false }, 1*time.Minute)

	_, isPartialErr := err.(PartialError)
	assert.True(t, isPartialErr, "Cancelled watcher should emit `await.PartialError`")
	assert.Equal(t, "Resource operation was cancelled for ''", err.Error())
}

func Test_Watcher_Interface_Timeout(t *testing.T) {
	// Timeout because the `WatchUntil` predicate always returns false.
	err := watcher.ForObject(context.Background(), &mockResourceInterface{}, "").
		WatchUntil(func(_ *unstructured.Unstructured) bool { return false }, 1*time.Second)

	_, isPartialErr := err.(PartialError)
	assert.True(t, isPartialErr, "Timed out watcher should emit `await.PartialError`")
	assert.Equal(t, "Timeout occurred polling for ''", err.Error())
}

func TestAwaiterInterfaceTimeout(t *testing.T) {
	awaiter, err := internal.NewAwaiter(internal.WithCondition(condition.NewNever(nil)))
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = awaiter.Await(ctx)
	_, isPartialErr := err.(PartialError)
	assert.True(t, isPartialErr, "Timed out watcher should emit `await.PartialError`")
}

// --------------------------------------------------------------------------

// Helpers

// --------------------------------------------------------------------------

func withSkipAwait(obj *unstructured.Unstructured) *unstructured.Unstructured {
	copy := obj.DeepCopy()
	copy.SetAnnotations(map[string]string{
		"pulumi.com/skipAwait": "true",
	})
	return copy
}

func withGenerateName(obj *unstructured.Unstructured) *unstructured.Unstructured {
	copy := obj.DeepCopy()
	copy.SetGenerateName(fmt.Sprintf("%s-", obj.GetName()))
	copy.SetName("")
	return copy
}

func withReadyCondition(obj *unstructured.Unstructured) *unstructured.Unstructured {
	copy := obj.DeepCopy()
	copy.Object["status"] = map[string]any{
		"conditions": []any{map[string]any{
			"type":   "Ready",
			"status": "True",
		}},
		"phase": "Running",
	}
	return copy
}

// --------------------------------------------------------------------------

// Mock implementations of Kubernetes client stuff.

// --------------------------------------------------------------------------

type mockResourceInterface struct{}

var _ dynamic.ResourceInterface = (*mockResourceInterface)(nil)

func (mri *mockResourceInterface) Create(
	ctx context.Context, obj *unstructured.Unstructured, options metav1.CreateOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Create not implemented")
}

func (mri *mockResourceInterface) Update(
	ctx context.Context, obj *unstructured.Unstructured, options metav1.UpdateOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Update not implemented")
}

func (mri *mockResourceInterface) UpdateStatus(
	ctx context.Context, obj *unstructured.Unstructured, options metav1.UpdateOptions,
) (*unstructured.Unstructured, error) {
	panic("UpdateStatus not implemented")
}

func (mri *mockResourceInterface) Delete(ctx context.Context, name string, options metav1.DeleteOptions, subresources ...string) error {
	panic("Delete not implemented")
}

func (mri *mockResourceInterface) DeleteCollection(
	ctx context.Context, deleteOptions metav1.DeleteOptions, listOptions metav1.ListOptions,
) error {
	panic("DeleteCollection not implemented")
}

func (mri *mockResourceInterface) Get(
	ctx context.Context, name string, options metav1.GetOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{Object: map[string]any{}}, nil
}

func (mri *mockResourceInterface) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	panic("List not implemented")
}

func (mri *mockResourceInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("Watch not implemented")
}

func (mri *mockResourceInterface) Patch(
	ctx context.Context, name string, pt types.PatchType, data []byte, options metav1.PatchOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Patch not implemented")
}

func (mri *mockResourceInterface) Apply(ctx context.Context, name string, obj *unstructured.Unstructured, options metav1.ApplyOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Apply not implemented")
}

func (mri *mockResourceInterface) ApplyStatus(ctx context.Context, name string, obj *unstructured.Unstructured, options metav1.ApplyOptions,
) (*unstructured.Unstructured, error) {
	panic("ApplyStatus not implemented")
}

func FlakyRESTMapper(mapper meta.ResettableRESTMapper, resettable error, extraErrors ...error) *fake.StubResettableRESTMapper {
	return &fake.StubResettableRESTMapper{
		ResettableRESTMapper: mapper,
		ResetF: func() {
			resettable = nil
			mapper.Reset()
		},
		RESTMappingF: func(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
			if resettable != nil {
				return nil, resettable
			}
			if len(extraErrors) == 0 {
				return mapper.RESTMapping(gk, versions...)
			}
			var e error
			e, extraErrors = extraErrors[0], extraErrors[1:]
			return nil, e
		},
	}
}

func FailedRESTMapper(mapper meta.ResettableRESTMapper, err error) *fake.StubResettableRESTMapper {
	return &fake.StubResettableRESTMapper{
		ResettableRESTMapper: mapper,
		RESTMappingF: func(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
			return nil, err
		},
	}
}

func fakeTypeConverter(t *testing.T) managedfields.TypeConverter {
	t.Helper()

	openapi, err := fake.LoadOpenAPISchema()
	require.NoError(t, err)

	swagger := spec.Swagger{}
	raw, err := openapi.YAMLValue("")
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(raw, &swagger))

	definitions := map[string]*spec.Schema{}
	for k, v := range swagger.Definitions {
		p := v
		definitions[k] = &p
	}

	tc, err := managedfields.NewTypeConverter(definitions, false)
	require.NoError(t, err)
	return tc
}

// TestSSAWithOldObjects is a regression test for
// https://github.com/pulumi/customer-support/issues/1837. An object is created
// and manipulated such that it no longer has any .metadata.managedFields, as
// is the case with things created prior to 1.18. We confirm this reproduces
// the issue and that our SSA upgrade logic handles it.
func TestSSAWithOldObjects(t *testing.T) {
	t.Parallel()

	tc := fakeTypeConverter(t)
	fm := managedfieldstest.NewTestFieldManager(tc, schema.FromAPIVersionAndKind("v1", "Service"))

	inputs := &unstructured.Unstructured{}
	in := `{
		"apiVersion": "v1",
		"kind": "Service",
		"metadata": {
			"labels": {
				"app.kubernetes.io/instance": "autoscaler",
				"app.kubernetes.io/managed-by": "pulumi",
				"app.kubernetes.io/name": "aws-cluster-autoscaler",
				"app.kubernetes.io/version": "1.28.2",
				"helm.sh/chart": "cluster-autoscaler-9.34.1"
			},
			"name": "cluster-autoscaler",
			"namespace": "kube-system"
		},
		"spec": {}
	}`
	require.NoError(t, json.Unmarshal([]byte(in), inputs))
	// We need the last-applied-config annotation in order to trigger kubectl's
	// graceful CSA->SSA upgrade path.
	last, err := inputs.MarshalJSON()
	require.NoError(t, err)
	inputs.SetAnnotations(map[string]string{
		"kubectl.kubernetes.io/last-applied-configuration": string(last),
	})

	// Create the object. As of 1.18 all objects are created with
	// managedFields -- even when using CSA.
	obj := inputs.DeepCopy()
	err = fm.Update(obj, "kubectl-create")
	require.NoError(t, err)
	require.NotEmpty(t, fm.ManagedFields())
	assert.Len(t, fm.ManagedFields(), 1)

	// However we can still disable managed fields after creating the object by
	// explicitly setting it to `[]`.
	obj = inputs.DeepCopy()
	obj.SetManagedFields([]metav1.ManagedFieldsEntry{})
	err = fm.Update(obj, "kubectl-update")
	require.NoError(t, err)
	assert.Empty(t, fm.ManagedFields())

	// Try to update a label on the object using a naive apply.
	obj = inputs.DeepCopy()
	obj.SetLabels(map[string]string{
		"helm.sh/chart": "cluster-autoscaler-9.36.0",
	})
	// Despite having no field managers, our apply still conflicts with the
	// legacy "before-first-apply" manager.
	err = fm.Apply(obj, "pulumi-kubernetes", false)
	assert.ErrorContains(t, err, `Apply failed with 1 conflict: conflict with "before-first-apply" using v1: .metadata.labels.helm.sh/chart`)

	// Now try again using our SSA upgrade logic -- this should succeed.
	cfg := &UpdateConfig{
		Inputs:  obj,
		Preview: false,
		ProviderConfig: ProviderConfig{
			URN:             resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"), tokens.Type(""), "v1/Service", "testresource"),
			FieldManager:    "pulumi-kubernetes",
			ServerSideApply: true,
		},
	}
	_, err = ssaUpdate(cfg, obj, fieldManagerPatcher{fm})
	require.NoError(t, err)
}

func TestSSAUpdate(t *testing.T) {
	tests := []struct {
		name         string
		obj          string
		preview      bool
		wantManagers []string
	}{
		{
			name: "we take ownership of kubectl CSA",
			obj: `apiVersion: v1
kind: Namespace
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Namespace","metadata":{"annotations":{},"labels":{"app.kubernetes.io/instance":"flux-system","app.kubernetes.io/part-of":"flux","pod-security.kubernetes.io/warn":"restricted","pod-security.kubernetes.io/warn-version":"latest"},"name":"flux-system"}}
  creationTimestamp: "2024-09-24T19:27:32Z"
  labels:
    app.kubernetes.io/instance: flux-system
    app.kubernetes.io/part-of: flux
    kubernetes.io/metadata.name: flux-system
    pod-security.kubernetes.io/warn: restricted
    pod-security.kubernetes.io/warn-version: latest
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .: {}
          f:kubectl.kubernetes.io/last-applied-configuration: {}
        f:labels:
          .: {}
          f:app.kubernetes.io/instance: {}
          f:app.kubernetes.io/part-of: {}
          f:kubernetes.io/metadata.name: {}
          f:pod-security.kubernetes.io/warn: {}
          f:pod-security.kubernetes.io/warn-version: {}
    manager: kubectl-client-side-apply
    operation: Update
    time: "2024-09-24T19:27:32Z"
  name: flux-system
  resourceVersion: "138234"
  uid: c14c35d8-ae5d-4f53-8391-791d47efe337
spec:
  finalizers:
  - kubernetes
status:
  phase: Active`,
			preview:      false,
			wantManagers: []string{"pulumi-kubernetes"},
		},
		{
			name: "we take ownership of kubectl SSA",
			obj: `apiVersion: v1
kind: Namespace
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Namespace","metadata":{"annotations":{},"labels":{"app.kubernetes.io/instance":"flux-system","app.kubernetes.io/part-of":"flux","pod-security.kubernetes.io/warn":"restricted","pod-security.kubernetes.io/warn-version":"latest"},"name":"flux-system"}}
  creationTimestamp: "2024-09-24T19:27:32Z"
  labels:
    app.kubernetes.io/instance: flux-system
    app.kubernetes.io/part-of: flux
    kubernetes.io/metadata.name: flux-system
    pod-security.kubernetes.io/warn: restricted
    pod-security.kubernetes.io/warn-version: latest
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .: {}
          f:kubectl.kubernetes.io/last-applied-configuration: {}
        f:labels:
          .: {}
          f:app.kubernetes.io/instance: {}
          f:app.kubernetes.io/part-of: {}
          f:kubernetes.io/metadata.name: {}
          f:pod-security.kubernetes.io/warn: {}
          f:pod-security.kubernetes.io/warn-version: {}
    manager: kubectl-client-side-apply
    operation: Update
    time: "2024-09-24T19:27:32Z"
  name: flux-system
  resourceVersion: "138234"
  uid: c14c35d8-ae5d-4f53-8391-791d47efe337
spec:
  finalizers:
  - kubernetes
status:
  phase: Active`,
			preview:      false,
			wantManagers: []string{"pulumi-kubernetes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj unstructured.Unstructured
			require.NoError(t, yaml.Unmarshal([]byte(tt.obj), &obj))

			typed, err := scheme.Scheme.New(obj.GroupVersionKind())
			require.NoError(t, err)

			require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, typed))

			// client, _, _, _ := fake.NewSimpleDynamicClient(fake.WithObjects(typed))
			//
			client := kfake.NewClientset(typed)
			c := client.CoreV1().Namespaces()

			inputs := obj.DeepCopy()
			inputs.SetLabels(nil)
			inputs.SetManagedFields(nil)
			cfg := &UpdateConfig{
				Inputs:  inputs,
				Preview: tt.preview,
				ProviderConfig: ProviderConfig{
					URN:             resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"), tokens.Type(""), "v1/Service", "testresource"),
					FieldManager:    "pulumi-kubernetes",
					ServerSideApply: true,
				},
			}
			live, err := ssaUpdate(cfg, &obj, untypedPatcher[*corev1.Namespace]{wrapped: c})
			require.NoError(t, err)
			assert.Len(t, live.GetManagedFields(), 1)
			for idx, want := range tt.wantManagers {
				assert.Equal(t, want, live.GetManagedFields()[idx].Manager)
			}
		})
	}
}

type typedPatcher[T runtime.Object] interface {
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (T, error)
}

type untypedPatcher[T runtime.Object] struct {
	wrapped typedPatcher[T]
}

func (p untypedPatcher[T]) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, options metav1.PatchOptions, subresources ...string) (*unstructured.Unstructured, error) {
	typed, err := p.wrapped.Patch(ctx, name, pt, data, options, subresources...)
	if err != nil {
		return nil, err
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(typed)
	return &unstructured.Unstructured{Object: obj}, err
}

type fieldManagerPatcher struct {
	fm managedfieldstest.TestFieldManager
}

func (p fieldManagerPatcher) Patch(_ context.Context, _ string, pt types.PatchType, data []byte, options metav1.PatchOptions, _ ...string) (*unstructured.Unstructured, error) {
	if pt != types.ApplyPatchType {
		return nil, fmt.Errorf("fieldManagerPatcher only handles Apply")
	}

	force := false
	if options.Force != nil {
		force = *options.Force
	}

	in, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(in, nil, nil)
	if err != nil {
		return nil, err
	}

	err = p.fm.Apply(obj, options.FieldManager, force)
	if err != nil {
		return nil, err
	}

	live := p.fm.Live()
	return live.(*unstructured.Unstructured), err
}
