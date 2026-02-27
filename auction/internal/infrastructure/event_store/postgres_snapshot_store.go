package event_store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"main/pkg"
)

type PostgresSnapshotStore struct {
}

func NewPostgresSnapshotStore() *PostgresSnapshotStore {
	return &PostgresSnapshotStore{}
}

func (s *PostgresSnapshotStore) Save(ctx context.Context, aggregateID string, snapshot *AuctionSnapshot) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	query := `
INSERT INTO auction.snapshot_store (aggregate_id, version, data, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (aggregate_id) DO UPDATE
SET version = excluded.version, data = excluded.data, created_at = excluded.created_at;
`
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	_, err = sqlTx.Tx().ExecContext(ctx, query, snapshot.AggregateID, snapshot.Version, data, snapshot.CreatedAt)
	if err != nil {
		return err
	}

	return nil

}

func (s *PostgresSnapshotStore) Load(ctx context.Context, aggregateID string) (*AuctionSnapshot, error) {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return nil, errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return nil, errors.New("transaction is not a sql transaction")
	}
	query := "SELECT data from auction.snapshot_store WHERE aggregate_id = $1 LIMIT 1"

	var data []byte
	err := sqlTx.Tx().QueryRowContext(ctx, query, aggregateID).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var snapshot *AuctionSnapshot
	err = json.Unmarshal(data, &snapshot)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
