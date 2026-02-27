package domain

import "time"

const (
	AuctionCreatedEventType   EventType = "AuctionCreated"
	AuctionStartedEventType   EventType = "AuctionStarted"
	AuctionCancelledEventType EventType = "AuctionCancelled"
	AuctionTimeoutEventType   EventType = "AuctionTimeout"
	AuctionCompletedEventType EventType = "AuctionCompleted"
	BidPlacedEventType        EventType = "BidPlaced"
	BidRejectedEventType      EventType = "BidRejected"
)

type AuctionCreated struct {
	AggregateID string        `json:"aggregate_id"`
	Version     int           `json:"version"`
	Timestamp   time.Time     `json:"timestamp"`
	LotID       string        `json:"lot_id"`
	StartPrice  int           `json:"start_price"`
	MinBidStep  int           `json:"min_bid_step"`
	SellerID    string        `json:"seller_id"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Timeout     time.Duration `json:"timeout"`
}

func (e AuctionCreated) GetType() EventType      { return AuctionCreatedEventType }
func (e AuctionCreated) GetAggregateID() string  { return e.AggregateID }
func (e AuctionCreated) GetVersion() int         { return e.Version }
func (e AuctionCreated) GetTimestamp() time.Time { return e.Timestamp }

type AuctionStarted struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
}

func (e AuctionStarted) GetType() EventType      { return AuctionStartedEventType }
func (e AuctionStarted) GetAggregateID() string  { return e.AggregateID }
func (e AuctionStarted) GetVersion() int         { return e.Version }
func (e AuctionStarted) GetTimestamp() time.Time { return e.Timestamp }

type BidPlaced struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	BidID       string    `json:"bid_id"`
	UserID      string    `json:"user_id"`
	Amount      int       `json:"amount"`
}

func (e BidPlaced) GetType() EventType      { return BidPlacedEventType }
func (e BidPlaced) GetAggregateID() string  { return e.AggregateID }
func (e BidPlaced) GetVersion() int         { return e.Version }
func (e BidPlaced) GetTimestamp() time.Time { return e.Timestamp }

type AuctionCancelled struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Reason      string    `json:"reason"`
}

func (e AuctionCancelled) GetType() EventType      { return AuctionCancelledEventType }
func (e AuctionCancelled) GetAggregateID() string  { return e.AggregateID }
func (e AuctionCancelled) GetVersion() int         { return e.Version }
func (e AuctionCancelled) GetTimestamp() time.Time { return e.Timestamp }

type AuctionTimeout struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
}

func (e AuctionTimeout) GetType() EventType      { return AuctionTimeoutEventType }
func (e AuctionTimeout) GetAggregateID() string  { return e.AggregateID }
func (e AuctionTimeout) GetVersion() int         { return e.Version }
func (e AuctionTimeout) GetTimestamp() time.Time { return e.Timestamp }

type AuctionCompleted struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	CompletedAt time.Time `json:"completed_at"`
	WinnerID    *string   `json:"winner_id"`
	FinalPrice  *int      `json:"final_price"`
}

func (e AuctionCompleted) GetType() EventType      { return AuctionCompletedEventType }
func (e AuctionCompleted) GetAggregateID() string  { return e.AggregateID }
func (e AuctionCompleted) GetVersion() int         { return e.Version }
func (e AuctionCompleted) GetTimestamp() time.Time { return e.Timestamp }

type BidRejected struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	//BidID       string    `json:"bid_id"`
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
	Error  string `json:"error"`
}

func (e BidRejected) GetType() EventType      { return BidRejectedEventType }
func (e BidRejected) GetAggregateID() string  { return e.AggregateID }
func (e BidRejected) GetVersion() int         { return e.Version }
func (e BidRejected) GetTimestamp() time.Time { return e.Timestamp }
