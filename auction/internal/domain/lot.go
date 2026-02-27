package domain

import (
	"errors"
	"time"
)

const (
	LotStatusCreated   LotStatus = "created"
	LotStatusPublished LotStatus = "published"
)

var (
	ErrLotPublished = errors.New("lot published")
)

type LotStatus string

type Lot struct {
	id               string
	name             string
	description      string
	ownerID          string
	status           LotStatus
	version          int
	createdAt        time.Time
	uncommitedEvents []Event
}

func NewLot(id, name, description, ownerID string) *Lot {
	lot := &Lot{
		version: 0,
	}

	event := &LotCreated{
		AggregateID: id,
		Version:     lot.version + 1,
		Timestamp:   time.Now(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Status:      LotStatusCreated,
	}

	lot.record(event)

	return lot
}

func (l *Lot) Update(name, description string) error {
	if l.status == LotStatusPublished {
		return ErrLotPublished
	}

	event := &LotUpdated{
		AggregateID: l.id,
		Version:     l.version + 1,
		Timestamp:   time.Now(),
		Name:        name,
		Description: description,
	}

	l.record(event)

	return nil
}

func (l *Lot) Publish() error {
	if l.status == LotStatusPublished {
		return ErrLotPublished
	}

	event := &LotPublished{
		AggregateID: l.id,
		Version:     l.version + 1,
		Timestamp:   time.Now(),
		Status:      LotStatusPublished,
	}

	l.record(event)

	return nil
}

func (l *Lot) Apply(event Event) {
	switch e := event.(type) {
	case *LotCreated:
		l.id = e.AggregateID
		l.name = e.Name
		l.description = e.Description
		l.ownerID = e.OwnerID
		l.status = LotStatusCreated
		l.createdAt = e.Timestamp
	case *LotPublished:
		l.status = LotStatusPublished
	case *LotUpdated:
		l.name = e.Name
		l.description = e.Description
	}

	l.version = event.GetVersion()
}

func (l *Lot) record(event Event) {
	l.uncommitedEvents = append(l.uncommitedEvents, event)
	l.Apply(event)
}

func (l *Lot) ClearUncommitedEvents() {
	l.uncommitedEvents = []Event{}
}

func ReconstructLotFromEvents(events []Event) *Lot {
	l := &Lot{}
	for _, event := range events {
		l.Apply(event)
	}

	return l
}

func (l *Lot) GetID() string                { return l.id }
func (l *Lot) GetName() string              { return l.name }
func (l *Lot) GetDescription() string       { return l.description }
func (l *Lot) GetOwnerID() string           { return l.ownerID }
func (l *Lot) GetStatus() LotStatus         { return l.status }
func (l *Lot) GetVersion() int              { return l.version }
func (l *Lot) GetCreatedAt() time.Time      { return l.createdAt }
func (l *Lot) GetTimestamp() time.Time      { return l.createdAt }
func (l *Lot) GetUncommitedEvents() []Event { return l.uncommitedEvents }
func (l *Lot) GetExpectedVersion() int      { return l.version - len(l.uncommitedEvents) }
