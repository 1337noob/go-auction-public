package repositories

import (
	"context"
	"errors"
	"log"
	"main/scheduler/internal/domain"
	"sync"
	"time"
)

type InMemoryTaskRepo struct {
	tasks map[string]*domain.Task
	mu    sync.Mutex
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *InMemoryTaskRepo) Create(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task

	return nil
}

func (r *InMemoryTaskRepo) GetInit(ctx context.Context, limit int) ([]*domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()

	var tasks []*domain.Task
	for _, task := range r.tasks {
		if task.Status == domain.TaskStatusInit {
			if now.After(task.ExecuteTime) {

				log.Println("task init with time")
				tasks = append(tasks, &domain.Task{
					ID:          task.ID,
					AggregateID: task.AggregateID,
					Command:     task.Command,
					Status:      domain.TaskStatusPending,
					ExecuteTime: task.ExecuteTime,
				})

				task.Status = domain.TaskStatusPending
				r.tasks[task.ID] = task
			}
		}

		if len(tasks) == limit {
			break
		}
	}

	return tasks, nil
}

func (r *InMemoryTaskRepo) MarkAsCompleted(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[id]
	if !ok {
		return errors.New("task not found")
	}

	now := time.Now()
	task.Status = domain.TaskStatusCompleted
	task.ExecutedAt = &now
	r.tasks[id] = task

	return nil
}
