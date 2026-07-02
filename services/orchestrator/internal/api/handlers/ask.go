package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type AskRequest struct {
	SessionID string `json:"session_id"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Content   string `json:"content"`
}

type AskResponse struct {
	Response string `json:"response"`
}

// AskHandler godoc
// @Summary Ask a question to the RAG model
// @Description Ask a question to the RAG model and get a response
// @Tags ask
// @Accept json
// @Produce json
// @Param request body AskRequest true "Ask request"
// @Success 200 {object} AskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/ask [post]
func (h *Handler) AskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req AskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.SessionID) == "" ||
		strings.TrimSpace(req.Owner) == "" ||
		strings.TrimSpace(req.Repo) == "" ||
		strings.TrimSpace(req.Content) == "" {
		WriteError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	res, err := h.rag.GenerateResponse(
		r.Context(),
		req.SessionID,
		req.Owner,
		req.Repo,
		req.Content,
	)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to generate response")
		return
	}

	if err := json.NewEncoder(w).Encode(AskResponse{Response: res}); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}
}
