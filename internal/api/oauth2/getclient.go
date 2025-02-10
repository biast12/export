package oauth2

import (
	"context"
	"github.com/TicketsBot/export/internal/api"
	"github.com/TicketsBot/export/internal/model"
	"github.com/TicketsBot/export/internal/repository"
	"github.com/TicketsBot/export/pkg/dto"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) GetClient(w http.ResponseWriter, r *http.Request) {
	clientId := chi.URLParam(r, "client_id")
	if clientId == "" {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid client ID"))
		return
	}

	client, err := repository.Exec1(r.Context(), a.Repository,
		func(ctx context.Context, tx repository.TransactionContext) (*model.OAuth2Client, error) {
			return tx.OAuth2().GetClient(ctx, clientId)
		})
	if err != nil {
		a.HandleError(r.Context(), w, api.NewError(err, http.StatusInternalServerError, "Failed to get client"))
		return
	}

	if client == nil {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusNotFound, "Client not found"))
		return
	}

	a.RespondJson(w, http.StatusOK, dto.OAuth2Client{
		ClientId: client.ClientId,
		Label:    client.Label,
	})
}
