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

// nolint: lll
package gen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v2"
)

// TestCreateGroups_IdentifyListKinds loads txtar files under testdata/identify-list-kinds and uses them to
// craft an unstructured map[string]any definitions file for createGroups.
// The goal of this test is to ensure we can accurately distinuguish between singletons and lists of kinds.
//
// The test files should contain the following files:
// - definitions: a JSON file containing the definitions
// - kinds: a YAML file containing a list of kinds that are singletons
// - listKinds: a YAML file containing a list of kinds that are lists/collections of singleton kinds
func TestCreateGroups_IdentifyListKinds(t *testing.T) {
	dir := filepath.Join("testdata/identify-list-kinds")
	tests, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			archive, err := txtar.ParseFile(filepath.Join(dir, tt.Name()))
			require.NoError(t, err)

			var definitions map[string]any
			var kinds, listKinds map[string]struct{}

			for _, f := range archive.Files {
				var parsed []string
				switch f.Name {
				case "definitions":
					err := json.Unmarshal(f.Data, &definitions)
					require.NoError(t, err, f.Name)
				case "kinds":
					err = yaml.Unmarshal(f.Data, &parsed)
					require.NoError(t, err, f.Name)
					kinds = sliceToSet(parsed)
				case "listKinds":
					err = yaml.Unmarshal(f.Data, &parsed)
					require.NoError(t, err, f.Name)
					listKinds = sliceToSet(parsed)
				default:
					t.Fatal("unrecognized filename", f.Name)
				}
			}

			configGroups := createGroups(definitions, true)

			// Loop through all parsed kinds and ensure they are accounted for.
			for _, g := range configGroups {
				for _, v := range g.versions {
					for _, kind := range v.kinds {
						gvk := gvkToString(kind.gvk.Group, kind.gvk.Version, kind.gvk.Kind)
						if kind.isList {
							delete(listKinds, gvk)
						} else {
							delete(kinds, gvk)
						}
					}
				}
			}

			assert.Equal(t, 0, len(kinds), "kinds not found while parsing: %v", kinds)
			assert.Equal(t, 0, len(listKinds), "listKinds not found while parsing: %v", listKinds)
		})
	}
}

func gvkToString(group, version, kind string) string {
	return group + "." + version + "." + kind
}

func sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}
