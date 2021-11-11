// Copyright 2016-2021, Pulumi Corporation.
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

var helmV2ChartResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: "Chart is a component representing a collection of resources described by an arbitrary Helm Chart.\n\nThe Chart can be fetched from any source that is accessible to the `helm` command line. Values in the `values.yml` file can be overridden using `ChartOpts.values` (equivalent to `--set` or having multiple `values.yml` files). Objects can be transformed arbitrarily by supplying callbacks to `ChartOpts.transformations`.\n\n`Chart` does not use Tiller. The Chart specified is copied and expanded locally; the semantics are equivalent to running `helm template` and then using Pulumi to manage the resulting YAML manifests. Any values that would be retrieved in-cluster are assigned fake values, and none of Tiller's server-side validity testing is executed.\n\n{{% examples %}}\n## Example Usage\n{{% example %}}\n### Local Chart Directory\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v2.Chart(\"nginx-ingress\", {\n    path: \"./nginx-ingress\",\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v2 import Chart, LocalChartOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    LocalChartOpts(\n        path=\"./nginx-ingress\",\n    ),\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V2;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new LocalChartArgs\n        {\n            Path = \"./nginx-ingress\",\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tPath: pulumi.String(\"./nginx-ingress\"),\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Remote Chart\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v2.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v2 import Chart, ChartOpts, FetchOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n    ),\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V2;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            }\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.run(func(ctx *pulumi.context) error {\n\t\t_, err := helm.newchart(ctx, \"nginx-ingress\", helm.chartargs{\n\t\t\tchart:   pulumi.string(\"nginx-ingress\"),\n\t\t\tversion: pulumi.string(\"1.24.4\"),\n\t\t\tfetchargs: helm.fetchargs{\n\t\t\t\trepo: pulumi.string(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n\n{{% /example %}}\n{{% example %}}\n### Set Chart Values\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v2.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n    values: {\n        controller: {\n            metrics: {\n                enabled: true,\n            }\n        }\n    },\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v2 import Chart, ChartOpts, FetchOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n        values={\n            \"controller\": {\n                \"metrics\": {\n                    \"enabled\": True,\n                },\n            },\n        },\n    ),\n)\n```\n```csharp\nusing System.Collections.Generic;\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V2;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var values = new Dictionary<string, object>\n        {\n            [\"controller\"] = new Dictionary<string, object>\n            {\n                [\"metrics\"] = new Dictionary<string, object>\n                {\n                    [\"enabled\"] = true\n                }\n            },\n        };\n\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            },\n            Values = values,\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tChart:   pulumi.String(\"nginx-ingress\"),\n\t\t\tVersion: pulumi.String(\"1.24.4\"),\n\t\t\tFetchArgs: helm.FetchArgs{\n\t\t\t\tRepo: pulumi.String(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t\tValues: pulumi.Map{\n\t\t\t\t\"controller\": pulumi.Map{\n\t\t\t\t\t\"metrics\": pulumi.Map{\n\t\t\t\t\t\t\"enabled\": pulumi.Bool(true),\n\t\t\t\t\t},\n\t\t\t\t},\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Deploy Chart into Namespace\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v2.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    namespace: \"test-namespace\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v2 import Chart, ChartOpts, FetchOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        namespace=\"test-namespace\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n    ),\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V2;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            Namespace = \"test-namespace\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            },\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tChart:     pulumi.String(\"nginx-ingress\"),\n\t\t\tVersion:   pulumi.String(\"1.24.4\"),\n\t\t\tNamespace: pulumi.String(\"test-namespace\"),\n\t\t\tFetchArgs: helm.FetchArgs{\n\t\t\t\tRepo: pulumi.String(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Chart with Transformations\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v2.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n    transformations: [\n        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Service\" && obj.apiVersion === \"v1\") {\n                if (obj.spec && obj.spec.type && obj.spec.type === \"LoadBalancer\") {\n                    obj.spec.type = \"ClusterIP\";\n                }\n            }\n        },\n\n        // Set a resource alias for a previous name.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Deployment\") {\n                opts.aliases = [{ name: \"oldName\" }]\n            }\n        },\n\n        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Pod\" && obj.metadata.name === \"test\") {\n                obj.apiVersion = \"v1\"\n                obj.kind = \"List\"\n            }\n        },\n    ],\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v2 import Chart, ChartOpts, FetchOpts\n\n# Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\ndef make_service_private(obj, opts):\n    if obj[\"kind\"] == \"Service\" and obj[\"apiVersion\"] == \"v1\":\n        try:\n            t = obj[\"spec\"][\"type\"]\n            if t == \"LoadBalancer\":\n                obj[\"spec\"][\"type\"] = \"ClusterIP\"\n        except KeyError:\n            pass\n\n\n# Set a resource alias for a previous name.\ndef alias(obj, opts):\n    if obj[\"kind\"] == \"Deployment\":\n        opts.aliases = [\"oldName\"]\n\n\n# Omit a resource from the Chart by transforming the specified resource definition to an empty List.\ndef omit_resource(obj, opts):\n    if obj[\"kind\"] == \"Pod\" and obj[\"metadata\"][\"name\"] == \"test\":\n        obj[\"apiVersion\"] = \"v1\"\n        obj[\"kind\"] = \"List\"\n\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n        transformations=[make_service_private, alias, omit_resource],\n    ),\n)\n```\n```csharp\nusing System.Collections.Generic;\nusing System.Collections.Immutable;\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V2;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            },\n            Transformations =\n            {\n                LoadBalancerToClusterIP,\n                ResourceAlias,\n                OmitTestPod,\n            }\n\n        });\n\n        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\n        ImmutableDictionary<string, object> LoadBalancerToClusterIP(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            if ((string)obj[\"kind\"] == \"Service\" && (string)obj[\"apiVersion\"] == \"v1\")\n            {\n                var spec = (ImmutableDictionary<string, object>)obj[\"spec\"];\n                if (spec != null && (string)spec[\"type\"] == \"LoadBalancer\")\n                {\n                    return obj.SetItem(\"spec\", spec.SetItem(\"type\", \"ClusterIP\"));\n                }\n            }\n\n            return obj;\n        }\n\n        // Set a resource alias for a previous name.\n        ImmutableDictionary<string, object> ResourceAlias(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            if ((string)obj[\"kind\"] == \"Deployment\")\n            {\n                opts.Aliases.Add(new Alias { Name = \"oldName\" });\n            }\n\n            return obj;\n        }\n\n        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.\n        ImmutableDictionary<string, object> OmitTestPod(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            var metadata = (ImmutableDictionary<string, object>)obj[\"metadata\"];\n            if ((string)obj[\"kind\"] == \"Pod\" && (string)metadata[\"name\"] == \"test\")\n            {\n                return new Dictionary<string, object>\n                {\n                    [\"apiVersion\"] = \"v1\",\n                    [\"kind\"] = \"List\",\n                    [\"items\"] = new Dictionary<string, object>(),\n                }.ToImmutableDictionary();\n            }\n\n            return obj;\n        }\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2\"\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tChart:   pulumi.String(\"nginx-ingress\"),\n\t\t\tVersion: pulumi.String(\"1.24.4\"),\n\t\t\tFetchArgs: helm.FetchArgs{\n\t\t\t\tRepo: pulumi.String(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t\tTransformations: []yaml.Transformation{\n\t\t\t\t// Make every service private to the cluster, i.e., turn all services into ClusterIP\n\t\t\t\t// instead of LoadBalancer.\n\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\tif state[\"kind\"] == \"Service\" {\n\t\t\t\t\t\tspec := state[\"spec\"].(map[string]interface{})\n\t\t\t\t\t\tspec[\"type\"] = \"ClusterIP\"\n\t\t\t\t\t}\n\t\t\t\t},\n\n\t\t\t\t// Set a resource alias for a previous name.\n\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\tif state[\"kind\"] == \"Deployment\" {\n\t\t\t\t\t\taliases := pulumi.Aliases([]pulumi.Alias{\n\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\tName: pulumi.String(\"oldName\"),\n\t\t\t\t\t\t\t},\n\t\t\t\t\t\t})\n\t\t\t\t\t\topts = append(opts, aliases)\n\t\t\t\t\t}\n\t\t\t\t},\n\n\t\t\t\t// Omit a resource from the Chart by transforming the specified resource definition\n\t\t\t\t// to an empty List.\n\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\tname := state[\"metadata\"].(map[string]interface{})[\"name\"]\n\t\t\t\t\tif state[\"kind\"] == \"Pod\" && name == \"test\" {\n\t\t\t\t\t\tstate[\"apiVersion\"] = \"core/v1\"\n\t\t\t\t\t\tstate[\"kind\"] = \"List\"\n\t\t\t\t\t}\n\t\t\t\t},\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% /examples %}}\n\n\n\n",
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Resources created by the Chart.",
			},
			"urn": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "The stable logical URN used to distinctly address a resource, both before and after deployments.",
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
				Ref: "#/types/kubernetes:helm.sh/v2:FetchOpts",
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

