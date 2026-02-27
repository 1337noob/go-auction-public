package integration_event_handlers

import (
	"context"
	"main/metrics/internal/domain"
	"main/pkg/integration_event_bus"
)

type OnAuctionCompletedHandler struct {
	repo domain.MetricsRepository
}

func NewOnAuctionCompletedHandler(repo domain.MetricsRepository) *OnAuctionCompletedHandler {
	return &OnAuctionCompletedHandler{
		repo: repo,
	}
}

func (h *OnAuctionCompletedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.AuctionCompleted)
	if !ok {
		return nil
	}

	metrics, err := h.repo.GetAuctionMetrics(e.AuctionID)
	if err != nil || metrics == nil {
		return err
	}

	metrics.CompletedAt = &e.CompletedAt
	metrics.FinalPrice = e.FinalPrice
	metrics.WinnerID = e.WinnerID
	metrics.Status = "completed"

	if metrics.StartedAt != nil {
		duration := e.CompletedAt.Sub(*metrics.StartedAt)
		metrics.Duration = &duration
	}

	err = h.repo.SaveAuctionMetrics(metrics)
	if err != nil {
		return err
	}

	if e.WinnerID != nil {
		err = h.repo.UpdateUserMetrics(*e.WinnerID, func(user *domain.UserMetrics) {
			user.WonAuctions++
		})
	}

	return err
}
