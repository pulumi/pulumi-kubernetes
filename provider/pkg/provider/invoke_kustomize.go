// Copyright 2016-2022, Pulumi Corporation.
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
	"fmt"
	"os"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// kustomizeDirectory takes a path to a kustomization directory, either a local directory or a folder in a git repo,
// and then returns a slice of untyped structs that can be marshalled into Pulumi RPC calls.
func kustomizeDirectory(ctx context.Context, directory string, clientSet *clients.DynamicClientSet) ([]any, error) {
	path := directory

	// If provided directory doesn't exist locally, assume it's a git repo link.
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		var err error

		// Create a temp dir.
		var temp string
		if temp, err = os.MkdirTemp("", "kustomize-"); err != nil {
			return nil, fmt.Errorf("failed to create temp directory for remote kustomize directory: %w", err)
		}
		defer contract.IgnoreError(os.RemoveAll(temp))

		path, err = workspace.RetrieveGitFolder(ctx, directory, temp)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve specified kustomize directory: %q: %w", directory, err)
		}
	}

	fSys := filesys.MakeFsOnDisk()
	opts := krusty.MakeDefaultOptions()
	opts.Reorder = krusty.ReorderOptionLegacy

	// TODO: kustomize helmChart support is currently enabled via an undocumented feature flag.
	//       See https://github.com/pulumi/pulumi-kubernetes/issues/2470 for additional details.
	enableHelmChartSupport := false
	helmPath := "helm" // TODO: support this as a parameter to kustomize.Directory; this won't work for Windows
	if v, ok := os.LookupEnv("PULUMI_K8S_KUSTOMIZE_HELM"); ok && cmdutil.IsTruthy(v) {
		enableHelmChartSupport = true
	}
	// Add support for helmCharts plugin
	// See https://github.com/kubernetes-sigs/kustomize/blob/v3.3.1/examples/chart.md for more details.
	if enableHelmChartSupport {
		opts.PluginConfig = types.EnabledPluginConfig(types.BploUseStaticallyLinked)
		opts.PluginConfig.HelmConfig.Command = helmPath
	}

	k := krusty.MakeKustomizer(opts)

	rm, err := k.Run(fSys, path)
	if err != nil {
		if enableHelmChartSupport && strings.Contains(err.Error(), `(is 'helm' installed?)`) {
			err = fmt.Errorf("the helmCharts feature requires %q binary to be on the system PATH", helmPath)
		}
		return nil, fmt.Errorf("kustomize failed for directory %q: %w", path, err)
	}

	yamlBytes, err := rm.AsYaml()
	if err != nil {
		return nil, fmt.Errorf("failed to convert kustomize result to YAML: %w", err)
	}

	return decodeYaml(string(yamlBytes), "", clientSet)
}
