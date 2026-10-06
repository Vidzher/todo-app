package users_transport_http

import (
	"net/http"
	"todo-api/internal/core/domain"
	core_loger "todo-api/internal/core/logger"
	core_http_request "todo-api/internal/core/transport/http/request"
	core_http_response "todo-api/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"fullName" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phoneNumber" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"fullName"`
	PhoneNumber *string `json:"phoneNumber"`
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_loger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	var req CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")

		return
	}

	userDomain := domainFromDTO(req)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := dtoFromDomain(userDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}

func dtoFromDomain(user domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}
