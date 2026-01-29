package config

import (
	"github.com/spf13/viper"
)

type ServiceConfiguration struct {
	Port  string
	Debug bool
}

type PostgresConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  bool
	TimeZone string
}

type RedisConfiguration struct {
	Addr     string
	Db       int
	Password string
}

type LoggerConfig struct {
	Filename string
	MaxSize  int
}

type JWTConfiguration struct {
	Secret string
	Expire int64 // seconds
}

type UploadConfiguration struct {
	MaxSize    int64    // bytes
	SavePath   string
	AllowTypes []string
}

func InitConfiguration(configName string, configPaths []string, config interface{}) error {
	vp := viper.New()
	vp.SetConfigName(configName)
	vp.SetConfigType("toml")
	vp.AutomaticEnv()

	for _, p := range configPaths {
		vp.AddConfigPath(p)
	}

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	if err := vp.Unmarshal(config); err != nil {
		return err
	}

	return nil
}
