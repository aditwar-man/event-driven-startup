package events

import (
	"context"
	"log"
	"sync"
)

// MemoryEventBus implements EventBus interface using in-memory storage
type MemoryEventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// NewMemoryEventBus creates a new in-memory event bus
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Publish sends an event to all subscribers of the topic
func (m *MemoryEventBus) Publish(ctx context.Context, topic string, event *Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	handlers, exists := m.subscribers[topic]
	if !exists {
		return nil // No subscribers for this topic
	}

	// Execute handlers concurrently
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				log.Printf("Error handling event %s: %v", event.Type, err)
			}
		}(handler)
	}

	log.Printf("Event published: %s to topic: %s", event.Type, topic)
	return nil
}

// Subscribe registers a handler for a topic
func (m *MemoryEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.subscribers[topic] = append(m.subscribers[topic], handler)
	log.Printf("Handler subscribed to topic: %s", topic)
	return nil
}

// Close cleans up the event bus (no-op for memory bus)
func (m *MemoryEventBus) Close() error {
	log.Println("Memory event bus closed")
	return nil
}
