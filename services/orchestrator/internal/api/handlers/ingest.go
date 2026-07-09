package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LucasM4r/repomind/internal/ingestor"
	// Importado para reconhecer a interface do Fetcher
)

type IngestRequest struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Provider string `json:"provider"`
}

type IngestResponse struct {
	Message  string `json:"message"`
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Provider string `json:"provider"`
}

// IngestHandler godoc
// @Summary Ingest a repository
// @Description Enqueue a repository for processing
// @Tags ingest
// @Accept json
// @Produce json
// @Param request body IngestRequest true "Ingest request body"
// @Success 202 {object} IngestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/ingest [post]
func (h *Handler) IngestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req IngestRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Owner == "" || req.Repo == "" || req.Provider == "" {
		WriteError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	providerClient, exists := h.enqueuer.GetProvider(req.Provider)
	if !exists {
		WriteError(w, http.StatusBadRequest, "unsupported provider: choose 'github'")
		return
	}

	h.enqueuer.EnqueueJob(ingestor.Job{
		Owner:    req.Owner,
		Repo:     req.Repo,
		Provider: providerClient,
	})

	res := IngestResponse{
		Message:  "Repository ingestion successfully enqueued in background",
		Owner:    req.Owner,
		Repo:     req.Repo,
		Provider: req.Provider,
	}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return
	}
}
