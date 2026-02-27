package repositories

import (
	"context"
	"errors"
	"log"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/auction/internal/infrastructure/event_store"
)

type AuctionRepository struct {
	eventStore    application.EventStore
	snapshotStore event_store.SnapshotStore
	policy        event_store.SnapshotPolicy
}

func NewAuctionRepository(eventStore application.EventStore, snapshotStore event_store.SnapshotStore, policy event_store.SnapshotPolicy) *AuctionRepository {
	return &AuctionRepository{
		eventStore:    eventStore,
		snapshotStore: snapshotStore,
		policy:        policy,
	}
}

func (r *AuctionRepository) Save(ctx context.Context, auction *domain.Auction) error {
	err := r.eventStore.Save(ctx, auction.GetID(), auction.GetUncommitedEvents(), auction.GetExpectedVersion())
	if err != nil {
		return err
	}

	if r.policy.ShouldTakeSnapshot(auction.GetVersion()) {
		log.Println("Saving snapshot for auction", auction.GetID())
		snapshot := r.toSnapshot(auction)
		err = r.snapshotStore.Save(ctx, auction.GetID(), snapshot)
		if err != nil {
			log.Printf("Error saving snapshot for auction %s: %v", auction.GetID(), err)
		}
	}

	return nil
}

func (r *AuctionRepository) FindByID(ctx context.Context, id string) (*domain.Auction, error) {
	snapshot, err := r.snapshotStore.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	if snapshot != nil {
		log.Println("Loading from auction snapshot ", id)
		auction := r.fromSnapshot(snapshot)
		events, err := r.eventStore.LoadFromVersion(ctx, id, snapshot.Version+1)
		if err != nil {
			return nil, err
		}
		for _, event := range events {
			auction.Apply(event)
		}
		return auction, nil
	}

	events, err := r.eventStore.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("auction not found")
	}

	auction := domain.ReconstructAuctionFromEvents(events)

	return auction, nil
}

func (r *AuctionRepository) fromSnapshot(snapshot *event_store.AuctionSnapshot) *domain.Auction {
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

func (r *AuctionRepository) toSnapshot(auction *domain.Auction) *event_store.AuctionSnapshot {
	snapshot := &event_store.AuctionSnapshot{
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
		snapshot.Bids = append(snapshot.Bids, &event_store.BidData{
			ID:        bid.ID,
			UserID:    bid.UserID,
			Amount:    bid.Amount,
			CreatedAt: bid.CreatedAt,
		})
	}

	if currentBid := auction.GetCurrentBid(); currentBid != nil {
		snapshot.CurrentBid = &event_store.BidData{
			ID:        currentBid.ID,
			UserID:    currentBid.UserID,
			Amount:    currentBid.Amount,
			CreatedAt: currentBid.CreatedAt,
		}
	}

	return snapshot
}
