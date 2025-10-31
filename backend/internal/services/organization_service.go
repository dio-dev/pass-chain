package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"pass-chain/backend/internal/models"
	"gorm.io/gorm"
)

// OrganizationService handles organization operations
type OrganizationService struct {
	db          *gorm.DB
	rbacService *RBACService
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(db *gorm.DB, rbacService *RBACService) *OrganizationService {
	return &OrganizationService{
		db:          db,
		rbacService: rbacService,
	}
}

// CreateOrganization creates a new organization with default roles and the creator as Owner
func (s *OrganizationService) CreateOrganization(ctx context.Context, name string, creatorUserID string) (*models.Organization, error) {
	slug := generateSlug(name)

	org := models.Organization{
		Name:       name,
		Slug:       slug,
		Plan:       "free",
		MaxMembers: 5,
		MaxVaults:  10,
		CreatedBy:  creatorUserID,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create organization
		if err := tx.Create(&org).Error; err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}

		// Initialize default roles
		if err := s.rbacService.InitializeOrgRoles(ctx, org.ID); err != nil {
			return fmt.Errorf("failed to initialize roles: %w", err)
		}

		// Get Owner role
		ownerRole, err := s.rbacService.GetRoleByName(ctx, org.ID, "Owner")
		if err != nil {
			return fmt.Errorf("failed to get owner role: %w", err)
		}

		// Add creator as Owner
		now := time.Now()
		member := models.OrganizationMember{
			OrgID:    org.ID,
			UserID:   creatorUserID,
			RoleID:   ownerRole.ID,
			Status:   "active",
			JoinedAt: &now,
		}
		if err := tx.Create(&member).Error; err != nil {
			return fmt.Errorf("failed to add creator as member: %w", err)
		}

		// Create default organization vault
		vault := models.Vault{
			OrgID:     &org.ID,
			VaultType: "org",
			Name:      "Default Vault",
			Description: "Default organization vault for shared credentials",
			CreatedBy: creatorUserID,
		}
		if err := tx.Create(&vault).Error; err != nil {
			return fmt.Errorf("failed to create default vault: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &org, nil
}

// GetOrganization retrieves an organization by ID
func (s *OrganizationService) GetOrganization(ctx context.Context, orgID string) (*models.Organization, error) {
	var org models.Organization
	err := s.db.WithContext(ctx).
		Preload("Members").
		Preload("Members.User").
		Preload("Members.Role").
		First(&org, "id = ?", orgID).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// ListUserOrganizations lists all organizations a user is a member of
func (s *OrganizationService) ListUserOrganizations(ctx context.Context, userID string) ([]models.Organization, error) {
	var orgs []models.Organization
	err := s.db.WithContext(ctx).
		Joins("JOIN organization_members ON organization_members.org_id = organizations.id").
		Where("organization_members.user_id = ? AND organization_members.status = ?", userID, "active").
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// InviteMember creates an invitation for a user to join the organization
func (s *OrganizationService) InviteMember(ctx context.Context, orgID, invitedBy, roleID string, email, walletAddress *string) (*models.Invitation, error) {
	// Check if inviter has permission
	canInvite, err := s.rbacService.CheckPermission(ctx, invitedBy, orgID, "member", "invite")
	if err != nil || !canInvite {
		return nil, fmt.Errorf("no permission to invite members")
	}

	// Check member limit
	var memberCount int64
	s.db.WithContext(ctx).
		Model(&models.OrganizationMember{}).
		Where("org_id = ? AND status IN (?)", orgID, []string{"active", "invited"}).
		Count(&memberCount)

	var org models.Organization
	s.db.WithContext(ctx).First(&org, "id = ?", orgID)
	if int(memberCount) >= org.MaxMembers {
		return nil, fmt.Errorf("member limit reached")
	}

	// Generate invite token
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	invitation := models.Invitation{
		OrgID:         orgID,
		Email:         email,
		WalletAddress: walletAddress,
		RoleID:        roleID,
		InvitedBy:     invitedBy,
		Token:         token,
		ExpiresAt:     time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.db.WithContext(ctx).Create(&invitation).Error; err != nil {
		return nil, err
	}

	return &invitation, nil
}

// AcceptInvitation accepts an organization invitation
func (s *OrganizationService) AcceptInvitation(ctx context.Context, token, userID string) error {
	var invitation models.Invitation
	err := s.db.WithContext(ctx).
		Where("token = ? AND accepted_at IS NULL AND expires_at > ?", token, time.Now()).
		First(&invitation).Error
	if err != nil {
		return fmt.Errorf("invalid or expired invitation")
	}

	now := time.Now()
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update invitation
		invitation.AcceptedAt = &now
		if err := tx.Save(&invitation).Error; err != nil {
			return err
		}

		// Add member
		member := models.OrganizationMember{
			OrgID:     invitation.OrgID,
			UserID:    userID,
			RoleID:    invitation.RoleID,
			Status:    "active",
			InvitedBy: invitation.InvitedBy,
			JoinedAt:  &now,
		}
		return tx.Create(&member).Error
	})
}

// RemoveMember removes a member from an organization
func (s *OrganizationService) RemoveMember(ctx context.Context, orgID, memberUserID, removedBy string) error {
	// Check permission
	canRemove, err := s.rbacService.CheckPermission(ctx, removedBy, orgID, "member", "remove")
	if err != nil || !canRemove {
		return fmt.Errorf("no permission to remove members")
	}

	// Don't allow removing the last owner
	var ownerRole models.Role
	s.db.WithContext(ctx).Where("org_id = ? AND name = ?", orgID, "Owner").First(&ownerRole)

	var ownerCount int64
	s.db.WithContext(ctx).
		Model(&models.OrganizationMember{}).
		Where("org_id = ? AND role_id = ? AND status = ?", orgID, ownerRole.ID, "active").
		Count(&ownerCount)

	var member models.OrganizationMember
	s.db.WithContext(ctx).
		Where("org_id = ? AND user_id = ? AND status = ?", orgID, memberUserID, "active").
		First(&member)

	if member.RoleID == ownerRole.ID && ownerCount <= 1 {
		return fmt.Errorf("cannot remove the last owner")
	}

	// Update member status
	return s.db.WithContext(ctx).
		Model(&models.OrganizationMember{}).
		Where("org_id = ? AND user_id = ?", orgID, memberUserID).
		Update("status", "removed").Error
}

// UpdateMemberRole updates a member's role
func (s *OrganizationService) UpdateMemberRole(ctx context.Context, orgID, memberUserID, newRoleID, updatedBy string) error {
	// Check permission
	canUpdate, err := s.rbacService.CheckPermission(ctx, updatedBy, orgID, "member", "update")
	if err != nil || !canUpdate {
		return fmt.Errorf("no permission to update members")
	}

	return s.db.WithContext(ctx).
		Model(&models.OrganizationMember{}).
		Where("org_id = ? AND user_id = ?", orgID, memberUserID).
		Update("role_id", newRoleID).Error
}

// CreateProject creates a new project in an organization
func (s *OrganizationService) CreateProject(ctx context.Context, orgID, name, description, createdBy string) (*models.Project, error) {
	// Check permission
	canCreate, err := s.rbacService.CheckPermission(ctx, createdBy, orgID, "project", "create")
	if err != nil || !canCreate {
		return nil, fmt.Errorf("no permission to create projects")
	}

	project := models.Project{
		OrgID:       orgID,
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&project).Error; err != nil {
			return err
		}

		// Create project vault
		vault := models.Vault{
			OrgID:       &orgID,
			ProjectID:   &project.ID,
			VaultType:   "project",
			Name:        fmt.Sprintf("%s Vault", name),
			Description: fmt.Sprintf("Vault for project %s", name),
			CreatedBy:   createdBy,
		}
		return tx.Create(&vault).Error
	})

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Helper functions

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

