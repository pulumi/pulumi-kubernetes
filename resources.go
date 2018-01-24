// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

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
	kubernetesPkg 			= "kubernetes"
	kubernetesCore          = "core"             //Resources
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
			"kubernetes_config_map": {Tok: kubernetesResource(kubernetesCore, "ConfigMap")},
			"kubernetes_horizontal_pod_autoscaler":      {Tok: kubernetesResource(kubernetesCore, "HorizontalPodAutoascaler")},
			"kubernetes_limit_range": {Tok: kubernetesResource(kubernetesCore, "LimitRange")},
			"kubernetes_namespace": {Tok: kubernetesResource(kubernetesCore, "Namespace")},
			"kubernetes_persistent_volume":    {Tok: kubernetesResource(kubernetesCore, "PersistentVolume")},
			"kubernetes_persistent_volume_claim":    {Tok: kubernetesResource(kubernetesCore, "PersistentVolumeClaim")},
			"kubernetes_pod": {Tok: kubernetesResource(kubernetesCore, "Pod")},
			"kubernetes_replication_controller":    {Tok: kubernetesResource(kubernetesCore, "ReplicationController")},
			"kubernetes_resource_quota":   {Tok: kubernetesResource(kubernetesCore, "ResourceQuota")},
			"kubernetes_secret": {Tok: kubernetesResource(kubernetesCore, "Secret")},
			"kubernetes_service": {Tok: kubernetesResource(kubernetesCore, "Service")},
			"kubernetes_service_account": {Tok: kubernetesResource(kubernetesCore, "ServiceAccount")},
			"kubernetes_storage_class":  {Tok: kubernetesResource(kubernetesCore, "StorageClass")},
			},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			"kubernetes_service":           {Tok: kubernetesDataSource(kubernetesCore, "getService")},
			"kubernetes_storage_class":                   {Tok: kubernetesDataSource(kubernetesCore, "getStorageClass")},
			},
		Overlay: &tfbridge.OverlayInfo{
			Files:        []string{},
			Modules:      map[string]*tfbridge.OverlayInfo{},
			Dependencies: map[string]string{},
		},
	}

	// For all resources with name properties, we will add an auto-name property.  Make sure to skip those that
	// already have a name mapping entry, since those may have custom overrides set above (e.g., for length).
	const kubernetesName = "name"
	for resname, res := range prov.Resources {
		if schema := p.ResourcesMap[resname]; schema != nil {
			if _, has := schema.Schema[kubernetesName]; has {
				if _, hasfield := res.Fields[kubernetesName]; !hasfield {
					if res.Fields == nil {
						res.Fields = make(map[string]*tfbridge.SchemaInfo)
					}
					res.Fields[kubernetesName] = tfbridge.AutoName(kubernetesName, 255)
				}
			}
		}
	}

	return prov
}
