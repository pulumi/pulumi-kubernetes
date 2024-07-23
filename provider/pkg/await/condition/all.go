package condition

import (
	"errors"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

var _ Satisfier = (*All)(nil)

// NewAll allows joining multiple Satisfiers of different GVKs.
// The first condition is the only one considered for reporting resource state in errors.
func NewAll(conditions ...Satisfier) (*All, error) {
	gvks := map[schema.GroupVersionKind][]Satisfier{}
	for _, c := range conditions {
		gvk := c.Object().GroupVersionKind()
		gvks[gvk] = append(gvks[gvk], c)
	}
	cond := &All{
		conditions: conditions,
		gvks:       gvks,
	}
	return cond, nil
}

type All struct {
	// mu         sync.Mutex
	conditions []Satisfier
	gvks       map[schema.GroupVersionKind][]Satisfier
}

// Satisfied returns true when all of the sub-conditions are true.
func (ac *All) Satisfied() (bool, error) {
	for _, c := range ac.conditions {
		done, err := c.Satisfied()
		if !done || err != nil {
			return false, err
		}
	}
	return true, nil
}

func (a *All) Observe(e watch.Event) error {
	// a.mu.Lock()
	// defer a.mu.Unlock()

	gvk := e.Object.GetObjectKind().GroupVersionKind()
	var err error
	for _, c := range a.gvks[gvk] {
		err = errors.Join(err, c.Observe(e))
	}
	return err
}

// Object returns the first condition's current state.
func (a *All) Object() *unstructured.Unstructured {
	// Not sure about this
	// panic("WHY DO I NEED THIS")
	return a.conditions[0].Object()
}

func (a *All) Range(yield func(watch.Event) bool) {
	wg := sync.WaitGroup{}
	// for gvks, conditions := range a.gvks[gvk] {
	// 	err = errors.Join(err, c.Observe(e))
	// }
	for _, c := range a.conditions {
		c := c
		// yy := func(e watch.Event) bool {
		// 	x := yield(e)
		// 	d, _ := a.Satisfied()
		// 	if d {
		// 		return false
		// 	}

		// 	return x
		// 	// return yield(e)
		// }
		wg.Add(1)
		go func(c Satisfier) {
			defer wg.Done()
			c.Range(yield)
			// c.Range(yy)
		}(c)
	}
	wg.Wait()

	// for _, c := range a.conditions {
	// 	c.Range(yy)
	// }
}
