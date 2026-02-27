package integration_event_handlers

import (
	"context"
	"errors"
	"log"
	integrationeventbus "main/pkg/integration_event_bus"
	"main/scheduler/internal/application"
	"main/scheduler/internal/domain"
	"time"

	"github.com/google/uuid"
)

type OnBidPlacedHandler struct {
	repo              domain.TaskRepository
	cachedAuctionRepo application.CachedAuctionRepository
}

func NewOnBidPlacedHandler(repo domain.TaskRepository, cachedAuctionRepo application.CachedAuctionRepository) *OnBidPlacedHandler {
	return &OnBidPlacedHandler{
		repo:              repo,
		cachedAuctionRepo: cachedAuctionRepo,
	}
}

func (h *OnBidPlacedHandler) Handle(ctx context.Context, event integrationeventbus.Event) error {
	e, ok := event.(*integrationeventbus.BidPlaced)
	if !ok {
		return errors.New("wrong event type")
	}

	auction, err := h.cachedAuctionRepo.FindByID(ctx, e.AuctionID)
	if err != nil {
		return err
	}

	now := time.Now()

	execTime := now.Add(auction.Timeout)

	task := &domain.Task{
		ID:          uuid.NewString(),
		AggregateID: e.AuctionID,
		Command:     domain.TimeoutReachedEventType,
		Status:      domain.TaskStatusInit,
		ExecuteTime: execTime,
		CreatedAt:   now,
	}

	err = h.repo.Create(ctx, task)
	if err != nil {
		return err
	}

	log.Println("task created: ", task)

	return nil
}
