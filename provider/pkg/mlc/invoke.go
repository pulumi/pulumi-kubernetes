package mlc

import (
	"io"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/clients"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DecodeYaml parses a YAML string, and then returns a slice of untyped structs that can be marshalled into
// Pulumi RPC calls. If a default namespace is specified, set that on the relevant decoded objects.
func DecodeYaml(text, defaultNamespace string, clientSet *clients.DynamicClientSet) ([]map[string]interface{}, error) {
	var resources []unstructured.Unstructured

	dec := yaml.NewYAMLOrJSONDecoder(io.NopCloser(strings.NewReader(text)), 128)
	for {
		var value map[string]interface{}
		if err := dec.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		resource := unstructured.Unstructured{Object: value}

		// Sometimes manifests include empty resources, so skip these.
		if len(resource.GetKind()) == 0 || len(resource.GetAPIVersion()) == 0 {
			continue
		}

		if len(defaultNamespace) > 0 {
			namespaced, err := clients.IsNamespacedKind(resource.GroupVersionKind(), clientSet)
			if err != nil {
				if clients.IsNoNamespaceInfoErr(err) {
					// Assume resource is namespaced.
					namespaced = true
				} else {
					return nil, err
				}
			}

			// Set namespace if resource Kind is namespaced and namespace is not already set.
			if namespaced && len(resource.GetNamespace()) == 0 {
				resource.SetNamespace(defaultNamespace)
			}
		}
		resources = append(resources, resource)
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, resource.Object)
	}

	return result, nil
}
