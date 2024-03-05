// Copyright 2016-2023, Pulumi Corporation.
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
	_ "embed" // Needed to support go:embed directive

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	v1 "k8s.io/api/core/v1"
)

var serviceSpec = pschema.ComplexTypeSpec{
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
}

var serviceSpecType = pschema.ComplexTypeSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		Type: "string",
	},
	Enum: []pschema.EnumValueSpec{
		{Value: v1.ServiceTypeExternalName},
		{Value: v1.ServiceTypeClusterIP},
		{Value: v1.ServiceTypeNodePort},
		{Value: v1.ServiceTypeLoadBalancer},
	},
}

//go:embed examples/overlays/chartV3.md
var helmV3ChartMD string

var helmV3ChartResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: helmV3ChartMD,
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Resources created by the Chart.",
			},
		},
		Type: "object",
	},
	InputProperties: map[string]pschema.PropertySpec{
		"chart": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "The name of the chart to deploy. If [repo] is provided, this chart name will be prefixed by the repo name. Example: repo: \"stable\", chart: \"nginx-ingress\" -> \"stable/nginx-ingress\" Example: chart: \"stable/nginx-ingress\" -> \"stable/nginx-ingress\"\n\nRequired if specifying `ChartOpts` for a remote chart.",
		},
		"fetchOpts": {
			TypeSpec: pschema.TypeSpec{
				Ref: "#/types/kubernetes:helm.sh/v3:FetchOpts",
			},
			Description: "Additional options to customize the fetching of the Helm chart.",
		},
		"path": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "The path to the chart directory which contains the `Chart.yaml` file.\n\nRequired if specifying `LocalChartOpts`.",
		},
		"namespace": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "The optional namespace to install chart resources into.",
		},
		"repo": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "The repository name of the chart to deploy. Example: \"stable\".\n\nUsed only when specifying options for a remote chart.",
		},
		"resourcePrefix": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
		},
		"transformations": {
			TypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "Optional array of transformations to apply to resources that will be created by this chart prior to creation. Allows customization of the chart behaviour without directly modifying the chart itself.",
		},
		"values": {
			TypeSpec: pschema.TypeSpec{
				Type: "object",
				AdditionalProperties: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "Overrides for chart values.",
		},
		"version": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "The version of the chart to deploy. If not provided, the latest version will be deployed.",
		},
	},
}

var helmV3FetchOpts = pschema.ComplexTypeSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: "Additional options to customize the fetching of the Helm chart.",
		Properties: map[string]pschema.PropertySpec{
			"caFile": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Verify certificates of HTTPS-enabled servers using this CA bundle.",
			},
			"certFile": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Identify HTTPS client using this SSL certificate file.",
			},
			"destination": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Location to write the chart. If this and tardir are specified, tardir is appended to this (default \".\").",
			},
			"devel": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Use development versions, too. Equivalent to version '>0.0.0-0'. If –version is set, this is ignored.",
			},
			"home": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Location of your Helm config. Overrides $HELM_HOME (default \"/Users/abc/.helm\").",
			},
			"keyFile": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Identify HTTPS client using this SSL key file.",
			},
			"keyring": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Keyring containing public keys (default “/Users/abc/.gnupg/pubring.gpg”).",
			},
			"password": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Chart repository password.",
			},
			"prov": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Fetch the provenance file, but don’t perform verification.",
			},
			"repo": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Chart repository url where to locate the requested chart.",
			},
			"untar": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "If set to false, will leave the chart as a tarball after downloading.",
			},
			"untardir": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "If untar is specified, this flag specifies the name of the directory into which the chart is expanded (default \".\").",
			},
			"username": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Chart repository username.",
			},
			"verify": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Verify the package against its signature.",
			},
			"version": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Specific version of a chart. Without this, the latest version is fetched.",
			},
		},
		Type: "object",
	},
}

var helmV3RepoOpts = pschema.ComplexTypeSpec{
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
}

var helmV3ReleaseStatus = pschema.ComplexTypeSpec{
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
}

