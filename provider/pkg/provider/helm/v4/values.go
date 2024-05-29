// Copyright 2024, Pulumi Corporation.
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

package v4

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
)

// readValues hydrates Assets and persists values on-disk in order to provide
// them to upstream's MergeValues logic.
//
// The returned function cleans up the on-disk values and should always be
// called, even on error.
func readValues(p getter.Providers, v map[string]any, files []pulumi.Asset) (values.Options, func(), error) {
	opts := values.Options{}
	tmp, err := os.MkdirTemp(os.TempDir(), "pulumi-kubernetes")
	if err != nil {
		return opts, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tmp)
	}

	valuesFiles := make([]string, 0, len(files)+1)

	persist := func(out []byte) error {
		fname := filepath.Join(tmp, fmt.Sprintf("values-%d.yaml", len(valuesFiles)))
		err := os.WriteFile(fname, out, 0o600)
		if err != nil {
			return err
		}
		valuesFiles = append(valuesFiles, fname)
		return nil
	}

	for _, f := range files {
		out, err := readAsset(p, f)
		if err != nil {
			return opts, cleanup, err
		}
		err = persist(out)
		if err != nil {
			return opts, cleanup, err
		}
	}

	values, err := readAssets(p, v)
	if err != nil {
		return opts, cleanup, err
	}
	out, err := yaml.Marshal(values)
	if err != nil {
		return opts, cleanup, err
	}
	err = persist(out)
	if err != nil {
		return opts, cleanup, err
	}

	opts.ValueFiles = valuesFiles

	return opts, cleanup, nil
}

// readAsset reads the content of a Pulumi asset.
func readAsset(p getter.Providers, asset pulumi.Asset) ([]byte, error) {
	switch {
	case asset.Text() != "":
		return []byte(asset.Text()), nil
	case asset.Path() != "":
		bytes, err := os.ReadFile(asset.Path())
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", asset.Path(), err)
		}
		return bytes, nil
	case asset.URI() != "":
		u, err := url.Parse(asset.URI())
		if err != nil {
			return nil, err
		}
		g, err := p.ByScheme(u.Scheme)
		if err != nil {
			return nil, fmt.Errorf("no protocol handler for uri %q", asset.URI())
		}
		data, err := g.Get(asset.URI(), getter.WithURL(asset.URI()))
		if err != nil {
			return nil, fmt.Errorf("failed to read uri %q: %w", asset.URI(), err)
		}
		return data.Bytes(), nil
	default:
		return nil, errors.New("unrecognized asset type")
	}
}

// readAssets converts Pulumi values to Helm values, hydrating Asset values
// along the way.
func readAssets(p getter.Providers, a map[string]interface{}) (map[string]interface{}, error) {
	var err error
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		if v, ok := v.(map[string]interface{}); ok {
			out[k], err = readAssets(p, v)
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
