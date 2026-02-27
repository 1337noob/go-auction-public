package event_store

import (
	"context"
	"encoding/json"
	"main/auction/internal/domain"
	"sync"
	"time"
)

type AuctionSnapshot struct {
	AggregateID string        `json:"aggregate_id"`
	Version     int           `json:"version"`
	LotID       string        `json:"lot_id"`
	StartPrice  int           `json:"start_price"`
	MinBidStep  int           `json:"min_bid_step"`
	SellerID    string        `json:"seller_id"`
	Bids        []*BidData    `json:"bids"`
	CurrentBid  *BidData      `json:"current_bid"`
	Status      string        `json:"status"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Timeout     time.Duration `json:"timeout"`
	CreatedAt   time.Time     `json:"created_at"`
}

type BidData struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type SnapshotRecord struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Data        []byte    `json:"data"`
	CreatedAt   time.Time `json:"created_at"`
}

type SnapshotStore interface {
	Save(ctx context.Context, aggregateID string, snapshot *AuctionSnapshot) error
	Load(ctx context.Context, aggregateID string) (*AuctionSnapshot, error)
}

type InMemorySnapshotStore struct {
	snapshots map[string]SnapshotRecord
	mu        sync.RWMutex
}

func NewInMemorySnapshotStore() *InMemorySnapshotStore {
	return &InMemorySnapshotStore{
		snapshots: make(map[string]SnapshotRecord),
	}
}

func (s *InMemorySnapshotStore) Save(ctx context.Context, aggregateID string, snapshot *AuctionSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	s.snapshots[aggregateID] = SnapshotRecord{
		AggregateID: aggregateID,
		Version:     snapshot.Version,
		Data:        data,
		CreatedAt:   time.Now(),
	}

	return nil
}

func (s *InMemorySnapshotStore) Load(ctx context.Context, aggregateID string) (*AuctionSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.snapshots[aggregateID]
	if !exists {
		return nil, nil
	}

	var snapshot AuctionSnapshot
	if err := json.Unmarshal(record.Data, &snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

func toSnapshot(auction *domain.Auction) *AuctionSnapshot {
	snapshot := &AuctionSnapshot{
		AggregateID: auction.GetID(),
		Version:     auction.GetVersion(),
		LotID:       auction.GetLotID(),
		StartPrice:  auction.GetStartPrice(),
		MinBidStep:  auction.GetMinBidStep(),
		SellerID:    auction.GetSellerID(),
		Status:      string(auction.GetStatus()),
		StartTime:   auction.GetStartTime(),
		EndTime:     auction.GetEndTime(),
		Timeout:     auction.GetTimeout(),
		CreatedAt:   auction.GetCreatedAt(),
	}

	for _, bid := range auction.GetBids() {
		snapshot.Bids = append(snapshot.Bids, &BidData{
			ID:        bid.ID,
			UserID:    bid.UserID,
			Amount:    bid.Amount,
			CreatedAt: bid.CreatedAt,
		})
	}

	if currentBid := auction.GetCurrentBid(); currentBid != nil {
		snapshot.CurrentBid = &BidData{
			ID:        currentBid.ID,
			UserID:    currentBid.UserID,
			Amount:    currentBid.Amount,
			CreatedAt: currentBid.CreatedAt,
		}
	}

	return snapshot
}

func fromSnapshot(snapshot *AuctionSnapshot) *domain.Auction {
	var bids []*domain.Bid
	for _, b := range snapshot.Bids {
		bids = append(bids, &domain.Bid{
			ID:        b.ID,
			UserID:    b.UserID,
			Amount:    b.Amount,
			CreatedAt: b.CreatedAt,
		})
	}

	var currentBid *domain.Bid
	if snapshot.CurrentBid != nil {
		currentBid = &domain.Bid{
			ID:        snapshot.CurrentBid.ID,
			UserID:    snapshot.CurrentBid.UserID,
			Amount:    snapshot.CurrentBid.Amount,
			CreatedAt: snapshot.CurrentBid.CreatedAt,
		}
	}

	return domain.ReconstructAuctionFromState(
		snapshot.AggregateID,
		snapshot.Version,
		snapshot.LotID,
		snapshot.StartPrice,
		snapshot.MinBidStep,
		snapshot.SellerID,
		bids,
		currentBid,
		domain.AuctionStatus(snapshot.Status),
		snapshot.StartTime,
		snapshot.EndTime,
		snapshot.Timeout,
		snapshot.CreatedAt,
	)
}
