// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"unicode"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/lumi/pkg/tokens"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
)

// all of the Azure token components used below.
const (
	// packages:
	azurePkg = "tf-azure"
	// modules; in general, we took naming inspiration from the Azure SDK for Go:
	// https://godoc.org/github.com/Azure/azure-sdk-for-go
	azureAppInsightsMod      = "appinsights"      // AppInsights
	azureContainerServiceMod = "containerservice" // Azure Container Service
	azureCDN                 = "cdn"              // CDN
	azureCosmosDB            = "cosmosdb"         // Cosmos DB
	azureDNS                 = "dns"              // DNS
	azureEventHub            = "eventhub"         // Event Hub
	azureExpressRoute        = "expressroute"     // Express Route
	azureKeyVault            = "keyvault"         // Key Vault
	azureLB                  = "lb"               // Load Balancer
	azureDisk                = "disk"             // Managed Disks
	azureNetwork             = "network"          // Networking
	azureRedis               = "redis"            // RedisCache
	azureResources           = "resources"        // Azure Resource Manager
	azureSearch              = "search"           // Search
	azureServiceBus          = "servicebus"       // ServiceBus
	azureSQL                 = "sql"              // SQL
	azureStorage             = "storage"          // Storage
	azureTrafficManager      = "trafficmanager"   // Traffic Manager
	azureVirtualMachine      = "virtualmachine"   // Virtual Machines
)

// azuretok manufactures a standard resource token given a module and resource name.  It automatically uses the Azure
// package and names the file by simply lower casing the resource's first character.
func azuretok(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return tokens.Type(azurePkg + ":" + mod + "/" + fn + ":" + res)
}

