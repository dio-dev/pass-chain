package handlers

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
)

type UserHandler struct {
	db          *gorm.DB
	orgService  *services.OrganizationService
	rbacService *services.RBACService
}

func NewUserHandler(db *gorm.DB, orgService *services.OrganizationService, rbacService *services.RBACService) *UserHandler {
	return &UserHandler{
		db:          db,
		orgService:  orgService,
		rbacService: rbacService,
	}
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security WalletAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create user if not exists (first time login)
			user = models.User{
				WalletAddress: walletAddress.(string),
				DisplayName:   formatWalletDisplay(walletAddress.(string)),
			}
			if err := h.db.Create(&user).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to create user"})
				return
			}

			// Create personal vault for new user
			if err := h.createPersonalVault(&user); err != nil {
				c.JSON(500, gin.H{"error": "Failed to create personal vault"})
				return
			}
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve user"})
			return
		}
	}

	c.JSON(200, user)
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param profile body UpdateProfileRequest true "Profile update"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.ENSName != nil {
		user.ENSName = *req.ENSName
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, user)
}

// GetMyOrganizations godoc
// @Summary List user's organizations
// @Description Get all organizations the current user is a member of
// @Tags users
// @Accept json
// @Produce json
// @Security WalletAuth
// @Success 200 {array} OrganizationWithRole
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/me/organizations [get]
func (h *UserHandler) GetMyOrganizations(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Get user's organizations with role info
	var memberships []models.OrganizationMember
	err := h.db.Preload("Organization").
		Preload("Role").
		Where("user_id = ? AND status = ?", user.ID, "active").
		Find(&memberships).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve organizations"})
		return
	}

	// Format response
	orgsWithRoles := make([]OrganizationWithRole, len(memberships))
	for i, member := range memberships {
		// Get member count
		var memberCount int64
		h.db.Model(&models.OrganizationMember{}).
			Where("org_id = ? AND status = ?", member.OrgID, "active").
			Count(&memberCount)

		// Get vault count
		var vaultCount int64
		h.db.Model(&models.Vault{}).
			Where("org_id = ?", member.OrgID).
			Count(&vaultCount)

		orgsWithRoles[i] = OrganizationWithRole{
			Organization: member.Organization,
			Role:         member.Role,
			JoinedAt:     member.JoinedAt,
			MemberCount:  int(memberCount),
			VaultCount:   int(vaultCount),
		}
	}

	c.JSON(200, orgsWithRoles)
}

// GetMyPermissions godoc
// @Summary Get user's permissions in an organization
// @Description Get the permissions the user has in a specific organization
// @Tags users
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param orgId query string true "Organization ID"
// @Success 200 {object} PermissionsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/me/permissions [get]
func (h *UserHandler) GetMyPermissions(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Query("orgId")
	if orgID == "" {
		c.JSON(400, gin.H{"error": "orgId is required"})
		return
	}

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Get user's membership
	var member models.OrganizationMember
	err := h.db.Preload("Role").
		Preload("Role.Permissions").
		Where("org_id = ? AND user_id = ? AND status = ?", orgID, user.ID, "active").
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Not a member of this organization"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve permissions"})
		}
		return
	}

	// Get accessible vaults
	var accessibleVaults []models.Vault
	h.db.Joins("LEFT JOIN vault_access ON vaults.id = vault_access.vault_id").
		Where("vaults.org_id = ? AND (vault_access.user_id = ? OR vaults.created_by = ?)", orgID, user.ID, user.ID).
		Distinct().
		Find(&accessibleVaults)

	response := PermissionsResponse{
		Role:              member.Role,
		Permissions:       member.Role.Permissions,
		AccessibleVaults:  accessibleVaults,
		AccessibleVaultIDs: make([]string, len(accessibleVaults)),
	}
	for i, vault := range accessibleVaults {
		response.AccessibleVaultIDs[i] = vault.ID
	}

	c.JSON(200, response)
}

// Helper functions

func (h *UserHandler) createPersonalVault(user *models.User) error {
	vault := models.Vault{
		OwnerUserID: &user.ID,
		VaultType:   "personal",
		Name:        "Personal Vault",
		Description: "Your personal password vault",
		CreatedBy:   user.ID,
	}
	return h.db.Create(&vault).Error
}

func formatWalletDisplay(wallet string) string {
	if len(wallet) > 10 {
		return wallet[:6] + "..." + wallet[len(wallet)-4:]
	}
	return wallet
}

// Request/Response types

type UpdateProfileRequest struct {
	DisplayName *string `json:"displayName"`
	ENSName     *string `json:"ensName"`
}

type OrganizationWithRole struct {
	Organization models.Organization `json:"organization"`
	Role         models.Role         `json:"role"`
	JoinedAt     *time.Time          `json:"joinedAt"`
	MemberCount  int                 `json:"memberCount"`
	VaultCount   int                 `json:"vaultCount"`
}

type PermissionsResponse struct {
	Role               models.Role         `json:"role"`
	Permissions        []models.Permission `json:"permissions"`
	AccessibleVaults   []models.Vault      `json:"accessibleVaults"`
	AccessibleVaultIDs []string            `json:"accessibleVaultIds"`
}

