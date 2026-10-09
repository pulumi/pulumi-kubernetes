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

// Package crdsdk drives extension SDK generation over the vendored
// CustomResourceDefinition manifests under tests/testdata/crds. Each language
// directory under tests/sdk runs the compile check for its own language, so
// that the base SDK and toolchain the check needs are the ones the job already
// installed.
package crdsdk

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

type Case struct {
	Name              string
	Manifests         []string
	ExpectedResources []string

	// Languages whose generated SDK is known not to compile, mapped to the
	// reason. A listed language must fail: when one starts to compile, the
	// test fails and says to drop the entry.
	KnownFailures map[string]string
}

// Codegen emits the same three broken SDKs for every CustomResourceDefinition,
// so every test case shares these.
var codegenDefects = map[string]string{
	"go": "generated go.mod omits the base SDK module, so the module does not build on its own",
	"dotnet": "TODO #pulumi/pulumi-kubernetes/4627: generated csproj omits the " +
		"Pulumi.Kubernetes package reference",
	"java": "TODO #pulumi/pulumi-kubernetes/4628: codegen puts the schema in one " +
		"string constant, over the class-file limit",
}

// Cases are the CustomResourceDefinitions of widely used operators. To move one
// to a new release, download the manifests from the URL above it again.
var Cases = []Case{
	{
		// cert-manager v1.19.1
		// https://github.com/cert-manager/cert-manager/releases/download/v1.19.1/cert-manager.crds.yaml
		Name:      "cert-manager",
		Manifests: []string{"crds.yaml"},
		ExpectedResources: []string{
			"cert-manager:cert-manager.io/v1:Certificate",
			"cert-manager:cert-manager.io/v1:ClusterIssuer",
			"cert-manager:cert-manager.io/v1:Issuer",
			"cert-manager:acme.cert-manager.io/v1:Challenge",
			"cert-manager:acme.cert-manager.io/v1:Order",
		},
		KnownFailures: codegenDefects,
	},
	{
		// Prometheus Operator v0.87.0
		// https://github.com/prometheus-operator/prometheus-operator/releases/download/v0.87.0/stripped-down-crds.yaml
		Name:      "prometheus-operator",
		Manifests: []string{"crds.yaml"},
		ExpectedResources: []string{
			"prometheus-operator:monitoring.coreos.com/v1:Alertmanager",
			"prometheus-operator:monitoring.coreos.com/v1:Prometheus",
			"prometheus-operator:monitoring.coreos.com/v1:PrometheusRule",
			"prometheus-operator:monitoring.coreos.com/v1:ServiceMonitor",
		},
		KnownFailures: codegenDefects,
	},
	{
		// Argo CD v3.2.0
		// https://raw.githubusercontent.com/argoproj/argo-cd/v3.2.0/manifests/crds/application-crd.yaml
		// https://raw.githubusercontent.com/argoproj/argo-cd/v3.2.0/manifests/crds/appproject-crd.yaml
		// https://raw.githubusercontent.com/argoproj/argo-cd/v3.2.0/manifests/crds/applicationset-crd.yaml
		Name:      "argocd",
		Manifests: []string{"application.yaml", "appproject.yaml", "applicationset.yaml"},
		ExpectedResources: []string{
			"argocd:argoproj.io/v1alpha1:Application",
			"argocd:argoproj.io/v1alpha1:ApplicationSet",
			"argocd:argoproj.io/v1alpha1:AppProject",
		},
		KnownFailures: codegenDefects,
	},
}

// RepoRoot resolves the repository root from any test directory under tests/sdk.
func RepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	root, err := filepath.Abs(filepath.Join(cwd, "..", "..", ".."))
	require.NoError(t, err)
	return root
}

// HostProject is a Pulumi project that declares no resources. The package
// commands need a project to run in.
func HostProject(t *testing.T) string {
	t.Helper()
	return filepath.Join(RepoRoot(t), "tests", "testdata", "crd-sdkgen")
}

func (testCase Case) ExtensionArg(t *testing.T) string {
	t.Helper()

	crdDir := filepath.Join(RepoRoot(t), "tests", "testdata", "crds", testCase.Name)
	parts := []string{"name=" + testCase.Name}
	for _, manifest := range testCase.Manifests {
		path := filepath.Join(crdDir, manifest)
		require.FileExists(t, path)
		parts = append(parts, "crd-manifest="+path)
	}
	return strings.Join(parts, " ")
}

