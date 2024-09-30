// Copyright 2016-2024, Pulumi Corporation.
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
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestAPIVersionComment(t *testing.T) {
	tests := []struct {
		gvk      schema.GroupVersionKind
		expected string
	}{
		{
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1beta1", Kind: "Deployment"},
			expected: "apps/v1beta1/Deployment is deprecated by apps/v1/Deployment and not supported by Kubernetes v1.16+ clusters.",
		},
		{
			gvk:      schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"},
			expected: "extensions/v1beta1/Ingress is deprecated by networking.k8s.io/v1beta1/Ingress and not supported by Kubernetes v1.20+ clusters.",
		},
		{
			gvk:      schema.GroupVersionKind{Group: "batch", Version: "v1beta1", Kind: "CronJob"},
			expected: "batch/v1beta1/CronJob is deprecated by batch/v1beta1/CronJob.",
		},
		{
			gvk:      schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"},
			expected: "core/v1/Pod is deprecated by core/v1/Pod.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			actual := APIVersionComment(tt.gvk)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

// TODO: The tests, and function itself, are rather brittle. We should find a way to make this more robust.
func TestPulumiComment(t *testing.T) {
	tests := []struct {
		kind     string
		expected string
	}{
		{
			kind: "Deployment",
			expected: "\n\nThis resource waits until its status is ready before registering success\n" +
				"for create/update, and populating output properties from the current state of the resource.\n" +
				"The following conditions are used to determine whether the resource creation has\n" +
				"succeeded or failed:\n\n" +
				"1. The Deployment has begun to be updated by the Deployment controller. If the current\n" +
				"   generation of the Deployment is > 1, then this means that the current generation must\n" +
				"   be different from the generation reported by the last outputs.\n" +
				"2. There exists a ReplicaSet whose revision is equal to the current revision of the\n" +
				"   Deployment.\n" +
				"3. The Deployment's '.status.conditions' has a status of type 'Available' whose 'status'\n" +
				"   member is set to 'True'.\n" +
				"4. If the Deployment has generation > 1, then '.status.conditions' has a status of type\n" +
				"   'Progressing', whose 'status' member is set to 'True', and whose 'reason' is\n" +
				"   'NewReplicaSetAvailable'. For generation <= 1, this status field does not exist,\n" +
				"   because it doesn't do a rollout (i.e., it simply creates the Deployment and\n" +
				"   corresponding ReplicaSet), and therefore there is no rollout to mark as 'Progressing'.\n\n" +
				"If the Deployment has not reached a Ready state after 10 minutes, it will\n" +
				"time out and mark the resource update as Failed. You can override the default timeout value\n" +
				"by setting the 'customTimeouts' option on the resource.",
		},
		{
			kind: "Secret",
			expected: "\n\nNote: While Pulumi automatically encrypts the 'data' and 'stringData'\n" +
				"fields, this encryption only applies to Pulumi's context, including the state file, \n" +
				"the Service, the CLI, etc. Kubernetes does not encrypt Secret resources by default,\n" +
				"and the contents are visible to users with access to the Secret in Kubernetes using\n" +
				"tools like 'kubectl'.\n\n" +
				"For more information on securing Kubernetes Secrets, see the following links:\n" +
				"https://kubernetes.io/docs/concepts/configuration/secret/#security-properties\n" +
				"https://kubernetes.io/docs/concepts/configuration/secret/#risks",
		},
		{
			kind: "Job",
			expected: "\n\nThis resource waits until its status is ready before registering success\n" +
				"for create/update, and populating output properties from the current state of the resource.\n" +
				"The following conditions are used to determine whether the resource creation has\n" +
				"succeeded or failed:\n\n" +
				"1. The Job's '.status.startTime' is set, which indicates that the Job has started running.\n" +
				"2. The Job's '.status.conditions' has a status of type 'Complete', and a 'status' set\n" +
				"   to 'True'.\n" +
				"3. The Job's '.status.conditions' do not have a status of type 'Failed', with a\n" +
				"	'status' set to 'True'. If this condition is set, we should fail the Job immediately.\n\n" +
				"If the Job has not reached a Ready state after 10 minutes, it will\n" +
				"time out and mark the resource update as Failed. You can override the default timeout value\n" +
				"by setting the 'customTimeouts' option on the resource.\n\n" +
				"By default, if a resource failed to become ready in a previous update, \n" +
				"Pulumi will continue to wait for readiness on the next update. If you would prefer\n" +
				"to schedule a replacement for an unready resource on the next update, you can add the\n" +
				"\"pulumi.com/replaceUnready\": \"true\" annotation to the resource definition.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			actual := PulumiComment(tt.kind)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
