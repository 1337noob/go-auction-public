package domain_event_handlers

import (
	"context"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
)

type OnAuctionCancelledHandler struct {
	readModelRepo read_model.AuctionReadModelRepo
}

func NewOnAuctionCancelledHandler(readModelRepo read_model.AuctionReadModelRepo) *OnAuctionCancelledHandler {
	return &OnAuctionCancelledHandler{
		readModelRepo: readModelRepo,
	}
}

func (h *OnAuctionCancelledHandler) Handle(ctx context.Context, event domain.Event) error {
	e, ok := event.(*domain.AuctionCancelled)
	if !ok {
		return ErrInvalidEvent
	}

	auction, err := h.readModelRepo.FindByID(ctx, e.AggregateID)
	if err != nil {
		return err
	}

	auction.Status = domain.AuctionStatusCancelled
	auction.Version = e.Version
	auction.UpdatedAt = e.Timestamp

	err = h.readModelRepo.Save(ctx, auction)
	if err != nil {
		return err
	}

	return nil
}
