package application

import (
	"context"
	"log"
	"main/pkg"
	"main/scheduler/internal/domain"
	"time"
)

type Worker struct {
	repo      domain.TaskRepository
	bus       EventBus
	txManager pkg.TransactionManager
	limit     int
}

func NewWorker(repo domain.TaskRepository, bus EventBus, txManager pkg.TransactionManager, limit int) *Worker {
	return &Worker{
		repo:      repo,
		bus:       bus,
		txManager: txManager,
		limit:     limit,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			//log.Println("start loop")
			// tr start
			ctx := context.Background()
			tx, err := w.txManager.Begin(ctx)
			if err != nil {
				log.Printf("failed to begin transaction: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			baseCtx := tx.Context()
			ctx = pkg.WithTransaction(baseCtx, tx)

			tasks, err := w.repo.GetInit(ctx, w.limit)
			if err != nil {
				log.Printf("Error getting init tasks: %v", err)
				// tr rollback
				tx.Rollback()
				time.Sleep(1 * time.Second)
				continue
			}
			//log.Println("tasks found:", len(tasks))
			now := time.Now()
			hasError := false
			for _, task := range tasks {
				if hasError {
					break
				}

				switch task.Command {
				case domain.StartTimeReachedEventType:
					log.Println("start time reached")
					event := &domain.TaskStartTimeReached{
						AggregateID: task.AggregateID,
						Timestamp:   now,
					}
					err = w.bus.Publish(ctx, event)
					if err != nil {
						log.Printf("Error publishing task start time reached: %v", err)
						hasError = true
					}
				case domain.TimeoutReachedEventType:
					log.Println("command timeout reached")
					event := &domain.TaskTimeoutReached{
						AggregateID: task.AggregateID,
						Timestamp:   now,
					}
					err = w.bus.Publish(ctx, event)
					if err != nil {
						log.Printf("Error publishing task timeout reached: %v", err)
						hasError = true
					}
				case domain.EndTimeReachedEventType:
					log.Println("command end reached")
					event := &domain.TaskEndTimeReached{
						AggregateID: task.AggregateID,
						Timestamp:   now,
					}
					err = w.bus.Publish(ctx, event)
					if err != nil {
						log.Printf("Error publishing task end time reached: %v", err)
						hasError = true
					}
				default:
					log.Printf("Unrecognized task command: %v", task.Command)
				}

				if hasError {
					break
				}

				err = w.repo.MarkAsCompleted(ctx, task.ID)
				if err != nil {
					log.Printf("failed to mark task as completed: %v", err)
					hasError = true
					break
				}
			}

			if hasError {
				tx.Rollback()
			} else {
				// TODO tr commit
				err = tx.Commit()
				if err != nil {
					log.Printf("failed to commit transaction: %v", err)
				} else {
					//log.Println("commit from scheduler worker")
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *Worker) Stop() error {
	return nil
}
