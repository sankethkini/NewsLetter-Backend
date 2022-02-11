//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/sankethkini/NewsLetter-Backend/internal/config"
	"github.com/sankethkini/NewsLetter-Backend/internal/service/user"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/database"
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

func IntializeUserRepo() (*user.UserServiceImpl, func(), error) {
	panic(wire.Build(
		JWTProviderSet,
		config.LoadDataBaseConfig,
		database.Open,
		user.NewDB,
		user.NewUserService,
	))
}
