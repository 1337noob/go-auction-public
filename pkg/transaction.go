package pkg

import (
	"context"
	"database/sql"
)

type TransactionManager interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
	Context() context.Context
}

type transactionContextKey struct{}

func WithTransaction(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, transactionContextKey{}, tx)
}

func TransactionFromContext(ctx context.Context) Transaction {
	tx, ok := ctx.Value(transactionContextKey{}).(Transaction)
	if !ok {
		return nil
	}
	return tx
}

type SQLTransactionManager struct {
	db *sql.DB
}

func NewSQLTransactionManager(db *sql.DB) *SQLTransactionManager {
	return &SQLTransactionManager{db: db}
}

type SQLTransaction struct {
	tx  *sql.Tx
	ctx context.Context
}

func (m *SQLTransactionManager) Begin(ctx context.Context) (Transaction, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &SQLTransaction{
		tx:  tx,
		ctx: ctx,
	}, nil
}

func (t *SQLTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *SQLTransaction) Context() context.Context {
	return t.ctx
}

func (t *SQLTransaction) Tx() *sql.Tx {
	return t.tx
}
