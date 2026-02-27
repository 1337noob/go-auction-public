package application

import "context"

type CachedAuctionRepository interface {
	Create(ctx context.Context, auction *CachedAuction) error
	FindByID(ctx context.Context, id string) (*CachedAuction, error)
}
