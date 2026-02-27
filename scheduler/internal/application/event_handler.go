package application

import (
	"context"
	"main/scheduler/internal/domain"
)

type EventHandler interface {
	Handle(ctx context.Context, event domain.Event) error
}
