package outbox

import (
	"context"
	"errors"
	"sync"
)

type OutboxRepository interface {
	Save(ctx context.Context, msg *OutboxMessage) error
	GetPending(ctx context.Context, limit int) ([]*OutboxMessage, error)
	MarkAsProcessing(ctx context.Context, id string) error
	MarkAsCompleted(ctx context.Context, id string) error
	MarkAsFailed(ctx context.Context, id string) error
}

type InMemoryOutboxRepository struct {
	messages map[string]*OutboxMessage
	mu       sync.Mutex
}

func NewInMemoryOutboxRepository() *InMemoryOutboxRepository {
	return &InMemoryOutboxRepository{
		messages: make(map[string]*OutboxMessage),
	}
}

func (r *InMemoryOutboxRepository) Save(ctx context.Context, msg *OutboxMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[msg.ID] = msg
	return nil
}

func (r *InMemoryOutboxRepository) GetPending(ctx context.Context, limit int) ([]*OutboxMessage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	messages := make([]*OutboxMessage, 0)
	for _, msg := range r.messages {
		if msg.Status == OutboxStatusPending {
			messages = append(messages, msg)
		}

		if len(messages) >= limit {
			break
		}
	}

	return messages, nil
}

func (r *InMemoryOutboxRepository) MarkAsProcessing(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg, ok := r.messages[id]
	if !ok {
		return errors.New("outbox message not found")
	}
	msg.Status = OutboxStatusProcessing
	r.messages[id] = msg
	return nil
}

func (r *InMemoryOutboxRepository) MarkAsCompleted(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg, ok := r.messages[id]
	if !ok {
		return errors.New("outbox message not found")
	}
	msg.Status = OutboxStatusCompleted
	r.messages[id] = msg
	return nil
}

func (r *InMemoryOutboxRepository) MarkAsFailed(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg, ok := r.messages[id]
	if !ok {
		return errors.New("outbox message not found")
	}
	msg.Status = OutboxStatusFailed
	r.messages[id] = msg
	return nil
}
