package model

import (
	"github.com/google/uuid"
	"time"
)

type Request struct {
	Id        uuid.UUID     `json:"id"`
	UserId    uint64        `json:"user_id,string"`
	Type      RequestType   `json:"type"`
	CreatedAt time.Time     `json:"created_at"`
	GuildId   *uint64       `json:"guild_id"`
	Status    RequestStatus `json:"status"`
}

type RequestType string

const (
	RequestTypeGuildTranscripts RequestType = "guild_transcripts"
	RequestTypeGuildData        RequestType = "guild_data"
)

type RequestStatus string

const (
	RequestStatusQueued    RequestStatus = "queued"
	RequestStatusFailed    RequestStatus = "failed"
	RequestStatusCompleted RequestStatus = "completed"
)
