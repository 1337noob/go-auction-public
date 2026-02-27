package integration_event_handlers

import (
	"context"
	"errors"
	"main/notification/internal/application"
	"main/pkg/integration_event_bus"
)

type OnBidPlacedHandler struct {
	hub *application.WsHub
}

func NewOnBidPlacedHandler(hub *application.WsHub) *OnBidPlacedHandler {
	return &OnBidPlacedHandler{
		hub: hub,
	}
}

func (h *OnBidPlacedHandler) Handle(ctx context.Context, e integration_event_bus.Event) error {
	event, ok := e.(*integration_event_bus.BidPlaced)
	if !ok {
		return errors.New("invalid event")
	}

	msg := &application.Message{
		EventType: string(event.GetType()),
		EventData: event,
		Broadcast: true,
	}
	h.hub.SendMessage(event.AuctionID, msg)

	return nil
}