var kubeClientSettings = pschema.ComplexTypeSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		Description: "Options for tuning the Kubernetes client used by a Provider.",
		Properties: map[string]pschema.PropertySpec{
			"burst": {
				Description: "Maximum burst for throttle. Default value is 10.",
				TypeSpec:    pschema.TypeSpec{Type: "integer"},
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_CLIENT_BURST",
					},
				},
			},
			"qps": {
				Description: "Maximum queries per second (QPS) to the API server from this client. Default value is 5.",
				TypeSpec:    pschema.TypeSpec{Type: "number"},
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_CLIENT_QPS",
					},
				},
			},
			"timeout": {
				Description: "Maximum time in seconds to wait before cancelling a HTTP request to the Kubernetes server. Default value is 32.",
				TypeSpec:    pschema.TypeSpec{Type: "integer"},
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_CLIENT_TIMEOUT",
					},
				},
			},
		},
		Type: "object",
	},
}

var helmReleaseSettings = pschema.ComplexTypeSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		Description: "Options to configure the Helm Release resource.",
		Properties: map[string]pschema.PropertySpec{
			"driver": {
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_HELM_DRIVER",
					},
				},
				Description: "The backend storage driver for Helm. Values are: configmap, secret, memory, sql.",
				TypeSpec:    pschema.TypeSpec{Type: "string"},
			},
			"pluginsPath": {
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_HELM_PLUGINS_PATH",
					},
				},
				Description: "The path to the helm plugins directory.",
				TypeSpec:    pschema.TypeSpec{Type: "string"},
			},
			"registryConfigPath": {
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_HELM_REGISTRY_CONFIG_PATH",
					},
				},
				Description: "The path to the registry config file.",
				TypeSpec:    pschema.TypeSpec{Type: "string"},
			},
			"repositoryConfigPath": {
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_HELM_REPOSITORY_CONFIG_PATH",
					},
				},
				Description: "The path to the file containing repository names and URLs.",
				TypeSpec:    pschema.TypeSpec{Type: "string"},
			},
			"repositoryCache": {
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_HELM_REPOSITORY_CACHE",
					},
				},
				Description: "The path to the file containing cached repository indexes.",
				TypeSpec:    pschema.TypeSpec{Type: "string"},
			},
		},
		Type: "object",
	},
}

//go:embed examples/overlays/helmRelease.md
var helmV3ReleaseMD string

var helmV3ReleaseResource = pschema.ResourceSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		Description: helmV3ReleaseMD,
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
				Description: "List of assets (raw yaml files). Content is read and merged with values (with values taking precedence).",
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
				Description: "Prevent CRD hooks from running, but run other hooks.  See helm install --no-crd-hook",
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
			"allowNullValues": {
				TypeSpec: pschema.TypeSpec{
					Type: "boolean",
				},
				Description: "Whether to allow Null values in helm chart configs.",
			},
		},
		Type: "object",
		Required: []string{
			"chart",
			"status",
		},
		Language: map[string]pschema.RawMessage{
			"nodejs": rawMessage(map[string][]string{
				"requiredOutputs": {
					"name",
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
			Description: "List of assets (raw yaml files). Content is read and merged with values.",
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
			Description: "Prevent CRD hooks from running, but run other hooks.  See helm install --no-crd-hook",
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
		"allowNullValues": {
			TypeSpec: pschema.TypeSpec{
				Type: "boolean",
			},
			Description: "Whether to allow Null values in helm chart configs.",
		},
	},
	RequiredInputs: []string{
		"chart",
	},
}

//go:embed examples/overlays/kustomizeDirectory.md
var kustomizeDirectoryMD string

var kustomizeDirectoryResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: kustomizeDirectoryMD,
		Properties: map[string]pschema.PropertySpec{
			"directory": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "The directory containing the kustomization to apply. The value can be a local directory or a folder in a\ngit repository.\nExample: ./helloWorld\nExample: https://github.com/kubernetes-sigs/kustomize/tree/master/examples/helloWorld",
			},
			"resourcePrefix": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
			},
			"transformations": {
				TypeSpec: pschema.TypeSpec{
					Type: "array",
					Items: &pschema.TypeSpec{
						Ref: "pulumi.json#/Any",
					},
				},
				Description: "A set of transformations to apply to Kubernetes resource definitions before registering with engine.",
			},
		},
		Type: "object",
		Required: []string{
			"directory",
		},
	},
	InputProperties: map[string]pschema.PropertySpec{
		"directory": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "The directory containing the kustomization to apply. The value can be a local directory or a folder in a\ngit repository.\nExample: ./helloWorld\nExample: https://github.com/kubernetes-sigs/kustomize/tree/master/examples/helloWorld",
		},
		"resourcePrefix": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
		},
		"transformations": {
			TypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "A set of transformations to apply to Kubernetes resource definitions before registering with engine.",
		},
	},
	RequiredInputs: []string{
		"directory",
	},
}

