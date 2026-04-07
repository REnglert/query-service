package api

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)
	mux.HandleFunc("/query", h.Query) 
}
