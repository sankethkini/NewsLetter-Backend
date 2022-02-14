package kproducer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(context.Context, []byte, []byte) error
}

type producer struct {
	kafka *kafka.Writer
}

func NewProducer(kf *kafka.Writer) Producer {
	return &producer{
		kafka: kf,
	}
}

func (p producer) Produce(ctx context.Context, key []byte, value []byte) error {
	err := p.kafka.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	return err
}
