package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Host               string `yaml:"host"`
	Name               string `yaml:"name"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	DisableTLS         bool   `yaml:"disable_tls"`
	Debug              bool   `yaml:"debug"`
}

type JWTConfig struct {
	Secret   string `yaml:"secret"`
	Duration int    `yaml:"duration"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	Database Database     `yaml:"database"`
	Jwt      JWTConfig    `yaml:"jwt"`
	Server   ServerConfig `yaml:"server"`
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

func LoadDataBaseConfig(app *AppConfig) Database {
	return app.Database
}

func LoadJWTConfig(app *AppConfig) JWTConfig {
	return app.Jwt
}

func LaodServerConfig(app *AppConfig) ServerConfig {
	return app.Server
}

func LoadAccessibleRoles() map[string][]string {
	// const userServicePath = "/userpb.v1.".

	return map[string][]string{
		// this will be done when needed and it will be like below.
		// laptopServicePath + "CreateUser": {"user"},.
	}
}
