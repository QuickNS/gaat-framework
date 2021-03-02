package helper

import (
	"fmt"
	"os"
)

const (
	// Reader role ID
	Reader = "acdd72a7-3385-48ef-bd42-f606fba81ae7"
	// Contributor role ID
	Contributor = "b24988ac-6180-42a0-ab88-20f7382dd24c"
	// ManagedIdentityOperator role ID
	ManagedIdentityOperator = "f1a07417-d97a-45cb-824c-7a7467783830"
	// AcrPull role ID
	AcrPull = "7f951dda-4ed3-4680-a7ca-43fe172d538d"
	// VirtualMachineContributor role ID
	VirtualMachineContributor = "9980e02c-c2be-4d73-94e8-173b1dc7cf3c"
	// NetworkContributor role ID
	NetworkContributor = "4d97b98b-1d4f-4787-a291-c67834d212e7"
)

// GetFullyQualifiedRoleDefinitionID returns a FQDN Role Definition ID
func GetFullyQualifiedRoleDefinitionID(roleID string) string {
	subscription := os.Getenv("ARM_SUBSCRIPTION_ID")
	return fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", subscription, roleID)
}
