package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/user-authentication-go/internal/api/handler"
	"github.com/imnzr/user-authentication-go/internal/api/middleware"
	"github.com/imnzr/user-authentication-go/internal/config"
	"github.com/imnzr/user-authentication-go/internal/database"
	"github.com/imnzr/user-authentication-go/internal/repository"
	"github.com/imnzr/user-authentication-go/internal/repository/redis"
	"github.com/imnzr/user-authentication-go/internal/service"
	"github.com/imnzr/user-authentication-go/pkg/auth"
	"go.uber.org/zap"
)

func New(cfg *config.Config, db *database.DB, logger *zap.Logger) *fiber.App {
	// Initialize auth manager
	authManager := auth.NewJWTManager(*cfg)

	// Initialize repository
	userRepo := repository.NewUserRepository(db.Primary)

	// Initialize transaction manager
	txManager := database.NewTxManager(db.Primary)

	// Initialize redis
	redisClient := redis.NewRedisClient(cfg.RedisCfg.RedisAddr, cfg.RedisCfg.RedisPass, cfg.RedisCfg.RedisDB)

	// Initialize services
	userService := service.NewUserService(userRepo, txManager, authManager, redisClient)

	// Initialize handle
	userHandler := handler.NewUserHandler(userService, logger, authManager)

	// Create Fiber APP
	app := fiber.New()

	// Global Middleware
	app.Use(middleware.CORS())

	// API Routes
	api := app.Group("/api/v1")

	// Auth Routes
	authRoutes := api.Group("/auth")
	authRoutes.Post("/signup", userHandler.CreateUser)
	authRoutes.Post("/signin", userHandler.LoginUser)
	authRoutes.Get("/profile", middleware.AuthMiddleware(context.Background(), authManager, *cfg, redisClient), userHandler.GetProfile)
	authRoutes.Get("/verify/:token", userHandler.VerifyEmail)
	authRoutes.Post("/logout", middleware.AuthMiddleware(context.Background(), authManager, *cfg, redisClient), userHandler.LogoutUser)

	return app
}
