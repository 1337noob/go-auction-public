package integration_event_handlers

import (
	"context"
	"main/metrics/internal/domain"
	"main/pkg/integration_event_bus"
	"time"
)

type OnAuctionCreatedHandler struct {
	repo domain.MetricsRepository
}

func NewOnAuctionCreatedHandler(repo domain.MetricsRepository) *OnAuctionCreatedHandler {
	return &OnAuctionCreatedHandler{
		repo: repo,
	}
}

func (h *OnAuctionCreatedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionCreated)
	if !ok {
		return nil
	}

	metrics := &domain.AuctionMetrics{
		AuctionID:     e.AuctionID,
		CreatedAt:     e.Timestamp,
		StartPrice:    e.StartPrice,
		MinBidStep:    e.MinBidStep,
		Status:        "created",
		BidAmounts:    make([]int, 0),
		BidTimestamps: make([]time.Time, 0),
	}

	return h.repo.SaveAuctionMetrics(metrics)
}
