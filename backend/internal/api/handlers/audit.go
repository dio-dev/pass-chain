package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"pass-chain/backend/internal/database"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/pkg/logger"
)

type AuditHandler struct {
	db     *database.Database
	logger *logger.Logger
}

func NewAuditHandler(db *database.Database, log *logger.Logger) *AuditHandler {
	return &AuditHandler{
		db:     db,
		logger: log,
	}
}

// GetAuditLogs handles GET /api/v1/audit-logs
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	// Get wallet address from context (set by AuthRequired middleware) or header
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		walletAddress = c.GetHeader("X-Wallet-Address")
		if walletAddress == "" || walletAddress.(string) == "" {
			walletAddress = c.Query("wallet")
		}
	}

	walletAddr, ok := walletAddress.(string)
	if !ok || walletAddr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address required"})
		return
	}

	// Optional: filter by credential ID
	credentialID := c.Query("credentialId")

	var auditLogs []models.AuditLog
	query := h.db.Where("wallet_address = ?", walletAddr)
	
	// If credentialId is provided, filter by it
	if credentialID != "" {
		query = query.Where("credential_id = ?", credentialID)
	}
	
	// Query audit logs, ordered by most recent first
	if err := query.Order("timestamp DESC").
		Limit(100).
		Find(&auditLogs).Error; err != nil {
		h.logger.Error("Failed to fetch audit logs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	// Enrich with credential names
	type AuditLogResponse struct {
		ID             string    `json:"id"`
		CredentialID   string    `json:"credentialId"`
		CredentialName string    `json:"credentialName"`
		Action         string    `json:"action"`
		Timestamp      time.Time `json:"timestamp"`
		TxHash         string    `json:"txHash,omitempty"`
		IPHash         string    `json:"ipHash,omitempty"`
	}

	var response []AuditLogResponse
	for _, log := range auditLogs {
		// Fetch credential name
		var cred models.Credential
		credName := "Unknown"
		if err := h.db.Where("id = ?", log.CredentialID).First(&cred).Error; err == nil {
			credName = cred.CredentialName
		}

		response = append(response, AuditLogResponse{
			ID:             log.ID,
			CredentialID:   log.CredentialID,
			CredentialName: credName,
			Action:         log.Action,
			Timestamp:      log.Timestamp,
			TxHash:         log.TxHash,
			IPHash:         log.IPAddress,
		})
	}

	h.logger.Info("Audit logs fetched", "wallet", walletAddr, "count", len(response), "credentialId", credentialID)

	c.JSON(http.StatusOK, gin.H{
		"logs":  response,
		"count": len(response),
	})
}

// GetCredentialAuditLogs handles GET /api/v1/credentials/:id/audit
func (h *AuditHandler) GetCredentialAuditLogs(c *gin.Context) {
	credentialID := c.Param("id")
	if credentialID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credential ID required"})
		return
	}

	// Get wallet address from context
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	walletAddr := walletAddress.(string)

	// Verify the credential belongs to this wallet
	var cred models.Credential
	if err := h.db.Where("id = ? AND wallet_address = ?", credentialID, walletAddr).First(&cred).Error; err != nil {
		h.logger.Error("Credential not found or unauthorized", "error", err, "credentialId", credentialID, "wallet", walletAddr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"})
		return
	}

	// Fetch audit logs for this credential
	var auditLogs []models.AuditLog
	if err := h.db.Where("credential_id = ? AND wallet_address = ?", credentialID, walletAddr).
		Order("timestamp DESC").
		Limit(100).
		Find(&auditLogs).Error; err != nil {
		h.logger.Error("Failed to fetch credential audit logs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	// Format response
	type AuditLogResponse struct {
		ID             string    `json:"id"`
		CredentialID   string    `json:"credentialId"`
		CredentialName string    `json:"credentialName"`
		Action         string    `json:"action"`
		Timestamp      time.Time `json:"timestamp"`
		TxHash         string    `json:"txHash,omitempty"`
		IPHash         string    `json:"ipHash,omitempty"`
	}

	var response []AuditLogResponse
	for _, log := range auditLogs {
		response = append(response, AuditLogResponse{
			ID:             log.ID,
			CredentialID:   log.CredentialID,
			CredentialName: cred.CredentialName,
			Action:         log.Action,
			Timestamp:      log.Timestamp,
			TxHash:         log.TxHash,
			IPHash:         log.IPAddress,
		})
	}

	h.logger.Info("Credential audit logs fetched", "credentialId", credentialID, "wallet", walletAddr, "count", len(response))

	c.JSON(http.StatusOK, gin.H{
		"logs":  response,
		"count": len(response),
	})
}

// GetStats handles GET /api/v1/stats
func (h *AuditHandler) GetStats(c *gin.Context) {
	// Get wallet address from context (set by AuthRequired middleware) or header
	walletAddress, exists := c.Get("walletAddress")
	if !exists {
		walletAddress = c.GetHeader("X-Wallet-Address")
		if walletAddress == "" || walletAddress.(string) == "" {
			walletAddress = c.Query("wallet")
		}
	}

	walletAddr, ok := walletAddress.(string)
	if !ok || walletAddr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address required"})
		return
	}

	// Count total credentials
	var totalCredentials int64
	h.db.Model(&models.Credential{}).Where("wallet_address = ?", walletAddr).Count(&totalCredentials)

	// Count total accesses
	var totalAccesses int64
	h.db.Model(&models.AuditLog{}).Where("wallet_address = ? AND action = ?", walletAddr, "read").Count(&totalAccesses)

	// Get last activity timestamp
	var lastLog models.AuditLog
	lastActivity := ""
	if err := h.db.Where("wallet_address = ?", walletAddr).
		Order("timestamp DESC").
		First(&lastLog).Error; err == nil {
		lastActivity = lastLog.Timestamp.Format(time.RFC3339)
	}

	// Vault shards = total credentials (each has one shard in Vault)
	vaultShards := totalCredentials

	// Blockchain shards = total credentials (each has one shard in blockchain)
	blockchainShards := totalCredentials

	h.logger.Info("Stats fetched", "wallet", walletAddr)

	c.JSON(http.StatusOK, gin.H{
		"totalCredentials": totalCredentials,
		"totalAccesses":    totalAccesses,
		"lastActivity":     lastActivity,
		"vaultShards":      vaultShards,
		"blockchainShards": blockchainShards,
	})
}

