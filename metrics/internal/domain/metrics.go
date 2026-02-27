package domain

import "time"

type AuctionMetrics struct {
	AuctionID     string
	CreatedAt     time.Time
	StartedAt     *time.Time
	CompletedAt   *time.Time
	Duration      *time.Duration
	FinalPrice    *int
	BidCount      int
	BidAmounts    []int
	BidTimestamps []time.Time
	WinnerID      *string
	Status        string
	StartPrice    int
	MinBidStep    int
}

type UserMetrics struct {
	UserID           string
	TotalBids        int
	WonAuctions      int
	TotalBidAmount   int
	AverageBidAmount float64
	LastBidAt        *time.Time
}

type GlobalMetrics struct {
	TotalAuctions          int
	CompletedAuctions      int
	CancelledAuctions      int
	TotalBids              int
	AverageAuctionDuration time.Duration
	AverageFinalPrice      float64
	AverageBidsPerAuction  float64
	AverageBidAmount       float64
}
