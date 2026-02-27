package integration_event_bus

type EventType string

const (
	AuctionCreatedEventType   EventType = "AuctionCreated"
	AuctionStartedEventType   EventType = "AuctionStarted"
	AuctionCancelledEventType EventType = "AuctionCancelled"
	AuctionCompletedEventType EventType = "AuctionCompleted"
	BidPlacedEventType        EventType = "BidPlaced"
	BidRejectedEventType      EventType = "BidRejected"

	LotCreatedEventType   EventType = "LotCreated"
	LotPublishedEventType EventType = "LotPublished"
	LotUpdatedEventType   EventType = "LotUpdated"

	AuctionStartTimeReachedEventType EventType = "AuctionStartTimeReached"
	AuctionEndTimeReachedEventType   EventType = "AuctionEndTimeReached"
	AuctionTimeoutReachedEventType   EventType = "AuctionTimeoutReached"
)

type Event interface {
	GetType() EventType
}
