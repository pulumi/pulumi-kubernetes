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

package provider

import (
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	lumirpc "github.com/pulumi/pulumi/sdk/proto/go"

	// Load auth plugins. Removing this will likely cause compilation error.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Serve launches the gRPC server for the Pulumi Kubernetes resource provider.
func Serve(providerName, version string) {
	// Start gRPC service.
	err := provider.Main(
		providerName, func(host *provider.HostClient) (lumirpc.ResourceProviderServer, error) {
			return makeKubeProvider(host, providerName, version)
		})

	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}
