//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/sankethkini/NewsLetter-Backend/internal/config"
	"github.com/sankethkini/NewsLetter-Backend/internal/kproducer"
	"github.com/sankethkini/NewsLetter-Backend/internal/service"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/admin"
	newsletter "github.com/sankethkini/NewsLetter-Backend/internal/service/news_letter"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/subscription"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/user"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/cache"
	"github.com/sankethkini/NewsLetter-Backend/pkg/database"
	"github.com/sankethkini/NewsLetter-Backend/pkg/email"
	kafkaservice "github.com/sankethkini/NewsLetter-Backend/pkg/kafka"
)

var JWTProviderSet = wire.NewSet(config.LoadConfig, config.LoadJWTConfig, auth.NewJWTManager)

func IntializeServerConfig() (config.ServerConfig, error) {
	panic(wire.Build(config.LoadConfig, config.LaodServerConfig))
}

func IntializeJWT() (*auth.AuthInterceptor, error) {
	panic(wire.Build(
		JWTProviderSet,
		config.LoadAccessibleRoles,
		auth.NewAuthInterceptor,
	))
}

func IntializeConsumer() (*kafkaservice.Consumer, error) {
	panic(wire.Build(
		config.LoadConfig,
		config.LoadKafkaConsumer,
		config.LoadEmailConfig,
		email.NewEmailServer,
		kafkaservice.NewConsumer,
	))
}

func IntializeServiceRegistry() (*service.Registry, func(), error) {
	panic(wire.Build(
		config.LoadConfig,
		config.LoadDataBaseConfig,
		config.LoadJWTConfig,
		config.LoadRedisConfig,
		config.LoadKafkaConfig,
		auth.NewJWTManager,
		database.Open,
		cache.NewRedisCache,
		kafkaservice.NewProducer,
		kproducer.NewProducer,
		user.NewDB,
		user.NewUserService,
		admin.NewRepo,
		admin.NewAdminService,
		subscription.NewSubRepo,
		subscription.NewSubService,
		newsletter.NewNewsRepo,
		newsletter.NewNewsService,
		service.NewRegistry,
	))
}
