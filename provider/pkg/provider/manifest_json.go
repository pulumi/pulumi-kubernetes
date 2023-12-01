// Copyright 2021, Pulumi Corporation.  All rights reserved.

package provider

import (
	"fmt"
	"strings"

	goset "github.com/deckarep/golang-set/v2"

	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"golang.org/x/crypto/sha3"
	"helm.sh/helm/v3/pkg/releaseutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// This file is a modified version of
// https://github.com/hashicorp/terraform-provider-helm/blob/main/helm/manifest_json.go

// convertYAMLManifestToJSON converts manifests provided s a string and returns
// a deserialized map representation of the manifest (with secrets masked), a map
// grouping resource names in the manifests by group version and any error encountered.
// Not, currently only kubernetes secret data is masked.
func convertYAMLManifestToJSON(manifest string) (map[string]any, map[string][]string, error) {
	releaseResources := map[string]goset.Set[string]{}
	m := map[string]any{}

	resources := releaseutil.SplitManifests(manifest)
	for _, resource := range resources {
		obj := new(unstructured.Unstructured)
		// decode YAML into unstructured.Unstructured
		dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		_, gvk, err := dec.Decode([]byte(resource), nil, obj)
		if err != nil {
			if runtime.IsMissingKind(err) {
				// Likely empty/nil resource. Ignore.
				continue
			}
			return nil, nil, err
		}

		resKey := fmt.Sprintf("%s/%s", gvk.GroupKind().String(), obj.GetAPIVersion())
		resVal, has := releaseResources[resKey]
		if !has {
			resVal = goset.NewSet[string]()
		}
		resName := obj.GetName()
		if namespace := obj.GetNamespace(); namespace != "" {
			resName = fmt.Sprintf("%s/%s", namespace, resName)
		}
		resVal.Add(resName)
		releaseResources[resKey] = resVal
		key := fmt.Sprintf("%s/%s/%s", strings.ToLower(gvk.GroupKind().String()),
			obj.GetAPIVersion(),
			obj.GetName())

		if namespace := obj.GetNamespace(); namespace != "" {
			key = fmt.Sprintf("%s/%s", namespace, key)
		}

		var o any = &obj.Object
		if gvk.Kind == "Secret" {
			var secret corev1.Secret
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &secret)
			if err != nil {
				return nil, nil, err
			}

			for k, v := range secret.Data {
				h := hashSensitiveValue(string(v))
				secret.Data[k] = []byte(h)
			}
			o = &secret
		}
		unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(o)
		if err != nil {
			return nil, nil, err
		}
		m[key] = unstructured
	}

	logger.V(9).Infof("Manifest: %#v", m)

	releaseResourcesGrouping := map[string][]string{}
	for k, v := range releaseResources {
		releaseResourcesGrouping[k] = goset.Sorted(v)
	}
	return m, releaseResourcesGrouping, nil
}

// hashSensitiveValue creates a hash of a sensitive value and returns the string
// "(sensitive value xxxxxxxx)". We have to do this because helm release manifests
// may end up embedding secrets. This allows us to try and render the manifests
// while avoiding the possibility of leaking sensitive values.
func hashSensitiveValue(v string) string {
	hash := make([]byte, 8)
	sha3.ShakeSum256(hash, []byte(v))
	return fmt.Sprintf("(sensitive value %x)", hash)
}
