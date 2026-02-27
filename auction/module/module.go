package module

import (
	"context"
	"database/sql"
	"main/auction/internal/application"
	"main/auction/internal/application/commands/cancel_auction"
	"main/auction/internal/application/commands/complete_auction"
	"main/auction/internal/application/commands/create_auction"
	"main/auction/internal/application/commands/create_lot"
	"main/auction/internal/application/commands/place_bid"
	"main/auction/internal/application/commands/publish_lot"
	"main/auction/internal/application/commands/start_auction"
	"main/auction/internal/application/commands/timeout_auction"
	"main/auction/internal/application/domain_event_handlers"
	"main/auction/internal/application/integration_event_handlers"
	"main/auction/internal/application/queries/find_auction_by_id"
	"main/auction/internal/application/read_model"
	"main/auction/internal/domain"
	"main/auction/internal/infrastructure/event_bus"
	"main/auction/internal/infrastructure/event_forwarders"
	"main/auction/internal/infrastructure/event_store"
	"main/auction/internal/infrastructure/projections"
	"main/auction/internal/infrastructure/repositories"
	"main/pkg"
	integrationeventbus "main/pkg/integration_event_bus"
	"main/pkg/outbox"
	"time"

	"github.com/google/uuid"
)

type AuctionModule struct {
	createAuctionHandler   *create_auction.CreateAuctionHandler
	startAuctionHandler    *start_auction.StartAuctionHandler
	placeBidHandler        *place_bid.PlaceBidHandler
	cancelAuctionHandler   *cancel_auction.CancelAuctionHandler
	completeAuctionHandler *complete_auction.CompleteAuctionHandler
	findAuctionByIdHandler *find_auction_by_id.FindAuctionByIDHandler
	createLotHandler       *create_lot.CreateLotHandler
	publishLotHandler      *publish_lot.PublishLotHandler
}

func NewAuctionModule(integrationEventBus integrationeventbus.IntegrationEventBus, outboxRepo outbox.OutboxRepository, txManager pkg.TransactionManager, db *sql.DB) *AuctionModule {
	domainEventMarshaller := application.NewJsonEventMarshaller()
	postgresEventStore := event_store.NewPostgresEventStore(domainEventMarshaller)
	countBasedSnapshotPolicy := event_store.NewCountBasedSnapshotPolicy(10)
	postgresSnapshotStore := event_store.NewPostgresSnapshotStore()
	repo := repositories.NewAuctionRepository(postgresEventStore, postgresSnapshotStore, countBasedSnapshotPolicy)
	eventBus := event_bus.NewDomainEventBus()
	integrationEventMarshaller := integrationeventbus.NewEventMarshaller()
	postgresReadModelRepo := projections.NewPostgresAuctionReadModelRepo(db)
	lotRepo := repositories.NewLotRepository(postgresEventStore)
	createAuctionHandler := create_auction.NewCreateAuctionHandler(repo, lotRepo, eventBus, txManager)
	startAuctionHandler := start_auction.NewStartAuctionHandler(repo, eventBus)
	placeBidHandler := place_bid.NewPlaceBidHandler(repo, eventBus, txManager)
	cancelAuctionHandler := cancel_auction.NewCancelAuctionHandler(repo, eventBus, txManager)
	completeAuctionHandler := complete_auction.NewCompleteAuctionHandler(repo, eventBus)
	findAuctionByIdHandler := find_auction_by_id.NewFindAuctionByIDHandler(postgresReadModelRepo)
	timeOutAuctionHandler := timeout_auction.NewTimeoutAuctionHandler(repo, eventBus)

	createLotHandler := create_lot.NewCreateLotHandler(lotRepo, eventBus, txManager)
	publishLotHandler := publish_lot.NewPublishLotHandler(lotRepo, eventBus, txManager)

	eventBus.Subscribe(domain.AuctionCreatedEventType, domain_event_handlers.NewOnAuctionCreatedHandler(postgresReadModelRepo, lotRepo))
	eventBus.Subscribe(domain.AuctionStartedEventType, domain_event_handlers.NewOnAuctionStartedHandler(postgresReadModelRepo))
	eventBus.Subscribe(domain.BidPlacedEventType, domain_event_handlers.NewOnBidPlacedHandler(postgresReadModelRepo))
	eventBus.Subscribe(domain.AuctionCancelledEventType, domain_event_handlers.NewOnAuctionCancelledHandler(postgresReadModelRepo))
	eventBus.Subscribe(domain.AuctionCompletedEventType, domain_event_handlers.NewOnAuctionCompletedHandler(postgresReadModelRepo))

	auctionEventForwarder := event_forwarders.NewAuctionEventForwarder(outboxRepo, integrationEventMarshaller)
	eventBus.Subscribe(domain.AuctionCreatedEventType, auctionEventForwarder)
	eventBus.Subscribe(domain.AuctionStartedEventType, auctionEventForwarder)
	eventBus.Subscribe(domain.BidPlacedEventType, auctionEventForwarder)
	eventBus.Subscribe(domain.BidRejectedEventType, auctionEventForwarder)
	eventBus.Subscribe(domain.AuctionCancelledEventType, auctionEventForwarder)
	eventBus.Subscribe(domain.AuctionCompletedEventType, auctionEventForwarder)

	lotEventForwarder := event_forwarders.NewLotEventForwarder(outboxRepo, integrationEventMarshaller)
	eventBus.Subscribe(domain.LotCreatedEventType, lotEventForwarder)
	eventBus.Subscribe(domain.LotPublishedEventType, lotEventForwarder)
	eventBus.Subscribe(domain.LotUpdatedEventType, lotEventForwarder)

	integrationEventBus.Subscribe(integrationeventbus.AuctionStartTimeReachedEventType, integration_event_handlers.NewOnAuctionStartTimeReachedHandler(startAuctionHandler))
	integrationEventBus.Subscribe(integrationeventbus.AuctionEndTimeReachedEventType, integration_event_handlers.NewOnAuctionEndTimeReachedHandler(completeAuctionHandler))
	integrationEventBus.Subscribe(integrationeventbus.AuctionTimeoutReachedEventType, integration_event_handlers.NewOnAuctionTimeoutReachedHandler(timeOutAuctionHandler))

	return &AuctionModule{
		createAuctionHandler:   createAuctionHandler,
		startAuctionHandler:    startAuctionHandler,
		placeBidHandler:        placeBidHandler,
		cancelAuctionHandler:   cancelAuctionHandler,
		completeAuctionHandler: completeAuctionHandler,
		findAuctionByIdHandler: findAuctionByIdHandler,
		createLotHandler:       createLotHandler,
		publishLotHandler:      publishLotHandler,
	}
}

