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

func TestMessageFuncService_SendMessage(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserIDKey, int64(1))

	tests := []struct {
		name      string
		req       dto.SendMessageRequest
	 storageFn func() *mocks.MockStorageMessage
		eventFn   func() *mocks.MockEventBus
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.SendMessageRequest{ChatID: 10, Content: "hello", MessageType: dto.MessageTypeText},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					CreateMessageFn: func(_ context.Context, _ int64, _ dto.SendMessageRequest) (models.Message, error) {
						return models.Message{ID: 1, ChatID: 10, SenderID: 1, Content: "hello", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
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
			req:  dto.SendMessageRequest{ChatID: 10, Content: "hello"},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					CreateMessageFn: func(_ context.Context, _ int64, _ dto.SendMessageRequest) (models.Message, error) {
						return models.Message{}, errors.New("db error")
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{}
			},
			wantErr: true,
		},
		{
			name: "publish error",
			req:  dto.SendMessageRequest{ChatID: 10, Content: "hello"},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					CreateMessageFn: func(_ context.Context, _ int64, _ dto.SendMessageRequest) (models.Message, error) {
						return models.Message{ID: 1, ChatID: 10}, nil
					},
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{
					PublishFn: func(_ context.Context, _ event.Event) error { return errors.New("redis error") },
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewMessageFuncService(zap.NewNop(), tt.storageFn(), tt.eventFn())
			msg, err := svc.SendMessage(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && msg.ID == 0 {
				t.Error("SendMessage() returned message with zero ID")
			}
		})
	}
}

func TestMessageFuncService_GetMessages(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.GetMessagesRequest
		storageFn func() *mocks.MockStorageMessage
		wantCount int32
		wantMore  bool
		wantErr   bool
	}{
		{
			name: "multiple messages",
			req:  dto.GetMessagesRequest{ChatID: 10, Limit: 50},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					SelectMessagesFn: func(_ context.Context, _ dto.GetMessagesRequest) ([]*models.Message, error) {
						return []*models.Message{{ID: 1}, {ID: 2}}, nil
					},
				}
			},
			wantCount: 2,
			wantMore:  true,
		},
		{
			name: "empty messages",
			req:  dto.GetMessagesRequest{ChatID: 10},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					SelectMessagesFn: func(_ context.Context, _ dto.GetMessagesRequest) ([]*models.Message, error) {
						return nil, nil
					},
				}
			},
			wantCount: 0,
			wantMore:  false,
		},
		{
			name: "storage error",
			req:  dto.GetMessagesRequest{ChatID: 10},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					SelectMessagesFn: func(_ context.Context, _ dto.GetMessagesRequest) ([]*models.Message, error) {
						return nil, errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewMessageFuncService(zap.NewNop(), tt.storageFn(), &mocks.MockEventBus{})
			resp, err := svc.GetMessages(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if resp.HasMore != tt.wantMore {
					t.Errorf("GetMessages() HasMore = %v, want %v", resp.HasMore, tt.wantMore)
				}
				if resp.TotalCount != tt.wantCount {
					t.Errorf("GetMessages() TotalCount = %v, want %v", resp.TotalCount, tt.wantCount)
				}
			}
		})
	}
}

func TestMessageFuncService_EditMessage(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserIDKey, int64(1))

	tests := []struct {
		name      string
		req       dto.EditMessageRequest
		storageFn func() *mocks.MockStorageMessage
		eventFn   func() *mocks.MockEventBus
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.EditMessageRequest{MessageID: 1, Content: "edited"},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					UpdateMessageFn: func(_ context.Context, _ dto.EditMessageRequest) (models.Message, error) {
						return models.Message{ID: 1, ChatID: 10, Content: "edited", UpdatedAt: time.Now()}, nil
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
			req:  dto.EditMessageRequest{MessageID: 1, Content: "edited"},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					UpdateMessageFn: func(_ context.Context, _ dto.EditMessageRequest) (models.Message, error) {
						return models.Message{}, errors.New("not found")
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
			svc := service.NewMessageFuncService(zap.NewNop(), tt.storageFn(), tt.eventFn())
			msg, err := svc.EditMessage(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && msg.Content != "edited" {
				t.Errorf("EditMessage() content = %v, want 'edited'", msg.Content)
			}
		})
	}
}

func TestMessageFuncService_DeleteMessage(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserIDKey, int64(1))

	tests := []struct {
		name      string
		req       dto.DeleteMessageRequest
		storageFn func() *mocks.MockStorageMessage
		eventFn   func() *mocks.MockEventBus
		want      bool
		wantErr   bool
	}{
		{
			name: "delete for everyone",
			req:  dto.DeleteMessageRequest{MessageID: 1, ForEveryone: true, ChatID: 10},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					DeleteMessageFn: func(_ context.Context, _ int64) error { return nil },
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{
					PublishFn: func(_ context.Context, _ event.Event) error { return nil },
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "delete for self",
			req:  dto.DeleteMessageRequest{MessageID: 1, ForEveryone: false, ChatID: 10},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					DeleteMessageForUserFn: func(_ context.Context, _ int64) error { return nil },
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{
					PublishFn: func(_ context.Context, _ event.Event) error { return nil },
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "storage error",
			req:  dto.DeleteMessageRequest{MessageID: 1, ForEveryone: true, ChatID: 10},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					DeleteMessageFn: func(_ context.Context, _ int64) error { return errors.New("db error") },
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{}
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewMessageFuncService(zap.NewNop(), tt.storageFn(), tt.eventFn())
			got, err := svc.DeleteMessage(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("DeleteMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageFuncService_GetMessage(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.GetMessageRequest
		storageFn func() *mocks.MockStorageMessage
		wantID    int64
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.GetMessageRequest{MessageID: 1},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					SelectMessageFn: func(_ context.Context, id int64) (models.Message, error) {
						return models.Message{ID: id, Content: "hello"}, nil
					},
				}
			},
			wantID:  1,
			wantErr: false,
		},
		{
			name: "not found",
			req:  dto.GetMessageRequest{MessageID: 999},
			storageFn: func() *mocks.MockStorageMessage {
				return &mocks.MockStorageMessage{
					SelectMessageFn: func(_ context.Context, _ int64) (models.Message, error) {
						return models.Message{}, errors.New("not found")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewMessageFuncService(zap.NewNop(), tt.storageFn(), &mocks.MockEventBus{})
			msg, err := svc.GetMessage(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && msg.ID != tt.wantID {
				t.Errorf("GetMessage() ID = %v, want %v", msg.ID, tt.wantID)
			}
		})
	}
}
