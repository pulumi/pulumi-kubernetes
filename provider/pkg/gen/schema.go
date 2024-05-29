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

// nolint: goconst
package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gen/examples"
	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// typeOverlays augment the types defined by the kubernetes schema.
var typeOverlays = map[string]pschema.ComplexTypeSpec{}

// resourceOverlays augment the resources defined by the kubernetes schema.
var resourceOverlays = map[string]pschema.ResourceSpec{}

// PulumiSchema will generate a Pulumi schema for the given k8s schema.
func PulumiSchema(swagger map[string]any) pschema.PackageSpec {
	pkg := pschema.PackageSpec{
		Name:        "kubernetes",
		Description: "A Pulumi package for creating and managing Kubernetes resources.",
		DisplayName: "Kubernetes",
		License:     "Apache-2.0",
		Keywords:    []string{"pulumi", "kubernetes", "category/cloud", "kind/native"},
		Homepage:    "https://pulumi.com",
		Publisher:   "Pulumi",
		Repository:  "https://github.com/pulumi/pulumi-kubernetes",

		Config: pschema.ConfigSpec{
			Variables: map[string]pschema.PropertySpec{
				"kubeconfig": {
					Description: "The contents of a kubeconfig file or the path to a kubeconfig file. If this is set," +
						" this config will be used instead of $KUBECONFIG.",
					TypeSpec: pschema.TypeSpec{Type: "string"},
					Language: map[string]pschema.RawMessage{
						"csharp": rawMessage(map[string]any{
							"name": "KubeConfig",
						}),
					},
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
				"deleteUnreachable": {
					Description: "If present and set to true, the provider will delete resources associated with an unreachable Kubernetes cluster from Pulumi state",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"skipUpdateUnreachable": {
					Description: "If present and set to true, the provider will skip resources update associated with an unreachable Kubernetes cluster from Pulumi state",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"enableServerSideApply": {
					Description: "If present and set to false, disable Server-Side Apply mode.\nSee https://github.com/pulumi/pulumi-kubernetes/issues/2011 for additional details.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"enableReplaceCRD": {
					Description:        "Obsolete. This option has no effect.",
					TypeSpec:           pschema.TypeSpec{Type: "boolean"},
					DeprecationMessage: "This option is deprecated, and will be removed in a future release.",
				},
				"enableConfigMapMutable": {
					Description: "BETA FEATURE - If present and set to true, allow ConfigMaps to be mutated.\nThis feature is in developer preview, and is disabled by default.\n\nThis config can be specified in the following ways using this precedence:\n1. This `enableConfigMapMutable` parameter.\n2. The `PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"renderYamlToDirectory": {
					Description: "BETA FEATURE - If present, render resource manifests to this directory. In this mode, resources will not\nbe created on a Kubernetes cluster, but the rendered manifests will be kept in sync with changes\nto the Pulumi program. This feature is in developer preview, and is disabled by default.\n\nNote that some computed Outputs such as status fields will not be populated\nsince the resources are not created on a Kubernetes cluster. These Output values will remain undefined,\nand may result in an error if they are referenced by other resources. Also note that any secret values\nused in these resources will be rendered in plaintext to the resulting YAML.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"suppressDeprecationWarnings": {
					Description: "If present and set to true, suppress apiVersion deprecation warnings from the CLI.\n\nThis config can be specified in the following ways, using this precedence:\n1. This `suppressDeprecationWarnings` parameter.\n2. The `PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"suppressHelmHookWarnings": {
					Description: "If present and set to true, suppress unsupported Helm hook warnings from the CLI.\n\nThis config can be specified in the following ways, using this precedence:\n1. This `suppressHelmHookWarnings` parameter.\n2. The `PULUMI_K8S_SUPPRESS_HELM_HOOK_WARNINGS` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"strictMode": {
					Description: "If present and set to true, the provider will use strict configuration mode. Recommended for production stacks. In this mode, the default Kubernetes provider is disabled, and the `kubeconfig` and `context` settings are required for Provider configuration. These settings unambiguously ensure that every Kubernetes resource is associated with a particular cluster.",
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
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"KUBECONFIG",
						},
					},
					Description: "The contents of a kubeconfig file or the path to a kubeconfig file.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
					Language: map[string]pschema.RawMessage{
						"csharp": rawMessage(map[string]any{
							"name": "KubeConfig",
						}),
					},
				},
				"context": {
					Description: "If present, the name of the kubeconfig context to use.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"cluster": {
					Description: "If present, the name of the kubeconfig cluster to use.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"deleteUnreachable": {
					Description: "If present and set to true, the provider will delete resources associated with an unreachable Kubernetes cluster from Pulumi state",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"PULUMI_K8S_DELETE_UNREACHABLE",
						},
					},
				},
				"skipUpdateUnreachable": {
					Description: "If present and set to true, the provider will skip resources update associated with an unreachable Kubernetes cluster from Pulumi state",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"PULUMI_K8S_SKIP_UPDATE_UNREACHABLE",
						},
					},
				},
				"namespace": {
					Description: "If present, the default namespace to use. This flag is ignored for cluster-scoped resources.\n\nA namespace can be specified in multiple places, and the precedence is as follows:\n1. `.metadata.namespace` set on the resource.\n2. This `namespace` parameter.\n3. `namespace` set for the active context in the kubeconfig.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"enableServerSideApply": {
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"PULUMI_K8S_ENABLE_SERVER_SIDE_APPLY",
						},
					},
					Description: "If present and set to false, disable Server-Side Apply mode.\nSee https://github.com/pulumi/pulumi-kubernetes/issues/2011 for additional details.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"enableConfigMapMutable": {
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE",
						},
					},
					Description: "BETA FEATURE - If present and set to true, allow ConfigMaps to be mutated.\nThis feature is in developer preview, and is disabled by default.\n\nThis config can be specified in the following ways using this precedence:\n1. This `enableConfigMapMutable` parameter.\n2. The `PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE` environment variable.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"renderYamlToDirectory": {
					Description: "BETA FEATURE - If present, render resource manifests to this directory. In this mode, resources will not\nbe created on a Kubernetes cluster, but the rendered manifests will be kept in sync with changes\nto the Pulumi program. This feature is in developer preview, and is disabled by default.\n\nNote that some computed Outputs such as status fields will not be populated\nsince the resources are not created on a Kubernetes cluster. These Output values will remain undefined,\nand may result in an error if they are referenced by other resources. Also note that any secret values\nused in these resources will be rendered in plaintext to the resulting YAML.",
					TypeSpec:    pschema.TypeSpec{Type: "string"},
				},
				"suppressDeprecationWarnings": {
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS",
						},
					},
					Description: "If present and set to true, suppress apiVersion deprecation warnings from the CLI.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
				"kubeClientSettings": {
					Description: "Options for tuning the Kubernetes client used by a Provider.",
					TypeSpec:    pschema.TypeSpec{Ref: "#/types/kubernetes:index:KubeClientSettings"},
				},
				"helmReleaseSettings": {
					Description: "Options to configure the Helm Release resource.",
					TypeSpec:    pschema.TypeSpec{Ref: "#/types/kubernetes:index:HelmReleaseSettings"},
				},
				"suppressHelmHookWarnings": {
					DefaultInfo: &pschema.DefaultSpec{
						Environment: []string{
							"PULUMI_K8S_SUPPRESS_HELM_HOOK_WARNINGS",
						},
					},
					Description: "If present and set to true, suppress unsupported Helm hook warnings from the CLI.",
					TypeSpec:    pschema.TypeSpec{Type: "boolean"},
				},
			},
		},

		Types:     map[string]pschema.ComplexTypeSpec{},
		Resources: map[string]pschema.ResourceSpec{},
		Functions: map[string]pschema.FunctionSpec{},
		Language:  map[string]pschema.RawMessage{},
	}

	goImportPath := "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"

	csharpNamespaces := map[string]string{
		"apiextensions": "ApiExtensions",
		"helm.sh/v2":    "Helm.V2",
		"helm.sh/v3":    "Helm.V3",
		"helm.sh/v4":    "Helm.V4",
		"kustomize/v2":  "Kustomize.V2",
		"yaml":          "Yaml",
		"yaml/v2":       "Yaml.V2",
		"":              "Provider",
	}
	javaPackages := map[string]string{
		"helm.sh/v2":   "helm.v2",
		"helm.sh/v3":   "helm.v3",
		"helm.sh/v4":   "helm.v4",
		"kustomize/v2": "kustomize.v2",
		"yaml/v2":      "yaml.v2",
	}
	modToPkg := map[string]string{
		"apiextensions.k8s.io": "apiextensions",
		"helm.sh":              "helm",
		"helm.sh/v2":           "helm/v2",
		"helm.sh/v3":           "helm/v3",
		"helm.sh/v4":           "helm/v4",
	}
	pkgImportAliases := map[string]string{
		"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3":      "helmv3",
		"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4":      "helmv4",
		"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize/v2": "kustomizev2",
		"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2":      "yamlv2",
	}

	definitions := swagger["definitions"].(map[string]any)
	groupsSlice := createGroups(definitions)

	for _, group := range groupsSlice {
		if group.Group() == "apiserverinternal" {
			continue
		}
		for _, version := range group.Versions() {
			for _, kind := range version.Kinds() {
				tok := fmt.Sprintf(`kubernetes:%s:%s`, kind.apiVersion, kind.kind)
				var patchTok string
				if !strings.HasSuffix(kind.kind, "List") {
					patchTok = fmt.Sprintf("%sPatch", tok)
				}

				csharpNamespaces[kind.apiVersion] = fmt.Sprintf("%s.%s", pascalCase(group.Group()), pascalCase(version.Version()))
				javaPackages[kind.apiVersion] = fmt.Sprintf("%s.%s", group.Group(), version.Version())
				modToPkg[kind.apiVersion] = kind.schemaPkgName
				pkgImportAliases[fmt.Sprintf("%s/%s", goImportPath, kind.schemaPkgName)] = strings.Replace(
					kind.schemaPkgName, "/", "", -1)

				objectSpec := pschema.ObjectTypeSpec{
					Description: kind.Comment() + kind.PulumiComment(),
					Type:        "object",
					Properties:  map[string]pschema.PropertySpec{},
					Language:    map[string]pschema.RawMessage{},
				}

				var propNames []string
				for _, p := range kind.Properties() {
					objectSpec.Properties[p.name] = genPropertySpec(p, kind.kind)
					propNames = append(propNames, p.name)
				}
				for _, p := range kind.RequiredInputProperties() {
					objectSpec.Required = append(objectSpec.Required, p.name)
				}
				objectSpec.Language["nodejs"] = rawMessage(map[string][]string{"requiredOutputs": propNames})

				// Check if the current type exists in the overlays and overwrite types accordingly.
				if overlaySpec, hasType := typeOverlays[tok]; hasType {
					for propName, overlayProp := range overlaySpec.Properties {
						// If overlay prop types are defined, overwrite the k8s schema prop type.
						if len(overlayProp.OneOf) > 1 {
							if k8sProp, propExists := objectSpec.Properties[propName]; propExists {
								k8sProp.OneOf = overlayProp.OneOf
								k8sProp.Ref = ""
								k8sProp.Type = ""

								objectSpec.Properties[propName] = k8sProp
							}
						}
					}
				}

				pkg.Types[tok] = pschema.ComplexTypeSpec{
					ObjectTypeSpec: objectSpec,
				}

				// Add "Patch" types corresponding to all non "List" resource kinds.
				// All fields except for name are made optional by leaving off the required schema field.
				var patchSpec pschema.ObjectTypeSpec
				if len(patchTok) > 0 {
					patchSpec = pschema.ObjectTypeSpec{
						Description: kind.Comment() + kind.PulumiComment(),
						Type:        "object",
						Properties:  map[string]pschema.PropertySpec{},
						Language:    map[string]pschema.RawMessage{},
					}

					var propNames []string
					for _, p := range kind.Properties() {
						patchSpec.Properties[p.name] = genPropertySpec(p, kind.kind+"Patch")
						propNames = append(propNames, p.name)
					}

					patchSpec.Language["nodejs"] = rawMessage(map[string][]string{"requiredOutputs": propNames})

					// Check if the current type exists in the overlays and overwrite types accordingly.
					if overlaySpec, hasType := typeOverlays[tok]; hasType {
						for propName, overlayProp := range overlaySpec.Properties {
							// If overlay prop types are defined, overwrite the k8s schema prop type.
							if len(overlayProp.OneOf) > 1 {
								if k8sProp, propExists := patchSpec.Properties[propName]; propExists {
									k8sProp.OneOf = overlayProp.OneOf
									k8sProp.Ref = ""
									k8sProp.Type = ""

									patchSpec.Properties[propName] = k8sProp
								}
							}
						}
					}

					pkg.Types[patchTok] = pschema.ComplexTypeSpec{
						ObjectTypeSpec: patchSpec,
					}
				}

				if kind.IsNested() {
					continue
				}

				resourceSpec := pschema.ResourceSpec{
					ObjectTypeSpec:     objectSpec,
					DeprecationMessage: kind.DeprecationComment(),
					InputProperties:    map[string]pschema.PropertySpec{},
				}

				for _, p := range kind.RequiredInputProperties() {
					resourceSpec.InputProperties[p.name] = genPropertySpec(p, kind.kind)
					resourceSpec.RequiredInputs = append(resourceSpec.RequiredInputs, p.name)
				}
				for _, p := range kind.OptionalInputProperties() {
					resourceSpec.InputProperties[p.name] = genPropertySpec(p, kind.kind)
				}
				var required []string
				for p := range resourceSpec.InputProperties {
					required = append(required, p)
				}
				sort.Strings(required)
				resourceSpec.Required = required

				for _, t := range kind.Aliases() {
					aliasedType := t
					resourceSpec.Aliases = append(resourceSpec.Aliases, pschema.AliasSpec{Type: &aliasedType})
				}

				// Check if the current resource exists in the overlays and overwrite types accordingly.
				if overlaySpec, hasResource := resourceOverlays[tok]; hasResource {
					for propName, overlayProp := range overlaySpec.InputProperties {
						// If overlay prop types are defined, overwrite the k8s schema prop type.
						if len(overlayProp.OneOf) > 1 {
							if k8sProp, propExists := objectSpec.Properties[propName]; propExists {
								k8sProp.OneOf = overlayProp.OneOf
								k8sProp.Ref = ""
								k8sProp.Type = ""

								resourceSpec.InputProperties[propName] = k8sProp
							}
						}
					}
				}

				pkg.Resources[tok] = resourceSpec

				// Add "Patch" resources corresponding to all non "List" resource kinds.
				// All fields except for metadata are made optional by leaving off the requiredInputs schema field.
				if len(patchTok) > 0 {
					patchResourceSpec := pschema.ResourceSpec{
						ObjectTypeSpec:     patchSpec,
						DeprecationMessage: kind.DeprecationComment(),
						InputProperties:    map[string]pschema.PropertySpec{},
					}

					for _, p := range kind.RequiredInputProperties() {
						patchResourceSpec.InputProperties[p.name] = genPropertySpec(p, kind.kind+"Patch")
					}
					for _, p := range kind.OptionalInputProperties() {
						patchResourceSpec.InputProperties[p.name] = genPropertySpec(p, kind.kind+"Patch")
					}

					for _, t := range kind.Aliases() {
						aliasedType := fmt.Sprintf("%sPatch", t)
						patchResourceSpec.Aliases = append(patchResourceSpec.Aliases, pschema.AliasSpec{Type: &aliasedType})
					}

					// Check if the current resource exists in the overlays and overwrite types accordingly.
					if overlaySpec, hasResource := resourceOverlays[patchTok]; hasResource {
						for propName, overlayProp := range overlaySpec.InputProperties {
							// If overlay prop types are defined, overwrite the k8s schema prop type.
							if len(overlayProp.OneOf) > 1 {
								if k8sProp, propExists := objectSpec.Properties[propName]; propExists {
									k8sProp.OneOf = overlayProp.OneOf
									k8sProp.Ref = ""
									k8sProp.Type = ""

									patchResourceSpec.InputProperties[propName] = k8sProp
								}
							}
						}
					}

					patchDescription := `Patch resources are used to modify existing Kubernetes resources by using
Server-Side Apply updates. The name of the resource must be specified, but all other properties are optional. More than 
one patch may be applied to the same resource, and a random FieldManager name will be used for each Patch resource. 
Conflicts will result in an error by default, but can be forced using the "pulumi.com/patchForce" annotation. See the
[Server-Side Apply Docs](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/) for
additional information about using Server-Side Apply to manage Kubernetes resources with Pulumi.`
					patchResourceSpec.Description = fmt.Sprintf(
						"%s\n%s", patchDescription, patchResourceSpec.Description)

					pkg.Resources[patchTok] = patchResourceSpec
				}
			}
		}

		// If there are types in the overlays that do not exist in the schema (i.e. enum types), add them.
		for tok, overlayType := range typeOverlays {
			if _, typeExists := pkg.Types[tok]; !typeExists {
				pkg.Types[tok] = overlayType
			}
		}

		// Finally, add overlay resources that weren't in the schema.
		for tok := range resourceOverlays {
			if _, resourceExists := pkg.Resources[tok]; !resourceExists {
				pkg.Resources[tok] = resourceOverlays[tok]
			}
		}
	}

	// Add examples to resources
	for k, v := range examples.ApiVersionToExample {
		if r, ok := pkg.Resources[k]; ok {
			r.Description += "\n\n" + v
			pkg.Resources[k] = r
		}
	}

	// Compatibility mode for Kubernetes 2.0 SDK
	const kubernetes20 = "kubernetes20"

	pkg.Language["csharp"] = rawMessage(map[string]any{
		"respectSchemaVersion": true,
		"packageReferences": map[string]string{
			"Glob":   "1.1.5",
			"Pulumi": "3.*",
		},
		"namespaces":             csharpNamespaces,
		"compatibility":          kubernetes20,
		"dictionaryConstructors": true,
	})

	pkg.Language["java"] = rawMessage(map[string]any{
		"packages": javaPackages,
		"dependencies": map[string]string{
			"net.bytebuddy:byte-buddy": "1.14.15",
			"com.google.guava:guava":   "32.1.2-jre",
		},
	})

	pkg.Language["go"] = rawMessage(map[string]any{
		"respectSchemaVersion":           true,
		"importBasePath":                 goImportPath,
		"moduleToPackage":                modToPkg,
		"packageImportAliases":           pkgImportAliases,
		"generateResourceContainerTypes": true,
		"generateExtraInputTypes":        true,
		"internalModuleName":             "utilities",
	})
	pkg.Language["nodejs"] = rawMessage(map[string]any{
		"respectSchemaVersion": true,
		"compatibility":        kubernetes20,
		"dependencies": map[string]string{
			"@pulumi/pulumi":    "^3.25.0",
			"shell-quote":       "^1.6.1",
			"tmp":               "^0.0.33",
			"@types/tmp":        "^0.0.33",
			"glob":              "^10.3.10",
			"node-fetch":        "^2.3.0",
			"@types/node-fetch": "^2.1.4",
		},
		"devDependencies": map[string]string{
			"mocha":              "^5.2.0",
			"@types/mocha":       "^5.2.5",
			"@types/shell-quote": "^1.6.0",
		},
		"moduleToPackage": modToPkg,
		"readme": `The Kubernetes provider package offers support for all Kubernetes resources and their properties.
Resources are exposed as types from modules based on Kubernetes API groups such as 'apps', 'core',
'rbac', and 'storage', among many others. Additionally, support for deploying Helm charts ('helm')
and YAML files ('yaml') is available in this package. Using this package allows you to
programmatically declare instances of any Kubernetes resources and any supported resource version
using infrastructure as code, which Pulumi then uses to drive the Kubernetes API.

If this is your first time using this package, these two resources may be helpful:

* [Kubernetes Getting Started Guide](https://www.pulumi.com/docs/quickstart/kubernetes/): Get up and running quickly.
* [Kubernetes Pulumi Setup Documentation](https://www.pulumi.com/docs/quickstart/kubernetes/configure/): How to configure Pulumi
    for use with your Kubernetes cluster.

Use the navigation below to see detailed documentation for each of the supported Kubernetes resources.
`,
	})
	pkg.Language["python"] = rawMessage(map[string]any{
		"respectSchemaVersion": true,
		"requires": map[string]string{
			"pulumi":   ">=3.109.0,<4.0.0",
			"requests": ">=2.21,<3.0",
		},
		"pyproject": map[string]bool{
			"enabled": true,
		},
		"moduleNameOverrides": modToPkg,
		"compatibility":       kubernetes20,
		"usesIOClasses":       true,
		"readme": `The Kubernetes provider package offers support for all Kubernetes resources and their properties.
Resources are exposed as types from modules based on Kubernetes API groups such as 'apps', 'core',
'rbac', and 'storage', among many others. Additionally, support for deploying Helm charts ('helm')
and YAML files ('yaml') is available in this package. Using this package allows you to
programmatically declare instances of any Kubernetes resources and any supported resource version
using infrastructure as code, which Pulumi then uses to drive the Kubernetes API.

If this is your first time using this package, these two resources may be helpful:

* [Kubernetes Getting Started Guide](https://www.pulumi.com/docs/quickstart/kubernetes/): Get up and running quickly.
* [Kubernetes Pulumi Setup Documentation](https://www.pulumi.com/docs/quickstart/kubernetes/configure/): How to configure Pulumi
    for use with your Kubernetes cluster.
`,
	})

	return pkg
}

