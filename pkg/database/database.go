package database

import (
	"fmt"

	"github.com/sankethkini/NewsLetter-Backend/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// should return close function.
func Open(cfg config.Database) (*gorm.DB, func(), error) {
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
