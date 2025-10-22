package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/application/services"
	"auth-service/internal/infrastructure/auth"
	"auth-service/internal/infrastructure/http/handlers"
	"auth-service/internal/infrastructure/middleware"
	"auth-service/internal/infrastructure/persistence"
	"shared/pkg/database"
	sharedEvents "shared/pkg/events"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"

	_ "auth-service/docs" // Swagger docs

	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SMM Platform - Auth Service
// @version 1.0
// @description Authentication and Authorization microservice for SMM Platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

func main() {
	// Initialize database connection
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize event bus (using memory bus)
	eventBus := sharedEvents.NewMemoryEventBus()
	defer eventBus.Close()

	// Initialize repositories
	userRepo := persistence.NewPostgresUserRepository(db)
	sessionRepo := persistence.NewPostgresSessionRepository(db)

	// Initialize auth services
	tokenService := auth.NewTokenService(auth.TokenConfig{
		SecretKey:       getEnv("JWT_SECRET", "your-super-secret-jwt-key-here-change-in-production"),
		AccessTokenExp:  15 * time.Minute,   // Short-lived access tokens
		RefreshTokenExp: 7 * 24 * time.Hour, // Longer-lived refresh tokens
	})

	sessionManager := auth.NewSessionManager(sessionRepo, tokenService, 24*time.Hour)

	// Initialize universal event publisher
	eventPublisher := sharedEvents.NewUniversalEventPublisher(eventBus)

	// Initialize application services
	authService := services.NewAuthService(userRepo, sessionManager, eventPublisher, tokenService)

	// Initialize HTTP handlers
	authHandler := handlers.NewAuthHandler(authService, sessionManager)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, sessionManager)

	// Setup HTTP router with security middleware
	r := gin.Default()

	// Add security middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimit(100, time.Minute)) // 100 requests per minute
	r.Use(middleware.CORS())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	public := r.Group("/api/v1/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
		public.POST("/logout", authHandler.Logout)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "auth", "timestamp": time.Now()})
	})

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.GET("/sessions", authHandler.GetSessions)
		protected.POST("/sessions/revoke", authHandler.RevokeSession)
		protected.POST("/sessions/revoke-all", authHandler.RevokeAllSessions)
		protected.POST("/upgrade-tier", authHandler.UpgradeTier)
	}

	// Start HTTP server
	port := getEnv("PORT", "8081")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Auth service running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down auth service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Auth service exited properly")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
