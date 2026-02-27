package integration_event_handlers

import (
	"context"
	"errors"
	"main/notification/internal/application"
	"main/pkg/integration_event_bus"
)

type OnBidRejectedHandler struct {
	hub *application.WsHub
}

func NewOnBidRejectedHandler(hub *application.WsHub) *OnBidRejectedHandler {
	return &OnBidRejectedHandler{
		hub: hub,
	}
}

func (h *OnBidRejectedHandler) Handle(ctx context.Context, e integration_event_bus.Event) error {
	event, ok := e.(*integration_event_bus.BidRejected)
	if !ok {
		return errors.New("invalid event")
	}

	msg := &application.Message{
		EventType: string(event.GetType()),
		EventData: event,
		UserID:    event.UserID,
		Broadcast: false,
	}
	h.hub.SendMessage(event.AuctionID, msg)

	return nil
}
