// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"unicode"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/lumi/pkg/tokens"

	"github.com/terraform-providers/terraform-provider-google/google"
)

// all of the Google token components used below.
const (
	// packages:
	gcpPkg = "gcp"
	// modules:
	bigqueryMod = "bigquery" // BigQuery
	bigtableMod = "bigtable" // BigTable
	gceMod      = "gce"      // Compute Engine (GCE)
	gkeMod      = "gke"      // Container Engine (GKE)
	dnsMod      = "dns"      // Domain Name Services (DNS)
	projectMod  = "project"  // Projects and Identities
	pubsubMod   = "pubsub"   // Pub/Sub
	sqlMod      = "sql"      // SQL
	storageMod  = "storage"  // Google Cloud Storage
)

// gcptok manufactures a standard resource token given a module and resource name.  It automatically uses the GCP
// package and names the file by simply lower casing the resource's first character.
func gcptok(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return tokens.Type(gcpPkg + ":" + mod + "/" + fn + ":" + res)
}

func gcpProvider() ProviderInfo {
	git, err := getGitInfo("google")
	if err != nil {
		panic(err)
	}
	p := google.Provider().(*schema.Provider)
	prov := ProviderInfo{
		P:   p,
		Git: git,
		Resources: map[string]ResourceInfo{
			// BigQuery
			"google_bigquery_dataset": {Tok: gcptok(bigqueryMod, "DataSet")},
			"google_bigquery_table":   {Tok: gcptok(bigqueryMod, "Table")},
			// BigTable
			"google_bigtable_instance": {Tok: gcptok(bigtableMod, "Instance")},
			"google_bigtable_table":    {Tok: gcptok(bigtableMod, "Table")},
			// GCE
			"google_compute_autoscaler":             {Tok: gcptok(gceMod, "AutoScaler")},
			"google_compute_address":                {Tok: gcptok(gceMod, "Address")},
			"google_compute_backend_bucket":         {Tok: gcptok(gceMod, "BackendBucket")},
			"google_compute_backend_service":        {Tok: gcptok(gceMod, "BackendService")},
			"google_compute_disk":                   {Tok: gcptok(gceMod, "Disk")},
			"google_compute_snapshot":               {Tok: gcptok(gceMod, "Snapshot")},
			"google_compute_firewall":               {Tok: gcptok(gceMod, "Firewall")},
			"google_compute_forwarding_rule":        {Tok: gcptok(gceMod, "ForwardingRule")},
			"google_compute_global_address":         {Tok: gcptok(gceMod, "GlobalAddress")},
			"google_compute_global_forwarding_rule": {Tok: gcptok(gceMod, "GlobalForwardingRule")},
			"google_compute_health_check":           {Tok: gcptok(gceMod, "HealthCheck")},
			"google_compute_http_health_check":      {Tok: gcptok(gceMod, "HttpHealthCheck")},
			"google_compute_https_health_check":     {Tok: gcptok(gceMod, "HttpsHealthCheck")},
			"google_compute_image":                  {Tok: gcptok(gceMod, "Image")},
			"google_compute_instance":               {Tok: gcptok(gceMod, "Instance")},
			"google_compute_instance_group":         {Tok: gcptok(gceMod, "InstanceGroup")},
			"google_compute_instance_group_manager": {Tok: gcptok(gceMod, "InstanceGroupManager")},
			"google_compute_instance_template":      {Tok: gcptok(gceMod, "InstanceTemplate")},
			"google_compute_network":                {Tok: gcptok(gceMod, "Network")},
			"google_compute_project_metadata":       {Tok: gcptok(gceMod, "ProjectMetadata")},
			"google_compute_region_backend_service": {Tok: gcptok(gceMod, "RegionBackendService")},
			"google_compute_route":                  {Tok: gcptok(gceMod, "Route")},
			"google_compute_router":                 {Tok: gcptok(gceMod, "Router")},
			"google_compute_router_interface":       {Tok: gcptok(gceMod, "RouterInterface")},
			"google_compute_router_peer":            {Tok: gcptok(gceMod, "RouterPeer")},
			"google_compute_ssl_certificate": {
				Tok: gcptok(gceMod, "SslCertificate"),
				Fields: map[string]SchemaInfo{
					"id": {Name: "certificateId"},
				},
			},
			"google_compute_subnetwork": {Tok: gcptok(gceMod, "SubNetwork")},
			"google_compute_target_http_proxy": {
				Tok: gcptok(gceMod, "TargetHttpProxy"),
				Fields: map[string]SchemaInfo{
					"id": {Name: "proxyId"},
				},
			},
			"google_compute_target_https_proxy": {
				Tok: gcptok(gceMod, "TargetHttpsProxy"),
				Fields: map[string]SchemaInfo{
					"id": {Name: "proxyId"},
				},
			},
			"google_compute_target_pool": {Tok: gcptok(gceMod, "TargetPool")},
			"google_compute_url_map": {
				Tok: gcptok(gceMod, "UrlMap"),
				Fields: map[string]SchemaInfo{
					"id": {Name: "mapId"},
				},
			},
			"google_compute_vpn_gateway": {Tok: gcptok(gceMod, "VpnGateway")},
			"google_compute_vpn_tunnel":  {Tok: gcptok(gceMod, "VpnTunnel")},
			// GKE
			"google_container_cluster":   {Tok: gcptok(gkeMod, "Cluster")},
			"google_container_node_pool": {Tok: gcptok(gkeMod, "NodePool")},
			// DNS
			"google_dns_managed_zone": {Tok: gcptok(dnsMod, "ManagedZone")},
			"google_dns_record_set":   {Tok: gcptok(dnsMod, "RecordSet")},
			// Projects and Identities
			"google_project":            {Tok: gcptok(projectMod, "Project")},
			"google_project_iam_policy": {Tok: gcptok(projectMod, "IamPolicy")},
			"google_project_services":   {Tok: gcptok(projectMod, "Services")},
			"google_service_account":    {Tok: gcptok(projectMod, "ServiceAccount")},
			// Pub/Sub
			"google_pubsub_topic":        {Tok: gcptok(pubsubMod, "Topic")},
			"google_pubsub_subscription": {Tok: gcptok(pubsubMod, "Subscription")},
			// SQL
			"google_sql_database":          {Tok: gcptok(sqlMod, "Database")},
			"google_sql_database_instance": {Tok: gcptok(sqlMod, "DatabaseInstance")},
			"google_sql_user":              {Tok: gcptok(sqlMod, "User")},
			// Google Cloud Storage
			"google_storage_bucket":        {Tok: gcptok(storageMod, "Bucket")},
			"google_storage_bucket_acl":    {Tok: gcptok(storageMod, "BucketAcl")},
			"google_storage_bucket_object": {Tok: gcptok(storageMod, "Object")},
			"google_storage_object_acl":    {Tok: gcptok(storageMod, "ObjectAcl")},
		},
	}

	// For all resources with name properties, we will add an auto-name property.
	for resname := range prov.Resources {
		if schema := p.ResourcesMap[resname]; schema != nil {
			if _, has := schema.Schema[NameProperty]; has {
				prov.Resources[resname] = autoName(prov.Resources[resname], -1)
			}
		}
	}

	return prov
}
