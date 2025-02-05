package worker

import (
	"context"
	"fmt"
	"github.com/TicketsBot/data-self-service/internal/model"
	"log/slog"
)

func (d *Daemon) handleGuildDataTask(ctx context.Context, task model.Task, request model.Request) error {
	if request.GuildId == nil || *request.GuildId == 0 {
		d.logger.Error("Guild ID is nil", slog.String("task_id", task.Id.String()))
		return fmt.Errorf("guild ID is nil")
	}

	guildId := *request.GuildId

	logger := d.logger.With(slog.Uint64("guild_id", guildId), "request_id", request.Id)

	guild
}
