package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pass-chain/backend/internal/database"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
	"pass-chain/backend/pkg/logger"
)

type CredentialHandler struct {
	db      *database.Database
	vault   *services.VaultService
	fabric  *services.FabricClient
	logger  *logger.Logger
}

func NewCredentialHandler(db *database.Database, vault *services.VaultService, fabric *services.FabricClient, log *logger.Logger) *CredentialHandler {
	return &CredentialHandler{
		db:     db,
		vault:  vault,
		fabric: fabric,
		logger: log,
	}
}

// CreateCredential handles POST /api/v1/credentials
func (h *CredentialHandler) CreateCredential(c *gin.Context) {
	var req models.CreateCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Verify wallet signature (from request body)
	if req.WalletAddress == "" || req.Signature == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing wallet authentication"})
		return
	}

	// TODO: Verify signature matches wallet address using EIP-191 or SIWE
	// For now, just log it
	h.logger.Info("Creating credential", "wallet", req.WalletAddress, "signature", req.Signature[:20]+"...")

	// Find or create user and personal vault for backward compatibility
	var user models.User
	if err := h.db.Where("wallet_address = ?", req.WalletAddress).FirstOrCreate(&user, models.User{
		WalletAddress: req.WalletAddress,
	}).Error; err != nil {
		h.logger.Error("Failed to find/create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	h.logger.Info("User found/created", "userId", user.ID, "wallet", user.WalletAddress)

	// Find or create personal vault
	var vault models.Vault
	if err := h.db.Where("owner_user_id = ? AND vault_type = ?", user.ID, "personal").FirstOrCreate(&vault, models.Vault{
		OwnerUserID: &user.ID,
		VaultType:   "personal",
		Name:        "Personal Vault",
		CreatedBy:   user.ID,
	}).Error; err != nil {
		h.logger.Error("Failed to find/create vault", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vault"})
		return
	}

	h.logger.Info("Vault found/created", "vaultId", vault.ID, "userId", user.ID)

	if vault.ID == "" {
		h.logger.Error("Vault ID is empty after FirstOrCreate")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vault - ID empty"})
		return
	}

	// Store Share1 in Vault
	// Vault KV v2 path format: secret/data/<path>
	vaultPath := "secret/data/passchain/" + req.WalletAddress + "/" + req.Name
	err := h.vault.WriteSecret(vaultPath, map[string]interface{}{
		"share1":     req.Share1,
		"created_at": "now",
	})
	if err != nil {
		h.logger.Error("Failed to store in Vault", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store encryption key"})
		return
	}

	// Store Share2 in Fabric blockchain (if available)
	var txID string
	if h.fabric != nil {
		txID, err = h.fabric.StoreShard(req.WalletAddress, req.Name, req.Share2)
		if err != nil {
			h.logger.Error("Failed to store in Fabric", "error", err)
			// Don't fail - fallback to DB storage
			txID = "fabric-unavailable"
		} else {
			h.logger.Info("Share2 stored in Fabric", "txID", txID)
		}
	} else {
		txID = "fabric-not-configured"
		h.logger.Warn("Fabric client not available, storing Share2 in database")
	}

	// Store credential in database
	credential := &models.Credential{
		ID:             uuid.New().String(),
		VaultID:        vault.ID,
		CredentialName: req.Name,
		Username:       req.Username,
		URL:            req.URL,
		EncryptedData:  req.EncryptedData,
		Nonce:          req.Nonce,
		WalletAddress:  req.WalletAddress,
		VaultPath:      vaultPath,
		BlockchainTxID: txID,
		Share2:         req.Share2, // Fallback storage in DB
		CreatedBy:      user.ID,
	}

	if err := h.db.Create(credential).Error; err != nil {
		h.logger.Error("Failed to save credential", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential"})
		return
	}

	// Log audit to database
	auditLog := &models.AuditLog{
		CredentialID:  credential.ID,
		WalletAddress: req.WalletAddress,
		Action:        "create",
		IPAddress:     hashIP(c.ClientIP()),
		Timestamp:     time.Now(),
		ResourceType:  "credential",
		Metadata:      "{}",
	}
	
	// Log to Fabric blockchain (optional)
	if h.fabric != nil {
		auditTxID, err := h.fabric.LogAccess(req.WalletAddress, credential.ID, "create", c.ClientIP())
		if err != nil {
			h.logger.Error("Failed to log to Fabric", "error", err)
		} else {
			auditLog.TxHash = auditTxID
			h.logger.Info("Audit logged to Fabric", "auditTxID", auditTxID)
		}
	}
	
	// Save audit log to database
	if err := h.db.Create(auditLog).Error; err != nil {
		h.logger.Error("Failed to save audit log", "error", err)
		// Don't fail - audit logging shouldn't break the operation
	}

	h.logger.Info("Credential created", "id", credential.ID, "name", credential.CredentialName, "txID", txID)

	c.JSON(http.StatusCreated, gin.H{
		"id":      credential.ID,
		"txId":    txID,
		"share3":  "client_backup", // Client should already have this
		"message": "Credential encrypted and saved successfully",
	})
}

// GetCredentials handles GET /api/v1/credentials
func (h *CredentialHandler) GetCredentials(c *gin.Context) {
	walletAddress := c.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		walletAddress = c.Query("wallet")
	}

	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address required"})
		return
	}

	var credentials []models.Credential
	if err := h.db.Where("wallet_address = ?", walletAddress).Find(&credentials).Error; err != nil {
		h.logger.Error("Failed to fetch credentials", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credentials"})
		return
	}

	h.logger.Info("Fetched credentials", "wallet", walletAddress, "count", len(credentials))

	c.JSON(http.StatusOK, credentials)
}

// GetCredentialByID handles GET /api/v1/credentials/:id
func (h *CredentialHandler) GetCredentialByID(c *gin.Context) {
	id := c.Param("id")
	walletAddress := c.GetHeader("X-Wallet-Address")

	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet authentication required"})
		return
	}

	var credential models.Credential
	if err := h.db.Where("id = ? AND wallet_address = ?", id, walletAddress).First(&credential).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"})
		return
	}

	// Retrieve Share1 from Vault
	vaultData, err := h.vault.ReadSecret(credential.VaultPath)
	if err != nil {
		h.logger.Error("Failed to read from Vault", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve encryption key"})
		return
	}

	share1, ok := vaultData["share1"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid vault data"})
		return
	}

	// Update last accessed
	now := time.Now()
	h.db.Model(&credential).Update("last_accessed", now)

	// Log audit to database
	auditLog := &models.AuditLog{
		CredentialID:  id,
		WalletAddress: walletAddress,
		Action:        "read",
		IPAddress:     hashIP(c.ClientIP()),
		Timestamp:     now,
		ResourceType:  "credential",
		Metadata:      "{}",
	}
	
	// Log access to Fabric (optional)
	if h.fabric != nil {
		auditTxID, err := h.fabric.LogAccess(walletAddress, id, "read", c.ClientIP())
		if err != nil {
			h.logger.Error("Failed to log access to Fabric", "error", err)
		} else {
			auditLog.TxHash = auditTxID
			h.logger.Info("Access logged to Fabric", "auditTxID", auditTxID)
		}
	}
	
	// Save audit log to database
	if err := h.db.Create(auditLog).Error; err != nil {
		h.logger.Error("Failed to save audit log", "error", err)
	}

	h.logger.Info("Credential accessed", "id", id, "wallet", walletAddress)

	c.JSON(http.StatusOK, gin.H{
		"id":            credential.ID,
		"name":          credential.CredentialName,
		"username":      credential.Username,
		"url":           credential.URL,
		"encryptedData": credential.EncryptedData,
		"nonce":         credential.Nonce,
		"walletAddress": credential.WalletAddress,
		"createdAt":     credential.CreatedAt,
		"lastAccessed":  credential.LastAccessed,
		"share1":        share1,
		"share2":        credential.Share2, // From blockchain (future)
		"share3":        "client_backup", // Client should provide from localStorage backup
	})
}

// RecoverCredential handles POST /api/v1/credentials/:id/recover with Share3
func (h *CredentialHandler) RecoverCredential(c *gin.Context) {
	id := c.Param("id")
	walletAddress := c.GetHeader("X-Wallet-Address")

	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet authentication required"})
		return
	}

	var req struct {
		Share3 string `json:"share3" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Share3 required"})
		return
	}

	var credential models.Credential
	if err := h.db.Where("id = ? AND wallet_address = ?", id, walletAddress).First(&credential).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"})
		return
	}

	// Retrieve Share1 from Vault
	vaultData, err := h.vault.ReadSecret(credential.VaultPath)
	if err != nil {
		h.logger.Error("Failed to read from Vault", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve encryption key"})
		return
	}

	share1, ok := vaultData["share1"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid vault data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            credential.ID,
		"name":          credential.CredentialName,
		"username":      credential.Username,
		"url":           credential.URL,
		"encryptedData": credential.EncryptedData,
		"nonce":         credential.Nonce,
		"share1":        share1,
		"share2":        credential.Share2,
		"share3":        req.Share3, // User-provided Share3
	})
}

// DeleteCredential handles DELETE /api/v1/credentials/:id
func (h *CredentialHandler) DeleteCredential(c *gin.Context) {
	id := c.Param("id")
	walletAddress := c.GetHeader("X-Wallet-Address")

	if walletAddress == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wallet authentication required"})
		return
	}

	var credential models.Credential
	if err := h.db.Where("id = ? AND wallet_address = ?", id, walletAddress).First(&credential).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"})
		return
	}

	// Delete from Vault
	if err := h.vault.DeleteSecret(credential.VaultPath); err != nil {
		h.logger.Error("Failed to delete from Vault", "error", err)
		// Continue anyway to delete from DB
	}

	// Delete from database
	if err := h.db.Delete(&credential).Error; err != nil {
		h.logger.Error("Failed to delete credential", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete credential"})
		return
	}

	// Log audit to database
	auditLog := &models.AuditLog{
		CredentialID:  id,
		WalletAddress: walletAddress,
		Action:        "delete",
		IPAddress:     hashIP(c.ClientIP()),
		Timestamp:     time.Now(),
		ResourceType:  "credential",
		Metadata:      "{}",
	}
	h.db.Create(auditLog)

	h.logger.Info("Credential deleted", "id", id, "wallet", walletAddress)

	c.JSON(http.StatusOK, gin.H{"message": "Credential deleted successfully"})
}

// hashIP creates a SHA-256 hash of the IP address for privacy
func hashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}