func azureProvider() ProviderInfo {
	git, err := getGitInfo("azurerm")
	if err != nil {
		panic(err)
	}
	return ProviderInfo{
		P:   azurerm.Provider().(*schema.Provider),
		Git: git,
		Resources: map[string]ResourceInfo{
			// AppInsights
			"azurerm_application_insights": {Tok: azuretok(azureAppInsightsMod, "Insights")},
			// Azure Container Service
			"azurerm_container_registry": {Tok: azuretok(azureContainerServiceMod, "Registry")},
			"azurerm_container_service":  {Tok: azuretok(azureContainerServiceMod, "Service")},
			// CDN
			"azurerm_cdn_endpoint": {Tok: azuretok(azureCDN, "Endpoint")},
			"azurerm_cdn_profile":  {Tok: azuretok(azureCDN, "Profile")},
			// CosmosDB
			"azurerm_cosmosdb_account": {Tok: azuretok(azureCosmosDB, "Account")},
			// DNS
			"azurerm_dns_a_record":     {Tok: azuretok(azureDNS, "ARecord")},
			"azurerm_dns_aaaa_record":  {Tok: azuretok(azureDNS, "AaaaRecord")},
			"azurerm_dns_cname_record": {Tok: azuretok(azureDNS, "CNameRecord")},
			"azurerm_dns_mx_record":    {Tok: azuretok(azureDNS, "MxRecord")},
			"azurerm_dns_ns_record":    {Tok: azuretok(azureDNS, "NsRecord")},
			"azurerm_dns_ptr_record":   {Tok: azuretok(azureDNS, "PtrRecord")},
			"azurerm_dns_srv_record":   {Tok: azuretok(azureDNS, "SrvRecord")},
			"azurerm_dns_txt_record":   {Tok: azuretok(azureDNS, "TxtRecord")},
			"azurerm_dns_zone":         {Tok: azuretok(azureDNS, "Zone")},
			// EventHub
			"azurerm_eventhub":                    {Tok: azuretok(azureEventHub, "EventHub")},
			"azurerm_eventhub_authorization_rule": {Tok: azuretok(azureEventHub, "AuthorizationRule")},
			"azurerm_eventhub_consumer_group":     {Tok: azuretok(azureEventHub, "ConsumerGroup")},
			"azurerm_eventhub_namespace":          {Tok: azuretok(azureEventHub, "Namespace")},
			// ExpressRoute
			"azurerm_express_route_circuit": {Tok: azuretok(azureExpressRoute, "Circuit")},
			// KeyVault
			"azurerm_key_vault": {Tok: azuretok(azureKeyVault, "KeyVault")},
			// LoadBalancer
			"azurerm_lb":                      {Tok: azuretok(azureLB, "LoadBalancer")},
			"azurerm_lb_backend_address_pool": {Tok: azuretok(azureLB, "BackendAddressPool")},
			"azurerm_lb_nat_rule":             {Tok: azuretok(azureLB, "NatRule")},
			"azurerm_lb_nat_pool":             {Tok: azuretok(azureLB, "NatPool")},
			"azurerm_lb_probe":                {Tok: azuretok(azureLB, "Probe")},
			"azurerm_lb_rule":                 {Tok: azuretok(azureLB, "Rule")},
			// Managed Disks
			"azurerm_managed_disk": {Tok: azuretok(azureDisk, "ManagedDisk")},
			// Network
			"azurerm_local_network_gateway":  {Tok: azuretok(azureNetwork, "LocalNetworkGateway")},
			"azurerm_network_interface":      {Tok: azuretok(azureNetwork, "NetworkInterface")},
			"azurerm_network_security_group": {Tok: azuretok(azureNetwork, "NetworkSecurityGroup")},
			"azurerm_network_security_rule":  {Tok: azuretok(azureNetwork, "NetworkSecurityRule")},
			"azurerm_public_ip":              {Tok: azuretok(azureNetwork, "PublicIp")},
			"azurerm_route":                  {Tok: azuretok(azureNetwork, "Route")},
			"azurerm_route_table":            {Tok: azuretok(azureNetwork, "RouteTable")},
			"azurerm_subnet":                 {Tok: azuretok(azureNetwork, "Subnet")},
			// Redis
			"azurerm_redis_cache": {Tok: azuretok(azureRedis, "Cache")},
			// ResourceManager
			"azurerm_resource_group":      {Tok: azuretok(azureResources, "ResourceGroup")},
			"azurerm_template_deployment": {Tok: azuretok(azureResources, "TemplateDeployment")},
			// Search
			"azurerm_search_service": {Tok: azuretok(azureSearch, "Service")},
			// ServiceBus
			"azurerm_servicebus_namespace":    {Tok: azuretok(azureServiceBus, "Namespace")},
			"azurerm_servicebus_queue":        {Tok: azuretok(azureServiceBus, "Queue")},
			"azurerm_servicebus_subscription": {Tok: azuretok(azureServiceBus, "Subscription")},
			"azurerm_servicebus_topic":        {Tok: azuretok(azureServiceBus, "Topic")},
			// SQL
			"azurerm_sql_elasticpool":   {Tok: azuretok(azureSQL, "ElasticPool")},
			"azurerm_sql_database":      {Tok: azuretok(azureSQL, "Database")},
			"azurerm_sql_firewall_rule": {Tok: azuretok(azureSQL, "FirewallRule")},
			"azurerm_sql_server":        {Tok: azuretok(azureSQL, "SqlServer")},
			// Storage
			"azurerm_storage_account":   {Tok: azuretok(azureStorage, "Account")},
			"azurerm_storage_blob":      {Tok: azuretok(azureStorage, "Blob")},
			"azurerm_storage_container": {Tok: azuretok(azureStorage, "Container")},
			"azurerm_storage_share":     {Tok: azuretok(azureStorage, "Share")},
			"azurerm_storage_queue":     {Tok: azuretok(azureStorage, "Queue")},
			"azurerm_storage_table":     {Tok: azuretok(azureStorage, "Table")},
			// Traffic Manager
			"azurerm_traffic_manager_endpoint": {Tok: azuretok(azureTrafficManager, "Endpoint")},
			"azurerm_traffic_manager_profile":  {Tok: azuretok(azureTrafficManager, "Profile")},
			// Virtual Machines
			"azurerm_availability_set":          {Tok: azuretok(azureVirtualMachine, "AvailabilitySet")},
			"azurerm_virtual_machine_extension": {Tok: azuretok(azureVirtualMachine, "Extension")},
			"azurerm_virtual_machine":           {Tok: azuretok(azureVirtualMachine, "VirtualMachine")},
			"azurerm_virtual_machine_scale_set": {Tok: azuretok(azureVirtualMachine, "ScaleSet")},
			"azurerm_virtual_network":           {Tok: azuretok(azureVirtualMachine, "Network")},
			"azurerm_virtual_network_peering":   {Tok: azuretok(azureVirtualMachine, "NetworkPeering")},
		},
	}
}
