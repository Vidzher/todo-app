package auth_transport_http

import (
	"context"
	"net/http"
	core_http_server "todo-api/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	AuthCodeURL(state string) string
	CompleteLogin(ctx context.Context, code, provider string) (any, error)
}

func NewAuthHTTPHandler(s AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{authService: s}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/auth/google/login",
			Handler: h.GoogleLogin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/auth/google/callback",
			Handler: h.GoogleCallback,
		},
	}
}
