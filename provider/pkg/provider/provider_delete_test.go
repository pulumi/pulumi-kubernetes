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

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = gk.Describe("RPC:Delete", func() {
	var k *kubeProvider
	var req *pulumirpc.DeleteRequest

	gk.BeforeEach(func() {
		k = pctx.NewProvider()
		req = &pulumirpc.DeleteRequest{}
	})

	gk.Context("when the cluster is unreachable", func() {
		gk.BeforeEach(func() {
			k.clusterUnreachable = true
			k.clusterUnreachableReason = "unable to load kubeconfig"
		})

		gk.Describe("Helm Release", func() {
			gk.BeforeEach(func() {
				req.Urn = "urn:pulumi:test::test::kubernetes:helm.sh/v3:Release::test"
				req.Id = "default/test"
			})

			gk.Context("when deleteUnreachable is disabled", func() {
				gk.BeforeEach(func() {
					k.deleteUnreachable = false
				})
				gk.It("should return an error mentioning PULUMI_K8S_DELETE_UNREACHABLE", func() {
					_, err := k.Delete(context.Background(), req)
					gm.Expect(err).Should(gm.MatchError(gm.ContainSubstring("PULUMI_K8S_DELETE_UNREACHABLE")))
				})
			})

			gk.Context("when deleteUnreachable is enabled", func() {
				gk.BeforeEach(func() {
					k.deleteUnreachable = true
				})
				gk.It("should remove the resource from state", func() {
					_, err := k.Delete(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
				})
			})
		})
	})
})
