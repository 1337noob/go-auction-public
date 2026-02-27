package outbox

import (
	"main/pkg/integration_event_bus"
	"time"
)

type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusCompleted  OutboxStatus = "completed"
	OutboxStatusFailed     OutboxStatus = "failed"
)

type OutboxMessage struct {
	ID            string                          `json:"id"`
	AggregateID   string                          `json:"aggregate_id"`
	AggregateType string                          `json:"aggregate_type"`
	EventType     integration_event_bus.EventType `json:"event_type"`
	EventData     []byte                          `json:"event_data"`
	Status        OutboxStatus                    `json:"status"`
	CreatedAt     time.Time                       `json:"created_at"`
	//RetryCount    int                             `json:"retry_count"` // TODO
	//LastError     string                          `json:"last_error"`  // TODO
	//Metadata any `json:"metadata"`
}
