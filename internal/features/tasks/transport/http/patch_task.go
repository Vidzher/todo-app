package tasks_transport_http

import (
	"net/http"
	"todo-api/internal/core/domain"
	core_loger "todo-api/internal/core/logger"
	core_http_request "todo-api/internal/core/transport/http/request"
	core_http_response "todo-api/internal/core/transport/http/response"
	core_http_utils "todo-api/internal/core/transport/http/utils"
)

type PatchTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
	Completed   *bool   `json:"completed" validate:"omitempty"`
}

type PatchTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) PatchTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_loger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	var req PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	taskPatch := taskPatchFromRequest(req)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(req PatchTaskRequest) domain.TaskPatch {
	return domain.TaskPatch{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}
}
