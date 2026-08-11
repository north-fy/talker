package redis

import (
	"context"
	"fmt"

	"github.com/north-fy/talker/services/chat/internal/config"
	"github.com/redis/go-redis/v9"
)

const CountTriesPing = 3

type Storage struct {
	client *redis.Client
}

func NewStorage(ctx context.Context, cfg config.RedisCfg) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		DB:           cfg.DB,
	})

	for range CountTriesPing {
		if err := client.Ping(ctx).Err(); err != nil {
			panic(err)
		}
	}

	return &Storage{client: client}
}

func (s *Storage) Close() error {
	return s.client.Close()
}
