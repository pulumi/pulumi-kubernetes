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
	"sort"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type yamlText struct {
	Name string
	Text string
}

// getYaml loads the contents of a variety of input sources, and returns a slice of named strings, or an error.
// The following input types are supported:
// 1. File path
// 2. File path glob (e.g., *.yaml)
// 3. URL
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

// isUrl returns true if the input string has a URL prefix, false otherwise.
func isUrl(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

// loadFromURL makes an HTTP GET request at the specified URL and returns the result as a string, or returns
// an error.
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

// loadFromFile returns the contents of the specified file as a string, or returns an error.
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
	var resources []unstructured.Unstructured

	dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(strings.NewReader(text)), 128)
	for {
		var value map[string]interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		resources = append(resources, unstructured.Unstructured{Object: value})
	}

	// Sort the resources by Kind to minimize retries on creation.
	ks := newKindSorter(resources)
	sort.Sort(ks)

	result := make([]interface{}, len(ks.resources))
	for _, resource := range ks.resources {
		result = append(result, resource.Object)
	}

	return result, nil
}

type kindSorter struct {
	ordering  map[string]int
	resources []unstructured.Unstructured
}

func newKindSorter(m []unstructured.Unstructured) *kindSorter {
	// This slice defines the sort order of resources, with higher priority given to resources that should be
	// created first. This code was derived from
	// https://github.com/helm/helm/blob/09ebcde05c5331743713af16777768de854ac972/pkg/releaseutil/kind_sorter.go
	s := []string{
		"Namespace",
		"NetworkPolicy",
		"ResourceQuota",
		"LimitRange",
		"PodSecurityPolicy",
		"PodDisruptionBudget",
		"Secret",
		"ConfigMap",
		"StorageClass",
		"PersistentVolume",
		"PersistentVolumeClaim",
		"ServiceAccount",
		"CustomResourceDefinition",
		"ClusterRole",
		"ClusterRoleList",
		"ClusterRoleBinding",
		"ClusterRoleBindingList",
		"Role",
		"RoleList",
		"RoleBinding",
		"RoleBindingList",
		"Service",
		"DaemonSet",
		"Pod",
		"ReplicationController",
		"ReplicaSet",
		"Deployment",
		"HorizontalPodAutoscaler",
		"StatefulSet",
		"Job",
		"CronJob",
		"Ingress",
		"APIService",
	}

	o := make(map[string]int, len(s))
	for v, k := range s {
		o[k] = v
	}

	return &kindSorter{
		resources: m,
		ordering:  o,
	}
}

func (k *kindSorter) Len() int { return len(k.resources) }

func (k *kindSorter) Swap(i, j int) { k.resources[i], k.resources[j] = k.resources[j], k.resources[i] }

func (k *kindSorter) Less(i, j int) bool {
	a := k.resources[i]
	b := k.resources[j]
	first, aok := k.ordering[a.GroupVersionKind().Kind]
	second, bok := k.ordering[b.GroupVersionKind().Kind]
	// if same kind (including unknown) sub sort alphanumeric
	if first == second {
		// if both are unknown and of different kind sort by kind alphabetically
		if !aok && !bok && a.GroupVersionKind().Kind != b.GroupVersionKind().Kind {
			return a.GroupVersionKind().Kind < b.GroupVersionKind().Kind
		}
		return a.GetName() < b.GetName()
	}
	// unknown kind is last
	if !aok {
		return false
	}
	if !bok {
		return true
	}
	// sort different kinds
	return first < second
}
