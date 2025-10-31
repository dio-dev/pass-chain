package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pass-chain/backend/internal/database"
	"pass-chain/backend/internal/models"
	"pass-chain/backend/internal/services"
	"pass-chain/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialHandler_CreateCredential(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test DB
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test services
	vaultService := &mockVaultService{}
	fabricClient := nil
	loggerInstance := logger.New("test")

	handler := NewCredentialHandler(db, vaultService, fabricClient, loggerInstance)

	// Test cases
	tests := []struct {
		name           string
		payload        models.CreateCredentialRequest
		setupFunc      func()
		expectedStatus int
		checkFunc      func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid credential creation",
			payload: models.CreateCredentialRequest{
				Name:          "Test Cred",
				Username:      "test@example.com",
				URL:           "https://example.com",
				EncryptedData: "encrypted",
				Nonce:         "nonce123",
				Share1:        "share1",
				Share2:        "share2",
				WalletAddress: "0x123",
				Signature:     "0xsig",
			},
			expectedStatus: http.StatusCreated,
			checkFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response["id"])
				assert.Equal(t, "Test Cred", response["name"])
			},
		},
		{
			name: "Missing required fields",
			payload: models.CreateCredentialRequest{
				Name: "Test",
				// Missing other required fields
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			// Create request
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/v1/credentials", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.CreateCredential(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkFunc != nil {
				tt.checkFunc(t, w)
			}
		})
	}
}

func TestCredentialHandler_GetCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, cleanup := setupTestDB(t)
	defer cleanup()

	vaultService := &mockVaultService{}
	handler := NewCredentialHandler(db, vaultService, nil, logger.New("test"))

	// Create test user and vault
	user := models.User{WalletAddress: "0x123"}
	db.Create(&user)

	vault := models.Vault{
		OwnerUserID: &user.ID,
		VaultType:   "personal",
		Name:        "Test Vault",
		CreatedBy:   user.ID,
	}
	db.Create(&vault)

	// Create test credential
	cred := models.Credential{
		VaultID:        vault.ID,
		CredentialName: "Test",
		Username:       "user",
		EncryptedData:  "data",
		Nonce:          "nonce",
		WalletAddress:  "0x123",
		CreatedBy:      user.ID,
	}
	db.Create(&cred)

	// Test request
	req := httptest.NewRequest("GET", "/api/v1/credentials", nil)
	req.Header.Set("X-Wallet-Address", "0x123")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetCredentials(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var creds []models.Credential
	err := json.Unmarshal(w.Body.Bytes(), &creds)
	require.NoError(t, err)
	assert.Len(t, creds, 1)
	assert.Equal(t, "Test", creds[0].CredentialName)
}

// Mock VaultService
type mockVaultService struct{}

func (m *mockVaultService) WriteSecret(path string, data map[string]interface{}) error {
	return nil
}

func (m *mockVaultService) ReadSecret(path string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"share1": "test_share1",
	}, nil
}

func (m *mockVaultService) DeleteSecret(path string) error {
	return nil
}

// Test helpers
func setupTestDB(t *testing.T) (*database.Database, func()) {
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

	cleanup := func() {
		// Clean up
		db.Exec("TRUNCATE users, organizations, vaults, credentials, audit_logs CASCADE")
		db.Close()
	}

	return db, cleanup
}

