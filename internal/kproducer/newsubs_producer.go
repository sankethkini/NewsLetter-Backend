package kproducer

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/pkg/apperrors"
	"github.com/segmentio/kafka-go"
)

const (
	errWrite = "kafka: error in writing message"
)

//go:generate mockgen -destination kafka_mock.go -package kproducer github.com/sankethkini/NewsLetter-Backend/internal/kproducer Producer
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
	return apperrors.E(ctx, errors.Wrapf(err, errWrite))
}
