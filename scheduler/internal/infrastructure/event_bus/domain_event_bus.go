package event_bus

import (
	"context"
	"log"
	"main/scheduler/internal/application"
	"main/scheduler/internal/domain"
	"sync"
)

type DomainEventBus struct {
	subscribers map[domain.EventType][]application.EventHandler
	mu          sync.Mutex
}

func NewDomainEventBus() *DomainEventBus {
	return &DomainEventBus{
		subscribers: make(map[domain.EventType][]application.EventHandler),
	}
}

func (b *DomainEventBus) Subscribe(eventType domain.EventType, handler application.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *DomainEventBus) Publish(ctx context.Context, event domain.Event) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, exists := b.subscribers[event.GetType()]
	if !exists {
		log.Printf("no handlers for event type %s", event.GetType())
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
