package application

import (
	"context"
	"main/scheduler/internal/domain"
)

type EventBus interface {
	Subscribe(eventType domain.EventType, handler EventHandler)
	Publish(ctx context.Context, event domain.Event) error
}
