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

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

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
