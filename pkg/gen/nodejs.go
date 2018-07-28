// Copyright 2016-2018, Pulumi Corporation.
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

package gen

import (
	"fmt"

	"github.com/cbroglie/mustache"
	providerVersion "github.com/pulumi/pulumi-kubernetes/pkg/version"
)

// --------------------------------------------------------------------------

// Main interface.

// --------------------------------------------------------------------------

// NodeJSClient will generate a Pulumi Kubernetes provider client SDK for nodejs.
func NodeJSClient(
	swagger map[string]interface{}, templateDir string,
) (inputsts, outputsts, providerts, helmts, packagejson string, err error) {
	definitions := swagger["definitions"].(map[string]interface{})

	groupsSlice := createGroups(definitions, inputsAPI)
	inputsts, err = mustache.RenderFile(fmt.Sprintf("%s/typesInput.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return "", "", "", "", "", err
	}

	groupsSlice = createGroups(definitions, outputsAPI)
	outputsts, err = mustache.RenderFile(fmt.Sprintf("%s/typesOutput.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return "", "", "", "", "", err
	}

	groupsSlice = createGroups(definitions, provider)
	providerts, err = mustache.RenderFile(fmt.Sprintf("%s/provider.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return "", "", "", "", "", err
	}

	helmts, err = mustache.RenderFile(fmt.Sprintf("%s/helm.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return "", "", "", "", "", err
	}

	packagejson, err = mustache.RenderFile(fmt.Sprintf("%s/package.json.mustache", templateDir),
		map[string]interface{}{
			"ProviderVersion": providerVersion.Version,
		})
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputsts, outputsts, providerts, helmts, packagejson, nil
}
