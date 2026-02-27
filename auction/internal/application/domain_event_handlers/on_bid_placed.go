package domain_event_handlers

import (
	"context"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
)

type OnBidPlacedHandler struct {
	readModelRepo read_model.AuctionReadModelRepo
}

func NewOnBidPlacedHandler(readModelRepo read_model.AuctionReadModelRepo) *OnBidPlacedHandler {
	return &OnBidPlacedHandler{
		readModelRepo: readModelRepo,
	}
}

func (h *OnBidPlacedHandler) Handle(ctx context.Context, event domain.Event) error {
	e, ok := event.(*domain.BidPlaced)
	if !ok {
		return ErrInvalidEvent
	}

	auction, err := h.readModelRepo.FindByID(ctx, e.AggregateID)
	if err != nil {
		return err
	}

	bid := read_model.BidReadModel{
		ID:        e.BidID,
		UserID:    e.UserID,
		Amount:    e.Amount,
		CreatedAt: e.Timestamp,
	}

	auction.Bids = append(auction.Bids, bid)
	auction.CurrentBid = &bid
	auction.Version = e.Version

	err = h.readModelRepo.Save(ctx, auction)
	if err != nil {
		return err
	}

	return nil
}
