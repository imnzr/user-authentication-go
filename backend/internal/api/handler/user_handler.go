package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/user-authentication-go/internal/domain/user"
	errorpkg "github.com/imnzr/user-authentication-go/internal/pkg/error_pkg"
	"github.com/imnzr/user-authentication-go/pkg/auth"
	"github.com/imnzr/user-authentication-go/pkg/request"
	"go.uber.org/zap"
)

type UserHandler struct {
	*BaseHandler
	userService user.Service
	jwtManager  auth.AuthManager
}

func NewUserHandler(userService user.Service, logger *zap.Logger, jwtManager auth.AuthManager) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(logger),
		userService: userService,
		jwtManager:  jwtManager,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req *request.UserCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	result, err := h.userService.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "token required",
		})
	}
	_, err := h.userService.VerifyEmail(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Message": "email verified, account activated",
	})
}

func (h *UserHandler) GetById(c *fiber.Ctx) error {
	userId := c.Locals("userID")
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Error": "User id not found in context",
		})
	}
	id, ok := userId.(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Invalid user ID type in context",
		})
	}
	userProfile, err := h.userService.GetById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(200).JSON(userProfile)
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req request.UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Error": "invalid request",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "validation error",
		})
	}

	resp, err := h.userService.LoginUser(c.Context(), &req)
	if err != nil {
		if errors.Is(err, errorpkg.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.Status(200).JSON(resp)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == nil {
		return c.Status(401).JSON(fiber.Map{
			"Error": "user id not found in context",
		})
	}
	id, ok := userId.(int)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"Error": "invalid user id type in contenxt",
		})
	}

	userProfile, err := h.userService.GetUserProfile(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Error": "Failed get user profile",
		})
	}

	return c.Status(200).JSON(userProfile)
}

func (h *UserHandler) LogoutUser(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(500).JSON(fiber.Map{
			"Error": "Missing token",
		})
	}
	token = strings.TrimPrefix(token, "Bearer ")

	// verify lagi biar dapat exp
	claims, err := h.jwtManager.VerifyToken(c.Context(), token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"Error": "Invalid or expired token",
		})
	}

	exp, ok := claims["expired"].(float64)
	fmt.Printf("Claims saat logout: %+v\n", claims)

	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"Error": "Failed to logout",
		})
	}

	if err := h.userService.LogoutUser(c.Context(), token, int64(exp)); err != nil {
		h.logger.Error("failed to logout user", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{
			"Error": "Failed to logout",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Success": "Logout successfully",
	})
}
