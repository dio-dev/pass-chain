package api

import (
	"github.com/gin-gonic/gin"
	"pass-chain/backend/internal/api/handlers"
	"pass-chain/backend/internal/config"
	"pass-chain/backend/internal/database"
	"pass-chain/backend/internal/middleware"
	"pass-chain/backend/internal/services"
	"pass-chain/backend/pkg/logger"
)

// NewRouter creates a new Gin router with all routes configured
func NewRouter(cfg *config.Config, db *database.Database, log *logger.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "pass-chain-api",
		})
	})

	// Initialize services
	vaultService, err := services.NewVaultService(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		log.Error("Failed to initialize Vault service", "error", err)
		// Continue without Vault for now
		vaultService = nil
	}

	// Initialize Fabric client (optional - gracefully handle if unavailable)
	var fabricClient *services.FabricClient
	fabricConfig := &services.FabricConfig{
		ConfigPath:   "./config/fabric-connection.yaml",
		ChannelID:    "passchain",
		ChaincodeID:  "credentials",
		OrgName:      "Org1",
		OrgUser:      "Admin",
		PeerEndpoint: "fabric-peer.fabric.svc.cluster.local:7051",
	}
	
	// Try to initialize Fabric - if it fails, log but continue
	fabricClient, err = services.NewFabricClient(fabricConfig, log)
	if err != nil {
		log.Warn("Fabric client not available, app will work without blockchain", "error", err)
		fabricClient = nil
	} else {
		log.Info("✅ Fabric client initialized - blockchain integration active!")
	}

	// Initialize handlers
	credHandler := handlers.NewCredentialHandler(db, vaultService, fabricClient, log)
	auditHandler := handlers.NewAuditHandler(db, log)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Credentials
		credentials := v1.Group("/credentials")
		{
			credentials.POST("", credHandler.CreateCredential)
			credentials.GET("", credHandler.GetCredentials)
			credentials.GET("/:id", credHandler.GetCredentialByID)
			credentials.DELETE("/:id", credHandler.DeleteCredential)
		}

		// Audit logs & blockchain explorer
		v1.GET("/audit-logs", auditHandler.GetAuditLogs)
		v1.GET("/stats", auditHandler.GetStats)
	}

	return router
}
