package event

import (
	"xm-microservice/pkg/logger"

	"github.com/segmentio/kafka-go"
)

// CreateTopic creates a Kafka topic if it doesn't already exist
func CreateTopic(brokerAddress, topic string, partitions int, replications int, log *logger.Logger) error {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Error(err, "Failed to connect to Kafka broker")
		return err
	}
	defer conn.Close()

	// Check if the topic already exists
	partitionsList, err := conn.ReadPartitions()
	if err == nil {
		for _, p := range partitionsList {
			if p.Topic == topic {
				log.Info("Topic '%s' already exists", topic)
				return nil
			}
		}
	}

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replications,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		log.Error(err, "Failed to create Kafka topic")
		return err
	}

	log.Info("Kafka topic '%s' created successfully", topic)
	return nil
}