// Schema returns the extension schema the provider serves for a test case.
func (testCase Case) Schema(t *testing.T) pulumischema.PackageSpec {
	t.Helper()

	binDir := filepath.Join(RepoRoot(t), "bin")
	test := pulumitest.NewPulumiTest(t, HostProject(t),
		opttest.SkipInstall(),
		opttest.LocalProviderPath("kubernetes", binDir))

	pulumi := test.CurrentStack().Workspace().PulumiCommand()
	stdout, stderr, _, err := pulumi.Run(test.Context(), test.WorkingDir(), nil, nil, nil, nil,
		"package", "get-schema", filepath.Join(binDir, "pulumi-resource-kubernetes"),
		"--extension", testCase.ExtensionArg(t))
	require.NoErrorf(t, err, "pulumi package get-schema failed: %s", stderr)

	var pkg pulumischema.PackageSpec
	require.NoError(t, json.Unmarshal([]byte(stdout), &pkg))
	return pkg
}

func (testCase Case) generateSDK(t *testing.T, language string) string {
	t.Helper()

	binDir := filepath.Join(RepoRoot(t), "bin")
	outDir := t.TempDir()

	test := pulumitest.NewPulumiTest(t, HostProject(t),
		opttest.SkipInstall(),
		opttest.LocalProviderPath("kubernetes", binDir))

	pulumi := test.CurrentStack().Workspace().PulumiCommand()
	_, stderr, _, err := pulumi.Run(test.Context(), test.WorkingDir(), nil, nil, nil, nil,
		"package", "gen-sdk", filepath.Join(binDir, "pulumi-resource-kubernetes"),
		"--extension", testCase.ExtensionArg(t), "--language", language, "-o", outDir)
	require.NoErrorf(t, err, "pulumi package gen-sdk failed: %s", stderr)

	return filepath.Join(outDir, language)
}

var compilers = map[string]func(t *testing.T, build sdkBuild) error{
	"nodejs": compileNodejsSDK,
	"python": compilePythonSDK,
	"go":     compileGoSDK,
	"dotnet": compileDotnetSDK,
	"java":   compileJavaSDK,
}

// RunCompileTest generates the extension SDK for every test case in one
// language and compiles it against this repository's own base SDK.
func RunCompileTest(t *testing.T, language string) {
	compile, supported := compilers[language]
	require.Truef(t, supported, "no compile check for %s", language)

	for _, testCase := range Cases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			build := sdkBuild{
				sdkDir:   testCase.generateSDK(t, language),
				repoRoot: RepoRoot(t),
			}

			err := compile(t, build)
			if reason, known := testCase.KnownFailures[language]; known {
				require.Errorf(t, err,
					"the %s SDK compiles now; drop the %s knownFailures entry for %s (%s)",
					language, language, testCase.Name, reason)
				t.Logf("failed as expected: %s", reason)
				return
			}
			require.NoError(t, err)
		})
	}
}

type sdkBuild struct {
	sdkDir   string
	repoRoot string
}

func runBuildCommand(t *testing.T, dir, name string, args ...string) error {
	t.Helper()

	command := exec.Command(name, args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	t.Logf("$ %s %s\n%s", name, strings.Join(args, " "), output)
	if err != nil {
		return fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, output)
	}
	return nil
}

func onlyMatch(t *testing.T, dir, pattern string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	require.NoError(t, err)
	require.Lenf(t, matches, 1, "expected exactly one %s in %s, found %v", pattern, dir, matches)
	return matches[0]
}

func compileGoSDK(t *testing.T, build sdkBuild) error {
	const baseModule = "github.com/pulumi/pulumi-kubernetes/sdk/v4"

	goMod, err := os.ReadFile(filepath.Join(build.sdkDir, "go.mod"))
	if err != nil {
		return err
	}
	if !strings.Contains(string(goMod), baseModule) {
		return fmt.Errorf("generated go.mod does not require %s:\n%s", baseModule, goMod)
	}

	replacement := baseModule + "=" + filepath.Join(build.repoRoot, "sdk")
	if err := runBuildCommand(t, build.sdkDir, "go", "mod", "edit", "-replace", replacement); err != nil {
		return err
	}
	if err := runBuildCommand(t, build.sdkDir, "go", "mod", "tidy"); err != nil {
		return err
	}
	return runBuildCommand(t, build.sdkDir, "go", "build", "./...")
}

