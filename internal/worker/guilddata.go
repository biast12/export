package worker

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/database"
	"github.com/TicketsBot/export/internal/metrics"
	"github.com/TicketsBot/export/internal/model"
	"github.com/TicketsBot/export/internal/repository"
	"github.com/TicketsBot/export/internal/utils"
	"github.com/TicketsBot/export/pkg/dto"
	"github.com/jackc/pgx/v4"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

type guildDataTask struct {
	Name string
	F    func(ctx context.Context, guildId uint64, guildData *dto.GuildData) error
}

func newGuildDataTask(name string, f func(ctx context.Context, guildId uint64, guildData *dto.GuildData) error) guildDataTask {
	return guildDataTask{
		Name: name,
		F:    f,
	}
}

func (d *Daemon) handleGuildDataTask(ctx context.Context, task model.Task, request model.Request) error {
	if request.GuildId == nil || *request.GuildId == 0 {
		d.logger.Error("Guild ID is nil", slog.String("task_id", task.Id.String()))
		return fmt.Errorf("guild ID is nil")
	}

	guildId := *request.GuildId

	logger := d.logger.With(slog.Uint64("guild_id", guildId), "request_id", request.Id)

	data := dto.GuildData{
		GuildId: guildId,
	}

	tasks := []guildDataTask{
		newGuildDataTask("Fetch Active Language", d.fetchActiveLanguage),
		newGuildDataTask("Fetch Archive Channel", d.fetchArchiveChannel),
		newGuildDataTask("Fetch Archive Messages", d.fetchArchiveMessages),
		newGuildDataTask("Fetch Autoclose Settings", d.fetchAutocloseSettings),
		newGuildDataTask("Fetch Autoclose Excluded", d.fetchAutocloseExcluded),
		newGuildDataTask("Fetch Guild Blacklisted Users", d.fetchGuildBlacklistedUsers),
		newGuildDataTask("Fetch Channel Category", d.fetchChannelCategory),
		newGuildDataTask("Fetch Claim Settings", d.fetchClaimSettings),
		newGuildDataTask("Fetch Close Confirmation Enabled", d.fetchCloseConfirmationEnabled),
		newGuildDataTask("Fetch Close Reasons", d.fetchCloseReasons),
		newGuildDataTask("Fetch Custom Colours", d.fetchCustomColours),
		newGuildDataTask("Fetch Embed Fields", d.fetchEmbedFields),
		newGuildDataTask("Fetch Embeds", d.fetchEmbeds),
		newGuildDataTask("Fetch Exit Survey Responses", d.fetchExitSurveyResponses),
		newGuildDataTask("Fetch Feedback Enabled", d.fetchFeedbackEnabled),
		newGuildDataTask("Fetch First Response Times", d.fetchFirstResponseTimes),
		newGuildDataTask("Fetch Form Inputs", d.fetchFormInputs),
		newGuildDataTask("Fetch Forms", d.fetchForms),
		newGuildDataTask("Fetch Guild Is Globally Blacklisted", d.fetchGuildIsGloballyBlacklisted),
		newGuildDataTask("Fetch Guild Metadata", d.fetchGuildMetadata),
		newGuildDataTask("Fetch Multi Panels", d.fetchMultiPanels),
		newGuildDataTask("Fetch Multi Panel Targets", d.fetchMultiPanelTargets),
		newGuildDataTask("Fetch Naming Scheme", d.fetchNamingScheme),
		newGuildDataTask("Fetch On Call Users", d.fetchOnCallUsers),
		newGuildDataTask("Fetch Panel Access Control Rules", d.fetchPanelAccessControlRules),
		newGuildDataTask("Fetch Panel Mention User", d.fetchPanelMentionUser),
		newGuildDataTask("Fetch Panel Role Mentions", d.fetchPanelRoleMentions),
		newGuildDataTask("Fetch Panels", d.fetchPanels),
		newGuildDataTask("Fetch Panel Teams", d.fetchPanelTeams),
		newGuildDataTask("Fetch Participants", d.fetchParticipants),
		newGuildDataTask("Fetch User Permissions", d.fetchUserPermissions),
		newGuildDataTask("Fetch Guild Blacklisted Roles", d.fetchGuildBlacklistedRoles),
		newGuildDataTask("Fetch Role Permissions", d.fetchRolePermissions),
		newGuildDataTask("Fetch Service Ratings", d.fetchServiceRatings),
		newGuildDataTask("Fetch Settings", d.fetchSettings),
		newGuildDataTask("Fetch Support Team Users", d.fetchSupportTeamUsers),
		newGuildDataTask("Fetch Support Team Roles", d.fetchSupportTeamRoles),
		newGuildDataTask("Fetch Support Teams", d.fetchSupportTeams),
		newGuildDataTask("Fetch Tags", d.fetchTags),
		newGuildDataTask("Fetch Ticket Claims", d.fetchTicketClaims),
		newGuildDataTask("Fetch Ticket Last Messages", d.fetchTicketLastMessages),
		newGuildDataTask("Fetch Ticket Limit", d.fetchTicketLimit),
		newGuildDataTask("Fetch Ticket Additional Members", d.fetchTicketAdditionalMembers),
		newGuildDataTask("Fetch Ticket Permissions", d.fetchTicketPermissions),
		newGuildDataTask("Fetch Tickets", d.fetchTickets),
		newGuildDataTask("Fetch Users Can Close", d.fetchUsersCanClose),
		newGuildDataTask("Fetch Welcome Message", d.fetchWelcomeMessage),
	}

	group, groupCtx := errgroup.WithContext(ctx)
	for _, task := range tasks {
		task := task

		group.Go(func() error {
			now := time.Now()

			logger.DebugContext(groupCtx, "Running task", "name", task.Name)
			if err := task.F(groupCtx, guildId, &data); err != nil {
				logger.ErrorContext(ctx, "Failed to run task", "error", err, "elapsed", time.Since(now))
				return err
			}

			logger.DebugContext(groupCtx, "Task completed", "name", task.Name, "elapsed", time.Since(now))
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		d.logger.Error("Failed to run tasks", "error", err)
		return err
	}

	logger.InfoContext(ctx, "All tasks completed")

	files := make(map[string][]byte)

	marshalled, err := json.Marshal(data)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to marshal data", "error", err)
		return err
	}

	files["data.json"] = marshalled
	files["data.json.sig"] = []byte(utils.Base64Encode(ed25519.Sign(d.privateKey, marshalled)))

	artifact, err := utils.BuildZip(files)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to build zip", "error", err)
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

	logger.InfoContext(ctx, "Uploading artifact", slog.Int64("size", artifactSize))

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

func (d *Daemon) fetchActiveLanguage(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.ActiveLanguage, d.database.ActiveLanguage.Get)
}

func (d *Daemon) fetchArchiveChannel(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchPtr(ctx, guildId, &guildData.ArchiveChannel, d.database.ArchiveChannel.Get)
}

func (d *Daemon) fetchArchiveMessages(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, channel_id, message_id
FROM archive_messages
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.ArchiveMessages, query,
		func(rows pgx.Rows) (dto.TicketUnion[database.ArchiveMessage], error) {
			var res dto.TicketUnion[database.ArchiveMessage]
			if err := rows.Scan(&res.TicketId, &res.Data.ChannelId, &res.Data.MessageId); err != nil {
				return dto.TicketUnion[database.ArchiveMessage]{}, err
			}

			return res, nil
		})
}

