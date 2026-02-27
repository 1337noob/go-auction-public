package queries

import (
	"main/metrics/internal/domain"
)

type GetMetricsHandler struct {
	repo domain.MetricsRepository
}

func NewGetMetricsHandler(repo domain.MetricsRepository) *GetMetricsHandler {
	return &GetMetricsHandler{
		repo: repo,
	}
}

func (h *GetMetricsHandler) GetGlobalMetrics() (*domain.GlobalMetrics, error) {
	return h.repo.GetGlobalMetrics()
}

func (h *GetMetricsHandler) GetAuctionMetrics(auctionID string) (*domain.AuctionMetrics, error) {
	return h.repo.GetAuctionMetrics(auctionID)
}

func (h *GetMetricsHandler) GetAllAuctionMetrics() ([]*domain.AuctionMetrics, error) {
	return h.repo.GetAllAuctionMetrics()
}

func (h *GetMetricsHandler) GetUserMetrics(userID string) (*domain.UserMetrics, error) {
	return h.repo.GetUserMetrics(userID)
}

func (h *GetMetricsHandler) GetAllUserMetrics() ([]*domain.UserMetrics, error) {
	return h.repo.GetAllUserMetrics()
}
