package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"shared/pkg/database"
	sharedEvents "shared/pkg/events"
	"user-service/internal/application/services"
	"user-service/internal/infrastructre/events"
	"user-service/internal/infrastructre/http/handlers"
	"user-service/internal/infrastructre/persistence"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "user-service/docs" // Swagger docs
)

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

	// Initialize infrastructure
	userRepo := persistence.NewPostgresUserRepository(db)

	// Initialize universal event publisher
	eventPublisher := sharedEvents.NewUniversalEventPublisher(eventBus)

	// Initialize application services
	userService := services.NewUserService(userRepo, eventPublisher)

	// Initialize event subscriber (now using universal subscriber)
	eventSubscriber := events.NewUniversalEventSubscriber(eventBus, userService, userRepo)

	// Start event consumers
	ctx := context.Background()
	if err := eventSubscriber.SubscribeToUserEvents(ctx); err != nil {
		log.Fatal("Failed to start event consumers:", err)
	}

	// Initialize HTTP handlers
	userHandler := handlers.NewUserHandler(userService)

	// Setup HTTP router
	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/users/:id", userHandler.GetUser)
		api.GET("/users/email/:email", userHandler.GetUserByEmail)
		api.POST("/users/:id/use-ai-description", userHandler.UseAIDescriptionQuota)
		api.POST("/users/:id/use-ai-video", userHandler.UseAIVideoQuota)
		api.POST("/users/:id/use-auto-posting", userHandler.UseAutoPostingQuota)
		api.POST("/users/:id/upgrade-pro", userHandler.UpgradeToPro)
		api.GET("/users/:id/check-ai-description-quota", userHandler.CheckAIDescriptionQuota)
		api.GET("/users/:id/check-ai-video-quota", userHandler.CheckAIVideoQuota)
		api.GET("/users/:id/check-auto-posting-quota", userHandler.CheckAutoPostingQuota)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user"})
	})

	// Admin routes
	admin := r.Group("/api/v1/admin")
	{
		admin.POST("/reset-monthly-quotas", userHandler.ResetMonthlyQuotas)
	}

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	go func() {
		log.Printf("User service running on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down user service...")
}