func (d *Daemon) fetchAutocloseSettings(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.AutocloseSettings, d.database.AutoClose.Get)
}

func (d *Daemon) fetchAutocloseExcluded(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id
FROM auto_close_exclude
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.AutocloseExcluded, query,
		func(rows pgx.Rows) (int, error) {
			var ticketId int
			if err := rows.Scan(&ticketId); err != nil {
				return 0, err
			}

			return ticketId, nil
		})
}

func (d *Daemon) fetchGuildBlacklistedUsers(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT user_id
FROM blacklist
WHERE guild_id = $1
ORDER BY user_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.GuildBlacklistedUsers, query, func(rows pgx.Rows) (uint64, error) {
		var userId uint64
		if err := rows.Scan(&userId); err != nil {
			return 0, err
		}

		return userId, nil
	})
}

func (d *Daemon) fetchChannelCategory(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.ChannelCategory, d.database.ChannelCategory.Get)
}

func (d *Daemon) fetchClaimSettings(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.ClaimSettings, d.database.ClaimSettings.Get)
}

func (d *Daemon) fetchCloseConfirmationEnabled(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.CloseConfirmationEnabled, d.database.CloseConfirmation.Get)
}

func (d *Daemon) fetchCloseReasons(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, close_reason, closed_by
FROM close_reason
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.CloseReasons, query,
		func(rows pgx.Rows) (dto.TicketUnion[database.CloseMetadata], error) {
			var data dto.TicketUnion[database.CloseMetadata]
			if err := rows.Scan(&data.TicketId, &data.Data.Reason, &data.Data.ClosedBy); err != nil {
				return dto.TicketUnion[database.CloseMetadata]{}, err
			}

			return data, nil
		})
}

