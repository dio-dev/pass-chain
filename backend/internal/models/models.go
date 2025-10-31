package models

import (
	"time"
)

// CreateCredentialRequest for API (backward compatible)
type CreateCredentialRequest struct {
	Name          string `json:"name" binding:"required"`
	Username      string `json:"username" binding:"required"`
	URL           string `json:"url"`
	EncryptedData string `json:"encryptedData" binding:"required"`
	Nonce         string `json:"nonce" binding:"required"`
	Share1        string `json:"share1" binding:"required"` // Will be stored in Vault
	Share2        string `json:"share2" binding:"required"` // Will be stored in Blockchain
	WalletAddress string `json:"walletAddress" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
}

// Helper fields for backward compatibility
type CredentialLegacyFields struct {
	WalletAddress  string     `gorm:"index" json:"walletAddress"`
	VaultPath      string     `json:"-"`                    // Path in Vault for Share1
	BlockchainTxID string     `gorm:"column:blockchain_tx_id" json:"txId"`  // Fabric transaction ID
	DeletedAt      *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

