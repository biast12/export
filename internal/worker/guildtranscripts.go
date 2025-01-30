package worker

import (
	"context"
	"fmt"
	"github.com/TicketsBot/data-self-service/internal/model"
	"log/slog"
)

func (d *Daemon) handleGuildTranscriptsTask(ctx context.Context, task model.Task, request model.Request) error {
	if request.GuildId == nil {
		d.logger.Error("Guild ID is nil", slog.String("task_id", task.Id.String()))
		return fmt.Errorf("guild ID is nil")
	}

	guildId := *request.GuildId

	transcripts, err := d.transcripts.GetTranscriptsForGuild(ctx, guildId)
	if err != nil {
		d.logger.Error("Failed to get transcripts for guild", err, slog.Uint64("guild_id", guildId))
		return err
	}

	for ticketId, transcript := range transcripts {
		fmt.Printf("Ticket %d: %s\n", ticketId, string(transcript))
	}

	return nil
}
