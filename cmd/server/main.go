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

	// Setup database config
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

	// Connect to the database
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseConnection(db)

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	productService := service.NewProductService(productRepo)

	// Initialize gRPC server
	server := grpc.NewServer(userService, productService, cfg.Server.Port)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logrus.Info("Shutting down server...")
		cancel()

		time.Sleep(cfg.Server.ShutdownTimeout)
		server.Stop()
	}()

	// Start the server
	if err := server.Start(); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}

	logrus.Info("Server shutdown complete")
}
