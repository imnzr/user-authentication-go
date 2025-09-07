package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/imnzr/user-authentication-go/internal/config"
	"github.com/imnzr/user-authentication-go/internal/repository/redis"
	"github.com/imnzr/user-authentication-go/pkg/auth"
)

func AuthMiddleware(ctx context.Context, jwtManager auth.AuthManager, cfg config.Config, redis *redis.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Error": "authorization header is missing",
			})
		}
		// Check redis blacklist
		if val, _ := redis.Client.Get(c.Context(), authHeader).Result(); val == "blacklisted" {
			return c.Status(500).JSON(fiber.Map{
				"Error": "Token revoked",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Error": "invalid authorization header format. Expected 'Bearer <token>'",
			})
		}
		tokenString := parts[1]

		claims, err := jwtManager.VerifyToken(ctx, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Erro": fmt.Sprintf("invalid or expired token: %v", err),
			})
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Error": "invalid user id in token claims",
			})
		}
		userID := int(userIdFloat)

		c.Locals("userId", userID)

		return c.Next()
	}
}

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowCredentials: true,
	})
}
