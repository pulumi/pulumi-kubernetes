package provider

import (
	"errors"
	"fmt"

	"github.com/yudai/gojsondiff"
)

// replacedFields produces a list of fields have have changed. They are specified as JSON paths from
// the root of the API object, e.g., `.metadata.name` or `.spec.containers`.
//
// TODO(hausdorff): Come back later and wire this back up to the `kubeProvider#Diff` when v0.12.0 is
// cut, and the provider needs to provide the diff.
func replacedFields(path string, deltas []gojsondiff.Delta) ([]string, error) {
	appendPath := func(component string) string {
		return path + "." + component
	}

	replaced := []string{}
	for _, delta := range deltas {
		switch delta.(type) {
		case *gojsondiff.Object:
			d := delta.(*gojsondiff.Object)
			fields, err := replacedFields(appendPath(d.String()), d.Deltas)
			if err != nil {
				return nil, err
			}
			replaced = append(replaced, fields...)
		case *gojsondiff.Array:
			d := delta.(*gojsondiff.Array)
			fields, err := replacedFields(appendPath(d.String()), d.Deltas)
			if err != nil {
				return nil, err
			}
			replaced = append(replaced, fields...)
		case *gojsondiff.Added:
			d := delta.(*gojsondiff.Added)
			replaced = append(replaced, appendPath(d.String()))
		case *gojsondiff.Modified:
			d := delta.(*gojsondiff.Modified)
			replaced = append(replaced, appendPath(d.String()))
		case *gojsondiff.TextDiff:
			d := delta.(*gojsondiff.TextDiff)
			replaced = append(replaced, appendPath(d.String()))
		case *gojsondiff.Deleted:
			d := delta.(*gojsondiff.Deleted)
			replaced = append(replaced, appendPath(d.String()))
		case *gojsondiff.Moved:
			return nil, errors.New("Delta type 'Move' is not supported in objects")
		default:
			return nil, fmt.Errorf("Unknown Delta type detected: %#v", delta)
		}
	}

	return replaced, nil
}
