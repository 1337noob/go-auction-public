package integration_event_handlers

import (
	"context"
	"errors"
	"main/auction/internal/application/commands/timeout_auction"
	"main/pkg/integration_event_bus"
)

type OnAuctionTimeoutReachedHandler struct {
	timeoutAuctionHandler *timeout_auction.TimeoutAuctionHandler
}

func NewOnAuctionTimeoutReachedHandler(timeoutAuctionHandler *timeout_auction.TimeoutAuctionHandler) *OnAuctionTimeoutReachedHandler {
	return &OnAuctionTimeoutReachedHandler{
		timeoutAuctionHandler: timeoutAuctionHandler,
	}
}

func (h *OnAuctionTimeoutReachedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionTimeoutReached)
	if !ok {
		return errors.New("invalid event type")
	}

	cmd := timeout_auction.TimeoutAuction{
		AggregateID: e.AuctionID,
	}

	err := h.timeoutAuctionHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
