package mocks

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
)

type MockStorageFeature struct {
	SearchMessagesFn     func(ctx context.Context, req dto.SearchMessagesRequest) ([]*models.Message, error)
	SetAsReadFn          func(ctx context.Context, req dto.MarkAsReadRequest) error
	SelectUnreadCountFn  func(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error)
	SelectLastMessageFn  func(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error)
	DeleteChatMessagesFn func(ctx context.Context, req dto.DeleteChatMessagesRequest) error
}

func (m *MockStorageFeature) SearchMessages(ctx context.Context, req dto.SearchMessagesRequest) ([]*models.Message, error) {
	if m.SearchMessagesFn != nil {
		return m.SearchMessagesFn(ctx, req)
	}
	return nil, nil
}

func (m *MockStorageFeature) SetAsRead(ctx context.Context, req dto.MarkAsReadRequest) error {
	if m.SetAsReadFn != nil {
		return m.SetAsReadFn(ctx, req)
	}
	return nil
}

func (m *MockStorageFeature) SelectUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error) {
	if m.SelectUnreadCountFn != nil {
		return m.SelectUnreadCountFn(ctx, req)
	}
	return dto.GetUnreadCountResponse{}, nil
}

func (m *MockStorageFeature) SelectLastMessage(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error) {
	if m.SelectLastMessageFn != nil {
		return m.SelectLastMessageFn(ctx, req)
	}
	return models.Message{}, nil
}

func (m *MockStorageFeature) DeleteChatMessages(ctx context.Context, req dto.DeleteChatMessagesRequest) error {
	if m.DeleteChatMessagesFn != nil {
		return m.DeleteChatMessagesFn(ctx, req)
	}
	return nil
}
