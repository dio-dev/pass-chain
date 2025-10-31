# 🔐 Audit Logs with Vault Backup

## Overview
This document explains how Pass Chain stores audit logs in **both PostgreSQL (primary) and HashiCorp Vault (backup/verification)** for enhanced security and immutability.

---

## 🎯 Why Store Audit Logs in Vault?

### Benefits:
1. **Immutability** - Vault writes are append-only
2. **Backup verification** - Cross-check DB logs against Vault
3. **Tamper detection** - Any DB tampering is detectable
4. **Compliance** - Vault audit logs provide additional audit trail
5. **Zero-trust** - Even if DB is compromised, Vault logs remain

### Architecture:
```
User Action (e.g., Read Credential)
  ↓
Backend creates audit log
  ├─→ Write to PostgreSQL (fast, queryable)
  └─→ Write to Vault KV (immutable backup)
```

---

## 📂 Vault Storage Structure

### Path Format:
```
passchain/audit/:orgId/:resourceType/:resourceId/:timestamp-:action
```

### Examples:
```
passchain/audit/org-123/credential/cred-456/2025-10-30T12:34:56Z-read
passchain/audit/org-123/vault/vault-789/2025-10-30T12:35:00Z-create
passchain/audit/org-123/member/user-101/2025-10-30T12:36:00Z-invite
```

### Data Structure:
```json
{
  "id": "audit-uuid",
  "org_id": "org-123",
  "user_id": "user-456",
  "wallet_address": "0x742d...",
  "action": "read",
  "resource_type": "credential",
  "resource_id": "cred-789",
  "ip_hash": "sha256(ip+salt)",
  "device_fingerprint": "hash",
  "timestamp": "2025-10-30T12:34:56Z",
  "metadata": {
    "credential_name": "Production API Key",
    "vault_name": "Backend Vault"
  },
  "db_written": true,
  "vault_written": true
}
```

---

## 🔧 Implementation

### 1. Enhanced Audit Service

**File:** `backend/internal/services/audit_service.go`

