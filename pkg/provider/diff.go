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
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func forceNewProperties(patch map[string]interface{}, gvk schema.GroupVersionKind) ([]string, error) {
	props := metadataForceNewProperties(".metadata")
	if group, groupExists := forceNew[gvk.Group]; groupExists {
		if version, versionExists := group[gvk.Version]; versionExists {
			if kindFields, kindExists := version[gvk.Kind]; kindExists {
				props = append(props, kindFields...)
			}
		}
	}

	return openapi.PatchPropertiesChanged(patch, props)
}

type groups map[string]versions
type versions map[string]kinds
type kinds map[string]properties
type properties []string

var forceNew = groups{
	"apps": versions{
		// NOTE: .spec.selector triggers a replacement in Deployment only AFTER v1beta1.
		"v1beta1": kinds{"StatefulSet": statefulSet},
		"v1beta2": kinds{
			"Deployment":  deployment,
			"StatefulSet": statefulSet},
		"v1": kinds{
			"Deployment":  deployment,
			"StatefulSet": statefulSet},
	},
	// List `core` under its canonical name and under it's legacy name (i.e., "", the empty string)
	// for compatibility purposes.
	"core": core,
	"":     core,
	"policy": versions{
		"v1beta1": kinds{"PodDisruptionBudget": podDisruptionBudget},
	},
	"rbac.authorization.k8s.io": versions{
		"v1alpha1": kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
		"v1beta1":  kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
		"v1":       kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
	},
	"storage.k8s.io": versions{
		"v1": kinds{
			"StorageClass": properties{
				".parameters",
				".provisioner",
			},
		},
	},
	"batch": versions{
		"v1beta1":  kinds{"Job": job},
		"v1":       kinds{"Job": job},
		"v2alpha1": kinds{"Job": job},
	},
}

var core = versions{
	"v1": kinds{
		"ConfigMap": properties{".binaryData", ".data"},
		"PersistentVolumeClaim": append(
			properties{
				".spec",
				".spec.accessModes",
				".spec.resources",
				".spec.resources.limits",
				".spec.resources.requests",
				".spec.selector",
				".spec.storageClassName",
				".spec.volumeName",
			},
			labelSelectorForceNewProperties(".spec")...),
		"Pod": containerForceNewProperties(".spec.containers[*]"),
		"ResourceQuota": properties{
			".spec.scopes",
		},
		"Secret": properties{
			".type",
			".stringData",
			".data",
		},
		"Service": properties{
			".spec.clusterIP",
			".spec.type",
		},
	},
}

var deployment = properties{
	".spec.selector",
}

var job = properties{
	".spec.selector",
	".spec.template",
}

var podDisruptionBudget = properties{
	".spec",
}

var roleBinding = properties{
	".roleRef",
}

var statefulSet = properties{
	".spec.podManagementPolicy",
	".spec.revisionHistoryLimit",
	".spec.selector",
	".spec.serviceName",
	".spec.volumeClaimTemplates",
}

func metadataForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".name",
		prefix + ".namespace",
	}
}

func containerForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".env",
		prefix + ".env.value",
		prefix + ".image",
		prefix + ".lifecycle",
		prefix + ".livenessProbe",
		prefix + ".readinessProbe",
		prefix + ".securityContext",
		prefix + ".terminationMessagePath",
		prefix + ".workingDir",
	}
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
