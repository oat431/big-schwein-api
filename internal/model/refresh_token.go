package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a stored refresh token for session management.
type RefreshToken struct {
	BaseEntity

	AuthID uuid.UUID `db:"auth_id" json:"auth_id"`

	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Revoked   bool      `db:"revoked" json:"revoked"`
}
