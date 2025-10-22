package events

import (
	"context"
	"log"
	"shared/pkg/domain"
	"time"
)

// UniversalEventPublisher works with any EventBus implementation
type UniversalEventPublisher struct {
	eventBus EventBus
}

// NewUniversalEventPublisher creates a new universal event publisher
func NewUniversalEventPublisher(eventBus EventBus) *UniversalEventPublisher {
	return &UniversalEventPublisher{
		eventBus: eventBus,
	}
}

// PublishUserRegistered publishes user registered event
func (u *UniversalEventPublisher) PublishUserRegistered(ctx context.Context, user interface{}) error {
	// Handle different user types
	var data UserRegisteredData

	switch user := user.(type) {
	case map[string]interface{}:
		data = UserRegisteredData{
			UserID:    user["id"].(string),
			Email:     user["email"].(string),
			FullName:  user["full_name"].(string),
			Tier:      user["tier"].(string),
			CreatedAt: user["created_at"].(string),
		}
	default:
		// If it's a domain user, convert it
		// For now, just log and return
		log.Printf("Unsupported user type: %T", user)
		return nil
	}

	event, err := NewEvent(
		UserRegisteredEvent,
		"auth-service",
		"1.0",
		data,
	)
	if err != nil {
		return err
	}

	return u.eventBus.Publish(ctx, "user-events", event)
}

// PublishUserTierUpgraded publishes user tier upgraded event
func (u *UniversalEventPublisher) PublishUserTierUpgraded(ctx context.Context, userID string, oldTier, newTier interface{}) error {
	data := UserTierUpgradedData{
		UserID:     userID,
		OldTier:    convertToString(oldTier),
		NewTier:    convertToString(newTier),
		UpgradedAt: time.Now().UTC().Format(time.RFC3339),
	}

	event, err := NewEvent(
		UserTierUpgradedEvent,
		"auth-service",
		"1.0",
		data,
	)
	if err != nil {
		return err
	}

	return u.eventBus.Publish(ctx, "user-events", event)
}

// PublishUserQuotaUpdated publishes user quota updated event
func (u *UniversalEventPublisher) PublishUserQuotaUpdated(ctx context.Context, userID string, quotas interface{}) error {
	// For now, just log that we received the event
	log.Printf("User quota updated event for user: %s", userID)

	// TODO: Implement proper quota data extraction
	// This is a placeholder implementation
	data := UserQuotaUpdatedData{
		UserID: userID,
		Quotas: QuotaInfoData{
			AIDescription: QuotaData{Used: 0, Limit: 0},
			AIVideo:       QuotaData{Used: 0, Limit: 0},
			AutoPosting:   QuotaData{Used: 0, Limit: 0},
		},
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	event, err := NewEvent(
		UserQuotaUpdatedEvent,
		"user-service",
		"1.0",
		data,
	)
	if err != nil {
		return err
	}

	return u.eventBus.Publish(ctx, "user-events", event)
}

// Helper function to convert various types to string
func convertToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case domain.UserTier:
		return string(v)
	default:
		return "unknown"
	}
}
