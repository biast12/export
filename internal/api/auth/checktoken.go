package auth

import (
	"github.com/TicketsBot/export/internal/utils"
	"net/http"
)

func (a *API) CheckToken(w http.ResponseWriter, r *http.Request) {
	a.RespondJson(w, http.StatusOK, utils.Map{
		"user_id": a.UserId(r.Context()),
	})
}
