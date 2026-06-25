package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/mocks"
	"github.com/north-fy/talker/services/message/internal/service"
	"go.uber.org/zap"
)

func TestWebSocketService_HandleClientMessage(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.EventRequest
		eventFn   func() *mocks.MockEventBus
		wantErr   bool
	}{
		{
			name: "subscribe error",
			req:  dto.EventRequest{ChatID: 10, UserID: 1},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{
					SubscribeFn: func(_ context.Context, _ int64) (<-chan event.Event, error) {
						return nil, errors.New("subscribe error")
					},
				}
			},
			wantErr: true,
		},
		{
			name: "context cancellation",
			req:  dto.EventRequest{ChatID: 10, UserID: 1},
			eventFn: func() *mocks.MockEventBus {
				ch := make(chan event.Event, 1)
				return &mocks.MockEventBus{
					SubscribeFn: func(_ context.Context, _ int64) (<-chan event.Event, error) {
						return ch, nil
					},
				}
			},
			wantErr: true,
		},
		{
			name: "invalid message data",
			req:  dto.EventRequest{ChatID: 10, UserID: 1},
			eventFn: func() *mocks.MockEventBus {
				ch := make(chan event.Event, 1)
				ch <- &event.MessageEvent{
					Type:    event.WebSocketMessage_NewMessage,
					ChatID:  10,
					Payload: json.RawMessage(`invalid json`),
				}
				return &mocks.MockEventBus{
					SubscribeFn: func(_ context.Context, _ int64) (<-chan event.Event, error) {
						return ch, nil
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWebSocketService(zap.NewNop(), tt.eventFn())
			sendChan := make(chan *messagev1.WebSocketMessage, 1)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.name == "context cancellation" {
				cancel()
			}

			err := svc.HandleClientMessage(ctx, tt.req, sendChan)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleClientMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
