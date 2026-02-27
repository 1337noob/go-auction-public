package application

import (
	"context"
	"main/auction/internal/domain"
)

type EventHandler interface {
	Handle(ctx context.Context, event domain.Event) error
}
