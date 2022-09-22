// Copyright 2021, Pulumi Corporation.  All rights reserved.

package await

import (
	"context"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/watcher"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

	_, isPartialErr := err.(PartialError)
	assert.True(t, isPartialErr, "Cancelled watcher should emit `await.PartialError`")
	assert.Equal(t, "Resource operation was cancelled for ''", err.Error())
}

func Test_Watcher_Interface_Timeout(t *testing.T) {
	// Timeout because the `WatchUntil` predicate always returns false.
	err := watcher.ForObject(context.Background(), &mockResourceInterface{}, "").
		WatchUntil(func(_ *unstructured.Unstructured) bool { return false }, 1*time.Second)

	_, isPartialErr := err.(PartialError)
	assert.True(t, isPartialErr, "Timed out watcher should emit `await.PartialError`")
	assert.Equal(t, "Timeout occurred polling for ''", err.Error())
}

// --------------------------------------------------------------------------

// Mock implementations of Kubernetes client stuff.

// --------------------------------------------------------------------------

type mockResourceInterface struct{}

var _ dynamic.ResourceInterface = (*mockResourceInterface)(nil)

func (mri *mockResourceInterface) Create(
	ctx context.Context, obj *unstructured.Unstructured, options metav1.CreateOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Create not implemented")
}
func (mri *mockResourceInterface) Update(
	ctx context.Context, obj *unstructured.Unstructured, options metav1.UpdateOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Update not implemented")
}
func (mri *mockResourceInterface) UpdateStatus(
	ctx context.Context, obj *unstructured.Unstructured, options metav1.UpdateOptions) (*unstructured.Unstructured, error) {
	panic("UpdateStatus not implemented")
}
func (mri *mockResourceInterface) Delete(ctx context.Context, name string, options metav1.DeleteOptions, subresources ...string) error {
	panic("Delete not implemented")
}
func (mri *mockResourceInterface) DeleteCollection(
	ctx context.Context, deleteOptions metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	panic("DeleteCollection not implemented")
}
func (mri *mockResourceInterface) Get(
	ctx context.Context, name string, options metav1.GetOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{Object: map[string]interface{}{}}, nil
}
func (mri *mockResourceInterface) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	panic("List not implemented")
}
func (mri *mockResourceInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("Watch not implemented")
}
func (mri *mockResourceInterface) Patch(
	ctx context.Context, name string, pt types.PatchType, data []byte, options metav1.PatchOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Patch not implemented")
}

func (mri *mockResourceInterface) Apply(ctx context.Context, name string, obj *unstructured.Unstructured, options metav1.ApplyOptions, subresources ...string,
) (*unstructured.Unstructured, error) {
	panic("Apply not implemented")
}

func (mri *mockResourceInterface) ApplyStatus(ctx context.Context, name string, obj *unstructured.Unstructured, options metav1.ApplyOptions,
) (*unstructured.Unstructured, error) {
	panic("ApplyStatus not implemented")
}
