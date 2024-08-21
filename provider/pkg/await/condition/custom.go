package condition

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

var _ Satisfier = (*Custom)(nil)

// Custom waits for a specific ".status.condition" matching a user-provided
// expression.
type Custom struct {
	observer        *ObjectObserver
	logger          logger
	conditionType   string
	conditionStatus string
}

// NewCustom creates a new Custom condition.
//
// The expression's syntax is identical to `kubectl wait --for=condition=`. A
// condition type is required, optionally followed by a value:
//
//	"condition=Foo" or "condition=Foo=Bar"
//
// The "condition=" prefix is also optional.
func NewCustom(
	ctx context.Context,
	source Source,
	logger logger,
	uns *unstructured.Unstructured,
	expr string,
) (*Custom, error) {
	expr = strings.TrimPrefix(expr, "condition=")

	condition := expr
	status := "True"

	if idx := strings.Index(condition, "="); idx != -1 {
		condition = expr[0:idx]
		status = expr[idx+1:]
	}

	cond := &Custom{
		observer:        NewObjectObserver(ctx, source, uns),
		logger:          logger,
		conditionType:   condition,
		conditionStatus: status,
	}
	return cond, nil
}

// Satisfied returns true when the object's last-known state matches the
// expected condition.
func (cc *Custom) Satisfied() (bool, error) {
	return checkCondition(cc.Object(), cc.logger, cc.conditionType, cc.conditionStatus)
}

// Observe is a passthrough to the underlying Observer.
func (cc *Custom) Observe(e watch.Event) error {
	return cc.observer.Observe(e)
}

// Range is a passthrough to the underlying Observer.
func (cc *Custom) Range(yield func(watch.Event) bool) {
	cc.observer.Range(yield)
}

// Object is a passthrough to the underlying Observer.
func (cc *Custom) Object() *unstructured.Unstructured {
	return cc.observer.Object()
}
