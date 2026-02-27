package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/pkg"
)

type PostgresOutboxRepository struct {
	tableName string
}

func NewPostgresOutboxRepository(tableName string) *PostgresOutboxRepository {
	return &PostgresOutboxRepository{
		tableName: tableName,
	}
}

func (r *PostgresOutboxRepository) Save(ctx context.Context, message *OutboxMessage) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	jsonEventData, err := json.Marshal(message.EventData)
	if err != nil {
		return err
	}

	query := "INSERT INTO %s (id, aggregate_id, aggregate_type, event_type, event_data, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = sqlTx.Tx().ExecContext(ctx, r.table(query), message.ID, message.AggregateID, message.AggregateType, message.EventType, jsonEventData, message.Status, message.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresOutboxRepository) GetPending(ctx context.Context, limit int) ([]*OutboxMessage, error) {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return nil, errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return nil, errors.New("transaction is not a sql transaction")
	}

	query := `
SELECT id, aggregate_id, aggregate_type, event_type, event_data, status, created_at
FROM %s WHERE status = $1
ORDER BY created_at DESC 
FOR UPDATE SKIP LOCKED
LIMIT $2
`

	rows, err := sqlTx.Tx().QueryContext(ctx, r.table(query), OutboxStatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*OutboxMessage
	for rows.Next() {
		var m OutboxMessage
		var eventDataJson []byte
		err = rows.Scan(&m.ID, &m.AggregateID, &m.AggregateType, &m.EventType, &eventDataJson, &m.Status, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(eventDataJson, &m.EventData)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	return messages, nil
}

func (r *PostgresOutboxRepository) MarkAsProcessing(ctx context.Context, id string) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	query := "UPDATE %s SET status = $1 WHERE id = $2"
	_, err := sqlTx.Tx().ExecContext(ctx, r.table(query), OutboxStatusProcessing, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresOutboxRepository) MarkAsCompleted(ctx context.Context, id string) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	query := "UPDATE %s SET status = $1 WHERE id = $2"
	_, err := sqlTx.Tx().ExecContext(ctx, r.table(query), OutboxStatusCompleted, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresOutboxRepository) MarkAsFailed(ctx context.Context, id string) error {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	query := "UPDATE %s SET status = $1 WHERE id = $2"
	_, err := sqlTx.Tx().ExecContext(ctx, r.table(query), OutboxStatusFailed, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresOutboxRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
