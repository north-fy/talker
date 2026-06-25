package mocks

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
)

type MockStorageMessage struct {
	CreateMessageFn       func(ctx context.Context, senderID int64, req dto.SendMessageRequest) (models.Message, error)
	SelectMessagesFn      func(ctx context.Context, req dto.GetMessagesRequest) ([]*models.Message, error)
	UpdateMessageFn       func(ctx context.Context, req dto.EditMessageRequest) (models.Message, error)
	DeleteMessageForUserFn func(ctx context.Context, id int64) error
	DeleteMessageFn       func(ctx context.Context, id int64) error
	SelectMessageFn       func(ctx context.Context, id int64) (models.Message, error)
}

func (m *MockStorageMessage) CreateMessage(ctx context.Context, senderID int64, req dto.SendMessageRequest) (models.Message, error) {
	if m.CreateMessageFn != nil {
		return m.CreateMessageFn(ctx, senderID, req)
	}
	return models.Message{}, nil
}

func (m *MockStorageMessage) SelectMessages(ctx context.Context, req dto.GetMessagesRequest) ([]*models.Message, error) {
	if m.SelectMessagesFn != nil {
		return m.SelectMessagesFn(ctx, req)
	}
	return nil, nil
}

func (m *MockStorageMessage) UpdateMessage(ctx context.Context, req dto.EditMessageRequest) (models.Message, error) {
	if m.UpdateMessageFn != nil {
		return m.UpdateMessageFn(ctx, req)
	}
	return models.Message{}, nil
}

func (m *MockStorageMessage) DeleteMessageForUser(ctx context.Context, id int64) error {
	if m.DeleteMessageForUserFn != nil {
		return m.DeleteMessageForUserFn(ctx, id)
	}
	return nil
}

func (m *MockStorageMessage) DeleteMessage(ctx context.Context, id int64) error {
	if m.DeleteMessageFn != nil {
		return m.DeleteMessageFn(ctx, id)
	}
	return nil
}

func (m *MockStorageMessage) SelectMessage(ctx context.Context, id int64) (models.Message, error) {
	if m.SelectMessageFn != nil {
		return m.SelectMessageFn(ctx, id)
	}
	return models.Message{}, nil
}
