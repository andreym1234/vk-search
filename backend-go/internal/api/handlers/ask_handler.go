package handlers

import (
	"encoding/json"
	"net/http"

	"vk-search/internal/domain"
)

type AskHandler struct {
	useCase domain.AskUseCase
}

func NewAskHandler(useCase domain.AskUseCase) *AskHandler {
	return &AskHandler{useCase: useCase}
}

func (h *AskHandler) Ask(w http.ResponseWriter, r *http.Request) {
	var req domain.AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		return
	}

	if req.Query == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "query cannot be empty"}`))
		return
	}

	userID, _ := r.Context().Value(domain.UserIDKey).(int64)

	res, err := h.useCase.Ask(r.Context(), req, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to process ask request"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}