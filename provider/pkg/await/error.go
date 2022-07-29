// Copyright 2021, Pulumi Corporation.  All rights reserved.

package await

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
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

// PreviewError represents a preview operation that failed.
type PreviewError interface {
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
	return fmt.Sprintf("Resource operation was cancelled for %q", ce.object.GetName())
}

// SubErrors returns the errors that were present when cancellation occurred.
func (ce *cancellationError) SubErrors() []string {
	return ce.subErrors
}

func (ce *cancellationError) Object() *unstructured.Unstructured {
	return ce.object
}

// namespaceError represents an operation that failed because the namespace didn't exist.
type namespaceError struct {
	object *unstructured.Unstructured
}

var _ error = (*namespaceError)(nil)
var _ PreviewError = (*namespaceError)(nil)

func (ne *namespaceError) Error() string {
	return fmt.Sprintf("namespace does not exist for %q", ne.object.GetName())
}

func (ne *namespaceError) Object() *unstructured.Unstructured {
	return ne.object
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

// IsNamespaceNotFoundErr returns true if the namespace wasn't found for a k8s client operation.
func IsNamespaceNotFoundErr(err error) bool {
	se, isStatusError := err.(*errors.StatusError)
	if !isStatusError {
		return false
	}

	return errors.IsNotFound(err) && se.Status().Details.Kind == "namespaces"
}

// IsResourceExistsErr returns true if the resource already exists on the k8s cluster.
func IsResourceExistsErr(err error) bool {
	_, isStatusError := err.(*errors.StatusError)
	if !isStatusError {
		return false
	}

	return errors.IsAlreadyExists(err)
}

// IsDeleteRequiredFieldErr is true if the user attempted to delete a Patch resource that was the sole manager of a
// required field.
func IsDeleteRequiredFieldErr(err error) bool {
	return errors.IsInvalid(err) && strings.Contains(err.Error(), "Required value")
}
