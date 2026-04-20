// Copyright 2016-2024, Pulumi Corporation.
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

	_ "embed"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
	structpb "google.golang.org/protobuf/types/known/structpb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

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

var _ = gk.Describe("RPC:Cancel", func() {
	var k *kubeProvider

	gk.JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	gk.It("should cancel any future operations", func() {
		_, err := k.Cancel(context.Background(), &pbempty.Empty{})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
		gm.Expect(k.canceler.context.Done()).Should(gm.BeClosed())
	})
})

var _ = gk.Describe("RPC:List", func() {
	var k *kubeProvider
	var stream *listResponseStream

	gk.JustBeforeEach(func() {
		configMap1 := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-1",
				Namespace: "ns-1",
				Labels: map[string]string{
					"env": "prod",
				},
			},
		}
		configMap2 := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-2",
				Namespace: "ns-2",
				Labels: map[string]string{
					"env": "dev",
				},
			},
		}

		k = pctx.NewProvider(WithObjects(configMap1, configMap2))
		clientSet, _, err := k.makeClient(context.Background(), nil)
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
		k.clientSet = clientSet

		stream = &listResponseStream{ctx: context.Background()}
	})

	gk.It("lists resources for a token", func() {
		err := k.List(&pulumirpc.ListRequest{
			Token: "kubernetes:core/v1:ConfigMap",
		}, stream)
		gm.Expect(err).ShouldNot(gm.HaveOccurred())

		var ids []string
		for _, response := range stream.responses {
			if result := response.GetResult(); result != nil {
				ids = append(ids, result.GetId())
			}
		}
		gm.Expect(ids).Should(gm.ConsistOf("ns-1/cm-1", "ns-2/cm-2"))
	})

	gk.It("filters list requests by query.namespace", func() {
		query, err := structpb.NewStruct(map[string]any{
			"namespace": "ns-1",
		})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())

		err = k.List(&pulumirpc.ListRequest{
			Token: "kubernetes:core/v1:ConfigMap",
			Query: query,
		}, stream)
		gm.Expect(err).ShouldNot(gm.HaveOccurred())

		var ids []string
		for _, response := range stream.responses {
			if result := response.GetResult(); result != nil {
				ids = append(ids, result.GetId())
			}
		}
		gm.Expect(ids).Should(gm.Equal([]string{"ns-1/cm-1"}))
	})

	gk.It("returns an error for unknown package tokens", func() {
		err := k.List(&pulumirpc.ListRequest{
			Token: "aws:s3/bucket:Bucket",
		}, stream)
		gm.Expect(err).Should(gm.HaveOccurred())
	})
})

var _ = gk.Describe("RPC:GetPluginInfo", func() {
	var k *kubeProvider

	gk.JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	gk.It("should return plugin info", func() {
		resp, err := k.GetPluginInfo(context.Background(), &pbempty.Empty{})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
		gm.Expect(resp.Version).Should(gm.Equal(testPluginVersion))
	})
})

var _ = gk.Describe("RPC:GetSchema", func() {
	var k *kubeProvider
	var req *pulumirpc.GetSchemaRequest

	gk.BeforeEach(func() {
		// initialize the GetSchemaRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.GetSchemaRequest{}
	})

	gk.JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	gk.Context("when the requested version is 0", func() {
		gk.BeforeEach(func() {
			req.Version = 0
		})

		gk.It("should return Pulumi schema info", func() {
			resp, err := k.GetSchema(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(resp.Schema).Should(gm.Equal(testPulumiSchema))
		})
	})

	gk.Context("when the requested version is not 0", func() {
		gk.BeforeEach(func() {
			req.Version = 1
		})

		gk.It("should fail with an invalid schema version", func() {
			_, err := k.GetSchema(context.Background(), req)
			gm.Expect(err).Should(gm.HaveOccurred())
		})
	})
})

var _ = gk.Describe("RPC:GetMapping", func() {
	var k *kubeProvider
	var req *pulumirpc.GetMappingRequest

	gk.BeforeEach(func() {
		// initialize the GetMappingRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.GetMappingRequest{}
	})

	gk.JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	gk.Context("when the requested mapping is for terraform", func() {
		gk.BeforeEach(func() {
			req.Key = "terraform"
		})

		gk.It("should return terraform mapping info", func() {
			resp, err := k.GetMapping(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(resp.Provider).Should(gm.Equal("kubernetes"))
			gm.Expect(resp.Data).Should(gm.Equal([]byte(testTerraformMapping)))
		})
	})

	gk.Context("when the requested mapping is not for terraform", func() {
		gk.BeforeEach(func() {
			req.Key = "foo"
		})

		gk.It("should return empty mapping info", func() {
			resp, err := k.GetMapping(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(resp.Provider).Should(gm.BeEmpty())
			gm.Expect(resp.Data).Should(gm.BeEmpty())
		})
	})
})
