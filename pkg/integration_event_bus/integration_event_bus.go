package integration_event_bus

import (
	"context"
	"errors"
	"log"
	"sync"
)

var (
	ErrNoHandlerFound = errors.New("no handler found")
)

type Handler interface {
	Handle(ctx context.Context, e Event) error
}

type IntegrationEventBus interface {
	Subscribe(eventType EventType, handler Handler)
	Publish(ctx context.Context, event Event) error
}

type InMemoryIntegrationEventBus struct {
	subscribers map[EventType][]Handler
	mu          sync.RWMutex
}

func NewInMemoryIntegrationEventBus() *InMemoryIntegrationEventBus {
	return &InMemoryIntegrationEventBus{
		subscribers: make(map[EventType][]Handler),
	}
}

func (b *InMemoryIntegrationEventBus) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *InMemoryIntegrationEventBus) Publish(ctx context.Context, event Event) error {
	log.Println("From IntegrationEventBus Publishing event: ", event.GetType())

	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, exists := b.subscribers[event.GetType()]
	if !exists {
		log.Println("No handlers found for event: ", event.GetType())
		return nil
	}

	for _, handler := range handlers {
		err := handler.Handle(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}
