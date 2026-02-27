package integration_event_bus

import (
	"encoding/json"
	"fmt"
)

type EventMarshaller struct {
}

func NewEventMarshaller() *EventMarshaller {
	return &EventMarshaller{}
}

func (m *EventMarshaller) MarshalEvent(event Event) ([]byte, error) {
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return eventData, nil
}

func (m *EventMarshaller) UnmarshalEvent(eventType EventType, data []byte) (Event, error) {
	switch eventType {
	case AuctionCreatedEventType:
		var auctionCreatedEvent AuctionCreated
		err := json.Unmarshal(data, &auctionCreatedEvent)
		if err != nil {
			return nil, err
		}
		return &auctionCreatedEvent, nil
	case AuctionStartedEventType:
		var auctionStartedEvent AuctionStarted
		err := json.Unmarshal(data, &auctionStartedEvent)
		if err != nil {
			return nil, err
		}
		return &auctionStartedEvent, nil
	case AuctionCancelledEventType:
		var auctionCancelledEvent AuctionCancelled
		err := json.Unmarshal(data, &auctionCancelledEvent)
		if err != nil {
			return nil, err
		}
		return &auctionCancelledEvent, nil
	case AuctionCompletedEventType:
		var auctionCompletedEvent AuctionCompleted
		err := json.Unmarshal(data, &auctionCompletedEvent)
		if err != nil {
			return nil, err
		}
		return &auctionCompletedEvent, nil
	case BidPlacedEventType:
		var bidPlacedEvent BidPlaced
		err := json.Unmarshal(data, &bidPlacedEvent)
		if err != nil {
			return nil, err
		}
		return &bidPlacedEvent, nil
	case BidRejectedEventType:
		var bidRejectedEvent BidRejected
		err := json.Unmarshal(data, &bidRejectedEvent)
		if err != nil {
			return nil, err
		}
		return &bidRejectedEvent, nil

	case LotCreatedEventType:
		var lotCreatedEvent LotCreated
		err := json.Unmarshal(data, &lotCreatedEvent)
		if err != nil {
			return nil, err
		}
		return &lotCreatedEvent, nil
	case LotPublishedEventType:
		var lotPublishedEvent LotPublished
		err := json.Unmarshal(data, &lotPublishedEvent)
		if err != nil {
			return nil, err
		}
		return &lotPublishedEvent, nil
	case LotUpdatedEventType:
		var lotUpdatedEvent LotUpdated
		err := json.Unmarshal(data, &lotUpdatedEvent)
		if err != nil {
			return nil, err
		}
		return &lotUpdatedEvent, nil
	case AuctionStartTimeReachedEventType:
		var taskStartTimeReachedEvent AuctionStartTimeReached
		err := json.Unmarshal(data, &taskStartTimeReachedEvent)
		if err != nil {
			return nil, err
		}
		return &taskStartTimeReachedEvent, nil
	case AuctionEndTimeReachedEventType:
		var taskEndTimeReachedEvent AuctionEndTimeReached
		err := json.Unmarshal(data, &taskEndTimeReachedEvent)
		if err != nil {
			return nil, err
		}
		return &taskEndTimeReachedEvent, nil
	case AuctionTimeoutReachedEventType:
		var taskTimeoutReachedEvent AuctionTimeoutReached
		err := json.Unmarshal(data, &taskTimeoutReachedEvent)
		if err != nil {
			return nil, err
		}
		return &taskTimeoutReachedEvent, nil
	default:
		return nil, fmt.Errorf("unknown event type: %v", eventType)
	}
}
