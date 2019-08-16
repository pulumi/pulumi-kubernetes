package await

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// AggregatedError represents an error with 0 or more sub-errors.
type AggregatedError interface {
	SubErrors() []string
}

// PartialError represents an object that failed to complete its current operation.
type PartialError interface {
	Object() *unstructured.Unstructured
}

// cancellationError represents an operation that failed because the user cancelled it.
type cancellationError struct {
	object    *unstructured.Unstructured
	subErrors []string
}

var _ error = (*cancellationError)(nil)
var _ AggregatedError = (*cancellationError)(nil)
var _ PartialError = (*cancellationError)(nil)

func (ce *cancellationError) Error() string {
	return fmt.Sprintf("Resource operation was cancelled for '%s'", ce.object.GetName())
}

// SubErrors returns the errors that were present when cancellation occurred.
func (ce *cancellationError) SubErrors() []string {
	return ce.subErrors
}

func (ce *cancellationError) Object() *unstructured.Unstructured {
	return ce.object
}

// timeoutError represents an operation that failed because it timed out.
type timeoutError struct {
	object    *unstructured.Unstructured
	subErrors []string
}

var _ error = (*timeoutError)(nil)
var _ AggregatedError = (*timeoutError)(nil)
var _ PartialError = (*timeoutError)(nil)

func (te *timeoutError) Error() string {
	// TODO(levi): May want to add a shortlink to more detailed troubleshooting docs.
	return fmt.Sprintf("'%s' timed out waiting to be Ready", te.object.GetName())
}

// SubErrors returns the errors that were present when timeout occurred.
func (te *timeoutError) SubErrors() []string {
	return te.subErrors
}

func (te *timeoutError) Object() *unstructured.Unstructured {
	return te.object
}

// initializationError occurs when we attempt to read a resource that failed to fully initialize.
type initializationError struct {
	subErrors []string
	object    *unstructured.Unstructured
}

var _ error = (*initializationError)(nil)
var _ AggregatedError = (*initializationError)(nil)
var _ PartialError = (*initializationError)(nil)

func (ie *initializationError) Error() string {
	return fmt.Sprintf("Resource '%s' was created but failed to initialize", ie.object.GetName())
}

// SubErrors returns the errors that were present when timeout occurred.
func (ie *initializationError) SubErrors() []string {
	return ie.subErrors
}

func (ie *initializationError) Object() *unstructured.Unstructured {
	return ie.object
}
