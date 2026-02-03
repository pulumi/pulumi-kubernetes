// Copyright 2016-2023, Pulumi Corporation.
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

// nolint:govet
// nolint:goconst
package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pulumi/providertest/grpclog"
	"github.com/pulumi/providertest/pulumitest"
	. "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// TestHelmUnknowns tests the handling of unknowns in the Helm provider.
// Test steps:
// 1. Preview a program that has computed inputs; expected computed outputs.
// 2. Deploy a program; expected real outputs.
// 3. Preview an update involving a change to the release name; expect replacement.
func TestHelmUnknowns(t *testing.T) {
	g := NewWithT(t)

	// Copy test_dir to temp directory, install deps and create "my-stack"
	test := pulumitest.NewPulumiTest(
		t,
		"helm-release-unknowns", /*opttest.LocalProviderPath("kubernetes", abs(t, "../../../bin"))*/
	)
	t.Logf("into %s", test.WorkingDir())

	urn := func(baseType tokens.Type, name string) string {
		return string(resource.NewURN("test", "helm-release-unknowns", "", baseType, name))
	}

	previewF := func(opts ...optpreview.Option) auto.PreviewResult {
		clearGrpcLog(t, test)
		preview := test.Preview(t, opts...)
		t.Log(preview.StdOut)
		return preview
	}
	upF := func() auto.UpResult {
		clearGrpcLog(t, test)
		up := test.Up(t)
		t.Log(up.StdOut)
		return up
	}

	lookup := func() grpclog.TypedEntry[rpc.CreateRequest, rpc.CreateResponse] {
		creates, err := test.GrpcLog(t).Creates()
		g.Expect(err).ToNot(HaveOccurred())
		release := findByUrn(t, creates, urn("kubernetes:helm.sh/v3:Release", "release"))
		g.Expect(release).ToNot(BeNil())
		logEntry(t, *release)
		return *release
	}

	// 1. Preview and then deploy a program that has computed inputs; expected computed outputs.
	previewF()
	release1 := lookup()
	outputs1 := unmarshalProperties(t, release1.Response.Properties)
	g.Expect(release1.Response.Id).To(BeEmpty(), "Previews should return empty IDs")
	g.Expect(outputs1).To(MatchProps(IgnoreExtras, Props{
		"name":          BeComputed(),
		"resourceNames": BeComputed(),
		"status":        BeComputed(),
		"values": MatchObject(IgnoreExtras, Props{
			"global": MatchObject(IgnoreExtras, Props{
				"redis": MatchObject(IgnoreExtras, Props{
					"password": BeComputed(),
				}),
			}),
		}),
	}))

	// 2. Deploy a program; expected real outputs.
	upF()
	release2 := lookup()
	outputs2 := unmarshalProperties(t, release2.Response.Properties)
	g.Expect(release2.Response.Id).To(MatchRegexp(`release-ns-\w+\/\w{8}`), "Ups should return proper IDs")
	g.Expect(outputs2).To(MatchProps(IgnoreExtras, Props{
		"name": MatchValue(MatchRegexp(`\w{8}`)),
		"resourceNames": MatchObject(IgnoreExtras, Props{
			"Service/v1": BeArray(),
		}),
		"status": MatchObject(IgnoreExtras, Props{
			"status": MatchValue(Equal("deployed")),
		}),
		"values": MatchObject(IgnoreExtras, Props{
			"global": MatchObject(IgnoreExtras, Props{
				"redis": MatchObject(IgnoreExtras, Props{
					"password": MatchSecret(MatchValue(Not(BeEmpty()))),
				}),
			}),
		}),
	}))

	// 3. Preview an update involving a change to the release name; expect replacement.
	previewF(optpreview.Replace([]string{urn("random:index/randomString:RandomString", "name")}))
	release3 := lookup()
	output3 := unmarshalProperties(t, release3.Response.Properties)
	g.Expect(release3.Response.Id).To(BeEmpty(), "Previews should return empty IDs")
	g.Expect(output3).To(MatchProps(IgnoreExtras, Props{
		"name":          BeComputed(),
		"resourceNames": BeComputed(),
		"status":        BeComputed(),
		"values": MatchObject(IgnoreExtras, Props{
			"global": MatchObject(IgnoreExtras, Props{
				"redis": MatchObject(IgnoreExtras, Props{
					"password": MatchSecret(MatchValue(Not(BeEmpty()))),
				}),
			}),
		}),
	}))
}

