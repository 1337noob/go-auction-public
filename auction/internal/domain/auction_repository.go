package domain

import (
	"context"
)

type AuctionRepository interface {
	Save(ctx context.Context, auction *Auction) error
	FindByID(ctx context.Context, id string) (*Auction, error)
}
