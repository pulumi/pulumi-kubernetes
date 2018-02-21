// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package kubernetes

import (
	"unicode"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
	"github.com/pulumi/pulumi/pkg/tokens"
	kubernetes "github.com/terraform-providers/terraform-provider-kubernetes/kubernetes"
)

// all of the Kubernetes token components used below.
const (
	// packages:
	kubernetesPkg  = "kubernetes"
	kubernetesCore = "" // All resources are placed in the root `kubernetes` namespace
)

// kubernetesMember manufactures a type token for the Kubernetes package and the given module and type.
func kubernetesMember(mod string, mem string) tokens.ModuleMember {
	return tokens.ModuleMember(kubernetesPkg + ":" + mod + ":" + mem)
}

// kubernetesType manufactures a type token for the Kubernetes package and the given module and type.
func kubernetesType(mod string, typ string) tokens.Type {
	return tokens.Type(kubernetesMember(mod, typ))
}

// kubernetesDataSource manufactures a standard resource token given a module and resource name.  It automatically uses the Kubernetes
// package and names the file by simply lower casing the data source's first character.
func kubernetesDataSource(mod string, res string) tokens.ModuleMember {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return kubernetesMember(mod+"/"+fn, res)
}

// kubernetesResource manufactures a standard resource token given a module and resource name.  It automatically uses the Kubernetes
// package and names the file by simply lower casing the resource's first character.
func kubernetesResource(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return kubernetesType(mod+"/"+fn, res)
}

// Provider returns additional overlaid schema and metadata associated with the Kubernetes package.
func Provider() tfbridge.ProviderInfo {
	p := kubernetes.Provider().(*schema.Provider)
	prov := tfbridge.ProviderInfo{
		P:    p,
		Name: "kubernetes",
		Resources: map[string]*tfbridge.ResourceInfo{
			"kubernetes_config_map":                {Tok: kubernetesResource(kubernetesCore, "ConfigMap")},
			"kubernetes_daemonset":                 {Tok: kubernetesResource(kubernetesCore, "DaemonSet")},
			"kubernetes_deployment":                {Tok: kubernetesResource(kubernetesCore, "Deployment")},
			"kubernetes_horizontal_pod_autoscaler": {Tok: kubernetesResource(kubernetesCore, "HorizontalPodAutoscaler")},
			"kubernetes_ingress":                   {Tok: kubernetesResource(kubernetesCore, "Ingres")},
			"kubernetes_job":                       {Tok: kubernetesResource(kubernetesCore, "Job")},
			"kubernetes_limit_range":               {Tok: kubernetesResource(kubernetesCore, "LimitRange")},
			"kubernetes_namespace":                 {Tok: kubernetesResource(kubernetesCore, "Namespace")},
			"kubernetes_persistent_volume":         {Tok: kubernetesResource(kubernetesCore, "PersistentVolume")},
			"kubernetes_persistent_volume_claim":   {Tok: kubernetesResource(kubernetesCore, "PersistentVolumeClaim")},
			"kubernetes_pod":                       {Tok: kubernetesResource(kubernetesCore, "Pod")},
			"kubernetes_replication_controller":    {Tok: kubernetesResource(kubernetesCore, "ReplicationController")},
			"kubernetes_resource_quota":            {Tok: kubernetesResource(kubernetesCore, "ResourceQuota")},
			"kubernetes_secret":                    {Tok: kubernetesResource(kubernetesCore, "Secret")},
			"kubernetes_service":                   {Tok: kubernetesResource(kubernetesCore, "Service")},
			"kubernetes_service_account":           {Tok: kubernetesResource(kubernetesCore, "ServiceAccount")},
			"kubernetes_stateful_set":              {Tok: kubernetesResource(kubernetesCore, "StatefulSet")},
			"kubernetes_storage_class":             {Tok: kubernetesResource(kubernetesCore, "StorageClass")},
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			"kubernetes_service":       {Tok: kubernetesDataSource(kubernetesCore, "getService")},
			"kubernetes_storage_class": {Tok: kubernetesDataSource(kubernetesCore, "getStorageClass")},
		},
		Overlay: &tfbridge.OverlayInfo{
			Files:   []string{},
			Modules: map[string]*tfbridge.OverlayInfo{},
			Dependencies: map[string]string{
				"@pulumi/pulumi": "^0.11.0-dev-23-g444ebdd1",
			},
		},
	}

	// TODO[pulumi/pulumi-kubernetes#10: Auto-populate `res.metadata.name`

	return prov
}
