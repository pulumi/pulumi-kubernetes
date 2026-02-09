package condition

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

var _ Satisfier = (*All)(nil)

// NewAll joins multiple Satisfiers and resolves when all of them are
// simultaneously satisfied. The conditions should all apply to the same object.
func NewAll(conditions ...Satisfier) (*All, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("requires a condition")
	}
	obj := conditions[0].Object()
	if obj == nil {
		return nil, fmt.Errorf("requires an object")
	}
	gvk := obj.GroupVersionKind()
	for _, c := range conditions {
		if c.Object().GroupVersionKind() != gvk {
			return nil, fmt.Errorf("GVK mismatch: %q != %q", c.Object().GroupVersionKind(), gvk)
		}
	}
	cond := &All{
		conditions: conditions,
	}
	return cond, nil
}

type All struct {
	conditions []Satisfier
}

// Satisfied returns true when all the sub-conditions are true.
func (ac *All) Satisfied() (bool, error) {
	for _, c := range ac.conditions {
		done, err := c.Satisfied()
		if !done || err != nil {
			return false, err
		}
	}
	return true, nil
}

// Observe sends the given event to all sub-conditions.
func (ac *All) Observe(e watch.Event) error {
	var err error
	for _, c := range ac.conditions {
		err = errors.Join(err, c.Observe(e))
	}
	return err
}

// Object returns the first condition's current state.
func (ac *All) Object() *unstructured.Unstructured {
	return ac.conditions[0].Object()
}

// Range iterates over the first condition's events since all conditions are
// assumed to be watching the same object.
func (ac *All) Range(yield func(watch.Event) bool) {
	ac.conditions[0].Range(yield)
}
