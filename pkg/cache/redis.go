package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
	ExpireIn int64  `yaml:"expire"`
}

//go:generate mockgen -destination redis_mock.go -package cache github.com/sankethkini/NewsLetter-Backend/pkg/cache Service
type Service interface {
	Set(ctx context.Context, key string, value []*subscriptionpb.Scheme)
	Get(ctx context.Context, key string) ([]*subscriptionpb.Scheme, error)
}

type service struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedisCache(cfg RedisConfig) Service {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &service{
		client: client,
		exp:    time.Duration(cfg.ExpireIn) * time.Second,
	}
}

func (s service) Set(ctx context.Context, key string, value []*subscriptionpb.Scheme) {
	logger := ctxzap.Extract(ctx)
	js, err := json.Marshal(value)
	if err != nil {
		logger.Sugar().Fatal("caanot marshal into cache")
		return
	}
	s.client.Set(ctx, key, string(js), s.exp)
}

func (s service) Get(ctx context.Context, key string) ([]*subscriptionpb.Scheme, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var res []*subscriptionpb.Scheme
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
