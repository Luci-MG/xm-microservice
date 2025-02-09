package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"xm-microservice/internal/user"
	"xm-microservice/pkg/logger"
	"xm-microservice/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	jwtService JWTService
	userRepo   *user.Repository
	logger     *logger.Logger
}

// NewAuthHandler initializes a new AuthHandler
func NewAuthHandler(jwtService JWTService, userRepo *user.Repository, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		userRepo:   userRepo,
		logger:     logger,
	}
}

// Login handles user authentication and JWT token generation
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Login handler invoked")

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.logger.Error(err, "Invalid input while decoding credentials")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// Fetch user from the database
	userData, err := h.userRepo.GetUserByUsername(creds.Username)
	if err != nil {
		h.logger.Error(err, "Unauthorized - user not found")
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized - user not found")
		return
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(creds.Password)); err != nil {
		h.logger.Error(err, "Unauthorized - invalid password")
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized - invalid password")
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(userData.Username)
	if err != nil {
		h.logger.Error(err, "Failed to generate token")
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	h.logger.Info("JWT token generated successfully for user: %s", userData.Username)

	// Extract token claims to include created_at and expires_at
	createdAt := time.Now().Unix()
	expiresAt := createdAt + (2 * 3600) // 2 hours from now

	// Return the JWT token with metadata
	response := map[string]interface{}{
		"token":      token,
		"created_at": createdAt,
		"expires_at": expiresAt,
	}

	utils.JSONResponse(w, http.StatusOK, response)
}