func (d *Daemon) fetchCustomColours(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT colour_id, colour_code
FROM custom_colours
WHERE guild_id = $1
ORDER BY colour_id ASC LIMIT $2 OFFSET $3;`

	return fetchMap(ctx, d.database, guildId, &guildData.CustomColors, query)
}

func (d *Daemon) fetchEmbedFields(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT f.id, f.embed_id, f.name, f.value, f.inline
FROM embed_fields f
INNER JOIN embeds e
ON e.id = f.embed_id
WHERE e.guild_id = $1
ORDER BY f.embed_id, f.id ASC
LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.EmbedFields, query, func(rows pgx.Rows) (database.EmbedField, error) {
		var field database.EmbedField
		if err := rows.Scan(&field.FieldId, &field.EmbedId, &field.Name, &field.Value, &field.Inline); err != nil {
			return database.EmbedField{}, err
		}

		return field, nil
	})
}

func (d *Daemon) fetchEmbeds(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT id, guild_id, title, description, colour, author_name, author_icon_url, author_url, image_url, thumbnail_url, footer_text, footer_icon_url, timestamp
FROM embeds
WHERE guild_id = $1
ORDER BY id ASC
LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.Embeds, query, func(rows pgx.Rows) (database.CustomEmbed, error) {
		var embed database.CustomEmbed
		if err := rows.Scan(&embed.Id, &embed.GuildId, &embed.Title, &embed.Description, &embed.Colour, &embed.AuthorName, &embed.AuthorIconUrl, &embed.AuthorUrl, &embed.ImageUrl, &embed.ThumbnailUrl, &embed.FooterText, &embed.FooterIconUrl, &embed.Timestamp); err != nil {
			return database.CustomEmbed{}, err
		}

		return embed, nil
	})
}

func (d *Daemon) fetchExitSurveyResponses(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, form_id, question_id, response
FROM exit_survey_responses
WHERE guild_id = $1
ORDER BY ticket_id, question_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.ExitSurveyResponses, query, func(rows pgx.Rows) (dto.TicketUnion[dto.ExitSurveyResponse], error) {
		var res dto.TicketUnion[dto.ExitSurveyResponse]
		if err := rows.Scan(&res.TicketId, &res.Data.FormId, &res.Data.QuestionId, &res.Data.Response); err != nil {
			return dto.TicketUnion[dto.ExitSurveyResponse]{}, err
		}

		return res, nil
	})
}

func (d *Daemon) fetchFeedbackEnabled(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.FeedbackEnabled, d.database.FeedbackEnabled.Get)
}

func (d *Daemon) fetchFirstResponseTimes(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, user_id, response_time
FROM first_response_time
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.FirstResponseTimes, query, func(rows pgx.Rows) (dto.FirstResponseTime, error) {
		var res dto.FirstResponseTime
		if err := rows.Scan(&res.TicketId, &res.UserId, &res.ResponseTime); err != nil {
			return dto.FirstResponseTime{}, err
		}

		return res, nil
	})
}

func (d *Daemon) fetchFormInputs(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT i.id, i.form_id, i.position, i.custom_id, i.style, i.label, i.placeholder, i.required, i.min_length, i.max_length
FROM form_input i
INNER JOIN forms f
ON i.form_id = f.form_id
WHERE f.guild_id = $1
ORDER BY i.id ASC
LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.FormInputs, query, func(rows pgx.Rows) (database.FormInput, error) {
		var input database.FormInput
		if err := rows.Scan(&input.Id, &input.FormId, &input.Position, &input.CustomId, &input.Style, &input.Label, &input.Placeholder, &input.Required, &input.MinLength, &input.MaxLength); err != nil {
			return database.FormInput{}, err
		}

		return input, nil
	})
}

func (d *Daemon) fetchForms(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT form_id, guild_id, title, custom_id
FROM forms
WHERE guild_id = $1
ORDER BY form_id ASC
LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.Forms, query, func(rows pgx.Rows) (database.Form, error) {
		var form database.Form
		if err := rows.Scan(&form.Id, &form.GuildId, &form.Title, &form.CustomId); err != nil {
			return database.Form{}, err
		}

		return form, nil
	})
}

func (d *Daemon) fetchGuildIsGloballyBlacklisted(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.GuildIsGloballyBlacklisted, func(ctx context.Context, guildId uint64) (bool, error) {
		isBlacklisted, _, err := d.database.ServerBlacklist.IsBlacklisted(ctx, guildId)
		return isBlacklisted, err
	})
}

func (d *Daemon) fetchGuildMetadata(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.GuildMetadata, d.database.GuildMetadata.Get)
}

func (d *Daemon) fetchMultiPanels(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.MultiPanels, d.database.MultiPanels.GetByGuild)
}

func (d *Daemon) fetchMultiPanelTargets(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT multi_panel_id, panel_id
FROM multi_panel_targets t
INNER JOIN multi_panels p
ON t.multi_panel_id = p.id
WHERE p.guild_id = $1
ORDER BY multi_panel_id, panel_id ASC LIMIT $2 OFFSET $3;
`

	type response struct {
		MultiPanelId int
		PanelId      int
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.MultiPanelId, &data.PanelId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]int)
	for _, r := range res {
		m[r.MultiPanelId] = append(m[r.MultiPanelId], r.PanelId)
	}

	guildData.MultiPanelTargets = m
	return nil
}

