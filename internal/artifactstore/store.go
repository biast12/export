package artifactstore

import (
	"context"
	"github.com/google/uuid"
)

type ArtifactStore interface {
	Fetch(ctx context.Context, requestId uuid.UUID, key string) ([]byte, error)
	Store(ctx context.Context, requestId uuid.UUID, key string, data []byte) error
}
