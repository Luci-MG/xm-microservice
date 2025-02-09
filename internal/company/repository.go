package company

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(company *Company) error
	Update(id uuid.UUID, company *Company) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*Company, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository initializes a new company repository with the database connection
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new company into the database
func (r *repository) Create(company *Company) error {
	query := `INSERT INTO companies (id, name, description, amount_of_employees, registered, type) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, company.ID, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type)
	return err
}

// Update modifies the details of an existing company
func (r *repository) Update(id uuid.UUID, company *Company) error {
	query := `UPDATE companies SET name=$1, description=$2, amount_of_employees=$3, registered=$4, type=$5 WHERE id=$6`
	_, err := r.db.Exec(query, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type, id)
	return err
}

// Delete removes a company from the database based on its ID
func (r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM companies WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetByID retrieves a company by its unique identifier
func (r *repository) GetByID(id uuid.UUID) (*Company, error) {
	query := `SELECT id, name, description, amount_of_employees, registered, type FROM companies WHERE id=$1`
	company := &Company{}
	err := r.db.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&company.Description,
		&company.AmountOfEmployees,
		&company.Registered,
		&company.Type,
	)
	if err != nil {
		return nil, err
	}
	return company, nil
}
