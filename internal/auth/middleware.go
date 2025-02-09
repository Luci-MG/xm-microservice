package auth

import (
	"net/http"
	"strings"
)

type Middleware struct {
	jwtService JWTService
}

// NewMiddleware initializes the JWT middleware with a secret key
func NewMiddleware(secretKey string) *Middleware {
	return &Middleware{
		jwtService: NewJWTService(secretKey),
	}
}

// GetJWTService returns the JWT service instance
func (m *Middleware) GetJWTService() JWTService {
	return m.jwtService
}

// ProtectMiddleware is a middleware function that protects routes by validating JWT tokens
func (m *Middleware) ProtectMiddleware(next http.Handler) http.Handler {
	return m.Protect(next.ServeHTTP)
}

// Protect validates the JWT token from the Authorization header
func (m *Middleware) Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized - no token provided", http.StatusUnauthorized)
			return
		}

		token, err := m.jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// extractToken extracts the JWT token from the Authorization header
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
