package event_store

import (
	"context"
	"fmt"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"sync"
	"time"
)

type EventRecord struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	EventType   string    `json:"event_type"`
	EventData   []byte    `json:"event_data"`
	Timestamp   time.Time `json:"timestamp"`
}

type InMemoryEventStore struct {
	events     map[string][]EventRecord
	mu         sync.Mutex
	marshaller application.EventMarshaller
}

func NewInMemoryEventStore(marshaller application.EventMarshaller) *InMemoryEventStore {
	return &InMemoryEventStore{
		events:     make(map[string][]EventRecord),
		marshaller: marshaller,
	}
}

func (es *InMemoryEventStore) Save(ctx context.Context, aggregateID string, events []domain.Event, expectedVersion int) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	currentVersion := 0
	_, ok := es.events[aggregateID]
	if ok {
		for _, record := range es.events[aggregateID] {
			if record.Version > currentVersion {
				currentVersion = record.Version
			}
		}
	}

	if currentVersion != expectedVersion {
		return fmt.Errorf("expected version %d but got %d", expectedVersion, currentVersion)
	}

	for _, event := range events {
		eventData, err := es.marshaller.Marshal(event)
		if err != nil {
			return err
		}

		record := EventRecord{
			AggregateID: aggregateID,
			Version:     event.GetVersion(),
			EventType:   string(event.GetType()),
			EventData:   eventData,
			Timestamp:   event.GetTimestamp(),
		}

		es.events[aggregateID] = append(es.events[aggregateID], record)
	}

	return nil
}

func (es *InMemoryEventStore) Load(ctx context.Context, aggregateID string) ([]domain.Event, error) {
	return es.LoadFromVersion(ctx, aggregateID, 1)
}

func (es *InMemoryEventStore) LoadFromVersion(ctx context.Context, aggregateID string, fromVersion int) ([]domain.Event, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	records, exists := es.events[aggregateID]
	if !exists {
		return nil, fmt.Errorf("aggregate not found: %s", aggregateID)
	}

	var events []domain.Event
	for _, record := range records {
		if record.Version < fromVersion {
			continue
		}
		event, err := es.marshaller.Unmarshal(domain.EventType(record.EventType), record.EventData)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
