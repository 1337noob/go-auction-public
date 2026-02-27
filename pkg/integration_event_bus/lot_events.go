package integration_event_bus

import "time"

type LotCreated struct {
	LotID       string    `json:"lot_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	Status      string    `json:"status"`
}

func (LotCreated) GetType() EventType { return LotCreatedEventType }

type LotPublished struct {
	LotID     string    `json:"lot_id"`
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

func (LotPublished) GetType() EventType { return LotPublishedEventType }

type LotUpdated struct {
	LotID       string    `json:"lot_id"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (LotUpdated) GetType() EventType { return LotUpdatedEventType }
