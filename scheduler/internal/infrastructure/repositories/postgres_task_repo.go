package repositories

import (
	"context"
	"errors"
	"main/pkg"
	"main/scheduler/internal/domain"
	"time"
)

type PostgresTaskRepo struct {
}

func NewPostgresTaskRepo() *PostgresTaskRepo {
	return &PostgresTaskRepo{}
}

func (r *PostgresTaskRepo) Create(ctx context.Context, task *domain.Task) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	query := "INSERT INTO scheduler.tasks (id, aggregate_id, command, status, execute_time, executed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := sqlTx.Tx().ExecContext(ctx, query, task.ID, task.AggregateID, task.Command, task.Status, task.ExecuteTime, task.ExecutedAt, task.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresTaskRepo) GetInit(ctx context.Context, limit int) ([]*domain.Task, error) {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return nil, errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return nil, errors.New("transaction is not a sql transaction")
	}

	query := `
SELECT id, aggregate_id, command, status, execute_time, executed_at, created_at
FROM scheduler.tasks WHERE status = $1 AND $2 >= execute_time
FOR UPDATE SKIP LOCKED
LIMIT $3
`

	rows, err := sqlTx.Tx().QueryContext(ctx, query, domain.TaskStatusInit, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.ID, &task.AggregateID, &task.Command, &task.Status, &task.ExecuteTime, &task.ExecutedAt, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *PostgresTaskRepo) MarkAsCompleted(ctx context.Context, id string) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	now := time.Now()
	query := "UPDATE scheduler.tasks SET status = $1, executed_at = $2 WHERE id = $3"
	_, err := sqlTx.Tx().ExecContext(ctx, query, domain.TaskStatusCompleted, now, id)
	if err != nil {
		return err
	}

	return nil
}
