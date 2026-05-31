package authdto

import "time"

type LoginResponse struct {
	Token     string       `json:"token"`
	TokenType string       `json:"token_type"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
