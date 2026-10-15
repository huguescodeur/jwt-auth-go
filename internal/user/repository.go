package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u *User) (*User, error) {
	q := `
		INSERT INTO users (full_name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, full_name, email
	`
	err := r.db.QueryRowContext(ctx, q, u.FullName, u.Email, u.PasswordHash).Scan(&u.ID, &u.FullName, &u.Email)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, ErrAlreadyExist
			}

		}

		return nil, err
	}

	return u, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	q := `
		SELECT id, full_name, email, password_hash
		FROM users
		WHERE email = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	q := `
		SELECT id, full_name, email, password_hash
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &u, nil
}
