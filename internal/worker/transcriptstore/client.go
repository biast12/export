package transcriptstore

import "context"

type Client interface {
	GetTranscriptsForGuild(ctx context.Context, guildId uint64) (*GetTranscriptsResponse, error)
}

type GetTranscriptsResponse struct {
	Transcripts map[int][]byte
	Failed      []int
}
