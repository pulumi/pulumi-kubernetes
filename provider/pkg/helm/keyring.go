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

package helm

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"helm.sh/helm/v3/pkg/getter"
)

// LocateKeyring locates a keyring file for Helm from the given asset.
func LocateKeyring(p getter.Providers, asset pulumi.Asset) (string, error) {
	path, _, err := downloadAsset(p, asset)
	return path, err
}
