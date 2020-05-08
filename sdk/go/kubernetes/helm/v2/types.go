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

package helm

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// FetchArgs specifies arguments for fetching the Helm chart.
type FetchArgs struct {
	// Specific version of a chart. If unset, the latest version is fetched.
	Version pulumi.StringInput
	// Verify certificates of HTTPS-enabled servers using this CA bundle.
	CAFile pulumi.StringInput
	// Identify HTTPS client using this SSL certificate file.
	CertFile pulumi.StringInput
	// Identify HTTPS client using this SSL key file.
	KeyFile pulumi.StringInput
	// Location to write the chart. If Destination and UntarDir are specified, UntarDir is
	// appended to Destination (default ".").
	Destination pulumi.StringInput
	// Keyring containing public keys (default "~/.gnupg/pubring.gpg").
	Keyring pulumi.StringInput
	// Chart repository password.
	Password pulumi.StringInput
	// Chart repository URL for the requested chart.
	Repo pulumi.StringInput
	// Location to expand the chart. (default ".").
	UntarDir pulumi.StringInput
	// Chart repository username.
	Username pulumi.StringInput
	// Location of your Helm config. Overrides $HELM_HOME (default "~/.helm").
	Home pulumi.StringInput
	// Use development versions, too. Equivalent to version '>0.0.0-0'. If Version is set,
	// Devel is ignored.
	Devel pulumi.BoolPtrInput
	// Fetch the provenance file, but don't perform verification.
	Prov pulumi.BoolPtrInput
	// If false, leave the chart as a tarball after downloading.
	Untar pulumi.BoolPtrInput
	// Verify the package against its signature.
	Verify pulumi.BoolPtrInput
}

// fetchArgs is a copy of FetchArgs but without using TInput in types.
type fetchArgs struct {
	Version     string `pulumi:"version"`
	CAFile      string `pulumi:"caFile"`
	CertFile    string `pulumi:"certFile"`
	KeyFile     string `pulumi:"keyFile"`
	Destination string `pulumi:"destination"`
	Keyring     string `pulumi:"keyring"`
	Password    string `pulumi:"password"`
	Repo        string `pulumi:"repo"`
	UntarDir    string `pulumi:"untarDir"`
	Username    string `pulumi:"username"`
	Home        string `pulumi:"home"`
	Devel       *bool  `pulumi:"devel"`
	Prov        *bool  `pulumi:"prov"`
	Untar       *bool  `pulumi:"untar"`
	Verify      *bool  `pulumi:"verify"`
}

type FetchArgsInput interface {
	pulumi.Input

	ToFetchArgsOutput() FetchArgsOutput
	ToFetchArgsOutputWithContext(context.Context) FetchArgsOutput
}

func (FetchArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*fetchArgs)(nil)).Elem()
}

func (i FetchArgs) ToFetchArgsOutput() FetchArgsOutput {
	return i.ToFetchArgsOutputWithContext(context.Background())
}

func (i FetchArgs) ToFetchArgsOutputWithContext(ctx context.Context) FetchArgsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FetchArgsOutput)
}

type FetchArgsOutput struct{ *pulumi.OutputState }

func (FetchArgsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*fetchArgs)(nil)).Elem()
}

func (o FetchArgsOutput) ToFetchArgsOutput() FetchArgsOutput {
	return o
}

func (o FetchArgsOutput) ToFetchArgsOutputWithContext(ctx context.Context) FetchArgsOutput {
	return o
}

// ChartArgs specifies arguments for constructing a Chart resource.
type ChartArgs struct {
	// The optional Kubernetes API versions used for Capabilities.APIVersions.
	APIVersions pulumi.StringArrayInput
	// The optional namespace to install chart resources into.
	Namespace pulumi.StringInput
	// Overrides for chart values.
	Values pulumi.MapInput
	// Transformations is an optional list of transformations to apply to Kubernetes resource definitions
	// before registering with the engine.
	Transformations []yaml.Transformation
	// ResourcePrefix is an optional prefix for the auto-generated resource names. For example, a resource named `bar`
	// created with resource prefix of `"foo"` would produce a resource named `"foo-bar"`.
	ResourcePrefix string

	// (Remote chart) The repository name of the chart to deploy. Example: "stable".
	Repo pulumi.StringInput
	// (Remote chart) The name of the chart to deploy.  If Repo is specified, this chart name will be prefixed
	// by the repo name.
	// Example: Repo: "stable", Chart: "nginx-ingress" -> "stable/nginx-ingress"
	// Example: Chart: "stable/nginx-ingress" -> "stable/nginx-ingress"
	Chart pulumi.StringInput
	// (Remote chart) The version of the chart to deploy. If not provided, the latest version will be deployed.
	Version pulumi.StringInput
	// (Remote chart) Additional options to customize the fetching of the Helm chart.
	FetchArgs FetchArgsInput

	// (Local chart) The path to the chart directory which contains the `Chart.yaml` file.
	// If Path is set, any remote chart args (Repo, Chart, Version, FetchArgs) will be ignored.
	Path pulumi.StringInput
}

// chartArgs is a copy of ChartArgs but without using TInput in types.
type chartArgs struct {
	APIVersions     []string               `pulumi:"apiVersions"`
	Namespace       string                 `pulumi:"namespace"`
	Values          map[string]interface{} `pulumi:"values"`
	Transformations []yaml.Transformation  `pulumi:"transformations"`
	ResourcePrefix  string                 `pulumi:"resourcePrefix"`
	Repo            string                 `pulumi:"repo"`
	Chart           string                 `pulumi:"chart"`
	Version         string                 `pulumi:"version"`
	FetchArgs       fetchArgs              `pulumi:"fetchArgs"`
	Path            string                 `pulumi:"path"`
}

type ChartArgsInput interface {
	pulumi.Input

	ToChartArgsOutput() ChartArgsOutput
	ToChartArgsOutputWithContext(context.Context) ChartArgsOutput
}

func (ChartArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*chartArgs)(nil)).Elem()
}

func (i ChartArgs) ToChartArgsOutput() ChartArgsOutput {
	return i.ToChartArgsOutputWithContext(context.Background())
}

func (i ChartArgs) ToChartArgsOutputWithContext(ctx context.Context) ChartArgsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ChartArgsOutput)
}

type ChartArgsOutput struct{ *pulumi.OutputState }

func (ChartArgsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*chartArgs)(nil)).Elem()
}

func (o ChartArgsOutput) ToChartArgsOutput() ChartArgsOutput {
	return o
}

func (o ChartArgsOutput) ToChartArgsOutputWithContext(ctx context.Context) ChartArgsOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(FetchArgsOutput{})
	pulumi.RegisterOutputType(ChartArgsOutput{})
}