var helmV3ChartResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: "Chart is a component representing a collection of resources described by an arbitrary Helm Chart.\n\nThe Chart can be fetched from any source that is accessible to the `helm` command line. Values in the `values.yml` file can be overridden using `ChartOpts.values` (equivalent to `--set` or having multiple `values.yml` files). Objects can be transformed arbitrarily by supplying callbacks to `ChartOpts.transformations`.\n\n`Chart` does not use Tiller. The Chart specified is copied and expanded locally; the semantics are equivalent to running `helm template` and then using Pulumi to manage the resulting YAML manifests. Any values that would be retrieved in-cluster are assigned fake values, and none of Tiller's server-side validity testing is executed.\n\n{{% examples %}}\n## Example Usage\n{{% example %}}\n### Local Chart Directory\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v3.Chart(\"nginx-ingress\", {\n    path: \"./nginx-ingress\",\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v3 import Chart, LocalChartOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    LocalChartOpts(\n        path=\"./nginx-ingress\",\n    ),\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V3;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new LocalChartArgs\n        {\n            Path = \"./nginx-ingress\",\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tPath: pulumi.String(\"./nginx-ingress\"),\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Remote Chart\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v3.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n    ),\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V3;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            }\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.run(func(ctx *pulumi.context) error {\n\t\t_, err := helm.newchart(ctx, \"nginx-ingress\", helm.chartargs{\n\t\t\tchart:   pulumi.string(\"nginx-ingress\"),\n\t\t\tversion: pulumi.string(\"1.24.4\"),\n\t\t\tfetchargs: helm.fetchargs{\n\t\t\t\trepo: pulumi.string(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n\n{{% /example %}}\n{{% example %}}\n### Set Chart Values\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v3.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n    values: {\n        controller: {\n            metrics: {\n                enabled: true,\n            }\n        }\n    },\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n        values={\n            \"controller\": {\n                \"metrics\": {\n                    \"enabled\": True,\n                },\n            },\n        },\n    ),\n)\n```\n```csharp\nusing System.Collections.Generic;\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V3;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var values = new Dictionary<string, object>\n        {\n            [\"controller\"] = new Dictionary<string, object>\n            {\n                [\"metrics\"] = new Dictionary<string, object>\n                {\n                    [\"enabled\"] = true\n                }\n            },\n        };\n\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            },\n            Values = values,\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tChart:   pulumi.String(\"nginx-ingress\"),\n\t\t\tVersion: pulumi.String(\"1.24.4\"),\n\t\t\tFetchArgs: helm.FetchArgs{\n\t\t\t\tRepo: pulumi.String(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t\tValues: pulumi.Map{\n\t\t\t\t\"controller\": pulumi.Map{\n\t\t\t\t\t\"metrics\": pulumi.Map{\n\t\t\t\t\t\t\"enabled\": pulumi.Bool(true),\n\t\t\t\t\t},\n\t\t\t\t},\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Deploy Chart into Namespace\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v3.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    namespace: \"test-namespace\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        namespace=\"test-namespace\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n    ),\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V3;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            Namespace = \"test-namespace\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            },\n        });\n\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tChart:     pulumi.String(\"nginx-ingress\"),\n\t\t\tVersion:   pulumi.String(\"1.24.4\"),\n\t\t\tNamespace: pulumi.String(\"test-namespace\"),\n\t\t\tFetchArgs: helm.FetchArgs{\n\t\t\t\tRepo: pulumi.String(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Chart with Transformations\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst nginxIngress = new k8s.helm.v3.Chart(\"nginx-ingress\", {\n    chart: \"nginx-ingress\",\n    version: \"1.24.4\",\n    fetchOpts:{\n        repo: \"https://charts.helm.sh/stable//\",\n    },\n    transformations: [\n        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Service\" && obj.apiVersion === \"v1\") {\n                if (obj.spec && obj.spec.type && obj.spec.type === \"LoadBalancer\") {\n                    obj.spec.type = \"ClusterIP\";\n                }\n            }\n        },\n\n        // Set a resource alias for a previous name.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Deployment\") {\n                opts.aliases = [{ name: \"oldName\" }]\n            }\n        },\n\n        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Pod\" && obj.metadata.name === \"test\") {\n                obj.apiVersion = \"v1\"\n                obj.kind = \"List\"\n            }\n        },\n    ],\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts\n\n# Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\ndef make_service_private(obj, opts):\n    if obj[\"kind\"] == \"Service\" and obj[\"apiVersion\"] == \"v1\":\n        try:\n            t = obj[\"spec\"][\"type\"]\n            if t == \"LoadBalancer\":\n                obj[\"spec\"][\"type\"] = \"ClusterIP\"\n        except KeyError:\n            pass\n\n\n# Set a resource alias for a previous name.\ndef alias(obj, opts):\n    if obj[\"kind\"] == \"Deployment\":\n        opts.aliases = [\"oldName\"]\n\n\n# Omit a resource from the Chart by transforming the specified resource definition to an empty List.\ndef omit_resource(obj, opts):\n    if obj[\"kind\"] == \"Pod\" and obj[\"metadata\"][\"name\"] == \"test\":\n        obj[\"apiVersion\"] = \"v1\"\n        obj[\"kind\"] = \"List\"\n\n\nnginx_ingress = Chart(\n    \"nginx-ingress\",\n    ChartOpts(\n        chart=\"nginx-ingress\",\n        version=\"1.24.4\",\n        fetch_opts=FetchOpts(\n            repo=\"https://charts.helm.sh/stable//\",\n        ),\n        transformations=[make_service_private, alias, omit_resource],\n    ),\n)\n```\n```csharp\nusing System.Collections.Generic;\nusing System.Collections.Immutable;\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Helm;\nusing Pulumi.Kubernetes.Helm.V3;\n\nclass HelmStack : Stack\n{\n    public HelmStack()\n    {\n        var nginx = new Chart(\"nginx-ingress\", new ChartArgs\n        {\n            Chart = \"nginx-ingress\",\n            Version = \"1.24.4\",\n            FetchOptions = new ChartFetchArgs\n            {\n                Repo = \"https://charts.helm.sh/stable//\"\n            },\n            Transformations =\n            {\n                LoadBalancerToClusterIP,\n                ResourceAlias,\n                OmitTestPod,\n            }\n\n        });\n\n        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\n        ImmutableDictionary<string, object> LoadBalancerToClusterIP(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            if ((string)obj[\"kind\"] == \"Service\" && (string)obj[\"apiVersion\"] == \"v1\")\n            {\n                var spec = (ImmutableDictionary<string, object>)obj[\"spec\"];\n                if (spec != null && (string)spec[\"type\"] == \"LoadBalancer\")\n                {\n                    return obj.SetItem(\"spec\", spec.SetItem(\"type\", \"ClusterIP\"));\n                }\n            }\n\n            return obj;\n        }\n\n        // Set a resource alias for a previous name.\n        ImmutableDictionary<string, object> ResourceAlias(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            if ((string)obj[\"kind\"] == \"Deployment\")\n            {\n                opts.Aliases.Add(new Alias { Name = \"oldName\" });\n            }\n\n            return obj;\n        }\n\n        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.\n        ImmutableDictionary<string, object> OmitTestPod(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            var metadata = (ImmutableDictionary<string, object>)obj[\"metadata\"];\n            if ((string)obj[\"kind\"] == \"Pod\" && (string)metadata[\"name\"] == \"test\")\n            {\n                return new Dictionary<string, object>\n                {\n                    [\"apiVersion\"] = \"v1\",\n                    [\"kind\"] = \"List\",\n                    [\"items\"] = new Dictionary<string, object>(),\n                }.ToImmutableDictionary();\n            }\n\n            return obj;\n        }\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3\"\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := helm.NewChart(ctx, \"nginx-ingress\", helm.ChartArgs{\n\t\t\tChart:   pulumi.String(\"nginx-ingress\"),\n\t\t\tVersion: pulumi.String(\"1.24.4\"),\n\t\t\tFetchArgs: helm.FetchArgs{\n\t\t\t\tRepo: pulumi.String(\"https://charts.helm.sh/stable//\"),\n\t\t\t},\n\t\t\tTransformations: []yaml.Transformation{\n\t\t\t\t// Make every service private to the cluster, i.e., turn all services into ClusterIP\n\t\t\t\t// instead of LoadBalancer.\n\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\tif state[\"kind\"] == \"Service\" {\n\t\t\t\t\t\tspec := state[\"spec\"].(map[string]interface{})\n\t\t\t\t\t\tspec[\"type\"] = \"ClusterIP\"\n\t\t\t\t\t}\n\t\t\t\t},\n\n\t\t\t\t// Set a resource alias for a previous name.\n\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\tif state[\"kind\"] == \"Deployment\" {\n\t\t\t\t\t\taliases := pulumi.Aliases([]pulumi.Alias{\n\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\tName: pulumi.String(\"oldName\"),\n\t\t\t\t\t\t\t},\n\t\t\t\t\t\t})\n\t\t\t\t\t\topts = append(opts, aliases)\n\t\t\t\t\t}\n\t\t\t\t},\n\n\t\t\t\t// Omit a resource from the Chart by transforming the specified resource definition\n\t\t\t\t// to an empty List.\n\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\tname := state[\"metadata\"].(map[string]interface{})[\"name\"]\n\t\t\t\t\tif state[\"kind\"] == \"Pod\" && name == \"test\" {\n\t\t\t\t\t\tstate[\"apiVersion\"] = \"core/v1\"\n\t\t\t\t\t\tstate[\"kind\"] = \"List\"\n\t\t\t\t\t}\n\t\t\t\t},\n\t\t\t},\n\t\t})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% /examples %}}\n\n",
		Properties: map[string]pschema.PropertySpec{
			"resources": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "Resources created by the Chart.",
			},
			"urn": {
				TypeSpec: pschema.TypeSpec{
					Type: "string",
				},
				Description: "The stable logical URN used to distinctly address a resource, both before and after deployments.",
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

var helmV3Release = pschema.ComplexTypeSpec{
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
}

var helmV2FetchOpts = pschema.ComplexTypeSpec{
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
		},
		Type: "object",
	},
}

