package oauth2

import (
	"github.com/TicketsBot/export/internal/api"
)

type API struct {
	*api.Core
}

func NewAPI(core *api.Core) *API {
	return &API{
		Core: core,
	}
}
