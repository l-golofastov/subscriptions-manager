// @title Subscriptions Manager API
// @version 1.0
// @description API for managing user subscriptions

// @contact.name Lev Golofastov
// @contact.email l.golofastov@mail.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /subscriptions
// @schemes http
package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/l-golofastov/subscriptions-manager/internal/config"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/subscriptions"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/sum"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/middleware"
	"github.com/l-golofastov/subscriptions-manager/internal/repository/postgres"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/l-golofastov/subscriptions-manager/docs"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger()

	log.Info("starting application")

	storage, err := postgres.NewStoragePostgres(cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer storage.Close()

	log.Info("connected to database")

	mux := http.NewServeMux()

	mux.HandleFunc("/subscriptions", subscriptions.NewSubscriptionsHandler(log, storage))
	mux.HandleFunc("/subscriptions/", subscriptions.NewSubscriptionByIDHandler(log, storage))
	mux.HandleFunc("/subscriptions/sum", sum.NewSumHandler(log, storage))
	mux.Handle("/swagger", httpSwagger.WrapHandler)

	var handler http.Handler = mux
	handler = middleware.NewLoggingMiddleware(handler, log)
	handler = middleware.NewRequestIDMiddleware(handler)
	handler = middleware.NewRecovererMiddleware(handler)

	log.Info("starting server")

	srv := &http.Server{
		Addr:         cfg.HTTPServer.ServerAddress,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start server", "error", err)
	}
}

func setupLogger() *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	return logger
}
