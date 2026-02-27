package event_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"main/auction/internal/application"
	"main/auction/internal/domain"
	"main/pkg"

	"github.com/google/uuid"
)

type PostgresEventStore struct {
	marshaller application.EventMarshaller
}

func NewPostgresEventStore(marshaller application.EventMarshaller) *PostgresEventStore {
	return &PostgresEventStore{
		marshaller: marshaller,
	}
}

func (es *PostgresEventStore) Save(ctx context.Context, aggregateID string, events []domain.Event, expectedVersion int) error {
	if len(events) == 0 {
		return nil
	}

	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return errors.New("transaction is not a sql transaction")
	}

	var currentVersion sql.NullInt64
	query := "SELECT MAX(version) FROM auction.event_store WHERE aggregate_id = $1"
	err := sqlTx.Tx().QueryRowContext(ctx, query, aggregateID).Scan(&currentVersion)
	if err != nil {
		return err
	}

	if currentVersion.Valid {
		current, err := currentVersion.Value()
		if err != nil {
			return err
		}
		if current != int64(expectedVersion) {
			return fmt.Errorf("expected version %d but got %d", expectedVersion, current)
		}
	} else {
		if expectedVersion != 0 {
			return errors.New("expected version is invalid")
		}
	}

	for _, event := range events {
		id := uuid.NewString()
		eventData, err := es.marshaller.Marshal(event)
		if err != nil {
			return err
		}
		query := "INSERT INTO auction.event_store (id, aggregate_id, version, event_type, event_data, timestamp) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err = sqlTx.Tx().ExecContext(ctx, query, id, event.GetAggregateID(), event.GetVersion(), event.GetType(), eventData, event.GetTimestamp())
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *PostgresEventStore) Load(ctx context.Context, aggregateID string) ([]domain.Event, error) {
	return es.LoadFromVersion(ctx, aggregateID, 1)
}

func (es *PostgresEventStore) LoadFromVersion(ctx context.Context, aggregateID string, fromVersion int) ([]domain.Event, error) {
	tx := pkg.TransactionFromContext(ctx)
	if tx == nil {
		return nil, errors.New("transaction not found in context")
	}
	sqlTx, ok := tx.(*pkg.SQLTransaction)
	if !ok {
		return nil, errors.New("transaction is not a sql transaction")
	}

	query := "SELECT event_type, event_data FROM auction.event_store WHERE aggregate_id = $1 AND version >= $2 ORDER BY version"

	rows, err := sqlTx.Tx().QueryContext(ctx, query, aggregateID, fromVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var eventType string
		var eventData []byte
		err = rows.Scan(&eventType, &eventData)
		if err != nil {
			return nil, err
		}

		event, err := es.marshaller.Unmarshal(domain.EventType(eventType), eventData)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return events, nil
}
