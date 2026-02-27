package repositories

import (
	"context"
	"errors"
	"main/scheduler/internal/application"
	"sync"
)

type InMemoryCachedAuctionRepository struct {
	auctions map[string]*application.CachedAuction
	mu       sync.Mutex
}

func NewInMemoryCachedAuctionRepository() *InMemoryCachedAuctionRepository {
	return &InMemoryCachedAuctionRepository{
		auctions: make(map[string]*application.CachedAuction),
	}
}

func (r *InMemoryCachedAuctionRepository) Create(ctx context.Context, auction *application.CachedAuction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auctions[auction.ID] = auction

	return nil
}

func (r *InMemoryCachedAuctionRepository) FindByID(ctx context.Context, id string) (*application.CachedAuction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	auction, ok := r.auctions[id]
	if !ok {
		return nil, errors.New("cached auction not found")
	}

	return auction, nil
}
