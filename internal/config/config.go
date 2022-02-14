package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sankethkini/NewsLetter-Backend/pkg/auth"
	"github.com/sankethkini/NewsLetter-Backend/pkg/cache"
	"github.com/sankethkini/NewsLetter-Backend/pkg/database"
	"github.com/sankethkini/NewsLetter-Backend/pkg/email"
	kafkaservice "github.com/sankethkini/NewsLetter-Backend/pkg/kafka"
	"github.com/sankethkini/NewsLetter-Backend/pkg/role"
)

type Role int

const (
	ADMIN Role = iota
	USER
)

func (r Role) String() string {
	return []string{"admin", "user"}[r]
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	Database      database.Database           `yaml:"database"`
	Jwt           auth.JWTConfig              `yaml:"jwt"`
	Server        ServerConfig                `yaml:"server"`
	Redis         cache.RedisConfig           `yaml:"redis"`
	KafkaProducer kafkaservice.KafkaConfig    `yaml:"kafkap"`
	KafkaConsumer kafkaservice.ConsumerConfig `yaml:"kafkac"`
	Email         email.EmailConfig           `yaml:"email"`
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig
	err := cleanenv.ReadConfig("application.yaml", &config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &config, nil
}

func LoadDataBaseConfig(app *AppConfig) database.Database {
	return app.Database
}

func LoadJWTConfig(app *AppConfig) auth.JWTConfig {
	return app.Jwt
}

func LaodServerConfig(app *AppConfig) ServerConfig {
	return app.Server
}

func LoadRedisConfig(app *AppConfig) cache.RedisConfig {
	return app.Redis
}

func LoadKafkaConfig(app *AppConfig) kafkaservice.KafkaConfig {
	return app.KafkaProducer
}

func LoadKafkaConsumer(app *AppConfig) kafkaservice.ConsumerConfig {
	return app.KafkaConsumer
}

func LoadEmailConfig(app *AppConfig) email.EmailConfig {
	return app.Email
}

func LoadAccessibleRoles() map[string][]string {
	const subsPath = "/subscriptionpb.v1.SubscriptionService/"
	const newsPath = "/newsletterpb.v1.NewsLetterService/"

	return map[string][]string{
		// this will be done when needed and it will be like below.

		newsPath + "CreateNewsLetter": {role.ADMIN.String()},
		newsPath + "AddSchemeToNews":  {role.ADMIN.String()},
		subsPath + "CreateScheme":     {role.ADMIN.String()},

		subsPath + "AddUser":    {role.USER.String()},
		subsPath + "RemoveUser": {role.USER.String()},
		subsPath + "Renew":      {role.USER.String()},
	}
}
