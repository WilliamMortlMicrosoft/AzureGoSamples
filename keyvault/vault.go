// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package keyvault

import (
	"context"
	"fmt"
	//"encoding/json"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"

	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/config"
	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/iam"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	uuid "github.com/satori/go.uuid"
)

func getVaultsClient() keyvault.VaultsClient {
	vaultsClient := keyvault.NewVaultsClient(config.SubscriptionID())
	a, _ := iam.GetResourceManagementAuthorizer()
	vaultsClient.Authorizer = a
	vaultsClient.AddToUserAgent(config.UserAgent())
	return vaultsClient
}

// CreateVault creates a new vault
func CreateVault(ctx context.Context, vaultName string) (keyvault.Vault, error) {
	vaultsClient := getVaultsClient()
	tenantID, err := uuid.FromString(config.TenantID())
	if err != nil {
		return keyvault.Vault{}, err
	}

	_, errRet := vaultsClient.CreateOrUpdate(
		ctx,
		config.GroupName(),
		vaultName,
		keyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(config.Location()),
			Properties: &keyvault.VaultProperties{
				TenantID: &tenantID,
				Sku: &keyvault.Sku{
					Family: to.StringPtr("A"),
					Name:   keyvault.Standard,
				},
				AccessPolicies: &[]keyvault.AccessPolicyEntry{},
			},
		},
	)

	var vaultRet keyvault.Vault
	if (errRet == nil) {
		vaultRet, errRet = GetVault(ctx, vaultName)
	}

	return vaultRet, errRet
}

// GetVault returns an existing vault
func GetVault(ctx context.Context, vaultName string) (keyvault.Vault, error) {
	vaultsClient := getVaultsClient()
	return vaultsClient.Get(ctx, config.GroupName(), vaultName)
}

// CreateVaultWithPolicies creates a new Vault with policies granting access to the specified user.
func CreateVaultWithPolicies(ctx context.Context, vaultName, userID string) (vault keyvault.Vault, err error) {
	vaultsClient := getVaultsClient()

	tenantID, err := uuid.FromString(config.TenantID())
	if err != nil {
		return
	}

	apList := []keyvault.AccessPolicyEntry{}
	ap := keyvault.AccessPolicyEntry{
		TenantID: &tenantID,
		Permissions: &keyvault.Permissions{
			Keys: &[]keyvault.KeyPermissions{
				keyvault.KeyPermissionsCreate,
			},
			Secrets: &[]keyvault.SecretPermissions{
				keyvault.SecretPermissionsSet,
			},
		},
	}
	if userID != "" {
		ap.ObjectID = to.StringPtr(userID)
		apList = append(apList, ap)
	}

	_, errRet := vaultsClient.CreateOrUpdate(
		ctx,
		config.GroupName(),
		vaultName,
		keyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(config.Location()),
			Properties: &keyvault.VaultProperties{
				AccessPolicies:           &apList,
				EnabledForDiskEncryption: to.BoolPtr(true),
				Sku: &keyvault.Sku{
					Family: to.StringPtr("A"),
					Name:   keyvault.Standard,
				},
				TenantID: &tenantID,
			},
		},
	)

	var vaultRet keyvault.Vault
	if (errRet == nil) {
		vaultRet, errRet = GetVault(ctx, vaultName)
	}

	return vaultRet, errRet
}

