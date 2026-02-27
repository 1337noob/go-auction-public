package integration_event_handlers

import (
	"context"
	"errors"
	"main/notification/internal/application"
	"main/pkg/integration_event_bus"
)

type OnAuctionCancelledHandler struct {
	hub *application.WsHub
}

func NewOnAuctionCancelledHandler(hub *application.WsHub) *OnAuctionCancelledHandler {
	return &OnAuctionCancelledHandler{
		hub: hub,
	}
}

func (h *OnAuctionCancelledHandler) Handle(ctx context.Context, e integration_event_bus.Event) error {
	event, ok := e.(*integration_event_bus.AuctionCancelled)
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
