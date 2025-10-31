package database

import (
	"fmt"
	"pass-chain/backend/internal/models"
)

// MigrateEnterprise runs enterprise-specific database migrations and backfills data
func (db *Database) MigrateEnterprise() error {
	// Step 0: Fix audit_logs credential_id type if needed
	if err := db.FixAuditLogCredentialId(); err != nil {
		return fmt.Errorf("failed to fix audit_logs credential_id: %w", err)
	}
	
	// Step 1: Check if we need to add vault_id column
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'credentials' AND column_name = 'vault_id'").Scan(&count)
	
	needsVaultIdMigration := count == 0
	
	if needsVaultIdMigration {
		// Run backfill BEFORE adding the constraint
		if err := db.BackfillExistingCredentials(); err != nil {
			return fmt.Errorf("failed to backfill credentials: %w", err)
		}
	}
	
	// Step 2: AutoMigrate all new enterprise models
	err := db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.Role{},
		&models.Permission{},
		&models.Vault{},
		&models.VaultAccess{},
		&models.Project{},
		&models.Invitation{},
		&models.Credential{}, // Now includes vault_id
		&models.KeyShard{},
		&models.AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate enterprise models: %w", err)
	}

	return nil
}

// BackfillExistingCredentials migrates existing credentials to personal vaults
func (db *Database) BackfillExistingCredentials() error {
	// First, check if credentials table exists
	var tableExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'credentials')").Scan(&tableExists)
	
	if !tableExists {
		// Table doesn't exist yet, nothing to backfill
		return nil
	}

	// Check if there are any credentials without vault_id
	var credentialsNeedingVault []struct {
		ID            string
		WalletAddress string
	}
	
	// Use raw SQL to avoid GORM validation issues
	if err := db.Raw("SELECT id, wallet_address FROM credentials WHERE wallet_address IS NOT NULL AND wallet_address != '' AND (vault_id IS NULL OR vault_id = '00000000-0000-0000-0000-000000000000')").Scan(&credentialsNeedingVault).Error; err != nil {
		return fmt.Errorf("failed to find credentials needing vault: %w", err)
	}

	if len(credentialsNeedingVault) == 0 {
		// No credentials to migrate
		return nil
	}

	// Process each credential
	for _, cred := range credentialsNeedingVault {
		// Find or create user
		user := models.User{WalletAddress: cred.WalletAddress}
		if err := db.Where(&user).FirstOrCreate(&user).Error; err != nil {
			return fmt.Errorf("failed to find or create user %s: %w", cred.WalletAddress, err)
		}

		// Find or create personal vault for the user
		vault := models.Vault{
			OwnerUserID: &user.ID,
			VaultType:   "personal",
			Name:        "Personal Vault",
		}
		if err := db.Where("owner_user_id = ? AND vault_type = ?", user.ID, "personal").FirstOrCreate(&vault).Error; err != nil {
			return fmt.Errorf("failed to find or create personal vault for user %s: %w", user.WalletAddress, err)
		}

		// Update credential to link to the vault (using raw SQL to avoid validation)
		if err := db.Exec("UPDATE credentials SET vault_id = ?, created_by = ? WHERE id = ?", vault.ID, user.ID, cred.ID).Error; err != nil {
			return fmt.Errorf("failed to update credential %s with vault ID %s: %w", cred.ID, vault.ID, err)
		}
	}
	
	return nil
}

// FixAuditLogCredentialId fixes the credential_id column type in audit_logs from bigint to uuid
func (db *Database) FixAuditLogCredentialId() error {
	// Check if audit_logs table exists
	var tableExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'audit_logs')").Scan(&tableExists)
	
	if !tableExists {
		// Table doesn't exist yet, nothing to fix
		return nil
	}
	
	// Check if credential_id column exists
	var columnExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'audit_logs' AND column_name = 'credential_id')").Scan(&columnExists)
	
	if !columnExists {
		// Column doesn't exist yet, nothing to fix
		return nil
	}
	
	// Check current type of credential_id column
	var dataType string
	db.Raw(`SELECT data_type FROM information_schema.columns 
		WHERE table_schema = CURRENT_SCHEMA() 
		AND table_name = 'audit_logs' 
		AND column_name = 'credential_id'`).Scan(&dataType)
	
	// If it's already uuid, nothing to do
	if dataType == "uuid" {
		return nil
	}
	
	// If it's bigint or integer, drop and recreate as uuid
	// We can't preserve data because bigint can't be meaningfully converted to uuid
	if dataType == "bigint" || dataType == "integer" {
		if err := db.Exec("ALTER TABLE audit_logs DROP COLUMN credential_id").Error; err != nil {
			return fmt.Errorf("failed to drop old credential_id column: %w", err)
		}
	}
	
	return nil
}
