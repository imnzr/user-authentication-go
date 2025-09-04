package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imnzr/user-authentication-go/internal/config"
)

type jwtManager struct {
	JWTSecretKey         []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(cfg config.Config) AuthManager {
	return &jwtManager{
		JWTSecretKey:         []byte(cfg.JSONWebToken.JWTSecretKey),
		accessTokenDuration:  cfg.JSONWebToken.AccessTokenDuration,
		refreshTokenDuration: cfg.JSONWebToken.RefreshTokenDuration,
	}
}

// VerifyToken implements AuthManager.
func (j *jwtManager) VerifyToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.JWTSecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GenerateAccessToken implements AuthManager.
func (j *jwtManager) GenerateAccessToken(ctx context.Context, userId int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"expired": time.Now().Add(j.accessTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWTSecretKey)
}

// GenerateRefreshToken implements AuthManager.
func (j *jwtManager) GenerateRefreshToken(ctx context.Context, userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"expired": time.Now().Add(j.refreshTokenDuration).Unix(),
		"type":    "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWTSecretKey)
}

// GenerateTokenVerif implements AuthManager.
func (j *jwtManager) GenerateTokenVerif(ctx context.Context, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"expired": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWTSecretKey)
}
