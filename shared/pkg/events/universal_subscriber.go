package events

import (
	"context"
)

// UniversalEventSubscriber works with any EventBus implementation
type UniversalEventSubscriber struct {
	eventBus EventBus
}

// NewUniversalEventSubscriber creates a new universal event subscriber
func NewUniversalEventSubscriber(eventBus EventBus) *UniversalEventSubscriber {
	return &UniversalEventSubscriber{
		eventBus: eventBus,
	}
}

// SubscribeToUserEvents subscribes to user-related events
func (u *UniversalEventSubscriber) SubscribeToUserEvents(ctx context.Context, handler EventHandler) error {
	return u.eventBus.Subscribe(ctx, "user-events", handler)
}

// SubscribeToProductEvents subscribes to product-related events
func (u *UniversalEventSubscriber) SubscribeToProductEvents(ctx context.Context, handler EventHandler) error {
	return u.eventBus.Subscribe(ctx, "product-events", handler)
}

// SubscribeToAIEvents subscribes to AI-related events
func (u *UniversalEventSubscriber) SubscribeToAIEvents(ctx context.Context, handler EventHandler) error {
	return u.eventBus.Subscribe(ctx, "ai-events", handler)
}
