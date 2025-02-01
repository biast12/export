package repository

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type ArtifactRepository struct {
	tx pgx.Tx
}

var (
	//go:embed sql/artifacts/create.sql
	queryArtifactsCreate string

	//go:embed sql/artifacts/get_global_size.sql
	queryArtifactsGetGlobalSize string
)

func NewArtifactRepository(tx pgx.Tx) *ArtifactRepository {
	return &ArtifactRepository{
		tx: tx,
	}
}

func (r *ArtifactRepository) Create(ctx context.Context, requestId uuid.UUID, key string, expiresAt time.Time, size int64) error {
	_, err := r.tx.Exec(ctx, queryArtifactsCreate, requestId, key, expiresAt, size)
	return err
}

func (r *ArtifactRepository) GetGlobalSize(ctx context.Context) (int64, error) {
	var size int64
	if err := r.tx.QueryRow(ctx, queryArtifactsGetGlobalSize).Scan(&size); err != nil {
		return 0, err
	}

	return size, nil
}
