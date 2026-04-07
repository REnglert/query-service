package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"query-service/internal/llm"
)

type Handler struct {
	llmClient llm.Client
}

func NewHandler(llmClient llm.Client) *Handler {
	return &Handler{
		llmClient: llmClient,
	}
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Result string `json:"result"`
}


func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func (h *Handler) Ready(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready\n"))
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Use a context with timeout for the LLM call
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result, err := h.llmClient.Query(ctx, req.Query)
	if err != nil {
		http.Error(w, "LLM query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := QueryResponse{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
