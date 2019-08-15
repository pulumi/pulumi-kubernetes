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

package gen

import (
	"github.com/pulumi/pulumi-kubernetes/pkg/kinds"
)

const DeploymentAwaitComment = `Pulumi uses "await logic" to determine if a Deployment is ready.
The following conditions are considered by this logic:
1. '.metadata.annotations["deployment.kubernetes.io/revision"]' in the current Deployment
  must have been incremented by the Deployment controller, i.e., it must not be equal to
  the revision number in the previous outputs. This number is used to indicate the the
  active ReplicaSet. Any time a change is made to the Deployment's Pod template, this
  revision is incremented, and a new ReplicaSet is created with a corresponding revision
  number in its own annotations. This condition overall is a test to make sure that the
  Deployment controller is making progress in rolling forward to the new generation.

2. '.status.conditions' has a status with 'type' equal to 'Progressing', a 'status' set to
  'True', and a 'reason' set to 'NewReplicaSetAvailable'. Though the condition is called
  "Progressing", this condition indicates that the new ReplicaSet has become healthy and
  available, and the Deployment controller is now free to delete the old ReplicaSet.

3. '.status.conditions' has a status with 'type' equal to 'Available', a 'status' equal to
  'True'. If the Deployment is not available, we should fail the Deployment immediately.`

const StatefulSetAwaitComment = `Pulumi uses "await logic" to determine if a StatefulSet is ready.
The following conditions are considered by this logic:
1. '.status.replicas', '.status.currentReplicas' and '.status.readyReplicas' match the
   value of '.spec.replicas'.
2. '.status.updateRevision' matches '.status.currentRevision'.`

const PodAwaitComment = `Pulumi uses "await logic" to determine if a Pod is ready.
The following conditions are considered by this logic:
1. The Pod is scheduled (PodScheduled condition is true).
2. The Pod is initialized (Initialized condition is true).
3. The Pod is ready (Ready condition is true) and the '.status.phase' is set to "Running".
Or (for Jobs): The Pod succeeded ('.status.phase' set to "Succeeded").`

const ServiceAwaitComment = `Pulumi uses "await logic" to determine if a Pod is ready.
The following conditions are considered by this logic:
1. Service object exists.
2. Related Endpoint objects are created. Each time we get an update, wait ~5-10 seconds
   for any stragglers.
3. The endpoints objects target some number of living objects.
4. External IP address is allocated (if Service is type 'LoadBalancer').`

func AwaitComment(kind string, opts groupOpts) string {
	var prefix string
	const suffix = "\n"

	switch opts.language {
	case typescript:
		prefix = "*\n"
	case python:
		prefix = "\n"
	}

	style := func(comment string) string {
		return prefix + comment + suffix
	}

	switch kinds.Kind(kind) {
	case kinds.Deployment:
		return style(DeploymentAwaitComment)
	case kinds.Pod:
		return style(PodAwaitComment)
	case kinds.Service:
		return style(ServiceAwaitComment)
	case kinds.StatefulSet:
		return style(StatefulSetAwaitComment)
	default:
		if opts.language == typescript {
			return "*"
		}
		return ""
	}
}
