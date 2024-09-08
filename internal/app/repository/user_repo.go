package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/internal/models"
	"log"
)

type UserRepo struct {
	db *sql.DB
	// You can insert a caching service here
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`
	_, err := u.db.ExecContext(ctx, query, email, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE user_id = $1`
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	err := u.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	return exists, nil
}

func (u *UserRepo) CheckIfUserExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE user_id = $1)`
	err := u.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
