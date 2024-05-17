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

package provider

// Note: These types must match the types defined in the Go SDK (sdk/go/kubernetes/pulumiTypes.go).
// Copying the types avoids having the provider depend on the SDK.

// Options to configure the Helm Release resource.
type HelmReleaseSettings struct {
	// The backend storage driver for Helm. Values are: configmap, secret, memory, sql.
	Driver *string `json:"driver"`
	// The path to the helm plugins directory.
	PluginsPath *string `json:"pluginsPath"`
	// The path to the registry config file.
	RegistryConfigPath *string `json:"registryConfigPath"`
	// The path to the directory containing cached repository indexes.
	RepositoryCache *string `json:"repositoryCache"`
	// The path to the file containing repository names and URLs.
	RepositoryConfigPath *string `json:"repositoryConfigPath"`
}

// Options for tuning the Kubernetes client used by a Provider.
type KubeClientSettings struct {
	// Maximum burst for throttle. Default value is 120.
	Burst *int `json:"burst"`
	// Maximum queries per second (QPS) to the API server from this client. Default value is 50.
	QPS *float64 `json:"qps"`
	// Maximum time in seconds to wait before cancelling a HTTP request to the Kubernetes server. Default value is 32.
	Timeout *int `json:"timeout"`
}
