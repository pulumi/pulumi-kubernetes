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
	"path/filepath"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type yamlText struct {
	Name string
	Text string
}

func getYaml(input string) ([]yamlText, error) {
	paths, err := filepath.Glob(input)
	if err != nil || len(paths) == 0 {
		// Not a valid glob; continue checking other possibilities
		paths = []string{input}
	}

	var yamlTexts []yamlText
	for _, path := range paths {
		var text string

		if isUrl(path) {
			text, err = loadFromURL(path)
		} else {
			text, err = loadFromFile(path)
		}

		if err != nil {
			return nil, err
		}

		yamlTexts = append(yamlTexts, yamlText{Name: path, Text: text})
	}

	return yamlTexts, nil
}

func isUrl(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

func loadFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", pkgerrors.Wrapf(err, "failed to fetch URL: %q", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", pkgerrors.Wrapf(err, "failed to read response from HTTP Get at URL: %q", url)
		}
		return string(bodyBytes), nil
	} else {
		return "", fmt.Errorf("HTTP Get for %q returned status: %s", url, resp.Status)
	}
}

func loadFromFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", pkgerrors.Wrapf(err, "failed to read file from path: %q", path)
	}

	return string(b), nil
}

// parseYaml parses a YAML string, and then returns a slice of untyped structs that can be marshalled into
// Pulumi RPC calls.
func parseYaml(text string) ([]interface{}, error) {
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
