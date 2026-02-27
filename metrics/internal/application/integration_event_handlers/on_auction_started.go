package integration_event_handlers

import (
	"context"
	"main/metrics/internal/domain"
	"main/pkg/integration_event_bus"
)

type OnAuctionStartedHandler struct {
	repo domain.MetricsRepository
}

func NewOnAuctionStartedHandler(repo domain.MetricsRepository) *OnAuctionStartedHandler {
	return &OnAuctionStartedHandler{
		repo: repo,
	}
}

func (h *OnAuctionStartedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionStarted)
	if !ok {
		return nil
	}

	metrics, err := h.repo.GetAuctionMetrics(e.AuctionID)
	if err != nil || metrics == nil {
		return err
	}

	metrics.StartedAt = &e.Timestamp
	metrics.Status = "started"

	return h.repo.SaveAuctionMetrics(metrics)
}
