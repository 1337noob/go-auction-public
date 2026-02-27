package create_auction

import "time"

type CreateAuction struct {
	AggregateID string        `json:"aggregate_id"`
	LotID       string        `json:"lot_id"`
	StartPrice  int           `json:"start_price"`
	MinBidStep  int           `json:"min_bid_step"`
	SellerID    string        `json:"seller_id"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Timeout     time.Duration `json:"timeout"`
}
