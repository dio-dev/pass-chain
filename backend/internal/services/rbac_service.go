package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pass-chain/backend/internal/models"
	"gorm.io/gorm"
)

// RBACService handles role-based access control
type RBACService struct {
	db *gorm.DB
}

// NewRBACService creates a new RBAC service
func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{db: db}
}

// DefaultRole represents system roles
type DefaultRole struct {
	Name        string
	Description string
	Permissions []PermissionDef
}

// PermissionDef defines a permission
type PermissionDef struct {
	ResourceType string // org, project, vault, credential, member
	Action       string // create, read, update, delete, share, rotate, invite, audit
}

// GetDefaultRoles returns the system default roles
func GetDefaultRoles() []DefaultRole {
	return []DefaultRole{
		{
			Name:        "Owner",
			Description: "Full control over organization",
			Permissions: []PermissionDef{
				// All permissions on all resources
				{ResourceType: "org", Action: "read"},
				{ResourceType: "org", Action: "update"},
				{ResourceType: "org", Action: "delete"},
				{ResourceType: "org", Action: "audit"},
				{ResourceType: "member", Action: "invite"},
				{ResourceType: "member", Action: "remove"},
				{ResourceType: "member", Action: "update"},
				{ResourceType: "project", Action: "create"},
				{ResourceType: "project", Action: "read"},
				{ResourceType: "project", Action: "update"},
				{ResourceType: "project", Action: "delete"},
				{ResourceType: "vault", Action: "create"},
				{ResourceType: "vault", Action: "read"},
				{ResourceType: "vault", Action: "update"},
				{ResourceType: "vault", Action: "delete"},
				{ResourceType: "credential", Action: "create"},
				{ResourceType: "credential", Action: "read"},
				{ResourceType: "credential", Action: "update"},
				{ResourceType: "credential", Action: "delete"},
				{ResourceType: "credential", Action: "share"},
				{ResourceType: "credential", Action: "rotate"},
			},
		},
		{
			Name:        "Admin",
			Description: "Manage members, roles, and vaults",
			Permissions: []PermissionDef{
				{ResourceType: "org", Action: "read"},
				{ResourceType: "member", Action: "invite"},
				{ResourceType: "member", Action: "update"},
				{ResourceType: "project", Action: "create"},
				{ResourceType: "project", Action: "read"},
				{ResourceType: "project", Action: "update"},
				{ResourceType: "vault", Action: "create"},
				{ResourceType: "vault", Action: "read"},
				{ResourceType: "vault", Action: "update"},
				{ResourceType: "credential", Action: "create"},
				{ResourceType: "credential", Action: "read"},
				{ResourceType: "credential", Action: "update"},
				{ResourceType: "credential", Action: "delete"},
				{ResourceType: "credential", Action: "share"},
			},
		},
		{
			Name:        "Security Officer",
			Description: "View audit logs and security settings",
			Permissions: []PermissionDef{
				{ResourceType: "org", Action: "read"},
				{ResourceType: "org", Action: "audit"},
				{ResourceType: "vault", Action: "read"},
				{ResourceType: "credential", Action: "read"},
			},
		},
		{
			Name:        "Member",
			Description: "Read/write assigned vaults",
			Permissions: []PermissionDef{
				{ResourceType: "org", Action: "read"},
				{ResourceType: "vault", Action: "read"},
				{ResourceType: "credential", Action: "create"},
				{ResourceType: "credential", Action: "read"},
				{ResourceType: "credential", Action: "update"},
				{ResourceType: "credential", Action: "share"},
			},
		},
		{
			Name:        "Auditor",
			Description: "Read-only access + audit logs",
			Permissions: []PermissionDef{
				{ResourceType: "org", Action: "read"},
				{ResourceType: "org", Action: "audit"},
				{ResourceType: "vault", Action: "read"},
				{ResourceType: "credential", Action: "read"},
			},
		},
		{
			Name:        "Guest",
			Description: "Limited read access to shared credentials",
			Permissions: []PermissionDef{
				{ResourceType: "org", Action: "read"},
				{ResourceType: "credential", Action: "read"},
			},
		},
	}
}

// InitializeOrgRoles creates default roles for a new organization
func (s *RBACService) InitializeOrgRoles(ctx context.Context, orgID string) error {
	return s.InitializeOrgRolesWithTx(ctx, s.db, orgID)
}

