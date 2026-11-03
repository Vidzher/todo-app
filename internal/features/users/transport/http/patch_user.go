package users_transport_http

import (
	"net/http"
	"todo-api/internal/core/domain"
	core_loger "todo-api/internal/core/logger"
	core_http_request "todo-api/internal/core/transport/http/request"
	core_http_response "todo-api/internal/core/transport/http/response"
	core_http_utils "todo-api/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    *string `json:"fullName" validate:"omitempty,min=3,max=100"`
	PhoneNumber *string `json:"phoneNumber" validate:"omitempty,min=10,max=15,startswith=+"`
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_loger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	var req PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(req PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}
}
