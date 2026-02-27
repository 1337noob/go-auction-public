package repositories

import (
	"main/metrics/internal/domain"
	"sync"
	"time"
)

type InMemoryMetricsRepository struct {
	auctions map[string]*domain.AuctionMetrics
	users    map[string]*domain.UserMetrics
	mu       sync.RWMutex
}

func NewInMemoryMetricsRepository() *InMemoryMetricsRepository {
	return &InMemoryMetricsRepository{
		auctions: make(map[string]*domain.AuctionMetrics),
		users:    make(map[string]*domain.UserMetrics),
	}
}

func (r *InMemoryMetricsRepository) SaveAuctionMetrics(metrics *domain.AuctionMetrics) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auctions[metrics.AuctionID] = metrics
	return nil
}

func (r *InMemoryMetricsRepository) GetAuctionMetrics(auctionID string) (*domain.AuctionMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics, exists := r.auctions[auctionID]
	if !exists {
		return nil, nil
	}
	return metrics, nil
}

func (r *InMemoryMetricsRepository) GetAllAuctionMetrics() ([]*domain.AuctionMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.AuctionMetrics, 0, len(r.auctions))
	for _, metrics := range r.auctions {
		result = append(result, metrics)
	}
	return result, nil
}

func (r *InMemoryMetricsRepository) UpdateUserMetrics(userID string, updateFn func(*domain.UserMetrics)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		user = &domain.UserMetrics{
			UserID: userID,
		}
		r.users[userID] = user
	}

	updateFn(user)

	if user.TotalBids > 0 {
		user.AverageBidAmount = float64(user.TotalBidAmount) / float64(user.TotalBids)
	}

	return nil
}

func (r *InMemoryMetricsRepository) GetUserMetrics(userID string) (*domain.UserMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics, exists := r.users[userID]
	if !exists {
		return nil, nil
	}
	return metrics, nil
}

func (r *InMemoryMetricsRepository) GetAllUserMetrics() ([]*domain.UserMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.UserMetrics, 0, len(r.users))
	for _, metrics := range r.users {
		result = append(result, metrics)
	}
	return result, nil
}

func (r *InMemoryMetricsRepository) GetGlobalMetrics() (*domain.GlobalMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allAuctions, _ := r.GetAllAuctionMetrics()

	global := &domain.GlobalMetrics{
		TotalAuctions: len(allAuctions),
	}

	var totalDuration time.Duration
	var totalFinalPrice int64
	var totalBids int
	var totalBidAmount int64
	var completedCount int
	var cancelledCount int

	for _, auction := range allAuctions {
		if auction.Duration != nil {
			totalDuration += *auction.Duration
		}
		if auction.FinalPrice != nil {
			totalFinalPrice += int64(*auction.FinalPrice)
		}
		totalBids += auction.BidCount
		for _, amount := range auction.BidAmounts {
			totalBidAmount += int64(amount)
		}

		switch auction.Status {
		case "completed":
			completedCount++
		case "cancelled":
			cancelledCount++
		}
	}

	global.CompletedAuctions = completedCount
	global.CancelledAuctions = cancelledCount
	global.TotalBids = totalBids

	if completedCount > 0 {
		global.AverageAuctionDuration = totalDuration / time.Duration(completedCount)
		global.AverageFinalPrice = float64(totalFinalPrice) / float64(completedCount)
	}

	if len(allAuctions) > 0 {
		global.AverageBidsPerAuction = float64(totalBids) / float64(len(allAuctions))
	}

	if totalBids > 0 {
		global.AverageBidAmount = float64(totalBidAmount) / float64(totalBids)
	}

	return global, nil
}
