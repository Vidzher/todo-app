package core_http_response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	h.w.WriteHeader(statusCode)

	response := map[string]string{
		"message":    msg,
		"statusCode": strconv.Itoa(statusCode),
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}
