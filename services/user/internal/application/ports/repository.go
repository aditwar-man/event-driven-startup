package ports

import (
	"context"

	sharedDomain "shared/pkg/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *sharedDomain.User) error
	FindByID(ctx context.Context, id string) (*sharedDomain.User, error)
	FindByEmail(ctx context.Context, email string) (*sharedDomain.User, error)
	Update(ctx context.Context, user *sharedDomain.User) error
	UpdateQuotas(ctx context.Context, userID string, quotas sharedDomain.QuotaInfo) error
	ResetAllMonthlyQuotas(ctx context.Context) error
}

type EventPublisher interface {
	PublishUserQuotaUpdated(ctx context.Context, userID string, quotas interface{}) error
}

type EventSubscriber interface {
	SubscribeToUserEvents(ctx context.Context) error
}
