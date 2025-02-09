package company

import (
	"github.com/google/uuid"
)

type CompanyType string

const (
	Corporation        CompanyType = "Corporation"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

type Company struct {
	ID                uuid.UUID   `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description,omitempty"`
	AmountOfEmployees *int        `json:"amount_of_employees"`
	Registered        *bool       `json:"registered"`
	Type              CompanyType `json:"type"`
}
