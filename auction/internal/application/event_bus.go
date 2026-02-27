package application

import (
	"context"
	"main/auction/internal/domain"
)

type EventBus interface {
	Subscribe(eventType domain.EventType, handler EventHandler)
	Publish(ctx context.Context, event domain.Event) error
}
