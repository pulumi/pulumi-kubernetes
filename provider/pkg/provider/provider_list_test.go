// Copyright 2016-2026, Pulumi Corporation.
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

package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	fakeclients "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
)

// Pulumi resource type tokens hoisted out to silence gosec G101 false
// positives on `Token: "..."` literals (the proto field name "Token" matches
// the linter's credential-name regex).
const (
	configMapType      = "kubernetes:core/v1:ConfigMap"
	configMapPatchType = "kubernetes:core/v1:ConfigMapPatch"
	namespaceType      = "kubernetes:core/v1:Namespace"
	awsBucketType      = "aws:s3/bucket:Bucket"
)

// listResponseStream is an in-memory ListResponse stream for testing.
type listResponseStream struct {
	ctx       context.Context
	responses []*pulumirpc.ListResponse
}

func (s *listResponseStream) Send(resp *pulumirpc.ListResponse) error {
	s.responses = append(s.responses, resp)
	return nil
}

func (s *listResponseStream) SetHeader(metadata.MD) error  { return nil }
func (s *listResponseStream) SendHeader(metadata.MD) error { return nil }
func (s *listResponseStream) SetTrailer(metadata.MD)       {}
func (s *listResponseStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}
func (s *listResponseStream) SendMsg(any) error { return nil }
func (s *listResponseStream) RecvMsg(any) error { return nil }

// newListTestProvider returns a *kubeProvider wired to a fake clientset pre-populated with objs.
func newListTestProvider(t *testing.T, objs ...runtime.Object) *kubeProvider {
	t.Helper()
	k, err := makeKubeProvider(nil, "kubernetes", "v0.0.0", []byte("{}"), []byte("{}"), []byte("{}"))
	require.NoError(t, err)
	cs, _, _, _ := fakeclients.NewSimpleDynamicClient(fakeclients.WithObjects(objs...))
	k.clientSet = cs
	return k
}

// startKubeProviderServer runs k as an in-process gRPC server and returns a connected client.
func startKubeProviderServer(t *testing.T, k *kubeProvider) pulumirpc.ResourceProviderClient {
	t.Helper()
	cancel := make(chan bool)
	t.Cleanup(func() { close(cancel) })

	handle, err := rpcutil.ServeWithOptions(rpcutil.ServeOptions{
		Cancel: cancel,
		Init: func(srv *grpc.Server) error {
			pulumirpc.RegisterResourceProviderServer(srv, k)
			return nil
		},
	})
	require.NoError(t, err)
	go func() {
		if err := <-handle.Done; err != nil {
			t.Errorf("kube provider gRPC server: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", handle.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return pulumirpc.NewResourceProviderClient(conn)
}

// drainList consumes a List response stream and returns the result IDs and the
// final continuation token (empty if none was sent).
func drainList(t *testing.T, stream grpc.ServerStreamingClient[pulumirpc.ListResponse]) ([]string, string) {
	t.Helper()
	var ids []string
	var contToken string
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		if r := resp.GetResult(); r != nil {
			ids = append(ids, r.GetId())
		}
		if c := resp.GetContinuation(); c != nil {
			contToken = c.GetContinuationToken()
		}
	}
	return ids, contToken
}

func configMap(ns, name string, labels map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, Labels: labels},
	}
}

func collectIDs(stream *listResponseStream) []string {
	var ids []string
	for _, r := range stream.responses {
		if res := r.GetResult(); res != nil {
			ids = append(ids, res.GetId())
		}
	}
	return ids
}

