package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pass-chain/backend/internal/api"
	"pass-chain/backend/internal/database"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
	"pass-chain/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ownerWallet  = "0xOwner123"
	memberWallet = "0xMember456"
	adminWallet  = "0xAdmin789"
)

// TestEndToEndFlow tests complete user journey
func TestEndToEndFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	db, router := setupE2ETest(t)
	defer cleanupE2ETest(t, db)

	// Flow: Owner creates organization → invites members → creates projects → manages credentials → RBAC
	t.Run("Complete User Journey", func(t *testing.T) {
		// Step 1: Owner registers (auto-create user)
		ownerID := testUserRegistration(t, router, ownerWallet)
		require.NotEmpty(t, ownerID)

		// Step 2: Owner creates organization
		orgID := testCreateOrganization(t, router, ownerWallet, "Acme Corp")
		require.NotEmpty(t, orgID)

		// Step 3: Owner creates a project
		projectID := testCreateProject(t, router, ownerWallet, orgID, "Production")
		require.NotEmpty(t, projectID)

		// Step 4: Owner invites admin
		inviteToken := testInviteMember(t, router, ownerWallet, orgID, adminWallet, "Admin")
		require.NotEmpty(t, inviteToken)

		// Step 5: Admin accepts invitation
		testAcceptInvitation(t, router, adminWallet, inviteToken)

		// Step 6: Admin invites member
		memberInvite := testInviteMember(t, router, adminWallet, orgID, memberWallet, "Member")
		require.NotEmpty(t, memberInvite)

		// Step 7: Member accepts invitation
		testAcceptInvitation(t, router, memberWallet, memberInvite)

		// Step 8: Admin creates project vault
		vaultID := testCreateVault(t, router, adminWallet, orgID, projectID, "Project Secrets")
		require.NotEmpty(t, vaultID)

		// Step 9: Admin creates credential in project vault
		credID := testCreateCredential(t, router, adminWallet, "AWS API Key", vaultID)
		require.NotEmpty(t, credID)

		// Step 10: Member tries to access credential (should succeed - read permission)
		testAccessCredential(t, router, memberWallet, credID, true)

		// Step 11: Member tries to delete credential (should fail - no delete permission)
		testDeleteCredential(t, router, memberWallet, credID, false)

		// Step 12: Admin deletes credential (should succeed)
		testDeleteCredential(t, router, adminWallet, credID, true)

		// Step 13: Owner removes member
		testRemoveMember(t, router, ownerWallet, orgID, memberWallet)

		// Step 14: Removed member tries to access org (should fail)
		testAccessOrganization(t, router, memberWallet, orgID, false)

		// Step 15: Owner views audit logs
		testViewAuditLogs(t, router, ownerWallet, orgID)
	})
}

func testUserRegistration(t *testing.T, router *gin.Engine, wallet string) string {
	t.Logf("👤 Step 1: User registration (%s)", wallet)

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "User registration failed: %s", w.Body.String())

	var user map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &user)
	
	userID := user["id"].(string)
	t.Logf("   ✓ User created: %s", userID)
	return userID
}

func testCreateOrganization(t *testing.T, router *gin.Engine, wallet, name string) string {
	t.Logf("🏢 Step 2: Create organization '%s'", name)

	payload := map[string]string{"name": name}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("   ❌ Response: %s", w.Body.String())
	}
	require.Equal(t, http.StatusCreated, w.Code, "Organization creation failed")

	var org map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &org)
	
	orgID := org["id"].(string)
	t.Logf("   ✓ Organization created: %s", orgID)
	return orgID
}

func testCreateProject(t *testing.T, router *gin.Engine, wallet, orgID, name string) string {
	t.Logf("📁 Step 3: Create project '%s'", name)

	payload := map[string]string{
		"name":        name,
		"description": "Test project",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/organizations/%s/projects", orgID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("   ⚠️  Response: %s", w.Body.String())
	}
	require.Equal(t, http.StatusCreated, w.Code, "Project creation failed")

	var project map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &project)
	
	projectID := project["id"].(string)
	t.Logf("   ✓ Project created: %s", projectID)
	return projectID
}

func testInviteMember(t *testing.T, router *gin.Engine, inviterWallet, orgID, memberWallet, role string) string {
	t.Logf("✉️  Step: Invite member (%s) as %s", memberWallet, role)

	payload := map[string]string{
		"walletAddress": memberWallet,
		"role":          role,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/organizations/%s/invite", orgID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wallet-Address", inviterWallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("   ⚠️  Response: %s", w.Body.String())
	}
	require.Equal(t, http.StatusCreated, w.Code, "Invitation failed")

	var invite map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &invite)
	
	token := invite["token"].(string)
	t.Logf("   ✓ Invitation sent: %s", token[:20]+"...")
	return token
}

func testAcceptInvitation(t *testing.T, router *gin.Engine, wallet, token string) {
	t.Logf("✅ Step: Accept invitation (%s)", wallet)

	req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/invitations/%s/accept", token), nil)
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Invitation acceptance failed: %s", w.Body.String())
	t.Logf("   ✓ Invitation accepted")
}

