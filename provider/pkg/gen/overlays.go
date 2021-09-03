// Copyright 2016-2020, Pulumi Corporation.
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
	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

// typeOverlays augment the types defined by the kubernetes schema.
var typeOverlays = map[string]pschema.ComplexTypeSpec{
	"kubernetes:core/v1:ServiceSpec": {
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Properties: map[string]pschema.PropertySpec{
				"type": {
					TypeSpec: pschema.TypeSpec{
						OneOf: []pschema.TypeSpec{
							{Type: "string"},
							{Ref: "#/types/kubernetes:core/v1:ServiceSpecType"},
						},
					},
				},
			},
		},
	},
	"kubernetes:core/v1:ServiceSpecType": {
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Type: "string",
		},
		Enum: []pschema.EnumValueSpec{
			{Value: "ExternalName"},
			{Value: "ClusterIP"},
			{Value: "NodePort"},
			{Value: "LoadBalancer"},
		},
	},
	"kubernetes:helm.sh/v3:Release": {
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Description: "A Release is an instance of a chart running in a Kubernetes cluster.\nA Chart is a Helm package. It contains all of the resource definitions necessary to run an application, tool, or service inside of a Kubernetes cluster.\nNote - Helm Release is currently in BETA and may change. Use in production environment is discouraged.",
			Properties: map[string]pschema.PropertySpec{
				"name": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Release name.",
				},
				"repositoryOpts": {
					TypeSpec: pschema.TypeSpec{
						Ref: "#/types/kubernetes:helm.sh/v3:RepositoryOpts",
					},
					Description: "Specification defining the Helm chart repository to use.",
				},
				"chart": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Chart name to be installed. A path may be used.",
				},
				"version": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Specify the exact chart version to install. If this is not specified, the latest version is installed.",
				},
				"devel": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Use chart development versions, too. Equivalent to version '>0.0.0-0'. If `version` is set, this is ignored.",
				},
				"valueYamlFiles": {
					TypeSpec: pschema.TypeSpec{
						Type: "array",
						Items: &pschema.TypeSpec{
							Ref: "pulumi.json#/Asset",
						},
					},
					Description: "List of assets (raw yaml files). Content is read and merged with values. Not yet supported.",
				},
				"values": {
					TypeSpec: pschema.TypeSpec{
						Type: "object",
						AdditionalProperties: &pschema.TypeSpec{
							Ref: "pulumi.json#/Any",
						},
					},
					Description: "Custom values set for the release.",
				},
				"manifest": {
					TypeSpec: pschema.TypeSpec{
						Type: "object",
						AdditionalProperties: &pschema.TypeSpec{
							Ref: "pulumi.json#/Any",
						},
					},
					Description: "The rendered manifests as JSON. Not yet supported.",
				},
				"resourceNames": {
					TypeSpec: pschema.TypeSpec{
						Type: "object",
						AdditionalProperties: &pschema.TypeSpec{
							Type: "array",
							Items: &pschema.TypeSpec{
								Type: "string",
							},
						},
					},
					Description: "Names of resources created by the release grouped by \"kind/version\".",
				},
				"namespace": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Namespace to install the release into.",
				},
				"verify": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Verify the package before installing it.",
				},
				"keyring": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Location of public keys used for verification. Used only if `verify` is true",
				},
				"timeout": {
					TypeSpec: pschema.TypeSpec{
						Type: "integer",
					},
					Description: "Time in seconds to wait for any individual kubernetes operation.",
				},
				"disableWebhooks": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Prevent hooks from running.",
				},
				"disableCRDHooks": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Prevent CRD hooks from, running, but run other hooks.  See helm install --no-crd-hook",
				},
				"reuseValues": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "When upgrading, reuse the last release's values and merge in any overrides. If 'resetValues' is specified, this is ignored",
				},
				"resetValues": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "When upgrading, reset the values to the ones built into the chart.",
				},
				"forceUpdate": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Force resource update through delete/recreate if needed.",
				},
				"recreatePods": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Perform pods restart during upgrade/rollback.",
				},
				"cleanupOnFail": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Allow deletion of new resources created in this upgrade when upgrade fails.",
				},
				"maxHistory": {
					TypeSpec: pschema.TypeSpec{
						Type: "integer",
					},
					Description: "Limit the maximum number of revisions saved per release. Use 0 for no limit.",
				},
				"atomic": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, installation process purges chart on fail. `skipAwait` will be disabled automatically if atomic is used.",
				},
				"skipCrds": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, no CRDs will be installed. By default, CRDs are installed if not already present.",
				},
				"renderSubchartNotes": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, render subchart notes along with the parent.",
				},
				"disableOpenapiValidation": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema",
				},
				"skipAwait": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "By default, the provider waits until all resources are in a ready state before marking the release as successful. Setting this to true will skip such await logic.",
				},
				"waitForJobs": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Will wait until all Jobs have been completed before marking the release as successful. This is ignored if `skipAwait` is enabled.",
				},
				"dependencyUpdate": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Run helm dependency update before installing the chart.",
				},
				"replace": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Re-use the given name, even if that name is already used. This is unsafe in production",
				},
				"description": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Add a custom description",
				},
				"createNamespace": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Create the namespace if it does not exist.",
				},
				"postrender": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Postrender command to run.",
				},
				"lint": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Run helm lint when planning.",
				},
				"status": {
					TypeSpec: pschema.TypeSpec{
						Ref: "#/types/kubernetes:helm.sh/v3:ReleaseStatus",
					},
					Description: "Status of the deployed release.",
				},
			},
			Type: "object",
			Required: []string{
				"chart",
				"repositoryOpts",
				"values",
				"status",
			},
			Language: map[string]pschema.RawMessage{
				"nodejs": rawMessage(map[string][]string{
					"requiredOutputs": {
						"name",
						"repositoryOpts",
						"chart",
						"version",
						"devel",
						"values",
						"set",
						"manifest",
						"namespace",
						"verify",
						"keyring",
						"timeout",
						"disableWebhooks",
						"disableCRDHooks",
						"reuseValues",
						"resetValues",
						"forceUpdate",
						"recreatePods",
						"cleanupOnFail",
						"maxHistory",
						"atomic",
						"skipCrds",
						"renderSubchartNotes",
						"disableOpenapiValidation",
						"skipAwait",
						"waitForJobs",
						"dependencyUpdate",
						"replace",
						"description",
						"createNamespace",
						"postrender",
						"lint",
						"status",
					},
				}),
			},
		},
	},
	"kubernetes:helm.sh/v3:RepositoryOpts": {
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Description: "Specification defining the Helm chart repository to use.",
			Properties: map[string]pschema.PropertySpec{
				"repo": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Repository where to locate the requested chart. If is a URL the chart is installed without installing the repository.",
				},
				"keyFile": { // TODO: Content or file
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "The repository's cert key file",
				},
				"certFile": { // TODO: Content or file
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "The repository's cert file",
				},
				"caFile": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "The Repository's CA File",
				},
				"username": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Username for HTTP basic authentication",
				},
				"password": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Secret:      true,
					Description: "Password for HTTP basic authentication",
				},
			},
			Language: map[string]pschema.RawMessage{
				"nodejs": rawMessage(map[string][]string{
					"requiredOutputs": {
						"repo",
						"keyFile",
						"certFile",
						"caFile",
						"username",
						"password",
					}}),
			},
			Type: "object",
		},
	},
	"kubernetes:helm.sh/v3:ReleaseStatus": {
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Required: []string{"status"},
			Properties: map[string]pschema.PropertySpec{
				"name": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Name is the name of the release.",
				},
				"revision": {
					TypeSpec: pschema.TypeSpec{
						Type: "integer",
					},
					Description: "Version is an int32 which represents the version of the release.",
				},
				"namespace": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Namespace is the kubernetes namespace of the release.",
				},
				"chart": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "The name of the chart.",
				},
				"version": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "A SemVer 2 conformant version string of the chart.",
				},
				"appVersion": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "The version number of the application being deployed.",
				},
				"status": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Status of the release.",
				},
			},
			Language: map[string]pschema.RawMessage{
				"nodejs": rawMessage(map[string][]string{
					"requiredOutputs": {
						"name",
						"revision",
						"namespace",
						"chart",
						"version",
						"appVersion",
						"values",
						"status",
					}}),
			},
			Type: "object",
		},
	},
}

