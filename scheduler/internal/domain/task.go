package domain

import "time"

type TaskStatus string

const (
	TaskStatusInit      TaskStatus = "init"
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string
	AggregateID string
	Command     EventType
	Status      TaskStatus
	ExecuteTime time.Time
	ExecutedAt  *time.Time
	CreatedAt   time.Time
}
