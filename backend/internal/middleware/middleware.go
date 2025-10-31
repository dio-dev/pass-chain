package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"pass-chain/backend/pkg/logger"
)

// Logger is a middleware that logs HTTP requests
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Infow("HTTP Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
		)
	}
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Allow requests from localhost (dev) and minikube IP (prod)
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
			"http://192.168.49.2:30300",
		}
		
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}
		
		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// In development, allow all origins
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Wallet-Address, X-Signature")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// WalletAuth middleware for wallet-based authentication
func WalletAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, just pass through
		// TODO: Implement actual wallet signature verification
		c.Next()
	}
}

// AuthRequired middleware for wallet-based authentication
// Extracts wallet address from header and sets it in context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get wallet address from header
		walletAddress := c.GetHeader("X-Wallet-Address")
		if walletAddress == "" {
			// Try query parameter (for websocket/alternative auth)
			walletAddress = c.Query("wallet")
		}

		if walletAddress == "" {
			c.JSON(401, gin.H{"error": "Unauthorized: wallet address required"})
			c.Abort()
			return
		}

		// TODO: Verify wallet signature
		// signature := c.GetHeader("X-Signature")
		// message := c.GetHeader("X-Message")
		// if !verifySignature(walletAddress, signature, message) {
		//     c.JSON(401, gin.H{"error": "Invalid signature"})
		//     c.Abort()
		//     return
		// }

		// Set wallet address in context
		c.Set("walletAddress", walletAddress)
		c.Next()
	}
}

