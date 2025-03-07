package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/tusharsingune/meeting-scheduler/internal/api"
	"github.com/tusharsingune/meeting-scheduler/internal/config"
	"github.com/tusharsingune/meeting-scheduler/internal/logger"
	"github.com/tusharsingune/meeting-scheduler/internal/middleware"
	"github.com/tusharsingune/meeting-scheduler/internal/repository"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log, err := logger.Initialize("development")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Log startup information
	log.Info("Starting Meeting Scheduler API")
	log.Info("Environment variables",
		zap.String("DB_HOST", os.Getenv("DB_HOST")),
		zap.String("DB_PORT", os.Getenv("DB_PORT")),
		zap.String("DB_USER", os.Getenv("DB_USER")),
		zap.String("DB_NAME", os.Getenv("DB_NAME")),
		zap.String("SERVER_PORT", os.Getenv("SERVER_PORT")))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Log loaded configuration
	log.Info("Configuration loaded",
		zap.String("server.port", cfg.Server.Port),
		zap.String("database.host", cfg.Database.Host),
		zap.String("database.port", cfg.Database.Port),
		zap.String("database.user", cfg.Database.User),
		zap.String("database.dbname", cfg.Database.DBName))

	// Initialize router
	router := mux.NewRouter()

	// Add a basic health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.CORS)
	router.Use(middleware.Recovery)

	// Initialize database with retry logic
	var db repository.Repository
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Info("Attempting to connect to database",
			zap.Int("attempt", i+1),
			zap.Int("max_attempts", maxRetries))

		db, err = repository.NewPostgresDB(cfg.Database)
		if err == nil {
			log.Info("Successfully connected to database")
			break
		}

		log.Error("Failed to connect to database",
			zap.Error(err),
			zap.Int("attempt", i+1),
			zap.Duration("retry_delay", retryDelay))

		if i < maxRetries-1 {
			log.Info("Retrying database connection", zap.Duration("delay", retryDelay))
			time.Sleep(retryDelay)
		} else {
			log.Fatal("Failed to connect to database after multiple attempts", zap.Error(err))
		}
	}

	// Initialize API handlers
	api.RegisterHandlers(router, db)

	// Register Swagger UI
	api.RegisterSwagger(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited properly")
}
