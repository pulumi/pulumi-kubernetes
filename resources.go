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
		P:           p,
		Name:        "kubernetes",
		Description: "A Pulumi package for creating and managing Kubernetes resources.",
		Keywords:    []string{"pulumi", "kubernetes"},
		Homepage:    "https://pulumi.io/kubernetes",
		Repository:  "https://github.com/pulumi/pulumi-kubernetes",
		Resources: map[string]*tfbridge.ResourceInfo{
			// TODO[pulumi/pulumi-kubernetes#10] Until we are auto-generating `metadata.name` with a random suffix by
			// default, we must mark all Kubernetes resources as delete-before-replace.  Without this, any change which
			// forces a replace will re-use the same resource name, which will fail.
			"kubernetes_config_map":                {Tok: kubernetesResource(kubernetesCore, "ConfigMap"), DeleteBeforeReplace: true},
			"kubernetes_daemonset":                 {Tok: kubernetesResource(kubernetesCore, "DaemonSet"), DeleteBeforeReplace: true},
			"kubernetes_deployment":                {Tok: kubernetesResource(kubernetesCore, "Deployment"), DeleteBeforeReplace: true},
			"kubernetes_horizontal_pod_autoscaler": {Tok: kubernetesResource(kubernetesCore, "HorizontalPodAutoscaler"), DeleteBeforeReplace: true},
			"kubernetes_ingress":                   {Tok: kubernetesResource(kubernetesCore, "Ingres"), DeleteBeforeReplace: true},
			"kubernetes_job":                       {Tok: kubernetesResource(kubernetesCore, "Job"), DeleteBeforeReplace: true},
			"kubernetes_limit_range":               {Tok: kubernetesResource(kubernetesCore, "LimitRange"), DeleteBeforeReplace: true},
			"kubernetes_namespace":                 {Tok: kubernetesResource(kubernetesCore, "Namespace"), DeleteBeforeReplace: true},
			"kubernetes_persistent_volume":         {Tok: kubernetesResource(kubernetesCore, "PersistentVolume"), DeleteBeforeReplace: true},
			"kubernetes_persistent_volume_claim":   {Tok: kubernetesResource(kubernetesCore, "PersistentVolumeClaim"), DeleteBeforeReplace: true},
			"kubernetes_pod":                       {Tok: kubernetesResource(kubernetesCore, "Pod"), DeleteBeforeReplace: true},
			"kubernetes_replication_controller":    {Tok: kubernetesResource(kubernetesCore, "ReplicationController"), DeleteBeforeReplace: true},
			"kubernetes_resource_quota":            {Tok: kubernetesResource(kubernetesCore, "ResourceQuota"), DeleteBeforeReplace: true},
			"kubernetes_secret":                    {Tok: kubernetesResource(kubernetesCore, "Secret"), DeleteBeforeReplace: true},
			"kubernetes_service":                   {Tok: kubernetesResource(kubernetesCore, "Service"), DeleteBeforeReplace: true},
			"kubernetes_service_account":           {Tok: kubernetesResource(kubernetesCore, "ServiceAccount"), DeleteBeforeReplace: true},
			"kubernetes_stateful_set":              {Tok: kubernetesResource(kubernetesCore, "StatefulSet"), DeleteBeforeReplace: true},
			"kubernetes_storage_class":             {Tok: kubernetesResource(kubernetesCore, "StorageClass"), DeleteBeforeReplace: true},
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			"kubernetes_service":       {Tok: kubernetesDataSource(kubernetesCore, "getService")},
			"kubernetes_storage_class": {Tok: kubernetesDataSource(kubernetesCore, "getStorageClass")},
		},
		Overlay: &tfbridge.OverlayInfo{
			Files:   []string{},
			Modules: map[string]*tfbridge.OverlayInfo{},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			PeerDependencies: map[string]string{
				"@pulumi/pulumi": "^0.11.0-dev-168-g7e14a09b",
			},
		},
		Python: &tfbridge.PythonInfo{
			Requires: map[string]string{
				"pulumi": ">=0.11.0",
			},
		},
	}

	// TODO[pulumi/pulumi-kubernetes#10: Auto-populate `res.metadata.name`

	return prov
}
