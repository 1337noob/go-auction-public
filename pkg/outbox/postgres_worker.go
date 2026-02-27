package outbox

import (
	"context"
	"log"
	"main/pkg"
	"sync/atomic"
	"time"
)

type PostgresOutboxWorker struct {
	txManager pkg.TransactionManager
	repo      OutboxRepository
	publisher OutboxEventPublisher
	interval  time.Duration
	limit     int
	running   atomic.Bool
}

func NewPostgresOutboxWorker(txManager pkg.TransactionManager, repo OutboxRepository, publisher OutboxEventPublisher, interval time.Duration, limit int) *PostgresOutboxWorker {
	return &PostgresOutboxWorker{
		txManager: txManager,
		repo:      repo,
		publisher: publisher,
		interval:  interval,
		limit:     limit,
	}
}

func (w *PostgresOutboxWorker) StartOutboxWorker() {
	w.running.Store(true)
	go w.run()
}

func (w *PostgresOutboxWorker) Stop() {
	w.running.Store(false)
}

func (w *PostgresOutboxWorker) run() {
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

func (w *PostgresOutboxWorker) processPending() error {
	ctx := context.Background()

	tx, err := w.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	baseCtx := tx.Context()
	ctx = pkg.WithTransaction(baseCtx, tx)

	messages, err := w.repo.GetPending(ctx, w.limit)
	if err != nil {
		tx.Rollback()
		return err
	}
	//log.Printf("Processing %d messages", len(messages))

	for _, msg := range messages {
		err = w.repo.MarkAsProcessing(ctx, msg.ID)
		err = w.publish(ctx, msg)
		if err != nil {
			log.Printf("Error marking message as processing: %v", err)
			tx.Rollback()
			return err
		}
		err = w.repo.MarkAsCompleted(ctx, msg.ID)
		if err != nil {
			log.Printf("Error marking message as completed: %v", err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	//log.Printf("Processed %d messages", len(messages))
	return nil
}

func (w *PostgresOutboxWorker) publish(ctx context.Context, msg *OutboxMessage) error {
	err := w.publisher.Publish(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
