package integration_event_bus

import "time"

type AuctionCreated struct {
	AuctionID  string        `json:"auction_id"`
	Version    int           `json:"version"`
	Timestamp  time.Time     `json:"timestamp"`
	LotID      string        `json:"lot_id"`
	StartPrice int           `json:"start_price"`
	MinBidStep int           `json:"min_bid_step"`
	SellerID   string        `json:"seller_id"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Timeout    time.Duration `json:"timeout"`
}

func (AuctionCreated) GetType() EventType { return AuctionCreatedEventType }

type AuctionStarted struct {
	AuctionID string    `json:"auction_id"`
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

func (AuctionStarted) GetType() EventType { return AuctionStartedEventType }

type BidPlaced struct {
	AuctionID string    `json:"auction_id"`
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	BidID     string    `json:"bid_id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
}

func (BidPlaced) GetType() EventType { return BidPlacedEventType }

type BidRejected struct {
	AuctionID string    `json:"auction_id"`
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	Error     string    `json:"error"`
}

func (BidRejected) GetType() EventType { return BidRejectedEventType }

type AuctionCancelled struct {
	AuctionID string    `json:"auction_id"`
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason"`
}

func (AuctionCancelled) GetType() EventType { return AuctionCancelledEventType }

type AuctionCompleted struct {
	AuctionID   string    `json:"auction_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	CompletedAt time.Time `json:"completed_at"`
	WinnerID    *string   `json:"winner_id"`
	FinalPrice  *int      `json:"final_price"`
}

func (AuctionCompleted) GetType() EventType { return AuctionCompletedEventType }
