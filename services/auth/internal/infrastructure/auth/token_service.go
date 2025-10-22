package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	secretKey       []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

type TokenConfig struct {
	SecretKey       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

func NewTokenService(cfg TokenConfig) *TokenService {
	return &TokenService{
		secretKey:       []byte(cfg.SecretKey),
		accessTokenExp:  cfg.AccessTokenExp,
		refreshTokenExp: cfg.RefreshTokenExp,
	}
}

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair creates both access and refresh tokens
func (s *TokenService) GenerateTokenPair(userID, email string) (*TokenPair, error) {
	sessionID, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, err := s.generateAccessToken(userID, email, sessionID)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(userID, email, sessionID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenExp.Seconds()),
	}, nil
}

func (s *TokenService) generateAccessToken(userID, email, sessionID string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "smm-platform",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) generateRefreshToken(userID, email, sessionID string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "smm-platform",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, "access")
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, "refresh")
}

func (s *TokenService) validateToken(tokenString, tokenType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation for token type
		if tokenType == "access" && time.Until(claims.ExpiresAt.Time) > s.refreshTokenExp {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateSecurePasswordHash creates a bcrypt hash with salt
func GenerateSecurePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecureToken creates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// PasswordStrengthCheck validates password strength
func PasswordStrengthCheck(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for uppercase, lowercase, number, and special character
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char >= 33 && char <= 126:
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, number, and special characters")
	}

	return nil
}
