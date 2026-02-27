package repositories

import (
	"context"
	"errors"
	"main/pkg"
	"main/scheduler/internal/application"
	"time"
)

type PostgresCachedAuctionRepo struct{}

func NewPostgresCachedAuctionRepo() *PostgresCachedAuctionRepo {
	return &PostgresCachedAuctionRepo{}
}

func (r *PostgresCachedAuctionRepo) Create(ctx context.Context, auction *application.CachedAuction) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	timeout := auction.Timeout.String()
	query := "INSERT INTO scheduler.cached_auction (id, timeout) VALUES ($1, $2)"
	_, err := sqlTx.Tx().ExecContext(ctx, query, auction.ID, timeout)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresCachedAuctionRepo) FindByID(ctx context.Context, id string) (*application.CachedAuction, error) {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return nil, errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return nil, errors.New("transaction is not a sql transaction")
	}

	query := "SELECT id, timeout  from scheduler.cached_auction WHERE id = $1"

	var timeoutString string
	var auction application.CachedAuction
	err := sqlTx.Tx().QueryRowContext(ctx, query, id).Scan(&auction.ID, &timeoutString)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		return nil, err
	}
	auction.Timeout = timeout

	return &auction, nil
}
