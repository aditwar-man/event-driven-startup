package middleware

import (
	"net/http"
	"strings"

	"auth-service/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService   *auth.TokenService
	sessionManager *auth.SessionManager
}

func NewAuthMiddleware(tokenService *auth.TokenService, sessionManager *auth.SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService:   tokenService,
		sessionManager: sessionManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c.Request)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := m.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Validate session
		session, err := m.sessionManager.ValidateSession(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("session_id", claims.SessionID)
		c.Set("session", session)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check authentication
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// TODO: Check user role from database
		// For now, we'll implement basic role checking
		userID, _ := c.Get("user_id")

		// Example role check - in real implementation, fetch user roles from database
		if !hasRole(userID.(string), role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractToken(r *http.Request) string {
	// Check Authorization header
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
		return bearerToken[7:]
	}

	// Check query parameter
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}

	return ""
}

func hasRole(userID string, role string) bool {
	// TODO: Implement proper role checking from database
	// This is a placeholder implementation
	return role == "user" // Basic implementation
}
