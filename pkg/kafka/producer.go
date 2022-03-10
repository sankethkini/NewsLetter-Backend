package kafkaservice

import (
	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

func NewProducer(cfg KafkaConfig) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	})
	return w
}
