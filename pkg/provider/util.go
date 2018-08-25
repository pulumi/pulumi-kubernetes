package provider

import (
	"github.com/pulumi/pulumi/pkg/resource"
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
