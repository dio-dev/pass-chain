package models

import (
	"time"

	"gorm.io/gorm"
)

// Credential represents an encrypted credential
type Credential struct {
	ID             string         `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	CredentialName string         `gorm:"column:credential_name;not null" json:"name"`
	Username       string         `gorm:"not null" json:"username"`
	URL            string         `json:"url,omitempty"`
	EncryptedData  string         `gorm:"type:text;not null" json:"encryptedData"`
	Nonce          string         `gorm:"not null" json:"nonce"`
	WalletAddress  string         `gorm:"index;not null" json:"walletAddress"`
	VaultPath      string         `gorm:"not null" json:"-"`                    // Path in Vault for Share1
	BlockchainTxID string         `gorm:"column:blockchain_tx_id" json:"txId"`  // Fabric transaction ID (optional for now)
	Share2         string         `gorm:"type:text" json:"-"`                   // Will be stored in blockchain
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"-"`
	LastAccessed   *time.Time     `json:"lastAccessed,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateCredentialRequest for API
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

// AuditLog for blockchain audit trail
type AuditLog struct {
	ID            string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	CredentialID  string    `gorm:"type:uuid;index" json:"credentialId"`
	WalletAddress string    `gorm:"index" json:"walletAddress"`
	Action        string    `json:"action"` // "create", "read", "delete"
	IPAddress     string    `json:"ipAddress,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	TxHash        string    `json:"txHash,omitempty"` // Blockchain transaction hash
}