//go:embed examples/overlays/configFile.md
var configFileMD string

var yamlConfigFileResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: configFileMD,
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Resources created by the ConfigFile.",
			},
		},
		Type: "object",
	},
	InputProperties: map[string]pschema.PropertySpec{
		"file": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "Path or a URL that uniquely identifies a file.",
		},
		"resourcePrefix": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
		},
		"transformations": {
			TypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "A set of transformations to apply to Kubernetes resource definitions before registering with engine.",
		},
	},
	RequiredInputs: []string{
		"file",
	},
}

//go:embed examples/overlays/configFileV2.md
var configFileV2MD string

var yamlConfigFileV2Resource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   false,
		Description: configFileV2MD,
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "array",
					Items: &pschema.TypeSpec{
						Ref: "pulumi.json#/Any",
					},
				},
				Description: "Resources created by the ConfigFile.",
			},
		},
		Type: "object",
	},
	InputProperties: map[string]pschema.PropertySpec{
		"file": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "Path or a URL that uniquely identifies a file.",
		},
		"resourcePrefix": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
		},
		"skipAwait": {
			TypeSpec: pschema.TypeSpec{
				Type: "boolean",
			},
			Description: "Indicates that child resources should skip the await logic.",
		},
	},
	RequiredInputs: []string{
		"file",
	},
}

//go:embed examples/overlays/configGroup.md
var configGroupMD string

var yamlConfigGroupResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: configGroupMD,
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Resources created by the ConfigGroup.",
			},
		},
		Type: "object",
	},
	InputProperties: map[string]pschema.PropertySpec{
		"files": {
			TypeSpec: pschema.TypeSpec{
				OneOf: []pschema.TypeSpec{
					{
						Type: "string",
					},
					{
						Type: "array",
						Items: &pschema.TypeSpec{
							Type: "string",
						},
					},
				},
			},
			Description: "Set of paths or a URLs that uniquely identify files.",
		},
		"objs": {
			TypeSpec: pschema.TypeSpec{
				OneOf: []pschema.TypeSpec{
					{
						Ref: "pulumi.json#/Any",
					},
					{
						Type: "array",
						Items: &pschema.TypeSpec{
							Ref: "pulumi.json#/Any",
						},
					},
				},
			},
			Description: "Objects representing Kubernetes resources.",
		},
		"resourcePrefix": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
		},
		"transformations": {
			TypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "A set of transformations to apply to Kubernetes resource definitions before registering with engine.",
		},
		"yaml": {
			TypeSpec: pschema.TypeSpec{
				OneOf: []pschema.TypeSpec{
					{
						Type: "string",
					},
					{
						Type: "array",
						Items: &pschema.TypeSpec{
							Type: "string",
						},
					},
				},
			},
			Description: "YAML text containing Kubernetes resource definitions.",
		},
	},
}

//go:embed examples/overlays/configGroupV2.md
var configGroupV2MD string

var yamlConfigGroupV2Resource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   false,
		Description: configGroupV2MD,
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "array",
					Items: &pschema.TypeSpec{
						Ref: "pulumi.json#/Any",
					},
				},
				Description: "Resources created by the ConfigGroup.",
			},
		},
		Type: "object",
	},
	InputProperties: map[string]pschema.PropertySpec{
		"files": {
			TypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Type: "string",
				},
			},
			Description: "Set of paths and/or URLs to Kubernetes manifest files. Supports glob patterns.",
		},
		"objs": {
			TypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "Objects representing Kubernetes resource configurations.",
		},
		"resourcePrefix": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "A prefix for the auto-generated resource names. Defaults to the name of the ConfigGroup. Example: A resource created with resourcePrefix=\"foo\" would produce a resource named \"foo-resourceName\".",
		},
		"skipAwait": {
			TypeSpec: pschema.TypeSpec{
				Type: "boolean",
			},
			Description: "Indicates that child resources should skip the await logic.",
		},
		"yaml": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "A Kubernetes YAML manifest containing Kubernetes resource configuration(s).",
		},
	},
}

