package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tsh96/easy-count/server/internal/config"
	"github.com/tsh96/easy-count/server/internal/db"
	"github.com/tsh96/easy-count/server/internal/handlers"
	"github.com/tsh96/easy-count/server/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize database schema
	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Apply security middleware
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimit())

	// Initialize handlers
	h := handlers.NewHandlers(database, cfg)

	// Public routes (authentication)
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}

	// Protected routes (require authentication)
	api := router.Group("/api")
	api.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		// Transactions
		api.GET("/transactions", h.GetTransactions)
		api.POST("/transactions", h.CreateTransaction)
		api.PUT("/transactions/:id", h.UpdateTransaction)
		api.DELETE("/transactions/:id", h.DeleteTransaction)
		api.POST("/transactions/bulk", h.BulkCreateTransactions)

		// Customer records
		api.GET("/customer-records/:type", h.GetCustomerRecords)
		api.POST("/customer-records/:type", h.CreateCustomerRecord)
		api.PUT("/customer-records/:type/:id", h.UpdateCustomerRecord)
		api.DELETE("/customer-records/:type/:id", h.DeleteCustomerRecord)
		api.POST("/customer-records/:type/bulk", h.BulkCreateCustomerRecords)
		api.POST("/customer-records/:type/replace-name", h.ReplaceCustomerName)

		// Backup and restore
		api.GET("/backup/transactions", h.BackupTransactions)
		api.POST("/restore/transactions", h.RestoreTransactions)
		api.GET("/backup/customer-records", h.BackupCustomerRecords)
		api.POST("/restore/customer-records", h.RestoreCustomerRecords)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
