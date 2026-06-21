package event

import "context"

type Event interface {
	Type() string
	AggregateID() string
	Payload() []byte
}

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, aggregateID string) (<-chan Event, error)
	Close() error
}