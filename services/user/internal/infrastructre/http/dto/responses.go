package dto

import (
	sharedDomain "shared/pkg/domain"
)

// UserResponse represents user data response
type UserResponse struct {
	Message string            `json:"message"`
	Data    *UserResponseData `json:"data"`
}

type UserResponseData struct {
	User      *sharedDomain.User      `json:"user"`
	QuotaInfo *sharedDomain.QuotaInfo `json:"quota_info"`
}

// QuotaUsageResponse represents quota usage response
type QuotaUsageResponse struct {
	Message string `json:"message"`
}

// TierUpgradeResponse represents tier upgrade response
type TierUpgradeResponse struct {
	Message string `json:"message"`
}

// QuotaCheckResponse represents quota check response
type QuotaCheckResponse struct {
	HasQuota bool `json:"has_quota"`
}

// AdminQuotaResetResponse represents admin quota reset response
type AdminQuotaResetResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}