var apiextentionsCustomResource = pschema.ResourceSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: "CustomResource represents an instance of a CustomResourceDefinition (CRD). For example, the\n CoreOS Prometheus operator exposes a CRD `monitoring.coreos.com/ServiceMonitor`; to\n instantiate this as a Pulumi resource, one could call `new CustomResource`, passing the\n `ServiceMonitor` resource definition as an argument.",
		Properties: map[string]pschema.PropertySpec{
			"apiVersion": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
			},
			"kind": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
			},
			"metadata": {
				TypeSpec: pschema.TypeSpec{
					Ref: "#/types/kubernetes:meta/v1:ObjectMeta",
				},
				Description: "Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.",
			},
		},
		Type: "object",
		Required: []string{
			"apiVersion",
			"kind",
		},
	},
	InputProperties: map[string]pschema.PropertySpec{
		"apiVersion": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
		},
		"kind": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
		},
		"metadata": {
			TypeSpec: pschema.TypeSpec{
				Ref: "#/types/kubernetes:meta/v1:ObjectMeta",
			},
			Description: "Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.",
		},
		"others": {
			TypeSpec: pschema.TypeSpec{
				Type: "object",
				AdditionalProperties: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "This field is not an actual property. It is used to represent custom property names and their values that can be passed in addition to the other input properties.",
		},
	},
	RequiredInputs: []string{
		"apiVersion",
		"kind",
	},
}

var apiextentionsCustomResourcePatch = pschema.ResourceSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: "CustomResourcePatch represents an instance of a CustomResourceDefinition (CRD). For example, the\n CoreOS Prometheus operator exposes a CRD `monitoring.coreos.com/ServiceMonitor`; to\n instantiate this as a Pulumi resource, one could call `new CustomResourcePatch`, passing the\n `ServiceMonitor` resource definition as an argument.",
		Properties: map[string]pschema.PropertySpec{
			"apiVersion": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
			},
			"kind": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
			},
			"metadata": {
				TypeSpec: pschema.TypeSpec{
					Ref: "#/types/kubernetes:meta/v1:ObjectMeta",
				},
				Description: "Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.",
			},
		},
		Type: "object",
		Required: []string{
			"apiVersion",
			"kind",
		},
	},
	InputProperties: map[string]pschema.PropertySpec{
		"apiVersion": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
		},
		"kind": {
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
		},
		"metadata": {
			TypeSpec: pschema.TypeSpec{
				Ref: "#/types/kubernetes:meta/v1:ObjectMeta",
			},
			Description: "Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.",
		},
		"others": {
			TypeSpec: pschema.TypeSpec{
				Type: "object",
				AdditionalProperties: &pschema.TypeSpec{
					Ref: "pulumi.json#/Any",
				},
			},
			Description: "This field is not an actual property. It is used to represent custom property names and their values that can be passed in addition to the other input properties.",
		},
	},
	RequiredInputs: []string{
		"apiVersion",
		"kind",
	},
}

func init() {
	typeOverlays["kubernetes:core/v1:ServiceSpec"] = serviceSpec
	typeOverlays["kubernetes:core/v1:ServiceSpecType"] = serviceSpecType
	typeOverlays["kubernetes:helm.sh/v3:FetchOpts"] = helmV3FetchOpts
	typeOverlays["kubernetes:helm.sh/v3:RepositoryOpts"] = helmV3RepoOpts
	typeOverlays["kubernetes:helm.sh/v3:ReleaseStatus"] = helmV3ReleaseStatus
	typeOverlays["kubernetes:index:KubeClientSettings"] = kubeClientSettings
	typeOverlays["kubernetes:index:HelmReleaseSettings"] = helmReleaseSettings

	resourceOverlays["kubernetes:apiextensions.k8s.io:CustomResource"] = apiextentionsCustomResource
	resourceOverlays["kubernetes:apiextensions.k8s.io:CustomResourcePatch"] = apiextentionsCustomResourcePatch
	resourceOverlays["kubernetes:helm.sh/v3:Chart"] = helmV3ChartResource
	resourceOverlays["kubernetes:helm.sh/v3:Release"] = helmV3ReleaseResource
	resourceOverlays["kubernetes:kustomize:Directory"] = kustomizeDirectoryResource
	resourceOverlays["kubernetes:yaml:ConfigFile"] = yamlConfigFileResource
	resourceOverlays["kubernetes:yaml/v2:ConfigFile"] = yamlConfigFileV2Resource
	resourceOverlays["kubernetes:yaml:ConfigGroup"] = yamlConfigGroupResource
	resourceOverlays["kubernetes:yaml/v2:ConfigGroup"] = yamlConfigGroupV2Resource
}
