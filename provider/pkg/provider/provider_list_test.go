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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	fakeclients "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
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
		{"negative limit", &pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap", Limit: -1}},
		{"negative page size", &pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap", PageSize: -1}},
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
	err = k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap"}, &listResponseStream{})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestList_UnknownPackageRejected(t *testing.T) {
	k := newListTestProvider(t)
	err := k.List(&pulumirpc.ListRequest{Token: "aws:s3/bucket:Bucket"}, &listResponseStream{})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestList_BasicReturnsAllResources(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", map[string]string{"env": "prod"}),
		configMap("ns-2", "cm-2", map[string]string{"env": "dev"}),
	)
	stream := &listResponseStream{}
	err := k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap"}, stream)
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
	err = k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap", Query: q}, stream)
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
	err = k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap", Query: q}, stream)
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
	err = k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMap", Query: q}, stream)
	require.NoError(t, err)
	assert.Equal(t, []string{"ns-1/cm-prod"}, collectIDs(stream))
}

func TestList_PatchSuffixStripped(t *testing.T) {
	k := newListTestProvider(t,
		configMap("ns-1", "cm-1", nil),
	)
	stream := &listResponseStream{}
	err := k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:ConfigMapPatch"}, stream)
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
	err = k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:Namespace", Query: q}, &listResponseStream{})
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
	err := k.List(&pulumirpc.ListRequest{Token: "kubernetes:core/v1:Namespace"}, stream)
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
		token    string
		ok       bool
		group    string
		version  string
		kind     string
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
