package model

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerifyToken represents a stored email verification token.
type EmailVerifyToken struct {
	BaseEntity

	AuthID    uuid.UUID `db:"auth_id" json:"auth_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}
