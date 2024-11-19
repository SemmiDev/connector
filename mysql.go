package g_learning_connector

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDatabase(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBDatabase,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	conn, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	conn.SetMaxOpenConns(config.DBPoolMax)
	conn.SetMaxIdleConns(config.DBPoolIdle)
	conn.SetConnMaxLifetime(config.DBPoolLifetime)

	return db, nil
}
