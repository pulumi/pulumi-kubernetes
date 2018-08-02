package await

import "fmt"

// AggregatedError represents an error with 0 or more sub-errors.
type AggregatedError interface {
	SubErrors() []string
}

// cancellationError represents an operation that failed because the user cancelled it.
type cancellationError struct {
	objectName string
	subErrors  []string
}

var _ error = (*cancellationError)(nil)
var _ AggregatedError = (*cancellationError)(nil)

func (ce *cancellationError) Error() string {
	return fmt.Sprintf("Resource operation was cancelled for '%s'", ce.objectName)
}

// SubErrors returns the errors that were present when cancellation occurred.
func (ce *cancellationError) SubErrors() []string {
	return ce.subErrors
}

// timeoutError represents an operation that failed because it timed out.
type timeoutError struct {
	objectName string
	subErrors  []string
}

var _ error = (*timeoutError)(nil)
var _ AggregatedError = (*timeoutError)(nil)

func (te *timeoutError) Error() string {
	return fmt.Sprintf("Timeout occurred for '%s'", te.objectName)
}

// SubErrors returns the errors that were present when timeout occurred.
func (te *timeoutError) SubErrors() []string {
	return te.subErrors
}
