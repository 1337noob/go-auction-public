package integration_event_handlers

import (
	"context"
	"errors"
	"main/auction/internal/application/commands/start_auction"
	"main/pkg/integration_event_bus"
)

type OnAuctionStartTimeReachedHandler struct {
	startAuctionHandler *start_auction.StartAuctionHandler
}

func NewOnAuctionStartTimeReachedHandler(startAuctionHandler *start_auction.StartAuctionHandler) *OnAuctionStartTimeReachedHandler {
	return &OnAuctionStartTimeReachedHandler{
		startAuctionHandler: startAuctionHandler,
	}
}

func (h *OnAuctionStartTimeReachedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionStartTimeReached)
	if !ok {
		return errors.New("invalid event type")
	}

	cmd := start_auction.StartAuction{
		AggregateID: e.AuctionID,
	}

	err := h.startAuctionHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
