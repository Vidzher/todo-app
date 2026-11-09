package auth_transport_http

import (
	"errors"
	"net/http"
	core_loger "todo-api/internal/core/logger"
	core_http_response "todo-api/internal/core/transport/http/response"

	"go.uber.org/zap"
)

func (h *AuthHTTPHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_loger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Error("google sent empty code", zap.String("code", code))
		responseHandler.ErrorResponse(errors.New("google sent empty code"), "empty code")
		return
	}

	user, err := h.authService.CompleteLogin(ctx, code)
	if err != nil {
		responseHandler.ErrorResponse(err, "get user info")
	}

	responseHandler.JSONResponse(user, http.StatusOK)
}
