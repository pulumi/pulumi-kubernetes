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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = Describe("RPC:Cancel", func() {
	var k *kubeProvider

	JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	It("should cancel any future operations", func() {
		_, err := k.Cancel(context.Background(), &pbempty.Empty{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(k.canceler.context.Done()).Should(BeClosed())
	})
})

var _ = Describe("RPC:GetPluginInfo", func() {
	var k *kubeProvider

	JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	It("should return plugin info", func() {
		resp, err := k.GetPluginInfo(context.Background(), &pbempty.Empty{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.Version).Should(Equal(testPluginVersion))
	})
})

var _ = Describe("RPC:GetSchema", func() {
	var k *kubeProvider
	var req *pulumirpc.GetSchemaRequest

	BeforeEach(func() {
		// initialize the GetSchemaRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.GetSchemaRequest{}
	})

	JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	Context("when the requested version is 0", func() {
		BeforeEach(func() {
			req.Version = 0
		})

		It("should return Pulumi schema info", func() {
			resp, err := k.GetSchema(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.Schema).Should(Equal(testPulumiSchema))
		})
	})

	Context("when the requested version is not 0", func() {
		BeforeEach(func() {
			req.Version = 1
		})

		It("should fail with an invalid schema version", func() {
			_, err := k.GetSchema(context.Background(), req)
			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("RPC:GetMapping", func() {
	var k *kubeProvider
	var req *pulumirpc.GetMappingRequest

	BeforeEach(func() {
		// initialize the GetMappingRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.GetMappingRequest{}
	})

	JustBeforeEach(func() {
		k = pctx.NewProvider()
	})

	Context("when the requested mapping is for terraform", func() {
		BeforeEach(func() {
			req.Key = "terraform"
		})

		It("should return terraform mapping info", func() {
			resp, err := k.GetMapping(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.Provider).Should(Equal("kubernetes"))
			Expect(resp.Data).Should(Equal([]byte(testTerraformMapping)))
		})
	})

	Context("when the requested mapping is not for terraform", func() {
		BeforeEach(func() {
			req.Key = "foo"
		})

		It("should return empty mapping info", func() {
			resp, err := k.GetMapping(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.Provider).Should(BeEmpty())
			Expect(resp.Data).Should(BeEmpty())
		})
	})
})
