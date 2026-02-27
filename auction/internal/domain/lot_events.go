package domain

import "time"

const (
	LotCreatedEventType   EventType = "LotCreated"
	LotPublishedEventType EventType = "LotPublished"
	LotUpdatedEventType   EventType = "LotUpdated"
)

type LotCreated struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	Status      LotStatus `json:"status"`
}

func (e LotCreated) GetType() EventType      { return LotCreatedEventType }
func (e LotCreated) GetAggregateID() string  { return e.AggregateID }
func (e LotCreated) GetVersion() int         { return e.Version }
func (e LotCreated) GetTimestamp() time.Time { return e.Timestamp }

type LotPublished struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Status      LotStatus `json:"status"`
}

func (e LotPublished) GetType() EventType      { return LotPublishedEventType }
func (e LotPublished) GetAggregateID() string  { return e.AggregateID }
func (e LotPublished) GetVersion() int         { return e.Version }
func (e LotPublished) GetTimestamp() time.Time { return e.Timestamp }

type LotUpdated struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (e LotUpdated) GetType() EventType      { return LotUpdatedEventType }
func (e LotUpdated) GetAggregateID() string  { return e.AggregateID }
func (e LotUpdated) GetVersion() int         { return e.Version }
func (e LotUpdated) GetTimestamp() time.Time { return e.Timestamp }
