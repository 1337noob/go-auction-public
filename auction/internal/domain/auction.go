package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuctionStatus string

const (
	AuctionStatusDraft     AuctionStatus = "draft"
	AuctionStatusCreated   AuctionStatus = "created"
	AuctionStatusStarted   AuctionStatus = "started"
	AuctionStatusCancelled AuctionStatus = "cancelled"
	AuctionStatusCompleted AuctionStatus = "completed"
)

type Auction struct {
	id               string
	version          int
	lotID            string
	startPrice       int
	minBidStep       int
	sellerID         string
	bids             []*Bid
	currentBid       *Bid
	status           AuctionStatus
	startTime        time.Time
	endTime          time.Time
	timeout          time.Duration
	createdAt        time.Time
	startedAt        *time.Time
	completedAt      *time.Time
	updatedAt        time.Time
	uncommitedEvents []Event
}

func NewAuction(
	aggregateID string,
	lotID string,
	startPrice int,
	minBidStep int,
	sellerID string,
	startTime time.Time,
	endTime time.Time,
	timeout time.Duration,
) (*Auction, error) {
	if startTime.Before(time.Now()) {
		return nil, ErrStartTimeInvalid
	}
	if endTime.Before(time.Now()) {
		return nil, ErrEndTimeInvalid
	}
	//if startTime.Equal(endTime) || startTime.After(endTime) {
	if !endTime.After(startTime) {
		return nil, ErrEndTimeInvalid
	}

	a := &Auction{
		version: 0,
	}

	event := &AuctionCreated{
		AggregateID: aggregateID,
		Version:     a.version + 1,
		Timestamp:   time.Now(),
		LotID:       lotID,
		StartPrice:  startPrice,
		MinBidStep:  minBidStep,
		SellerID:    sellerID,
		StartTime:   startTime,
		EndTime:     endTime,
		Timeout:     timeout,
	}

	a.record(event)

	return a, nil
}

func (a *Auction) Apply(event Event) {
	switch e := event.(type) {

	case *AuctionCreated:
		a.id = e.AggregateID
		a.lotID = e.LotID
		a.startPrice = e.StartPrice
		a.minBidStep = e.MinBidStep
		a.sellerID = e.SellerID
		a.startTime = e.StartTime
		a.endTime = e.EndTime
		a.timeout = e.Timeout
		a.status = AuctionStatusCreated
		a.createdAt = e.Timestamp
		a.updatedAt = e.Timestamp

	case *AuctionStarted:
		a.status = AuctionStatusStarted
		a.startedAt = &e.Timestamp
		a.updatedAt = e.Timestamp

	case *BidPlaced:
		newBid := &Bid{
			ID:        e.BidID,
			UserID:    e.UserID,
			Amount:    e.Amount,
			CreatedAt: e.Timestamp,
		}
		a.bids = append(a.bids, newBid)
		a.currentBid = newBid
		a.updatedAt = e.Timestamp

	case *AuctionCancelled:
		a.status = AuctionStatusCancelled
		a.updatedAt = e.Timestamp

	case *AuctionCompleted:
		a.status = AuctionStatusCompleted
		a.completedAt = &e.Timestamp
		a.updatedAt = e.Timestamp
	}

	a.version = event.GetVersion()
}

func (a *Auction) StartAuction() error {
	if a.status == AuctionStatusStarted {
		return ErrAuctionAlreadyStarted
	}
	if a.status == AuctionStatusCancelled {
		return ErrAuctionCancelled
	}
	if a.status == AuctionStatusCompleted {
		return ErrAuctionCompleted
	}
	if time.Since(a.startTime) > time.Second*10 {
		// TODO warning or cancel ???
	}
	if time.Now().After(a.endTime) {
		// TODO very bad
	}

	event := &AuctionStarted{
		AggregateID: a.id,
		Version:     a.version + 1,
		Timestamp:   time.Now(),
	}

	a.record(event)

	return nil
}

func (a *Auction) PlaceBid(userID string, amount int) error {
	if a.status == AuctionStatusCompleted {
		return ErrAuctionCompleted
	}
	if a.status == AuctionStatusCancelled {
		return ErrAuctionCancelled
	}
	if a.status != AuctionStatusStarted {
		return ErrAuctionNotStarted
	}

	minBidAmount := a.startPrice
	if a.currentBid != nil {
		minBidAmount = a.currentBid.Amount + a.minBidStep
	}

	if amount < minBidAmount {
		event := &BidRejected{
			AggregateID: a.id,
			Version:     a.version + 1,
			Timestamp:   time.Now(),
			UserID:      userID,
			Amount:      amount,
			Error:       ErrBidTooLow.Error(),
		}
		a.record(event)
		return nil
	}

	event := &BidPlaced{
		AggregateID: a.id,
		Version:     a.version + 1,
		Timestamp:   time.Now(),
		BidID:       uuid.NewString(),
		UserID:      userID,
		Amount:      amount,
	}

	a.record(event)

	return nil
}

