package integration_event_bus

import "time"

type AuctionStartTimeReached struct {
	AuctionID string    `json:"auction_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (AuctionStartTimeReached) GetType() EventType { return AuctionStartTimeReachedEventType }

type AuctionEndTimeReached struct {
	AuctionID string    `json:"auction_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (AuctionEndTimeReached) GetType() EventType { return AuctionEndTimeReachedEventType }

type AuctionTimeoutReached struct {
	AuctionID string    `json:"auction_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (AuctionTimeoutReached) GetType() EventType { return AuctionTimeoutReachedEventType }
