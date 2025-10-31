package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pass-chain/backend/internal/api"
	"pass-chain/backend/internal/database"
	"pass-chain/backend/internal/services"
	"pass-chain/backend/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testWalletAddress = "0x7C1409F8280144B3e3FBD3F977A097f492b4289e"
	testSignature     = "0xtest123"
)

// TestIntegration runs full integration tests
func TestIntegration(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Setup services
	vaultService := services.NewVaultService("http://localhost:8200", "test-token")
	fabricClient := nil // Mock or skip Fabric for tests
	loggerInstance := logger.New("test")
	rbacService := services.NewRBACService(db.DB)
	orgService := services.NewOrganizationService(db.DB, rbacService)

	// Create router
	router := api.NewRouter(db, vaultService, fabricClient, loggerInstance, rbacService, orgService)

	// Run test suite
	t.Run("Health", func(t *testing.T) {
		testHealth(t, router)
	})

	t.Run("User", func(t *testing.T) {
		testUserEndpoints(t, router)
	})

	t.Run("Organization", func(t *testing.T) {
		testOrganizationEndpoints(t, router)
	})

	t.Run("Credentials", func(t *testing.T) {
		testCredentialEndpoints(t, router)
	})

	t.Run("AuditLogs", func(t *testing.T) {
		testAuditLogEndpoints(t, router)
	})
}

func testHealth(t *testing.T, router http.Handler) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func testUserEndpoints(t *testing.T, router http.Handler) {
	t.Run("GET /api/v1/me", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/me", nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var user map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &user)
		require.NoError(t, err)
		assert.Equal(t, testWalletAddress, user["walletAddress"])
	})

	t.Run("GET /api/v1/me/organizations", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/me/organizations", nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func testOrganizationEndpoints(t *testing.T, router http.Handler) {
	var orgID string

	t.Run("POST /api/v1/organizations - Create", func(t *testing.T) {
		payload := map[string]string{
			"name": "Test Organization",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var org map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &org)
		require.NoError(t, err)
		require.NotEmpty(t, org["id"])
		orgID = org["id"].(string)
		assert.Equal(t, "Test Organization", org["name"])
	})

	t.Run("GET /api/v1/organizations - List", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/organizations", nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var orgs []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &orgs)
		require.NoError(t, err)
		assert.NotEmpty(t, orgs)
	})

	t.Run("GET /api/v1/organizations/:id - Get by ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/organizations/"+orgID, nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func testCredentialEndpoints(t *testing.T, router http.Handler) {
	var credID string

	t.Run("POST /api/v1/credentials - Create", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":          "Test Credential",
			"username":      "test@example.com",
			"url":           "https://example.com",
			"encryptedData": "encrypted_data_here",
			"nonce":         "random_nonce",
			"share1":        "share1_data",
			"share2":        "share2_data",
			"walletAddress": testWalletAddress,
			"signature":     testSignature,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/credentials", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should create successfully
		if w.Code != http.StatusCreated {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var cred map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &cred)
		require.NoError(t, err)
		require.NotEmpty(t, cred["id"])
		credID = cred["id"].(string)
	})

	t.Run("GET /api/v1/credentials - List", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/credentials", nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var creds []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &creds)
		require.NoError(t, err)
		assert.NotEmpty(t, creds)
	})

	t.Run("GET /api/v1/credentials/:id - Get by ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/credentials/"+credID, nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// May fail if Vault is not available
		if w.Code != http.StatusOK {
			t.Logf("Warning: Credential retrieval failed (Vault may not be available): %s", w.Body.String())
		}
	})

	t.Run("DELETE /api/v1/credentials/:id - Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/credentials/"+credID, nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func testAuditLogEndpoints(t *testing.T, router http.Handler) {
	t.Run("GET /api/v1/audit-logs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/audit-logs", nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "logs")
		assert.Contains(t, response, "count")
	})

	t.Run("GET /api/v1/stats", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/stats", nil)
		req.Header.Set("X-Wallet-Address", testWalletAddress)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var stats map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		require.NoError(t, err)
		assert.Contains(t, stats, "totalCredentials")
	})
}

// Helper functions

func setupTestDB(t *testing.T) *database.Database {
	// Use test database configuration
	config := database.Config{
		Host:     "localhost",
		Port:     5432,
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

	return db
}

func teardownTestDB(t *testing.T, db *database.Database) {
	// Clean up test data
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	db.Close()
}

