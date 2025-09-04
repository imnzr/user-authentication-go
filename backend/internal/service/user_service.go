package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imnzr/user-authentication-go/internal/database"
	"github.com/imnzr/user-authentication-go/internal/domain/user"
	"github.com/imnzr/user-authentication-go/pkg/auth"
	"github.com/imnzr/user-authentication-go/pkg/request"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo    user.Repository
	txManager   database.TxManager
	authManager auth.AuthManager
}

func NewUserService(userRepo user.Repository, txManager database.TxManager, authManager auth.AuthManager) user.Service {
	return &service{
		userRepo:    userRepo,
		txManager:   txManager,
		authManager: authManager,
	}
}

func (s *service) ValidateCreateUser(req request.UserCreateRequest) error {
	if req.Username == "" {
		return fmt.Errorf("invalid username")
	}
	if req.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if req.Password == "" {
		return fmt.Errorf("invalid password")
	}

	return nil
}

// Create implements user.Service.
func (s *service) Create(ctx context.Context, req *request.UserCreateRequest) (*user.User, error) {
	if err := s.ValidateCreateUser(*req); err != nil {
		return nil, err
	}

	var createdUser *user.User

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Check if user already exists
		existing, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if existing != nil {
			return fmt.Errorf("user already exists")
		}

		// Hash a password
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash user password: %w", err)
		}

		// Create user
		newUser := &user.User{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashPassword),
			Status:   "pending",
		}

		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		token, err := s.authManager.GenerateTokenVerif(ctx, newUser.Email)
		if err != nil {
			return fmt.Errorf("error generate token: %w", err)
		}
		verifyLink := fmt.Sprintf("http://localhost:8080/api/v1/auth/verify/%s", token)
		fmt.Println("Verification link: ", verifyLink)

		createdUser = newUser

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// GetByEmail implements user.Service.
func (s *service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetById implements user.Service.
func (s *service) GetById(ctx context.Context, userId int) (*user.User, error) {
	return s.userRepo.GetById(ctx, userId)
}

// VerifyEmail implements user.Service.
func (s *service) VerifyEmail(ctx context.Context, tokenString string) (jwt.MapClaims, error) {

	claims, err := s.authManager.VerifyToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email in token")
	}

	if err := s.userRepo.ActivateByEmail(ctx, email); err != nil {
		return nil, err
	}

	return claims, nil

}
