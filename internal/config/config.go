package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	JWTSecret              string
	KafkaBroker            string
	KafkaPartitions        int
	KafkaReplicationFactor int
	KafkaTopicCompany      string
}

// LoadConfig loads the configuration from environment variables or uses default values
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "postgresql://user:password@localhost:5432/dbname?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "secretkey")
	kafkaBroker := getEnv("KAFKA_BROKER", "9092")
	kafkaPartitions := getEnvAsInt("KAFKA_PARTITIONS", 3)
	kafkaReplicationFactor := getEnvAsInt("KAFKA_REPLICATION_FACTOR", 1)
	kafkaTopicCompany := getEnv("KAFKA_TOPIC_COMPANY", "company-events")

	return &Config{
		Port:                   port,
		DatabaseURL:            databaseURL,
		JWTSecret:              jwtSecret,
		KafkaBroker:            kafkaBroker,
		KafkaPartitions:        kafkaPartitions,
		KafkaReplicationFactor: kafkaReplicationFactor,
		KafkaTopicCompany:      kafkaTopicCompany,
	}
}

// getEnv retrieves the environment variable value or returns the default if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("%s not set, using default: %s", key, defaultValue)
	return defaultValue
}

// getEnvAsInt retrieves the environment variable value as an integer or returns the default if not set or invalid
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("%s not set or invalid, using default: %d", key, defaultValue)
	return defaultValue
}
