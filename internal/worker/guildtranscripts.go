package worker

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/TicketsBot/data-self-service/internal/metrics"
	"github.com/TicketsBot/data-self-service/internal/model"
	"github.com/TicketsBot/data-self-service/internal/repository"
	"github.com/TicketsBot/data-self-service/internal/utils"
	"log/slog"
	"time"
)

const transcriptExpiry = time.Hour * 24 * 3
const maxActiveSize = 250 * 1024 * 1024 * 1024

func (d *Daemon) handleGuildTranscriptsTask(ctx context.Context, task model.Task, request model.Request) error {
	if request.GuildId == nil || *request.GuildId == 0 {
		d.logger.Error("Guild ID is nil", slog.String("task_id", task.Id.String()))
		return fmt.Errorf("guild ID is nil")
	}

	guildId := *request.GuildId

	transcripts, err := d.transcripts.GetTranscriptsForGuild(ctx, guildId)
	if err != nil {
		d.logger.Error("Failed to get transcripts for guild", err, slog.Uint64("guild_id", guildId))
		return err
	}

	files := make(map[string][]byte)

	guildIdStr := fmt.Sprintf("%d", guildId)
	files["guild_id.txt"] = []byte(guildIdStr)
	files["guild_id.txt.sig"] = []byte(utils.Base64Encode(ed25519.Sign(d.privateKey, []byte(guildIdStr))))

	for ticketId, transcript := range transcripts {
		sigData := make([]byte, 0, len(transcript)+len(guildIdStr)+2+6)
		sigData = append(sigData, []byte(guildIdStr)...)
		sigData = append(sigData, byte('|'))
		sigData = append(sigData, []byte(fmt.Sprintf("%d", ticketId))...)
		sigData = append(sigData, byte('|'))
		sigData = append(sigData, transcript...)

		files[fmt.Sprintf("transcripts/%d.json", ticketId)] = transcript
		files[fmt.Sprintf("transcripts/%d.json.sig", ticketId)] =
			[]byte(utils.Base64Encode(ed25519.Sign(d.privateKey, sigData)))
	}

	artifact, err := utils.BuildZip(files)
	if err != nil {
		d.logger.Error("Failed to build zip", err)
		return err
	}

	artifactSize := int64(len(artifact))

	var globalArtifactSize int64
	if err := d.repository.Tx(ctx, func(ctx context.Context, tx repository.TransactionContext) (err error) {
		globalArtifactSize, err = tx.Artifacts().GetGlobalSize(ctx)
		return err
	}); err != nil {
		d.logger.Error("Failed to get global artifact size", "error", err)
		return err
	}

	if globalArtifactSize+artifactSize > maxActiveSize {
		d.logger.Error("Artifact size exceeds maximum", slog.Int64("size", globalArtifactSize+artifactSize))
		return fmt.Errorf("artifact size exceeds maximum")
	}

	key := utils.RandomString(32)
	expiresAt := time.Now().Add(transcriptExpiry)
	if err := d.artifacts.Store(ctx, request.Id, key, expiresAt, artifact); err != nil {
		d.logger.Error("Failed to store artifact", "error", err)
		return err
	}

	metrics.ArtifactsUploaded.WithLabelValues(request.Type.String()).Inc()
	metrics.ArtifactsUploadedBytes.WithLabelValues(request.Type.String()).Add(float64(artifactSize))

	if err := d.repository.Tx(ctx, func(ctx context.Context, tx repository.TransactionContext) error {
		if err := tx.Requests().SetStatus(ctx, request.Id, model.RequestStatusCompleted); err != nil {
			return err
		}

		if err := tx.Artifacts().Create(ctx, request.Id, key, expiresAt, artifactSize); err != nil {
			return err
		}

		return nil
	}); err != nil {
		d.logger.Error("Failed to update request status", "error", err)
		return err
	}

	return nil
}
