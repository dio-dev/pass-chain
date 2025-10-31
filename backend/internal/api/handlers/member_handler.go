package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
)

type MemberHandler struct {
	db          *gorm.DB
	orgService  *services.OrganizationService
	rbacService *services.RBACService
}

func NewMemberHandler(db *gorm.DB, orgService *services.OrganizationService, rbacService *services.RBACService) *MemberHandler {
	return &MemberHandler{
		db:          db,
		orgService:  orgService,
		rbacService: rbacService,
	}
}

// ListMembers godoc
// @Summary List organization members
// @Description Get all members of an organization
// @Tags members
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Param status query string false "Filter by status (active, invited, suspended)"
// @Success 200 {array} MemberResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id}/members [get]
func (h *MemberHandler) ListMembers(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")
	statusFilter := c.Query("status")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check permission
	canRead, err := h.rbacService.CheckPermission(c.Request.Context(), user.ID, orgID, "org", "read")
	if err != nil || !canRead {
		c.JSON(403, gin.H{"error": "No permission to view members"})
		return
	}

	// Get members
	query := h.db.Preload("User").
		Preload("Role").
		Where("org_id = ?", orgID)

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var members []models.OrganizationMember
	if err := query.Find(&members).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve members"})
		return
	}

	// Format response
	response := make([]MemberResponse, len(members))
	for i, member := range members {
		response[i] = MemberResponse{
			ID:            member.ID,
			User:          member.User,
			Role:          member.Role,
			Status:        member.Status,
			JoinedAt:      member.JoinedAt,
			CreatedAt:     member.CreatedAt,
		}
	}

	c.JSON(200, response)
}

// InviteMember godoc
// @Summary Invite a member to the organization
// @Description Send an invitation to join the organization
// @Tags members
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Param invitation body InviteMemberRequest true "Invitation details"
// @Success 201 {object} models.Invitation
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id}/invite [post]
func (h *MemberHandler) InviteMember(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")

	var req InviteMemberRequest
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

	// Create invitation
	var email, wallet *string
	if req.Email != "" {
		email = &req.Email
	}
	if req.WalletAddress != "" {
		wallet = &req.WalletAddress
	}

	invitation, err := h.orgService.InviteMember(
		c.Request.Context(),
		orgID,
		user.ID,
		req.RoleID,
		email,
		wallet,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, invitation)
}

// PreviewInvitation godoc
// @Summary Preview an invitation
// @Description Get invitation details before accepting
// @Tags members
// @Accept json
// @Produce json
// @Param token path string true "Invitation Token"
// @Success 200 {object} InvitationPreview
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/invitations/{token} [get]
func (h *MemberHandler) PreviewInvitation(c *gin.Context) {
	token := c.Param("token")

	var invitation models.Invitation
	err := h.db.Preload("Organization").
		Preload("Role").
		Preload("Inviter").
		Where("token = ? AND accepted_at IS NULL AND expires_at > ?", token, time.Now()).
		First(&invitation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Invitation not found or expired"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve invitation"})
		}
		return
	}

	// Get org stats
	var memberCount int64
	h.db.Model(&models.OrganizationMember{}).
		Where("org_id = ? AND status = ?", invitation.OrgID, "active").
		Count(&memberCount)

	preview := InvitationPreview{
		Organization: invitation.Organization,
		Role:         invitation.Role,
		InvitedBy:    invitation.Inviter,
		ExpiresAt:    invitation.ExpiresAt,
		MemberCount:  int(memberCount),
	}

	c.JSON(200, preview)
}

// AcceptInvitation godoc
// @Summary Accept an invitation
// @Description Accept an invitation to join an organization
// @Tags members
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param token path string true "Invitation Token"
// @Success 200 {object} models.OrganizationMember
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/invitations/{token}/accept [post]
func (h *MemberHandler) AcceptInvitation(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	token := c.Param("token")

	// Get user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Accept invitation
	err := h.orgService.AcceptInvitation(c.Request.Context(), token, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get the created membership
	var invitation models.Invitation
	h.db.Where("token = ?", token).First(&invitation)

	var member models.OrganizationMember
	h.db.Preload("Role").
		Where("org_id = ? AND user_id = ?", invitation.OrgID, user.ID).
		First(&member)

	c.JSON(200, member)
}

// UpdateMemberRole godoc
// @Summary Update member's role
// @Description Change a member's role in the organization
// @Tags members
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Param userId path string true "User ID"
// @Param role body UpdateMemberRoleRequest true "New role"
// @Success 200 {object} models.OrganizationMember
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id}/members/{userId}/role [put]
func (h *MemberHandler) UpdateMemberRole(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")
	targetUserID := c.Param("userId")

	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Get current user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Update role
	err := h.orgService.UpdateMemberRole(c.Request.Context(), orgID, targetUserID, req.RoleID, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get updated member
	var member models.OrganizationMember
	h.db.Preload("Role").
		Where("org_id = ? AND user_id = ?", orgID, targetUserID).
		First(&member)

	c.JSON(200, member)
}

// RemoveMember godoc
// @Summary Remove a member from the organization
// @Description Remove a member from the organization
// @Tags members
// @Accept json
// @Produce json
// @Security WalletAuth
// @Param id path string true "Organization ID"
// @Param userId path string true "User ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/organizations/{id}/members/{userId} [delete]
func (h *MemberHandler) RemoveMember(c *gin.Context) {
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	orgID := c.Param("id")
	targetUserID := c.Param("userId")

	// Get current user
	var user models.User
	if err := h.db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Remove member
	err := h.orgService.RemoveMember(c.Request.Context(), orgID, targetUserID, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

// Request/Response types

type InviteMemberRequest struct {
	Email         string `json:"email"`
	WalletAddress string `json:"walletAddress"`
	RoleID        string `json:"roleId" binding:"required"`
}

type UpdateMemberRoleRequest struct {
	RoleID string `json:"roleId" binding:"required"`
}

type MemberResponse struct {
	ID        string              `json:"id"`
	User      models.User         `json:"user"`
	Role      models.Role         `json:"role"`
	Status    string              `json:"status"`
	JoinedAt  *time.Time          `json:"joinedAt"`
	CreatedAt time.Time           `json:"createdAt"`
}

type InvitationPreview struct {
	Organization models.Organization `json:"organization"`
	Role         models.Role         `json:"role"`
	InvitedBy    models.User         `json:"invitedBy"`
	ExpiresAt    time.Time           `json:"expiresAt"`
	MemberCount  int                 `json:"memberCount"`
}

