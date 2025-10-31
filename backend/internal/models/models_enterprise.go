package models

import (
	"time"

	"github.com/lib/pq"
)

// ==================== Core User & Identity ====================

// User represents a wallet-based user identity
type User struct {
	ID                string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	WalletAddress     string    `gorm:"uniqueIndex;not null" json:"walletAddress"`
	ENSName           string    `json:"ensName,omitempty"`
	DisplayName       string    `json:"displayName,omitempty"`
	DeviceFingerprint string    `json:"-"` // Hidden from JSON
	CreatedAt         time.Time `json:"createdAt"`
	LastLoginAt       *time.Time `json:"lastLoginAt,omitempty"`
}

// ==================== Organization ====================

// Organization represents a company/team container
type Organization struct {
	ID         string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Slug       string    `gorm:"uniqueIndex;not null" json:"slug"`
	Plan       string    `gorm:"default:'free'" json:"plan"` // free, pro, enterprise
	MaxMembers int       `gorm:"default:5" json:"maxMembers"`
	MaxVaults  int       `gorm:"default:10" json:"maxVaults"`
	CreatedBy  string    `gorm:"type:uuid" json:"createdBy"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	
	// Relationships
	Creator User                  `gorm:"foreignKey:CreatedBy" json:"-"`
	Members []OrganizationMember  `gorm:"foreignKey:OrgID" json:"members,omitempty"`
	Vaults  []Vault               `gorm:"foreignKey:OrgID" json:"vaults,omitempty"`
}

// OrganizationMember represents user membership in an organization
type OrganizationMember struct {
	ID         string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	OrgID      string    `gorm:"type:uuid;not null;index:idx_org_user" json:"orgId"`
	UserID     string    `gorm:"type:uuid;not null;index:idx_org_user" json:"userId"`
	RoleID     string    `gorm:"type:uuid" json:"roleId"`
	Status     string    `gorm:"default:'invited'" json:"status"` // invited, active, suspended, removed
	InvitedBy  string    `gorm:"type:uuid" json:"invitedBy"`
	JoinedAt   *time.Time `json:"joinedAt,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	
	// Relationships
	Organization Organization `gorm:"foreignKey:OrgID" json:"-"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role         Role         `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Inviter      User         `gorm:"foreignKey:InvitedBy" json:"-"`
}

// ==================== RBAC (Roles & Permissions) ====================

// Role represents a role within an organization
type Role struct {
	ID          string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	OrgID       string    `gorm:"type:uuid;index:idx_org_role" json:"orgId"`
	Name        string    `gorm:"not null;index:idx_org_role" json:"name"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `gorm:"default:false" json:"isSystem"` // true for default roles
	CreatedAt   time.Time `json:"createdAt"`
	
	// Relationships
	Organization Organization `gorm:"foreignKey:OrgID" json:"-"`
	Permissions  []Permission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
}

// Permission represents a specific permission for a role
type Permission struct {
	ID           string `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	RoleID       string `gorm:"type:uuid;not null;index" json:"roleId"`
	ResourceType string `gorm:"not null" json:"resourceType"` // org, project, vault, credential
	Action       string `gorm:"not null" json:"action"` // create, read, update, delete, share, rotate, invite, audit
	
	// Relationships
	Role Role `gorm:"foreignKey:RoleID" json:"-"`
}

// ==================== Projects ====================

// Project represents a team workspace within an organization
type Project struct {
	ID          string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	OrgID       string    `gorm:"type:uuid;not null;index" json:"orgId"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `gorm:"type:uuid" json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relationships
	Organization Organization `gorm:"foreignKey:OrgID" json:"-"`
	Creator      User         `gorm:"foreignKey:CreatedBy" json:"-"`
	Vault        *Vault       `gorm:"foreignKey:ProjectID" json:"vault,omitempty"`
}

// ==================== Vaults ====================

