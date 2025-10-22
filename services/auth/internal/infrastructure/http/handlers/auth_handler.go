package handlers

import (
	"net/http"

	"auth-service/internal/application/services"
	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/auth"
	"auth-service/internal/infrastructure/http/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService    *services.AuthService
	sessionManager *auth.SessionManager
}

func NewAuthHandler(authService *services.AuthService, sessionManager *auth.SessionManager) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionManager: sessionManager,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := services.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	response, err := h.authService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		switch err.(type) {
		case *domain.DomainError:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal server error",
				Message: err.Error(),
			})
		}
		return
	}

	resp := dto.RegisterResponse{
		Message: "User registered successfully",
		Data: &dto.RegisterResponseData{
			User:      response.User,
			TokenPair: response.TokenPair,
		},
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Extract user agent and IP address for session tracking
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	response, session, err := h.authService.Login(c.Request.Context(), serviceReq, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Set session cookie (optional)
	c.SetCookie("refresh_token", response.TokenPair.RefreshToken, 7*24*3600, "/", "", false, true)

	resp := dto.LoginResponse{
		Message: "Login successful",
		Data: &dto.LoginResponseData{
			User:      response.User,
			TokenPair: response.TokenPair,
		},
	}

	if session != nil {
		resp.Session = &dto.SessionInfo{
			ID:        session.ID,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			CreatedAt: session.CreatedAt,
			ExpiresAt: session.ExpiresAt,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid refresh token"})
		return
	}

	resp := dto.RefreshTokenResponse{
		Message: "Token refreshed successfully",
		Data:    tokenPair,
	}

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary User logout
// @Description Invalidate user session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Logout request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// If no session ID provided, try to get from context (for current session)
	sessionID := req.SessionID
	if sessionID == "" {
		if val, exists := c.Get("session_id"); exists {
			sessionID = val.(string)
		}
	}

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Session ID required"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to logout"})
		return
	}

	// Clear cookies
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Logout successful"})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ProfileResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Not authenticated"})
		return
	}

	userEmail, _ := c.Get("user_email")
	// In a real implementation, you would fetch complete user profile from service

	resp := dto.ProfileResponse{
		UserID:   userID.(string),
		Email:    userEmail.(string),
		FullName: "User Full Name", // Fetch from service
		Tier:     "free",           // Fetch from service
	}

	c.JSON(http.StatusOK, resp)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user's password
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Not authenticated"})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID.(string), req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Password changed successfully"})
}

// GetSessions godoc
// @Summary Get user sessions
// @Description Get all active sessions for current user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SessionsListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions [get]
func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Not authenticated"})
		return
	}

	sessions, err := h.sessionManager.ListByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch sessions"})
		return
	}

	var sessionResponses []*dto.SessionResponse
	for _, session := range sessions {
		sessionResponses = append(sessionResponses, &dto.SessionResponse{
			ID:        session.ID,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			CreatedAt: session.CreatedAt,
			ExpiresAt: session.ExpiresAt,
		})
	}

	resp := dto.SessionsListResponse{Sessions: sessionResponses}
	c.JSON(http.StatusOK, resp)
}

// RevokeSession godoc
// @Summary Revoke a session
// @Description Revoke a specific user session
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RevokeSessionRequest true "Revoke session request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /sessions/revoke [post]
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	var req dto.RevokeSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.sessionManager.RevokeSession(c.Request.Context(), req.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Session revoked successfully"})
}

// RevokeAllSessions godoc
// @Summary Revoke all sessions
// @Description Revoke all active sessions for current user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/revoke-all [post]
func (h *AuthHandler) RevokeAllSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Not authenticated"})
		return
	}

	if err := h.sessionManager.RevokeAllUserSessions(c.Request.Context(), userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to revoke sessions"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "All sessions revoked successfully"})
}

// UpgradeTier godoc
// @Summary Upgrade to PRO tier
// @Description Upgrade user account to PRO tier
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /upgrade-tier [post]
func (h *AuthHandler) UpgradeTier(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Not authenticated"})
		return
	}

	req := services.UpgradeTierRequest{
		UserID: userID.(string),
	}

	if err := h.authService.UpgradeTier(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to upgrade tier"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Tier upgraded successfully"})
}

// ListUsers is an admin-only endpoint
func (h *AuthHandler) ListUsers(c *gin.Context) {
	// TODO: Implement admin user listing
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse{Error: "Admin user listing not implemented yet"})
}