func (m *AuctionModule) CreateAuction(
	ctx context.Context,
	lotID string,
	startPrice int,
	minBidStep int,
	sellerID string,
	startTime time.Time,
	endTime time.Time,
	timeout time.Duration,
) (string, error) {
	id := uuid.NewString()

	createCmd := create_auction.CreateAuction{
		AggregateID: id,
		LotID:       lotID,
		StartPrice:  startPrice,
		MinBidStep:  minBidStep,
		SellerID:    sellerID,
		StartTime:   startTime,
		EndTime:     endTime,
		Timeout:     timeout,
	}

	err := m.createAuctionHandler.Handle(ctx, createCmd)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *AuctionModule) StartAuction(ctx context.Context, id string) error {
	cmd := start_auction.StartAuction{
		AggregateID: id,
	}

	err := m.startAuctionHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (m *AuctionModule) PlaceBid(ctx context.Context, auctionID string, userID string, amount int) error {
	cmd := place_bid.PlaceBid{
		AggregateID: auctionID,
		UserID:      userID,
		Amount:      amount,
	}

	err := m.placeBidHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (m *AuctionModule) CancelAuction(ctx context.Context, auctionID string, sellerID string, reason string) error {
	cmd := cancel_auction.CancelAuction{
		AggregateID: auctionID,
		SellerID:    sellerID,
		Reason:      reason,
	}

	err := m.cancelAuctionHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (m *AuctionModule) CompleteAuction(ctx context.Context, auctionID string) error {
	cmd := complete_auction.CompleteAuction{
		AggregateID: auctionID,
	}

	err := m.completeAuctionHandler.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (m *AuctionModule) FindAuctionByID(ctx context.Context, id string) (*read_model.AuctionReadModel, error) {
	a, err := m.findAuctionByIdHandler.Handle(
		ctx,
		find_auction_by_id.FindAuctionByID{AuctionID: id},
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *AuctionModule) GetAllAuctions(ctx context.Context) []*read_model.AuctionReadModel {
	// TODO not implemented
	return nil
}

func (m *AuctionModule) CreateLot(ctx context.Context, name, description, ownerID string) (string, error) {
	id := uuid.NewString()

	createCmd := &create_lot.CreateLot{
		AggregateID: id,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}

	err := m.createLotHandler.Handle(ctx, createCmd)
	if err != nil {
		return "", err
	}

	// publish for testing

	publishCmd := &publish_lot.PublishLot{
		AggregateID: id,
	}

	err = m.publishLotHandler.Handle(ctx, publishCmd)
	if err != nil {
		return "", err
	}

	return id, nil
}
