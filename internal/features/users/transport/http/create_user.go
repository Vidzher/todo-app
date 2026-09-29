package users_transport_http

import (
	"encoding/json"
	"net/http"
	core_loger "todo-api/internal/core/logger"
)

type CreateUserRequest struct {
	FullName    string  `json:"fullName"`
	PhoneNumber *string `json:"phoneNumber"`
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

	log.Debug("invoke CreateUsers handler")

	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