func (a *Auction) CancelAuction(sellerID string, reason string) error {
	if a.status == AuctionStatusCancelled {
		return ErrAuctionCancelled
	}
	if a.status == AuctionStatusCompleted {
		return ErrAuctionCompleted
	}
	if a.sellerID != sellerID {
		return ErrSellerIDInvalid
	}

	event := &AuctionCancelled{
		AggregateID: a.id,
		Version:     a.version + 1,
		Timestamp:   time.Now(),
		Reason:      reason,
	}

	a.record(event)

	return nil
}

func (a *Auction) TimeoutAuction() error {
	if a.status == AuctionStatusCreated {
		return ErrAuctionNotStarted
	}
	if a.status == AuctionStatusCancelled {
		return ErrAuctionCancelled
	}
	if a.status == AuctionStatusCompleted {
		return ErrAuctionCompleted
	}

	now := time.Now()
	if a.currentBid != nil {
		timeoutAt := a.GetCurrentBid().CreatedAt.Add(a.GetTimeout())
		if now.Before(timeoutAt) {
			return ErrTimeoutTooEarly
		}
	}

	err := a.CompleteAuction()
	if err != nil {
		return err
	}

	return nil
}

func (a *Auction) CompleteAuction() error {
	if a.status == AuctionStatusCreated {
		return ErrAuctionNotStarted
	}
	if a.status == AuctionStatusCompleted {
		return ErrAuctionCompleted
	}
	if a.status == AuctionStatusCancelled {
		return ErrAuctionCancelled
	}

	event := &AuctionCompleted{
		AggregateID: a.id,
		Version:     a.version + 1,
		Timestamp:   time.Now(),
		CompletedAt: time.Now(),
	}

	if a.currentBid != nil {
		event.WinnerID = &a.currentBid.UserID
		event.FinalPrice = &a.currentBid.Amount
	}

	a.record(event)

	return nil
}

func (a *Auction) record(event Event) {
	a.uncommitedEvents = append(a.uncommitedEvents, event)
	a.Apply(event)
}

func (a *Auction) ClearUncommitedEvents() {
	a.uncommitedEvents = []Event{}
}

func ReconstructAuctionFromEvents(events []Event) *Auction {
	a := &Auction{}
	for _, event := range events {
		a.Apply(event)
	}

	return a
}

func (a *Auction) GetID() string                { return a.id }
func (a *Auction) GetVersion() int              { return a.version }
func (a *Auction) GetCurrentBid() *Bid          { return a.currentBid }
func (a *Auction) GetStatus() AuctionStatus     { return a.status }
func (a *Auction) GetBids() []*Bid              { return a.bids }
func (a *Auction) GetUncommitedEvents() []Event { return a.uncommitedEvents }
func (a *Auction) GetExpectedVersion() int      { return a.version - len(a.uncommitedEvents) }
func (a *Auction) GetLotID() string             { return a.lotID }
func (a *Auction) GetStartPrice() int           { return a.startPrice }
func (a *Auction) GetMinBidStep() int           { return a.minBidStep }
func (a *Auction) GetSellerID() string          { return a.sellerID }
func (a *Auction) GetStartTime() time.Time      { return a.startTime }
func (a *Auction) GetEndTime() time.Time        { return a.endTime }
func (a *Auction) GetTimeout() time.Duration    { return a.timeout }
func (a *Auction) GetCreatedAt() time.Time      { return a.createdAt }

func ReconstructAuctionFromState(
	id string,
	version int,
	lotID string,
	startPrice int,
	minBidStep int,
	sellerID string,
	bids []*Bid,
	currentBid *Bid,
	status AuctionStatus,
	startTime time.Time,
	endTime time.Time,
	timeout time.Duration,
	createdAt time.Time,
) *Auction {
	return &Auction{
		id:         id,
		version:    version,
		lotID:      lotID,
		startPrice: startPrice,
		minBidStep: minBidStep,
		sellerID:   sellerID,
		bids:       bids,
		currentBid: currentBid,
		status:     status,
		startTime:  startTime,
		endTime:    endTime,
		timeout:    timeout,
		createdAt:  createdAt,
	}
}
