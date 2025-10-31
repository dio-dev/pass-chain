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

	// Initialize enterprise services
	rbacService := services.NewRBACService(db.DB)
	orgService := services.NewOrganizationService(db.DB, rbacService)

	// Initialize handlers
	credHandler := handlers.NewCredentialHandler(db, vaultService, fabricClient, log)
	auditHandler := handlers.NewAuditHandler(db, log)
	
	// Enterprise handlers
	userHandler := handlers.NewUserHandler(db.DB, orgService, rbacService)
	orgHandler := handlers.NewOrganizationHandler(db.DB, orgService, rbacService)
	memberHandler := handlers.NewMemberHandler(db.DB, orgService, rbacService)
	projectHandler := handlers.NewProjectHandler(db.DB, orgService, rbacService)
	vaultHandler := handlers.NewVaultHandler(db.DB, rbacService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User/Profile routes
		v1.GET("/me", middleware.AuthRequired(), userHandler.GetMe)
		v1.PUT("/me", middleware.AuthRequired(), userHandler.UpdateMe)
		v1.GET("/me/organizations", middleware.AuthRequired(), userHandler.GetMyOrganizations)
		v1.GET("/me/permissions", middleware.AuthRequired(), userHandler.GetMyPermissions)

		// Organizations
		orgs := v1.Group("/organizations", middleware.AuthRequired())
		{
			orgs.POST("", orgHandler.CreateOrganization)
			orgs.GET("", orgHandler.ListOrganizations)
			orgs.GET("/:id", orgHandler.GetOrganization)
			orgs.PUT("/:id", orgHandler.UpdateOrganization)
			orgs.DELETE("/:id", orgHandler.DeleteOrganization)

			// Organization members
			orgs.GET("/:id/members", memberHandler.ListMembers)
			orgs.POST("/:id/invite", memberHandler.InviteMember)
			orgs.PUT("/:id/members/:userId/role", memberHandler.UpdateMemberRole)
			orgs.DELETE("/:id/members/:userId", memberHandler.RemoveMember)

			// Organization projects
			orgs.POST("/:id/projects", projectHandler.CreateProject)
			orgs.GET("/:id/projects", projectHandler.ListProjects)

			// Organization audit
			orgs.GET("/:id/audit", auditHandler.GetAuditLogs) // Will be enhanced
		}

		// Invitations (public routes - no auth middleware)
		invitations := v1.Group("/invitations")
		{
			invitations.GET("/:token", memberHandler.PreviewInvitation)
			invitations.POST("/:token/accept", middleware.AuthRequired(), memberHandler.AcceptInvitation)
		}

		// Projects
		projects := v1.Group("/projects", middleware.AuthRequired())
		{
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
		}

		// Vaults
		vaults := v1.Group("/vaults", middleware.AuthRequired())
		{
			vaults.POST("", vaultHandler.CreateVault)
			vaults.GET("", vaultHandler.ListVaults)
			vaults.GET("/:id", vaultHandler.GetVault)
			vaults.PUT("/:id", vaultHandler.UpdateVault)
			vaults.DELETE("/:id", vaultHandler.DeleteVault)

			// Vault access control
			vaults.POST("/:id/access", vaultHandler.GrantVaultAccess)
			vaults.DELETE("/:id/access/:userId", vaultHandler.RevokeVaultAccess)

			// Vault audit (will be enhanced)
			vaults.GET("/:id/audit", auditHandler.GetAuditLogs)
		}

	// Credentials (updated for vault-scoping)
	credentials := v1.Group("/credentials", middleware.AuthRequired())
	{
		credentials.POST("", credHandler.CreateCredential)
		credentials.GET("", credHandler.GetCredentials)
		credentials.GET("/:id", credHandler.GetCredentialByID)
		credentials.POST("/:id/recover", credHandler.RecoverCredential) // Recovery endpoint with Share3
		credentials.DELETE("/:id", credHandler.DeleteCredential)
		
		// Credential-specific audit logs
		credentials.GET("/:id/audit", auditHandler.GetCredentialAuditLogs)
	}

		// Audit logs & blockchain explorer
		v1.GET("/audit-logs", middleware.AuthRequired(), auditHandler.GetAuditLogs)
		v1.GET("/stats", middleware.AuthRequired(), auditHandler.GetStats)
	}

	return router
}
