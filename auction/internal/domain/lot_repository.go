package domain

import (
	"context"
)

type LotRepository interface {
	Save(ctx context.Context, lot *Lot) error
	FindByID(ctx context.Context, id string) (*Lot, error)
}
