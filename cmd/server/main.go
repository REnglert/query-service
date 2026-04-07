package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"query-service/internal/config"
	"query-service/internal/server"
	"query-service/internal/shutdown"
)

func main() {
	ctx := context.Background()

	shutdownCtx := shutdown.WaitForShutdown(ctx)

	// ---- Load configuration ----
	cfg, err := config.Load()
	if err != nil {

		panic(err)
	}

	// ---- Initialize logger ----
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	logger.Info("starting service",
		"service", cfg.ServiceName,
		"port", cfg.Port,
	)

	// ---- Create HTTP server ----
	srv := server.New(cfg, logger)

	// ---- Run server in background ----
	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("server stopped with error", "error", err)
			os.Exit(1)
		}
	}()

	// ---- Block until shutdown signal ----
	<-shutdownCtx.Done()
	logger.Info("shutdown signal received")

	// ---- Perform graceful shutdown with timeout ----
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("service shut down cleanly")

}
