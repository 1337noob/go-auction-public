package event_forwarders

import (
	"context"
	"log"
	"main/pkg/integration_event_bus"
	"main/pkg/outbox"
	"main/scheduler/internal/domain"
	"time"

	"github.com/google/uuid"
)

type SchedulerEventForwarder struct {
	outbox     outbox.OutboxRepository
	marshaller *integration_event_bus.EventMarshaller
}

func NewScheduleEventForwarder(outboxRepo outbox.OutboxRepository, marshaller *integration_event_bus.EventMarshaller) *SchedulerEventForwarder {
	return &SchedulerEventForwarder{
		outbox:     outboxRepo,
		marshaller: marshaller,
	}
}

func (f *SchedulerEventForwarder) Handle(ctx context.Context, event domain.Event) error {
	var appEvent integration_event_bus.Event

	switch e := event.(type) {
	case *domain.TaskStartTimeReached:
		appEvent = &integration_event_bus.AuctionStartTimeReached{
			AuctionID: e.AggregateID,
			Timestamp: e.Timestamp,
		}
	case *domain.TaskTimeoutReached:
		appEvent = &integration_event_bus.AuctionTimeoutReached{
			AuctionID: e.AggregateID,
			Timestamp: e.Timestamp,
		}
	case *domain.TaskEndTimeReached:
		appEvent = &integration_event_bus.AuctionEndTimeReached{
			AuctionID: e.AggregateID,
			Timestamp: e.Timestamp,
		}
	default:
		log.Printf("unknown event type: %T", e)
	}

	eventData, err := f.marshaller.MarshalEvent(appEvent)
	if err != nil {
		return err
	}

	msg := &outbox.OutboxMessage{
		ID:            uuid.NewString(),
		AggregateID:   event.GetAggregateID(),
		AggregateType: "scheduler",
		EventType:     appEvent.GetType(),
		EventData:     eventData,
		Status:        outbox.OutboxStatusPending,
		CreatedAt:     time.Now(),
	}

	err = f.outbox.Save(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
