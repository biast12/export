package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type TransactionContext interface {
	Rollback(ctx context.Context) error
	Requests() *RequestRepository
	Tasks() *TaskRepository
	Artifacts() *ArtifactRepository
	Downloads() *DownloadRepository
	OAuth2() *OAuth2Repository
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

func (t *PostgresTransactionContext) Artifacts() *ArtifactRepository {
	return NewArtifactRepository(t.tx)
}

func (t *PostgresTransactionContext) Downloads() *DownloadRepository {
	return NewDownloadRepository(t.tx)
}

func (t *PostgresTransactionContext) OAuth2() *OAuth2Repository {
	return NewOAuth2Repository(t.tx)
}
