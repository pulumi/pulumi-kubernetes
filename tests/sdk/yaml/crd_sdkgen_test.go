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

package test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

const crdTestdataDir = "testdata/crds"

type crdManifest struct {
	File   string `json:"file"`
	Source string `json:"source"`
}

type crdTestCase struct {
	Name              string            `json:"name"`
	Version           string            `json:"version"`
	Manifests         []crdManifest     `json:"manifests"`
	ExpectedResources []string          `json:"expectedResources"`
	KnownFailures     map[string]string `json:"knownFailures"`

	dir string
}

func (testCase crdTestCase) extensionArg() string {
	parts := []string{"name=" + testCase.Name}
	for _, manifest := range testCase.Manifests {
		parts = append(parts, "crd-manifest="+filepath.Join(testCase.dir, manifest.File))
	}
	return strings.Join(parts, " ")
}

func loadCRDTestCases(t *testing.T) []crdTestCase {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join(crdTestdataDir, "*", "testcase.yaml"))
	require.NoError(t, err)
	require.NotEmptyf(t, paths, "no test cases found under %s", crdTestdataDir)

	testCases := make([]crdTestCase, 0, len(paths))
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		require.NoError(t, err)

		var testCase crdTestCase
		require.NoErrorf(t, yaml.Unmarshal(raw, &testCase), "parse %s", path)

		testCase.dir, err = filepath.Abs(filepath.Dir(path))
		require.NoError(t, err)
		require.Equalf(t, filepath.Base(testCase.dir), testCase.Name,
			"%s declares name %q but lives in directory %q", path, testCase.Name, filepath.Base(testCase.dir))
		require.NotEmptyf(t, testCase.Manifests, "%s lists no manifests", path)
		require.NotEmptyf(t, testCase.ExpectedResources, "%s lists no expected resources", path)

		for _, manifest := range testCase.Manifests {
			require.FileExistsf(t, filepath.Join(testCase.dir, manifest.File),
				"%s references a manifest that is not vendored", path)
		}

		testCases = append(testCases, testCase)
	}
	return testCases
}

type sdkBuild struct {
	sdkDir   string
	repoRoot string
}

type sdkLanguage struct {
	name    string
	tool    string
	compile func(t *testing.T, build sdkBuild) error
}

var sdkLanguages = []sdkLanguage{
	{name: "nodejs", tool: "npm", compile: compileNodejsSDK},
	{name: "python", tool: "python3", compile: compilePythonSDK},
	{name: "go", tool: "go", compile: compileGoSDK},
	{name: "dotnet", tool: "dotnet", compile: compileDotnetSDK},
	{name: "java", tool: "gradle", compile: compileJavaSDK},
}

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	root, err := filepath.Abs(filepath.Join(cwd, "..", "..", ".."))
	require.NoError(t, err)
	return root
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

// The generated SDK pins the base provider at the version of the provider that
// generated it. A development build is never published, so the compile checks
// resolve the base SDK from this repository instead of from a package registry.
func requireLocalBaseSDK(t *testing.T, path, makeTarget string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Skipf("local base SDK missing at %s; run `make %s` first", path, makeTarget)
	}
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

func compileNodejsSDK(t *testing.T, build sdkBuild) error {
	baseSDK := filepath.Join(build.repoRoot, "sdk", "nodejs", "bin")
	requireLocalBaseSDK(t, baseSDK, "nodejs_sdk")

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
	requireLocalBaseSDK(t, baseSDK, "python_sdk")

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
	requireLocalBaseSDK(t, baseSDK, "dotnet_sdk")

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

	baseSDK := filepath.Join(build.repoRoot, "sdk", "java", "build", "libs")
	requireLocalBaseSDK(t, baseSDK, "java_sdk")
	return runBuildCommand(t, build.sdkDir, "gradle", "--console=plain", "compileJava")
}

func generateExtensionSDK(t *testing.T, testCase crdTestCase, language string) string {
	t.Helper()

	root := repoRoot(t)
	binDir := filepath.Join(root, "bin")
	outDir := t.TempDir()

	test := pulumitest.NewPulumiTest(t, "testdata/crd-sdkgen",
		opttest.SkipInstall(),
		opttest.LocalProviderPath("kubernetes", binDir))

	pulumi := test.CurrentStack().Workspace().PulumiCommand()
	_, stderr, _, err := pulumi.Run(test.Context(), test.WorkingDir(), nil, nil, nil, nil,
		"package", "gen-sdk", filepath.Join(binDir, "pulumi-resource-kubernetes"),
		"--extension", testCase.extensionArg(), "--language", language, "-o", outDir)
	require.NoErrorf(t, err, "pulumi package gen-sdk failed: %s", stderr)

	return filepath.Join(outDir, language)
}

func TestCRDExtensionSchemaHasExpectedResources(t *testing.T) {
	for _, testCase := range loadCRDTestCases(t) {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			binDir := filepath.Join(repoRoot(t), "bin")
			test := pulumitest.NewPulumiTest(t, "testdata/crd-sdkgen",
				opttest.SkipInstall(),
				opttest.LocalProviderPath("kubernetes", binDir))

			pulumi := test.CurrentStack().Workspace().PulumiCommand()
			stdout, stderr, _, err := pulumi.Run(test.Context(), test.WorkingDir(), nil, nil, nil, nil,
				"package", "get-schema", filepath.Join(binDir, "pulumi-resource-kubernetes"),
				"--extension", testCase.extensionArg())
			require.NoErrorf(t, err, "pulumi package get-schema failed: %s", stderr)

			var pkg pulumischema.PackageSpec
			require.NoError(t, json.Unmarshal([]byte(stdout), &pkg))

			for _, token := range testCase.ExpectedResources {
				require.Contains(t, pkg.Resources, token)
			}
		})
	}
}

func TestCRDExtensionSDKCompiles(t *testing.T) {
	for _, testCase := range loadCRDTestCases(t) {
		for _, language := range sdkLanguages {
			t.Run(testCase.Name+"/"+language.name, func(t *testing.T) {
				t.Parallel()

				if _, err := exec.LookPath(language.tool); err != nil {
					t.Skipf("%s is not installed", language.tool)
				}

				build := sdkBuild{
					sdkDir:   generateExtensionSDK(t, testCase, language.name),
					repoRoot: repoRoot(t),
				}

				err := language.compile(t, build)
				if reason, known := testCase.KnownFailures[language.name]; known {
					require.Errorf(t, err,
						"the %s SDK compiles now; drop the knownFailures entry from %s/testcase.yaml (%s)",
						language.name, testCase.Name, reason)
					t.Logf("failed as expected: %s", reason)
					return
				}
				require.NoError(t, err)
			})
		}
	}
}
