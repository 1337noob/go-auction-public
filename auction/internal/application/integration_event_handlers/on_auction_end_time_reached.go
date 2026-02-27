package integration_event_handlers

import (
	"context"
	"errors"
	"main/auction/internal/application/commands/complete_auction"
	"main/pkg/integration_event_bus"
)

type OnAuctionEndTimeReachedHandler struct {
	completeAuctionHandler *complete_auction.CompleteAuctionHandler
}

func NewOnAuctionEndTimeReachedHandler(completeAuctionHandler *complete_auction.CompleteAuctionHandler) *OnAuctionEndTimeReachedHandler {
	return &OnAuctionEndTimeReachedHandler{
		completeAuctionHandler: completeAuctionHandler,
	}
}

func (h *OnAuctionEndTimeReachedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionEndTimeReached)
	if !ok {
		return errors.New("invalid event type")
	}

	cmd := complete_auction.CompleteAuction{
		AggregateID: e.AuctionID,
	}

	err := h.completeAuctionHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
