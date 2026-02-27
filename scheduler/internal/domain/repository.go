package domain

import "context"

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetInit(ctx context.Context, limit int) ([]*Task, error)
	MarkAsCompleted(ctx context.Context, id string) error
}
