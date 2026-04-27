package model

import (
	"time"

	"oat431/big-shwein-api/pkg/common"

	"github.com/google/uuid"
)

type RefreshToken struct {
	common.BaseEntity

	AuthID uuid.UUID `db:"auth_id" json:"auth_id"`

	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Revoked   bool      `db:"revoked" json:"revoked"`
}
