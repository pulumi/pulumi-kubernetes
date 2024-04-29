// Copyright 2021, Pulumi Corporation.  All rights reserved.

package provider

import (
	"encoding/json"
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

// annotateComputed copies the "computedness" from the ins to the outs. If there are values with the same keys for the
// outs and the ins, if they are both objects, they are transformed recursively. Likewise for arrays.
// Otherwise, if the value in the ins is computed, the value in the outs is marked as computed.
// If the value in the ins is a secret, its underlying value is checked for computedness.
func annotateComputed(outs, ins resource.PropertyMap) {
	if outs == nil || ins == nil {
		return
	}
	for key, inValue := range ins {
		outValue, has := outs[key]
		if !has {
			continue
		}
		outs[key] = annotateComputedValue(outValue, inValue)
	}
}

func annotateComputedValue(outValue, inValue resource.PropertyValue) resource.PropertyValue {
	if inValue.IsSecret() {
		return annotateComputedValue(outValue, inValue.SecretValue().Element)
	}
	if !outValue.IsComputed() && inValue.IsComputed() {
		return resource.MakeComputed(resource.NewStringProperty(""))
	}
	if outValue.IsObject() && inValue.IsObject() {
		annotateComputed(outValue.ObjectValue(), inValue.ObjectValue())
	} else if outValue.IsArray() && inValue.IsArray() {
		annotateComputedArray(outValue.ArrayValue(), inValue.ArrayValue())
	}
	return outValue
}

func annotateComputedArray(outs, ins []resource.PropertyValue) {
	for i := range ins {
		if i < len(outs) {
			outs[i] = annotateComputedValue(outs[i], ins[i])
		}
	}
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

func writeKubeconfigToFile(config *clientapi.Config) (string, error) {
	file, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		return "", err
	}
	err = clientcmd.WriteToFile(*config, file.Name())
	return file.Name(), err
}

// parseKubeconfigString takes a string that contains either a path to a kubeconfig file
// or the contents of a kubeconfig (YAML or JSON).
func parseKubeconfigString(pathOrContents string) (*clientapi.Config, error) {
	var contents string

	if pathOrContents == "" {
		return &clientapi.Config{}, nil
	}

	// Handle the '~' character if it is set in the config string. Normally, this would be expanded by the shell
	// into the user's home directory, but we have to do that manually if it is set in a config value.
	homeDir := func() string {
		// Ignore errors. The filepath will be checked later, so we can handle failures there.
		usr, _ := user.Current()
		return usr.HomeDir
	}
	if pathOrContents == "~" {
		// In case of "~", which won't be caught by the "else if"
		pathOrContents = homeDir()
	} else if strings.HasPrefix(pathOrContents, "~/") {
		pathOrContents = filepath.Join(homeDir(), pathOrContents[2:])
	}

	// If the variable is a valid filepath, load the file and parse the contents as a k8s config.
	_, err := os.Stat(pathOrContents)
	if err == nil {
		b, err := os.ReadFile(pathOrContents)
		if err != nil {
			return nil, err
		}
		contents = string(b)
	} else { // Assume the contents are a k8s config.
		contents = pathOrContents
	}

	return clientcmd.Load([]byte(contents))
}

// parseKubeconfigPropertyValue takes a PropertyValue that possibly contains a raw kubeconfig
// (YAML or JSON) string or map and attempts to unmarshal it into a Config struct. If the property value
// is empty, an empty Config is returned.
func parseKubeconfigPropertyValue(kubeconfig resource.PropertyValue) (*clientapi.Config, error) {
	if kubeconfig.IsNull() {
		return &clientapi.Config{}, nil
	}

	var cfg string
	if kubeconfig.IsString() {
		cfg = kubeconfig.StringValue()
	} else if kubeconfig.IsObject() {
		raw := kubeconfig.ObjectValue().Mappable()
		jsonBytes, err := json.Marshal(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal kubeconfig: %v", err)
		}
		cfg = string(jsonBytes)
	} else {
		return nil, fmt.Errorf("unexpected kubeconfig format, type: %v", kubeconfig.TypeString())
	}
	config, err := parseKubeconfigString(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %v", err)
	}
	return config, nil
}

// getActiveClusterFromConfig gets the current cluster from a kubeconfig, accounting for provider overrides.
// If the config is nil or the active cluster could not be found, false is returned.
func getActiveClusterFromConfig(config *clientapi.Config, overrides resource.PropertyMap) (*clientapi.Cluster, bool) {
	if config == nil || len(config.Clusters) == 0 {
		return &clientapi.Cluster{}, false
	}

	currentContext := config.CurrentContext
	if val := overrides["context"]; !val.IsNull() {
		currentContext = val.StringValue()
	}

	activeContext := config.Contexts[currentContext]
	if activeContext == nil {
		return &clientapi.Cluster{}, false
	}
	activeClusterName := activeContext.Cluster

	activeCluster := config.Clusters[activeClusterName]
	if val := overrides["cluster"]; !val.IsNull() {
		activeCluster = config.Clusters[val.StringValue()]
	}
	if activeCluster == nil {
		return &clientapi.Cluster{}, false
	}

	return activeCluster, true
}

// --------------------------------------------------------------------------
// Unstructured helpers.
// --------------------------------------------------------------------------

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
