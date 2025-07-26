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
	"fmt"
	"strings"
	"time"

	"github.com/pulumi/cloud-ready-checks/pkg/checker"
	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/cloud-ready-checks/pkg/kubernetes/job"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

// ------------------------------------------------------------------------------------------------

// Await logic for batch/v1/Job.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes
// Job as it runs. The idea is that if something goes wrong early, we can cancel the operation or
// alert the user that something is going wrong.
//
// A Job is a construct that allows users to run a workload as a Pod that terminates with a
// success or failure.
//
// A Job is considered "ready" if the following conditions are true:
//
//   1. `.status.startTime` is set, which indicates that the Job has started running.
//   2. `.status.conditions` has a status with `type` equal to `Complete`, and a
//   	`status` set to `True`.
//   3. `.status.conditions` do not have a status with `type` equal to `Failed`, with a
//   	`status` equal to `True`. If this condition is set, we should fail the Job immediately.
//
// The event loop depends on the following channels:
//
//   1. The Job channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any Job it knows about.
//   2. The PodAggregator channel, which monitors Pods related to the Job, and reports any
//		warnings/errors produced by those Pods.
//   3. A timeout channel, which fires after some minutes.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//
// The `jobInitAwaiter` will synchronously process events from the union of all these channels.
// Any time the success conditions described above are reached, we will terminate the awaiter.
//
// The opportunity to display intermediate results will typically appear after a container in the
// Pod fails, (e.g., volume fails to mount, image fails to pull, exited with code 1, etc.).
//
//
// x-refs:
//   * https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/

// --------------------------------------------------------------------------

const (
	DefaultJobTimeoutMins = 10
)

type jobInitAwaiter struct {
	job      *unstructured.Unstructured
	config   awaitConfig
	checker  *checker.StateChecker
	errors   logging.TimeOrderedLogSet
	resource *unstructured.Unstructured
	ready    bool
}

func makeJobInitAwaiter(c awaitConfig) *jobInitAwaiter {
	return &jobInitAwaiter{
		config:   c,
		job:      c.currentOutputs,
		checker:  job.NewJobChecker(),
		resource: c.currentOutputs,
	}
}

func (jia *jobInitAwaiter) Await() error {
	jobEvents := make(chan watch.Event)
	jobInformer, err := jia.config.factory.Subscribe(
		batchv1.SchemeGroupVersion.WithResource("jobs"),
		jobEvents,
	)
	if err != nil {
		return err
	}
	defer jobInformer.Close()

	podEvents := make(chan watch.Event)
	podInformer, err := jia.config.factory.Subscribe(
		corev1.SchemeGroupVersion.WithResource("pods"),
		podEvents,
	)
	if err != nil {
		return err
	}
	defer podInformer.Close()

	podAggregator := NewPodAggregator(jia.job, podInformer)
	podAggregator.Start(podEvents)
	defer podAggregator.Stop()

	timeout := jia.config.getTimeout(DefaultJobTimeoutMins * 60)
	for {
		if jia.ready {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-jia.config.ctx.Done():
			return &cancellationError{
				object:    jia.job,
				subErrors: jia.errorMessages(),
			}
		case <-time.After(timeout):
			return &timeoutError{
				object:    jia.job,
				subErrors: jia.errorMessages(),
			}
		case event := <-jobEvents:
			err := jia.processJobEvent(event)
			if err != nil {
				return err
			}
		case messages := <-podAggregator.ResultChan():
			jia.processPodMessages(messages)
		}
	}
}

func (jia *jobInitAwaiter) Read() error {
	jobClient, err := clients.ResourceClient(kinds.Job, jia.config.currentOutputs.GetNamespace(), jia.config.clientSet)
	if err != nil {
		return fmt.Errorf("Could not make client to get Job %q: %w", jia.config.currentOutputs.GetName(), err)
	}
	// Get live version of Job.
	job, err := jobClient.Get(jia.config.ctx, jia.config.currentOutputs.GetName(), metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the Pod as having been deleted.
		return err
	}

	_ = jia.processJobEvent(watchAddedEvent(job))

	// Check whether we've succeeded.
	if jia.ready {
		return nil
	}

	podInformer, err := jia.config.factory.Subscribe(
		corev1.SchemeGroupVersion.WithResource("pods"),
		nil, // Informer is only used for listing.
	)
	if err != nil {
		return err
	}
	defer podInformer.Close()

	podAggregator := NewPodAggregator(jia.job, podInformer)
	messages := podAggregator.Read()
	for _, message := range messages {
		jia.errors.Add(message)
		jia.config.logger.LogStatus(message.Severity, message.S)
	}

	return &initializationError{
		subErrors: jia.errorMessages(),
		object:    job,
	}
}

func (jia *jobInitAwaiter) processJobEvent(event watch.Event) error {
	if event.Object == nil {
		logger.V(3).Infof("received event with nil Object: %#v", event)
		return nil
	}
	job, err := clients.JobFromUnstructured(event.Object.(*unstructured.Unstructured))
	if err != nil {
		logger.V(3).Infof("Failed to unmarshal Job event: %v", err)
		return nil
	}

	// Do nothing if this is not the job we're waiting for.
	if job.GetName() != jia.config.currentOutputs.GetName() {
		return nil
	}

	var results checker.Results
	jia.ready, results = jia.checker.ReadyDetails(job)
	messages := results.Messages()
	for _, message := range messages.MessagesWithSeverity(diag.Warning, diag.Error) {
		jia.errors.Add(message)
	}
	for _, result := range results {
		jia.config.logger.LogStatus(diag.Info, result.Description)
	}

	if len(messages.Errors()) > 0 {
		return &initializationError{
			subErrors: jia.errorMessages(),
			object:    jia.job,
		}
	}

	return nil
}

func (jia *jobInitAwaiter) processPodMessages(messages checkerlog.Messages) {
	for _, message := range messages {
		jia.errors.Add(message)

		// The unready status condition always occurs as a normal part of a Job running, so don't print
		// this as a warning. If the Job fails to complete, this warning will be included in the subErrors.
		if strings.Contains(message.S, "containers with unready status") {
			continue
		}
		jia.config.logger.LogStatus(message.Severity, message.S)
	}
}

func (jia *jobInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)
	for _, message := range jia.errors.Messages {
		messages = append(messages, message.S)
	}

	return messages
}
