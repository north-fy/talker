package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"go.uber.org/zap"
)

const inviteBaseURL = "https://talker.app/join/"

const inviteCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type InviteService struct {
	log     *zap.Logger
	storage InviteStorage
	cache   Cache
}

type InviteStorage interface {
	SelectChat(ctx context.Context, chatID int64) (models.Chat, error)
	InsertInvite(ctx context.Context, req dto.CreateInviteLinkRequest, createdBy int64, code string) (models.InviteLink, error)
	SelectInvite(ctx context.Context, id int64) (models.InviteLink, error)
	SelectInviteByCode(ctx context.Context, code string) (models.InviteLink, error)
	IncrementUsedCount(ctx context.Context, id int64) (models.InviteLink, error)
	DeactivateInvite(ctx context.Context, chatID, inviteID int64) error
	AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
}

func NewInviteService(log *zap.Logger, storage InviteStorage, cache Cache) *InviteService {
	return &InviteService{
		log:     log,
		storage: storage,
		cache:   cache,
	}
}

func (s *InviteService) CreateInviteLink(ctx context.Context, req dto.CreateInviteLinkRequest) (models.InviteLink, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.InviteLink{}, err
	}

	if _, err := s.storage.SelectChat(ctx, req.ChatID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.InviteLink{}, domain.ErrChatNotFound
		}

		log.Error("failed to select chat", zap.Error(err))
		return models.InviteLink{}, domain.ErrInternalStorage
	}

	// TODO: взять создателя из контекста аутентификации
	const createdBy = 0

	code, err := generateInviteCode()
	if err != nil {
		log.Error("failed to generate invite code", zap.Error(err))
		return models.InviteLink{}, domain.ErrInternalStorage
	}

	invite, err := s.storage.InsertInvite(ctx, req, createdBy, code)
	if err != nil {
		log.Error("failed to insert invite", zap.Error(err))
		return models.InviteLink{}, domain.ErrInternalStorage
	}

	invite.URL = inviteBaseURL + invite.Code

	return invite, nil
}

func (s *InviteService) JoinChatByInvite(ctx context.Context, req dto.JoinChatByInviteRequest) (models.Chat, error) {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Chat{}, err
	}

	invite, err := s.storage.SelectInviteByCode(ctx, req.InviteCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Chat{}, domain.ErrInviteNotFound
		}

		log.Error("failed to select invite", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	if err := s.checkInviteValid(invite); err != nil {
		return models.Chat{}, err
	}

	if _, err := s.storage.IncrementUsedCount(ctx, invite.ID); err != nil {
		// Гонка: лимит мог быть исчерпан в другом запросе
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Chat{}, domain.ErrInviteMaxUsesReached
		}

		log.Error("failed to increment invite used count", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	if _, err := s.storage.GetMember(ctx, dto.GetMemberRequest{ChatID: invite.ChatID, UserID: req.UserID}); err == nil {
		return models.Chat{}, domain.ErrMemberAlreadyInChat
	}

	_, err = s.storage.AddMember(ctx, dto.AddMemberRequest{
		ChatID: invite.ChatID,
		UserID: req.UserID,
		Role:   dto.RoleMember,
	})
	if err != nil {
		log.Error("failed to add member by invite", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	if err := s.cache.DeleteUserChats(ctx, req.UserID); err != nil {
		log.Error("failed to invalidate user chats cache", zap.Error(err))
	}

	chat, err := s.storage.SelectChat(ctx, invite.ChatID)
	if err != nil {
		log.Error("failed to select chat", zap.Error(err))
		return models.Chat{}, domain.ErrInternalStorage
	}

	return chat, nil
}

func (s *InviteService) RevokeInviteLink(ctx context.Context, req dto.RevokeInviteLinkRequest) error {
	log := s.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return err
	}

	if err := s.storage.DeactivateInvite(ctx, req.ChatID, req.InviteID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInviteNotFound
		}

		log.Error("failed to revoke invite", zap.Error(err))
		return domain.ErrInternalStorage
	}

	return nil
}

func (s *InviteService) checkInviteValid(invite models.InviteLink) error {
	if !invite.IsActive {
		return domain.ErrInviteRevoked
	}

	if !invite.ExpiresAt.IsZero() && invite.ExpiresAt.Before(time.Now()) {
		return domain.ErrInviteExpired
	}

	if invite.MaxUses > 0 && invite.UsedCount >= invite.MaxUses {
		return domain.ErrInviteMaxUsesReached
	}

	return nil
}

func generateInviteCode() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	for i := range buf {
		buf[i] = inviteCodeCharset[int(buf[i])%len(inviteCodeCharset)]
	}

	return string(buf), nil
}
