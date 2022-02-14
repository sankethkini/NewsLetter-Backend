package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

// should return close function.
func Open(cfg Database) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			panic(err)
		}
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)

	return db, cleanup, nil
}
