package read_model

import (
	"context"
)

type AuctionReadModelRepo interface {
	Save(ctx context.Context, auction *AuctionReadModel) error
	FindByID(ctx context.Context, id string) (*AuctionReadModel, error)
}
