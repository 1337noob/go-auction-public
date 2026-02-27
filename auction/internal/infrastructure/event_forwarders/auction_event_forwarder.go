package event_forwarders

import (
	"context"
	"fmt"
	"main/auction/internal/domain"
	"main/pkg/integration_event_bus"
	"main/pkg/outbox"
	"time"

	"github.com/google/uuid"
)

type AuctionEventForwarder struct {
	outbox     outbox.OutboxRepository
	marshaller *integration_event_bus.EventMarshaller
}

func NewAuctionEventForwarder(outbox outbox.OutboxRepository, marshaller *integration_event_bus.EventMarshaller) *AuctionEventForwarder {
	return &AuctionEventForwarder{
		outbox:     outbox,
		marshaller: marshaller,
	}
}

func (p *AuctionEventForwarder) Handle(ctx context.Context, event domain.Event) error {
	var appEvent integration_event_bus.Event
	switch e := event.(type) {
	case *domain.AuctionCreated:
		appEvent = &integration_event_bus.AuctionCreated{
			AuctionID:  e.GetAggregateID(),
			Version:    e.GetVersion(),
			Timestamp:  e.GetTimestamp(),
			LotID:      e.LotID,
			StartPrice: e.StartPrice,
			MinBidStep: e.MinBidStep,
			SellerID:   e.SellerID,
			StartTime:  e.StartTime,
			EndTime:    e.EndTime,
			Timeout:    e.Timeout,
		}
	case *domain.AuctionStarted:
		appEvent = &integration_event_bus.AuctionStarted{
			AuctionID: e.GetAggregateID(),
			Version:   e.GetVersion(),
			Timestamp: e.GetTimestamp(),
		}
	case *domain.BidPlaced:
		appEvent = &integration_event_bus.BidPlaced{
			AuctionID: e.GetAggregateID(),
			Version:   e.GetVersion(),
			Timestamp: e.GetTimestamp(),
			BidID:     e.BidID,
			UserID:    e.UserID,
			Amount:    e.Amount,
		}
	case *domain.BidRejected:
		appEvent = &integration_event_bus.BidRejected{
			AuctionID: e.GetAggregateID(),
			Version:   e.GetVersion(),
			Timestamp: e.GetTimestamp(),
			UserID:    e.UserID,
			Amount:    e.Amount,
			Error:     e.Error,
		}
	case *domain.AuctionCancelled:
		appEvent = &integration_event_bus.AuctionCancelled{
			AuctionID: e.GetAggregateID(),
			Version:   e.GetVersion(),
			Timestamp: e.GetTimestamp(),
			Reason:    e.Reason,
		}
	case *domain.AuctionCompleted:
		appEvent = &integration_event_bus.AuctionCompleted{
			AuctionID:   e.GetAggregateID(),
			Version:     e.GetVersion(),
			Timestamp:   e.GetTimestamp(),
			CompletedAt: e.CompletedAt,
			WinnerID:    e.WinnerID,
			FinalPrice:  e.FinalPrice,
		}
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}

	eventData, err := p.marshaller.MarshalEvent(appEvent)
	if err != nil {
		return err
	}

	msg := &outbox.OutboxMessage{
		ID:            uuid.NewString(),
		AggregateID:   event.GetAggregateID(),
		AggregateType: "auction",
		EventType:     appEvent.GetType(),
		EventData:     eventData,
		Status:        outbox.OutboxStatusPending,
		CreatedAt:     time.Now(),
	}

	err = p.outbox.Save(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