```go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/pass-chain/backend/internal/models"
	"gorm.io/gorm"
)

type AuditService struct {
	db           *gorm.DB
	vaultClient  *VaultClient
}

type AuditEntry struct {
	ID                string                 `json:"id"`
	OrgID             string                 `json:"org_id"`
	UserID            string                 `json:"user_id"`
	WalletAddress     string                 `json:"wallet_address"`
	Action            string                 `json:"action"`
	ResourceType      string                 `json:"resource_type"`
	ResourceID        string                 `json:"resource_id"`
	IPHash            string                 `json:"ip_hash"`
	DeviceFingerprint string                 `json:"device_fingerprint"`
	Timestamp         time.Time              `json:"timestamp"`
	Metadata          map[string]interface{} `json:"metadata"`
	DBWritten         bool                   `json:"db_written"`
	VaultWritten      bool                   `json:"vault_written"`
	FabricTxHash      string                 `json:"fabric_tx_hash,omitempty"`
}

// LogAction logs an action to both DB and Vault
func (s *AuditService) LogAction(ctx context.Context, entry AuditEntry) error {
	entry.Timestamp = time.Now()

	// 1. Write to PostgreSQL (primary)
	dbLog := models.AuditLog{
		OrgID:             &entry.OrgID,
		UserID:            &entry.UserID,
		Action:            entry.Action,
		ResourceType:      entry.ResourceType,
		IPAddress:         entry.IPHash,
		DeviceFingerprint: entry.DeviceFingerprint,
		Timestamp:         entry.Timestamp,
	}

	if entry.ResourceType == "credential" {
		dbLog.CredentialID = &entry.ResourceID
	} else if entry.ResourceType == "vault" {
		dbLog.VaultID = &entry.ResourceID
	}

	if err := s.db.WithContext(ctx).Create(&dbLog).Error; err != nil {
		return fmt.Errorf("failed to write audit to DB: %w", err)
	}

	entry.ID = dbLog.ID
	entry.DBWritten = true

	// 2. Write to Vault (backup)
	vaultPath := fmt.Sprintf(
		"passchain/audit/%s/%s/%s/%s-%s",
		entry.OrgID,
		entry.ResourceType,
		entry.ResourceID,
		entry.Timestamp.Format(time.RFC3339),
		entry.Action,
	)

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit entry: %w", err)
	}

	if err := s.vaultClient.WriteKV(vaultPath, map[string]interface{}{
		"data": string(data),
	}); err != nil {
		// Log error but don't fail (Vault is backup)
		fmt.Printf("Warning: failed to write audit to Vault: %v\n", err)
	} else {
		entry.VaultWritten = true
	}

	// 3. Future: Anchor to Fabric (optional)
	// if s.fabricClient != nil {
	//     txHash, _ := s.fabricClient.LogAudit(entry)
	//     entry.FabricTxHash = txHash
	// }

	return nil
}

// GetAuditLogs retrieves audit logs from DB
func (s *AuditService) GetAuditLogs(ctx context.Context, filters AuditFilter) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := s.db.WithContext(ctx).Model(&models.AuditLog{})

	if filters.OrgID != "" {
		query = query.Where("org_id = ?", filters.OrgID)
	}
	if filters.UserID != "" {
		query = query.Where("user_id = ?", filters.UserID)
	}
	if filters.CredentialID != "" {
		query = query.Where("credential_id = ?", filters.CredentialID)
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("timestamp >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("timestamp <= ?", filters.EndDate)
	}

	if err := query.Order("timestamp DESC").Limit(filters.Limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditFromVault retrieves audit logs from Vault (verification)
func (s *AuditService) GetAuditFromVault(ctx context.Context, orgID string, year, month int) ([]AuditEntry, error) {
	// List all keys under path
	basePath := fmt.Sprintf("passchain/audit/%s", orgID)
	
	keys, err := s.vaultClient.ListKV(basePath)
	if err != nil {
		return nil, err
	}

	var entries []AuditEntry
	for _, key := range keys {
		data, err := s.vaultClient.ReadKV(key)
		if err != nil {
			continue
		}

		var entry AuditEntry
		if dataStr, ok := data["data"].(string); ok {
			if err := json.Unmarshal([]byte(dataStr), &entry); err == nil {
				// Filter by year/month
				if entry.Timestamp.Year() == year && int(entry.Timestamp.Month()) == month {
					entries = append(entries, entry)
				}
			}
		}
	}

	return entries, nil
}

// VerifyAuditIntegrity checks if DB and Vault logs match
func (s *AuditService) VerifyAuditIntegrity(ctx context.Context, credentialID string) (bool, error) {
	// Get DB logs
	dbLogs, err := s.GetAuditLogs(ctx, AuditFilter{
		CredentialID: credentialID,
		Limit:        1000,
	})
	if err != nil {
		return false, err
	}

	// For each DB log, verify in Vault
	mismatches := 0
	for _, dbLog := range dbLogs {
		// Construct expected Vault path
		vaultPath := fmt.Sprintf(
			"passchain/audit/%s/credential/%s/%s-%s",
			*dbLog.OrgID,
			*dbLog.CredentialID,
			dbLog.Timestamp.Format(time.RFC3339),
			dbLog.Action,
		)

		data, err := s.vaultClient.ReadKV(vaultPath)
		if err != nil {
			mismatches++
			continue
		}

		// Verify data matches
		var vaultEntry AuditEntry
		if dataStr, ok := data["data"].(string); ok {
			json.Unmarshal([]byte(dataStr), &vaultEntry)
			if vaultEntry.ID != dbLog.ID {
				mismatches++
			}
		}
	}

	return mismatches == 0, nil
}

type AuditFilter struct {
	OrgID        string
	UserID       string
	CredentialID string
	VaultID      string
	Action       string
	StartDate    time.Time
	EndDate      time.Time
	Limit        int
}
```

---

### 2. Vault Client Methods

**File:** `backend/internal/services/vault.go`

Add these methods to `VaultClient`:

```go
// WriteKV writes a key-value pair to Vault
func (c *VaultClient) WriteKV(path string, data map[string]interface{}) error {
	_, err := c.client.Logical().Write(path, data)
	return err
}

// ReadKV reads a key-value pair from Vault
func (c *VaultClient) ReadKV(path string) (map[string]interface{}, error) {
	secret, err := c.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no data found at path: %s", path)
	}
	return secret.Data, nil
}

// ListKV lists keys under a path
func (c *VaultClient) ListKV(path string) ([]string, error) {
	secret, err := c.client.Logical().List(path)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	result := make([]string, len(keys))
	for i, k := range keys {
		result[i] = k.(string)
	}
	return result, nil
}
```

---

### 3. API Endpoints

**File:** `backend/internal/handlers/audit_handler.go`

