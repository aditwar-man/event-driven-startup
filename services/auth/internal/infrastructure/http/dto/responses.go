package dto

import (
	"auth-service/internal/infrastructure/auth"
	sharedDomain "shared/pkg/domain"
	"time"
)

// RegisterResponse represents registration response
type RegisterResponse struct {
	Message string                `json:"message"`
	Data    *RegisterResponseData `json:"data"`
}

type RegisterResponseData struct {
	User      *sharedDomain.User `json:"user"`
	TokenPair *auth.TokenPair    `json:"tokens"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Message string             `json:"message"`
	Data    *LoginResponseData `json:"data"`
	Session *SessionInfo       `json:"session,omitempty"`
}

type LoginResponseData struct {
	User      *sharedDomain.User `json:"user"`
	TokenPair *auth.TokenPair    `json:"tokens"`
}

type SessionInfo struct {
	ID        string    `json:"id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshTokenResponse represents token refresh response
type RefreshTokenResponse struct {
	Message string          `json:"message"`
	Data    *auth.TokenPair `json:"data"`
}

// ProfileResponse represents user profile response
type ProfileResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Tier      string `json:"tier"`
	CreatedAt string `json:"created_at"`
}

// SessionResponse represents session information
type SessionResponse struct {
	ID        string    `json:"id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionsListResponse represents list of sessions response
type SessionsListResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
}

// SuccessResponse represents generic success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
