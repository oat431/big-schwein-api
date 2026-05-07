package repository

import (
	"context"
	"fmt"

	"oat431/big-shwein-api/internal/model"

	"github.com/jmoiron/sqlx"
)

// RefreshTokenRepository is the PostgreSQL implementation for refresh token data access.
type RefreshTokenRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository backed by PostgreSQL.
func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, refreshToken model.RefreshToken) error {
	query := `INSERT INTO tb_refresh_tokens (
				id, created_at, updated_at, deleted_at, auth_id, token, expires_at, revoked
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		refreshToken.ID,
		refreshToken.CreatedAt,
		refreshToken.UpdatedAt,
		refreshToken.DeletedAt,
		refreshToken.AuthID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.Revoked)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE tb_refresh_tokens SET revoked = true, updated_at = NOW() WHERE token = $1`
	if _, err := r.db.ExecContext(ctx, query, token); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, auth_id, token, expires_at, revoked FROM tb_refresh_tokens WHERE token = $1 AND revoked = false AND deleted_at IS NULL`
	var refreshToken model.RefreshToken
	if err := r.db.GetContext(ctx, &refreshToken, query, token); err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return &refreshToken, nil
}
