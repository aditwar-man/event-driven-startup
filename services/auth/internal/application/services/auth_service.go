package services

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/application/ports"
	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/auth"
	sharedDomain "shared/pkg/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo       ports.UserRepository
	sessionManager *auth.SessionManager
	eventPublisher ports.EventPublisher
	tokenService   *auth.TokenService
}

func NewAuthService(
	userRepo ports.UserRepository,
	sessionManager *auth.SessionManager,
	eventPublisher ports.EventPublisher,
	tokenService *auth.TokenService,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		sessionManager: sessionManager,
		eventPublisher: eventPublisher,
		tokenService:   tokenService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

type RegisterResponse struct {
	User      *sharedDomain.User `json:"user"`
	TokenPair *auth.TokenPair    `json:"tokens"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User      *sharedDomain.User `json:"user"`
	TokenPair *auth.TokenPair    `json:"tokens"`
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// Check password strength
	if err := auth.PasswordStrengthCheck(req.Password); err != nil {
		return nil, err
	}

	// Check if user already exists
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, domain.ErrUserNotFound // Don't reveal that user exists
	}

	// Create password hash
	passwordHash, err := auth.GenerateSecurePasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user using shared domain
	user := &sharedDomain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Tier:         sharedDomain.UserTierFree,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens (without session for registration)
	tokenPair, err := s.tokenService.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	// Publish user registered event
	if err := s.eventPublisher.PublishUserRegistered(ctx, user); err != nil {
		// Log error but don't fail registration
	}

	return &RegisterResponse{
		User:      user,
		TokenPair: tokenPair,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*LoginResponse, *auth.Session, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, domain.ErrUserNotFound
	}

	// Validate password
	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, nil, domain.ErrUserNotFound
	}

	// Create session
	session, tokenPair, err := s.sessionManager.CreateSession(ctx, user.ID.String(), userAgent, ipAddress)
	if err != nil {
		return nil, nil, err
	}

	return &LoginResponse{
		User:      user,
		TokenPair: tokenPair,
	}, session, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	return s.sessionManager.RefreshSession(ctx, refreshToken)
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionManager.RevokeSession(ctx, sessionID)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !auth.VerifyPassword(currentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Check new password strength
	if err := auth.PasswordStrengthCheck(newPassword); err != nil {
		return err
	}

	// Generate new password hash
	newPasswordHash, err := auth.GenerateSecurePasswordHash(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Revoke all sessions for security
	if err := s.sessionManager.RevokeAllUserSessions(ctx, userID); err != nil {
		// Log error but don't fail password change
	}

	return nil
}

type UpgradeTierRequest struct {
	UserID string `json:"user_id"`
}

func (s *AuthService) UpgradeTier(ctx context.Context, req UpgradeTierRequest) error {
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	oldTier := user.Tier
	user.Tier = sharedDomain.UserTierPro
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Publish tier upgraded event
	if err := s.eventPublisher.PublishUserTierUpgraded(ctx, user.ID.String(), oldTier, user.Tier); err != nil {
		// Log error but don't fail operation
	}

	return nil
}
