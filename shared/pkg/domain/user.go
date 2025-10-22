package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserTier string

const (
	UserTierFree UserTier = "free"
	UserTierPro  UserTier = "pro"
)

// User represents the core user entity shared across services
type User struct {
	ID                      uuid.UUID `json:"id"`
	Email                   string    `json:"email"`
	PasswordHash            string    `json:"-"` // Never expose in JSON
	FullName                string    `json:"full_name"`
	Tier                    UserTier  `json:"tier"`
	AIDescriptionQuotaUsed  int       `json:"ai_description_quota_used"`
	AIDescriptionQuotaLimit int       `json:"ai_description_quota_limit"`
	AIVideoQuotaUsed        int       `json:"ai_video_quota_used"`
	AIVideoQuotaLimit       int       `json:"ai_video_quota_limit"`
	AutoPostingQuotaUsed    int       `json:"auto_posting_quota_used"`
	AutoPostingQuotaLimit   int       `json:"auto_posting_quota_limit"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// Quota represents usage limits
type Quota struct {
	Used  int `json:"used"`
	Limit int `json:"limit"`
}

// QuotaInfo contains all quota information
type QuotaInfo struct {
	AIDescription Quota `json:"ai_description"`
	AIVideo       Quota `json:"ai_video"`
	AutoPosting   Quota `json:"auto_posting"`
}

// Helper methods for quota management
func (u *User) CanGenerateAIDescription() bool {
	return u.AIDescriptionQuotaUsed < u.AIDescriptionQuotaLimit
}

func (u *User) CanGenerateAIVideo() bool {
	return u.AIVideoQuotaUsed < u.AIVideoQuotaLimit
}

func (u *User) CanAutoPost() bool {
	return u.AutoPostingQuotaUsed < u.AutoPostingQuotaLimit
}

func (u *User) UseAIDescriptionQuota() error {
	if !u.CanGenerateAIDescription() {
		return ErrQuotaExceeded
	}
	u.AIDescriptionQuotaUsed++
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) UseAIVideoQuota() error {
	if !u.CanGenerateAIVideo() {
		return ErrQuotaExceeded
	}
	u.AIVideoQuotaUsed++
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) UseAutoPostingQuota() error {
	if !u.CanAutoPost() {
		return ErrQuotaExceeded
	}
	u.AutoPostingQuotaUsed++
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) UpgradeToPro() {
	u.Tier = UserTierPro
	u.AIDescriptionQuotaLimit = 100
	u.AIVideoQuotaLimit = 10
	u.AutoPostingQuotaLimit = 1000
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) ResetMonthlyQuotas() {
	u.AIDescriptionQuotaUsed = 0
	u.AIVideoQuotaUsed = 0
	u.AutoPostingQuotaUsed = 0
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) GetQuotaInfo() QuotaInfo {
	return QuotaInfo{
		AIDescription: Quota{
			Used:  u.AIDescriptionQuotaUsed,
			Limit: u.AIDescriptionQuotaLimit,
		},
		AIVideo: Quota{
			Used:  u.AIVideoQuotaUsed,
			Limit: u.AIVideoQuotaLimit,
		},
		AutoPosting: Quota{
			Used:  u.AutoPostingQuotaUsed,
			Limit: u.AutoPostingQuotaLimit,
		},
	}
}

// Domain errors
var (
	ErrQuotaExceeded = NewDomainError("quota exceeded")
)

type DomainError struct {
	message string
}

func NewDomainError(message string) *DomainError {
	return &DomainError{message: message}
}

func (d *DomainError) Error() string {
	return d.message
}
