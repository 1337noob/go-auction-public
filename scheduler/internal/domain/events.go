package domain

import "time"

type EventType string

const (
	StartTimeReachedEventType EventType = "StartTimeReached"
	EndTimeReachedEventType   EventType = "EndTimeReached"
	TimeoutReachedEventType   EventType = "TimeoutReached"
)

type Event interface {
	GetAggregateID() string
	GetType() EventType
	GetTimestamp() time.Time
}

type TaskStartTimeReached struct {
	AggregateID string
	Timestamp   time.Time
}

func (e TaskStartTimeReached) GetAggregateID() string  { return e.AggregateID }
func (e TaskStartTimeReached) GetType() EventType      { return StartTimeReachedEventType }
func (e TaskStartTimeReached) GetTimestamp() time.Time { return e.Timestamp }

type TaskEndTimeReached struct {
	AggregateID string
	Timestamp   time.Time
}

func (e TaskEndTimeReached) GetAggregateID() string  { return e.AggregateID }
func (e TaskEndTimeReached) GetType() EventType      { return EndTimeReachedEventType }
func (e TaskEndTimeReached) GetTimestamp() time.Time { return e.Timestamp }

type TaskTimeoutReached struct {
	AggregateID string
	Timestamp   time.Time
}

func (e TaskTimeoutReached) GetAggregateID() string  { return e.AggregateID }
func (e TaskTimeoutReached) GetType() EventType      { return TimeoutReachedEventType }
func (e TaskTimeoutReached) GetTimestamp() time.Time { return e.Timestamp }
