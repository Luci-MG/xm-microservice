package user

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// Repository handles database operations related to users
type Repository struct {
	db *sql.DB
}

// NewRepository initializes a new Repository with a database connection
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user into the database
func (r *Repository) CreateUser(user *User) error {
	query := `INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.ID, user.Username, user.PasswordHash)
	return err
}

// GetUserByID retrieves a user by their unique ID
func (r *Repository) GetUserByID(id uuid.UUID) (*User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user's information
func (r *Repository) UpdateUser(id uuid.UUID, user *User) error {
	query := `UPDATE users SET username = $1, password_hash = $2 WHERE id = $3`
	_, err := r.db.Exec(query, user.Username, user.PasswordHash, id)
	return err
}

// DeleteUser removes a user from the database by their ID
func (r *Repository) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetUserByUsername retrieves a user by their username
func (r *Repository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	row := r.db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
