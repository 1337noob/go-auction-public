package domain_event_handlers

import (
	"context"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
)

type OnAuctionCompletedHandler struct {
	readModelRepo read_model.AuctionReadModelRepo
}

func NewOnAuctionCompletedHandler(readModelRepo read_model.AuctionReadModelRepo) *OnAuctionCompletedHandler {
	return &OnAuctionCompletedHandler{
		readModelRepo: readModelRepo,
	}
}

func (h *OnAuctionCompletedHandler) Handle(ctx context.Context, event domain.Event) error {
	e, ok := event.(*domain.AuctionCompleted)
	if !ok {
		return ErrInvalidEvent
	}

	auction, err := h.readModelRepo.FindByID(ctx, e.AggregateID)
	if err != nil {
		return err
	}

	auction.WinnerID = e.WinnerID
	auction.FinalPrice = e.FinalPrice
	auction.Status = domain.AuctionStatusCompleted
	auction.Version = e.Version
	auction.CompletedAt = &e.Timestamp
	auction.UpdatedAt = e.Timestamp

	err = h.readModelRepo.Save(ctx, auction)
	if err != nil {
		return err
	}

	return nil
}
