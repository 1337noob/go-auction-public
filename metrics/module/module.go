package module

import (
	"main/metrics/internal/application/integration_event_handlers"
	"main/metrics/internal/application/queries"
	"main/metrics/internal/infrastructure/repositories"
	"main/pkg/integration_event_bus"
)

type MetricsModule struct {
	queryHandler *queries.GetMetricsHandler
}

func NewMetricsModule(integrationEventBus integration_event_bus.IntegrationEventBus) *MetricsModule {
	repo := repositories.NewInMemoryMetricsRepository()

	integrationEventBus.Subscribe(integration_event_bus.AuctionCreatedEventType, integration_event_handlers.NewOnAuctionCreatedHandler(repo))
	integrationEventBus.Subscribe(integration_event_bus.AuctionStartedEventType, integration_event_handlers.NewOnAuctionStartedHandler(repo))
	integrationEventBus.Subscribe(integration_event_bus.BidPlacedEventType, integration_event_handlers.NewOnBidPlacedHandler(repo))
	integrationEventBus.Subscribe(integration_event_bus.AuctionCompletedEventType, integration_event_handlers.NewOnAuctionCompletedHandler(repo))
	integrationEventBus.Subscribe(integration_event_bus.AuctionCancelledEventType, integration_event_handlers.NewOnAuctionCancelledHandler(repo))

	queryHandler := queries.NewGetMetricsHandler(repo)

	return &MetricsModule{
		queryHandler: queryHandler,
	}
}

func (m *MetricsModule) GetQueryHandler() *queries.GetMetricsHandler {
	return m.queryHandler
}
