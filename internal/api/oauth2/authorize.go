package oauth2

import (
	"context"
	"github.com/TicketsBot/export/internal/api"
	"github.com/TicketsBot/export/internal/api/constants"
	"github.com/TicketsBot/export/internal/model"
	"github.com/TicketsBot/export/internal/repository"
	"github.com/TicketsBot/export/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (a *API) Authorize(w http.ResponseWriter, r *http.Request) {
	// Validate response_type
	if r.URL.Query().Get("response_type") != "code" {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid response type"))
		return
	}

	clientId := r.URL.Query().Get("client_id")
	if clientId == "" {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid client ID"))
		return
	}

	redirectUri := r.URL.Query().Get("redirect_uri")
	if redirectUri == "" {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid redirect URI"))
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

	// Validate redirect_uri
	validRedirectUri, err := repository.Exec1(r.Context(), a.Repository,
		func(ctx context.Context, tx repository.TransactionContext) (bool, error) {
			return tx.OAuth2().ValidateRedirectUri(ctx, clientId, redirectUri)
		})
	if err != nil {
		a.HandleError(r.Context(), w, api.NewError(err, http.StatusInternalServerError, "Failed to validate redirect URI"))
		return
	}

	if !validRedirectUri {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid redirect URI"))
		return
	}

	scopes := strings.Split(r.URL.Query().Get("scope"), " ")
	if len(scopes) == 0 {
		a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid scope"))
		return
	}

	ownedGuilds := a.OwnedGuilds(r.Context())
	authorities := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		split := strings.Split(scope, ":")
		if len(split) != 2 {
			a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid scope"))
			return
		}

		guildId, err := strconv.ParseUint(split[0], 10, 64)
		if err != nil {
			a.HandleError(r.Context(), w, api.NewError(err, http.StatusBadRequest, "Invalid scope"))
			return
		}

		if !utils.Contains(ownedGuilds, guildId) {
			a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid scope"))
			return
		}

		authority := split[1]
		if !utils.Contains(constants.Scopes, authority) {
			a.HandleError(r.Context(), w, api.NewError(nil, http.StatusBadRequest, "Invalid scope"))
			return
		}

		authorities = append(authorities, authority)
	}

	code := model.OAuth2CodeData{
		Code:      utils.RandomString(32),
		ClientId:  client.ClientId,
		UserId:    a.UserId(r.Context()),
		CreatedAt: time.Now(),
	}

	if err := a.Repository.Tx(r.Context(), func(ctx context.Context, tx repository.TransactionContext) error {
		if err := tx.OAuth2().CreateCode(ctx, code); err != nil {
			return err
		}

		for _, authority := range authorities {
			if err := tx.OAuth2().CreateCodeAuthority(ctx, code.Code, authority); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		a.HandleError(r.Context(), w, api.NewError(err, http.StatusInternalServerError, "Failed to create code"))
		return
	}

	a.RespondJson(w, http.StatusOK, utils.Map{
		"code": code.Code,
	})
}
