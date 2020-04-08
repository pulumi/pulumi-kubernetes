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

//nolint: goconst
package gen

import (
	"encoding/json"
	"fmt"
	"strings"

	providerVersion "github.com/pulumi/pulumi-kubernetes/pkg/version"

	pschema "github.com/pulumi/pulumi/pkg/codegen/schema"
	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
)

// PulumiSchema will generate a Pulumi schema for the given k8s schema.
func PulumiSchema(swagger map[string]interface{}) pschema.PackageSpec {
	pkg := pschema.PackageSpec{
		Name:        "kubernetes",
		Version:     providerVersion.Version,
		Description: "A Pulumi package for creating and managing Kubernetes resources.",
		License:     "Apache-2.0",
		Keywords:    []string{"pulumi", "kubernetes"},
		Homepage:    "https://pulumi.com",
		Repository:  "https://github.com/pulumi/pulumi-kubernetes",

		Config: pschema.ConfigSpec{
			Variables: map[string]pschema.PropertySpec{
				"kubeconfig": {
					Description: "The contents of a kubeconfig file. If this is set, this config will be used instead of $KUBECONFIG.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"context": {
					Description: "If present, the name of the kubeconfig context to use.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"cluster": {
					Description: "If present, the name of the kubeconfig cluster to use.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"namespace": {
					Description: "If present, the default namespace to use. This flag is ignored for cluster-scoped resources.\n\nA namespace can be specified in multiple places, and the precedence is as follows:\n1. `.metadata.namespace` set on the resource.\n2. This `namespace` parameter.\n3. `namespace` set for the active context in the kubeconfig.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"enableDryRun": {
					Description: "BETA FEATURE - If present and set to true, enable server-side diff calculations.\nThis feature is in developer preview, and is disabled by default.\n\nThis config can be specified in the following ways, using this precedence:\n1. This `enableDryRun` parameter.\n2. The `PULUMI_K8S_ENABLE_DRY_RUN` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"suppressDeprecationWarnings": {
					Description: "If present and set to true, suppress apiVersion deprecation warnings from the CLI.\n\nThis config can be specified in the following ways, using this precedence:\n1. This `suppressDeprecationWarnings` parameter.\n2. The `PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
			},
		},

		Provider: pschema.ResourceSpec{
			ObjectTypeSpec: pschema.ObjectTypeSpec{
				Description: "The provider type for the kubernetes package.",
				Type:        "object",
			},
			InputProperties: map[string]pschema.PropertySpec{
				"kubeconfig": {
					Description: "The contents of a kubeconfig file. If this is set, this config will be used instead of $KUBECONFIG.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"context": {
					Description: "If present, the name of the kubeconfig context to use.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"cluster": {
					Description: "If present, the name of the kubeconfig cluster to use.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"namespace": {
					Description: "If present, the default namespace to use. This flag is ignored for cluster-scoped resources.\n\nA namespace can be specified in multiple places, and the precedence is as follows:\n1. `.metadata.namespace` set on the resource.\n2. This `namespace` parameter.\n3. `namespace` set for the active context in the kubeconfig.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"enableDryRun": {
					Description: "BETA FEATURE - If present and set to true, enable server-side diff calculations.\nThis feature is in developer preview, and is disabled by default.\n\nThis config can be specified in the following ways, using this precedence:\n1. This `enableDryRun` parameter.\n2. The `PULUMI_K8S_ENABLE_DRY_RUN` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"suppressDeprecationWarnings": {
					Description: "If present and set to true, suppress apiVersion deprecation warnings from the CLI.\n\nThis config can be specified in the following ways, using this precedence:\n1. This `suppressDeprecationWarnings` parameter.\n2. The `PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
			},
		},

		Types:     map[string]pschema.ObjectTypeSpec{},
		Resources: map[string]pschema.ResourceSpec{},
		Functions: map[string]pschema.FunctionSpec{},
		Language:  map[string]json.RawMessage{},
	}

	goImportPath := "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes"

	csharpNamespaces := map[string]string{}
	modToPkg := map[string]string{}
	pkgImportAliases := map[string]string{}

	definitions := swagger["definitions"].(map[string]interface{})
	groupsSlice := createGroups(definitions, schemaOpts())
	for _, group := range groupsSlice {
		for _, version := range group.Versions() {
			mod := version.gv.String()
			csharpNamespaces[mod] = fmt.Sprintf("%s.%s", pascalCase(group.Group()), pascalCase(version.Version()))

			for _, kind := range version.Kinds() {
				tok := fmt.Sprintf(`kubernetes:%s:%s`, kind.canonicalGV, kind.kind)

				modToPkg[kind.canonicalGV] = kind.schemaPkgName
				pkgImportAliases[fmt.Sprintf("%s/%s", goImportPath, kind.schemaPkgName)] = strings.Replace(
					kind.schemaPkgName, "/", "", -1)

				objectSpec := pschema.ObjectTypeSpec{
					Description: kind.Comment() + kind.PulumiComment(),
					Type:        "object",
					Properties:  map[string]pschema.PropertySpec{},
				}

				for _, p := range kind.Properties() {
					// JSONSchema type includes `$ref` and `$schema` properties, and $ is an invalid character in
					// the generated names.
					if strings.HasPrefix(p.name, "$") {
						p.name = "t_" + p.name[1:]
					}
					objectSpec.Properties[p.name] = genPropertySpec(p, kind.canonicalGV, kind.kind)
					//objectSpec.Required = append(objectSpec.Required, p.name)
				}

				pkg.Types[tok] = objectSpec
				if kind.IsNested() {
					continue
				}

				resourceSpec := pschema.ResourceSpec{
					ObjectTypeSpec:     objectSpec,
					DeprecationMessage: kind.DeprecationComment(),
					InputProperties:    map[string]pschema.PropertySpec{},
				}

				for _, p := range kind.RequiredInputProperties() {
					resourceSpec.InputProperties[p.name] = genPropertySpec(p, kind.canonicalGV, kind.kind)
					resourceSpec.RequiredInputs = append(resourceSpec.RequiredInputs, p.name)
				}
				for _, p := range kind.OptionalInputProperties() {
					resourceSpec.InputProperties[p.name] = genPropertySpec(p, kind.canonicalGV, kind.kind)
				}

				for _, t := range kind.Aliases() {
					aliasedType := t
					resourceSpec.Aliases = append(resourceSpec.Aliases, pschema.AliasSpec{Type: &aliasedType})
				}

				pkg.Resources[tok] = resourceSpec
			}
		}
	}

	pkg.Language["csharp"] = rawMessage(map[string]interface{}{
		"packageReferences": map[string]string{
			"Pulumi":                       "1.7.0-preview-alpha.1574743805",
			"System.Collections.Immutable": "1.6.0",
		},
		"namespaces": csharpNamespaces,
	})
	pkg.Language["go"] = rawMessage(map[string]interface{}{
		"importBasePath":       goImportPath,
		"moduleToPackage":      modToPkg,
		"packageImportAliases": pkgImportAliases,
	})
	pkg.Language["nodejs"] = rawMessage(map[string]interface{}{
		"dependencies": map[string]string{
			"@pulumi/pulumi":    "^1.14.0",
			"shell-quote":       "^1.6.1",
			"tmp":               "^0.0.33",
			"@types/tmp":        "^0.0.33",
			"glob":              "^7.1.2",
			"@types/glob":       "^5.0.35",
			"node-fetch":        "^2.3.0",
			"@types/node-fetch": "^2.1.4",
		},
		"devDependencies": map[string]string{
			"mocha":              "^5.2.0",
			"@types/mocha":       "^5.2.5",
			"@types/shell-quote": "^1.6.0",
		},
		"moduleToPackage": modToPkg,
	})
	pkg.Language["python"] = rawMessage(map[string]interface{}{
		"requires": map[string]string{
			"pulumi":   ">=1.11.0,<2.0.0",
			"requests": ">=2.21.0,<2.22.0",
			"pyyaml":   ">=5.1,<5.2",
			"semver":   ">=2.8.1",
			"parver":   ">=0.2.1",
		},
	})

	return pkg
}

func genPropertySpec(p Property, resourceGV string, resourceKind string) pschema.PropertySpec {
	var typ pschema.TypeSpec
	err := json.Unmarshal([]byte(p.ProviderType()), &typ)
	contract.Assert(err == nil)

	constValue := func() *string {
		if p.name == "apiVersion" {
			if strings.HasPrefix(resourceGV, "core/") {
				dv := strings.TrimPrefix(resourceGV, "core/")
				return &dv
			}
			return &resourceGV
		}
		if p.name == "kind" {
			return &resourceKind
		}

		return nil
	}
	defaultValue := func() *string {

		return nil
	}

	propertySpec := pschema.PropertySpec{
		Description: p.Comment(),
		TypeSpec:    typ,
	}
	if cv := constValue(); cv != nil {
		propertySpec.Const = *cv
	}
	if dv := defaultValue(); dv != nil {
		propertySpec.Default = *dv
	}
	return propertySpec
}

func rawMessage(v interface{}) json.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return json.RawMessage(bytes)
}
