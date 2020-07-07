package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	yamlExtras "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"
)

// ResourceDefinition contains the minimum ammount of information of a
// Kubernetes CRD to generate appropiate constructors/interfaces/classes.
type ResourceDefinition struct {
	json        string
	kind        string
	versionName string
	group       string
	fields      []Field
}

func (r ResourceDefinition) apiVersion() string {
	return fmt.Sprintf("%s/%s", r.group, r.versionName)
}

func (r ResourceDefinition) argsName() string {
	return r.kind + "Args"
}

func (r ResourceDefinition) specArgsName() string {
	return r.kind + "SpecArgs"
}

func (r ResourceDefinition) definitionName() string {
	return r.kind + "Definition"
}

type Field struct {
	name     string
	openType string
	required bool
}

// This is a utility function meant to linearly scan and check if a string is
// contained within a string slice. Since the CRD spec provides us with a list
// of required fields, we need to call this function for every field we process
// to check if it is required. This makes it O(n^2) time then, so if it turns
// out that some CRD required field lists are extremely large, then I'll
// implement some set class to get it down to O(n).
func contains(list []string, itemToFind string) bool {
	for _, item := range list {
		if item == itemToFind {
			return true
		}
	}
	return false
}

// NewResourceDefinition returns a ResourceDefinition by parsing out all
// required information from a NestedMap.
func NewResourceDefinition(yamlFile []byte) ResourceDefinition {
	var m map[interface{}]interface{}
	err := yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	crd := NestedMap{m}
	rawJSON, err := yamlExtras.YAMLToJSON(yamlFile)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, rawJSON, "\t\t", "\t"); err != nil {
		log.Fatalf("err: %v\n", err)
	}
	json := dst.String()

	kind := crd.Get("spec", "names", "kind").String()
	versionName := crd.Get("spec", "versions").Index(0).Get("name").String()
	group := crd.Get("spec", "group").String()

	spec := crd.Get("spec", "versions").Index(0).Get("schema", "openAPIV3Schema", "properties", "spec")
	requiredFields := spec.Get("required").StringArray()
	properties := spec.Get("properties").Map()
	fields := make([]Field, 0, len(properties))

	for k, v := range properties {
		name := k.(string)
		openType := v.(map[interface{}]interface{})["type"].(string)
		required := contains(requiredFields, name)
		fields = append(fields, Field{name, openType, required})
	}

	return ResourceDefinition{json, kind, versionName, group, fields}
}