func genPropertySpec(p Property, resourceKind string) pschema.PropertySpec {
	var typ pschema.TypeSpec
	err := json.Unmarshal([]byte(p.SchemaType()), &typ)
	contract.AssertNoErrorf(err, "unexpected error while unmarshalling JSON")

	if strings.HasSuffix(resourceKind, "Patch") {
		if len(typ.Ref) > 0 && !strings.Contains(typ.Ref, "pulumi.json#") {
			typ.Ref += "Patch"
		}
		if typ.Type == "array" && typ.Items != nil {
			if len(typ.Items.Ref) > 0 && !strings.Contains(typ.Items.Ref, "pulumi.json#") {
				typ.Items.Ref += "Patch"
			}
		}
	}

	constValue := func() *string {
		cv := p.ConstValue()
		if len(cv) != 0 {
			return &cv
		}

		return nil
	}
	defaultValue := func() *string {
		dv := p.DefaultValue()
		if len(dv) != 0 {
			return &dv
		}
		return nil
	}

	propertySpec := pschema.PropertySpec{
		Description: p.Comment(),
		TypeSpec:    typ,
	}
	if cv := constValue(); cv != nil {
		propertySpec.Const = *cv
	} else if dv := defaultValue(); dv != nil {
		propertySpec.Default = *dv
	}
	languageName := strings.ToUpper(p.name[:1]) + p.name[1:]
	if languageName == resourceKind {
		// .NET does not allow properties to be the same as the enclosing class - so special case these
		propertySpec.Language = map[string]pschema.RawMessage{
			"csharp": rawMessage(map[string]any{
				"name": languageName + "Value",
			}),
		}
	}
	// JSONSchema type includes `$ref` and `$schema` properties, and $ is an invalid character in
	// the generated names. Replace them with `Ref` and `Schema`.
	if strings.HasPrefix(p.name, "$") {
		propertySpec.Language = map[string]pschema.RawMessage{
			"csharp": rawMessage(map[string]any{
				"name": strings.ToUpper(p.name[1:2]) + p.name[2:],
			}),
		}
	}
	if resourceKind == "Secret" {
		switch p.Name() {
		case "data", "stringData":
			propertySpec.Secret = true
		}
	}
	return propertySpec
}

func rawMessage(v any) pschema.RawMessage {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	contract.AssertNoErrorf(err, "unexpected error while encoding JSON")
	return out.Bytes()
}
