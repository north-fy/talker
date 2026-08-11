package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/pkg/convert"
	"go.uber.org/zap"
)

type FeatureService struct {
	log           *zap.Logger
	storage       FeatStorage
	messageClient messagev1.MessageServiceClient
	cache         Cache
}

type FeatStorage interface {
	GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
	GetChatsByUser(ctx context.Context, userID int64) ([]*dto.UserChatDB, error)
}

func NewFeatureService(log *zap.Logger, storage FeatStorage, messageClient messagev1.MessageServiceClient, cache Cache) *FeatureService {
	return &FeatureService{
		log:           log,
		storage:       storage,
		messageClient: messageClient,
		cache:         cache,
	}
}

func (s *FeatureService) IsMember(ctx context.Context, req dto.IsMemberRequest) (dto.IsMemberResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.IsMemberResponse{}, err
	}

	storageRequest := dto.GetMemberRequest{
		ChatID: req.ChatID,
		UserID: req.UserID,
	}

	if cached, err := s.cache.GetMember(ctx, req.ChatID, req.UserID); err == nil && cached != nil {
		return dto.IsMemberResponse{
			IsMember: true,
			Role:     cached.Role,
		}, nil
	}

	member, err := s.storage.GetMember(ctx, storageRequest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.IsMemberResponse{}, nil
		}

		log.Error("failed to get member", zap.Error(err))
		return dto.IsMemberResponse{}, domain.ErrInternalStorage
	}

	if err := s.cache.SetMember(ctx, &member); err != nil {
		log.Error("failed to set member cache", zap.Error(err))
	}

	return dto.IsMemberResponse{
		IsMember: true,
		Role:     member.Role,
	}, nil
}

func (s *FeatureService) GetUserChats(ctx context.Context, req dto.GetUserChatsRequest) (dto.GetUserChatsResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.GetUserChatsResponse{}, err
	}

	// Кешируем только вариант без последнего сообщения, т.к. оно меняется часто.
	if !req.IncludeLastMessage {
		if cached, err := s.cache.GetUserChats(ctx, req.UserID); err == nil && cached != nil {
			return *cached, nil
		}
	}

	userChats, err := s.storage.GetChatsByUser(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get chats by user", zap.Error(err))
		return dto.GetUserChatsResponse{}, domain.ErrInternalStorage
	}

	resp := s.buildUserChatsResponse(ctx, userChats, req.IncludeLastMessage)

	if !req.IncludeLastMessage {
		if err := s.cache.SetUserChats(ctx, req.UserID, &resp); err != nil {
			log.Error("failed to set user chats cache", zap.Error(err))
		}
	}

	return resp, nil
}

func (s *FeatureService) buildUserChatsResponse(ctx context.Context, userChats []*dto.UserChatDB, includeLastMessage bool) dto.GetUserChatsResponse {
	resp := dto.GetUserChatsResponse{
		UserChats:  make([]*dto.UserChatResponse, 0, len(userChats)),
		TotalCount: int64(len(userChats)),
	}

	for _, uc := range userChats {
		item := &dto.UserChatResponse{
			Chat:        &uc.Chat,
			MemberInfo:  convert.ConvertMemberDBToModel(&uc.Member),
			UnreadCount: uc.Member.UnreadCount,
		}

		if includeLastMessage {
			item.LastMessage = s.getLastMessage(ctx, uc.Chat.ID)
		}

		resp.UserChats = append(resp.UserChats, item)
	}

	return resp
}

func (s *FeatureService) getLastMessage(ctx context.Context, chatID int64) *dto.MessageResponse {
	lm, err := s.messageClient.GetLastMessage(ctx, &messagev1.GetLastMessageRequest{
		ChatId: chatID,
	})
	if err != nil {
		s.log.Debug("failed to get last message",
			zap.Int64("chat_id", chatID),
			zap.Error(err))
		return nil
	}

	return &dto.MessageResponse{
		ID:        lm.Id,
		ChatID:    lm.ChatId,
		SenderID:  lm.SenderId,
		Content:   lm.Content,
		CreatedAt: lm.CreatedAt.AsTime(),
	}
}
