package repository

import (
	"context"
	"fmt"
	"time"

	"oat431/big-shwein-api/internal/model"
	"oat431/big-shwein-api/internal/payload/request"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AuthRepository is the PostgreSQL implementation for auth data access.
type AuthRepository struct {
	db *sqlx.DB
}

// NewAuthRepository creates a new AuthRepository backed by PostgreSQL.
func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Register(ctx context.Context, request request.RegisterRequest) (*model.Auth, error) {
	query := `INSERT INTO tb_auth (
				id,
				created_at,
				updated_at,
				deleted_at,
				username,
				email,
				"password",
				is_verified
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	id := uuid.New()
	currentTime := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, currentTime, currentTime, nil, request.Username, request.Email, request.Password, false)
	if err != nil {
		return nil, fmt.Errorf("register auth: %w", err)
	}
	return &model.Auth{
		BaseEntity: model.BaseEntity{
			ID:        id,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
		Username:   request.Username,
		Email:      request.Email,
		Password:   request.Password,
		IsVerified: false,
	}, nil
}

func (r *AuthRepository) GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error) {
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			deleted_at,
			username,
			email,
			"password",
			is_verified
		FROM
			tb_auth
		WHERE
			username = $1`
	var auth model.Auth
	if err := r.db.GetContext(ctx, &auth, query, username); err != nil {
		return nil, fmt.Errorf("get auth by username %q: %w", username, err)
	}
	return &auth, nil
}

func (r *AuthRepository) GetAuthByID(ctx context.Context, id uuid.UUID) (*model.Auth, error) {
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			deleted_at,
			username,
			email,
			"password",
			is_verified
		FROM
			tb_auth
		WHERE
			id = $1`
	var auth model.Auth
	if err := r.db.GetContext(ctx, &auth, query, id); err != nil {
		return nil, fmt.Errorf("get auth by id %s: %w", id, err)
	}
	return &auth, nil
}

func (r *AuthRepository) GetAuthByEmail(ctx context.Context, email string) (*model.Auth, error) {
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			deleted_at,
			username,
			email,
			"password",
			is_verified
		FROM
			tb_auth
		WHERE
			email = $1`
	var auth model.Auth
	if err := r.db.GetContext(ctx, &auth, query, email); err != nil {
		return nil, fmt.Errorf("get auth by email %q: %w", email, err)
	}
	return &auth, nil
}

func (r *AuthRepository) MarkAsVerified(ctx context.Context, authID uuid.UUID) error {
	query := `UPDATE tb_auth SET is_verified = true, updated_at = NOW() WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, authID); err != nil {
		return fmt.Errorf("mark auth %s as verified: %w", authID, err)
	}
	return nil
}
