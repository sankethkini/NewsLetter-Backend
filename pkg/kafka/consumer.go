package kafkaservice

import (
	"context"

	newsletter "github.com/sankethkini/NewsLetter-Backend/internal/service/news_letter"
	"github.com/sankethkini/NewsLetter-Backend/pkg/email"
	kafka "github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"groupID"`
	Brokers []string `yaml:"brokers"`
}

type Consumer struct {
	reader *kafka.Reader
	email  *email.Email
}

func NewConsumer(cfg ConsumerConfig, email *email.Email) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})
	return &Consumer{reader: r, email: email}
}

func (c Consumer) Consume(ctx context.Context) {
	// logger := ctxzap.Extract(ctx)

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			panic(err)
		}

		err = c.PostConsume(ctx, msg)
		if err != nil {
			panic(err)
		}
	}
}

// nolint: govet
func (c Consumer) PostConsume(ctx context.Context, msg kafka.Message) error {
	msg1, err := newsletter.ToModel(string(msg.Value))
	if err != nil {
		return err
	}
	emails, err := c.email.GetEmails(ctx, msg1)
	if err != nil {
		return err
	}
	m := email.NewMail(emails, msg1.Letter.Title, msg1.Letter.Body)
	err = c.email.SendEmail(m)
	if err != nil {
		return err
	}
	return nil
}
