package service

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"go.uber.org/zap"
)

type ChatFuncService struct {
	log     *zap.Logger
	storage ChatStorage
}

type ChatStorage interface {
	InsertChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error)
	SelectChat(ctx context.Context, chatID int64) (models.Chat, error)
	SelectChats(ctx context.Context, filter dto.ChatFilter) (dto.GetChatsResponse, error)
	UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error)
}

func NewChatFuncService(log *zap.Logger, storage ChatStorage) *ChatFuncService {
	return &ChatFuncService{
		log:     log,
		storage: storage,
	}
}

// TODO: добавить валидацию на все функции
// TODO: добавить обработку ошибок

func (s *ChatFuncService) CreateChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	chat, err := s.storage.InsertChat(ctx, req)
	if err != nil {
		log.Error("failed to create chat", zap.Error(err))
		return models.Chat{}, err
	}

	return chat, err
}

func (s *ChatFuncService) GetChat(ctx context.Context, req dto.GetChatRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	chat, err := s.storage.SelectChat(ctx, req.ChatID)
	if err != nil {
		log.Error("failed to select chat", zap.Error(err))
		return models.Chat{}, err
	}

	if !req.IncludeMembers {
		chat.MembersCount = 0
	}

	return chat, nil
}

func (s *ChatFuncService) GetChats(ctx context.Context, req dto.GetChatsRequest) (dto.GetChatsResponse, error) {
	log := s.log.With(zap.Any("request", req))

	chats, err := s.storage.SelectChats(ctx, req.Filter)
	if err != nil {
		log.Error("failed to select chats", zap.Error(err))
		return dto.GetChatsResponse{}, err
	}

	return chats, nil
}

func (s *ChatFuncService) UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	chat, err := s.storage.UpdateChat(ctx, req)
	if err != nil {
		log.Error("failed to update chat", zap.Error(err))
		return models.Chat{}, err
	}

	return chat, nil
}