// SetVaultPermissions adds an access policy permitting this app's Client ID to manage keys and secrets.
func SetVaultPermissions(ctx context.Context, vaultName string) (keyvault.Vault, error) {
	vaultsClient := getVaultsClient()

	tenantID, err := uuid.FromString(config.TenantID())
	if err != nil {
		return keyvault.Vault{}, err
	}

	clientID := config.ClientID()

	myObjID := "204b2841-e25d-455c-afd8-630a1244041e"

	valueForUpdate := keyvault.VaultCreateOrUpdateParameters{
		Location: to.StringPtr(config.Location()),
		Properties: &keyvault.VaultProperties{
			TenantID: &tenantID,
			Sku: &keyvault.Sku{
				Family: to.StringPtr("A"),
				Name:   keyvault.Standard,
			},
			AccessPolicies: &[]keyvault.AccessPolicyEntry{
				{
					ObjectID: &clientID,
					TenantID: &tenantID,
					Permissions: &keyvault.Permissions{
						Keys: &[]keyvault.KeyPermissions{
							keyvault.KeyPermissionsGet,
							keyvault.KeyPermissionsList,
							keyvault.KeyPermissionsCreate,
						},
						Secrets: &[]keyvault.SecretPermissions{
							keyvault.SecretPermissionsGet,
							keyvault.SecretPermissionsSet,
							keyvault.SecretPermissionsList,
						},
					},
				},
				{
					ObjectID: &myObjID,
					TenantID: &tenantID,
					Permissions: &keyvault.Permissions{
						Keys: &[]keyvault.KeyPermissions{
							keyvault.KeyPermissionsGet,
							keyvault.KeyPermissionsList,
							keyvault.KeyPermissionsCreate,
						},
						Secrets: &[]keyvault.SecretPermissions{
							keyvault.SecretPermissionsGet,
							keyvault.SecretPermissionsSet,
							keyvault.SecretPermissionsList,
						},
					},
				},
			},
		},
	}

	//buffN, _ := json.MarshalIndent(&valueForUpdate, "", " ")
	//fmt.Printf(string(buffN))

	_, errRet := vaultsClient.CreateOrUpdate(
		ctx,
		config.GroupName(),
		vaultName,
		valueForUpdate,
	)

	var vaultRet keyvault.Vault
	if (errRet == nil) {
		vaultRet, errRet = GetVault(ctx, vaultName)
	}

	return vaultRet, errRet
}

// SetVaultPermissionsForDeployment updates a key vault to enable deployments and add permissions to the application
func SetVaultPermissionsForDeployment(ctx context.Context, vaultName string) (keyvault.Vault, error) {
	vaultsClient := getVaultsClient()
	tenantID, err := uuid.FromString(config.TenantID())
	if err != nil {
		return keyvault.Vault{}, err
	}
	clientID := config.ClientID()

	_, errRet := vaultsClient.CreateOrUpdate(
		ctx,
		config.GroupName(),
		vaultName,
		keyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(config.Location()),
			Properties: &keyvault.VaultProperties{
				TenantID:                     &tenantID,
				EnabledForDeployment:         to.BoolPtr(true),
				EnabledForTemplateDeployment: to.BoolPtr(true),
				Sku: &keyvault.Sku{
					Family: to.StringPtr("A"),
					Name:   keyvault.Standard,
				},
				AccessPolicies: &[]keyvault.AccessPolicyEntry{
					{
						ObjectID: to.StringPtr(clientID),
						TenantID: &tenantID,
						Permissions: &keyvault.Permissions{
							Keys: &[]keyvault.KeyPermissions{
								keyvault.KeyPermissionsGet,
								keyvault.KeyPermissionsList,
								keyvault.KeyPermissionsCreate,
							},
							Secrets: &[]keyvault.SecretPermissions{
								keyvault.SecretPermissionsGet,
								keyvault.SecretPermissionsSet,
								keyvault.SecretPermissionsList,
							},
						},
					},
					{
						TenantID: &tenantID,
						Permissions: &keyvault.Permissions{
							Keys: &[]keyvault.KeyPermissions{
								keyvault.KeyPermissionsGet,
								keyvault.KeyPermissionsList,
								keyvault.KeyPermissionsCreate,
							},
							Secrets: &[]keyvault.SecretPermissions{
								keyvault.SecretPermissionsGet,
								keyvault.SecretPermissionsSet,
								keyvault.SecretPermissionsList,
							},
						},
					},
				},
			},
		},
	)

	var vaultRet keyvault.Vault
	if (errRet == nil) {
		vaultRet, errRet = GetVault(ctx, vaultName)
	}

	return vaultRet, errRet
}

// GetVaults lists all key vaults in a subscription
func GetVaults() {
	vaultsClient := getVaultsClient()

	fmt.Println("Getting all vaults in subscription")
	for subList, err := vaultsClient.ListComplete(context.Background(), nil); subList.NotDone(); err = subList.Next() {
		if err != nil {
			log.Printf("failed to get list of vaults: %v", err)
		}
		fmt.Printf("\t%s\n", *subList.Value().Name)
	}

	fmt.Println("Getting all vaults in resource group")
	for rgList, err := vaultsClient.ListByResourceGroupComplete(context.Background(), config.GroupName(), nil); rgList.NotDone(); err = rgList.Next() {
		if err != nil {
			log.Printf("failed to get list of vaults: %v", err)
		}
		fmt.Printf("\t%s\n", *rgList.Value().Name)
	}
}

// DeleteVault deletes an existing vault
func DeleteVault(ctx context.Context, vaultName string) (autorest.Response, error) {
	vaultsClient := getVaultsClient()
	return vaultsClient.Delete(ctx, config.GroupName(), vaultName)
}