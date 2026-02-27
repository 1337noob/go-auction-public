package projections

import (
	"context"
	"main/auction/internal/application/read_model"
	"sync"
)

type InMemoryAuctionReadModelRepo struct {
	auctions map[string]*read_model.AuctionReadModel
	mu       sync.Mutex
}

func NewInMemoryAuctionReadModelRepo() *InMemoryAuctionReadModelRepo {
	return &InMemoryAuctionReadModelRepo{
		auctions: make(map[string]*read_model.AuctionReadModel),
	}
}

func (r *InMemoryAuctionReadModelRepo) Save(ctx context.Context, auction *read_model.AuctionReadModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auctions[auction.ID] = auction

	return nil
}

func (r *InMemoryAuctionReadModelRepo) FindByID(ctx context.Context, id string) (*read_model.AuctionReadModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	auction, ok := r.auctions[id]
	if !ok {
		return nil, read_model.ErrAuctionReadModelNotFound
	}

	return auction, nil
}
