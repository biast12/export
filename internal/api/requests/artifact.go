package requests

import (
	"context"
	"github.com/TicketsBot/data-self-service/internal/api"
	"github.com/TicketsBot/data-self-service/internal/metrics"
	"github.com/TicketsBot/data-self-service/internal/model"
	"github.com/TicketsBot/data-self-service/internal/repository"
	"github.com/TicketsBot/data-self-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func (a *API) GetArtifact(w http.ResponseWriter, r *http.Request) {
	userId := a.userId(r.Context())
	ownedGuilds := a.ownedGuilds(r.Context())

	requestId, err := uuid.Parse(chi.URLParam(r, "requestId"))
	if err != nil {
		a.RespondJson(w, http.StatusBadRequest, utils.Map{
			"error": "Invalid request ID",
		})
		return
	}

	var request *model.RequestWithArtifact
	if err := a.Repository.Tx(r.Context(), func(ctx context.Context, tx repository.TransactionContext) (err error) {
		request, err = tx.Requests().GetById(ctx, requestId)
		return
	}); err != nil {
		a.HandleError(r.Context(), w, api.NewError(err, http.StatusInternalServerError, "Failed to fetch request"))
		return
	}

	if request == nil || request.Artifact == nil {
		a.RespondJson(w, http.StatusNotFound, utils.Map{
			"error": "Data export not found",
		})
		return
	}

	if request.Request.UserId != userId {
		a.RespondJson(w, http.StatusForbidden, utils.Map{
			"error": "You do not own this request",
		})
		return
	}

	if request.Artifact.ExpiresAt.Before(time.Now()) {
		a.RespondJson(w, http.StatusGone, utils.Map{
			"error": "Artifact has expired",
		})
		return
	}

	if request.Request.GuildId == nil || !utils.Contains(ownedGuilds, *request.Request.GuildId) {
		a.RespondJson(w, http.StatusForbidden, utils.Map{
			"error": "User does not own this guild",
		})
		return
	}

	a.Logger.Info("Fetching artifact", "requestId", requestId, "user_id", userId)

	// Get artifact
	bytes, err := a.Artifacts.Fetch(r.Context(), request.Artifact.RequestId, request.Artifact.Key)
	if err != nil {
		a.HandleError(r.Context(), w, api.NewError(err, http.StatusInternalServerError, "Failed to fetch artifact"))
		return
	}

	metrics.ArtifactsDownloaded.WithLabelValues(request.Request.Type.String()).Inc()
	metrics.ArtifactsDownloadedBytes.WithLabelValues(request.Request.Type.String()).Add(float64(len(bytes)))

	w.Header().Add("Content-Type", "application/zip")
	w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Add("Content-Disposition", "attachment; filename=transcripts.zip")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
