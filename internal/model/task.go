package model

import "github.com/google/uuid"

type Task struct {
	Id        uuid.UUID `json:"id"`
	RequestId uuid.UUID `json:"request_id"`
}
