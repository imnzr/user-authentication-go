package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imnzr/user-authentication-go/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{
		db: db,
	}
}

// Helper to extract transaction from context
func getTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value("db_tx").(*sql.Tx)
	return tx, ok
}

// Create implements user.Repository.
func (u *userRepository) Create(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO users(username, email, password, created_at, updated_at)
		VALUES (?,?,?,NOW(),NOW())
	`
	var result sql.Result
	var err error

	// Transaction on context ? using this
	if tx, ok := getTxFromContext(ctx); ok {
		result, err = tx.ExecContext(ctx, query,
			user.Username,
			user.Email,
			user.Password,
		)
		// If don't have transaction on context
	} else {
		result, err = u.db.ExecContext(ctx, query,
			user.Username,
			user.Email,
			user.Password,
		)
	}
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	user.Id = int(id)
	return nil

}

// GetByEmail implements user.Repository.
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, username, email, password FROM users WHERE email = ?
	`
	user := &user.User{}
	var err error

	if tx, ok := getTxFromContext(ctx); ok {
		err = tx.QueryRowContext(ctx, query, email).Scan(
			&user.Id, &user.Username, &user.Email, &user.Password,
		)
	} else {
		err = u.db.QueryRowContext(ctx, query, email).Scan(
			&user.Id, &user.Username, &user.Email, &user.Password,
		)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetById implements user.Repository.
func (u *userRepository) GetById(ctx context.Context, userId int) (*user.User, error) {
	query := `
		SELECT id, username, email, password FROM users WHERE id = ?
	`

	user := &user.User{}
	var err error

	if tx, ok := getTxFromContext(ctx); ok {
		err = tx.QueryRowContext(ctx, query, userId).Scan(
			&user.Id, &user.Username, &user.Email, &user.Password,
		)
	} else {
		err = u.db.QueryRowContext(ctx, query, userId).Scan(
			&user.Id, &user.Username, &user.Email, &user.Password,
		)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// ActivateByEmail implements user.Repository.
func (u *userRepository) ActivateByEmail(ctx context.Context, email string) error {
	query := "UPDATE users SET status='active' WHERE email=?"
	res, err := u.db.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no user found or already exist")
	}

	return nil
}

// ResetPassword implements user.Repository.
func (u *userRepository) ResetPassword(ctx context.Context, email string) error {
	query := "UPDATE users SET password = ? WHERE email = ?"
	res, err := u.db.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no user found or already exist")
	}

	return nil
}
