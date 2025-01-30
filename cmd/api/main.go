package main

import (
	"context"
	"github.com/TicketsBot/data-self-service/internal/api/router"
	"github.com/TicketsBot/data-self-service/internal/config"
	"github.com/TicketsBot/data-self-service/internal/repository"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New[config.ApiConfig]()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	setupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repository, err := repository.Connect(setupCtx, cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	cancel()

	router := router.New(logger, cfg, repository)

	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	closed := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch

		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server failed to shut down", "error", err)
			os.Exit(1)
		}

		close(closed)
	}()

	logger.Info("Starting server...", "address", cfg.Server.Address)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	<-closed
}