func (d *Daemon) fetchNamingScheme(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.NamingScheme, d.database.NamingScheme.Get)
}

func (d *Daemon) fetchOnCallUsers(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.OnCallUsers, d.database.OnCall.GetUsersOnCall)
}

func (d *Daemon) fetchPanelAccessControlRules(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.PanelAccessControlRules, d.database.PanelAccessControlRules.GetAllForGuild)
}

func (d *Daemon) fetchPanelMentionUser(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT m.panel_id, m.should_mention_user
FROM panel_user_mentions m
INNER JOIN panels p
ON m.panel_id = p.panel_id
WHERE p.guild_id = $1
ORDER BY m.panel_id ASC LIMIT $2 OFFSET $3;`

	return fetchMap(ctx, d.database, guildId, &guildData.PanelMentionUser, query)
}

func (d *Daemon) fetchPanelRoleMentions(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT m.panel_id, m.role_id
FROM panel_role_mentions m
INNER JOIN panels p
ON m.panel_id = p.panel_id
WHERE p.guild_id = $1
ORDER BY m.panel_id ASC LIMIT $2 OFFSET $3;`

	type response struct {
		PanelId int
		RoleId  uint64
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.PanelId, &data.RoleId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]uint64)
	for _, r := range res {
		m[r.PanelId] = append(m[r.PanelId], r.RoleId)
	}

	guildData.PanelRoleMentions = m
	return nil
}

func (d *Daemon) fetchPanels(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.Panels, d.database.Panel.GetByGuild)
}

func (d *Daemon) fetchPanelTeams(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT t.panel_id, t.team_id
FROM panel_teams t
INNER JOIN panels p
ON t.panel_id = p.panel_id
WHERE p.guild_id = $1
ORDER BY t.panel_id ASC LIMIT $2 OFFSET $3;`

	type response struct {
		PanelId int
		TeamId  int
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.PanelId, &data.TeamId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]int)
	for _, r := range res {
		m[r.PanelId] = append(m[r.PanelId], r.TeamId)
	}

	guildData.PanelTeams = m
	return nil
}

func (d *Daemon) fetchParticipants(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, user_id
FROM participant
WHERE guild_id = $1
ORDER BY ticket_id, user_id ASC LIMIT $2 OFFSET $3;`

	type response struct {
		TicketId int
		UserId   uint64
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.TicketId, &data.UserId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]uint64)
	for _, r := range res {
		m[r.TicketId] = append(m[r.TicketId], r.UserId)
	}

	guildData.Participants = m
	return nil
}

func (d *Daemon) fetchUserPermissions(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT user_id, support, admin
FROM permissions
WHERE guild_id = $1
ORDER BY user_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.UserPermissions, query, func(rows pgx.Rows) (dto.Permission, error) {
		var permission dto.Permission
		if err := rows.Scan(&permission.Snowflake, &permission.IsSupport, &permission.IsAdmin); err != nil {
			return dto.Permission{}, err
		}

		return permission, nil
	})
}

