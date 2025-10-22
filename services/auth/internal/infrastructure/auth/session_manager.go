package auth

import (
	"context"
	"errors"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, sessionID string) (*Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteByUserID(ctx context.Context, userID string) error
	ListByUserID(ctx context.Context, userID string) ([]*Session, error)
}

type SessionManager struct {
	repo          SessionRepository
	tokenService  *TokenService
	sessionExpiry time.Duration
}

func NewSessionManager(repo SessionRepository, tokenService *TokenService, sessionExpiry time.Duration) *SessionManager {
	return &SessionManager{
		repo:          repo,
		tokenService:  tokenService,
		sessionExpiry: sessionExpiry,
	}
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID, userAgent, ipAddress string) (*Session, *TokenPair, error) {
	// Generate tokens
	tokenPair, err := sm.tokenService.GenerateTokenPair(userID, "")
	if err != nil {
		return nil, nil, err
	}

	// Create session
	session := &Session{
		ID:           tokenPair.AccessToken, // Using access token as session ID for simplicity
		UserID:       userID,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(sm.sessionExpiry),
		CreatedAt:    time.Now(),
	}

	if err := sm.repo.Create(ctx, session); err != nil {
		return nil, nil, err
	}

	return session, tokenPair, nil
}

func (sm *SessionManager) ValidateSession(ctx context.Context, sessionID string) (*Session, error) {
	session, err := sm.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		sm.repo.Delete(ctx, sessionID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (sm *SessionManager) RefreshSession(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := sm.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return sm.tokenService.GenerateTokenPair(claims.UserID, claims.Email)
}

func (sm *SessionManager) RevokeSession(ctx context.Context, sessionID string) error {
	return sm.repo.Delete(ctx, sessionID)
}

func (sm *SessionManager) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return sm.repo.DeleteByUserID(ctx, userID)
}

// ListByUserID returns all active sessions for a user
func (sm *SessionManager) ListByUserID(ctx context.Context, userID string) ([]*Session, error) {
	return sm.repo.ListByUserID(ctx, userID)
}