type hasURN interface {
	GetUrn() string
}

func unmarshalProperties(t *testing.T, props *structpb.Struct) resource.PropertyMap {
	pm, err := plugin.UnmarshalProperties(props, plugin.MarshalOptions{
		KeepUnknowns: true,
		KeepSecrets:  true,
	})
	if err != nil {
		t.Fatalf("failed to unmarshal properties: %s", err)
	}
	return pm
}

func findByUrn[TRequest any, TResponse any](_ *testing.T, entries []grpclog.TypedEntry[TRequest, TResponse],
	urn string) *grpclog.TypedEntry[TRequest, TResponse] {
	for _, e := range entries {
		var eI any = &e.Request
		if hasUrn, ok := eI.(hasURN); ok {
			if hasUrn.GetUrn() == urn {
				return &e
			}
		}
	}
	return nil
}

func logEntry[TRequest any, TResponse any](t *testing.T, entries ...grpclog.TypedEntry[TRequest, TResponse]) {
	for _, e := range entries {
		var req any = &e.Request
		if m, ok := req.(proto.Message); ok {
			t.Log(protojson.Format(m))
		}
		var resp any = &e.Response
		if m, ok := resp.(proto.Message); ok {
			t.Log(protojson.Format(m))
		}
	}
}

func clearGrpcLog(t *testing.T, pt *pulumitest.PulumiTest) {
	env := pt.CurrentStack().Workspace().GetEnvVars()
	if env == nil || env["PULUMI_DEBUG_GRPC"] == "" {
		t.Log("can't clear gRPC log: PULUMI_DEBUG_GRPC env var not set")
		return
	}
	logPath := env["PULUMI_DEBUG_GRPC"]
	if _, err := os.Stat(logPath); err == nil {
		if err := os.Truncate(logPath, 0); err != nil {
			t.Fatalf("failed to clear gRPC log: %s", err)
		}
	}
}

func TestPreviewWithUnreachableCluster(t *testing.T) {
	t.Parallel()

	test := pulumitest.NewPulumiTest(t, "helm-preview-unreachable")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	test.Preview(t)
}

func TestHelmChartV4WithYamlRenderModeRendersWithoutClusterConnection(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// Create a temporary directory to hold rendered YAML manifests.
	dir, err := os.MkdirTemp("", "helm-chart-v4-render-yaml-test")
	g.Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(dir)

	test := pulumitest.NewPulumiTest(t, "helm-chart-v4-render-yaml")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	// Set the config value for renderDir
	err = test.CurrentStack().SetConfig(context.Background(), "renderDir", auto.ConfigValue{
		Value:  dir,
		Secret: false,
	})
	g.Expect(err).ToNot(HaveOccurred())

	preview := test.Preview(t)
	t.Log(preview.StdOut)

	up := test.Up(t)
	t.Log(up.StdOut)

	// Verify that YAML directory was created and contains files
	files, err := os.ReadDir(dir)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(files)).To(BeNumerically(">", 0), "YAML directory should contain rendered files")

	manifestDir := filepath.Join(dir, "1-manifest")
	if _, err := os.Stat(manifestDir); err == nil {
		manifestFiles, err := os.ReadDir(manifestDir)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(len(manifestFiles)).
			To(BeNumerically(">", 0), "Manifest directory should contain rendered Helm chart resources")
	}
}
