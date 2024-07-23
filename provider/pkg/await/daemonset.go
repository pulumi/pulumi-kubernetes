// Copyright 2024, Pulumi Corporation.
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
	"reflect"
	"time"

	"github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

const (
	_defaultDaemonSetTimeout = 10 * time.Minute
)

// dsAwaiter manages await logic for extensions/v1beta1/DaemonSet,
// apps/v1beta1/DaemonSet, apps/v1beta2/DaemonSet, and apps/v1/DaemonSet. It
// handles create, update, read, and delete operations.
//
// A DaemonSet is a construct that allows users to run at most one pod per
// node. DaemonSets operate in two different update modes, depending on the
// specified .spec.updateStrategy.type:
//
//  1. RollingUpdate (default) - After the DaemonSet is updated, old pods will
//     be killed and new ones will be created in a controlled fashion. At most
//     one pod will be running on each node during the whole update process
//     (unless .spec.updateStrategy.maxSurge is specified).
//
//  2. OnDelete - After the DaemonSet is updated, new pods will only be
//     created when the user manually deletes old DaemonSet pods.
//
// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/daemon-set-v1/#DaemonSetStatus
//
// The success conditions are the same regardless of the update strategy and
// are determined by
// https://pkg.go.dev/sigs.k8s.io/cli-utils/pkg/kstatus/status.

// Importantly, this means OnDelete rollouts will wait until pods have been
// manually cleaned up unless the skipAwait annotation is present.
type dsAwaiter struct {
	config  updateAwaitConfig
	ds      *unstructured.Unstructured
	deleted bool
}

// newDaemonSetAwaiter returns a new dsAwaiter.
func newDaemonSetAwaiter(c updateAwaitConfig) *dsAwaiter {
	return &dsAwaiter{
		config: c,
		ds:     c.currentOutputs,
	}
}

// Await blocks until a DaemonSet is ready or encounters an error.
func (dsa *dsAwaiter) Await() error {
	return dsa.await(dsa.rolloutComplete)
}

// Delete blocks until a DaemonSet has been deleted. Returns nil if the
// DaemonSet does not exist.
func (dsa *dsAwaiter) Delete() error {
	// Perform a lookup in case the object has already been deleted
	client, _, err := dsa.clients()
	if err != nil {
		return err
	}
	ds := dsa.config.currentOutputs
	_, err = client.Get(dsa.config.ctx, ds.GetName(), metav1.GetOptions{})
	if is404(err) {
		return nil
	}
	if err != nil {
		logger.V(3).Infof("Received error deleting DaemonSet %q: %#v", ds.GetName(), err)
		return err
	}

	// Otherwise wait for a deletion event.
	deleted := func() bool {
		if dsa.deleted {
			return true
		}
		misscheduled, _ := openapi.Pluck(dsa.ds.Object, "status", "numberMisscheduled")
		dsa.config.logStatus(
			diag.Info,
			fmt.Sprintf(
				"DaemonSet %q still exists (%d pods misscheduled)",
				ds.GetName(),
				misscheduled,
			),
		)
		return false
	}

	return dsa.await(deleted)
}

// Read returns the current state of the DaemonSet and returns an error if it
// is not in a ready state.
func (dsa *dsAwaiter) Read() error {
	dsClient, podClient, err := dsa.clients()
	if err != nil {
		return err
	}

	// Get live versions of the DaemonSet.
	ds, err := dsClient.Get(
		dsa.config.ctx,
		dsa.config.currentOutputs.GetName(),
		metav1.GetOptions{},
	)
	if err != nil {
		return err
	}
	// Update the awaiter's state to reflect the live object.
	dsa.processDaemonSetEvent(watchAddedEvent(ds))

	// Grab any error messages from pods for more helpful debugging.
	pods, err := podClient.List(dsa.config.ctx, metav1.ListOptions{})
	if err != nil {
		logger.V(3).Infof(
			"Error retrieving Pod list for DaemonSet %q: %v",
			ds.GetName(),
			err,
		)
		pods = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}
	pa := NewPodAggregator(ds, &staticLister{pods})
	messages := pa.Read()
	dsa.processPodMessages(messages)

	suberrors := []string{}
	for _, e := range messages.Errors() {
		suberrors = append(suberrors, e.String())
	}

	if dsa.rolloutComplete() {
		return nil
	}

	return &initializationError{
		object:    ds,
		subErrors: suberrors,
	}
}

