package module

import (
	"main/pkg"
	integrationeventbus "main/pkg/integration_event_bus"
	"main/pkg/outbox"
	"main/scheduler/internal/application"
	"main/scheduler/internal/application/integration_event_handlers"
	"main/scheduler/internal/domain"
	"main/scheduler/internal/infrastructure/event_bus"
	"main/scheduler/internal/infrastructure/event_forwarders"
	"main/scheduler/internal/infrastructure/repositories"
)

type SchedulerModule struct {
	worker *application.Worker
}

func NewSchedulerModule(integrationEventBus integrationeventbus.IntegrationEventBus, outboxRepo outbox.OutboxRepository, txManager pkg.TransactionManager) *SchedulerModule {
	postgresTaskRepo := repositories.NewPostgresTaskRepo()
	cachedAuctionRepo := repositories.NewPostgresCachedAuctionRepo()
	limit := 100

	eventMarshaller := integrationeventbus.NewEventMarshaller()
	integrationEventForwarder := event_forwarders.NewScheduleEventForwarder(outboxRepo, eventMarshaller)
	domainEventBus := event_bus.NewDomainEventBus()
	domainEventBus.Subscribe(domain.StartTimeReachedEventType, integrationEventForwarder)
	domainEventBus.Subscribe(domain.TimeoutReachedEventType, integrationEventForwarder)
	domainEventBus.Subscribe(domain.EndTimeReachedEventType, integrationEventForwarder)

	worker := application.NewWorker(postgresTaskRepo, domainEventBus, txManager, limit)
	worker.Start()

	integrationEventBus.Subscribe(integrationeventbus.AuctionCreatedEventType, integration_event_handlers.NewOnAuctionCreatedHandler(postgresTaskRepo, cachedAuctionRepo))
	integrationEventBus.Subscribe(integrationeventbus.BidPlacedEventType, integration_event_handlers.NewOnBidPlacedHandler(postgresTaskRepo, cachedAuctionRepo))

	return &SchedulerModule{
		worker: worker,
	}
}

func (s *SchedulerModule) Start() {
	s.worker.Start()
}

func (s *SchedulerModule) Stop() {
	s.worker.Stop()
}
