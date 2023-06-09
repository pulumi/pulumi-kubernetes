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
	"os"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/clients"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// kustomizeDirectory takes a path to a kustomization directory, either a local directory or a folder in a git repo,
// and then returns a slice of untyped structs that can be marshalled into Pulumi RPC calls.
func kustomizeDirectory(directory string, clientSet *clients.DynamicClientSet) ([]any, error) {
	path := directory

	// If provided directory doesn't exist locally, assume it's a git repo link.
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		var err error

		// Create a temp dir.
		var temp string
		if temp, err = os.MkdirTemp("", "kustomize-"); err != nil {
			return nil, errors.Wrap(err, "failed to create temp directory for remote kustomize directory")
		}
		defer contract.IgnoreError(os.RemoveAll(temp))

		path, err = workspace.RetrieveGitFolder(directory, temp)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve specified kustomize directory: %q", directory)
		}
	}

	fSys := filesys.MakeFsOnDisk()
	opts := krusty.MakeDefaultOptions()
	opts.Reorder = krusty.ReorderOptionLegacy

	k := krusty.MakeKustomizer(opts)

	rm, err := k.Run(fSys, path)
	if err != nil {
		return nil, errors.Wrapf(err, "kustomize failed for directory %q", path)
	}

	yamlBytes, err := rm.AsYaml()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert kustomize result to YAML")
	}

	return decodeYaml(string(yamlBytes), "", clientSet)
}
