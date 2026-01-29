package client

import (
	"context"
	"fmt"

	"simple_im/pkg/common/config"

	"github.com/redis/go-redis/v9"
)

func RedisClient(conf config.RedisConfiguration) (*redis.Client, error) {
	if conf.Addr == "" {
		conf.Addr = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
