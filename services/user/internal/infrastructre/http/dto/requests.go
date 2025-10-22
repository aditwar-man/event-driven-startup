package dto

// UseAIDescriptionQuotaRequest represents AI description quota usage request
type UseAIDescriptionQuotaRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// UseAIVideoQuotaRequest represents AI video quota usage request
type UseAIVideoQuotaRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// UseAutoPostingQuotaRequest represents auto posting quota usage request
type UseAutoPostingQuotaRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// UpgradeToProRequest represents tier upgrade request
type UpgradeToProRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
