package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AuthManager interface {
	// Verifify User Create
	VerifyToken(ctx context.Context, tokenString string) (jwt.MapClaims, error)
	GenerateTokenVerif(ctx context.Context, email string) (string, error)
	GenerateAccessToken(ctx context.Context, userId int, email string) (string, error)
	GenerateRefreshToken(ctx context.Context, userId int) (string, error)
}
