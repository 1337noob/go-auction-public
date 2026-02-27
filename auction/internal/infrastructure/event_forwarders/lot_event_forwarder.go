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

type LotEventForwarder struct {
	outbox     outbox.OutboxRepository
	marshaller *integration_event_bus.EventMarshaller
}

func NewLotEventForwarder(outbox outbox.OutboxRepository, marshaller *integration_event_bus.EventMarshaller) *LotEventForwarder {
	return &LotEventForwarder{
		outbox:     outbox,
		marshaller: marshaller,
	}
}

func (p *LotEventForwarder) Handle(ctx context.Context, event domain.Event) error {
	var appEvent integration_event_bus.Event
	switch e := event.(type) {
	case *domain.LotCreated:
		appEvent = &integration_event_bus.LotCreated{
			LotID:       e.GetAggregateID(),
			Version:     e.GetVersion(),
			Timestamp:   e.GetTimestamp(),
			Name:        e.Name,
			Description: e.Description,
			OwnerID:     e.OwnerID,
			Status:      string(e.Status),
		}
	case *domain.LotPublished:
		appEvent = &integration_event_bus.LotPublished{
			LotID:     e.GetAggregateID(),
			Version:   e.GetVersion(),
			Timestamp: e.GetTimestamp(),
			Status:    string(e.Status),
		}
	case *domain.LotUpdated:
		appEvent = &integration_event_bus.LotUpdated{
			LotID:       e.GetAggregateID(),
			Version:     e.GetVersion(),
			Timestamp:   e.GetTimestamp(),
			Name:        e.Name,
			Description: e.Description,
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
		AggregateType: "lot",
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
