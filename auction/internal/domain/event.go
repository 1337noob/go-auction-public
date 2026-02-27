package domain

import "time"

type EventType string

type Event interface {
	GetAggregateID() string
	GetType() EventType
	GetVersion() int
	GetTimestamp() time.Time
}
