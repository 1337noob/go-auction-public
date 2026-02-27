package projections

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"main/auction/internal/application/read_model"
	"main/pkg"
	"time"
)

type PostgresAuctionReadModelRepo struct {
	db *sql.DB
}

func NewPostgresAuctionReadModelRepo(db *sql.DB) *PostgresAuctionReadModelRepo {
	return &PostgresAuctionReadModelRepo{
		db: db,
	}
}

func (r *PostgresAuctionReadModelRepo) Save(ctx context.Context, auction *read_model.AuctionReadModel) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	var currentBid []byte
	var bids []byte
	var err error

	if auction.CurrentBid != nil {
		currentBid, err = json.Marshal(auction.CurrentBid)
		if err != nil {
			return err
		}
	}
	if auction.Bids != nil {
		bids, err = json.Marshal(auction.Bids)
		if err != nil {
			return err
		}
	}

	timeout := auction.Timeout.String()

	query := `
INSERT INTO auction.auction_read_model (
id, lot_id, lot_name, start_price, min_bid_step,
seller_id, current_bid, bids, winner_id, final_price,
status, start_time, end_time, timeout, created_at,
started_at, completed_at, updated_at, version
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) 
ON CONFLICT (id) DO UPDATE SET
lot_id=EXCLUDED.lot_id,
lot_name=EXCLUDED.lot_name,
start_price=EXCLUDED.start_price,
min_bid_step=EXCLUDED.min_bid_step,
seller_id=EXCLUDED.seller_id,
current_bid=EXCLUDED.current_bid,
bids=EXCLUDED.bids,
winner_id=EXCLUDED.winner_id,
final_price=EXCLUDED.final_price,
status=EXCLUDED.status,
start_time=EXCLUDED.start_time,
end_time=EXCLUDED.end_time,
timeout=EXCLUDED.timeout,
started_at=EXCLUDED.started_at,
completed_at=EXCLUDED.completed_at,
updated_at=EXCLUDED.updated_at,
version=EXCLUDED.version;
`
	_, err = sqlTx.Tx().ExecContext(
		ctx,
		query,
		auction.ID,
		auction.LotID,
		auction.LotName,
		auction.StartPrice,
		auction.MinBidStep,
		auction.SellerID,
		currentBid,
		bids,
		auction.WinnerID,
		auction.FinalPrice,
		auction.Status,
		auction.StartTime,
		auction.EndTime,
		timeout,
		auction.CreatedAt,
		auction.StartedAt,
		auction.CompletedAt,
		auction.UpdatedAt,
		auction.Version,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresAuctionReadModelRepo) FindByID(ctx context.Context, id string) (*read_model.AuctionReadModel, error) {
	var readModel read_model.AuctionReadModel
	var currentBidData []byte
	var bidsData []byte
	var timeoutString string

	query := `SELECT id, lot_id, lot_name, start_price, min_bid_step,
seller_id, current_bid, bids, winner_id, final_price,
status, start_time, end_time, timeout, created_at,
started_at, completed_at, updated_at, version
from auction.auction_read_model WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&readModel.ID,
		&readModel.LotID,
		&readModel.LotName,
		&readModel.StartPrice,
		&readModel.MinBidStep,
		&readModel.SellerID,
		&currentBidData,
		&bidsData,
		&readModel.WinnerID,
		&readModel.FinalPrice,
		&readModel.Status,
		&readModel.StartTime,
		&readModel.EndTime,
		&timeoutString,
		&readModel.CreatedAt,
		&readModel.StartedAt,
		&readModel.CompletedAt,
		&readModel.UpdatedAt,
		&readModel.Version,
	)
	if err != nil {
		return nil, err
	}
	var currentBid read_model.BidReadModel
	if currentBidData != nil {
		err = json.Unmarshal(currentBidData, &currentBid)
		if err != nil {
			return nil, err
		}
		readModel.CurrentBid = &currentBid
	}
	var bids []read_model.BidReadModel
	if bidsData != nil {
		err = json.Unmarshal(bidsData, &bids)
		if err != nil {
			return nil, err
		}
		readModel.Bids = bids
	}

	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		return nil, err
	}

	readModel.Timeout = timeout
	//readModel.StartTime = readModel.StartTime.UTC()
	//readModel.EndTime = readModel.EndTime.UTC()
	//readModel.CreatedAt = readModel.CreatedAt.UTC()
	//readModel.UpdatedAt = readModel.UpdatedAt.UTC()

	return &readModel, nil
}
