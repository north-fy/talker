package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/north-fy/talker/services/message/internal/config"
	"github.com/redis/go-redis/v9"
)

const CountTriesPing = 3

type Storage struct {
	client *redis.Client
	mu     *sync.RWMutex
	subs   map[int64]*redis.PubSub // key chatID -> value *redis.PubSub
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

	return &Storage{
		client: client,
		mu:     &sync.RWMutex{},
		subs:   make(map[int64]*redis.PubSub),
	}
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sub := range s.subs {
		sub.Close()
		delete(s.subs, id)
	}

	return s.client.Close()
}
