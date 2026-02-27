package outbox

import (
	"context"
	integrationeventbus "main/pkg/integration_event_bus"
)

type OutboxEventPublisher interface {
	Publish(ctx context.Context, msg *OutboxMessage) error
}

type SimpleOutboxEventPublisher struct {
	bus        integrationeventbus.IntegrationEventBus
	marshaller *integrationeventbus.EventMarshaller
}

func NewSimpleEventPublisher(bus integrationeventbus.IntegrationEventBus, marshaller *integrationeventbus.EventMarshaller) *SimpleOutboxEventPublisher {
	return &SimpleOutboxEventPublisher{
		bus:        bus,
		marshaller: marshaller,
	}
}

func (p *SimpleOutboxEventPublisher) Publish(ctx context.Context, msg *OutboxMessage) error {
	event, err := p.marshaller.UnmarshalEvent(msg.EventType, msg.EventData)
	if err != nil {
		return err
	}

	err = p.bus.Publish(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
