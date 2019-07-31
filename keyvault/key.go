package keyvault

import (
	"context"

	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/config"
	"github.com/WilliamMortlMicrosoft/AzureGoSamples/internal/iam"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/to"
)

func getKeysClient() keyvault.BaseClient {
	keyClient := keyvault.New()
	a, _ := iam.GetKeyvaultAuthorizer()
	keyClient.Authorizer = a
	keyClient.AddToUserAgent(config.UserAgent())
	return keyClient
}

// CreateKeyBundle creates a key in the specified keyvault
func CreateKey(ctx context.Context, vaultName, keyName string) (key keyvault.KeyBundle, err error) {
	vaultsClient := getVaultsClient()
	vault, err := vaultsClient.Get(ctx, config.GroupName(), vaultName)
	if err != nil {
		return
	}
	vaultURL := *vault.Properties.VaultURI

	keyClient := getKeysClient()
	return keyClient.CreateKey(
		ctx,
		vaultURL,
		keyName,
		keyvault.KeyCreateParameters{
			KeyAttributes: &keyvault.KeyAttributes{
				Enabled: to.BoolPtr(true),
			},
			KeySize: to.Int32Ptr(2048), // As of writing this sample, 2048 is the only supported KeySize.
			KeyOps: &[]keyvault.JSONWebKeyOperation{
				keyvault.Encrypt,
				keyvault.Decrypt,
				keyvault.Sign,
				keyvault.Verify,
				keyvault.WrapKey,
				keyvault.UnwrapKey,
			},
			Kty: keyvault.RSA,
		})
}
