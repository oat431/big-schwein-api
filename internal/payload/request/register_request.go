package request

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Email    string `json:"email" validate:"required,email"`
}
