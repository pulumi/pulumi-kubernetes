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

package examples

import _ "embed" // Needed to support go:embed directive

//go:embed upstream/deployment.md
var appsV1DeploymentMD string

//go:embed upstream/service.md
var coreV1ServiceMD string

// ApiVersionToExample contains Markdown-formatted examples corresponding to a k8s apiVersion.
var ApiVersionToExample = map[string]string{
	"kubernetes:apps/v1:Deployment": appsV1DeploymentMD,
	"kubernetes:core/v1:Service":    coreV1ServiceMD,
}
