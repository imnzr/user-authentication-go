package domain

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Repository interface {
	Create(ctx context.Context, tx *sql.Tx, user *CreateUserRequest) error
	FindUserByEmail(ctx context.Context, tx *sql.Tx, email string) error
	FindUserById(ctx context.Context, tx *sql.Tx, userId int) error
}
type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserById(ctx context.Context, userId string) (*User, error)
}
