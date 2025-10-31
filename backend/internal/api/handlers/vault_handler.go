package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
)

type VaultHandler struct {
	db          *gorm.DB
	rbacService *services.RBACService
}

func NewVaultHandler(db *gorm.DB, rbacService *services.RBACService) *VaultHandler {
	return &VaultHandler{
		db:          db,
		rbacService: rbacService,
	}
}

// CreateVault godoc
// @Summary Create a new vault
// @Description Create a new vault in an organization
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param vault body CreateVaultRequest true "Vault details"
// @Success 201 {object} models.Vault
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults [post]
func (h *VaultHandler) CreateVault(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check permission if org vault
	if req.OrgID != nil {
		canCreate, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, *req.OrgID, "vault", "create")
		if err != nil || !canCreate {
			c.JSON(403, gin.H{"error": "No permission to create vault"})
			return
		}
	}

	// Create vault
	vault := models.Vault{
		OrgID:       req.OrgID,
		ProjectID:   req.ProjectID,
		VaultType:   req.VaultType,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   user.ID,
	}

	// Set owner for personal vaults
	if req.VaultType == "personal" {
		vault.OwnerUserID = &user.ID
	}

	if err := h.db.Create(&vault).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create vault"})
		return
	}

	c.JSON(201, vault)
}

// ListVaults godoc
// @Summary List accessible vaults
// @Description Get all vaults the user can access
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param orgId query string false "Filter by organization"
// @Param type query string false "Filter by type (personal, org, project)"
// @Success 200 {array} VaultWithStats
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults [get]
func (h *VaultHandler) ListVaults(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgIDFilter := c.Query("orgId")
	typeFilter := c.Query("type")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Build query for accessible vaults
	query := h.db.Model(&models.Vault{})

	// Personal vaults (owner)
	personalQuery := h.db.Where("owner_user_id = ?", user.ID)

	// Org vaults (via org membership + permissions)
	orgQuery := h.db.Joins("LEFT JOIN organizations ON vaults.org_id = organizations.id").
		Joins("LEFT JOIN organization_members ON organizations.id = organization_members.org_id").
		Where("organization_members.user_id = ? AND organization_members.status = ?", user.ID, "active")

	// Explicit vault access
	accessQuery := h.db.Joins("LEFT JOIN vault_access ON vaults.id = vault_access.vault_id").
		Where("vault_access.user_id = ?", user.ID)

	// Combine queries
	query = query.Where(personalQuery).Or(orgQuery).Or(accessQuery).Distinct()

	// Apply filters
	if orgIDFilter != "" {
		query = query.Where("org_id = ?", orgIDFilter)
	}
	if typeFilter != "" {
		query = query.Where("vault_type = ?", typeFilter)
	}

	var vaults []models.Vault
	if err := query.Find(&vaults).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve vaults"})
		return
	}

	// Add stats
	response := make([]VaultWithStats, len(vaults))
	for i, vault := range vaults {
		var credentialCount int64
		h.db.Model(&models.Credential{}).
			Where("vault_id = ?", vault.ID).
			Count(&credentialCount)

		response[i] = VaultWithStats{
			Vault:           vault,
			CredentialCount: int(credentialCount),
		}
	}

	c.JSON(200, response)
}

// GetVault godoc
// @Summary Get vault details
// @Description Get details of a specific vault
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Vault ID"
// @Success 200 {object} VaultDetail
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults/{id} [get]
func (h *VaultHandler) GetVault(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	vaultID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check access
	canAccess, err := h.rbacService.CheckVaultAccess(c.Request.Context(), user.ID, vaultID, "read")
	if err != nil || !canAccess {
		c.JSON(403, gin.H{"error": "No access to this vault"})
		return
	}

	// Get vault
	var vault models.Vault
	err = h.db.Preload("Organization").
		Preload("Project").
		First(&vault, "id = ?", vaultID).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "Vault not found"})
		return
	}

	// Get stats
	var credentialCount int64
	h.db.Model(&models.Credential{}).
		Where("vault_id = ?", vaultID).
		Count(&credentialCount)

	response := VaultDetail{
		Vault:           vault,
		CredentialCount: int(credentialCount),
	}

	c.JSON(200, response)
}

