package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"github.com/north-fy/talker/services/message/internal/mocks"
	"github.com/north-fy/talker/services/message/internal/service"
	"go.uber.org/zap"
)

func TestFeatureService_SearchMessages(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.SearchMessagesRequest
		storageFn func() *mocks.MockStorageFeature
		wantCount int
		wantMore  bool
		wantErr   bool
	}{
		{
			name: "found messages",
			req:  dto.SearchMessagesRequest{ChatID: 10, Query: "hello"},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SearchMessagesFn: func(_ context.Context, _ dto.SearchMessagesRequest) ([]*models.Message, error) {
						return []*models.Message{{ID: 1}, {ID: 2}}, nil
					},
				}
			},
			wantCount: 2,
			wantMore:  true,
		},
		{
			name: "no messages",
			req:  dto.SearchMessagesRequest{ChatID: 10, Query: "nonexistent"},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SearchMessagesFn: func(_ context.Context, _ dto.SearchMessagesRequest) ([]*models.Message, error) {
						return nil, nil
					},
				}
			},
			wantCount: 0,
			wantMore:  false,
		},
		{
			name: "storage error",
			req:  dto.SearchMessagesRequest{ChatID: 10, Query: "hello"},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SearchMessagesFn: func(_ context.Context, _ dto.SearchMessagesRequest) ([]*models.Message, error) {
						return nil, errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewFeatureService(zap.NewNop(), tt.storageFn(), &mocks.MockEventBus{})
			resp, err := svc.SearchMessages(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if resp.HasMore != tt.wantMore {
					t.Errorf("SearchMessages() HasMore = %v, want %v", resp.HasMore, tt.wantMore)
				}
				if len(resp.Messages) != tt.wantCount {
					t.Errorf("SearchMessages() count = %v, want %v", len(resp.Messages), tt.wantCount)
				}
			}
		})
	}
}

func TestFeatureService_MarkAsRead(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.MarkAsReadRequest
		storageFn func() *mocks.MockStorageFeature
		eventFn   func() *mocks.MockEventBus
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.MarkAsReadRequest{ChatID: 10, UserID: 1, UpToMessageID: 5},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SetAsReadFn: func(_ context.Context, _ dto.MarkAsReadRequest) error { return nil },
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
			req:  dto.MarkAsReadRequest{ChatID: 10, UserID: 1, UpToMessageID: 5},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SetAsReadFn: func(_ context.Context, _ dto.MarkAsReadRequest) error { return errors.New("db error") },
				}
			},
			eventFn: func() *mocks.MockEventBus {
				return &mocks.MockEventBus{}
			},
			wantErr: true,
		},
		{
			name: "publish error",
			req:  dto.MarkAsReadRequest{ChatID: 10, UserID: 1, UpToMessageID: 5},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SetAsReadFn: func(_ context.Context, _ dto.MarkAsReadRequest) error { return nil },
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
			svc := service.NewFeatureService(zap.NewNop(), tt.storageFn(), tt.eventFn())
			err := svc.MarkAsRead(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkAsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeatureService_GetUnreadCount(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.GetUnreadCountRequest
		storageFn func() *mocks.MockStorageFeature
		wantCount int32
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.GetUnreadCountRequest{ChatID: 10, UserID: 1},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SelectUnreadCountFn: func(_ context.Context, _ dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error) {
						return dto.GetUnreadCountResponse{Count: 5, LastMessageID: 20}, nil
					},
				}
			},
			wantCount: 5,
		},
		{
			name: "storage error",
			req:  dto.GetUnreadCountRequest{ChatID: 10, UserID: 1},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SelectUnreadCountFn: func(_ context.Context, _ dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error) {
						return dto.GetUnreadCountResponse{}, errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewFeatureService(zap.NewNop(), tt.storageFn(), &mocks.MockEventBus{})
			resp, err := svc.GetUnreadCount(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUnreadCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && resp.Count != tt.wantCount {
				t.Errorf("GetUnreadCount() Count = %v, want %v", resp.Count, tt.wantCount)
			}
		})
	}
}

func TestFeatureService_GetLastMessage(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.GetLastMessageRequest
		storageFn func() *mocks.MockStorageFeature
		wantID    int64
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.GetLastMessageRequest{ChatID: 10},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SelectLastMessageFn: func(_ context.Context, _ dto.GetLastMessageRequest) (models.Message, error) {
						return models.Message{ID: 42, Content: "last msg"}, nil
					},
				}
			},
			wantID:  42,
			wantErr: false,
		},
		{
			name: "storage error",
			req:  dto.GetLastMessageRequest{ChatID: 10},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					SelectLastMessageFn: func(_ context.Context, _ dto.GetLastMessageRequest) (models.Message, error) {
						return models.Message{}, errors.New("db error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewFeatureService(zap.NewNop(), tt.storageFn(), &mocks.MockEventBus{})
			msg, err := svc.GetLastMessage(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && msg.ID != tt.wantID {
				t.Errorf("GetLastMessage() ID = %v, want %v", msg.ID, tt.wantID)
			}
		})
	}
}

func TestFeatureService_DeleteChatMessages(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.DeleteChatMessagesRequest
		storageFn func() *mocks.MockStorageFeature
		wantErr   bool
	}{
		{
			name: "success",
			req:  dto.DeleteChatMessagesRequest{ChatID: 10},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					DeleteChatMessagesFn: func(_ context.Context, _ dto.DeleteChatMessagesRequest) error { return nil },
				}
			},
			wantErr: false,
		},
		{
			name: "storage error",
			req:  dto.DeleteChatMessagesRequest{ChatID: 10},
			storageFn: func() *mocks.MockStorageFeature {
				return &mocks.MockStorageFeature{
					DeleteChatMessagesFn: func(_ context.Context, _ dto.DeleteChatMessagesRequest) error { return errors.New("db error") },
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewFeatureService(zap.NewNop(), tt.storageFn(), &mocks.MockEventBus{})
			err := svc.DeleteChatMessages(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteChatMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