var helmReleaseSettings = pschema.ComplexTypeSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		Description: "BETA FEATURE - Options to configure the Helm Release resource.",
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
			"suppressBetaWarning": {
				DefaultInfo: &pschema.DefaultSpec{
					Environment: []string{
						"PULUMI_K8S_SUPPRESS_HELM_RELEASE_BETA_WARNING",
					},
				},
				Description: "While Helm Release provider is in beta, by default 'pulumi up' will log a warning if the resource is used. If present and set to \"true\", this warning is omitted.",
				TypeSpec:    pschema.TypeSpec{Type: "boolean"},
			},
		},
		Type: "object",
	},
}

var helmV3ReleaseResource = pschema.ResourceSpec{
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
	},
}

var kustomizeDirectoryResource = pschema.ResourceSpec{
	IsComponent: true,
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		IsOverlay:   true,
		Description: "Directory is a component representing a collection of resources described by a kustomize directory (kustomization).\n\n{{% examples %}}\n## Example Usage\n{{% example %}}\n### Local Kustomize Directory\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst helloWorld = new k8s.kustomize.Directory(\"helloWorldLocal\", {\n    directory: \"./helloWorld\",\n});\n```\n```python\nfrom pulumi_kubernetes.kustomize import Directory\n\nhello_world = Directory(\n    \"hello-world-local\",\n    directory=\"./helloWorld\",\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Kustomize;\n\nclass KustomizeStack : Stack\n{\n    public KustomizeStack()\n    {\n        var helloWorld = new Directory(\"helloWorldLocal\", new DirectoryArgs\n        {\n            Directory = \"./helloWorld\",\n        });\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/kustomize\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := kustomize.NewDirectory(ctx, \"helloWorldLocal\",\n\t\t\tkustomize.DirectoryArgs{\n\t\t\t\tDirectory: pulumi.String(\"./helloWorld\"),\n\t\t\t},\n\t\t)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Kustomize Directory from a Git Repo\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst helloWorld = new k8s.kustomize.Directory(\"helloWorldRemote\", {\n    directory: \"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\",\n});\n```\n```python\nfrom pulumi_kubernetes.kustomize import Directory\n\nhello_world = Directory(\n    \"hello-world-remote\",\n    directory=\"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\",\n)\n```\n```csharp\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Kustomize;\n\nclass KustomizeStack : Stack\n{\n    public KustomizeStack()\n    {\n        var helloWorld = new Directory(\"helloWorldRemote\", new DirectoryArgs\n        {\n            Directory = \"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\",\n        });\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/kustomize\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := kustomize.NewDirectory(ctx, \"helloWorldRemote\",\n\t\t\tkustomize.DirectoryArgs{\n\t\t\t\tDirectory: pulumi.String(\"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\"),\n\t\t\t},\n\t\t)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% example %}}\n### Kustomize Directory with Transformations\n\n```typescript\nimport * as k8s from \"@pulumi/kubernetes\";\n\nconst helloWorld = new k8s.kustomize.Directory(\"helloWorldRemote\", {\n    directory: \"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\",\n    transformations: [\n        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Service\" && obj.apiVersion === \"v1\") {\n                if (obj.spec && obj.spec.type && obj.spec.type === \"LoadBalancer\") {\n                    obj.spec.type = \"ClusterIP\";\n                }\n            }\n        },\n\n        // Set a resource alias for a previous name.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Deployment\") {\n                opts.aliases = [{ name: \"oldName\" }]\n            }\n        },\n\n        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.\n        (obj: any, opts: pulumi.CustomResourceOptions) => {\n            if (obj.kind === \"Pod\" && obj.metadata.name === \"test\") {\n                obj.apiVersion = \"v1\"\n                obj.kind = \"List\"\n            }\n        },\n    ],\n});\n```\n```python\nfrom pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts\n\n# Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\ndef make_service_private(obj, opts):\n    if obj[\"kind\"] == \"Service\" and obj[\"apiVersion\"] == \"v1\":\n        try:\n            t = obj[\"spec\"][\"type\"]\n            if t == \"LoadBalancer\":\n                obj[\"spec\"][\"type\"] = \"ClusterIP\"\n        except KeyError:\n            pass\n\n\n# Set a resource alias for a previous name.\ndef alias(obj, opts):\n    if obj[\"kind\"] == \"Deployment\":\n        opts.aliases = [\"oldName\"]\n\n\n# Omit a resource from the Chart by transforming the specified resource definition to an empty List.\ndef omit_resource(obj, opts):\n    if obj[\"kind\"] == \"Pod\" and obj[\"metadata\"][\"name\"] == \"test\":\n        obj[\"apiVersion\"] = \"v1\"\n        obj[\"kind\"] = \"List\"\n\n\nhello_world = Directory(\n    \"hello-world-remote\",\n    directory=\"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\",\n    transformations=[make_service_private, alias, omit_resource],\n)\n```\n```csharp\nusing System.Collections.Generic;\nusing System.Collections.Immutable;\nusing System.Threading.Tasks;\nusing Pulumi;\nusing Pulumi.Kubernetes.Kustomize;\n\nclass KustomizeStack : Stack\n{\n    public KustomizeStack()\n    {\n        var helloWorld = new Directory(\"helloWorldRemote\", new DirectoryArgs\n        {\n            Directory = \"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\",\n            Transformations =\n              {\n                  LoadBalancerToClusterIP,\n                  ResourceAlias,\n                  OmitTestPod,\n              }\n        });\n\n        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.\n        ImmutableDictionary<string, object> LoadBalancerToClusterIP(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            if ((string)obj[\"kind\"] == \"Service\" && (string)obj[\"apiVersion\"] == \"v1\")\n            {\n                var spec = (ImmutableDictionary<string, object>)obj[\"spec\"];\n                if (spec != null && (string)spec[\"type\"] == \"LoadBalancer\")\n                {\n                    return obj.SetItem(\"spec\", spec.SetItem(\"type\", \"ClusterIP\"));\n                }\n            }\n\n            return obj;\n        }\n\n        // Set a resource alias for a previous name.\n        ImmutableDictionary<string, object> ResourceAlias(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            if ((string)obj[\"kind\"] == \"Deployment\")\n            {\n                opts.Aliases.Add(new Alias { Name = \"oldName\" });\n            }\n\n            return obj;\n        }\n\n        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.\n        ImmutableDictionary<string, object> OmitTestPod(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)\n        {\n            var metadata = (ImmutableDictionary<string, object>)obj[\"metadata\"];\n            if ((string)obj[\"kind\"] == \"Pod\" && (string)metadata[\"name\"] == \"test\")\n            {\n                return new Dictionary<string, object>\n                {\n                    [\"apiVersion\"] = \"v1\",\n                    [\"kind\"] = \"List\",\n                    [\"items\"] = new Dictionary<string, object>(),\n                }.ToImmutableDictionary();\n            }\n\n            return obj;\n        }\n    }\n}\n```\n```go\npackage main\n\nimport (\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/kustomize\"\n\t\"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml\"\n\t\"github.com/pulumi/pulumi/sdk/v3/go/pulumi\"\n)\n\nfunc main() {\n\tpulumi.Run(func(ctx *pulumi.Context) error {\n\t\t_, err := kustomize.NewDirectory(ctx, \"helloWorldRemote\",\n\t\t\tkustomize.DirectoryArgs{\n\t\t\t\tDirectory: pulumi.String(\"https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld\"),\n\t\t\t\tTransformations: []yaml.Transformation{\n\t\t\t\t\t// Make every service private to the cluster, i.e., turn all services into ClusterIP\n\t\t\t\t\t// instead of LoadBalancer.\n\t\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\t\tif state[\"kind\"] == \"Service\" {\n\t\t\t\t\t\t\tspec := state[\"spec\"].(map[string]interface{})\n\t\t\t\t\t\t\tspec[\"type\"] = \"ClusterIP\"\n\t\t\t\t\t\t}\n\t\t\t\t\t},\n\n\t\t\t\t\t// Set a resource alias for a previous name.\n\t\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\t\tif state[\"kind\"] == \"Deployment\" {\n\t\t\t\t\t\t\taliases := pulumi.Aliases([]pulumi.Alias{\n\t\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\t\tName: pulumi.String(\"oldName\"),\n\t\t\t\t\t\t\t\t},\n\t\t\t\t\t\t\t})\n\t\t\t\t\t\t\topts = append(opts, aliases)\n\t\t\t\t\t\t}\n\t\t\t\t\t},\n\n\t\t\t\t\t// Omit a resource from the Chart by transforming the specified resource definition\n\t\t\t\t\t// to an empty List.\n\t\t\t\t\tfunc(state map[string]interface{}, opts ...pulumi.ResourceOption) {\n\t\t\t\t\t\tname := state[\"metadata\"].(map[string]interface{})[\"name\"]\n\t\t\t\t\t\tif state[\"kind\"] == \"Pod\" && name == \"test\" {\n\t\t\t\t\t\t\tstate[\"apiVersion\"] = \"core/v1\"\n\t\t\t\t\t\t\tstate[\"kind\"] = \"List\"\n\t\t\t\t\t\t}\n\t\t\t\t\t},\n\t\t\t\t},\n\t\t\t},\n\t\t)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\treturn nil\n\t})\n}\n```\n{{% /example %}}\n{{% /examples %}}\n",
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

//go:embed examples/configFile.md
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

//go:embed examples/configGroup.md
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

func init() {
	typeOverlays["kubernetes:core/v1:ServiceSpec"] = serviceSpec
	typeOverlays["kubernetes:core/v1:ServiceSpecType"] = serviceSpecType
	typeOverlays["kubernetes:helm.sh/v2:FetchOpts"] = helmV2FetchOpts
	typeOverlays["kubernetes:helm.sh/v3:FetchOpts"] = helmV2FetchOpts // v2 fetch opts are identical to v3
	typeOverlays["kubernetes:helm.sh/v3:Release"] = helmV3Release
	typeOverlays["kubernetes:helm.sh/v3:RepositoryOpts"] = helmV3RepoOpts
	typeOverlays["kubernetes:helm.sh/v3:ReleaseStatus"] = helmV3ReleaseStatus
	typeOverlays["kubernetes:index:KubeClientSettings"] = kubeClientSettings
	typeOverlays["kubernetes:index:HelmReleaseSettings"] = helmReleaseSettings

	resourceOverlays["kubernetes:apiextensions.k8s.io:CustomResource"] = apiextentionsCustomResource
	resourceOverlays["kubernetes:helm.sh/v2:Chart"] = helmV2ChartResource
	resourceOverlays["kubernetes:helm.sh/v3:Chart"] = helmV3ChartResource
	resourceOverlays["kubernetes:helm.sh/v3:Release"] = helmV3ReleaseResource
	resourceOverlays["kubernetes:kustomize:Directory"] = kustomizeDirectoryResource
	resourceOverlays["kubernetes:yaml:ConfigFile"] = yamlConfigFileResource
	resourceOverlays["kubernetes:yaml:ConfigGroup"] = yamlConfigGroupResource
}