func TestList_NilRequest(t *testing.T) {
	k := newListTestProvider(t)
	err := k.List(nil, &listResponseStream{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestList_RejectsInvalidArgs(t *testing.T) {
	k := newListTestProvider(t)
	cases := []struct {
		name string
		req  *pulumirpc.ListRequest
	}{
		{"empty token", &pulumirpc.ListRequest{}},
		{"negative limit", &pulumirpc.ListRequest{Token: configMapType, Limit: -1}},
		{"negative page size", &pulumirpc.ListRequest{Token: configMapType, PageSize: -1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := k.List(tc.req, &listResponseStream{})
			require.Error(t, err)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}

func TestList_NotConfigured(t *testing.T) {
	k, err := makeKubeProvider(nil, "kubernetes", "v0.0.0", []byte("{}"), []byte("{}"), []byte("{}"))
	require.NoError(t, err)
	err = k.List(&pulumirpc.ListRequest{Token: configMapType}, &listResponseStream{})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestList_UnknownPackageRejected(t *testing.T) {
	k := newListTestProvider(t)
	err := k.List(&pulumirpc.ListRequest{Token: awsBucketType}, &listResponseStream{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestList_BasicReturnsAllResources(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", map[string]string{"env": "prod"}),
		configMap("ns-2", "cm-2", map[string]string{"env": "dev"}),
	)
	stream := &listResponseStream{}
	err := k.List(&pulumirpc.ListRequest{Token: configMapType}, stream)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"ns-1/cm-1", "ns-2/cm-2"}, collectIDs(stream))
}

func TestList_NamespaceFilter(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", nil),
		configMap("ns-2", "cm-2", nil),
	)
	q, err := structpb.NewStruct(map[string]any{"namespace": "ns-1"})
	require.NoError(t, err)
	stream := &listResponseStream{}
	err = k.List(&pulumirpc.ListRequest{Token: configMapType, Query: q}, stream)
	require.NoError(t, err)
	assert.Equal(t, []string{"ns-1/cm-1"}, collectIDs(stream))
}

func TestList_MetadataNamespaceFallback(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", nil),
		configMap("ns-2", "cm-2", nil),
	)
	q, err := structpb.NewStruct(map[string]any{"metadata": map[string]any{"namespace": "ns-2"}})
	require.NoError(t, err)
	stream := &listResponseStream{}
	err = k.List(&pulumirpc.ListRequest{Token: configMapType, Query: q}, stream)
	require.NoError(t, err)
	assert.Equal(t, []string{"ns-2/cm-2"}, collectIDs(stream))
}

func TestList_LabelSelector(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-prod", map[string]string{"env": "prod"}),
		configMap("ns-1", "cm-dev", map[string]string{"env": "dev"}),
	)
	q, err := structpb.NewStruct(map[string]any{"labelSelector": "env=prod"})
	require.NoError(t, err)
	stream := &listResponseStream{}
	err = k.List(&pulumirpc.ListRequest{Token: configMapType, Query: q}, stream)
	require.NoError(t, err)
	assert.Equal(t, []string{"ns-1/cm-prod"}, collectIDs(stream))
}

func TestList_PatchSuffixStripped(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", nil),
	)
	stream := &listResponseStream{}
	err := k.List(&pulumirpc.ListRequest{Token: configMapPatchType}, stream)
	require.NoError(t, err)
	assert.Equal(t, []string{"ns-1/cm-1"}, collectIDs(stream))
}

func TestList_ClusterScopedRejectsNamespaceQuery(t *testing.T) {
	k := newListTestProvider(t,
		&corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: "ns-1"},
		},
	)
	q, err := structpb.NewStruct(map[string]any{"namespace": "x"})
	require.NoError(t, err)
	err = k.List(&pulumirpc.ListRequest{Token: namespaceType, Query: q}, &listResponseStream{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestList_ClusterScopedListsWithoutNamespace(t *testing.T) {
	k := newListTestProvider(t,
		&corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: "ns-1"},
		},
		&corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{Name: "ns-2"},
		},
	)
	stream := &listResponseStream{}
	err := k.List(&pulumirpc.ListRequest{Token: namespaceType}, stream)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"ns-1", "ns-2"}, collectIDs(stream))
}

func TestParseListQuery_NilReturnsZero(t *testing.T) {
	q, err := parseListQuery(nil)
	require.NoError(t, err)
	assert.Equal(t, listQuery{}, q)
}

