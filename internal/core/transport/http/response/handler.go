package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	core_errors "todo-api/internal/core/errors"
	core_loger "todo-api/internal/core/logger"

	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_loger.Logger
	w   http.ResponseWriter
}

func NewHttpResponseHandler(log *core_loger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{log: log, w: w}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))
	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) JSONResponse(rBody any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(rBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))

		return
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(response, statusCode)
}
