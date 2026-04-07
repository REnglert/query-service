package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"query-service/internal/api"
	"query-service/internal/config"
	"query-service/internal/llm"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	llmClient := llm.NewRealClient(cfg.LLMBaseURL)
	handler := api.NewHandler(llmClient)

	api.RegisterRoutes(mux, handler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("http server listening", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down http server")
	return s.httpServer.Shutdown(ctx)
}