func (d *Daemon) fetchGuildBlacklistedRoles(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT role_id
FROM role_blacklist
WHERE guild_id = $1
ORDER BY role_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.GuildBlacklistedRoles, query, func(rows pgx.Rows) (uint64, error) {
		var roleId uint64
		if err := rows.Scan(&roleId); err != nil {
			return 0, err
		}

		return roleId, nil
	})
}

func (d *Daemon) fetchRolePermissions(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT role_id, support, admin
FROM role_permissions
WHERE guild_id = $1
ORDER BY role_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.RolePermissions, query, func(rows pgx.Rows) (dto.Permission, error) {
		var permission dto.Permission
		if err := rows.Scan(&permission.Snowflake, &permission.IsSupport, &permission.IsAdmin); err != nil {
			return dto.Permission{}, err
		}

		return permission, nil
	})
}

func (d *Daemon) fetchServiceRatings(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, rating
FROM service_ratings
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.ServiceRatings, query, func(rows pgx.Rows) (dto.TicketUnion[int16], error) {
		var res dto.TicketUnion[int16]
		if err := rows.Scan(&res.TicketId, &res.Data); err != nil {
			return dto.TicketUnion[int16]{}, err
		}

		return res, nil
	})
}

func (d *Daemon) fetchSettings(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.Settings, d.database.Settings.Get)
}

func (d *Daemon) fetchSupportTeamUsers(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT team_id, user_id
FROM support_team_members
INNER JOIN support_team
ON support_team_members.team_id = support_team.id
WHERE support_team.guild_id = $1
ORDER BY team_id, user_id ASC LIMIT $2 OFFSET $3;`

	type response struct {
		TeamId int
		UserId uint64
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.TeamId, &data.UserId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]uint64)
	for _, r := range res {
		m[r.TeamId] = append(m[r.TeamId], r.UserId)
	}

	guildData.SupportTeamUsers = m
	return nil
}

func (d *Daemon) fetchSupportTeamRoles(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT team_id, role_id
FROM support_team_roles
INNER JOIN support_team
ON support_team_roles.team_id = support_team.id
WHERE support_team.guild_id = $1
ORDER BY team_id, role_id ASC LIMIT $2 OFFSET $3;`

	type response struct {
		TeamId int
		RoleId uint64
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.TeamId, &data.RoleId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]uint64)
	for _, r := range res {
		m[r.TeamId] = append(m[r.TeamId], r.RoleId)
	}

	guildData.SupportTeamRoles = m
	return nil
}

func (d *Daemon) fetchSupportTeams(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.SupportTeams, d.database.SupportTeam.Get)
}

func (d *Daemon) fetchTags(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.Tags, func(ctx context.Context, guildId uint64) ([]database.Tag, error) {
		tags, err := d.database.Tag.GetByGuild(ctx, guildId)
		if err != nil {
			return nil, err
		}

		values := make([]database.Tag, 0, len(tags))
		for _, tag := range tags {
			values = append(values, tag)
		}

		return values, nil
	})
}

func (d *Daemon) fetchTicketClaims(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, user_id
FROM ticket_claims
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.TicketClaims, query, func(rows pgx.Rows) (dto.TicketUnion[uint64], error) {
		var res dto.TicketUnion[uint64]
		if err := rows.Scan(&res.TicketId, &res.Data); err != nil {
			return dto.TicketUnion[uint64]{}, err
		}

		return res, nil
	})
}

func (d *Daemon) fetchTicketLastMessages(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, last_message_id, last_message_time, user_id, user_is_staff
FROM ticket_last_message
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.TicketLastMessages, query, func(rows pgx.Rows) (dto.TicketUnion[database.TicketLastMessage], error) {
		var res dto.TicketUnion[database.TicketLastMessage]
		if err := rows.Scan(&res.TicketId, &res.Data.LastMessageId, &res.Data.LastMessageTime, &res.Data.UserId, &res.Data.UserIsStaff); err != nil {
			return dto.TicketUnion[database.TicketLastMessage]{}, err
		}

		return res, nil
	})
}

func (d *Daemon) fetchTicketLimit(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.TicketLimit, func(ctx context.Context, guildId uint64) (int, error) {
		limit, err := d.database.TicketLimit.Get(ctx, guildId)
		return int(limit), err
	})
}

