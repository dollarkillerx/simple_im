package client

import (
	"fmt"

	"simple_im/pkg/common/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func PostgresClient(conf config.PostgresConfiguration, gormConfig *gorm.Config) (*gorm.DB, error) {
	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}
	if conf.Port == 0 {
		conf.Port = 5432
	}
	if conf.TimeZone == "" {
		conf.TimeZone = "Asia/Shanghai"
	}

	sslMode := "disable"
	if conf.SSLMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, sslMode, conf.TimeZone,
	)

	if gormConfig == nil {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}
