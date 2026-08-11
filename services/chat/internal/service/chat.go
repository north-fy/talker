package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"go.uber.org/zap"
)

type ChatFuncService struct {
	log     *zap.Logger
	storage ChatStorage
	cache   Cache
}

type ChatStorage interface {
	InsertChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error)
	SelectChat(ctx context.Context, chatID int64) (models.Chat, error)
	SelectChats(ctx context.Context, filter dto.ChatFilter) (dto.GetChatsResponse, error)
	UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error)
	SelectMemberIDs(ctx context.Context, chatID int64) ([]int64, error)
}

func NewChatFuncService(log *zap.Logger, storage ChatStorage, cache Cache) *ChatFuncService {
	return &ChatFuncService{
		log:     log,
		storage: storage,
		cache:   cache,
	}
}

func (s *ChatFuncService) CreateChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Chat{}, err
	}

	chat, err := s.storage.InsertChat(ctx, req)
	if err != nil {
		log.Error("failed to create chat", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	if err := s.cache.SetChat(ctx, &chat); err != nil {
		log.Error("failed to set chat cache", zap.Error(err))
	}

	for _, userID := range req.MemberIDs {
		if err := s.cache.DeleteUserChats(ctx, userID); err != nil {
			log.Error("failed to invalidate user chats cache", zap.Error(err))
		}
	}

	return chat, nil
}

func (s *ChatFuncService) GetChat(ctx context.Context, req dto.GetChatRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Chat{}, err
	}

	if cached, err := s.cache.GetChat(ctx, req.ChatID); err == nil && cached != nil {
		if !req.IncludeMembers {
			cached.MembersCount = 0
		}

		return *cached, nil
	}

	chat, err := s.storage.SelectChat(ctx, req.ChatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("chat not found", zap.Error(err))
			return models.Chat{}, domain.ErrChatNotFound
		}

		log.Error("failed to select chat", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	if err := s.cache.SetChat(ctx, &chat); err != nil {
		log.Error("failed to set chat cache", zap.Error(err))
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
		return dto.GetChatsResponse{}, domain.ErrInternalStorage
	}

	return chats, nil
}

func (s *ChatFuncService) UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Chat{}, err
	}

	chat, err := s.storage.UpdateChat(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("chat not found", zap.Error(err))
			return models.Chat{}, domain.ErrChatNotFound
		}

		log.Error("failed to update chat", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	if err := s.cache.SetChat(ctx, &chat); err != nil {
		log.Error("failed to set chat cache", zap.Error(err))
	}

	memberIDs, err := s.storage.SelectMemberIDs(ctx, req.ChatID)
	if err == nil {
		for _, userID := range memberIDs {
			if err := s.cache.DeleteUserChats(ctx, userID); err != nil {
				log.Error("failed to invalidate user chats cache", zap.Error(err))
			}
		}
	}

	return chat, nil
}
