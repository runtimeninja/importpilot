package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/runtimeninja/importpilot/internal/config"
	"github.com/runtimeninja/importpilot/internal/db"
	"github.com/runtimeninja/importpilot/internal/observability"
)

func main() {
	cfg := config.Load()

	logger := observability.NewLogger(cfg.AppEnv)
	observability.InitGlobalLogger(logger)

	slog.Info("application starting",
		"env", cfg.AppEnv,
		"port", cfg.AppPort,
	)

	postgresDB, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		log.Fatal(err)
	}
	defer postgresDB.Close()

	mux := http.NewServeMux()
	registerHealthRoutes(mux, postgresDB)

	err = http.ListenAndServe(":"+cfg.AppPort, mux)
	if err != nil {
		slog.Error("server failed", "error", err)
		log.Fatal(err)
	}
}

func registerHealthRoutes(mux *http.ServeMux, postgresDB *sql.DB) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := postgresDB.PingContext(ctx); err != nil {
			slog.Error("health check failed: database unreachable", "error", err)
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}

		slog.Debug("health check called")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