func (d *Daemon) fetchTicketAdditionalMembers(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT ticket_id, user_id
FROM ticket_members
WHERE guild_id = $1
ORDER BY ticket_id ASC LIMIT $2 OFFSET $3;`

	type response struct {
		TicketId int
		UserId   uint64
	}

	var res []response
	if err := fetchCustomPaginated(ctx, d.database, guildId, &res, query, func(rows pgx.Rows) (response, error) {
		var data response
		if err := rows.Scan(&data.TicketId, &data.UserId); err != nil {
			return response{}, err
		}

		return data, nil
	}); err != nil {
		return err
	}

	m := make(map[int][]uint64)
	for _, r := range res {
		m[r.TicketId] = append(m[r.TicketId], r.UserId)
	}

	guildData.TicketAdditionalMembers = m
	return nil
}

func (d *Daemon) fetchTicketPermissions(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.TicketPermissions, d.database.TicketPermissions.Get)
}

func (d *Daemon) fetchTickets(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	query := `
SELECT id, guild_id, channel_id, user_id, open, open_time, welcome_message_id, panel_id, has_transcript, close_time, is_thread, join_message_id, notes_thread_id, status
FROM tickets
WHERE guild_id = $1
ORDER BY id ASC LIMIT $2 OFFSET $3;`

	return fetchCustomPaginated(ctx, d.database, guildId, &guildData.Tickets, query, func(rows pgx.Rows) (database.Ticket, error) {
		var ticket database.Ticket
		if err := rows.Scan(
			&ticket.Id,
			&ticket.GuildId,
			&ticket.ChannelId,
			&ticket.UserId,
			&ticket.Open,
			&ticket.OpenTime,
			&ticket.WelcomeMessageId,
			&ticket.PanelId,
			&ticket.HasTranscript,
			&ticket.CloseTime,
			&ticket.IsThread,
			&ticket.JoinMessageId,
			&ticket.NotesThreadId,
			&ticket.Status,
		); err != nil {
			return database.Ticket{}, err
		}

		return ticket, nil
	})
}

func (d *Daemon) fetchUsersCanClose(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchValNoPtr(ctx, guildId, &guildData.UsersCanClose, d.database.UsersCanClose.Get)
}

func (d *Daemon) fetchWelcomeMessage(ctx context.Context, guildId uint64, guildData *dto.GuildData) error {
	return fetchVal(ctx, guildId, &guildData.WelcomeMessage, d.database.WelcomeMessages.Get)
}

func fetchVal[T comparable](
	ctx context.Context,
	guildId uint64,
	ptr **T,
	f func(context.Context, uint64) (T, error),
) error {
	data, err := f(ctx, guildId)
	if err != nil {
		return err
	}

	if data == *new(T) {
		return nil
	}

	*ptr = &data
	return nil
}

func fetchValNoPtr[T any](
	ctx context.Context,
	guildId uint64,
	ptr *T,
	f func(context.Context, uint64) (T, error),
) error {
	data, err := f(ctx, guildId)
	if err != nil {
		return err
	}

	*ptr = data
	return nil
}

func fetchPtr[T any](
	ctx context.Context,
	guildId uint64,
	ptr **T,
	f func(context.Context, uint64) (*T, error),
) error {
	data, err := f(ctx, guildId)
	if err != nil {
		return err
	}

	*ptr = data
	return nil
}

const paginationLimit = 2_500

func fetchCustomPaginated[T any](
	ctx context.Context,
	db *database.Database,
	guildId uint64,
	ptr *[]T,
	query string,
	de func(rows pgx.Rows) (T, error),
) error {
	res := make([]T, 0)
	count := 0
	hasMore := true
	for hasMore {
		rows, err := db.ArchiveMessages.Query(ctx, query, guildId, paginationLimit, count)
		if err != nil {
			return err
		}

		thisCount := 0
		for rows.Next() {
			row, err := de(rows)
			if err != nil {
				return err
			}

			res = append(res, row)
			thisCount++
		}

		count += thisCount
		hasMore = thisCount == paginationLimit
	}

	*ptr = res
	return nil
}

func fetchMap[T comparable, U any](
	ctx context.Context,
	db *database.Database,
	guildId uint64,
	ptr *map[T]U,
	query string,
) error {
	res := make(map[T]U)
	var count int
	hasMore := true
	for hasMore {
		rows, err := db.ArchiveMessages.Query(ctx, query, guildId, paginationLimit, count)
		if err != nil {
			return err
		}

		thisCount := 0
		for rows.Next() {
			var key T
			var value U
			if err := rows.Scan(&key, &value); err != nil {
				return err
			}

			res[key] = value
			thisCount++
		}

		count += thisCount
		hasMore = thisCount == paginationLimit
	}

	*ptr = res
	return nil
}
