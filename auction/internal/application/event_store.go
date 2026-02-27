package application

import (
	"context"
	"main/auction/internal/domain"
	"time"
)

type Filters struct {
	FromDate time.Time
	ToDate   time.Time
	Version  int
}

type EventStore interface {
	Save(ctx context.Context, aggregateID string, events []domain.Event, expectedVersion int) error
	Load(ctx context.Context, aggregateID string) ([]domain.Event, error)
	LoadFromVersion(ctx context.Context, aggregateID string, fromVersion int) ([]domain.Event, error)
}
