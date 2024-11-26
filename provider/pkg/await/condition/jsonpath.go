package condition

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/jsonpath"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

var _ Satisfier = (*JSONPath)(nil)

// JSONPath waits for the observed object to match a user-provided JSONPath
// expression.
type JSONPath struct {
	observer *ObjectObserver
	logger   logger
	jsp      *jsonpath.Parsed
}

// NewJSONPath creates a new JSONPath condition.
func NewJSONPath(
	ctx context.Context,
	source Source,
	logger logger,
	obj *unstructured.Unstructured,
	jsp *jsonpath.Parsed,
) (Satisfier, error) {
	cond := &JSONPath{
		observer: NewObjectObserver(ctx, source, obj),
		jsp:      jsp,
		logger:   logger,
	}
	return cond, nil
}

// Satisfied returns true with the observed object's last-known state matches
// the provided JSONPath expression.
func (jp *JSONPath) Satisfied() (bool, error) {
	message := fmt.Sprintf("Waiting for %s", jp.jsp)

	// Ensure the .status.observedGeneration is current, if present.
	generation, found, _ := unstructured.NestedInt64(jp.Object().Object, "metadata", "generation")
	if found {
		observedGeneration, ok, _ := unstructured.NestedInt64(jp.Object().Object, "status", "observedGeneration")
		if ok && observedGeneration < generation {
			return false, nil
		}
	}

	result, err := jp.jsp.Matches(jp.Object())
	if result.Found != "" {
		message = fmt.Sprintf("%s (found %q)", message, result.Found)
	}
	if result.Message != "" {
		message = result.Message
	}
	if result.Matched {
		jp.logger.LogStatus(diag.Info, "Found "+jp.jsp.String())
	} else {
		jp.logger.LogStatus(diag.Info, message)
	}
	return result.Matched, err
}

// Observe is a passthrough to the underlying Observer.
func (jp *JSONPath) Observe(e watch.Event) error {
	return jp.observer.Observe(e)
}

// Range is a passthrough to the underlying Observer.
func (jp *JSONPath) Range(yield func(watch.Event) bool) {
	jp.observer.Range(yield)
}

// Object is a passthrough to the underlying Observer.
func (jp *JSONPath) Object() *unstructured.Unstructured {
	return jp.observer.Object()
}
