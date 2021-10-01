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

package provider

// Note: These types must match the types defined in the Go SDK (sdk/go/kubernetes/config/pulumiTypes.go).
// Copying the types avoids having the provider depend on the SDK.

// Options for tuning the Kubernetes client used by a Provider.
type KubeClientSettings struct {
	// Maximum burst for throttle. Default value is 10.
	Burst *int `json:"burst"`
	// Maximum queries per second (QPS) to the API server from this client. Default value is 5.
	Qps *float64 `json:"qps"`
}
