package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
)

type OrganizationHandler struct {
	db          *gorm.DB
	orgService  *services.OrganizationService
	rbacService *services.RBACService
}

func NewOrganizationHandler(db *gorm.DB, orgService *services.OrganizationService, rbacService *services.RBACService) *OrganizationHandler {
	return &OrganizationHandler{
		db:          db,
		orgService:  orgService,
		rbacService: rbacService,
	}
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Create a new organization with the current user as owner
// @Tags organizations
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param organization body CreateOrganizationRequest true "Organization details"
// @Success 201 {object} models.Organization
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateOrganizationRequest
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

	// Create organization
	org, err := h.orgService.CreateOrganization(c.Request.Context(), req.Name, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, org)
}

// ListOrganizations godoc
// @Summary List all organizations for the current user
// @Description Get all organizations the current user is a member of
// @Tags organizations
// @Accept json
// @Produce json
// @Security WalletAuth
// @Success 200 {array} models.Organization
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations [get]
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
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

	// Get organizations
	orgs, err := h.orgService.ListUserOrganizations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve organizations"})
		return
	}

	c.JSON(200, orgs)
}

// GetOrganization godoc
// @Summary Get organization details
// @Description Get details of a specific organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} models.Organization
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check if user is member
	canRead, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, orgID, "org", "read")
	if err != nil || !canRead {
		c.JSON(403, gin.H{"error": "No permission to view organization"})
		return
	}

	// Get organization
	org, err := h.orgService.GetOrganization(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Organization not found"})
		return
	}

	// Get stats
	var vaultCount, credentialCount, projectCount int64
	h.db.Model(&models.Vault{}).Where("org_id = ?", orgID).Count(&vaultCount)
	h.db.Model(&models.Credential{}).
		Joins("JOIN vaults ON credentials.vault_id = vaults.id").
		Where("vaults.org_id = ?", orgID).
		Count(&credentialCount)
	h.db.Model(&models.Project{}).Where("org_id = ?", orgID).Count(&projectCount)

	response := OrganizationDetailResponse{
		Organization:    *org,
		VaultCount:      int(vaultCount),
		CredentialCount: int(credentialCount),
		ProjectCount:    int(projectCount),
	}

	c.JSON(200, response)
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update organization details
// @Tags organizations
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Param organization body UpdateOrganizationRequest true "Organization update"
// @Success 200 {object} models.Organization
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check permission
	canUpdate, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, orgID, "org", "update")
	if err != nil || !canUpdate {
		c.JSON(403, gin.H{"error": "No permission to update organization"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Get organization
	var org models.Organization
	if err := h.db.First(&org, "id = ?", orgID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Organization not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Plan != nil {
		org.Plan = *req.Plan
	}
	if req.MaxMembers != nil {
		org.MaxMembers = *req.MaxMembers
	}
	if req.MaxVaults != nil {
		org.MaxVaults = *req.MaxVaults
	}

	if err := h.db.Save(&org).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(200, org)
}

// DeleteOrganization godoc
// @Summary Delete organization
// @Description Delete an organization (owner only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check permission (only owner can delete)
	canDelete, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, orgID, "org", "delete")
	if err != nil || !canDelete {
		c.JSON(403, gin.H{"error": "Only owners can delete organizations"})
		return
	}

	// Delete organization (cascade will handle related records)
	if err := h.db.Delete(&models.Organization{}, "id = ?", orgID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.Status(204)
}

// Request/Response types

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateOrganizationRequest struct {
	Name       *string `json:"name"`
	Plan       *string `json:"plan"`
	MaxMembers *int    `json:"maxMembers"`
	MaxVaults  *int    `json:"maxVaults"`
}

type OrganizationDetailResponse struct {
	models.Organization
	VaultCount      int `json:"vaultCount"`
	CredentialCount int `json:"credentialCount"`
	ProjectCount    int `json:"projectCount"`
}

