package model

// Auth represents a user authentication record.
type Auth struct {
	BaseEntity

	Username   string `db:"username" json:"username"`
	Email      string `db:"email" json:"email"`
	Password   string `db:"password" json:"-"`
	IsVerified bool   `db:"is_verified" json:"is_verified"`
}