func TestParseListQuery_TopLevelOverridesMetadata(t *testing.T) {
	s, err := structpb.NewStruct(map[string]any{
		"namespace": "top",
		"metadata":  map[string]any{"namespace": "meta"},
	})
	require.NoError(t, err)
	q, err := parseListQuery(s)
	require.NoError(t, err)
	assert.Equal(t, "top", q.namespace)
}

func TestParseListQuery_RejectsWrongTypes(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]any
	}{
		{"namespace non-string", map[string]any{"namespace": 1}},
		{"name non-string", map[string]any{"name": true}},
		{"labelSelector non-string", map[string]any{"labelSelector": []any{"x"}}},
		{"fieldSelector non-string", map[string]any{"fieldSelector": 3.14}},
		{"metadata non-object", map[string]any{"metadata": "x"}},
		{"metadata.namespace non-string", map[string]any{"metadata": map[string]any{"namespace": 5}}},
		{"metadata.name non-string", map[string]any{"metadata": map[string]any{"name": false}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := structpb.NewStruct(tc.in)
			require.NoError(t, err)
			_, err = parseListQuery(s)
			require.Error(t, err)
		})
	}
}

func TestGvkFromTypeToken(t *testing.T) {
	k := &kubeProvider{providerPackage: "kubernetes"}
	cases := []struct {
		token   string
		ok      bool
		group   string
		version string
		kind    string
	}{
		{"kubernetes:core/v1:ConfigMap", true, "", "v1", "ConfigMap"},
		{"kubernetes:apps/v1:Deployment", true, "apps", "v1", "Deployment"},
		{"kubernetes:core/v1:ConfigMapPatch", true, "", "v1", "ConfigMapPatch"},
		{"aws:s3/bucket:Bucket", false, "", "", ""},
		{"kubernetes:core:ConfigMap", false, "", "", ""},
		{"kubernetes:core/v1:", false, "", "", ""},
		{"kubernetes::ConfigMap", false, "", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.token, func(t *testing.T) {
			gvk, err := k.gvkFromTypeToken(tc.token)
			if !tc.ok {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.group, gvk.Group)
			assert.Equal(t, tc.version, gvk.Version)
			assert.Equal(t, tc.kind, gvk.Kind)
		})
	}
}

// TestList_ClusterUnreachable verifies that List surfaces a descriptive error when the provider
// recorded the cluster as unreachable during Configure.
func TestList_ClusterUnreachable(t *testing.T) {
	k := newListTestProvider(t)
	k.clusterUnreachable = true
	k.clusterUnreachableReason = "test-only: cluster down"

	err := k.List(&pulumirpc.ListRequest{Token: configMapType}, &listResponseStream{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "test-only: cluster down")
}

// TestList_NameAndFieldSelectorCompose verifies that query.name and query.fieldSelector are
// composed (the name filter is appended to any caller-supplied fieldSelector).
func TestList_NameAndFieldSelectorCompose(t *testing.T) {
	k := newListTestProvider(t, configMap("ns-1", "cm-target", nil))
	q, err := structpb.NewStruct(map[string]any{
		"namespace":     "ns-1",
		"name":          "cm-target",
		"fieldSelector": "status.phase=Active",
	})
	require.NoError(t, err)
	stream := &listResponseStream{}
	err = k.List(&pulumirpc.ListRequest{Token: configMapType, Query: q}, stream)
	require.NoError(t, err)
	assert.Equal(t, []string{"ns-1/cm-target"}, collectIDs(stream))
}

// TestList_OverGRPC exercises List end-to-end through a real gRPC server stream,
// confirming the proto oneof wrappers and stream.Send/Recv work on the wire.
func TestList_OverGRPC(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", nil),
		configMap("ns-2", "cm-2", nil),
	)
	client := startKubeProviderServer(t, k)

	stream, err := client.List(context.Background(), &pulumirpc.ListRequest{
		Token: configMapType,
	})
	require.NoError(t, err)

	ids, contToken := drainList(t, stream)
	assert.ElementsMatch(t, []string{"ns-1/cm-1", "ns-2/cm-2"}, ids)
	assert.Empty(t, contToken, "two results fit in one page; no continuation token expected")
}

