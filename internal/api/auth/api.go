package auth

import (
	"github.com/TicketsBot/data-self-service/internal/api"
	"net/http"
	"time"
)

type API struct {
	*api.Core
	client *http.Client
}

func NewAPI(core *api.Core) *API {
	return &API{
		Core: core,
		client: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}