// InitializeOrgRolesWithTx creates default roles with a transaction
func (s *RBACService) InitializeOrgRolesWithTx(ctx context.Context, tx *gorm.DB, orgID string) error {
	defaultRoles := GetDefaultRoles()

	for _, dr := range defaultRoles {
		// Create role
		role := models.Role{
			ID:          uuid.New().String(), // Explicitly generate UUID
			OrgID:       orgID,
			Name:        dr.Name,
			Description: dr.Description,
			IsSystem:    true,
		}
		if err := tx.WithContext(ctx).Create(&role).Error; err != nil {
			return fmt.Errorf("failed to create role %s: %w", dr.Name, err)
		}

		// Create permissions
		for _, perm := range dr.Permissions {
			permission := models.Permission{
				ID:           uuid.New().String(), // Explicitly generate UUID
				RoleID:       role.ID,
				ResourceType: perm.ResourceType,
				Action:       perm.Action,
			}
			if err := tx.WithContext(ctx).Create(&permission).Error; err != nil {
				return fmt.Errorf("failed to create permission for role %s: %w", dr.Name, err)
			}
		}
	}

	return nil
}

// GetRoleByName gets a role by name within an organization
func (s *RBACService) GetRoleByName(ctx context.Context, orgID, roleName string) (*models.Role, error) {
	var role models.Role
	err := s.db.WithContext(ctx).
		Preload("Permissions").
		Where("org_id = ? AND name = ?", orgID, roleName).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CheckPermission checks if a user has a specific permission
func (s *RBACService) CheckPermission(ctx context.Context, userID, orgID, resourceType, action string) (bool, error) {
	// Get user's role in the organization
	var member models.OrganizationMember
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND org_id = ? AND status = ?", userID, orgID, "active").
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// Check if role has the permission
	var count int64
	err = s.db.WithContext(ctx).
		Model(&models.Permission{}).
		Where("role_id = ? AND resource_type = ? AND action = ?", member.RoleID, resourceType, action).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckVaultAccess checks if user can access a vault
func (s *RBACService) CheckVaultAccess(ctx context.Context, userID, vaultID, action string) (bool, error) {
	var vault models.Vault
	err := s.db.WithContext(ctx).First(&vault, "id = ?", vaultID).Error
	if err != nil {
		return false, err
	}

	// Personal vault - only owner can access
	if vault.VaultType == "personal" {
		if vault.OwnerUserID != nil && *vault.OwnerUserID == userID {
			return true, nil
		}
		return false, nil
	}

	// Check explicit vault access
	var vaultAccess models.VaultAccess
	err = s.db.WithContext(ctx).
		Where("vault_id = ? AND user_id = ?", vaultID, userID).
		First(&vaultAccess).Error
	if err == nil {
		// Has explicit access
		return true, nil
	}

	// Check org permission
	if vault.OrgID != nil {
		return s.CheckPermission(ctx, userID, *vault.OrgID, "vault", action)
	}

	return false, nil
}

// CheckCredentialAccess checks if user can access a credential
func (s *RBACService) CheckCredentialAccess(ctx context.Context, userID, credentialID, action string) (bool, error) {
	var credential models.Credential
	err := s.db.WithContext(ctx).
		Preload("Vault").
		First(&credential, "id = ?", credentialID).Error
	if err != nil {
		return false, err
	}

	// Check vault access first
	return s.CheckVaultAccess(ctx, userID, credential.VaultID, action)
}

// GrantVaultAccess grants a user access to a vault
func (s *RBACService) GrantVaultAccess(ctx context.Context, vaultID, userID, grantedBy string, roleID *string) error {
	access := models.VaultAccess{
		VaultID:   vaultID,
		UserID:    userID,
		RoleID:    roleID,
		GrantedBy: grantedBy,
	}
	return s.db.WithContext(ctx).Create(&access).Error
}

// RevokeVaultAccess removes a user's access to a vault
func (s *RBACService) RevokeVaultAccess(ctx context.Context, vaultID, userID string) error {
	return s.db.WithContext(ctx).
		Where("vault_id = ? AND user_id = ?", vaultID, userID).
		Delete(&models.VaultAccess{}).Error
}

// CreateCustomRole creates a custom role for an organization
func (s *RBACService) CreateCustomRole(ctx context.Context, orgID, name, description string, permissions []PermissionDef) (*models.Role, error) {
	role := models.Role{
		OrgID:       orgID,
		Name:        name,
		Description: description,
		IsSystem:    false,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		for _, perm := range permissions {
			permission := models.Permission{
				RoleID:       role.ID,
				ResourceType: perm.ResourceType,
				Action:       perm.Action,
			}
			if err := tx.Create(&permission).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &role, nil
}

