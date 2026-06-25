package event

import (
	"context"
)

type Event interface {
	GetType() int32
	GetChatID() int64
	GetData() []byte
}

type EventBus interface {
	Publish(ctx context.Context, ev Event) error
	Subscribe(ctx context.Context, chatID int64) (<-chan Event, error)
	Close() error
}

