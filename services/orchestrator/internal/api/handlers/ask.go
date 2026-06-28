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

func (h *Handler) AskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req AskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.SessionID) == "" ||
		strings.TrimSpace(req.Owner) == "" ||
		strings.TrimSpace(req.Repo) == "" ||
		strings.TrimSpace(req.Content) == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
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
		http.Error(w, "failed to generate response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(AskResponse{Response: res}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
