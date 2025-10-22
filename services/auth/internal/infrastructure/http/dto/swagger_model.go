package dto

import (
	"time"
)

// ============ SWAGGER MODELS ============
// These models are for Swagger documentation only

// User represents user model for Swagger
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Tier      string    `json:"tier"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TokenPair represents token pair for Swagger
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// QuotaInfo represents quota information for Swagger
type QuotaInfo struct {
	AIDescription Quota `json:"ai_description"`
	AIVideo       Quota `json:"ai_video"`
	AutoPosting   Quota `json:"auto_posting"`
}

// Quota represents quota for Swagger
type Quota struct {
	Used  int `json:"used"`
	Limit int `json:"limit"`
}
