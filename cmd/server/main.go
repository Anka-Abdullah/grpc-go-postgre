package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grpc-exmpl/api/grpc"
	"grpc-exmpl/internal/config"
	"grpc-exmpl/internal/repository"
	"grpc-exmpl/internal/service"
	"grpc-exmpl/pkg/database"
	"grpc-exmpl/pkg/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/app.yaml")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Log.Level, cfg.Log.Format); err != nil {
		logrus.Fatalf("Failed to initialize logger: %v", err)
	}

	logrus.Info("Starting gRPC server application...")

	// Initialize database
	dbConfig := &database.Config{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		Database:     cfg.Database.Database,
		SSLMode:      cfg.Database.SSLMode,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxLifetime:  cfg.Database.MaxLifetime,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseConnection(db)

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)

	// Initialize gRPC server
	server := grpc.NewServer(userService, cfg.Server.Port)

	// Create context for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logrus.Info("Shutting down server...")
		cancel()

		// Give some time for graceful shutdown
		time.Sleep(cfg.Server.ShutdownTimeout)
		server.Stop()
	}()

	// Start server
	if err := server.Start(); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}

	logrus.Info("Server shutdown complete")
}
