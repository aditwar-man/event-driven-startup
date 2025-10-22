package ports

import (
	"context"

	"auth-service/internal/infrastructure/auth"
	sharedDomain "shared/pkg/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *sharedDomain.User) error
	FindByID(ctx context.Context, id string) (*sharedDomain.User, error)
	FindByEmail(ctx context.Context, email string) (*sharedDomain.User, error)
	Update(ctx context.Context, user *sharedDomain.User) error
}

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, user interface{}) error
	PublishUserTierUpgraded(ctx context.Context, userID string, oldTier, newTier interface{}) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *auth.Session) error
	GetByID(ctx context.Context, sessionID string) (*auth.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteByUserID(ctx context.Context, userID string) error
	ListByUserID(ctx context.Context, userID string) ([]*auth.Session, error)
}
