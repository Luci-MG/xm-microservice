package company

import (
	"errors"

	"github.com/google/uuid"
)

// Service defines the business logic interface for companies
type Service interface {
	CreateCompany(company *Company) error
	UpdateCompany(id uuid.UUID, company *Company) error
	DeleteCompany(id uuid.UUID) error
	GetCompanyByID(id uuid.UUID) (*Company, error)
}

type service struct {
	repo Repository
}

// NewService initializes a new company service with the repository
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateCompany validates and creates a new company
func (s *service) CreateCompany(company *Company) error {
	if err := validateCompany(company); err != nil {
		return err
	}

	company.ID = uuid.New()
	return s.repo.Create(company)
}

// UpdateCompany validates and updates an existing company
func (s *service) UpdateCompany(id uuid.UUID, company *Company) error {
	if err := validateCompany(company); err != nil {
		return err
	}
	return s.repo.Update(id, company)
}

// DeleteCompany removes a company by its ID
func (s *service) DeleteCompany(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// GetCompanyByID retrieves a company by its ID
func (s *service) GetCompanyByID(id uuid.UUID) (*Company, error) {
	return s.repo.GetByID(id)
}

// validateCompany checks the business rules for company creation and updates
func validateCompany(company *Company) error {
	if company.Name == "" || len(company.Name) > 15 {
		return errors.New("invalid company name: must be non-empty and up to 15 characters")
	}

	if company.AmountOfEmployees == nil {
		return errors.New("amount of employees is required")
	}

	if *company.AmountOfEmployees < 0 {
		return errors.New("invalid amount of employees: cannot be negative")
	}

	if company.Registered == nil {
		return errors.New("registered status is required")
	}

	if company.Type == "" {
		return errors.New("company type is required")
	}

	switch company.Type {
	case "Corporation", "NonProfit", "Cooperative", "Sole Proprietorship":
	default:
		return errors.New("invalid company type")
	}

	if len(company.Description) > 3000 {
		return errors.New("invalid description: must be up to 3000 characters")
	}

	return nil
}
