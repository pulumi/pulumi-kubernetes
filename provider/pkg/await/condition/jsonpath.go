package condition

import (
	"context"
	"fmt"

	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/jsonpath"
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

	result, err := jp.jsp.Matches(jp.Object())
	if result.Found != "" {
		message = fmt.Sprintf("%s (found %q)", message, result.Found)
	}
	if result.Message != "" {
		message = result.Message
	}
	if result.Matched {
		jp.logger.LogMessage(checkerlog.StatusMessage("Found " + jp.jsp.String()))
	} else {
		jp.logger.LogMessage(checkerlog.StatusMessage(message))
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
