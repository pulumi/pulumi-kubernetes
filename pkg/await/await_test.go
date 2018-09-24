package await

import (
	"context"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/pkg/watcher"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

func Test_Watcher_Interface_Cancel(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// Cancel should occur before `WatchUntil` because predicate always returns false.
	err := watcher.ForObject(cancelCtx, &mockResourceInterface{}, "").
		WatchUntil(func(_ *unstructured.Unstructured) bool { return false }, 1*time.Minute)

	_, isInitErr := err.(InitializationError)
	assert.True(t, isInitErr, "Cancelled watcher should emit `await.InitializationError`")
	assert.Equal(t, "Resource operation was cancelled for ''", err.Error())
}

func Test_Watcher_Interface_Timeout(t *testing.T) {
	// Timeout because the `WatchUntil` predicate always returns false.
	err := watcher.ForObject(context.Background(), &mockResourceInterface{}, "").
		WatchUntil(func(_ *unstructured.Unstructured) bool { return false }, 1*time.Second)

	_, isInitErr := err.(InitializationError)
	assert.True(t, isInitErr, "Timed out watcher should emit `await.InitializationError`")
	assert.Equal(t, "Timeout occurred polling for ''", err.Error())
}

// --------------------------------------------------------------------------

// Mock implementations of Kubernetes client stuff.

// --------------------------------------------------------------------------

type mockResourceInterface struct{}

var _ dynamic.ResourceInterface = (*mockResourceInterface)(nil)

func (mri *mockResourceInterface) List(opts metav1.ListOptions) (runtime.Object, error) {
	panic("List not implemented")
}
func (mri *mockResourceInterface) Get(
	name string, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{Object: map[string]interface{}{}}, nil
}
func (mri *mockResourceInterface) Delete(name string, opts *metav1.DeleteOptions) error {
	panic("Delete not implemented")
}
func (mri *mockResourceInterface) DeleteCollection(
	deleteOptions *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	panic("DeleteCollection not implemented")
}
func (mri *mockResourceInterface) Create(
	obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("Create not implemented")
}
func (mri *mockResourceInterface) Update(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	panic("Update not implemented")
}
func (mri *mockResourceInterface) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	panic("Watch not implemented")
}
func (mri *mockResourceInterface) Patch(
	name string, pt types.PatchType, data []byte) (*unstructured.Unstructured, error) {
	panic("Patch not implemented")
}
