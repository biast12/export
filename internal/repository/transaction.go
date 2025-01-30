package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type TransactionContext interface {
	Rollback(ctx context.Context) error
	Requests() *RequestRepository
	Tasks() *TaskRepository
}

type PostgresTransactionContext struct {
	tx      pgx.Tx
	aborted bool
}

var _ TransactionContext = (*PostgresTransactionContext)(nil)

func (t *PostgresTransactionContext) Rollback(ctx context.Context) error {
	t.aborted = true
	return t.tx.Rollback(ctx)
}

func (t *PostgresTransactionContext) Requests() *RequestRepository {
	return NewRequestRepository(t.tx)
}

func (t *PostgresTransactionContext) Tasks() *TaskRepository {
	return NewTaskRepository(t.tx)
}
