// Copyright 2016-2018, Pulumi Corporation.
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

package provider

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
)

func (k *kubeProvider) forceNewProperties(obj *unstructured.Unstructured) []string {
	gvk := obj.GroupVersionKind()
	props := metadataForceNewProperties(".metadata")
	if group, groupExists := forceNew[gvk.Group]; groupExists {
		if version, versionExists := group[gvk.Version]; versionExists {
			if clients.IsSecret(obj) && (!k.enableSecretMutable || clients.IsImmutable(obj)) {
				props = append(props, properties{".type", ".stringData", ".data"}...)
			} else if kindFields, kindExists := version[gvk.Kind]; kindExists {
				props = append(props, kindFields...)
			} else if clients.IsConfigMap(obj) && (!k.enableConfigMapMutable || clients.IsImmutable(obj)) {
				props = append(props, properties{".binaryData", ".data"}...)
			}
		}
	}
	return props
}

type _groups map[string]_versions
type _versions map[string]_kinds
type _kinds map[string]properties
type properties []string

var forceNew = _groups{
	"apps": _versions{
		// NOTE: .spec.selector triggers a replacement in Deployment/DaemonSet only AFTER v1beta1.
		"v1beta1": _kinds{"StatefulSet": statefulSet},
		"v1beta2": _kinds{
			"DaemonSet":   daemonset,
			"Deployment":  deployment,
			"StatefulSet": statefulSet},
		"v1": _kinds{
			"DaemonSet":   daemonset,
			"Deployment":  deployment,
			"StatefulSet": statefulSet},
	},
	// List `core` under its canonical name and under it's legacy name (i.e., "", the empty string)
	// for compatibility purposes.
	"core": core,
	"":     core,
	"policy": _versions{
		"v1beta1": _kinds{"PodDisruptionBudget": podDisruptionBudget},
	},
	"rbac.authorization.k8s.io": _versions{
		"v1alpha1": _kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
		"v1beta1":  _kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
		"v1":       _kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
	},
	"storage.k8s.io": _versions{
		"v1": _kinds{
			"StorageClass": properties{
				".parameters",
				".provisioner",
			},
		},
	},
	"batch": _versions{
		"v1beta1":  _kinds{"Job": job},
		"v1":       _kinds{"Job": job},
		"v2alpha1": _kinds{"Job": job},
	},
}

var core = _versions{
	"v1": _kinds{
		"PersistentVolume": properties{
			".spec.awsElasticBlockStore",
			".spec.azureDisk",
			".spec.azureFile",
			".spec.cephfs",
			".spec.cinder",
			".spec.csi",
			".spec.fc",
			".spec.flexVolume",
			".spec.flocker",
			".spec.gcePersistentDisk",
			".spec.glusterfs",
			".spec.hostPath",
			".spec.iscsi",
			".spec.local",
			".spec.nfs",
			".spec.photonPersistentDisk",
			".spec.portworxVolume",
			".spec.quobyte",
			".spec.rbd",
			".spec.scaleIO",
			".spec.storageos",
			".spec.vsphereVolume",
		},
		"PersistentVolumeClaim": append(
			properties{
				".spec.accessModes",
				".spec.selector",
				".spec.storageClassName",
				".spec.volumeName",
			},
			labelSelectorForceNewProperties(".spec.selector")...),
		"Pod": append(
			properties{
				".spec.affinity",
				".spec.automountServiceAccountToken",
				".spec.dnsConfig",
				".spec.dnsPolicy",
				".spec.enableServiceLinks",
				".spec.hostAliases",
				".spec.hostIPC",
				".spec.hostNetwork",
				".spec.hostPID",
				".spec.hostname",
				".spec.imagePullSecrets",
				".spec.imagePullSecrets",
				".spec.nodeName",
				".spec.nodeSelector",
				".spec.overhead",
				".spec.preemptionPolicy",
				".spec.priority",
				".spec.priorityClassName",
				".spec.readinessGates",
				".spec.restartPolicy",
				".spec.runtimeClassName",
				".spec.schedulerName",
				".spec.securityContext",
				".spec.serviceAccount",
				".spec.serviceAccountName",
				".spec.shareProcessNamespace",
				".spec.subdomain",
				".spec.terminationGracePeriodSeconds",
				".spec.volumes",
			},
			containerForceNewProperties(".spec.containers[*]", ".spec.initContainers[*]")...),
		"ResourceQuota": properties{
			".spec.scopes",
		},
		"Secret": properties{
			".type",
		},
		"Service": properties{
			".spec.clusterIP",
		},
	},
}

var daemonset = append(
	properties{
		".spec.selector",
	},
	labelSelectorForceNewProperties(".spec.selector")...,
)

var deployment = append(
	properties{
		".spec.selector",
	},
	labelSelectorForceNewProperties(".spec.selector")...,
)

var job = append(
	properties{
		".spec.completions",
		".spec.parallelism",
		".spec.selector",
		".spec.template",
	},
	labelSelectorForceNewProperties(".spec.selector")...,
)

var podDisruptionBudget = properties{
	".spec",
}

var roleBinding = properties{
	".roleRef",
}

var statefulSet = append(
	properties{
		".spec.podManagementPolicy",
		".spec.revisionHistoryLimit",
		".spec.selector",
		".spec.serviceName",
		".spec.volumeClaimTemplates",
	},
	labelSelectorForceNewProperties(".spec.selector")...,
)

func metadataForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".name",
		prefix + ".namespace",
	}
}

func containerForceNewProperties(prefixes ...string) properties {
	var props properties
	for _, prefix := range prefixes {
		props = append(props, properties{
			prefix + ".args",
			prefix + ".command",
			prefix + ".env",
			prefix + ".env.value",
			prefix + ".image",
			prefix + ".imagePullPolicy",
			prefix + ".lifecycle",
			prefix + ".livenessProbe",
			prefix + ".ports",
			prefix + ".readinessProbe",
			prefix + ".resources",
			prefix + ".securityContext",
			prefix + ".stdin",
			prefix + ".stdinOnce",
			prefix + ".terminationMessagePath",
			prefix + ".terminationMessagePolicy",
			prefix + ".tty",
			prefix + ".volumeDevices",
			prefix + ".volumeMounts",
			prefix + ".workingDir",
		}...)
	}
	return props
}

func labelSelectorForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".matchExpressions",
		prefix + ".matchExpressions.key",
		prefix + ".matchExpressions.operator",
		prefix + ".matchExpressions.values",
		prefix + ".matchLabels",
	}
}
