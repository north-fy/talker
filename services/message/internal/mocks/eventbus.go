package mocks

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/event"
)

type MockEventBus struct {
	PublishFn    func(ctx context.Context, ev event.Event) error
	SubscribeFn  func(ctx context.Context, chatID int64) (<-chan event.Event, error)
	CloseFn      func() error
}

func (m *MockEventBus) Publish(ctx context.Context, ev event.Event) error {
	if m.PublishFn != nil {
		return m.PublishFn(ctx, ev)
	}
	return nil
}

func (m *MockEventBus) Subscribe(ctx context.Context, chatID int64) (<-chan event.Event, error) {
	if m.SubscribeFn != nil {
		return m.SubscribeFn(ctx, chatID)
	}
	ch := make(chan event.Event, 1)
	return ch, nil
}

func (m *MockEventBus) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}
