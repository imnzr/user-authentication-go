package user

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imnzr/user-authentication-go/pkg/request"
	"github.com/imnzr/user-authentication-go/pkg/response"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, userId int) (*User, error)

	// Verifify User Create
	ActivateByEmail(ctx context.Context, email string) error
}

type Service interface {
	Create(ctx context.Context, req *request.UserCreateRequest) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, userId int) (*User, error)
	GetUserProfile(ctx context.Context, userId int) (*response.UserProfileResponse, error)
	LoginUser(ctx context.Context, req *request.UserLoginRequest) (*response.TokenResponse, error)
	LogoutUser(ctx context.Context, token string, exp int64) error
	VerifyEmail(ctx context.Context, tokenString string) (jwt.MapClaims, error)
}

type Controller interface {
	Create(controller *fiber.Ctx) error
	GetByEmail(controller *fiber.Ctx) error
	GetById(controller *fiber.Ctx) error

	// Verify User Create
}
