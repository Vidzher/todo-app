package tasks_transport_http

import (
	"net/http"
	core_loger "todo-api/internal/core/logger"
	core_http_response "todo-api/internal/core/transport/http/response"
	core_http_utils "todo-api/internal/core/transport/http/utils"
)

type GetTaskReponse TaskDTOResponse

func (h *TasksHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_loger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	taskId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskId path value")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	response := GetTaskReponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
