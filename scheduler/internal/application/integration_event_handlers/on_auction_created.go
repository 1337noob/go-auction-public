package integration_event_handlers

import (
	"context"
	"errors"
	"log"
	"main/pkg/integration_event_bus"
	"main/scheduler/internal/application"
	"main/scheduler/internal/domain"
	"time"

	"github.com/google/uuid"
)

type OnAuctionCreatedHandler struct {
	repo              domain.TaskRepository
	cachedAuctionRepo application.CachedAuctionRepository
}

func NewOnAuctionCreatedHandler(repo domain.TaskRepository, cachedAuctionRepo application.CachedAuctionRepository) *OnAuctionCreatedHandler {
	return &OnAuctionCreatedHandler{
		repo:              repo,
		cachedAuctionRepo: cachedAuctionRepo,
	}
}

func (h *OnAuctionCreatedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionCreated)
	if !ok {
		return errors.New("wrong event type")
	}

	err := h.cachedAuctionRepo.Create(ctx, &application.CachedAuction{
		ID:      e.AuctionID,
		Timeout: e.Timeout,
	})

	now := time.Now()

	task1 := &domain.Task{
		ID:          uuid.NewString(),
		AggregateID: e.AuctionID,
		Command:     domain.StartTimeReachedEventType,
		Status:      domain.TaskStatusInit,
		ExecuteTime: e.StartTime,
		CreatedAt:   now,
	}

	err = h.repo.Create(ctx, task1)
	if err != nil {
		return err
	}
	log.Println("task created: ", task1)

	task2 := &domain.Task{
		ID:          uuid.NewString(),
		AggregateID: e.AuctionID,
		Command:     domain.EndTimeReachedEventType,
		Status:      domain.TaskStatusInit,
		ExecuteTime: e.EndTime,
		CreatedAt:   now,
	}

	err = h.repo.Create(ctx, task2)
	if err != nil {
		return err
	}
	log.Println("task created: ", task2)

	task3 := &domain.Task{
		ID:          uuid.NewString(),
		AggregateID: e.AuctionID,
		Command:     domain.TimeoutReachedEventType,
		Status:      domain.TaskStatusInit,
		ExecuteTime: e.StartTime.Add(e.Timeout),
		CreatedAt:   now,
	}

	err = h.repo.Create(ctx, task3)
	if err != nil {
		return err
	}
	log.Println("task created: ", task3)

	return nil
}
