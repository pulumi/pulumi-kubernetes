// Copyright 2021, Pulumi Corporation.  All rights reserved.

package provider

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
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

func loadKubeconfig(pathOrContents string, overrides *clientcmd.ConfigOverrides) (clientcmd.ClientConfig, *clientapi.Config, error) {
	homeDir := func() string {
		// Ignore errors. The filepath will be checked later, so we can handle failures there.
		usr, _ := user.Current()
		return usr.HomeDir
	}

	// Note: the Python SDK was setting the kubeconfig value to "" by default, so explicitly check for empty string.
	if pathOrContents != "" {
		var contents string

		// Handle the '~' character if it is set in the config string. Normally, this would be expanded by the shell
		// into the user's home directory, but we have to do that manually if it is set in a config value.
		if pathOrContents == "~" {
			// In case of "~", which won't be caught by the "else if"
			pathOrContents = homeDir()
		} else if strings.HasPrefix(pathOrContents, "~/") {
			pathOrContents = filepath.Join(homeDir(), pathOrContents[2:])
		}

		// If the variable is a valid filepath, load the file and parse the contents as a k8s config.
		_, err := os.Stat(pathOrContents)
		if err == nil || filepath.IsAbs(pathOrContents) || strings.HasPrefix(pathOrContents, ".") {
			b, err := os.ReadFile(pathOrContents)
			if err != nil {
				return nil, nil, err
			}
			contents = string(b)
		} else { // Assume the contents are a k8s config.
			contents = pathOrContents
		}

		// Load the contents of the k8s config.
		apiConfig, err := clientcmd.Load([]byte(contents))
		if err != nil {
			return nil, nil, err
		}
		kubeconfig := clientcmd.NewDefaultClientConfig(*apiConfig, overrides)

		// double-check that the kubeconfig is semantically valid w.r.t. context and cluster configuration.
		_, err = kubeconfig.ClientConfig()
		if err != nil {
			return nil, nil, err
		}
		return kubeconfig, apiConfig, nil
	}

	// Use client-go to resolve the final configuration values for the client. Typically, these
	// values would reside in the $KUBECONFIG file, but can also be altered in several
	// places, including in env variables, client-go default values, and (if we allowed it) CLI
	// flags.
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	apiConfig, err := kubeconfig.RawConfig()
	if err != nil {
		return nil, nil, err
	}
	// double-check that the kubeconfig is semantically valid w.r.t. context and cluster configuration.
	_, err = kubeconfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	return kubeconfig, &apiConfig, nil
}

// pruneMap builds a pruned map by recursively copying elements from the source map that have a matching key in the
// target map. This is useful as a preprocessing step for live resource state before comparing it to program inputs.
func pruneMap(source, target map[string]any) map[string]any {
	// If either map is nil, return nil.
	if target == nil || source == nil {
		return nil
	}

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
	// If either slice is nil, return nil.
	if target == nil || source == nil {
		return nil
	}

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
			if nestedResult != nil {
				result = append(result, nestedResult)
			}
		case reflect.Slice:
			nestedResult := pruneSlice(value.([]any), targetValue.([]any))
			if nestedResult != nil {
				result = append(result, nestedResult)
			}
		default:
			result = append(result, value)
		}
	}

	return result
}