// Vault represents a container for credentials
type Vault struct {
	ID          string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	OrgID       *string   `gorm:"type:uuid;index" json:"orgId,omitempty"` // NULL for personal vaults
	OwnerUserID *string   `gorm:"type:uuid;index" json:"ownerUserId,omitempty"` // Set for personal vaults
	ProjectID   *string   `gorm:"type:uuid;index" json:"projectId,omitempty"` // Optional project link
	VaultType   string    `gorm:"not null" json:"vaultType"` // personal, org, project
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `gorm:"type:uuid" json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relationships
	Organization *Organization  `gorm:"foreignKey:OrgID" json:"-"`
	Owner        *User          `gorm:"foreignKey:OwnerUserID" json:"owner,omitempty"`
	Project      *Project       `gorm:"foreignKey:ProjectID" json:"-"`
	Creator      User           `gorm:"foreignKey:CreatedBy" json:"-"`
	Credentials  []Credential   `gorm:"foreignKey:VaultID" json:"credentials,omitempty"`
	Access       []VaultAccess  `gorm:"foreignKey:VaultID" json:"access,omitempty"`
}

// VaultAccess represents explicit access grant to a vault
type VaultAccess struct {
	ID        string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	VaultID   string    `gorm:"type:uuid;not null;index:idx_vault_user" json:"vaultId"`
	UserID    string    `gorm:"type:uuid;not null;index:idx_vault_user" json:"userId"`
	RoleID    *string   `gorm:"type:uuid" json:"roleId,omitempty"`
	GrantedBy string    `gorm:"type:uuid" json:"grantedBy"`
	GrantedAt time.Time `json:"grantedAt"`
	
	// Relationships
	Vault    Vault `gorm:"foreignKey:VaultID" json:"-"`
	User     User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role     *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Granter  User  `gorm:"foreignKey:GrantedBy" json:"-"`
}

// ==================== Credentials (Updated) ====================

// Credential represents a secure record (now vault-scoped)
type Credential struct {
	ID               string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	VaultID          string         `gorm:"type:uuid;index" json:"vaultId"` // Removed NOT NULL for migration
	CredentialName   string         `gorm:"column:credential_name;not null" json:"name"`
	Username         string         `json:"username,omitempty"`
	URL              string         `json:"url,omitempty"`
	EncryptedData    string         `gorm:"not null" json:"encryptedData"`
	Nonce            string         `gorm:"not null" json:"nonce"`
	StorageRef       string         `json:"storageRef,omitempty"` // IPFS/S3 ref
	FabricCommitHash *string        `json:"fabricCommitHash,omitempty"`
	Share2           string         `json:"-"` // Fallback DB storage
	Tags             pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	CredentialType   string         `gorm:"default:'password'" json:"credentialType"` // password, api_key, ssh_key, cert, note
	CreatedBy        string         `gorm:"type:uuid" json:"createdBy"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	LastAccessed     *time.Time     `json:"lastAccessed,omitempty"`
	DeletedAt        *time.Time     `gorm:"index" json:"deletedAt,omitempty"`
	
	// Backward compatibility fields
	WalletAddress  string  `gorm:"index" json:"walletAddress,omitempty"`
	VaultPath      string  `json:"-"`                                        // Legacy Vault path
	BlockchainTxID string  `gorm:"column:blockchain_tx_id" json:"txId,omitempty"`  // Legacy blockchain TX
	
	// Relationships
	Vault   Vault      `gorm:"foreignKey:VaultID" json:"-"`
	Creator User       `gorm:"foreignKey:CreatedBy" json:"-"`
	Shards  []KeyShard `gorm:"foreignKey:CredentialID" json:"shards,omitempty"`
}

// KeyShard represents a reference to a key shard location
type KeyShard struct {
	ID           string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	CredentialID string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_cred_shard" json:"credentialId"`
	ShardIndex   int       `gorm:"not null;uniqueIndex:idx_cred_shard" json:"shardIndex"` // 1, 2, or 3
	Location     string    `gorm:"not null" json:"location"` // vault, fabric, client
	StorageRef   string    `gorm:"not null" json:"storageRef"` // Vault path or Fabric TX ID
	CreatedAt    time.Time `json:"createdAt"`
	
	// Relationships
	Credential Credential `gorm:"foreignKey:CredentialID" json:"-"`
}

// ==================== Audit ====================

// AuditLog represents an immutable access/action record (updated for multi-tenant)
type AuditLog struct {
	ID                string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	OrgID             *string   `gorm:"type:uuid;index" json:"orgId,omitempty"`
	UserID            *string   `gorm:"type:uuid;index" json:"userId,omitempty"`
	VaultID           *string   `gorm:"type:uuid;index" json:"vaultId,omitempty"`
	CredentialID      string    `gorm:"type:uuid;index" json:"credentialId"` // Changed to string for backward compat
	Action            string    `gorm:"index" json:"action"` // create, read, update, delete, share, rotate, invite, remove_member
	ResourceType      string    `json:"resourceType,omitempty"` // credential, vault, org, member
	Metadata          string    `gorm:"type:jsonb" json:"metadata,omitempty"` // Additional context
	IPAddress         string    `json:"ipAddress,omitempty"` // Hashed
	DeviceFingerprint string    `json:"-"`
	Timestamp         time.Time `gorm:"index" json:"timestamp"`
	FabricTxHash      *string   `json:"fabricTxHash,omitempty"`
	
	// Backward compatibility
	WalletAddress string `gorm:"index" json:"walletAddress,omitempty"` // Legacy field
	TxHash        string `json:"txHash,omitempty"` // Legacy blockchain hash
	
	// Relationships
	Organization *Organization `gorm:"foreignKey:OrgID" json:"-"`
	User         *User         `gorm:"foreignKey:UserID" json:"-"`
	Vault        *Vault        `gorm:"foreignKey:VaultID" json:"-"`
	Credential   *Credential   `gorm:"foreignKey:CredentialID" json:"-"`
}

// ==================== Invitations ====================

// Invitation represents an invite to join an organization
type Invitation struct {
	ID            string     `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	OrgID         string     `gorm:"type:uuid;not null;index" json:"orgId"`
	Email         *string    `json:"email,omitempty"`
	WalletAddress *string    `gorm:"index" json:"walletAddress,omitempty"`
	RoleID        string     `gorm:"type:uuid" json:"roleId"`
	InvitedBy     string     `gorm:"type:uuid" json:"invitedBy"`
	Token         string     `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt     time.Time  `json:"expiresAt"`
	AcceptedAt    *time.Time `json:"acceptedAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	
	// Relationships
	Organization Organization `gorm:"foreignKey:OrgID" json:"-"`
	Role         Role         `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Inviter      User         `gorm:"foreignKey:InvitedBy" json:"-"`
}

// ==================== Legacy Support ====================

// Keep old Payment struct for backwards compatibility (to be migrated)
type Payment struct {
	ID              string    `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	WalletAddress   string    `gorm:"index" json:"walletAddress"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Type            string    `json:"type"` // storage, access
	CredentialID    string    `gorm:"type:uuid;index" json:"credentialId,omitempty"`
	BlockchainTxID  string    `json:"blockchainTxId,omitempty"`
	Status          string    `gorm:"default:'pending'" json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}

