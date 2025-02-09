package main

import (
	"net/http"

	"xm-microservice/internal/auth"
	"xm-microservice/internal/company"
	"xm-microservice/internal/config"
	"xm-microservice/internal/database"
	"xm-microservice/internal/event"
	"xm-microservice/internal/health"
	"xm-microservice/internal/user"
	"xm-microservice/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting the application...")

	// Load application configuration
	cfg := config.LoadConfig()

	// Connect to the database
	db, err := database.Connect(cfg.DatabaseURL, appLogger)
	if err != nil {
		appLogger.Fatal(err)
	}
	defer db.Close()
	appLogger.Info("Connected to the database successfully")

	// Create Kafka topic and initialize producer
	kafkaTopic := cfg.KafkaTopicCompany
	if err := event.CreateTopic(cfg.KafkaBroker, kafkaTopic, cfg.KafkaPartitions, cfg.KafkaReplicationFactor, appLogger); err != nil {
		appLogger.Fatal(err)
	}
	appLogger.Info("Kafka topic created successfully")

	kafkaProducer := event.NewProducer(cfg.KafkaBroker, kafkaTopic, appLogger)

	// Initialize Company service and handler with logger
	companyRepo := company.NewRepository(db)
	companyService := company.NewService(companyRepo)
	companyHandler := company.NewHandler(companyService, kafkaProducer, appLogger)

	// Initialize User service and handler with logger
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, appLogger)

	// Initialize Authentication middleware and handler with logger
	authMiddleware := auth.NewMiddleware(cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authMiddleware.GetJWTService(), userRepo, appLogger)

	// Set up the HTTP router
	router := mux.NewRouter()

	// Register health check endpoint
	router.HandleFunc("/health", health.HealthHandler).Methods("GET")

	// Authentication route
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Public user registration route
	userRoutes := router.PathPrefix("/api/users").Subrouter()
	userRoutes.HandleFunc("", userHandler.CreateUser).Methods("POST")

	// Public route for retrieving company details
	router.HandleFunc("/api/companies/{id}", companyHandler.GetCompany).Methods("GET")

	// Protected routes for creating, updating, and deleting companies
	companyRoutes := router.PathPrefix("/api/companies").Subrouter()
	companyRoutes.Use(authMiddleware.ProtectMiddleware)
	companyRoutes.HandleFunc("", companyHandler.CreateCompany).Methods("POST")
	companyRoutes.HandleFunc("/{id}", companyHandler.UpdateCompany).Methods("PATCH")
	companyRoutes.HandleFunc("/{id}", companyHandler.DeleteCompany).Methods("DELETE")

	// Start the HTTP server
	appLogger.Info("Server is running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		appLogger.Fatal(err)
	}
}
