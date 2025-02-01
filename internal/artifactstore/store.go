package artifactstore

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ArtifactStore interface {
	Fetch(ctx context.Context, requestId uuid.UUID, key string) ([]byte, error)
	Store(ctx context.Context, requestId uuid.UUID, key string, expiresAt time.Time, data []byte) error
}
