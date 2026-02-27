package repositories

import (
	"context"
	"errors"
	"main/auction/internal/application"
	"main/auction/internal/domain"
)

type LotRepository struct {
	eventStore application.EventStore
}

func NewLotRepository(eventStore application.EventStore) *LotRepository {
	return &LotRepository{
		eventStore: eventStore,
	}
}

func (r *LotRepository) Save(ctx context.Context, lot *domain.Lot) error {
	err := r.eventStore.Save(ctx, lot.GetID(), lot.GetUncommitedEvents(), lot.GetExpectedVersion())
	if err != nil {
		return err
	}

	return nil
}

func (r *LotRepository) FindByID(ctx context.Context, id string) (*domain.Lot, error) {
	events, err := r.eventStore.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("lot not found")
	}

	lot := domain.ReconstructLotFromEvents(events)

	return lot, nil
}
