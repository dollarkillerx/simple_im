package main

import (
	"flag"
	"log"
	"strings"

	"simple_im/internal/conf"
	"simple_im/internal/server"
	"simple_im/pkg/common/client"
	"simple_im/pkg/common/config"
	logs "simple_im/pkg/common/log"
)

func main() {
	configName := flag.String("c", "config", "config file name (without extension)")
	configPaths := flag.String("cPath", "./,./configs/", "config file search paths (comma separated)")
	flag.Parse()

	paths := strings.Split(*configPaths, ",")

	var appConfig conf.Config
	if err := config.InitConfiguration(*configName, paths, &appConfig); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logs.InitLog(appConfig.LoggerConfiguration)

	db, err := client.PostgresClient(appConfig.PostgresConfiguration, nil)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	redisClient, err := client.RedisClient(appConfig.RedisConfiguration)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	s := server.NewServer(db, redisClient, appConfig)
	if err := s.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
