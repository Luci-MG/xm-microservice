package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{secretKey: secretKey}
}

// GenerateToken generates a JWT token for a given user ID with additional claims
func (j *jwtService) GenerateToken(userID string) (string, error) {
	createdAt := time.Now().Unix()
	expiresAt := time.Now().Add(2 * time.Hour).Unix() // Token expires in 2 hours

	claims := jwt.MapClaims{
		"user_id":    userID,
		"created_at": createdAt,
		"expires_at": expiresAt,
		"exp":        expiresAt, // Standard expiration claim
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates the provided JWT token and returns the parsed token
func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})
}
