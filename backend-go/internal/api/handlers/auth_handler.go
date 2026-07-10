package handlers

import (
	"encoding/json"
	"net/http"
	"vk-search/internal/domain"
)

type AuthHandler struct {
	authUC domain.AuthUseCase
}

func NewAuthHandler(authUC domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authUC.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}

func (h *AuthHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}