package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	User                  string        `yaml:"user"`
	Password              string        `yaml:"password"`
	Host                  string        `yaml:"host"`
	Name                  string        `yaml:"name"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
	DisableTLS            bool          `yaml:"disable_tls"`
	Debug                 bool          `yaml:"debug"`
}

//should return close function
func Open() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "1234", "localhost:3306", "newdb")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, err
	// }

	// sqlDB.SetConnMaxIdleTime(cfg.MaxConnectionIdleTime)
	// sqlDB.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)
	// sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	// sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)

	return db, nil

}
