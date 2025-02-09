package company

import (
	"encoding/json"
	"net/http"
	"strings"

	"xm-microservice/internal/event"
	"xm-microservice/pkg/logger"
	"xm-microservice/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	service  Service
	producer *event.Producer
	logger   *logger.Logger
}

// NewHandler initializes the company handler with the service, Kafka producer, and logger
func NewHandler(service Service, producer *event.Producer, logger *logger.Logger) *Handler {
	return &Handler{
		service:  service,
		producer: producer,
		logger:   logger,
	}
}

// Event represents the structure of Kafka messages
type Event struct {
	Action  string      `json:"action"`
	Company interface{} `json:"company"`
}

// CreateCompany handles the creation of a new company
func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CreateCompany handler invoked")

	var company Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		h.logger.Error(err, "Invalid input while decoding company data")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := h.service.CreateCompany(&company); err != nil {
		h.logger.Error(err, "Failed to create company")

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			utils.ErrorResponse(w, http.StatusBadRequest, "Company name already exists")
			return
		}

		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.produceEvent("create", company)
	h.logger.Info("Company created successfully with ID: %s", company.ID)
	utils.JSONResponse(w, http.StatusCreated, company)
}

// UpdateCompany handles updating an existing company
func (h *Handler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UpdateCompany handler invoked")

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error(err, "Invalid UUID")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var company Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		h.logger.Error(err, "Invalid input while decoding company data")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	company.ID = id
	if err := h.service.UpdateCompany(id, &company); err != nil {
		h.logger.Error(err, "Failed to update company")
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.produceEvent("update", company)
	h.logger.Info("Company updated successfully with ID: %s", id)
	utils.JSONResponse(w, http.StatusOK, company)
}

// DeleteCompany handles deleting an existing company
func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DeleteCompany handler invoked")

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error(err, "Invalid UUID")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.service.DeleteCompany(id); err != nil {
		h.logger.Error(err, "Failed to delete company")
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.produceEvent("delete", map[string]string{"id": id.String()})
	h.logger.Info("Company deleted successfully with ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

// GetCompany retrieves a company by its ID
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("GetCompany handler invoked")

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error(err, "Invalid UUID")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	company, err := h.service.GetCompanyByID(id)
	if err != nil {
		h.logger.Error(err, "Company not found")
		utils.ErrorResponse(w, http.StatusNotFound, "Company not found")
		return
	}

	h.logger.Info("Company retrieved successfully with ID: %s", id)
	utils.JSONResponse(w, http.StatusOK, company)
}

// produceEvent sends events to Kafka based on the action performed
func (h *Handler) produceEvent(action string, company interface{}) {
	event := Event{
		Action:  action,
		Company: company,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		h.logger.Error(err, "Failed to marshal event data")
		return
	}

	if err := h.producer.PublishMessage(action, string(eventData)); err != nil {
		h.logger.Error(err, "Failed to publish Kafka message")
	} else {
		h.logger.Info("Kafka message published successfully for action: %s", action)
	}
}
