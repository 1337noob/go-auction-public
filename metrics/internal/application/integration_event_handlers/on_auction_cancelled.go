package integration_event_handlers

import (
	"context"
	"main/metrics/internal/domain"
	"main/pkg/integration_event_bus"
)

type OnAuctionCancelledHandler struct {
	repo domain.MetricsRepository
}

func NewOnAuctionCancelledHandler(repo domain.MetricsRepository) *OnAuctionCancelledHandler {
	return &OnAuctionCancelledHandler{
		repo: repo,
	}
}

func (h *OnAuctionCancelledHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionCancelled)
	if !ok {
		return nil
	}

	metrics, err := h.repo.GetAuctionMetrics(e.AuctionID)
	if err != nil || metrics == nil {
		return err
	}

	now := e.Timestamp
	metrics.CompletedAt = &now
	metrics.Status = "cancelled"

	if metrics.StartedAt != nil {
		duration := now.Sub(*metrics.StartedAt)
		metrics.Duration = &duration
	}

	return h.repo.SaveAuctionMetrics(metrics)
}
