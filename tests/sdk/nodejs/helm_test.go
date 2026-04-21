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

	gm "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/pulumi/providertest/grpclog"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pgm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
)

// TestHelmUnknowns tests the handling of unknowns in the Helm provider.
// Test steps:
// 1. Preview a program that has computed inputs; expected computed outputs.
// 2. Deploy a program; expected real outputs.
// 3. Preview an update involving a change to the release name; expect replacement.
func TestHelmUnknowns(t *testing.T) {
	g := gm.NewWithT(t)

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
		g.Expect(err).ToNot(gm.HaveOccurred())
		release := findByUrn(t, creates, urn("kubernetes:helm.sh/v3:Release", "release"))
		g.Expect(release).ToNot(gm.BeNil())
		logEntry(t, *release)
		return *release
	}

	// 1. Preview and then deploy a program that has computed inputs; expected computed outputs.
	previewF()
	release1 := lookup()
	outputs1 := unmarshalProperties(t, release1.Response.Properties)
	g.Expect(release1.Response.Id).To(gm.BeEmpty(), "Previews should return empty IDs")
	g.Expect(outputs1).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
		"name":          pgm.BeComputed(),
		"resourceNames": pgm.BeComputed(),
		"status":        pgm.BeComputed(),
		"values": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
			"global": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
				"redis": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
					"password": pgm.BeComputed(),
				}),
			}),
		}),
	}))

	// 2. Deploy a program; expected real outputs.
	upF()
	release2 := lookup()
	outputs2 := unmarshalProperties(t, release2.Response.Properties)
	g.Expect(release2.Response.Id).To(gm.MatchRegexp(`release-ns-\w+\/\w{8}`), "Ups should return proper IDs")
	g.Expect(outputs2).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
		"name": pgm.MatchValue(gm.MatchRegexp(`\w{8}`)),
		"resourceNames": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
			"Service/v1": pgm.BeArray(),
		}),
		"status": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
			"status": pgm.MatchValue(gm.Equal("deployed")),
		}),
		"values": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
			"global": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
				"redis": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
					"password": pgm.MatchSecret(pgm.MatchValue(gm.Not(gm.BeEmpty()))),
				}),
			}),
		}),
	}))

	// 3. Preview an update involving a change to the release name; expect replacement.
	previewF(optpreview.Replace([]string{urn("random:index/randomString:RandomString", "name")}))
	release3 := lookup()
	output3 := unmarshalProperties(t, release3.Response.Properties)
	g.Expect(release3.Response.Id).To(gm.BeEmpty(), "Previews should return empty IDs")
	g.Expect(output3).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
		"name":          pgm.BeComputed(),
		"resourceNames": pgm.BeComputed(),
		"status":        pgm.BeComputed(),
		"values": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
			"global": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
				"redis": pgm.MatchObject(gs.IgnoreExtras, pgm.Props{
					"password": pgm.MatchSecret(pgm.MatchValue(gm.Not(gm.BeEmpty()))),
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

// TestHelmNullValues verifies that setting a Helm chart value to null deletes
// the chart's default for that key (https://github.com/pulumi/pulumi-kubernetes/issues/2997).
func TestHelmNullValues(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("helm-release-null-values", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Step 1: both chart defaults should be present in the ConfigMap.
			cm := stackInfo.Outputs["configMapData"].(map[string]any)
			assert.Equal(t, "default-alpha", cm["alpha"])
			assert.Equal(t, "default-beta", cm["beta"])
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("helm-release-null-values", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Step 2: alpha should be deleted by null, beta remains.
					cm := stackInfo.Outputs["configMapData"].(map[string]any)
					assert.NotContains(t, cm, "alpha", "alpha should be deleted by null")
					assert.Equal(t, "default-beta", cm["beta"])
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

// TestNullValues verifies that explicit null values in native Kubernetes resource
// specs survive the provider's Check/Create/Update pipeline (#2997).
func TestNullValues(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("null-values", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Step 1: replicas should be 2.
			replicas := stackInfo.Outputs["replicas"].(float64)
			assert.Equal(t, float64(2), replicas)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("null-values", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Step 2: replicas=null lets server default to 1.
					replicas := stackInfo.Outputs["replicas"].(float64)
					assert.Equal(t, float64(1), replicas)
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
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
	g := gm.NewWithT(t)

	// Create a temporary directory to hold rendered YAML manifests.
	dir, err := os.MkdirTemp("", "helm-chart-v4-render-yaml-test")
	g.Expect(err).ToNot(gm.HaveOccurred())
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
	g.Expect(err).ToNot(gm.HaveOccurred())

	preview := test.Preview(t)
	t.Log(preview.StdOut)

	up := test.Up(t)
	t.Log(up.StdOut)

	// Verify that YAML directory was created and contains files
	files, err := os.ReadDir(dir)
	g.Expect(err).ToNot(gm.HaveOccurred())
	g.Expect(len(files)).To(gm.BeNumerically(">", 0), "YAML directory should contain rendered files")

	manifestDir := filepath.Join(dir, "1-manifest")
	if _, err := os.Stat(manifestDir); err == nil {
		manifestFiles, err := os.ReadDir(manifestDir)
		g.Expect(err).ToNot(gm.HaveOccurred())
		g.Expect(len(manifestFiles)).
			To(gm.BeNumerically(">", 0), "Manifest directory should contain rendered Helm chart resources")
	}
}
