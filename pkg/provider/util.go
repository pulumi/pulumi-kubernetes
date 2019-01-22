package provider

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/pkg/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func hasComputedValue(obj *unstructured.Unstructured) bool {
	if obj == nil || obj.Object == nil {
		return false
	}

	objects := []map[string]interface{}{obj.Object}
	var curr map[string]interface{}

	for {
		if len(objects) == 0 {
			break
		}
		curr, objects = objects[0], objects[1:]
		for _, v := range curr {
			if _, isComputed := v.(resource.Computed); isComputed {
				return true
			}
			if field, isMap := v.(map[string]interface{}); isMap {
				objects = append(objects, field)
			}
		}
	}

	return false
}

// --------------------------------------------------------------------------
// Names and namespaces.
// --------------------------------------------------------------------------

// FqObjName returns "namespace.name"
func FqObjName(o metav1.Object) string {
	return FqName(o.GetNamespace(), o.GetName())
}

// ParseFqName will parse a fully-qualified Kubernetes object name.
func ParseFqName(id string) (namespace, name string) {
	split := strings.Split(id, "/")
	if len(split) == 1 {
		return "", split[0]
	}
	namespace, name = split[0], split[1]
	return
}

// FqName returns "namespace/name"
func FqName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s/%s", namespace, name)
}
