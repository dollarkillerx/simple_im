package conf

import "simple_im/pkg/common/config"

type Config struct {
	ServiceConfiguration  config.ServiceConfiguration
	PostgresConfiguration config.PostgresConfiguration
	RedisConfiguration    config.RedisConfiguration
	LoggerConfiguration   config.LoggerConfig
	JWTConfiguration      config.JWTConfiguration
	UploadConfiguration   config.UploadConfiguration
}
