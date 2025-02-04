package repository

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DownloadRepository struct {
	tx pgx.Tx
}

var (
	//go:embed sql/downloads/create.sql
	queryDownloadsCreate string

	//go:embed sql/downloads/get_user_daily_bytes.sql
	queryDownloadsGetUserDailyBytes string

	//go:embed sql/downloads/get_daily_bytes.sql
	queryDownloadsGetDailyBytes string
)

func NewDownloadRepository(tx pgx.Tx) *DownloadRepository {
	return &DownloadRepository{
		tx: tx,
	}
}

func (r *DownloadRepository) Create(ctx context.Context, userId uint64, artifactId uuid.UUID) error {
	_, err := r.tx.Exec(ctx, queryDownloadsCreate, userId, artifactId)
	return err
}

func (r *DownloadRepository) GetUserDailyBytes(ctx context.Context, userId uint64) (int64, error) {
	var size int64
	if err := r.tx.QueryRow(ctx, queryDownloadsGetUserDailyBytes, userId).Scan(&size); err != nil {
		return 0, err
	}

	return size, nil
}

func (r *DownloadRepository) GetDailyBytes(ctx context.Context) (int64, error) {
	var size int64
	if err := r.tx.QueryRow(ctx, queryDownloadsGetDailyBytes).Scan(&size); err != nil {
		return 0, err
	}

	return size, nil
}
