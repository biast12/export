package transcriptstore

import "context"

type Client interface {
	GetTranscriptsForGuild(ctx context.Context, guildId uint64) (map[int][]byte, error)
}
