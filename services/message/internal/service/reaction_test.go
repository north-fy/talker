package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"github.com/north-fy/talker/services/message/internal/mocks"
	"github.com/north-fy/talker/services/message/internal/service"
	"go.uber.org/zap"
)

func TestReactionService_AddReaction(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserIDKey, int64(1))

	tests := []struct {
		name      string
		req       dto.AddReactionRequest
		storageFn func() *mocks.MockStorageReaction
		eventFn   func() *mocks.MockEventBus
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.AddReactionRequest{MessageID: 10, Reaction: "👍"},
			storageFn: func() *mocks.MockStorageReaction {
				return &mocks.MockStorageReaction{
					InsertReactionFn: func(_ context.Context, _ dto.AddReactionRequest) (dto.Reaction, error) {
						return dto.Reaction{
							MessageID: 10,
							UserID:    1,
							Reaction:  `{"👍":3}`,
							CreatedAt: time.Now(),
						}, nil
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{
					PublishFn: func(_ context.Context, _ event.Event) error { return nil },
				}
			},
			wantErr: false,
		},
		{
			name: "storage error",
			req:  dto.AddReactionRequest{MessageID: 10, Reaction: "👍"},
			storageFn: func() *mocks.MockStorageReaction {
				return &mocks.MockStorageReaction{
					InsertReactionFn: func(_ context.Context, _ dto.AddReactionRequest) (dto.Reaction, error) {
						return dto.Reaction{}, errors.New("db error")
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{}
			},
			wantErr: true,
		},
		{
			name: "invalid reaction json",
			req:  dto.AddReactionRequest{MessageID: 10, Reaction: "👍"},
			storageFn: func() *mocks.MockStorageReaction {
				return &mocks.MockStorageReaction{
					InsertReactionFn: func(_ context.Context, _ dto.AddReactionRequest) (dto.Reaction, error) {
						return dto.Reaction{Reaction: "invalid json"}, nil
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewReactionService(zap.NewNop(), tt.storageFn(), tt.eventFn())
			_, err := svc.AddReaction(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddReaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReactionService_RemoveReaction(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserIDKey, int64(1))

	tests := []struct {
		name      string
		req       dto.RemoveReactionRequest
		storageFn func() *mocks.MockStorageReaction
		eventFn   func() *mocks.MockEventBus
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.RemoveReactionRequest{MessageID: 10, Reaction: "👍"},
			storageFn: func() *mocks.MockStorageReaction {
				return &mocks.MockStorageReaction{
					DeleteReactionFn: func(_ context.Context, _ dto.RemoveReactionRequest) (string, error) {
						return `{"👍":2}`, nil
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{
					PublishFn: func(_ context.Context, _ event.Event) error { return nil },
				}
			},
			wantErr: false,
		},
		{
			name: "storage error",
			req:  dto.RemoveReactionRequest{MessageID: 10, Reaction: "👍"},
			storageFn: func() *mocks.MockStorageReaction {
				return &mocks.MockStorageReaction{
					DeleteReactionFn: func(_ context.Context, _ dto.RemoveReactionRequest) (string, error) {
						return "", errors.New("db error")
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewReactionService(zap.NewNop(), tt.storageFn(), tt.eventFn())
			err := svc.RemoveReaction(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveReaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
