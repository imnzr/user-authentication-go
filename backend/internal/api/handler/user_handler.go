package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/user-authentication-go/internal/domain/user"
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