```go
// GetAuditLogsFromVault retrieves audit logs from Vault (verification endpoint)
func (h *AuditHandler) GetAuditLogsFromVault(c *fiber.Ctx) error {
	orgID := c.Params("orgId")
	year := c.QueryInt("year", time.Now().Year())
	month := c.QueryInt("month", int(time.Now().Month()))

	// Check permission
	userID := c.Locals("userId").(string)
	canAudit, err := h.rbacService.CheckPermission(c.Context(), userID, orgID, "org", "audit")
	if err != nil || !canAudit {
		return c.Status(403).JSON(fiber.Map{"error": "No permission to view audit logs"})
	}

	// Get logs from Vault
	logs, err := h.auditService.GetAuditFromVault(c.Context(), orgID, year, month)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve audit logs from Vault"})
	}

	return c.JSON(fiber.Map{
		"logs":   logs,
		"count":  len(logs),
		"source": "vault",
		"year":   year,
		"month":  month,
	})
}

// VerifyAuditIntegrity verifies DB and Vault logs match
func (h *AuditHandler) VerifyAuditIntegrity(c *fiber.Ctx) error {
	credentialID := c.Params("credentialId")

	// Check access
	userID := c.Locals("userId").(string)
	canAccess, err := h.rbacService.CheckCredentialAccess(c.Context(), userID, credentialID, "read")
	if err != nil || !canAccess {
		return c.Status(403).JSON(fiber.Map{"error": "No access to credential"})
	}

	// Verify integrity
	isValid, err := h.auditService.VerifyAuditIntegrity(c.Context(), credentialID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Verification failed"})
	}

	return c.JSON(fiber.Map{
		"credential_id": credentialID,
		"integrity":     isValid,
		"verified_at":   time.Now(),
	})
}
```

**Routes:**
```go
// Add to routes.go
auditGroup := api.Group("/audit")
auditGroup.Get("/organizations/:orgId/vault", auditHandler.GetAuditLogsFromVault)
auditGroup.Get("/credentials/:credentialId/verify", auditHandler.VerifyAuditIntegrity)
```

---

## 🔍 Usage Examples

### 1. Log a Credential Access
```go
auditService.LogAction(ctx, AuditEntry{
    OrgID:        "org-123",
    UserID:       "user-456",
    WalletAddress: "0x742d...",
    Action:       "read",
    ResourceType: "credential",
    ResourceID:   "cred-789",
    IPHash:       hashIP(clientIP),
    Metadata: map[string]interface{}{
        "credential_name": "Production API Key",
        "vault_name": "Backend Vault",
    },
})
```

### 2. Retrieve Audit Logs from DB
```bash
GET /api/v1/audit-logs?orgId=org-123&action=read&startDate=2025-10-01
```

### 3. Retrieve Audit Logs from Vault (Backup)
```bash
GET /api/v1/audit/organizations/org-123/vault?year=2025&month=10
```

### 4. Verify Audit Integrity
```bash
GET /api/v1/audit/credentials/cred-789/verify
```

Response:
```json
{
  "credential_id": "cred-789",
  "integrity": true,
  "verified_at": "2025-10-30T12:45:00Z"
}
```

---

## 🛡️ Security Considerations

### Vault Configuration:
```hcl
# Enable audit device in Vault
vault audit enable file file_path=/vault/logs/audit.log

# Enable KV v2 for audit storage
vault secrets enable -path=passchain/audit kv-v2
```

### Access Control:
```hcl
# Policy for audit writes (backend service)
path "passchain/audit/*" {
  capabilities = ["create", "update", "read", "list"]
}

# Policy for audit reads (security officers)
path "passchain/audit/:org_id/*" {
  capabilities = ["read", "list"]
}
```

---

## 📊 Monitoring & Alerts

### Metrics to Track:
- Audit write failures (DB vs Vault)
- Integrity check failures
- Missing Vault entries
- Vault API latency

### Alerts:
```yaml
- alert: AuditVaultWriteFailure
  expr: rate(audit_vault_write_errors[5m]) > 0.01
  annotations:
    summary: "High rate of Vault audit write failures"

- alert: AuditIntegrityMismatch
  expr: audit_integrity_check_failed > 0
  annotations:
    summary: "Audit log integrity check failed"
```

---

## 🔄 Backup & Retention

### Vault Backup:
```bash
# Periodic Vault snapshot
vault operator raft snapshot save backup-$(date +%Y%m%d).snap

# Store in S3
aws s3 cp backup-*.snap s3://passchain-vault-backups/
```

### Retention Policy:
- **PostgreSQL**: 90 days (hot data)
- **Vault**: 1 year (warm data)
- **S3 Archive**: 7 years (cold data, compliance)

---

## ✅ Implementation Checklist

- [ ] Update `AuditService` to write to Vault
- [ ] Add Vault KV methods to `VaultClient`
- [ ] Create audit verification endpoint
- [ ] Update credential handlers to use new audit service
- [ ] Add frontend "Verify Integrity" button
- [ ] Configure Vault audit device
- [ ] Set up monitoring alerts
- [ ] Document security officer workflow

---

**AUUUUFFFF!** 🔥

