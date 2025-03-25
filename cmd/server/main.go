package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alper.meric/messaging-system/api"
	"github.com/alper.meric/messaging-system/api/handlers"
	"github.com/alper.meric/messaging-system/clients"
	"github.com/alper.meric/messaging-system/config"
	"github.com/alper.meric/messaging-system/repository"
	"github.com/alper.meric/messaging-system/services"
	"github.com/gofiber/fiber/v2"
)

// @title Messaging System API
// @version 1.0
// @description Automatic message sending and tracking system API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	log.Println("Starting messaging system...")

	// Load configuration
	cfg := config.LoadConfig()

	// Set up PostgreSQL repository
	postgresRepo, err := repository.NewPostgresRepository(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL repository: %v", err)
	}
	log.Println("PostgreSQL repository created successfully")

	// Set up Redis repository
	var redisRepo *repository.RedisRepository
	redisRepo, err = repository.NewRedisRepository(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Printf("Warning: Failed to create Redis repository (caching will be disabled): %v", err)
		redisRepo = nil
	} else {
		log.Println("Redis repository created successfully")
	}

	// HTTP client oluşturma
	messageClient := clients.NewMessageClient(
		cfg.App.WebhookURL,
		cfg.App.MessageSendDryRun,
	)
	log.Println("HTTP client created successfully")

	// Prepare message sending service with repositories
	messageService := services.NewMessageService(cfg, postgresRepo, redisRepo, messageClient)

	// HTTP sunucusu ve API oluşturma
	app := fiber.New(fiber.Config{
		AppName:      "Messaging System",
		ErrorHandler: api.ErrorHandler,
	})

	// Controller sadece service'e bağımlı olmalı, repository'ye değil
	messageController := handlers.NewMessageController(messageService)

	// API endpoint'leri
	api.SetupRoutes(app, messageController)

	// Create a channel for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("Server listening on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	log.Println("Shutting down server...")

	// Give 5 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
