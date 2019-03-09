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

type GroupTS struct {
	Versions map[string]*VersionTS
	Index    string
}

type VersionTS struct {
	Kinds map[string]string
	Index string
}

// NodeJSClient will generate a Pulumi Kubernetes provider client SDK for nodejs.
func NodeJSClient(swagger map[string]interface{}, templateDir string,
) (inputsts, outputsts, providerts, helmts, indexts, packagejson string, groupsts map[string]*GroupTS, err error) {
	definitions := swagger["definitions"].(map[string]interface{})

	groupsSlice := createGroups(definitions, nodeJSInputs())
	inputsts, err = mustache.RenderFile(fmt.Sprintf("%s/typesInput.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return
	}

	groupsSlice = createGroups(definitions, nodeJSOutputs())
	outputsts, err = mustache.RenderFile(fmt.Sprintf("%s/typesOutput.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return
	}

	groupsSlice = createGroups(definitions, nodeJSProvider())
	groupsts = make(map[string]*GroupTS)
	for _, group := range groupsSlice {
		groupTS := &GroupTS{}
		for _, version := range group.Versions() {
			if groupTS.Versions == nil {
				groupTS.Versions = make(map[string]*VersionTS)
			}
			versionTS := &VersionTS{}
			for _, kind := range version.Kinds() {
				if versionTS.Kinds == nil {
					versionTS.Kinds = make(map[string]string)
				}
				kindts, err := mustache.RenderFile(fmt.Sprintf("%s/kind.ts.mustache", templateDir),
					map[string]interface{}{
						"Comment":    kind.Comment(),
						"Group":      group.Group(),
						"Kind":       kind.Kind(),
						"Properties": kind.Properties(),
						"Version":    version.Version(),
					})
				if err != nil {
					return "", "", "", "", "", "", nil, err
				}
				versionTS.Kinds[kind.Kind()] = kindts
			}

			kindIndexTS, err := mustache.RenderFile(fmt.Sprintf("%s/kindIndex.ts.mustache", templateDir),
				map[string]interface{}{
					"Kinds": version.Kinds(),
				})
			if err != nil {
				return "", "", "", "", "", "", nil, err
			}
			versionTS.Index = kindIndexTS
			groupTS.Versions[version.Version()] = versionTS
		}

		versionIndexTS, err := mustache.RenderFile(fmt.Sprintf("%s/versionIndex.ts.mustache", templateDir),
			map[string]interface{}{
				"Versions": group.Versions(),
			})
		if err != nil {
			return "", "", "", "", "", "", nil, err
		}
		groupTS.Index = versionIndexTS
		groupsts[group.Group()] = groupTS
	}

	providerts, err = mustache.RenderFile(fmt.Sprintf("%s/provider.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return
	}

	helmts, err = mustache.RenderFile(fmt.Sprintf("%s/helm.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return
	}

	packagejson, err = mustache.RenderFile(fmt.Sprintf("%s/package.json.mustache", templateDir),
		map[string]interface{}{
			"ProviderVersion": providerVersion.Version,
		})
	if err != nil {
		return
	}

	indexts, err = mustache.RenderFile(fmt.Sprintf("%s/providerIndex.ts.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return
	}

	return inputsts, outputsts, providerts, helmts, indexts, packagejson, groupsts, nil
}
