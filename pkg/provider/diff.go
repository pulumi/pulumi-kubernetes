package provider

import (
	"errors"
	"fmt"

	st "github.com/golang/protobuf/ptypes/struct"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/yudai/gojsondiff"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func computedProperties(obj, liveObj *unstructured.Unstructured) (*st.Struct, error) {
	diff := gojsondiff.New().CompareObjects(obj.Object, liveObj.Object)
	computed := resource.PropertyMap{}
	err := changedFields("", computed, diff.Deltas())
	if err != nil {
		return nil, err
	}

	return plugin.MarshalProperties(computed, plugin.MarshalOptions{
		KeepUnknowns: true, SkipNulls: true,
	})
}

// changedFields produces a list of fields have have changed. They are specified as JSON paths from
// the root of the API object, e.g., `.metadata.name` or `.spec.containers`.
//
// TODO(hausdorff): Come back later and wire this back up to the `kubeProvider#Diff` when v0.12.0 is
// cut, and the provider needs to provide the diff.
func changedFields(path string, values resource.PropertyMap, deltas []gojsondiff.Delta) error {
	appendPath := func(component string) string {
		return path + "." + component
	}

	for _, delta := range deltas {
		switch delta.(type) {
		case *gojsondiff.Object:
			d := delta.(*gojsondiff.Object)
			err := changedFields(appendPath(d.String()), values, d.Deltas)
			if err != nil {
				return err
			}
		case *gojsondiff.Array:
			d := delta.(*gojsondiff.Array)
			err := changedFields(appendPath(d.String()), values, d.Deltas)
			if err != nil {
				return err
			}
		case *gojsondiff.Added:
			d := delta.(*gojsondiff.Added)
			values[resource.PropertyKey(appendPath(d.String()))] = resource.NewPropertyValue(d.Value)
		case *gojsondiff.Modified:
			d := delta.(*gojsondiff.Modified)
			values[resource.PropertyKey(appendPath(d.String()))] = resource.NewPropertyValue(d.NewValue)
		case *gojsondiff.TextDiff:
			d := delta.(*gojsondiff.TextDiff)
			values[resource.PropertyKey(appendPath(d.String()))] = resource.NewPropertyValue(d.NewValue)
		case *gojsondiff.Deleted:
			break
		case *gojsondiff.Moved:
			return errors.New("Delta type 'Move' is not supported in objects")
		default:
			return fmt.Errorf("Unknown Delta type detected: %#v", delta)
		}
	}

	return nil
}
