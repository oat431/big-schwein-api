package repository

import (
	"context"
	"fmt"

	"oat431/big-shwein-api/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// EmailVerifyTokenRepository is the PostgreSQL implementation for email verification token data access.
type EmailVerifyTokenRepository struct {
	db *sqlx.DB
}

// NewEmailVerifyTokenRepository creates a new EmailVerifyTokenRepository backed by PostgreSQL.
func NewEmailVerifyTokenRepository(db *sqlx.DB) *EmailVerifyTokenRepository {
	return &EmailVerifyTokenRepository{db: db}
}

func (r *EmailVerifyTokenRepository) Save(ctx context.Context, token model.EmailVerifyToken) error {
	query := `INSERT INTO tb_verify_tokens (
				id, created_at, updated_at, deleted_at, auth_id, token, expires_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.CreatedAt,
		token.UpdatedAt,
		token.DeletedAt,
		token.AuthID,
		token.Token,
		token.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("save email verify token: %w", err)
	}
	return nil
}

func (r *EmailVerifyTokenRepository) FindByToken(ctx context.Context, token string) (*model.EmailVerifyToken, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, auth_id, token, expires_at
			  FROM tb_verify_tokens
			  WHERE token = $1 AND deleted_at IS NULL`
	var result model.EmailVerifyToken
	if err := r.db.GetContext(ctx, &result, query, token); err != nil {
		return nil, fmt.Errorf("find email verify token: %w", err)
	}
	return &result, nil
}

func (r *EmailVerifyTokenRepository) DeleteByAuthID(ctx context.Context, authID uuid.UUID) error {
	query := `DELETE FROM tb_verify_tokens WHERE auth_id = $1`
	if _, err := r.db.ExecContext(ctx, query, authID); err != nil {
		return fmt.Errorf("delete verify tokens for auth %s: %w", authID, err)
	}
	return nil
}
