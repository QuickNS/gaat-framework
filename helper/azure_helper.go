package helper

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/msi/mgmt/msi"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/recoveryservices/mgmt/backup"
	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-04-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2019-05-01/frontdoor"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	kv "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	mysql "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-04-01/network"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2017-05-01-preview/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	sqlmi "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2014-04-01/sql"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	// SubscriptionIDEnvName azure sub name
	SubscriptionIDEnvName = "ARM_SUBSCRIPTION_ID"
)

/********************************
		Authorization
*********************************/

// NewAuthorizer will return Authorizer
// using credentials previously set with az login
func NewAuthorizer() (*autorest.Authorizer, error) {
	authorizer, err := auth.NewAuthorizerFromCLI()
	return &authorizer, err
}

/********************************
		Resource Groups
*********************************/

// GetResourceGroupE will return Group object and an error object
func GetResourceGroupE(resourceGroupName string) (*resources.Group, error) {

	client, err := GetGroupsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	resourceGroup, err := client.Get(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}
	return &resourceGroup, nil
}

// GetGroupsClientE creates a GroupsClient
func GetGroupsClientE(subscriptionID string) (*resources.GroupsClient, error) {
	client := resources.NewGroupsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Virtual Machines
*********************************/

// GetVirtualMachineE will return Group object and an error object
func GetVirtualMachineE(resourceGroupName string, virtualMachineName string) (*compute.VirtualMachine, error) {

	client, err := GetVirtualMachinesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	virtualMachine, err := client.Get(context.Background(), resourceGroupName, virtualMachineName, "")
	if err != nil {
		return nil, err
	}
	return &virtualMachine, nil
}

// GetVirtualMachinesClientE creates a VirtualMachinesClient
func GetVirtualMachinesClientE(subscriptionID string) (*compute.VirtualMachinesClient, error) {
	client := compute.NewVirtualMachinesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Disks
*********************************/

// GetDiskE will return Disk object and an error object
func GetDiskE(resourceGroupName string, diskName string) (*compute.Disk, error) {

	client, err := GetDisksClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	disk, err := client.Get(context.Background(), resourceGroupName, diskName)
	if err != nil {
		return nil, err
	}
	return &disk, nil
}

// GetDisksClientE creates a DisksClient
func GetDisksClientE(subscriptionID string) (*compute.DisksClient, error) {
	client := compute.NewDisksClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Network Interfaces
*********************************/

// GetInterfaceE gets NetworkInterface object
func GetInterfaceE(resourceGroupName, networkInterfaceName string) (*network.Interface, error) {

	client, err := GetInterfacesClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	nic, err := client.Get(context.Background(), resourceGroupName, networkInterfaceName, "")
	if err != nil {
		return nil, err
	}
	return &nic, nil
}

// GetInterfacesClient creates a virtual network client
func GetInterfacesClient(subscriptionID string) (*network.InterfacesClient, error) {
	client := network.NewInterfacesClient(subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Storage Accounts
*********************************/

// GetStorageAccountE gets storage.Account object
func GetStorageAccountE(resourceGroupName, storageAccountName string) (*storage.Account, error) {

	client, err := GetStorageAccountsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	account, err := client.GetProperties(context.Background(), resourceGroupName, storageAccountName, "")
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetStorageAccountKeysE gets NetworkInterface object
func GetStorageAccountKeysE(resourceGroupName, storageAccountName string) (*[]storage.AccountKey, error) {

	client, err := GetStorageAccountsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	keys, err := client.ListKeys(context.Background(), resourceGroupName, storageAccountName, "")
	if err != nil {
		return nil, err
	}
	return keys.Keys, nil
}

// GetStorageAccountsClientE creates a virtual network client
func GetStorageAccountsClientE(subscriptionID string) (*storage.AccountsClient, error) {
	client := storage.NewAccountsClient(subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Blob Containers
*********************************/

// ListBlobContainersForAccountE will return Group object and an error object
func ListBlobContainersForAccountE(resourceGroupName string, storageAccountName string) (*[]storage.ListContainerItem, error) {

	client, err := GetBlobContainersClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), resourceGroupName, storageAccountName, "", "", "")
	if err != nil {
		return nil, err
	}

	containerList := result.Values()
	return &containerList, nil
}

// GetBlobContainersClientE creates a GroupsClient
func GetBlobContainersClientE(subscriptionID string) (*storage.BlobContainersClient, error) {
	client := storage.NewBlobContainersClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		File Shares
*********************************/

// ListFileSharesForAccountE will return Group object and an error object
func ListFileSharesForAccountE(resourceGroupName string, storageAccountName string) (*[]storage.FileShareItem, error) {

	client, err := GetFileSharesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), resourceGroupName, storageAccountName, "", "", "")
	if err != nil {
		return nil, err
	}

	shareList := result.Values()
	return &shareList, nil
}

// GetFileSharesClientE creates a GroupsClient
func GetFileSharesClientE(subscriptionID string) (*storage.FileSharesClient, error) {
	client := storage.NewFileSharesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		File Services
*********************************/

// ListFileServicesForAccountE will return Group object and an error object
func ListFileServicesForAccountE(resourceGroupName string, storageAccountName string) (*storage.FileServiceItems, error) {

	client, err := GetFileServicesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	fileService, err := client.List(context.Background(), resourceGroupName, storageAccountName)
	if err != nil {
		return nil, err
	}
	return &fileService, nil
}

// GetFileServicesClientE creates a GroupsClient
func GetFileServicesClientE(subscriptionID string) (*storage.FileServicesClient, error) {
	client := storage.NewFileServicesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Blob Services
*********************************/

// ListBlobServicesForAccountE will return Group object and an error object
func ListBlobServicesForAccountE(resourceGroupName string, storageAccountName string) (*storage.BlobServiceItems, error) {

	client, err := GetBlobServicesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	fileService, err := client.List(context.Background(), resourceGroupName, storageAccountName)
	if err != nil {
		return nil, err
	}
	return &fileService, nil
}

// GetBlobServicesClientE creates a BlobServicesClient
func GetBlobServicesClientE(subscriptionID string) (*storage.BlobServicesClient, error) {
	client := storage.NewBlobServicesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Virtual Network
*********************************/

// GetVirtualNetworkE gets virtual network object
func GetVirtualNetworkE(resourceGroupName, virtualNetworkName string) (*network.VirtualNetwork, error) {

	client, err := GetVirtualNetworksClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	virtualNetwork, err := client.Get(context.Background(), resourceGroupName, virtualNetworkName, "")
	if err != nil {
		return nil, err
	}
	return &virtualNetwork, nil
}

// GetVirtualNetworksClient creates a virtual network client
func GetVirtualNetworksClient(subscriptionID string) (*network.VirtualNetworksClient, error) {
	vnClient := network.NewVirtualNetworksClient(subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}

	vnClient.Authorizer = *authorizer
	return &vnClient, nil
}

/********************************
		Virtual Network Peerings
*********************************/

// GetVirtualNetworkPeeringE gets a VirtualNetworkPeering object
func GetVirtualNetworkPeeringE(resourceGroupName, virtualNetworkName string, virtualNetworkPeeringName string) (*network.VirtualNetworkPeering, error) {

	client, err := GetVirtualNetworkPeeringsClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, virtualNetworkName, virtualNetworkPeeringName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListVirtualNetworkPeeringE gets an array of VirtualNetworkPeering objects
func ListVirtualNetworkPeeringE(resourceGroupName, virtualNetworkName string) ([]network.VirtualNetworkPeering, error) {

	client, err := GetVirtualNetworkPeeringsClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), resourceGroupName, virtualNetworkName)
	if err != nil {
		return nil, err
	}
	return result.Values(), nil
}

// GetVirtualNetworkPeeringsClient creates a VirtualNetworkPeeringsClient
func GetVirtualNetworkPeeringsClient(subscriptionID string) (*network.VirtualNetworkPeeringsClient, error) {
	vnClient := network.NewVirtualNetworkPeeringsClient(subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}

	vnClient.Authorizer = *authorizer
	return &vnClient, nil
}

/********************************
		Subnets
*********************************/

// GetSubnetE gets virtual network object
func GetSubnetE(resourceGroupName, virtualNetworkName string, subnetName string) (*network.Subnet, error) {

	client, err := GetSubNetClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	subnet, err := client.Get(context.Background(), resourceGroupName, virtualNetworkName, subnetName, "")
	if err != nil {
		return nil, err
	}
	return &subnet, nil
}

//GetSubnetAddressesForVirtualNetworkE gets all virtual network subclients name, and address prefix
func GetSubnetAddressesForVirtualNetworkE(resourceGroupName, virtualNetworkName string) (*map[string]string, error) {
	client, err := GetSubNetClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}

	subnets, err := client.List(context.Background(), resourceGroupName, virtualNetworkName)
	if err != nil {
		return nil, err
	}

	subNetDetails := make(map[string]string)
	for _, v := range subnets.Values() {
		subnetName := v.Name
		subNetAddressPrefix := v.AddressPrefix
		subNetDetails[*subnetName] = *subNetAddressPrefix
	}
	return &subNetDetails, nil
}

//GetSubnetSecurityGroupsForVirtualNetworkE gets all virtual network subclients name, and security group IDs
func GetSubnetSecurityGroupsForVirtualNetworkE(resourceGroupName, virtualNetworkName string) (*map[string]string, error) {
	client, err := GetSubNetClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}

	subnets, err := client.List(context.Background(), resourceGroupName, virtualNetworkName)
	if err != nil {
		return nil, err
	}

	subNetDetails := make(map[string]string)
	for _, v := range subnets.Values() {
		subnetName := v.Name
		if v.NetworkSecurityGroup != nil {
			subNetDetails[*subnetName] = *v.NetworkSecurityGroup.ID
		} else {
			subNetDetails[*subnetName] = ""
		}

	}
	return &subNetDetails, nil
}

// GetSubNetClient creates a virtual network subnet client
func GetSubNetClient(subscriptionID string) (*network.SubnetsClient, error) {
	subNetClient := network.NewSubnetsClient(subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}

	subNetClient.Authorizer = *authorizer
	return &subNetClient, nil
}

/********************************
		Key Vault Management methods
*********************************/

// GetKeyVaultE will return a Vault object and an error object
func GetKeyVaultE(resourceGroupName, keyVaultName string) (*kv.Vault, error) {
	client, err := GetKeyVaultManagementClientE(os.Getenv(SubscriptionIDEnvName))

	if err != nil {
		return nil, err
	}
	keyVault, err := client.Get(context.Background(), resourceGroupName, keyVaultName)
	if err != nil {
		return nil, err
	}
	return &keyVault, nil
}

// GetKeyVaultManagementClientE creates a VaultsClient client
func GetKeyVaultManagementClientE(subscriptionID string) (*kv.VaultsClient, error) {
	client := kv.NewVaultsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Key Vault Client methods
*********************************/

// NewKeyVaultAuthorizer witll return Authorizer for KeyVault
func NewKeyVaultAuthorizer() (*autorest.Authorizer, error) {
	authorizer, err := kvauth.NewAuthorizerFromCLI()
	return &authorizer, err
}

// GetKeyVaultClientE creates a KeyVault client
func GetKeyVaultClientE() (*keyvault.BaseClient, error) {
	kvClient := keyvault.New()
	authorizer, err := NewKeyVaultAuthorizer()

	if err != nil {
		return nil, err
	}

	kvClient.Authorizer = *authorizer
	return &kvClient, nil
}

// GetKeyVaultSecretCurrentVersion gets the current version of the KeyVault
// e.g. https://foo.vault.azure.net/secrets/BAR/194bd7da9aa54944ab316faebd9120d0 -> 194bd7da9aa54944ab316faebd9120d0
func GetKeyVaultSecretCurrentVersion(keyVaultName, secretName string) (string, error) {
	client, err := GetKeyVaultClientE()
	if err != nil {
		return "", err
	}
	var maxVersionsCount int32 = 25
	versions, err := client.GetSecretVersions(context.Background(),
		fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName),
		secretName,
		&maxVersionsCount)
	if err != nil {
		return "", err
	}
	items := versions.Values()
	sort.Slice(items, func(i, j int) bool {
		return (*items[i].Attributes.Updated).Duration().Milliseconds() > (*items[j].Attributes.Updated).Duration().Milliseconds()
	})
	nonVersion := fmt.Sprintf("https://%s.vault.azure.net/secrets/%s/", keyVaultName, secretName)
	return strings.Replace(*items[0].ID, nonVersion, "", 1), nil
}

// GetKeyVaultSecretWithVersion is get secret from the specific key vault.
func GetKeyVaultSecretWithVersion(keyVaultName, secretName, version string) (string, error) {
	client, err := GetKeyVaultClientE()
	if err != nil {
		return "", err
	}
	secret, err := client.GetSecret(context.Background(), fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName), secretName, version)
	if err != nil {
		return "", err
	}
	return *secret.Value, err
}

// GetKeyVaultSecret returns current secret
func GetKeyVaultSecret(keyVaultName, secretName string) (string, error) {
	version, err := GetKeyVaultSecretCurrentVersion(keyVaultName, secretName)
	if err != nil {
		return "", err
	}
	secret, err := GetKeyVaultSecretWithVersion(keyVaultName, secretName, version)
	if err != nil {
		return "", err
	}
	return secret, nil
}

/*****************************************
		Private Endpoints
******************************************/

// GetPrivateEndpointE will return a PrivateEndpoint object and an error object
func GetPrivateEndpointE(resourceGroupName, endpointName string) (*network.PrivateEndpoint, error) {
	client, err := GetPrivateEndpointsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	endpoint, err := client.Get(context.Background(), resourceGroupName, endpointName, "")
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// GetPrivateEndpointsClientE creates a PrivateEndpointsClient client
func GetPrivateEndpointsClientE(subscriptionID string) (*network.PrivateEndpointsClient, error) {
	client := network.NewPrivateEndpointsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Private DNS Zone Groups
******************************************/

// GetPrivateDNSZoneGroupE will return a PrivateDNSZoneGroup object and an error object
func GetPrivateDNSZoneGroupE(resourceGroupName, endpointName string, groupName string) (*network.PrivateDNSZoneGroup, error) {
	client, err := GetPrivateDNSZoneGroupsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, endpointName, groupName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListPrivateDNSZoneGroupsE will return an array of PrivateDNSZoneGroup object and an error object
func ListPrivateDNSZoneGroupsE(resourceGroupName, endpointName string) (*[]network.PrivateDNSZoneGroup, error) {
	client, err := GetPrivateDNSZoneGroupsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), endpointName, resourceGroupName)
	if err != nil {
		return nil, err
	}
	groups := result.Values()
	return &groups, nil
}

// GetPrivateDNSZoneGroupsClientE creates a PrivateDNSZoneGroupsClient client
func GetPrivateDNSZoneGroupsClientE(subscriptionID string) (*network.PrivateDNSZoneGroupsClient, error) {
	client := network.NewPrivateDNSZoneGroupsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Private DNS Zones
******************************************/

// GetPrivateDNSZoneE will return a PrivateDNSZoneGroup object and an error object
func GetPrivateDNSZoneE(resourceGroupName string, dnsZoneName string) (*privatedns.PrivateZone, error) {
	client, err := GetPrivateDNSZonesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	zone, err := client.Get(context.Background(), resourceGroupName, dnsZoneName)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetPrivateDNSZonesClientE creates a PrivateZonesClient
func GetPrivateDNSZonesClientE(subscriptionID string) (*privatedns.PrivateZonesClient, error) {
	client := privatedns.NewPrivateZonesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Record Sets
******************************************/

// ListRecordSetsE will return a PrivateDNSZoneGroup object and an error object
func ListRecordSetsE(resourceGroupName string, dnsZoneName string, numberToRetrieve int32) (*[]privatedns.RecordSet, error) {
	client, err := GetRecordSetsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), resourceGroupName, dnsZoneName, &numberToRetrieve, "")
	if err != nil {
		return nil, err
	}
	recordSets := make([]privatedns.RecordSet, 0, numberToRetrieve)
	for _, n := range result.Values() {
		recordSets = append(recordSets, n)
	}
	return &recordSets, nil
}

// GetRecordSetsClientE creates a PrivateZonesClient
func GetRecordSetsClientE(subscriptionID string) (*privatedns.RecordSetsClient, error) {
	client := privatedns.NewRecordSetsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Virtual Network Link
******************************************/

// ListVirtualNetworkLinkE will return an array of VirtualNetworkLink objects and an error object
func ListVirtualNetworkLinkE(resourceGroupName string, dnsZoneName string, numberToRetrieve int32) (*[]privatedns.VirtualNetworkLink, error) {
	client, err := GetVirtualNetworkLinksClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}

	result, err := client.List(context.Background(), resourceGroupName, dnsZoneName, &numberToRetrieve)
	if err != nil {
		return nil, err
	}

	networkLinks := make([]privatedns.VirtualNetworkLink, 0, numberToRetrieve)
	for _, n := range result.Values() {
		networkLinks = append(networkLinks, n)
	}
	return &networkLinks, nil
}

// GetVirtualNetworkLinkE will return a privatedns.VirtualNetworkLink object and an error object
func GetVirtualNetworkLinkE(resourceGroupName string, dnsZoneName string, name string) (*privatedns.VirtualNetworkLink, error) {
	client, err := GetVirtualNetworkLinksClient(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}

	result, err := client.Get(context.Background(), resourceGroupName, dnsZoneName, name)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetVirtualNetworkLinksClient creates a PrivateZonesClient
func GetVirtualNetworkLinksClient(subscriptionID string) (*privatedns.VirtualNetworkLinksClient, error) {
	client := privatedns.NewVirtualNetworkLinksClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Route Tables
******************************************/

// GetRouteTableE will return a RouteTable object and an error object
func GetRouteTableE(resourceGroupName, routeTableName string) (*network.RouteTable, error) {
	client, err := GetRouteTableClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	table, err := client.Get(context.Background(), resourceGroupName, routeTableName, "")
	if err != nil {
		return nil, err
	}
	return &table, nil
}

// GetRouteTableClientE creates a RouteTablesClient
func GetRouteTableClientE(subscriptionID string) (*network.RouteTablesClient, error) {
	client := network.NewRouteTablesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Network Security Groups (NSG)
******************************************/

// GetSecurityGroupE will return a SecurityGroup object and an error object
func GetSecurityGroupE(resourceGroupName, securityGroupName string) (*network.SecurityGroup, error) {
	client, err := GetSecurityGroupsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	securityGroup, err := client.Get(context.Background(), resourceGroupName, securityGroupName, "")
	if err != nil {
		return nil, err
	}
	return &securityGroup, nil
}

// GetSecurityGroupsClientE creates a SecurityGroupClient client
func GetSecurityGroupsClientE(subscriptionID string) (*network.SecurityGroupsClient, error) {
	client := network.NewSecurityGroupsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/*****************************************
		Container Registry
******************************************/

// GetRegistryE will return a Registry object and an error object
func GetRegistryE(resourceGroupName, containerRegistry string) (*containerregistry.Registry, error) {
	client, err := GetRegistriesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	securityGroup, err := client.Get(context.Background(), resourceGroupName, containerRegistry)
	if err != nil {
		return nil, err
	}
	return &securityGroup, nil
}

// GetRegistryPoliciesE will return a RegistryPolicies object and an error object
func GetRegistryPoliciesE(resourceGroupName, containerRegistry string) (*containerregistry.RegistryPolicies, error) {
	client, err := GetRegistriesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	securityGroup, err := client.ListPolicies(context.Background(), resourceGroupName, containerRegistry)
	if err != nil {
		return nil, err
	}
	return &securityGroup, nil
}

// GetRegistriesClientE creates a RegistriesClient client
func GetRegistriesClientE(subscriptionID string) (*containerregistry.RegistriesClient, error) {
	client := containerregistry.NewRegistriesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Bastion
*********************************/

// GetBastionHostE will return BastionHost object and an error object
func GetBastionHostE(resourceGroupName string, bastionHostName string) (*network.BastionHost, error) {
	client, err := GetBastionHostClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	bastionHost, err := client.Get(context.Background(), resourceGroupName, bastionHostName)
	if err != nil {
		return nil, err
	}
	return &bastionHost, nil
}

// GetBastionHostClientE creates a BastionHostsClient client
func GetBastionHostClientE(subscriptionID string) (*network.BastionHostsClient, error) {
	client := network.NewBastionHostsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Application Gateway
*********************************/

// GetApplicationGatewayE will return ApplicationGateway object and an error object
func GetApplicationGatewayE(resourceGroupName, applicationGatewayName string) (*network.ApplicationGateway, error) {
	client, err := GetApplicationGatewayClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	applicationGateway, err := client.Get(context.Background(), resourceGroupName, applicationGatewayName)
	if err != nil {
		return nil, err
	}
	return &applicationGateway, nil
}

// GetApplicationGatewayClientE creates a ApplicationGatewaysClient client
func GetApplicationGatewayClientE(subscriptionID string) (*network.ApplicationGatewaysClient, error) {
	client := network.NewApplicationGatewaysClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		PublicIP
*********************************/

// GetPublicIPAddressE will return PublicIPAddress object and an error object
func GetPublicIPAddressE(resourceGroupName, publicIPAddressName string) (*network.PublicIPAddress, error) {
	client, err := GetPublicIPAddressClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	publicIPAddress, err := client.Get(context.Background(), resourceGroupName, publicIPAddressName, "")
	if err != nil {
		return nil, err
	}
	return &publicIPAddress, nil
}

// GetPublicIPAddressClientE creates a PublicIPAddresses client
func GetPublicIPAddressClientE(subscriptionID string) (*network.PublicIPAddressesClient, error) {
	client := network.NewPublicIPAddressesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		UserAssignedIdentity
*********************************/

// GetUserAssignedIdentityE will return Identity object and an error object
func GetUserAssignedIdentityE(resourceGroupName, identityName string) (*msi.Identity, error) {
	client, err := GetUserAssignedIdentitiesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	identity, err := client.Get(context.Background(), resourceGroupName, identityName)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

// GetUserAssignedIdentitiesClientE creates a UserAssignedIdentitiesClient
func GetUserAssignedIdentitiesClientE(subscriptionID string) (*msi.UserAssignedIdentitiesClient, error) {
	client := msi.NewUserAssignedIdentitiesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Role Assignments
*********************************/

// ListRoleAssignmentsForPrincipalID will return Identity object and an error object
func ListRoleAssignmentsForPrincipalID(resourceGroupName string, principalID string) (*[]string, error) {
	client, err := GetRoleAssignmentsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("principalId eq '{%s}'", principalID)
	result, err := client.ListForResourceGroup(context.Background(), resourceGroupName, filter)
	if err != nil {
		return nil, err
	}

	assignments := make([]string, 0, len(result.Values()))
	for _, v := range result.Values() {
		role := strings.ToLower(*v.Properties.RoleDefinitionID)
		assignments = append(assignments, role)
	}
	return &assignments, nil
}

// GetRoleAssignmentsClientE creates a RoleAssignmentsClient
func GetRoleAssignmentsClientE(subscriptionID string) (*authorization.RoleAssignmentsClient, error) {
	client := authorization.NewRoleAssignmentsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Role Definitions
*********************************/

// GetRoleDefinitionE will return RoleDefinition object and an error object
func GetRoleDefinitionE(roleDefinitionID string) (*authorization.RoleDefinition, error) {
	client, err := GetRoleDefinitionsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	identity, err := client.GetByID(context.Background(), roleDefinitionID)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

// GetRoleDefinitionsClientE creates a RoleDefinitionsClient
func GetRoleDefinitionsClientE(subscriptionID string) (*authorization.RoleDefinitionsClient, error) {
	client := authorization.NewRoleDefinitionsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		AKS (Managed Cluster)
*********************************/

// GetManagedClusterE will return ContainerService object and an error object
func GetManagedClusterE(resourceGroupName, clusterName string) (*containerservice.ManagedCluster, error) {
	client, err := GetManagedClustersClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	managedCluster, err := client.Get(context.Background(), resourceGroupName, clusterName)
	if err != nil {
		return nil, err
	}
	return &managedCluster, nil
}

// GetManagedClustersClientE creates a ContainerServicesClient client
func GetManagedClustersClientE(subscriptionID string) (*containerservice.ManagedClustersClient, error) {
	client := containerservice.NewManagedClustersClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

// GetClusterAdminCredentialsE returns credential information includes kubeconfig
func GetClusterAdminCredentialsE(resourceGroupName, clusterName string) (*containerservice.CredentialResults, error) {
	client, err := GetManagedClustersClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}

	credentials, err := client.ListClusterAdminCredentials(context.Background(), resourceGroupName, clusterName)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

// WriteKubeconfigFromCredentialsE writes a kubeconfig file
func WriteKubeconfigFromCredentialsE(credentialResults *containerservice.CredentialResults, filePath string) error {
	kubeconfig := (*(*credentialResults.Kubeconfigs)[0].Value)
	return ioutil.WriteFile(filePath, kubeconfig, os.ModePerm)
}

/********************************
		Availability Sets
*********************************/

// GetAvailabilitySetE will return AvailabilitySet object and an error object
func GetAvailabilitySetE(resourceGroupName string, availabilitySetName string) (*compute.AvailabilitySet, error) {
	client, err := GetAvailabilitySetsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, availabilitySetName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAvailabilitySetsClientE creates a AvailabilitySetsClient
func GetAvailabilitySetsClientE(subscriptionID string) (*compute.AvailabilitySetsClient, error) {
	client := compute.NewAvailabilitySetsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		SQL Managed Instance
*********************************/

// GetManagedInstanceE will return ManagedInstance object and an error object
func GetManagedInstanceE(resourceGroupName string, managedInstanceName string) (*sqlmi.ManagedInstance, error) {
	client, err := GetManagedInstancesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	instance, err := client.Get(context.Background(), resourceGroupName, managedInstanceName)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// GetManagedInstancesClientE creates a ManagedInstancesClient
func GetManagedInstancesClientE(subscriptionID string) (*sqlmi.ManagedInstancesClient, error) {
	client := sqlmi.NewManagedInstancesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Log Analytics Workspace
*********************************/

// ListLogAnalyticsWorkspacesByResourceGroupE will return a map[string]string with Workspace IDs and Names and an error object
func ListLogAnalyticsWorkspacesByResourceGroupE(resourceGroupName string) (map[string]string, error) {
	client, err := GetLogAnalyticsWorkspacesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListByResourceGroup(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}
	workspaces := make(map[string]string)
	for _, w := range *result.Value {
		workspaces[strings.ToLower(*w.ID)] = *w.Name
	}
	return workspaces, nil
}

// GetLogAnalyticsWorkspaceE will return Workspace object and an error object
func GetLogAnalyticsWorkspaceE(resourceGroupName string, workspaceName string) (*operationalinsights.Workspace, error) {
	client, err := GetLogAnalyticsWorkspacesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	workspace, err := client.Get(context.Background(), resourceGroupName, workspaceName)
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

// GetLogAnalyticsWorkspacesClientE creates a WorkspacesClient
func GetLogAnalyticsWorkspacesClientE(subscriptionID string) (*operationalinsights.WorkspacesClient, error) {
	client := operationalinsights.NewWorkspacesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Solutions (OperationsManagement)
*********************************/

// GetSolutionE will return Solution object and an error object
func GetSolutionE(resourceGroupName string, solutionName string) (*operationsmanagement.Solution, error) {
	client, err := GetSolutionsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	solution, err := client.Get(context.Background(), resourceGroupName, solutionName)
	if err != nil {
		return nil, err
	}
	return &solution, nil
}

// GetSolutionsClientE creates a SolutionsClient
func GetSolutionsClientE(subscriptionID string) (*operationsmanagement.SolutionsClient, error) {
	client := operationsmanagement.NewSolutionsClient(subscriptionID, "", "", "")
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Azure Firewall
*********************************/

// GetFirewallE will return Firewall object and an error object
func GetFirewallE(resourceGroupName string, firewallName string) (*network.AzureFirewall, error) {
	client, err := GetAzureFirewallsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	firewall, err := client.Get(context.Background(), resourceGroupName, firewallName)
	if err != nil {
		return nil, err
	}
	return &firewall, nil
}

// GetAzureFirewallsClientE creates a FirewallsClient
func GetAzureFirewallsClientE(subscriptionID string) (*network.AzureFirewallsClient, error) {
	client := network.NewAzureFirewallsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		DDOS Plan
*********************************/

// GetDdosProtectionPlanE will return DdosProtectionPlan object and an error object
func GetDdosProtectionPlanE(resourceGroupName string, planName string) (*network.DdosProtectionPlan, error) {
	client, err := GetDdosProtectionPlansClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	ddos, err := client.Get(context.Background(), resourceGroupName, planName)
	if err != nil {
		return nil, err
	}
	return &ddos, nil
}

// GetDdosProtectionPlansClientE creates a DdosProtectionPlansClient
func GetDdosProtectionPlansClientE(subscriptionID string) (*network.DdosProtectionPlansClient, error) {
	client := network.NewDdosProtectionPlansClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Recovery Services Vault
*********************************/

// GetRecoveryServicesVaultE will return Vault object and an error object
func GetRecoveryServicesVaultE(resourceGroupName string, vaultName string) (*recoveryservices.Vault, error) {
	client, err := GetVaultsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, vaultName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVaultsClientE creates a VaultsClient
func GetVaultsClientE(subscriptionID string) (*recoveryservices.VaultsClient, error) {
	client := recoveryservices.NewVaultsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Backup Policies
*********************************/

// ListBackupPoliciesE will return list of backup policies associated with Recovery Services Vault.
func ListBackupPoliciesE(resourceGroupName string, vaultName string) ([]backup.ProtectionPolicyResource, error) {
	client, err := GetPoliciesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), vaultName, resourceGroupName, "")
	if err != nil {
		return nil, err
	}

	return result.Values(), nil
}

// GetPoliciesClientE creates a PoliciesClient
func GetPoliciesClientE(subscriptionID string) (*backup.PoliciesClient, error) {
	client := backup.NewPoliciesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Protected Items
*********************************/

// GetProtectedItemE will return ProtectedItemsResource object and an error object
func GetProtectedItemE(resourceGroupName string, vaultName string, fabricName string, containerName string, protectedItemName string) (*backup.ProtectedItemResource, error) {
	client, err := GetProtectedItemsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), vaultName, resourceGroupName, fabricName, containerName, protectedItemName, "")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProtectedItemsClientE creates a ProtectedItemsClient
func GetProtectedItemsClientE(subscriptionID string) (*backup.ProtectedItemsClient, error) {
	client := backup.NewProtectedItemsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Redis Cache
*********************************/

// GetRedisE will return redis.ResourceStype object and an error object
func GetRedisE(resourceGroupName string, redisName string) (*redis.ResourceType, error) {
	client, err := GetRedisClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, redisName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRedisAccessKeysE will return redis.AccessKeys object and an error object
func GetRedisAccessKeysE(resourceGroupName string, redisName string) (*redis.AccessKeys, error) {
	client, err := GetRedisClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListKeys(context.Background(), resourceGroupName, redisName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRedisClientE creates a redis.Client object
func GetRedisClientE(subscriptionID string) (*redis.Client, error) {
	client := redis.NewClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Virtual Network Gateway
*********************************/

// GetVirtualNetworkGatewayE will return redis.VirtualNetworkGateway object and an error object
func GetVirtualNetworkGatewayE(resourceGroupName string, virtualNetworkGatewayName string) (*network.VirtualNetworkGateway, error) {
	client, err := GetVirtualNetworkGatewaysClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, virtualNetworkGatewayName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVirtualNetworkGatewaysClientE creates a network.VirtualNetworksClient object
func GetVirtualNetworkGatewaysClientE(subscriptionID string) (*network.VirtualNetworkGatewaysClient, error) {
	client := network.NewVirtualNetworkGatewaysClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Diagnostic Settings
*********************************/

// GetDiagnosticSettingsE will return insights.DiagnosticSettings object and an error object
func GetDiagnosticSettingsE(resourceURI string, name string) (*insights.DiagnosticSettings, error) {
	client, err := GetDiagnosticSettingsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceURI, name)
	if err != nil {
		return nil, err
	}
	return result.DiagnosticSettings, nil
}

// ListDiagnosticSettingsE will return a list of insights.DiagnosticSettingsResource object and an error object
func ListDiagnosticSettingsE(resourceURI string) (*[]insights.DiagnosticSettingsResource, error) {
	client, err := GetDiagnosticSettingsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.List(context.Background(), resourceURI)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}

// GetDiagnosticSettingsClientE creates a insights.DiagnosticSettingsClient  object
func GetDiagnosticSettingsClientE(subscriptionID string) (*insights.DiagnosticSettingsClient, error) {
	client := insights.NewDiagnosticSettingsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		MySql Server
*********************************/

// GetMySQLServerE will return insights.DiagnosticSettings object and an error object
func GetMySQLServerE(resourceGroupName string, serverName string) (*mysql.Server, error) {
	client, err := GetMySQLServersClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, serverName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMySQLServersClientE creates a mysql.ServersClient  object
func GetMySQLServersClientE(subscriptionID string) (*mysql.ServersClient, error) {
	client := mysql.NewServersClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		MySql Databases
*********************************/

// GetMySQLDatabaseE will return insights.DiagnosticSettings object and an error object
func GetMySQLDatabaseE(resourceGroupName string, serverName string, databaseName string) (*mysql.Database, error) {
	client, err := GetMySQLDatabasesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, serverName, databaseName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMySQLDatabasesClientE creates a mysql.DatabasesClient  object
func GetMySQLDatabasesClientE(subscriptionID string) (*mysql.DatabasesClient, error) {
	client := mysql.NewDatabasesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		MySql Server Configurations
*********************************/

// ListMySQLServerConfigE will return mysql.ConfigurationListResult object and an error object
func ListMySQLServerConfigE(resourceGroupName string, serverName string) (*[]mysql.Configuration, error) {
	client, err := GetMySQLConfigsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListByServer(context.Background(), resourceGroupName, serverName)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}

// GetMySQLConfigsClientE creates a mysql.ConfigurationsClient   object
func GetMySQLConfigsClientE(subscriptionID string) (*mysql.ConfigurationsClient, error) {
	client := mysql.NewConfigurationsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		MySql Server Configurations
*********************************/

// ListMySQLVirtualNetworkRulesE will return mysql.ConfigurationListResult object and an error object
func ListMySQLVirtualNetworkRulesE(resourceGroupName string, serverName string) (*[]mysql.VirtualNetworkRule, error) {
	client, err := GetMySQLVirtualNetworkRulesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListByServer(context.Background(), resourceGroupName, serverName)
	if err != nil {
		return nil, err
	}
	rules := result.Values()
	return &rules, nil
}

// GetMySQLVirtualNetworkRulesClientE creates a mysql.ConfigurationsClient   object
func GetMySQLVirtualNetworkRulesClientE(subscriptionID string) (*mysql.VirtualNetworkRulesClient, error) {
	client := mysql.NewVirtualNetworkRulesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Cosmos Database Account
*********************************/

// GetCosmosDatabaseAccountE will return documentdb.DatabaseAccountGetResults object and an error object
func GetCosmosDatabaseAccountE(resourceGroupName string, accountName string) (*documentdb.DatabaseAccountGetResults, error) {
	client, err := GetDatabaseAccountsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, accountName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCosmosKeysE will return documentdb.DatabaseAccountGetResults object and an error object
func GetCosmosKeysE(resourceGroupName string, accountName string) (*documentdb.DatabaseAccountListKeysResult, error) {
	client, err := GetDatabaseAccountsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListKeys(context.Background(), resourceGroupName, accountName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDatabaseAccountsClientE creates a documentdb.DatabaseAccountsClient object
func GetDatabaseAccountsClientE(subscriptionID string) (*documentdb.DatabaseAccountsClient, error) {
	client := documentdb.NewDatabaseAccountsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Cassandra Resources
*********************************/

// GetCassandraKeySpaceE will return documentdb.CassandraKeyspaceGetResults object and an error object
func GetCassandraKeySpaceE(resourceGroupName string, accountName string, keySpaceName string) (*documentdb.CassandraKeyspaceGetResults, error) {
	client, err := GetCassandraResourcesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.GetCassandraKeyspace(context.Background(), resourceGroupName, accountName, keySpaceName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCassandraResourcesClientE creates a documentdb.CassandraResourcesClient  object
func GetCassandraResourcesClientE(subscriptionID string) (*documentdb.CassandraResourcesClient, error) {
	client := documentdb.NewCassandraResourcesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Event Hub Namespace
*********************************/

// GetEventHubNamespaceE will return documentdb.EHNamespace object and an error object
func GetEventHubNamespaceE(resourceGroupName string, namespaceName string) (*eventhub.EHNamespace, error) {
	client, err := GetEventHubNamespacesClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, namespaceName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEventHubNamespacesClientE creates a mysql.ServersClient  object
func GetEventHubNamespacesClientE(subscriptionID string) (*eventhub.NamespacesClient, error) {
	client := eventhub.NewNamespacesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Event Hub
*********************************/

// GetEventHubE will return documentdb.EHNamespace object and an error object
func GetEventHubE(resourceGroupName string, namespaceName string, eventHubName string) (*eventhub.Model, error) {
	client, err := GetEventHubsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, namespaceName, eventHubName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEventHubsClientE creates a mysql.ServersClient  object
func GetEventHubsClientE(subscriptionID string) (*eventhub.EventHubsClient, error) {
	client := eventhub.NewEventHubsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		App Service Plans
*********************************/

// GetAppServicePlanE will return web.AppServicePlan object and an error object
func GetAppServicePlanE(resourceGroupName string, name string) (*web.AppServicePlan, error) {
	client, err := GetAppServicePlansClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAppServicePlansClientE creates a mysql.ServersClient  object
func GetAppServicePlansClientE(subscriptionID string) (*web.AppServicePlansClient, error) {
	client := web.NewAppServicePlansClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
		Web Apps
*********************************/

// GetSiteE will return web.Site object and an error object
func GetSiteE(resourceGroupName string, name string) (*web.Site, error) {
	client, err := GetAppsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListApplicationSettingsE will return web.StringDictionary
func ListApplicationSettingsE(resourceGroupName string, name string) (*web.StringDictionary, error) {
	client, err := GetAppsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListApplicationSettings(context.Background(), resourceGroupName, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListFunctionSettingsE will return an array of web.SiteConfigResource
func ListSiteConfigurationsE(resourceGroupName string, name string) ([]web.SiteConfigResource, error) {
	client, err := GetAppsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.ListConfigurations(context.Background(), resourceGroupName, name)
	if err != nil {
		return nil, err
	}

	return result.Values(), nil
}

// GetFunctionE will return web.AppServicePlan object and an error object
func GetFunctionE(resourceGroupName string, name string, functionName string) (*web.FunctionEnvelope, error) {
	client, err := GetAppsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.GetFunction(context.Background(), resourceGroupName, name, functionName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//GetSwiftVirtualNetworkConnectionE will return web.SwiftVirtualNetwork object and an error object
func GetSwiftVirtualNetworkConnectionE(resourceGroupName string, name string) (*web.SwiftVirtualNetwork, error) {
	client, err := GetAppsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.GetSwiftVirtualNetworkConnection(context.Background(), resourceGroupName, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAppsClientE creates a mysql.ServersClient  object
func GetAppsClientE(subscriptionID string) (*web.AppsClient, error) {
	client := web.NewAppsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/********************************
	SQL Server
*********************************/

// GetSQLServerE will return sql.Server object and an error object
func GetSQLServerE(resourceGroupName string, serverName string) (*sql.Server, error) {
	client, err := GetSQLServersClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, serverName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSQLServersClientE creates a sql.ServersClient
func GetSQLServersClientE(subscriptionID string) (*sql.ServersClient, error) {
	client := sql.NewServersClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/************************************
	Virtual Machine Scale Set (VMSS)
*************************************/

// GetVMScaleSetE will return compute.VirtualMachiuneScaleSet object and an error object
func GetVMScaleSetE(resourceGroupName string, vmScaleSetName string) (*compute.VirtualMachineScaleSet, error) {
	client, err := GetVMScaleSetsClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, vmScaleSetName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVMScaleSetsClientE creates a sql.ServersClient
func GetVMScaleSetsClientE(subscriptionID string) (*compute.VirtualMachineScaleSetsClient, error) {
	client := compute.NewVirtualMachineScaleSetsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}

/************************************
	frontdoor.FrontDoor
*************************************/

// GetFrontDoorE will return  frontdoor.FrontDoor object and an error object
func GetFrontDoorE(resourceGroupName string, frontDoorName string) (*frontdoor.FrontDoor, error) {
	client, err := GetFrontDoorClientE(os.Getenv(SubscriptionIDEnvName))
	if err != nil {
		return nil, err
	}
	result, err := client.Get(context.Background(), resourceGroupName, frontDoorName)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFrontDoorClientE creates a frontdoor.FrontDoorsClient
func GetFrontDoorClientE(subscriptionID string) (*frontdoor.FrontDoorsClient, error) {
	client := frontdoor.NewFrontDoorsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}
