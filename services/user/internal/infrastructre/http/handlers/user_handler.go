package handlers

import (
	"net/http"

	"user-service/internal/application/services"
	"user-service/internal/domain"
	"user-service/internal/infrastructre/http/dto"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details and quota information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	response, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	resp := dto.UserResponse{
		Message: "User retrieved successfully",
		Data: &dto.UserResponseData{
			User:      response.User,
			QuotaInfo: &response.QuotaInfo,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Get user details and quota information by email
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	response, err := h.userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	resp := dto.UserResponse{
		Message: "User retrieved successfully",
		Data: &dto.UserResponseData{
			User:      response.User,
			QuotaInfo: &response.QuotaInfo,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// UseAIDescriptionQuota godoc
// @Summary Use AI description quota
// @Description Use one AI description generation quota
// @Tags quotas
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.QuotaUsageResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/use-ai-description [post]
func (h *UserHandler) UseAIDescriptionQuota(c *gin.Context) {
	userID := c.Param("id")

	req := services.UseAIDescriptionQuotaRequest{
		UserID: userID,
	}

	if err := h.userService.UseAIDescriptionQuota(c.Request.Context(), req); err != nil {
		switch err {
		case domain.ErrQuotaExceeded:
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "AI description quota exceeded"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.QuotaUsageResponse{Message: "AI description quota used successfully"})
}

// UseAIVideoQuota godoc
// @Summary Use AI video quota
// @Description Use one AI video generation quota
// @Tags quotas
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.QuotaUsageResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/use-ai-video [post]
func (h *UserHandler) UseAIVideoQuota(c *gin.Context) {
	userID := c.Param("id")

	req := services.UseAIVideoQuotaRequest{
		UserID: userID,
	}

	if err := h.userService.UseAIVideoQuota(c.Request.Context(), req); err != nil {
		switch err {
		case domain.ErrQuotaExceeded:
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "AI video quota exceeded"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.QuotaUsageResponse{Message: "AI video quota used successfully"})
}

// UseAutoPostingQuota godoc
// @Summary Use auto posting quota
// @Description Use one auto posting quota
// @Tags quotas
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.QuotaUsageResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/use-auto-posting [post]
func (h *UserHandler) UseAutoPostingQuota(c *gin.Context) {
	userID := c.Param("id")

	req := services.UseAutoPostingQuotaRequest{
		UserID: userID,
	}

	if err := h.userService.UseAutoPostingQuota(c.Request.Context(), req); err != nil {
		switch err {
		case domain.ErrQuotaExceeded:
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Auto posting quota exceeded"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.QuotaUsageResponse{Message: "Auto posting quota used successfully"})
}

// UpgradeToPro godoc
// @Summary Upgrade to PRO tier
// @Description Upgrade user account to PRO tier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.TierUpgradeResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/upgrade-pro [post]
func (h *UserHandler) UpgradeToPro(c *gin.Context) {
	userID := c.Param("id")

	req := services.UpgradeToProRequest{
		UserID: userID,
	}

	if err := h.userService.UpgradeToPro(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to upgrade tier"})
		return
	}

	c.JSON(http.StatusOK, dto.TierUpgradeResponse{Message: "Tier upgraded to PRO successfully"})
}

// CheckAIDescriptionQuota godoc
// @Summary Check AI description quota
// @Description Check if user has AI description quota available
// @Tags quotas
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.QuotaCheckResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id}/check-ai-description-quota [get]
func (h *UserHandler) CheckAIDescriptionQuota(c *gin.Context) {
	userID := c.Param("id")

	hasQuota, err := h.userService.CheckAIDescriptionQuota(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.QuotaCheckResponse{HasQuota: hasQuota})
}

// CheckAIVideoQuota godoc
// @Summary Check AI video quota
// @Description Check if user has AI video quota available
// @Tags quotas
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.QuotaCheckResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id}/check-ai-video-quota [get]
func (h *UserHandler) CheckAIVideoQuota(c *gin.Context) {
	userID := c.Param("id")

	hasQuota, err := h.userService.CheckAIVideoQuota(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.QuotaCheckResponse{HasQuota: hasQuota})
}

// CheckAutoPostingQuota godoc
// @Summary Check auto posting quota
// @Description Check if user has auto posting quota available
// @Tags quotas
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.QuotaCheckResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id}/check-auto-posting-quota [get]
func (h *UserHandler) CheckAutoPostingQuota(c *gin.Context) {
	userID := c.Param("id")

	hasQuota, err := h.userService.CheckAutoPostingQuota(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.QuotaCheckResponse{HasQuota: hasQuota})
}

// ResetMonthlyQuotas godoc
// @Summary Reset monthly quotas (Admin)
// @Description Reset all users' monthly quotas (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} dto.AdminQuotaResetResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/reset-monthly-quotas [post]
func (h *UserHandler) ResetMonthlyQuotas(c *gin.Context) {
	if err := h.userService.ResetMonthlyQuotas(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to reset quotas"})
		return
	}

	c.JSON(http.StatusOK, dto.AdminQuotaResetResponse{Message: "Monthly quotas reset successfully"})
}
