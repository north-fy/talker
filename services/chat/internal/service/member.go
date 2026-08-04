package service

import (
	"context"

	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/pkg/convert"
	"go.uber.org/zap"
)

type MemberService struct {
	log        *zap.Logger
	storage    Storage
	userClient userv1.UserServiceClient
}

type Storage interface {
	AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error)
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error
	UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
	GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersDBResponse, error)
}

func NewMemberService(log *zap.Logger, storage Storage) *MemberService {
	return &MemberService{
		log:     log,
		storage: storage,
	}
}

func (m *MemberService) AddMember(ctx context.Context, req dto.AddMemberRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	member, err := m.storage.AddMember(ctx, req)
	if err != nil {
		log.Error("failed to add member", zap.Error(err))
		return models.Member{}, err
	}

	userResp, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: []int64{req.UserID},
	})
	if err != nil {
		log.Error("failed to get user", zap.Error(err))
		return models.Member{}, err
	}

	return *convert.ConvertMemberToDTO(&member, userResp.Users[0]), nil
}

func (m *MemberService) RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error {
	log := m.log.With(zap.Any("request", req))

	if err := m.storage.RemoveMember(ctx, req); err != nil {
		log.Error("failed to remove member", zap.Error(err))
		return err
	}

	return nil
}

func (m *MemberService) UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	member, err := m.storage.UpdateMemberRole(ctx, req)
	if err != nil {
		log.Error("failed to update member role", zap.Error(err))
		return models.Member{}, err
	}

	userResp, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: []int64{req.UserID},
	})
	if err != nil {
		log.Error("failed to get info user", zap.Error(err))
		return models.Member{}, err
	}

	return *convert.ConvertMemberToDTO(&member, userResp.Users[0]), nil
}

func (m *MemberService) GetMember(ctx context.Context, req dto.GetMemberRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	member, err := m.storage.GetMember(ctx, req)
	if err != nil {
		log.Error("failed to get member", zap.Error(err))
		return models.Member{}, err
	}

	userResp, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: []int64{req.UserID},
	})
	if err != nil {
		log.Error("failed to get info user", zap.Error(err))
		return models.Member{}, err
	}

	return *convert.ConvertMemberToDTO(&member, userResp.Users[0]), nil
}

func (m *MemberService) GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersResponse, error) {
	log := m.log.With(zap.Any("request", req))

	members, err := m.storage.GetMembers(ctx, req)
	if err != nil {
		log.Error("failed to add member", zap.Error(err))
		return dto.GetMembersResponse{}, err
	}

	users, err := m.userClient.GetUsers(ctx, &userv1.GetUsersRequest{
		UserIds: members.IDs,
	})
	if err != nil {
		log.Error("failed to get info users", zap.Error(err))
		return dto.GetMembersResponse{}, err
	}

	DTOmembers := make([]*models.Member, 0, len(members.Members))
	for i, member := range members.Members {
		DTOmembers = append(DTOmembers, convert.ConvertMemberToDTO(member, users.Users[i]))
	}

	return dto.GetMembersResponse{
		Members:    DTOmembers,
		TotalCount: int64(len(DTOmembers)),
	}, nil
}
