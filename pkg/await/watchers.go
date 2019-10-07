// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package await

import (
	"sync"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi-kubernetes/pkg/await/states"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/pkg/logging"
	"github.com/pulumi/pulumi/pkg/diag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// PodAggregator tracks status for any Pods related to the owner resource, and writes
// warning/error messages to a channel that can be consumed by a resource awaiter.
type PodAggregator struct {
	// Synchronization
	sync.Mutex
	stopped bool

	// Owner identity
	owner ResourceId

	// Pod checker
	checker *states.StateChecker

	// Clients
	client  dynamic.ResourceInterface
	watcher watch.Interface

	// Messages
	messages chan logging.Messages
}

// NewPodAggregator returns an initialized PodAggregator.
func NewPodAggregator(owner ResourceId, clientset *clients.DynamicClientSet) (*PodAggregator, error) {
	client, err := clients.ResourceClient(kinds.Pod, owner.Namespace, clientset)
	if err != nil {
		return nil, err
	}

	watcher, err := client.Watch(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pa := &PodAggregator{
		stopped:  false,
		owner:    owner,
		checker:  states.NewPodChecker(),
		client:   client,
		watcher:  watcher,
		messages: make(chan logging.Messages),
	}
	go pa.run()

	return pa, nil
}

// run contains the aggregation logic and is executed as a goroutine when a PodAggregator
// is initialized.
func (pa *PodAggregator) run() {
	defer close(pa.messages)

	checkPod := func(object runtime.Object) {
		pod, err := clients.PodFromUnstructured(object.(*unstructured.Unstructured))
		if err != nil {
			glog.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return
		}
		if relatedResource(pa.owner, pod) {
			messages := pa.checker.Update(pod).MessagesWithSeverity(diag.Warning, diag.Error)
			if len(messages) > 0 {
				pa.messages <- messages
			}
		}
	}

	// Get existing Pods.
	pods, err := pa.client.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list existing Pods: %v", err)
	} else {
		// Log errors and move on.
		_ = pods.EachListItem(func(object runtime.Object) error {
			checkPod(object)
			return nil
		})
	}

	for {
		if pa.stopping() {
			return
		}
		event := <-pa.watcher.ResultChan()
		if event.Object == nil {
			continue
		}
		checkPod(event.Object)
	}
}

// Read lists existing Pods and returns any related warning/error messages.
func (pa *PodAggregator) Read() logging.Messages {
	var messages logging.Messages
	checkPod := func(object runtime.Object) {
		pod, err := clients.PodFromUnstructured(object.(*unstructured.Unstructured))
		if err != nil {
			glog.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return
		}
		if relatedResource(pa.owner, pod) {
			messages = append(messages, pa.checker.Update(pod).MessagesWithSeverity(diag.Warning, diag.Error)...)
		}
	}

	// Get existing Pods.
	pods, err := pa.client.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Failed to list existing Pods: %v", err)
	} else {
		// Log errors and move on.
		_ = pods.EachListItem(func(object runtime.Object) error {
			checkPod(object)
			return nil
		})
	}

	return messages
}

// Stop safely stops a PodAggregator and underlying watch client.
func (pa *PodAggregator) Stop() {
	pa.Lock()
	defer pa.Unlock()
	if !pa.stopped {
		pa.stopped = true
		pa.watcher.Stop()
	}
}

// stopping returns true if Stop() was called previously.
func (pa *PodAggregator) stopping() bool {
	pa.Lock()
	defer pa.Unlock()
	return pa.stopped
}

// ResultChan returns a reference to the message channel used by the PodAggregator to
// communicate warning/error messages to a resource awaiter.
func (pa *PodAggregator) ResultChan() <-chan logging.Messages {
	return pa.messages
}

// relatedResource returns true if a resource ownerReference and metadata matches the provided owner
// ResourceId, false otherwise.
//
// Example ownerReference:
// {
//     "apiVersion": "batch/v1",
//     "blockOwnerDeletion": true,
//     "controller": true,
//     "kind": "Job",
//     "name": "foo",
//     "uid": "14ba58cc-cf83-11e9-8c3a-025000000001"
// }
func relatedResource(owner ResourceId, object metav1.Object) bool {
	if owner.Namespace != object.GetNamespace() {
		return false
	}
	if owner.Generation != object.GetGeneration() {
		return false
	}
	for _, ref := range object.GetOwnerReferences() {
		if ref.APIVersion != owner.GVK.GroupVersion().String() {
			continue
		}
		if ref.Kind != owner.GVK.Kind {
			continue
		}
		if ref.Name != owner.Name {
			continue
		}
		return true
	}
	return false
}
