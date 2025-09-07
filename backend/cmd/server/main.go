package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imnzr/user-authentication-go/internal/api/router"
	"github.com/imnzr/user-authentication-go/internal/config"
	"github.com/imnzr/user-authentication-go/internal/database"
	"github.com/imnzr/user-authentication-go/internal/repository/redis"
	"github.com/imnzr/user-authentication-go/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize Logger
	logger := logger.New()
	defer logger.Sync()

	logger.Info("Starting application...")

	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuratio: %v", err)
	}

	// Initialize Database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := redis.NewRedisClient(cfg.RedisCfg.RedisAddr, cfg.RedisCfg.RedisPass, cfg.RedisCfg.RedisDB)
	if err := redisClient.Ping(context.Background()); err != nil {
		logger.Fatal("failed to connect redis", zap.Error(err))
	}

	// Logger info database connected successfully
	logger.Info("Database connected successfully")

	// Run migration if enabled

	// Build router (Fiber App)
	app := router.New(cfg, db, logger)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		logger.Info("Starting Fiber server" + "address" + addr)

		if err := app.Listen(addr); err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Fiber server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Info("Server stopped gracefully")
	// Shutdown with timeout
}
