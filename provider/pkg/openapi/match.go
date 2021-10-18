// Copyright 2021, Pulumi Corporation.  All rights reserved.

package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/client-go/util/jsonpath"
)

// PropertiesChanged compares two versions of an object to see if any path specified in `paths` has
// been changed. Paths are specified as JSONPaths, e.g., `.spec.accessModes` refers to `{spec:
// {accessModes: {}}}`.
func PropertiesChanged(oldObj, newObj map[string]interface{}, paths []string) ([]string, error) {
	patch, err := mergePatchObj(oldObj, newObj)
	if err != nil {
		return nil, err
	}
	return PatchPropertiesChanged(patch, paths)
}

// PatchPropertiesChanged scrapes the given patch object to see if any path specified in `paths` has
// been changed. Paths are specified as JSONPaths, e.g., `.spec.accessModes` refers to `{spec:
// {accessModes: {}}}`.
func PatchPropertiesChanged(patch map[string]interface{}, paths []string) ([]string, error) {
	j := jsonpath.New("")
	matches := []string{}
	for _, path := range paths {
		j.AllowMissingKeys(false) // Explicitly handle any returned errors.
		err := j.Parse(fmt.Sprintf("{%s}", path))
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		err = j.Execute(buf, patch)
		if err != nil && strings.Contains(err.Error(), "not found") {
			continue
		}

		// If no error was returned, then this is a match.
		matches = append(matches, path)
	}

	return matches, nil
}

// mergePatchObj takes a two objects and returns an object that is the union of all
// fields that were changed (e.g., were deleted, were added, and so on) between the two.
//
// For example, say we have {a: 1, c:3} and {a:1, b:2}. This function would then return {b:2, c:3}.
//
// This is useful so that we can (e.g.) use jsonpath to see which fields were altered.
func mergePatchObj(oldObj, newObj map[string]interface{}) (map[string]interface{}, error) {
	oldJSON, err := json.Marshal(oldObj)
	if err != nil {
		return nil, err
	}

	newJSON, err := json.Marshal(newObj)
	if err != nil {
		return nil, err
	}

	patchBytes, err := jsonpatch.CreateMergePatch(oldJSON, newJSON)
	if err != nil {
		return nil, err
	}

	patch := map[string]interface{}{}
	err = json.Unmarshal(patchBytes, &patch)
	if err != nil {
		return nil, err
	}

	return patch, nil
}
