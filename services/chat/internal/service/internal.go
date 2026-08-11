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

type InternalService struct {
	log     *zap.Logger
	storage InternalStorage
}

type InternalStorage interface {
	SelectChat(ctx context.Context, chatID int64) (models.Chat, error)
	SelectChatsByIDs(ctx context.Context, ids []int64) ([]*models.Chat, error)
	SelectMemberIDs(ctx context.Context, chatID int64) ([]int64, error)
	SelectMembersByChatIDs(ctx context.Context, chatIDs []int64) (map[int64][]int64, error)
	SelectSettings(ctx context.Context, chatID int64) (dto.ChatSettings, error)
	SelectSettingsByChatIDs(ctx context.Context, chatIDs []int64) (map[int64]dto.ChatSettings, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
}

func NewInternalService(log *zap.Logger, storage InternalStorage) *InternalService {
	return &InternalService{
		log:     log,
		storage: storage,
	}
}

func (s *InternalService) GetChatInternal(ctx context.Context, req dto.GetChatInternalRequest) (dto.ChatInternalResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.ChatInternalResponse{}, err
	}

	chat, err := s.storage.SelectChat(ctx, req.ChatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ChatInternalResponse{}, domain.ErrChatNotFound
		}

		log.Error("failed to select chat", zap.Error(err))
		return dto.ChatInternalResponse{}, domain.ErrInternalStorage
	}

	resp := dto.ChatInternalResponse{
		ID:       chat.ID,
		Name:     chat.Name,
		Type:     dto.ChatType(chat.Type),
		IsActive: chat.IsActive,
	}

	if req.IncludeMembers {
		ids, err := s.storage.SelectMemberIDs(ctx, req.ChatID)
		if err != nil {
			log.Error("failed to select member ids", zap.Error(err))
			return dto.ChatInternalResponse{}, domain.ErrInternalStorage
		}
		resp.MemberIDs = ids
	}

	settings, err := s.storage.SelectSettings(ctx, req.ChatID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error("failed to select chat settings", zap.Error(err))
		return dto.ChatInternalResponse{}, domain.ErrInternalStorage
	}
	resp.Settings = settings

	return resp, nil
}

func (s *InternalService) GetChatsInternal(ctx context.Context, req dto.GetChatsInternalRequest) (dto.GetChatsInternalResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.GetChatsInternalResponse{}, err
	}

	chats, err := s.storage.SelectChatsByIDs(ctx, req.ChatIDs)
	if err != nil {
		log.Error("failed to select chats", zap.Error(err))
		return dto.GetChatsInternalResponse{}, domain.ErrInternalStorage
	}

	ids := make([]int64, 0, len(chats))
	for _, chat := range chats {
		ids = append(ids, chat.ID)
	}

	memberMap, err := s.storage.SelectMembersByChatIDs(ctx, ids)
	if err != nil {
		log.Error("failed to select members", zap.Error(err))
		return dto.GetChatsInternalResponse{}, domain.ErrInternalStorage
	}

	settingsMap, err := s.storage.SelectSettingsByChatIDs(ctx, ids)
	if err != nil {
		log.Error("failed to select settings", zap.Error(err))
		return dto.GetChatsInternalResponse{}, domain.ErrInternalStorage
	}

	result := make(map[int64]*dto.ChatInternalResponse, len(chats))
	for _, chat := range chats {
		result[chat.ID] = &dto.ChatInternalResponse{
			ID:        chat.ID,
			Name:      chat.Name,
			Type:      dto.ChatType(chat.Type),
			IsActive:  chat.IsActive,
			MemberIDs: memberMap[chat.ID],
			Settings:  settingsMap[chat.ID],
		}
	}

	return dto.GetChatsInternalResponse{Chats: result}, nil
}

func (s *InternalService) ValidateMemberAccess(ctx context.Context, req dto.ValidateMemberAccessRequest) (dto.ValidateMemberAccessResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.ValidateMemberAccessResponse{}, err
	}

	member, err := s.storage.GetMember(ctx, dto.GetMemberRequest{ChatID: req.ChatID, UserID: req.UserID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ValidateMemberAccessResponse{
				HasAccess: false,
				Reason:    "user is not a member of this chat",
			}, nil
		}

		log.Error("failed to get member", zap.Error(err))
		return dto.ValidateMemberAccessResponse{}, domain.ErrInternalStorage
	}

	if hasPermission(member.Role, req.RequiredPermission) {
		return dto.ValidateMemberAccessResponse{
			HasAccess: true,
			Role:      member.Role,
		}, nil
	}

	return dto.ValidateMemberAccessResponse{
		HasAccess: false,
		Role:      member.Role,
		Reason:    "insufficient role to perform this action",
	}, nil
}

func hasPermission(role dto.Role, permission dto.PermissionType) bool {
	switch permission {
	case dto.PermissionRead, dto.PermissionWrite:
		return role >= dto.RoleMember
	case dto.PermissionDelete:
		return role >= dto.RoleModerator
	case dto.PermissionManageMembers, dto.PermissionUpdateSettings:
		return role >= dto.RoleAdmin
	default:
		return false
	}
}
