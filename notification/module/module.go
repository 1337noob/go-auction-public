package module

import (
	"main/notification/internal/application"
	"main/notification/internal/application/integration_event_handlers"
	"main/pkg/integration_event_bus"

	"github.com/gorilla/websocket"
)

type NotificationModule struct {
	hub *application.WsHub
}

func NewNotificationModule(integrationEventBus integration_event_bus.IntegrationEventBus) *NotificationModule {
	hub := application.NewWsHub()
	integrationEventBus.Subscribe(integration_event_bus.AuctionStartedEventType, integration_event_handlers.NewOnAuctionStartedHandler(hub))
	integrationEventBus.Subscribe(integration_event_bus.BidPlacedEventType, integration_event_handlers.NewOnBidPlacedHandler(hub))
	integrationEventBus.Subscribe(integration_event_bus.BidRejectedEventType, integration_event_handlers.NewOnBidRejectedHandler(hub))
	integrationEventBus.Subscribe(integration_event_bus.AuctionCancelledEventType, integration_event_handlers.NewOnAuctionCancelledHandler(hub))
	integrationEventBus.Subscribe(integration_event_bus.AuctionCompletedEventType, integration_event_handlers.NewOnAuctionCompletedHandler(hub))

	return &NotificationModule{
		hub: hub,
	}
}

func (m *NotificationModule) AddConnection(auctionID string, userID string, conn *websocket.Conn) {
	m.hub.AddConnection(auctionID, userID, conn)
}

func (m *NotificationModule) RemoveConnection(auctionID string, userID string) {
	m.hub.RemoveConnection(auctionID, userID)
}

func (m *NotificationModule) SendMessage(auctionID string, message *application.Message) {
	m.hub.SendMessage(auctionID, message)
}
