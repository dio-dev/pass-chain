package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
)

type ProjectHandler struct {
	db          *gorm.DB
	orgService  *services.OrganizationService
	rbacService *services.RBACService
}

func NewProjectHandler(db *gorm.DB, orgService *services.OrganizationService, rbacService *services.RBACService) *ProjectHandler {
	return &ProjectHandler{
		db:          db,
		orgService:  orgService,
		rbacService: rbacService,
	}
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project in an organization
// @Tags projects
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Param project body CreateProjectRequest true "Project details"
// @Success 201 {object} ProjectWithVault
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id}/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")

	var req CreateProjectRequest
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

	// Create project
	project, err := h.orgService.CreateProject(c.Request.Context(), orgID, req.Name, req.Description, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get project vault
	var vault models.Vault
	h.db.Where("project_id = ?", project.ID).First(&vault)

	response := ProjectWithVault{
		Project: *project,
		Vault:   vault,
	}

	c.JSON(201, response)
}

// ListProjects godoc
// @Summary List organization projects
// @Description Get all projects in an organization
// @Tags projects
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Success 200 {array} ProjectWithStats
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id}/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
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
	canRead, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, orgID, "project", "read")
	if err != nil || !canRead {
		c.JSON(403, gin.H{"error": "No permission to view projects"})
		return
	}

	// Get projects
	var projects []models.Project
	if err := h.db.Preload("Vault").Where("org_id = ?", orgID).Find(&projects).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	// Add stats
	response := make([]ProjectWithStats, len(projects))
	for i, project := range projects {
		var credentialCount int64
		if project.Vault != nil {
			h.db.Model(&models.Credential{}).
				Where("vault_id = ?", project.Vault.ID).
				Count(&credentialCount)
		}

		response[i] = ProjectWithStats{
			Project:         project,
			CredentialCount: int(credentialCount),
		}
	}

	c.JSON(200, response)
}

// GetProject godoc
// @Summary Get project details
// @Description Get details of a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Project ID"
// @Success 200 {object} ProjectDetail
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Get project
	var project models.Project
	err := h.db.Preload("Vault").
		Preload("Organization").
		First(&project, "id = ?", projectID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Project not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve project"})
		}
		return
	}

	// Check permission
	canRead, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, project.OrgID, "project", "read")
	if err != nil || !canRead {
		c.JSON(403, gin.H{"error": "No permission to view project"})
		return
	}

	// Get stats
	var credentialCount int64
	if project.Vault != nil {
		h.db.Model(&models.Credential{}).
			Where("vault_id = ?", project.Vault.ID).
			Count(&credentialCount)
	}

	response := ProjectDetail{
		Project:         project,
		CredentialCount: int(credentialCount),
	}

	c.JSON(200, response)
}

// UpdateProject godoc
// @Summary Update project
// @Description Update project details
// @Tags projects
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Project ID"
// @Param project body UpdateProjectRequest true "Project update"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")

	var req UpdateProjectRequest
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

	// Get project
	var project models.Project
	if err := h.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	// Check permission
	canUpdate, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, project.OrgID, "project", "update")
	if err != nil || !canUpdate {
		c.JSON(403, gin.H{"error": "No permission to update project"})
		return
	}

	// Update fields
	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}

	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(200, project)
}

// DeleteProject godoc
// @Summary Delete project
// @Description Delete a project and its vault
// @Tags projects
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Project ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := c.Param("id")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Get project
	var project models.Project
	if err := h.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	// Check permission
	canDelete, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, project.OrgID, "project", "delete")
	if err != nil || !canDelete {
		c.JSON(403, gin.H{"error": "No permission to delete project"})
		return
	}

	// Delete project (cascade will handle vault)
	if err := h.db.Delete(&project).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(204)
}

// Request/Response types

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ProjectWithVault struct {
	models.Project
	Vault models.Vault `json:"vault"`
}

type ProjectWithStats struct {
	models.Project
	CredentialCount int `json:"credentialCount"`
}

type ProjectDetail struct {
	models.Project
	CredentialCount int `json:"credentialCount"`
}

