// Copyright 2021, Pulumi Corporation.  All rights reserved.

package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/deepcopy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
)

func hasComputedValue(obj *unstructured.Unstructured) bool {
	if obj == nil || obj.Object == nil {
		return false
	}

	objects := []map[string]any{obj.Object}
	var curr map[string]any

	for {
		if len(objects) == 0 {
			break
		}
		curr, objects = objects[0], objects[1:]
		for _, v := range curr {
			switch field := v.(type) {
			case resource.Computed:
				return true
			case map[string]any:
				objects = append(objects, field)
			case []any:
				for _, v := range field {
					objects = append(objects, map[string]any{"": v})
				}
			case []map[string]any:
				objects = append(objects, field...)
			}
		}
	}

	return false
}

// --------------------------------------------------------------------------
// Names and namespaces.
// --------------------------------------------------------------------------

// fqObjName returns "namespace.name"
func fqObjName(o metav1.Object) string {
	return fqName(o.GetNamespace(), o.GetName())
}

// parseFqName will parse a fully-qualified Kubernetes object name.
func parseFqName(id string) (namespace, name string) {
	split := strings.Split(id, "/")
	if len(split) == 1 {
		return "", split[0]
	}
	namespace, name = split[0], split[1]
	return
}

// fqName returns "namespace/name"
func fqName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s/%s", namespace, name)
}

// --------------------------------------------------------------------------
// Kubeconfig helpers.
// --------------------------------------------------------------------------

// parseKubeconfigPropertyValue takes a PropertyValue that possibly contains a raw kubeconfig
// (YAML or JSON) string or map and attempts to unmarshal it into a Config struct. If the property value
// is empty, an empty Config is returned.
func parseKubeconfigPropertyValue(kubeconfig resource.PropertyValue) (*clientapi.Config, error) {
	if kubeconfig.IsNull() {
		return &clientapi.Config{}, nil
	}

	var cfg []byte
	if kubeconfig.IsString() {
		cfg = []byte(kubeconfig.StringValue())
	} else if kubeconfig.IsObject() {
		raw := kubeconfig.ObjectValue().Mappable()
		jsonBytes, err := json.Marshal(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal kubeconfig: %v", err)
		}
		cfg = jsonBytes
	} else {
		return nil, fmt.Errorf("unexpected kubeconfig format, type: %v", kubeconfig.TypeString())
	}
	config, err := clientcmd.Load(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %v", err)
	}

	return config, nil
}

// getActiveClusterFromConfig gets the current cluster from a kubeconfig, accounting for provider overrides.
func getActiveClusterFromConfig(config *clientapi.Config, overrides resource.PropertyMap) *clientapi.Cluster {
	if config == nil || len(config.Clusters) == 0 {
		return &clientapi.Cluster{}
	}

	currentContext := config.CurrentContext
	if val := overrides["context"]; !val.IsNull() {
		currentContext = val.StringValue()
	}

	activeContext := config.Contexts[currentContext]
	if activeContext == nil {
		return &clientapi.Cluster{}
	}
	activeClusterName := activeContext.Cluster

	activeCluster := config.Clusters[activeClusterName]
	if val := overrides["cluster"]; !val.IsNull() {
		activeCluster = config.Clusters[val.StringValue()]
	}
	if activeCluster == nil {
		return &clientapi.Cluster{}
	}

	return activeCluster
}

// pruneMap builds a pruned map by recursively copying elements from the source map that have a matching key in the
// target map. This is useful as a preprocessing step for live resource state before comparing it to program inputs.
func pruneMap(source, target map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range source {
		valueT := reflect.TypeOf(value)

		if targetValue, ok := target[key]; ok {
			targetValueT := reflect.TypeOf(targetValue)

			if valueT == nil || targetValueT == nil || valueT != targetValueT {
				result[key] = value
				continue
			}

			switch valueT.Kind() {
			case reflect.Map:
				nestedResult := pruneMap(value.(map[string]any), targetValue.(map[string]any))
				result[key] = nestedResult
			case reflect.Slice:
				nestedResult := pruneSlice(value.([]any), targetValue.([]any))
				result[key] = nestedResult
			default:
				result[key] = value
			}
		}
	}

	return result
}

// pruneSlice builds a pruned slice by copying elements from the source slice that have a matching element in the
// target slice.
func pruneSlice(source, target []any) []any {
	result := make([]any, 0, len(target))

	// If either slice is empty, return an empty slice.
	if len(source) == 0 || len(target) == 0 {
		return result
	}

	valueT := reflect.TypeOf(source[0])
	targetValueT := reflect.TypeOf(target[0])

	// If slices are of different types, return a copy of the source.
	if valueT != targetValueT {
		return deepcopy.Copy(source).([]any)
	}

	for i, targetValue := range target {
		if i+1 > len(source) {
			break
		}
		value := source[i]

		if value == nil || targetValue == nil {
			result = append(result, value)
			continue
		}

		switch valueT.Kind() {
		case reflect.Map:
			nestedResult := pruneMap(value.(map[string]any), targetValue.(map[string]any))
			result = append(result, nestedResult)
		case reflect.Slice:
			nestedResult := pruneSlice(value.([]any), targetValue.([]any))
			result = append(result, nestedResult)
		default:
			result = append(result, value)
		}
	}

	return result
}
