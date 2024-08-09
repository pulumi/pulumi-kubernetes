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
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pulumi/cloud-ready-checks/pkg/checker"
	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"

	"github.com/pulumi/cloud-ready-checks/pkg/kubernetes/pod"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

var _ condition.Observer = (*Aggregator[*corev1.Event])(nil)

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
	messages chan checkerlog.Messages
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
		messages: make(chan checkerlog.Messages),
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
		if isOwnedBy(obj, pa.owner) {
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
func (pa *PodAggregator) Read() checkerlog.Messages {
	var messages checkerlog.Messages
	checkPod := func(object runtime.Object) {
		obj := object.(*unstructured.Unstructured)
		pod, err := clients.PodFromUnstructured(obj)
		if err != nil {
			logger.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return
		}
		if isOwnedBy(obj, pa.owner) {
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
func (pa *PodAggregator) ResultChan() <-chan checkerlog.Messages {
	return pa.messages
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

// Aggregator is a generic, stateless condition.Observer intended for reporting
// informational messages about related resources during an Await.
type Aggregator[T runtime.Object] struct {
	observer condition.Observer
	callback func(logMessager, T) error
	logger   logMessager
}

// NewAggregator creates a new Aggregator for the given runtime type. The
// provided condition.Observer must be configured for the corresponding GVK.
func NewAggregator[T runtime.Object](
	observer condition.Observer,
	logger logMessager,
	callback func(logMessager, T) error,
) *Aggregator[T] {
	return &Aggregator[T]{
		observer: observer,
		callback: callback,
		logger:   logger,
	}
}

func (i *Aggregator[T]) Observe(e watch.Event) error {
	obj, ok := e.Object.(*unstructured.Unstructured)
	if !ok {
		return nil
	}
	var t T
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &t)
	if err != nil {
		return err
	}
	return i.callback(i.logger, t)
}

func (i *Aggregator[T]) Range(yield func(watch.Event) bool) {
	i.observer.Range(yield)
}

// NewEventAggregator creates a new condition.Observer subscribed to Kubernetes
// events related to the owner object. Event messages are logged at WARN
// severity if the event has type Warning; Normal events are discarded to
// reduce noise.
func NewEventAggregator(
	ctx context.Context,
	source condition.Source,
	logger logMessager,
	owner *unstructured.Unstructured,
) condition.Observer {
	observer := condition.NewObserver(ctx,
		source,
		corev1.SchemeGroupVersion.WithKind("Event"),
		relatedEvents(owner),
	)
	return NewAggregator(observer, logger,
		func(l logMessager, e *corev1.Event) error {
			if e == nil {
				return nil
			}
			msg := fmt.Sprintf(
				"[%s/%s] %s: %s",
				strings.ToLower(e.InvolvedObject.Kind),
				e.InvolvedObject.Name,
				e.Reason,
				e.Message,
			)
			m := checkerlog.WarningMessage(msg)
			switch e.Type {
			case corev1.EventTypeWarning:
				logger.LogStatus(diag.Warning, m.S)
			default:
				logger.LogStatus(diag.Debug, m.S)
			}
			return nil
		},
	)
}

type logMessager interface {
	Log(diag.Severity, string)
	LogStatus(diag.Severity, string)
}

func relatedEvents(owner *unstructured.Unstructured) func(*unstructured.Unstructured) bool {
	return func(obj *unstructured.Unstructured) bool {
		uid, _, _ := unstructured.NestedString(obj.Object, "involvedObject", "uid")
		return uid == string(owner.GetUID())
	}
}
