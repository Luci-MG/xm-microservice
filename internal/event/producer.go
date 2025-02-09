package event

import (
	"context"

	"xm-microservice/pkg/logger"

	"github.com/segmentio/kafka-go"
)

// Producer represents a Kafka message producer
type Producer struct {
	writer *kafka.Writer
	log    *logger.Logger
}

// NewProducer initializes a new Kafka producer
func NewProducer(brokerAddress, topic string, log *logger.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: writer, log: log}
}

// PublishMessage sends a message to the Kafka topic
func (p *Producer) PublishMessage(key, value string) error {
	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	if err := p.writer.WriteMessages(context.Background(), message); err != nil {
		p.log.Error(err, "Failed to publish message to Kafka")
		return err
	}

	p.log.Info("Message published successfully: key=%s, value=%s", key, value)
	return nil
}

// Close gracefully closes the Kafka writer
func (p *Producer) Close() error {
	err := p.writer.Close()
	if err != nil {
		p.log.Error(err, "Failed to close Kafka producer")
		return err
	}

	p.log.Info("Kafka producer closed successfully")
	return nil
}
