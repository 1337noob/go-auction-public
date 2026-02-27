package integration_event_handlers

import (
	"context"
	"main/metrics/internal/domain"
	"main/pkg/integration_event_bus"
)

type OnBidPlacedHandler struct {
	repo domain.MetricsRepository
}

func NewOnBidPlacedHandler(repo domain.MetricsRepository) *OnBidPlacedHandler {
	return &OnBidPlacedHandler{
		repo: repo,
	}
}

func (h *OnBidPlacedHandler) Handle(ctx context.Context, event integration_event_bus.Event) error {
	e, ok := event.(*integration_event_bus.BidPlaced)
	if !ok {
		return nil
	}

	metrics, err := h.repo.GetAuctionMetrics(e.AuctionID)
	if err != nil {
		return err
	}
	if metrics == nil {
		return nil
	}

	metrics.BidCount++
	metrics.BidAmounts = append(metrics.BidAmounts, e.Amount)
	metrics.BidTimestamps = append(metrics.BidTimestamps, e.Timestamp)

	err = h.repo.SaveAuctionMetrics(metrics)
	if err != nil {
		return err
	}

	err = h.repo.UpdateUserMetrics(e.UserID, func(user *domain.UserMetrics) {
		user.TotalBids++
		user.TotalBidAmount += e.Amount
		user.LastBidAt = &e.Timestamp
	})

	return err
}
