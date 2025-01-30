package router

import (
	"github.com/TicketsBot/data-self-service/internal/api"
	"github.com/TicketsBot/data-self-service/internal/api/auth"
	"github.com/TicketsBot/data-self-service/internal/api/health"
	"github.com/TicketsBot/data-self-service/internal/api/middleware"
	"github.com/TicketsBot/data-self-service/internal/api/requests"
	"github.com/TicketsBot/data-self-service/internal/config"
	"github.com/TicketsBot/data-self-service/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	slogchi "github.com/samber/slog-chi"
	"log/slog"
	"net/http"
)

func New(logger *slog.Logger, config config.ApiConfig, repository *repository.Repository) *chi.Mux {
	core := api.NewCore(logger, config, repository)

	r := chi.NewRouter()

	// Middleware
	r.Use(slogchi.New(logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	// Routes
	r.Get("/livez", health.Livez)

	// /auth
	r.Group(func(r chi.Router) {
		api := auth.NewAPI(core)

		r.Post("/auth/exchange", api.Exchange)
		r.Post("/auth/guilds", api.FetchGuilds)
	})

	// /requests
	r.Group(func(r chi.Router) {
		api := requests.NewAPI(core)

		r.Use(middleware.Authenticate(core))

		r.Get("/requests", api.ListRequests)
		r.Post("/requests", api.CreateRequest)
	})

	return r
}