// UpdateVault godoc
// @Summary Update vault
// @Description Update vault details
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Vault ID"
// @Param vault body UpdateVaultRequest true "Vault update"
// @Success 200 {object} models.Vault
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults/{id} [put]
func (h *VaultHandler) UpdateVault(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	vaultID := c.Param("id")

	var req UpdateVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check access
	canUpdate, err := h.rbacService.CheckVaultAccess(c.Request.Context(), user.ID, vaultID, "update")
	if err != nil || !canUpdate {
		c.JSON(403, gin.H{"error": "No permission to update vault"})
		return
	}

	// Get vault
	var vault models.Vault
	if err := h.db.First(&vault, "id = ?", vaultID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Vault not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		vault.Name = *req.Name
	}
	if req.Description != nil {
		vault.Description = *req.Description
	}

	if err := h.db.Save(&vault).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update vault"})
		return
	}

	c.JSON(200, vault)
}

// DeleteVault godoc
// @Summary Delete vault
// @Description Delete a vault and all its credentials
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Vault ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults/{id} [delete]
func (h *VaultHandler) DeleteVault(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	vaultID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check access
	canDelete, err := h.rbacService.CheckVaultAccess(c.Request.Context(), user.ID, vaultID, "delete")
	if err != nil || !canDelete {
		c.JSON(403, gin.H{"error": "No permission to delete vault"})
		return
	}

	// Delete vault (cascade will handle credentials)
	if err := h.db.Delete(&models.Vault{}, "id = ?", vaultID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete vault"})
		return
	}

	c.Status(204)
}

// GrantVaultAccess godoc
// @Summary Grant vault access
// @Description Grant a user access to a vault
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Vault ID"
// @Param access body GrantAccessRequest true "Access details"
// @Success 201 {object} models.VaultAccess
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults/{id}/access [post]
func (h *VaultHandler) GrantVaultAccess(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	vaultID := c.Param("id")

	var req GrantAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check if user can share vault
	canShare, err := h.rbacService.CheckVaultAccess(c.Request.Context(), user.ID, vaultID, "update")
	if err != nil || !canShare {
		c.JSON(403, gin.H{"error": "No permission to share vault"})
		return
	}

	// Grant access
	err = h.rbacService.GrantVaultAccess(c.Request.Context(), vaultID, req.UserID, user.ID, req.RoleID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get created access record
	var access models.VaultAccess
	h.db.Preload("User").
		Preload("Role").
		Where("vault_id = ? AND user_id = ?", vaultID, req.UserID).
		First(&access)

	c.JSON(201, access)
}

// RevokeVaultAccess godoc
// @Summary Revoke vault access
// @Description Revoke a user's access to a vault
// @Tags vaults
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Vault ID"
// @Param userId path string true "User ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vaults/{id}/access/{userId} [delete]
func (h *VaultHandler) RevokeVaultAccess(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	vaultID := c.Param("id")
	targetUserID := c.Param("userId")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check permission
	canShare, err := h.rbacService.CheckVaultAccess(c.Request.Context(), user.ID, vaultID, "update")
	if err != nil || !canShare {
		c.JSON(403, gin.H{"error": "No permission to revoke access"})
		return
	}

	// Revoke access
	err = h.rbacService.RevokeVaultAccess(c.Request.Context(), vaultID, targetUserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

// Request/Response types

type CreateVaultRequest struct {
	OrgID       *string `json:"orgId"`
	ProjectID   *string `json:"projectId"`
	VaultType   string  `json:"vaultType" binding:"required"` // personal, org, project
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
}

type UpdateVaultRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type GrantAccessRequest struct {
	UserID string  `json:"userId" binding:"required"`
	RoleID *string `json:"roleId"`
}

type VaultWithStats struct {
	models.Vault
	CredentialCount int `json:"credentialCount"`
}

type VaultDetail struct {
	models.Vault
	CredentialCount int `json:"credentialCount"`
}

