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
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"helm.sh/helm/v3/pkg/getter"
)

func LocateKeyring(p getter.Providers, asset pulumi.Asset) (string, error) {

	makeTemp := func(data []byte) (string, error) {
		file, err := os.CreateTemp("", "keyring")
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := file.Write(data); err != nil {
			return "", err
		}
		return file.Name(), err
	}

	switch {
	case asset.Text() != "":
		return makeTemp([]byte(asset.Text()))
	case asset.Path() != "":
		return asset.Path(), nil
	case asset.URI() != "":
		u, err := url.Parse(asset.URI())
		if err != nil {
			return "", err
		}
		g, err := p.ByScheme(u.Scheme)
		if err != nil {
			return "", errors.Wrapf(err, "no protocol handler for uri %q", asset.URI())
		}
		data, err := g.Get(asset.URI(), getter.WithURL(asset.URI()))
		if err != nil {
			return "", errors.Wrapf(err, "failed to read uri %q", asset.URI())
		}
		return makeTemp(data.Bytes())
	default:
		return "", errors.New("unrecognized asset type")
	}
}
