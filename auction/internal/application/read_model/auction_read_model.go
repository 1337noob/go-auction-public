package read_model

import (
	"main/auction/internal/domain"
	"time"
)

type AuctionReadModel struct {
	ID          string               `json:"id"`
	LotID       string               `json:"lot_id"`
	LotName     string               `json:"lot_name"`
	StartPrice  int                  `json:"start_price"`
	MinBidStep  int                  `json:"min_bid_step"`
	SellerID    string               `json:"seller_id"`
	CurrentBid  *BidReadModel        `json:"current_bid"`
	Bids        []BidReadModel       `json:"bids"`
	WinnerID    *string              `json:"winner_id"`
	FinalPrice  *int                 `json:"final_price"`
	Status      domain.AuctionStatus `json:"status"`
	StartTime   time.Time            `json:"start_time"`
	EndTime     time.Time            `json:"end_time"`
	Timeout     time.Duration        `json:"timeout"`
	CreatedAt   time.Time            `json:"created_at"`
	StartedAt   *time.Time           `json:"started_at"`
	CompletedAt *time.Time           `json:"completed_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Version     int                  `json:"version"`
}

type BidReadModel struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
