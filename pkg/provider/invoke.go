// Copyright 2016-2019, Pulumi Corporation.
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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// loadYaml loads a YAML document from either a file or a URL, decodes it, and then returns a slice of
// untyped structs that can be marshalled easily into Pulumi RPC calls.
func loadYaml(path string) ([]interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path for loadYaml")
	}

	var text string

	isUrl := func(path string) bool {
		return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
	}
	if isUrl(path) {
		resp, err := http.Get(path)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to fetch URL: %q", path)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, pkgerrors.Wrapf(err, "failed to read response from HTTP Get at URL: %q", path)
			}
			text = string(bodyBytes)
		} else {
			return nil, fmt.Errorf("HTTP Get for %q returned status: %s", path, resp.Status)
		}
	} else {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to read file from path: %q", path)
		}

		text = string(b)
	}

	var result []interface{}

	dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(strings.NewReader(text)), 128)
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		result = append(result, value)
	}

	return result, nil
}
