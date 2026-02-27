package domain_event_handlers

import (
	"context"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
)

type OnAuctionStartedHandler struct {
	readModelRepo read_model.AuctionReadModelRepo
}

func NewOnAuctionStartedHandler(readModelRepo read_model.AuctionReadModelRepo) *OnAuctionStartedHandler {
	return &OnAuctionStartedHandler{
		readModelRepo: readModelRepo,
	}
}

func (h *OnAuctionStartedHandler) Handle(ctx context.Context, event domain.Event) error {
	e, ok := event.(*domain.AuctionStarted)
	if !ok {
		return ErrInvalidEvent
	}

	auction, err := h.readModelRepo.FindByID(ctx, e.AggregateID)
	if err != nil {
		return err
	}

	auction.Status = domain.AuctionStatusStarted
	auction.Version = e.Version
	auction.StartedAt = &e.Timestamp
	auction.UpdatedAt = e.Timestamp

	err = h.readModelRepo.Save(ctx, auction)
	if err != nil {
		return err
	}

	return nil
}
