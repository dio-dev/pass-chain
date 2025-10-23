package services

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type VaultService struct {
	client *vault.Client
}

func NewVaultService(address, token string) (*VaultService, error) {
	config := vault.DefaultConfig()
	config.Address = address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(token)

	return &VaultService{
		client: client,
	}, nil
}

// WriteSecret writes data to Vault KV v2
func (v *VaultService) WriteSecret(path string, data map[string]interface{}) error {
	// Vault KV v2 requires data to be wrapped in "data" key
	wrappedData := map[string]interface{}{
		"data": data,
	}

	_, err := v.client.Logical().Write(path, wrappedData)
	if err != nil {
		return fmt.Errorf("failed to write secret: %w", err)
	}

	return nil
}

// ReadSecret reads data from Vault KV v2
func (v *VaultService) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("secret not found")
	}

	// Vault KV v2 wraps data in "data" key
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format")
	}

	return data, nil
}

// DeleteSecret deletes a secret from Vault
func (v *VaultService) DeleteSecret(path string) error {
	_, err := v.client.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// ListSecrets lists secrets at a path
func (v *VaultService) ListSecrets(path string) ([]string, error) {
	secret, err := v.client.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	result := make([]string, len(keys))
	for i, k := range keys {
		result[i] = k.(string)
	}

	return result, nil
}

