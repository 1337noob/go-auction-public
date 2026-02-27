package application

import (
	"encoding/json"
	"errors"
	"main/auction/internal/domain"
)

type EventMarshaller interface {
	Marshal(event domain.Event) ([]byte, error)
	Unmarshal(eventType domain.EventType, data []byte) (domain.Event, error)
}

type JsonEventMarshaller struct{}

func NewJsonEventMarshaller() *JsonEventMarshaller {
	return &JsonEventMarshaller{}
}

func (m *JsonEventMarshaller) Marshal(event domain.Event) ([]byte, error) {
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return eventData, nil
}

func (m *JsonEventMarshaller) Unmarshal(eventType domain.EventType, data []byte) (domain.Event, error) {
	switch eventType {
	case domain.AuctionCreatedEventType:
		var event domain.AuctionCreated
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case domain.AuctionStartedEventType:
		var event domain.AuctionStarted
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case domain.AuctionCancelledEventType:
		var event domain.AuctionCancelled
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case domain.AuctionTimeoutEventType:
		var event domain.AuctionTimeout
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case domain.AuctionCompletedEventType:
		var event domain.AuctionCompleted
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case domain.BidPlacedEventType:
		var event domain.BidPlaced
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case domain.BidRejectedEventType:
		var event domain.BidRejected
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil

	case domain.LotCreatedEventType:
		var event domain.LotCreated
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, err
	case domain.LotPublishedEventType:
		var event domain.LotPublished
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, err
	case domain.LotUpdatedEventType:
		var event domain.LotUpdated
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		return &event, err

	default:
		return nil, errors.New("unknown event")
	}
}
