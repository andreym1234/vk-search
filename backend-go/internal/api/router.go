package api

import (
	"net/http"

	"vk-search/internal/api/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	return mux
}
