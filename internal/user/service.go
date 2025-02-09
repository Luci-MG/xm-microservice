package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service handles business logic related to users
type Service struct {
	repo *Repository
}

// NewService initializes a new Service with a user repository
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateUser handles the creation of a new user, including password hashing
func (s *Service) CreateUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their unique ID
func (s *Service) GetUserByID(id uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(id)
}

// UpdateUser updates an existing user's username and password
func (s *Service) UpdateUser(id uuid.UUID, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.UpdateUser(id, user)
}

// DeleteUser removes a user from the system by their ID
func (s *Service) DeleteUser(id uuid.UUID) error {
	return s.repo.DeleteUser(id)
}
