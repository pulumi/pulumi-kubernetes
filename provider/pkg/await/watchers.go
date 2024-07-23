// Copyright 2016-2022, Pulumi Corporation.
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

	"github.com/pulumi/cloud-ready-checks/pkg/checker"
	"github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/cloud-ready-checks/pkg/kubernetes/pod"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

// PodAggregator tracks status for any Pods related to the owner resource, and writes
// warning/error messages to a channel that can be consumed by a resource awaiter.
type PodAggregator struct {
	// Synchronization
	sync.Mutex
	stopped bool

	// Owner identity
	owner *unstructured.Unstructured

	// Pod checker
	checker *checker.StateChecker

	// Clients
	lister lister

	// Messages
	messages chan logging.Messages
}

// lister lists resources matching a label selector.
type lister interface {
	List(selector labels.Selector) (ret []runtime.Object, err error)
}

// NewPodAggregator returns an initialized PodAggregator.
func NewPodAggregator(owner *unstructured.Unstructured, lister lister) *PodAggregator {
	pa := &PodAggregator{
		owner:    owner,
		lister:   lister,
		checker:  pod.NewPodChecker(),
		messages: make(chan logging.Messages),
	}
	return pa
}

// Start initiates the aggregation logic and is executed as a goroutine which should be
// stopped through a call to Stop
func (pa *PodAggregator) Start(informChan <-chan watch.Event) {
	go pa.run(informChan)
}

func (pa *PodAggregator) run(informChan <-chan watch.Event) {
	defer close(pa.messages)

	checkPod := func(object runtime.Object) {
		obj := object.(*unstructured.Unstructured)
		pod, err := clients.PodFromUnstructured(obj)
		if err != nil {
			logger.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return
		}
		if relatedResource(pa.owner, obj) {
			_, results := pa.checker.ReadyDetails(pod)
			messages := results.Messages().MessagesWithSeverity(diag.Warning, diag.Error)
			if len(messages) > 0 {
				pa.messages <- messages
			}
		}
	}

	for {
		if pa.stopping() {
			return
		}
		event := <-informChan
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
		obj := object.(*unstructured.Unstructured)
		pod, err := clients.PodFromUnstructured(obj)
		if err != nil {
			logger.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return
		}
		if relatedResource(pa.owner, obj) {
			_, results := pa.checker.ReadyDetails(pod)
			messages = results.Messages().MessagesWithSeverity(diag.Warning, diag.Error)
		}
	}

	// Get existing Pods.
	pods, err := pa.lister.List(labels.Everything())
	if err != nil {
		logger.V(3).Infof("Failed to list existing Pods: %v", err)
	} else {
		for _, pod := range pods {
			checkPod(pod)
		}
	}

	return messages
}

// Stop safely stops a PodAggregator and underlying watch client.
func (pa *PodAggregator) Stop() {
	pa.Lock()
	defer pa.Unlock()
	if !pa.stopped {
		pa.stopped = true
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
// ResourceID, false otherwise.
//
// Example ownerReference:
//
//	{
//	    "apiVersion": "batch/v1",
//	    "blockOwnerDeletion": true,
//	    "controller": true,
//	    "kind": "Job",
//	    "name": "foo",
//	    "uid": "14ba58cc-cf83-11e9-8c3a-025000000001"
//	}
func relatedResource(owner *unstructured.Unstructured, obj *unstructured.Unstructured) bool {
	if owner.GetNamespace() != obj.GetNamespace() {
		return false
	}
	return isOwnedBy(obj, owner)
}

// staticLister can be used as a lister in situations where results are already
// known and only need to be aggregated.
type staticLister struct {
	list *unstructured.UnstructuredList
}

// List returns the staticLister's static contents.
func (s *staticLister) List(_ labels.Selector) (ret []runtime.Object, err error) {
	objects := []runtime.Object{}
	for _, l := range s.list.Items {
		objects = append(objects, l.DeepCopyObject())
	}
	return objects, nil
}
