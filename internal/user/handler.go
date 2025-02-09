package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"xm-microservice/pkg/logger"
	"xm-microservice/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UserResponse represents the sanitized user data to return in API responses
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// Handler manages HTTP requests related to user operations
type Handler struct {
	service *Service
	logger  *logger.Logger
}

// NewHandler initializes a new Handler for user services
func NewHandler(service *Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// sanitizeUser removes sensitive fields like password_hash from the response
func sanitizeUser(user *User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}

// CreateUser handles the creation of a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CreateUser handler invoked")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err, "Invalid input while decoding user data")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if req.Username == "" || req.Password == "" {
		h.logger.Info("Validation failed: Username or password is empty")
		utils.ErrorResponse(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	user, err := h.service.CreateUser(req.Username, req.Password)
	if err != nil {
		h.logger.Error(err, "Failed to create user")

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			utils.ErrorResponse(w, http.StatusBadRequest, "User already exists")
			return
		}

		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.logger.Info("User created successfully with ID: %s", user.ID)
	utils.JSONResponse(w, http.StatusCreated, sanitizeUser(user))
}

// GetUser retrieves a user by their unique ID
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetUser handler invoked")

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error(err, "Invalid ID")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		h.logger.Error(err, "User not found")
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	h.logger.Info("User retrieved successfully with ID: %s", id)
	utils.JSONResponse(w, http.StatusOK, sanitizeUser(user))
}

// UpdateUser handles updating an existing user's details
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UpdateUser handler invoked")

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error(err, "Invalid ID")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err, "Invalid input while decoding user data")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if req.Username == "" || req.Password == "" {
		h.logger.Info("Validation failed: Username or password is empty")
		utils.ErrorResponse(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	if err := h.service.UpdateUser(id, req.Username, req.Password); err != nil {
		h.logger.Error(err, "Failed to update user")
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.logger.Info("User updated successfully with ID: %s", id)
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// DeleteUser handles the deletion of a user by their ID
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DeleteUser handler invoked")

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error(err, "Invalid ID")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		h.logger.Error(err, "Failed to delete user")
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.logger.Info("User deleted successfully with ID: %s", id)
	utils.JSONResponse(w, http.StatusNoContent, nil)
}
