package model

import (
	"github.com/google/uuid"
	"time"
)

type Artifact struct {
	Id        uuid.UUID `json:"id"`
	RequestId uuid.UUID `json:"request_id"`
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}
