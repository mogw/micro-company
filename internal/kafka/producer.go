package kafka

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(topic string, key, value []byte) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *producer) Produce(topic string, key, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

func (p *producer) Close() error {
	return p.writer.Close()
}

type CompanyEvent struct {
	Type    string  `json:"type"`
	Company Company `json:"company"`
}

type Company struct {
	UUID              uuid.UUID `json:"uuid"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	AmountOfEmployees int       `json:"amount_of_employees"`
	Registered        bool      `json:"registered"`
	Type              string    `json:"type"`
}