// paginatingFake wraps dynamic.Interface and intercepts List on a per-GVR item pool
// to honor opts.Limit and opts.Continue. All other operations delegate to the embedded interface.
type paginatingFake struct {
	dynamic.Interface
	pools map[schema.GroupVersionResource][]unstructured.Unstructured
}

func (f *paginatingFake) Resource(gvr schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	return &paginatingNamespaceable{
		NamespaceableResourceInterface: f.Interface.Resource(gvr),
		pool:                           f.pools[gvr],
	}
}

type paginatingNamespaceable struct {
	dynamic.NamespaceableResourceInterface
	pool []unstructured.Unstructured
}

func (n *paginatingNamespaceable) Namespace(ns string) dynamic.ResourceInterface {
	return &paginatingResource{
		ResourceInterface: n.NamespaceableResourceInterface.Namespace(ns),
		pool:              n.pool,
	}
}

func (n *paginatingNamespaceable) List(
	_ context.Context, opts metav1.ListOptions,
) (*unstructured.UnstructuredList, error) {
	return paginatePool(n.pool, opts), nil
}

type paginatingResource struct {
	dynamic.ResourceInterface
	pool []unstructured.Unstructured
}

func (r *paginatingResource) List(_ context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return paginatePool(r.pool, opts), nil
}

// paginatePool faithfully implements server-side pagination: slice the pool by
// opts.Continue (encoded as "page-N", where N is the absolute starting offset),
// cap to opts.Limit items per call, and emit a continue token when more remain.
func paginatePool(pool []unstructured.Unstructured, opts metav1.ListOptions) *unstructured.UnstructuredList {
	start := 0
	if opts.Continue != "" {
		_, _ = fmt.Sscanf(opts.Continue, "page-%d", &start)
	}
	end := len(pool)
	if opts.Limit > 0 {
		candidate := start + int(opts.Limit)
		if candidate < end {
			end = candidate
		}
	}
	list := &unstructured.UnstructuredList{}
	list.SetAPIVersion("v1")
	list.SetKind("ConfigMapList")
	if start < end {
		list.Items = append([]unstructured.Unstructured{}, pool[start:end]...)
	}
	if end < len(pool) {
		list.SetContinue(fmt.Sprintf("page-%d", end))
	}
	return list
}

// newSessionTestProvider returns a *kubeProvider whose ConfigMap REST calls are
// served by the paginatingFake over the given pool of items.
func newSessionTestProvider(t *testing.T, ns string, count int) *kubeProvider {
	t.Helper()
	k, err := makeKubeProvider(nil, "kubernetes", "v0.0.0", []byte("{}"), []byte("{}"), []byte("{}"))
	require.NoError(t, err)
	cs, _, _, _ := fakeclients.NewSimpleDynamicClient()
	cs.GenericClient = &paginatingFake{
		Interface: cs.GenericClient,
		pools: map[schema.GroupVersionResource][]unstructured.Unstructured{
			{Group: "", Version: "v1", Resource: "configmaps"}: makeNConfigMaps(ns, count),
		},
	}
	k.clientSet = cs
	return k
}

func makeNConfigMaps(ns string, count int) []unstructured.Unstructured {
	items := make([]unstructured.Unstructured, count)
	for i := range items {
		items[i].SetAPIVersion("v1")
		items[i].SetKind("ConfigMap")
		items[i].SetNamespace(ns)
		items[i].SetName(fmt.Sprintf("cm-%03d", i))
	}
	return items
}

