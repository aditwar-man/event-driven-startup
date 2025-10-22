package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sharedDomain "shared/pkg/domain"
	sharedEvents "shared/pkg/events"
	"user-service/internal/application/services"
	"user-service/internal/infrastructre/persistence"

	"github.com/google/uuid"
)

type UniversalEventSubscriber struct {
	eventBus    sharedEvents.EventBus
	userService *services.UserService
	userRepo    *persistence.PostgresUserRepository
}

func NewUniversalEventSubscriber(
	eventBus sharedEvents.EventBus,
	userService *services.UserService,
	userRepo *persistence.PostgresUserRepository,
) *UniversalEventSubscriber {
	return &UniversalEventSubscriber{
		eventBus:    eventBus,
		userService: userService,
		userRepo:    userRepo,
	}
}

func (u *UniversalEventSubscriber) SubscribeToUserEvents(ctx context.Context) error {
	subscriber := sharedEvents.NewUniversalEventSubscriber(u.eventBus)
	return subscriber.SubscribeToUserEvents(ctx, u.handleUserEvent)
}

func (u *UniversalEventSubscriber) handleUserEvent(ctx context.Context, event *sharedEvents.Event) error {
	switch event.Type {
	case sharedEvents.UserRegisteredEvent:
		return u.handleUserRegistered(ctx, event)
	case sharedEvents.UserTierUpgradedEvent:
		return u.handleUserTierUpgraded(ctx, event)
	default:
		fmt.Printf("Unknown event type: %s\n", event.Type)
		return nil
	}
}

func (u *UniversalEventSubscriber) handleUserRegistered(ctx context.Context, event *sharedEvents.Event) error {
	var data sharedEvents.UserRegisteredData
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal user registered data: %w", err)
	}

	userID, err := uuid.Parse(data.UserID)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	var tier sharedDomain.UserTier
	switch data.Tier {
	case "free":
		tier = sharedDomain.UserTierFree
	case "pro":
		tier = sharedDomain.UserTierPro
	default:
		tier = sharedDomain.UserTierFree
	}

	// Create user in user service database using shared domain
	user := &sharedDomain.User{
		ID:        userID,
		Email:     data.Email,
		FullName:  data.FullName,
		Tier:      tier,
		CreatedAt: parseTime(data.CreatedAt),
		UpdatedAt: parseTime(data.CreatedAt),
	}

	// Set default quotas based on tier
	if user.Tier == sharedDomain.UserTierFree {
		user.AIDescriptionQuotaLimit = 5
		user.AIVideoQuotaLimit = 0
		user.AutoPostingQuotaLimit = 5
	} else {
		user.AIDescriptionQuotaLimit = 100
		user.AIVideoQuotaLimit = 10
		user.AutoPostingQuotaLimit = 1000
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user in user service: %w", err)
	}

	fmt.Printf("User created in user service: %s (%s)\n", user.ID, user.Email)
	return nil
}

func (u *UniversalEventSubscriber) handleUserTierUpgraded(ctx context.Context, event *sharedEvents.Event) error {
	var data sharedEvents.UserTierUpgradedData
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal user tier upgraded data: %w", err)
	}

	req := services.UpgradeToProRequest{
		UserID: data.UserID,
	}

	if err := u.userService.UpgradeToPro(ctx, req); err != nil {
		return fmt.Errorf("failed to upgrade user tier: %w", err)
	}

	fmt.Printf("User tier upgraded: %s from %s to %s\n", data.UserID, data.OldTier, data.NewTier)
	return nil
}

// Helper function to parse time strings
func parseTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Now().UTC()
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}
