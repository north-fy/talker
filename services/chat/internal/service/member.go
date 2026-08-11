package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/pkg/convert"
	"go.uber.org/zap"
)

type MemberService struct {
	log        *zap.Logger
	storage    MemberStorage
	userClient userv1.UserServiceClient
	cache      Cache
}

type MemberStorage interface {
	AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error)
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error
	UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
	GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersDBResponse, error)
	SelectChat(ctx context.Context, chatID int64) (models.Chat, error)
}

func NewMemberService(log *zap.Logger, storage MemberStorage, userClient userv1.UserServiceClient, cache Cache) *MemberService {
	return &MemberService{
		log:        log,
		storage:    storage,
		userClient: userClient,
		cache:      cache,
	}
}

func (m *MemberService) AddMember(ctx context.Context, req dto.AddMemberRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Member{}, err
	}

	if _, err := m.storage.SelectChat(ctx, req.ChatID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Member{}, domain.ErrChatNotFound
		}

		log.Error("failed to select chat", zap.Error(err))
		return models.Member{}, domain.ErrInternalStorage
	}

	member, err := m.storage.AddMember(ctx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Member{}, domain.ErrMemberAlreadyInChat
		}

		log.Error("failed to add member", zap.Error(err))
		return models.Member{}, domain.ErrInternalStorage
	}

	if err := m.cache.SetMember(ctx, &member); err != nil {
		log.Error("failed to set member cache", zap.Error(err))
	}
	if err := m.cache.DeleteUserChats(ctx, req.UserID); err != nil {
		log.Error("failed to invalidate user chats cache", zap.Error(err))
	}

	userResp, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: []int64{req.UserID},
	})
	if err != nil {
		log.Error("failed to get user", zap.Error(err))
		return models.Member{}, err
	}

	return *convert.ConvertMemberToDTO(&member, m.userByID(userResp, req.UserID)), nil
}

func (m *MemberService) RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error {
	log := m.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return err
	}

	current, err := m.storage.GetMember(ctx, dto.GetMemberRequest{ChatID: req.ChatID, UserID: req.UserID})
	if err == nil && current.Role == dto.RoleOwner {
		return domain.ErrCannotRemoveOwner
	}

	if err := m.storage.RemoveMember(ctx, req); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrMemberNotInChat
		}

		log.Error("failed to remove member", zap.Error(err))
		return domain.ErrInternalStorage
	}

	if err := m.cache.DeleteMember(ctx, req.ChatID, req.UserID); err != nil {
		log.Error("failed to invalidate member cache", zap.Error(err))
	}
	if err := m.cache.DeleteUserChats(ctx, req.UserID); err != nil {
		log.Error("failed to invalidate user chats cache", zap.Error(err))
	}

	return nil
}

func (m *MemberService) UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Member{}, err
	}

	current, err := m.storage.GetMember(ctx, dto.GetMemberRequest{ChatID: req.ChatID, UserID: req.UserID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Member{}, domain.ErrMemberNotInChat
		}

		log.Error("failed to get member", zap.Error(err))
		return models.Member{}, domain.ErrInternalStorage
	}

	if current.Role == dto.RoleOwner {
		return models.Member{}, domain.ErrCannotModifyOwner
	}

	member, err := m.storage.UpdateMemberRole(ctx, req)
	if err != nil {
		log.Error("failed to update member role", zap.Error(err))
		return models.Member{}, domain.ErrInternalStorage
	}

	if err := m.cache.SetMember(ctx, &member); err != nil {
		log.Error("failed to set member cache", zap.Error(err))
	}
	if err := m.cache.DeleteUserChats(ctx, req.UserID); err != nil {
		log.Error("failed to invalidate user chats cache", zap.Error(err))
	}

	userResp, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: []int64{req.UserID},
	})
	if err != nil {
		log.Error("failed to get info user", zap.Error(err))
		return models.Member{}, err
	}

	return *convert.ConvertMemberToDTO(&member, m.userByID(userResp, req.UserID)), nil
}

func (m *MemberService) GetMember(ctx context.Context, req dto.GetMemberRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Member{}, err
	}

	member, err := m.getMemberCached(ctx, req)
	if err != nil {
		return models.Member{}, err
	}

	userResp, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: []int64{req.UserID},
	})
	if err != nil {
		log.Error("failed to get info user", zap.Error(err))
		return models.Member{}, err
	}

	return *convert.ConvertMemberToDTO(member, m.userByID(userResp, req.UserID)), nil
}

func (m *MemberService) GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersResponse, error) {
	log := m.log.With(zap.Any("request", req))

	if err := validateStruct(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.GetMembersResponse{}, err
	}

	members, err := m.storage.GetMembers(ctx, req)
	if err != nil {
		log.Error("failed to get members", zap.Error(err))
		return dto.GetMembersResponse{}, domain.ErrInternalStorage
	}

	users, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: members.IDs,
	})
	if err != nil {
		log.Error("failed to get info users", zap.Error(err))
		return dto.GetMembersResponse{}, err
	}

	userMap := make(map[int64]*userv1.User, len(users.Users))
	for _, u := range users.Users {
		userMap[u.UserId] = u
	}

	DTOmembers := make([]*models.Member, 0, len(members.Members))
	for _, member := range members.Members {
		user, ok := userMap[member.UserID]
		if !ok {
			user = &userv1.User{}
		}

		model := convert.ConvertMemberToDTO(member, user)
		if req.Filter.Search != "" && !strings.Contains(strings.ToLower(model.FullName), strings.ToLower(req.Filter.Search)) {
			continue
		}

		DTOmembers = append(DTOmembers, model)
	}

	return dto.GetMembersResponse{
		Members:    DTOmembers,
		TotalCount: int64(len(DTOmembers)),
	}, nil
}

// getMemberCached возвращает участника через cache-aside.
func (m *MemberService) getMemberCached(ctx context.Context, req dto.GetMemberRequest) (*dto.MemberDB, error) {
	log := m.log.With(zap.Any("request", req))

	if cached, err := m.cache.GetMember(ctx, req.ChatID, req.UserID); err == nil && cached != nil {
		return cached, nil
	}

	member, err := m.storage.GetMember(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMemberNotInChat
		}

		log.Error("failed to get member", zap.Error(err))
		return nil, domain.ErrInternalStorage
	}

	if err := m.cache.SetMember(ctx, &member); err != nil {
		log.Error("failed to set member cache", zap.Error(err))
	}

	return &member, nil
}

func (m *MemberService) userByID(resp *userv1.GetUsersResponse, userID int64) *userv1.User {
	for _, u := range resp.Users {
		if u.UserId == userID {
			return u
		}
	}

	return &userv1.User{}
}
