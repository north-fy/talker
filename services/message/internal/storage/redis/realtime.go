package redis

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
)

func (s *Storage) PubMessage(ctx context.Context, channel string, data interface{}) error {
	return s.client.Publish(ctx, channel, data).Err()
}

func (s *Storage) SubMessage(ctx context.Context, channels ...string) dto.Subscribe {
	sub := s.client.Subscribe(ctx, channels...)

	_, err := sub.Receive(ctx)
	if err != nil {
		return dto.Subscribe{}
	}

	return dto.Subscribe{
		Data:  sub.Channel(),
		Close: sub.Close,
	}
}
