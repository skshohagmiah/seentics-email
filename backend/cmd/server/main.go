package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shohag/seentics-email/internal/config"
	"github.com/shohag/seentics-email/internal/database"
	"github.com/shohag/seentics-email/internal/handlers"
	"github.com/shohag/seentics-email/internal/middleware"
	"github.com/shohag/seentics-email/internal/postal"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		redisClient = nil // Disable Redis features
	} else {
		log.Println("Redis connected successfully")
	}

	// Initialize Postal client
	postalClient := postal.NewClient(cfg.PostalAPIURL, cfg.PostalAPIKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	apiKeyHandler := handlers.NewAPIKeyHandler()
	emailHandler := handlers.NewEmailHandler(postalClient)
	domainHandler := handlers.NewDomainHandler()
	webhookHandler := handlers.NewWebhookHandler()

	// Initialize middleware
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(redisClient)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login", authHandler.Login)
	}

	// Webhook endpoint (public, but verified)
	router.POST("/webhooks/postal", webhookHandler.HandlePostalWebhook)

	// Protected routes (JWT authentication)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// User profile
		api.GET("/profile", authHandler.GetProfile)

		// API Keys
		api.GET("/keys", apiKeyHandler.ListAPIKeys)
		api.POST("/keys", apiKeyHandler.CreateAPIKey)
		api.PUT("/keys/:id", apiKeyHandler.UpdateAPIKey)
		api.DELETE("/keys/:id", apiKeyHandler.DeleteAPIKey)

		// Emails (requires JWT)
		api.GET("/emails", emailHandler.ListEmails)
		api.GET("/emails/:id", emailHandler.GetEmail)

		// Domains
		api.GET("/domains", domainHandler.ListDomains)
		api.POST("/domains", domainHandler.AddDomain)
		api.GET("/domains/:id/verify", domainHandler.GetDomainVerification)
		api.POST("/domains/:id/verify", domainHandler.VerifyDomain)
		api.DELETE("/domains/:id", domainHandler.DeleteDomain)

		// Webhooks
		api.GET("/webhooks", webhookHandler.ListWebhooks)
		api.POST("/webhooks", webhookHandler.CreateWebhook)
		api.DELETE("/webhooks/:id", webhookHandler.DeleteWebhook)
	}

	// Email sending endpoint (API key authentication)
	send := router.Group("/api")
	send.Use(apiKeyMiddleware.Validate())
	{
		send.POST("/send", emailHandler.SendEmail)
	}

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
