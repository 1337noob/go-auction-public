package outbox

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

type OutboxWorker interface {
	StartOutboxWorker()
}

type SimpleOutboxWorker struct {
	repo      OutboxRepository
	publisher OutboxEventPublisher
	interval  time.Duration
	limit     int
	running   atomic.Bool
}

func NewSimpleOutboxWorker(repo OutboxRepository, publisher OutboxEventPublisher, interval time.Duration, limit int) *SimpleOutboxWorker {
	return &SimpleOutboxWorker{
		repo:      repo,
		publisher: publisher,
		interval:  interval,
		limit:     limit,
	}
}

func (w *SimpleOutboxWorker) StartOutboxWorker() {
	w.running.Store(true)
	go w.run()
}

func (w *SimpleOutboxWorker) Stop() {
	w.running.Store(false)
}

func (w *SimpleOutboxWorker) run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for w.running.Load() == true {
		select {
		case <-ticker.C:
			err := w.processPending()
			if err != nil {
				log.Printf("Error processing pending messages: %v", err)
			}
		}
	}
}

func (w *SimpleOutboxWorker) processPending() error {
	//log.Println("Processing pending messages")
	ctx := context.Background()
	messages, err := w.repo.GetPending(context.Background(), w.limit)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		err = w.repo.MarkAsProcessing(ctx, msg.ID)
		err = w.publish(ctx, msg)
		if err != nil {
			return err
		}
		err = w.repo.MarkAsCompleted(ctx, msg.ID)
	}

	//log.Printf("Processed %d messages", len(messages))
	return nil
}

func (w *SimpleOutboxWorker) publish(ctx context.Context, msg *OutboxMessage) error {
	err := w.publisher.Publish(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
