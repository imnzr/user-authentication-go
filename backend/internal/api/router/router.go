package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/user-authentication-go/internal/api/handler"
	"github.com/imnzr/user-authentication-go/internal/config"
	"github.com/imnzr/user-authentication-go/internal/database"
	"github.com/imnzr/user-authentication-go/internal/repository"
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

	// Initialize services
	userService := service.NewUserService(userRepo, txManager, authManager)

	// Initialize handle
	userHandler := handler.NewUserHandler(userService, logger, authManager)

	// Create Fiber APP
	app := fiber.New()

	// Global Middleware

	// API Routes
	api := app.Group("/api/v1")

	// Auth Routes
	authRoutes := api.Group("/auth")
	authRoutes.Post("/signup", userHandler.CreateUser)
	authRoutes.Get("/verify/:token", userHandler.VerifyEmail)

	return app
}
