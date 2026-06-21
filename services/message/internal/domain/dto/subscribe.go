package dto

import "github.com/redis/go-redis/v9"

type Subscribe struct {
	Data <-chan *redis.Message
	Close   func() error
}

func (s *Subscribe) GetData() <-chan *redis.Message {
	return s.Data
}
