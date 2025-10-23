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
	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		walletAddress = c.Query("wallet")
	}

	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address required"})
		return
	}

	var auditLogs []models.AuditLog
	
	// Query audit logs for this wallet, ordered by most recent first
	if err := h.db.Where("wallet_address = ?", walletAddress).
		Order("timestamp DESC").
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

	h.logger.Info("Audit logs fetched", "wallet", walletAddress, "count", len(response))

	c.JSON(http.StatusOK, gin.H{
		"logs":  response,
		"count": len(response),
	})
}

// GetStats handles GET /api/v1/stats
func (h *AuditHandler) GetStats(c *gin.Context) {
	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		walletAddress = c.Query("wallet")
	}

	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address required"})
		return
	}

	// Count total credentials
	var totalCredentials int64
	h.db.Model(&models.Credential{}).Where("wallet_address = ?", walletAddress).Count(&totalCredentials)

	// Count total accesses
	var totalAccesses int64
	h.db.Model(&models.AuditLog{}).Where("wallet_address = ? AND action = ?", walletAddress, "read").Count(&totalAccesses)

	// Get last activity timestamp
	var lastLog models.AuditLog
	lastActivity := ""
	if err := h.db.Where("wallet_address = ?", walletAddress).
		Order("timestamp DESC").
		First(&lastLog).Error; err == nil {
		lastActivity = lastLog.Timestamp.Format(time.RFC3339)
	}

	// Vault shards = total credentials (each has one shard in Vault)
	vaultShards := totalCredentials

	// Blockchain shards = total credentials (each has one shard in blockchain)
	blockchainShards := totalCredentials

	h.logger.Info("Stats fetched", "wallet", walletAddress)

	c.JSON(http.StatusOK, gin.H{
		"totalCredentials": totalCredentials,
		"totalAccesses":    totalAccesses,
		"lastActivity":     lastActivity,
		"vaultShards":      vaultShards,
		"blockchainShards": blockchainShards,
	})
}

