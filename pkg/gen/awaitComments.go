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
1. The Deployment has begun to be updated by the Deployment controller. If the current
   generation of the Deployment is > 1, then this means that the current generation must
   be different from the generation reported by the last outputs.
2. There exists a ReplicaSet whose revision is equal to the current revision of the
   Deployment.
3. The Deployment's '.status.conditions' has a status of type 'Available' whose 'status'
   member is set to 'True'.
4. If the Deployment has generation > 1, then '.status.conditions' has a status of type
   'Progressing', whose 'status' member is set to 'True', and whose 'reason' is
   'NewReplicaSetAvailable'. For generation <= 1, this status field does not exist,
   because it doesn't do a rollout (i.e., it simply creates the Deployment and
   corresponding ReplicaSet), and therefore there is no rollout to mark as 'Progressing'.`

const StatefulSetAwaitComment = `Pulumi uses "await logic" to determine if a StatefulSet is ready.
The following conditions are considered by this logic:
1. The value of 'spec.replicas' matches '.status.replicas', '.status.currentReplicas',
   and '.status.readyReplicas'.
2. The value of '.status.updateRevision' matches '.status.currentRevision'.`

const PodAwaitComment = `Pulumi uses "await logic" to determine if a Pod is ready.
The following conditions are considered by this logic:
1. The Pod is scheduled (PodScheduled condition is true).
2. The Pod is initialized (Initialized condition is true).
3. The Pod is ready (Ready condition is true) and the '.status.phase' is set to "Running".
Or (for Jobs): The Pod succeeded ('.status.phase' set to "Succeeded").`

const ServiceAwaitComment = `Pulumi uses "await logic" to determine if a Service is ready.
The following conditions are considered by this logic:
1. Service object exists.
2. Related Endpoint objects are created. Each time we get an update, wait ~5-10 seconds
   for any stragglers.
3. The endpoints objects target some number of living objects.
4. External IP address is allocated (if Service is type 'LoadBalancer').`

const IngressAwaitComment = `Pulumi uses "await logic" to determine if a Ingress is ready.
The following conditions are considered by this logic:
1.  Ingress object exists.
2.  Endpoint objects exist with matching names for each Ingress path (except when Service
    type is ExternalName).
3.  Ingress entry exists for '.status.loadBalancer.ingress'.`

func AwaitComment(kind string) string {
	const prefix = "\n\n"
	const suffix = ""

	style := func(comment string) string {
		return prefix + comment + suffix
	}

	switch kinds.Kind(kind) {
	case kinds.Deployment:
		return style(DeploymentAwaitComment)
	case kinds.Ingress:
		return style(IngressAwaitComment)
	case kinds.Pod:
		return style(PodAwaitComment)
	case kinds.Service:
		return style(ServiceAwaitComment)
	case kinds.StatefulSet:
		return style(StatefulSetAwaitComment)
	default:
		return ""
	}
}