// The generated SDK pins the base provider at the version of the provider that
// generated it. A development build is never published, so each compile check
// resolves the base SDK from this repository instead of from a package registry.
func compileNodejsSDK(t *testing.T, build sdkBuild) error {
	baseSDK := filepath.Join(build.repoRoot, "sdk", "nodejs", "bin")

	manifestPath := filepath.Join(build.sdkDir, "package.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	var manifest map[string]any
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return err
	}
	dependencies, ok := manifest["dependencies"].(map[string]any)
	if !ok {
		return fmt.Errorf("%s declares no dependencies", manifestPath)
	}
	if _, declared := dependencies["@pulumi/kubernetes"]; !declared {
		return fmt.Errorf("%s does not depend on @pulumi/kubernetes", manifestPath)
	}
	dependencies["@pulumi/kubernetes"] = "file:" + baseSDK
	patched, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(manifestPath, patched, 0o600); err != nil {
		return err
	}

	if err := runBuildCommand(t, build.sdkDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
		return err
	}
	return runBuildCommand(t, build.sdkDir, "npx", "--no-install", "tsc", "--noEmit")
}

func compilePythonSDK(t *testing.T, build sdkBuild) error {
	baseSDK := filepath.Join(build.repoRoot, "sdk", "python", "bin")

	// Resolve the module before pip runs, because installing adds an egg-info
	// directory that matches the same pattern.
	module := filepath.Base(onlyMatch(t, build.sdkDir, "pulumi_*"))

	venv := filepath.Join(t.TempDir(), "venv")
	if err := runBuildCommand(t, build.sdkDir, "python3", "-m", "venv", venv); err != nil {
		return err
	}
	pip := filepath.Join(venv, "bin", "pip")
	if err := runBuildCommand(t, build.sdkDir, pip, "install", "--quiet", baseSDK); err != nil {
		return err
	}
	if err := runBuildCommand(t, build.sdkDir, pip, "install", "--quiet", "--no-deps", "."); err != nil {
		return err
	}

	python := filepath.Join(venv, "bin", "python")
	return runBuildCommand(t, build.sdkDir, python, "-c", "import "+module)
}

func compileDotnetSDK(t *testing.T, build sdkBuild) error {
	project := onlyMatch(t, build.sdkDir, "*.csproj")
	projectFile, err := os.ReadFile(project)
	if err != nil {
		return err
	}
	if !strings.Contains(string(projectFile), "Pulumi.Kubernetes") {
		return fmt.Errorf("%s has no Pulumi.Kubernetes package reference:\n%s",
			filepath.Base(project), projectFile)
	}

	baseSDK := filepath.Join(build.repoRoot, "sdk", "dotnet", "bin", "Debug")

	// The local base SDK keeps the same version across rebuilds, so a shared
	// package folder serves a stale copy. Give each run its own.
	nugetConfig := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <config>
    <add key="globalPackagesFolder" value="%s" />
  </config>
  <packageSources>
    <add key="local-base-sdk" value="%s" />
  </packageSources>
</configuration>
`, filepath.Join(t.TempDir(), "nuget"), baseSDK)
	if err := os.WriteFile(filepath.Join(build.sdkDir, "nuget.config"), []byte(nugetConfig), 0o600); err != nil {
		return err
	}

	return runBuildCommand(t, build.sdkDir, "dotnet", "build", project)
}

// A class file cannot hold a string constant longer than this. Codegen puts
// each constant on one line, so an over-long line is an over-long constant.
const javaConstantLimit = 65535

func checkJavaConstantLength(sdkDir string) error {
	var sources []string
	err := filepath.WalkDir(sdkDir, func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && filepath.Ext(path) == ".java" {
			sources = append(sources, path)
		}
		return err
	})
	if err != nil {
		return err
	}

	for _, path := range sources {
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for number, line := range strings.Split(string(source), "\n") {
			if len(line) > javaConstantLimit {
				relative, _ := filepath.Rel(sdkDir, path)
				return fmt.Errorf("%s:%d is %d bytes, over the %d byte constant limit",
					relative, number+1, len(line), javaConstantLimit)
			}
		}
	}
	return nil
}

func compileJavaSDK(t *testing.T, build sdkBuild) error {
	if err := checkJavaConstantLength(build.sdkDir); err != nil {
		return err
	}
	return runBuildCommand(t, build.sdkDir, "gradle", "--console=plain", "compileJava")
}