func testCreateVault(t *testing.T, router *gin.Engine, wallet, orgID, projectID, name string) string {
	t.Logf("🔐 Step: Create vault '%s'", name)

	payload := map[string]interface{}{
		"name":        name,
		"vaultType":   "project",
		"projectId":   projectID,
		"description": "Test vault",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/vaults", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("   ⚠️  Response: %s", w.Body.String())
	}
	require.Equal(t, http.StatusCreated, w.Code, "Vault creation failed")

	var vault map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &vault)
	
	vaultID := vault["id"].(string)
	t.Logf("   ✓ Vault created: %s", vaultID)
	return vaultID
}

func testCreateCredential(t *testing.T, router *gin.Engine, wallet, name, vaultID string) string {
	t.Logf("🔑 Step: Create credential '%s'", name)

	payload := map[string]interface{}{
		"name":          name,
		"username":      "admin",
		"url":           "https://aws.amazon.com",
		"encryptedData": "encrypted_key_data",
		"nonce":         "random_nonce",
		"share1":        "share1_data",
		"share2":        "share2_data",
		"walletAddress": wallet,
		"signature":     "0xsig",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/credentials", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("   ❌ Response: %s", w.Body.String())
	}
	require.Equal(t, http.StatusCreated, w.Code, "Credential creation failed")

	var cred map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cred)
	
	credID := cred["id"].(string)
	t.Logf("   ✓ Credential created: %s", credID)
	return credID
}

func testAccessCredential(t *testing.T, router *gin.Engine, wallet, credID string, shouldSucceed bool) {
	t.Logf("👁️  Step: Access credential (wallet: %s, should succeed: %v)", wallet, shouldSucceed)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/credentials/%s", credID), nil)
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if shouldSucceed {
		if w.Code != http.StatusOK {
			t.Logf("   ⚠️  Expected success but got: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code, "Credential access should succeed")
		t.Logf("   ✓ Access granted")
	} else {
		assert.NotEqual(t, http.StatusOK, w.Code, "Credential access should fail")
		t.Logf("   ✓ Access denied (as expected)")
	}
}

func testDeleteCredential(t *testing.T, router *gin.Engine, wallet, credID string, shouldSucceed bool) {
	t.Logf("🗑️  Step: Delete credential (wallet: %s, should succeed: %v)", wallet, shouldSucceed)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/credentials/%s", credID), nil)
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if shouldSucceed {
		assert.Equal(t, http.StatusOK, w.Code, "Credential deletion should succeed: %s", w.Body.String())
		t.Logf("   ✓ Deleted successfully")
	} else {
		assert.NotEqual(t, http.StatusOK, w.Code, "Credential deletion should fail")
		t.Logf("   ✓ Deletion denied (as expected)")
	}
}

func testRemoveMember(t *testing.T, router *gin.Engine, ownerWallet, orgID, memberWallet string) {
	t.Logf("👋 Step: Remove member (%s)", memberWallet)

	// Get member user ID first
	var member models.User
	// ... implementation

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/organizations/%s/members/%s", orgID, "memberID"), nil)
	req.Header.Set("X-Wallet-Address", ownerWallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Member removal failed: %s", w.Body.String())
	t.Logf("   ✓ Member removed")
}

func testAccessOrganization(t *testing.T, router *gin.Engine, wallet, orgID string, shouldSucceed bool) {
	t.Logf("🏢 Step: Access organization (wallet: %s, should succeed: %v)", wallet, shouldSucceed)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/organizations/%s", orgID), nil)
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if shouldSucceed {
		assert.Equal(t, http.StatusOK, w.Code, "Organization access should succeed")
		t.Logf("   ✓ Access granted")
	} else {
		assert.NotEqual(t, http.StatusOK, w.Code, "Organization access should fail")
		t.Logf("   ✓ Access denied (as expected)")
	}
}

func testViewAuditLogs(t *testing.T, router *gin.Engine, wallet, orgID string) {
	t.Logf("📜 Step: View audit logs")

	req := httptest.NewRequest("GET", "/api/v1/audit-logs", nil)
	req.Header.Set("X-Wallet-Address", wallet)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Audit log retrieval failed: %s", w.Body.String())

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	logs := response["logs"].([]interface{})
	t.Logf("   ✓ Audit logs retrieved: %d entries", len(logs))
}

// Setup helpers

func setupE2ETest(t *testing.T) (*database.Database, *gin.Engine) {
	config := database.Config{
		Host:     "localhost",
		Port:     5433, // Test database port
		User:     "passchain",
		Password: "passchain123",
		DBName:   "passchain_test",
		SSLMode:  "disable",
	}

	db, err := database.New(config)
	require.NoError(t, err)

	// Run migrations
	err = db.Migrate()
	require.NoError(t, err)

	// Setup services
	vaultService := services.NewVaultService("http://localhost:8201", "test-token")
	rbacService := services.NewRBACService(db.DB)
	orgService := services.NewOrganizationService(db.DB, rbacService)
	loggerInstance := logger.New("test")

	// Create router
	router := api.NewRouter(db, vaultService, nil, loggerInstance, rbacService, orgService)

	return db, router
}

func cleanupE2ETest(t *testing.T, db *database.Database) {
	// Truncate all tables
	db.Exec("TRUNCATE users, organizations, organization_members, roles, permissions, vaults, vault_access, projects, invitations, credentials, key_shards, audit_logs CASCADE")
	db.Close()
}