// await watches the DaemonSet and its Pods until the `done` condition is met
// or until a timeout (or other error) occurs.
func (dsa *dsAwaiter) await(done func() bool) error {
	timeout := _defaultDaemonSetTimeout
	if dsa.config.timeout != nil {
		timeout = *dsa.config.timeout
	}
	ctx, cancel := context.WithCancelCause(dsa.config.ctx)
	defer cancel(context.Canceled)
	go func() {
		dsa.config.Clock().Sleep(timeout)
		cancel(context.DeadlineExceeded)
	}()

	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory := informers.NewInformerFactory(
		dsa.config.clientSet,
		informers.WithNamespaceOrDefault(dsa.config.currentOutputs.GetNamespace()),
	)
	informerFactory.Start(stopper)

	dsEvents := make(chan watch.Event)
	dsInformer, err := informers.New(
		informerFactory,
		informers.ForGVR(v1.SchemeGroupVersion.WithResource("daemonsets")),
		informers.WithEventChannel(dsEvents),
	)
	if err != nil {
		return err
	}
	go dsInformer.Informer().Run(stopper)

	podEvents := make(chan watch.Event)
	podInformer, err := informers.New(
		informerFactory,
		informers.ForPods(),
		informers.WithEventChannel(podEvents),
	)
	if err != nil {
		return err
	}
	go podInformer.Informer().Run(stopper)

	podAggregator := NewPodAggregator(dsa.ds, podInformer.Lister())
	podAggregator.Start(podEvents)
	defer podAggregator.Stop()

	for {
		if done() {
			return nil
		}
		select {
		case <-ctx.Done():
			return wait.ErrorInterrupted(nil)
		case event := <-dsEvents:
			dsa.processDaemonSetEvent(event)
		case messages := <-podAggregator.ResultChan():
			dsa.processPodMessages(messages)
		}
	}
}

// rolloutComplete checks whether we've succeeded, and logs the result as a
// status message to the provider.
func (dsa *dsAwaiter) rolloutComplete() bool {
	res, err := status.Compute(dsa.ds)
	if err != nil {
		dsa.config.logStatus(diag.Error, err.Error())
		return false
	}

	done := res.Status == status.CurrentStatus

	if done {
		dsa.config.logStatus(diag.Info, fmt.Sprintf("%s%s", cmdutil.EmojiOr("✅ ", ""), res.Message))
		return true
	}

	dsa.config.logStatus(diag.Info, res.Message)
	return false
}

// processDaemonSetEvent updates dsAwaiter's state to reflect the DS watch event.
func (dsa *dsAwaiter) processDaemonSetEvent(event watch.Event) {
	name := dsa.config.currentOutputs.GetName()

	ds, ok := event.Object.(*unstructured.Unstructured)
	if !ok {
		logger.V(3).Infof(
			"DaemonSet watch received unknown object type %q",
			reflect.TypeOf(ds),
		)
		return
	}

	// Do nothing if this is not the DaemonSet we're waiting for.
	if ds.GetName() != name {
		return
	}

	if event.Type == watch.Deleted {
		dsa.deleted = true
		return
	}

	dsa.ds = ds
}

// processPodMessages logs pod messages from a PodAggregator.
func (dsa *dsAwaiter) processPodMessages(messages logging.Messages) {
	for _, message := range messages {
		dsa.config.logMessage(message)
	}
}

// clients returns clients for the DaemonSet and its Pods.
func (dsa *dsAwaiter) clients() (
	dsClient, podClient dynamic.ResourceInterface, err error,
) {
	dsClient, err = clients.ResourceClient(
		kinds.DaemonSet, dsa.config.currentOutputs.GetNamespace(), dsa.config.clientSet)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not make client to watch DaemonSet %q: %w",
			dsa.config.currentOutputs.GetName(), err)
	}
	podClient, err = clients.ResourceClient(
		kinds.Pod, dsa.config.currentOutputs.GetNamespace(), dsa.config.clientSet)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not make client to watch Pods associated with DaemonSet %q: %w",
			dsa.config.currentOutputs.GetName(), err)
	}

	return dsClient, podClient, nil
}