// runListSession simulates the engine's pagination loop: call List repeatedly,
// chaining continuation tokens, until the provider stops emitting them.
// Returns all IDs received across calls and the total call count.
func runListSession(
	t *testing.T, client pulumirpc.ResourceProviderClient, base *pulumirpc.ListRequest,
) ([]string, int) {
	t.Helper()
	var allIDs []string
	var contToken string
	callCount := 0
	for {
		callCount++
		require.Less(t, callCount, 50, "infinite-loop guard — provider should signal stop by now")

		// Build a fresh ListRequest per call rather than copying *base (which contains
		// a sync.Mutex inside its proto MessageState — copylocks flags struct copies).
		stream, err := client.List(context.Background(), &pulumirpc.ListRequest{
			Token:             base.GetToken(),
			Query:             base.GetQuery(),
			Limit:             base.GetLimit(),
			PageSize:          base.GetPageSize(),
			ContinuationToken: contToken,
		})
		require.NoError(t, err)
		ids, nextToken := drainList(t, stream)
		allIDs = append(allIDs, ids...)

		if nextToken == "" {
			break
		}
		contToken = nextToken
	}
	return allIDs, callCount
}

func assertAllUnique(t *testing.T, ids []string) {
	t.Helper()
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID across paginated calls: %s", id)
		}
		seen[id] = true
	}
}

func TestList_Session_PageEvenlyDividesLimit(t *testing.T) {
	k := newSessionTestProvider(t, "ns-1", 100)
	client := startKubeProviderServer(t, k)

	ids, callCount := runListSession(t, client, &pulumirpc.ListRequest{
		Token:    configMapType,
		PageSize: 25,
		Limit:    50,
	})
	assert.Equal(t, 2, callCount, "page_size=25 + limit=50 should complete in 2 calls")
	assert.Len(t, ids, 50, "must respect session cap")
	assertAllUnique(t, ids)
}

func TestList_Session_PageDoesNotEvenlyDivideLimit(t *testing.T) {
	k := newSessionTestProvider(t, "ns-1", 100)
	client := startKubeProviderServer(t, k)

	ids, callCount := runListSession(t, client, &pulumirpc.ListRequest{
		Token:    configMapType,
		PageSize: 20,
		Limit:    50,
	})
	assert.Equal(t, 3, callCount, "page_size=20 + limit=50 should complete in 3 calls (20+20+10)")
	assert.Len(t, ids, 50, "must respect session cap exactly, not overshoot")
	assertAllUnique(t, ids)
}

func TestList_Session_K8sRunsOutBeforeLimit(t *testing.T) {
	k := newSessionTestProvider(t, "ns-1", 30)
	client := startKubeProviderServer(t, k)

	ids, callCount := runListSession(t, client, &pulumirpc.ListRequest{
		Token:    configMapType,
		PageSize: 10,
		Limit:    50,
	})
	assert.Equal(t, 3, callCount)
	assert.Len(t, ids, 30, "cluster only has 30 items; cap is unreachable")
	assertAllUnique(t, ids)
}

func TestList_Session_NoLimit(t *testing.T) {
	k := newSessionTestProvider(t, "ns-1", 50)
	client := startKubeProviderServer(t, k)

	ids, callCount := runListSession(t, client, &pulumirpc.ListRequest{
		Token:    configMapType,
		PageSize: 10,
	})
	assert.Equal(t, 5, callCount, "with no cap, paginate until K8s runs out")
	assert.Len(t, ids, 50)
	assertAllUnique(t, ids)
}

func TestList_Session_LimitOnly_OneCall(t *testing.T) {
	k := newSessionTestProvider(t, "ns-1", 100)
	client := startKubeProviderServer(t, k)

	ids, callCount := runListSession(t, client, &pulumirpc.ListRequest{
		Token: configMapType,
		Limit: 50,
	})
	assert.Equal(t, 1, callCount, "no page_size means one call, sized to the limit")
	assert.Len(t, ids, 50)
	assertAllUnique(t, ids)
}
