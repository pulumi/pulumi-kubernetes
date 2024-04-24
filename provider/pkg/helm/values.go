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
	"sigs.k8s.io/yaml"
)

// ValueOpts handles merging of chart values from various sources.
type ValueOpts struct {
	// ValuesFiles is a list of Helm values files encapsulated as Pulumi assets.
	ValuesFiles []pulumi.Asset
	// Values is a map of Pulumi values.
	Values map[string]any
}

// MergeValues merges the values in Helm's priority order.
func (opts *ValueOpts) MergeValues(p getter.Providers) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, asset := range opts.ValuesFiles {
		currentMap := map[string]interface{}{}

		bytes, err := readAsset(p, asset)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return nil, err
		}
		// Merge with the previous map
		base = MergeMaps(base, currentMap)
	}

	// User specified a literal value map (possibly containing assets)
	values, err := marshalValues(p, opts.Values)
	if err != nil {
		return nil, err
	}
	base = MergeMaps(base, values)

	return base, nil
}

// readAsset reads the content of a Pulumi asset.
func readAsset(p getter.Providers, asset pulumi.Asset) ([]byte, error) {
	switch {
	case asset.Text() != "":
		return []byte(asset.Text()), nil
	case asset.Path() != "":
		bytes, err := os.ReadFile(asset.Path())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file %q", asset.Path())
		}
		return bytes, nil
	case asset.URI() != "":
		u, err := url.Parse(asset.URI())
		if err != nil {
			return nil, err
		}
		g, err := p.ByScheme(u.Scheme)
		if err != nil {
			return nil, errors.Wrapf(err, "no protocol handler for uri %q", asset.URI())
		}
		data, err := g.Get(asset.URI(), getter.WithURL(asset.URI()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read uri %q", asset.URI())
		}
		return data.Bytes(), nil
	default:
		return nil, errors.New("unrecognized asset type")
	}
}

// marshalValues converts Pulumi values to Helm values.
// - Expands assets to their content (to support --set-file).
func marshalValues(p getter.Providers, a map[string]interface{}) (map[string]interface{}, error) {
	var err error
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		if v, ok := v.(map[string]interface{}); ok {
			out[k], err = marshalValues(p, v)
			if err != nil {
				return nil, err
			}
			continue
		}
		if v, ok := v.(pulumi.Asset); ok {
			bytes, err := readAsset(p, v)
			if err != nil {
				return nil, err
			}
			out[k] = string(bytes)
			continue
		}
		if _, ok := v.(pulumi.Archive); ok {
			return nil, errors.New("Archive values are not supported as a Helm value")
		}
		if _, ok := v.(pulumi.Resource); ok {
			return nil, errors.New("Resource values are not supported as a Helm value")
		}
		out[k] = v
	}
	return out, nil
}
