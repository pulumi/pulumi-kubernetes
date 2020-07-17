package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UnmarshalYaml decodes a yamlFile into a Unstructured type with a value of
// map[string]interface{}.
func UnmarshalYaml(yamlFile []byte) (unstruct.Unstructured, error) {
	var value map[string]interface{}
	dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(bytes.NewReader(yamlFile)), 128)
	if err := dec.Decode(&value); err != nil {
		return unstruct.Unstructured{Object: nil}, err
	}
	return unstruct.Unstructured{Object: value}, nil
}

// NestedMapSlice returns a copy of []map[string]interface{} value of a nested field.
// Returns false if value is not found and an error if not a []interface{} or contains non-map items in the slice.
// Notice that if the value is found but not of type []interface{}, this still returns true.
// The unstructured package only had NestedSlice and NestedStringSlice, so I had to manually implement this.
func NestedMapSlice(obj map[string]interface{}, fields ...string) ([]map[string]interface{}, bool, error) {
	val, found, err := unstruct.NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	m, ok := val.([]interface{})
	if !ok {
		return nil, true, fmt.Errorf("%v accessor error: %v is of the type %T, expected []interface{}", jsonPath(fields), val, val)
	}
	mapSlice := make([]map[string]interface{}, 0, len(m))
	for _, v := range m {
		if strMap, ok := v.(map[string]interface{}); ok {
			mapSlice = append(mapSlice, strMap)
		} else {
			return nil, false, fmt.Errorf("%v accessor error: contains non-map key in the slice: %v is of the type %T, expected map[string]interface{}", jsonPath(fields), v, v)
		}
	}
	return mapSlice, true, nil
}

func jsonPath(fields []string) string {
	return "." + strings.Join(fields, ".")
}
