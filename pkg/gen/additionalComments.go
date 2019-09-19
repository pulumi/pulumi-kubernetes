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
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-kubernetes/pkg/await"
	"github.com/pulumi/pulumi-kubernetes/pkg/kinds"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func awaitComments(kind kinds.Kind) string {
	const preamble = `This resource waits until it is ready before registering success for
create/update and populating output properties from the current state of the resource.
The following conditions are used to determine whether the resource creation has
succeeded or failed:
`
	timeoutComment := func(kind kinds.Kind) string {
		const timeoutOverride = `setting the 'customTimeouts' option on the resource.`

		var v int
		switch kind {
		case kinds.Deployment:
			v = await.DefaultDeploymentTimeoutMins
		case kinds.Ingress:
			v = await.DefaultIngressTimeoutMins
		case kinds.Pod:
			v = await.DefaultPodTimeoutMins
		case kinds.Service:
			v = await.DefaultServiceTimeoutMins
		case kinds.StatefulSet:
			v = await.DefaultStatefulSetTimeoutMins
		default:
			// No timeout defined for other resource Kinds.
			return ""
		}
		timeoutStr := strconv.Itoa(v) + " minutes"

		return fmt.Sprintf(`
If the %s has not reached a Ready state after %s, it will
time out and mark the resource update as Failed. You can override the default timeout value
by %s`, kind, timeoutStr, timeoutOverride)
	}

	comment := preamble
	switch kind {
	case kinds.Deployment:
		comment += `
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
   corresponding ReplicaSet), and therefore there is no rollout to mark as 'Progressing'.
`
	case kinds.Ingress:
		comment += `
1.  Ingress object exists.
2.  Endpoint objects exist with matching names for each Ingress path (except when Service
    type is ExternalName).
3.  Ingress entry exists for '.status.loadBalancer.ingress'.
`
	case kinds.Job:
		comment = `This resource currently does not wait until it is ready before registering
success for create/update and populating output properties from the current
state of the resource. Work to add readiness checks is in progress [1].

[1] https://github.com/pulumi/pulumi-kubernetes/pull/633
`
	case kinds.Pod:
		comment += `
1. The Pod is scheduled ("PodScheduled"" '.status.condition' is true).
2. The Pod is initialized ("Initialized" '.status.condition' is true).
3. The Pod is ready ("Ready" '.status.condition' is true) and the '.status.phase' is
   set to "Running".
Or (for Jobs): The Pod succeeded ('.status.phase' set to "Succeeded").
`
	case kinds.Service:
		comment += `
1. Service object exists.
2. Related Endpoint objects are created. Each time we get an update, wait 10 seconds
   for any stragglers.
3. The endpoints objects target some number of living objects (unless the Service is
   an "empty headless" Service [1] or a Service with '.spec.type: ExternalName').
4. External IP address is allocated (if Service has '.spec.type: LoadBalancer').

Known limitations: 
Services targeting ReplicaSets (and, by extension, Deployments,
StatefulSets, etc.) with '.spec.replicas' set to 0 are not handled, and will time
out. To work around this limitation, set 'pulumi.com/skipAwait: "true"' on
'.metadata.annotations' for the Service. Work to handle this case is in progress [2].

[1] https://kubernetes.io/docs/concepts/services-networking/service/#headless-services
[2] https://github.com/pulumi/pulumi-kubernetes/pull/703
`
	case kinds.StatefulSet:
		comment += `
1. The value of 'spec.replicas' matches '.status.replicas', '.status.currentReplicas',
   and '.status.readyReplicas'.
2. The value of '.status.updateRevision' matches '.status.currentRevision'.
`
	default:
		panic("awaitComments: unhandled kind")
	}

	comment += timeoutComment(kind)
	return comment
}

func helpfulLinkComments(kind kinds.Kind) string {
	switch kind {
	case kinds.Secret:
		return `Note: While Pulumi automatically encrypts the 'data' and 'stringData'
fields, this encryption only applies to Pulumi's context, including the state file, 
the Service, the CLI, etc. Kubernetes does not encrypt Secret resources by default,
and the contents are visible to users with access to the Secret in Kubernetes using
tools like 'kubectl'.

For more information on securing Kubernetes Secrets, see the following links:
https://kubernetes.io/docs/concepts/configuration/secret/#security-properties
https://kubernetes.io/docs/concepts/configuration/secret/#risks`
	default:
		return ""
	}
}

// PulumiComment adds additional information to the docs generated automatically from the OpenAPI specs.
// This includes information about Pulumi's await behavior, deprecation information, etc.
func PulumiComment(kind string) string {
	const prefix = "\n\n"

	k := kinds.Kind(kind)
	switch k {
	case kinds.Deployment, kinds.Ingress, kinds.Job, kinds.Pod, kinds.Service, kinds.StatefulSet:
		return prefix + awaitComments(k)
	case kinds.Secret:
		return prefix + helpfulLinkComments(k)
	default:
		return ""
	}
}

func ApiVersionComment(gvk schema.GroupVersionKind) string {
	const template = `%s is not supported by Kubernetes 1.16+ clusters. Use %s instead.

`
	gvkStr := gvk.GroupVersion().String() + "/" + gvk.Kind
	return fmt.Sprintf(template, gvkStr, kinds.SuggestedApiVersion(gvk))
}
