package response

import "github.com/gofiber/fiber/v2"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// Response error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

// SendSuccess sends a success response
func SendSuccess(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	return c.Status(statusCode).JSON(response)
}

// SendError sends a error response
func SendError(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	}
	return c.Status(statusCode).JSON(response)
}

// Send Token Response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