// resourceOverlays augment the resources defined by the kubernetes schema.
var resourceOverlays = map[string]pschema.ResourceSpec{
	"kubernetes:helm.sh/v3:Release": {
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Description: "A Release is an instance of a chart running in a Kubernetes cluster.\n\nA Chart is a Helm package. It contains all of the resource definitions necessary to run an application, tool, or service inside of a Kubernetes cluster.",
			Properties: map[string]pschema.PropertySpec{
				"name": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Release name.",
				},
				"repositoryOpts": {
					TypeSpec: pschema.TypeSpec{
						Ref: "#/types/kubernetes:helm.sh/v3:RepositoryOpts",
					},
					Description: "Specification defining the Helm chart repository to use.",
				},

				"chart": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Chart name to be installed. A path may be used.",
				},
				"version": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Specify the exact chart version to install. If this is not specified, the latest version is installed.",
				},
				"devel": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Use chart development versions, too. Equivalent to version '>0.0.0-0'. If `version` is set, this is ignored.",
				},
				"valueYamlFiles": {
					TypeSpec: pschema.TypeSpec{
						Type: "array",
						Items: &pschema.TypeSpec{
							Ref: "pulumi.json#/Asset",
						},
					},
					Description: "List of assets (raw yaml files). Content is read and merged with values. Not yet supported.",
				},
				"values": {
					TypeSpec: pschema.TypeSpec{
						Type: "object",
						AdditionalProperties: &pschema.TypeSpec{
							Ref: "pulumi.json#/Any",
						},
					},
					Description: "Custom values set for the release.",
				},
				"manifest": {
					TypeSpec: pschema.TypeSpec{
						Type: "object",
						AdditionalProperties: &pschema.TypeSpec{
							Ref: "pulumi.json#/Any",
						},
					},
					Description: "The rendered manifests as JSON. Not yet supported.",
				},
				"resourceNames": {
					TypeSpec: pschema.TypeSpec{
						Type: "object",
						AdditionalProperties: &pschema.TypeSpec{
							Type: "array",
							Items: &pschema.TypeSpec{
								Type: "string",
							},
						},
					},
					Description: "Names of resources created by the release grouped by \"kind/version\".",
				},
				"namespace": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Namespace to install the release into.",
				},
				"verify": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Verify the package before installing it.",
				},
				"keyring": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Location of public keys used for verification. Used only if `verify` is true",
				},
				"timeout": {
					TypeSpec: pschema.TypeSpec{
						Type: "integer",
					},
					Description: "Time in seconds to wait for any individual kubernetes operation.",
				},
				"disableWebhooks": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Prevent hooks from running.",
				},
				"disableCRDHooks": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Prevent CRD hooks from, running, but run other hooks.  See helm install --no-crd-hook",
				},
				"reuseValues": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "When upgrading, reuse the last release's values and merge in any overrides. If 'resetValues' is specified, this is ignored",
				},
				"resetValues": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "When upgrading, reset the values to the ones built into the chart.",
				},
				"forceUpdate": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Force resource update through delete/recreate if needed.",
				},
				"recreatePods": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Perform pods restart during upgrade/rollback.",
				},
				"cleanupOnFail": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Allow deletion of new resources created in this upgrade when upgrade fails.",
				},
				"maxHistory": {
					TypeSpec: pschema.TypeSpec{
						Type: "integer",
					},
					Description: "Limit the maximum number of revisions saved per release. Use 0 for no limit.",
				},
				"atomic": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, installation process purges chart on fail. `skipAwait` will be disabled automatically if atomic is used.",
				},
				"skipCrds": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, no CRDs will be installed. By default, CRDs are installed if not already present.",
				},
				"renderSubchartNotes": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, render subchart notes along with the parent.",
				},
				"disableOpenapiValidation": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "If set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema",
				},
				"skipAwait": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "By default, the provider waits until all resources are in a ready state before marking the release as successful. Setting this to true will skip such await logic.",
				},
				"waitForJobs": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Will wait until all Jobs have been completed before marking the release as successful. This is ignored if `skipAwait` is enabled.",
				},
				"dependencyUpdate": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Run helm dependency update before installing the chart.",
				},
				"replace": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Re-use the given name, even if that name is already used. This is unsafe in production",
				},
				"description": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Add a custom description",
				},
				"createNamespace": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Create the namespace if it does not exist.",
				},
				"postrender": {
					TypeSpec: pschema.TypeSpec{
						Type: "string",
					},
					Description: "Postrender command to run.",
				},
				"lint": {
					TypeSpec: pschema.TypeSpec{
						Type: "boolean",
					},
					Description: "Run helm lint when planning.",
				},
				"status": {
					TypeSpec: pschema.TypeSpec{
						Ref: "#/types/kubernetes:helm.sh/v3:ReleaseStatus",
					},
					Description: "Status of the deployed release.",
				},
			},
			Type: "object",
			Required: []string{
				"chart",
				"repositoryOpts",
				"values",
				"status",
			},
			Language: map[string]pschema.RawMessage{
				"nodejs": rawMessage(map[string][]string{
					"requiredOutputs": {
						"name",
						"repositoryOpts",
						"chart",
						"version",
						"devel",
						"values",
						"set",
						"manifest",
						"namespace",
						"verify",
						"keyring",
						"timeout",
						"disableWebhooks",
						"disableCRDHooks",
						"reuseValues",
						"resetValues",
						"forceUpdate",
						"recreatePods",
						"cleanupOnFail",
						"maxHistory",
						"atomic",
						"skipCrds",
						"renderSubchartNotes",
						"disableOpenapiValidation",
						"skipAwait",
						"waitForJobs",
						"dependencyUpdate",
						"replace",
						"description",
						"createNamespace",
						"postrender",
						"lint",
						"status",
					},
				}),
			},
		},
		InputProperties: map[string]pschema.PropertySpec{
			"name": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Release name.",
			},
			"repositoryOpts": {
				TypeSpec: pschema.TypeSpec{
					Ref: "#/types/kubernetes:helm.sh/v3:RepositoryOpts",
				},
				Description: "Specification defining the Helm chart repository to use.",
			},

			"chart": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Chart name to be installed. A path may be used.",
			},
			"version": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Specify the exact chart version to install. If this is not specified, the latest version is installed.",
			},
			"devel": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Use chart development versions, too. Equivalent to version '>0.0.0-0'. If `version` is set, this is ignored.",
			},
			"valueYamlFiles": {
				TypeSpec: pschema.TypeSpec{
					Type: "array",
					Items: &pschema.TypeSpec{
						Ref: "pulumi.json#/Asset",
					},
				},
				Description: "List of assets (raw yaml files). Content is read and merged with values. Not yet supported.",
			},
			"values": {
				TypeSpec: pschema.TypeSpec{
					Type: "object",
					AdditionalProperties: &pschema.TypeSpec{
						Ref: "pulumi.json#/Any",
					},
				},
				Description: "Custom values set for the release.",
			},
			"manifest": {
				TypeSpec: pschema.TypeSpec{
					Type: "object",
					AdditionalProperties: &pschema.TypeSpec{
						Ref: "pulumi.json#/Any",
					},
				},
				Description: "The rendered manifests as JSON. Not yet supported.",
			},
			"resourceNames": {
				TypeSpec: pschema.TypeSpec{
					Type: "object",
					AdditionalProperties: &pschema.TypeSpec{
						Type: "array",
						Items: &pschema.TypeSpec{
							Type: "string",
						},
					},
				},
				Description: "Names of resources created by the release grouped by \"kind/version\".",
			},
			"namespace": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Namespace to install the release into.",
			},
			"verify": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Verify the package before installing it.",
			},
			"keyring": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Location of public keys used for verification. Used only if `verify` is true",
			},
			"timeout": {
				TypeSpec: pschema.TypeSpec{
					Type: "integer",
				},
				Description: "Time in seconds to wait for any individual kubernetes operation.",
			},
			"disableWebhooks": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Prevent hooks from running.",
			},
			"disableCRDHooks": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Prevent CRD hooks from, running, but run other hooks.  See helm install --no-crd-hook",
			},
			"reuseValues": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "When upgrading, reuse the last release's values and merge in any overrides. If 'resetValues' is specified, this is ignored",
			},
			"resetValues": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "When upgrading, reset the values to the ones built into the chart.",
			},
			"forceUpdate": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Force resource update through delete/recreate if needed.",
			},
			"recreatePods": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Perform pods restart during upgrade/rollback.",
			},
			"cleanupOnFail": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Allow deletion of new resources created in this upgrade when upgrade fails.",
			},
			"maxHistory": {
				TypeSpec: pschema.TypeSpec{
					Type: "integer",
				},
				Description: "Limit the maximum number of revisions saved per release. Use 0 for no limit.",
			},
			"atomic": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "If set, installation process purges chart on fail. `skipAwait` will be disabled automatically if atomic is used.",
			},
			"skipCrds": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "If set, no CRDs will be installed. By default, CRDs are installed if not already present.",
			},
			"renderSubchartNotes": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "If set, render subchart notes along with the parent.",
			},
			"disableOpenapiValidation": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "If set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema",
			},
			"skipAwait": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "By default, the provider waits until all resources are in a ready state before marking the release as successful. Setting this to true will skip such await logic.",
			},
			"waitForJobs": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Will wait until all Jobs have been completed before marking the release as successful. This is ignored if `skipAwait` is enabled.",
			},
			"dependencyUpdate": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Run helm dependency update before installing the chart.",
			},
			"replace": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Re-use the given name, even if that name is already used. This is unsafe in production",
			},
			"description": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Add a custom description",
			},
			"createNamespace": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Create the namespace if it does not exist.",
			},
			"postrender": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Postrender command to run.",
			},
			"lint": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Run helm lint when planning.",
			},
			"compat": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Const: "true",
			},
		},
		RequiredInputs: []string{
			"chart",
			"repositoryOpts",
			"values",
		},
	},
}
