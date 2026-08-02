package service

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"go.uber.org/zap"
)

type MemberService struct {
	log     *zap.Logger
	storage Storage
}

type Storage interface {
	AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error)
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error
	UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
	GetMembers(ctx context.Context, req dto.GetMembersRequest) ([]*dto.MemberDB, error)
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

	return member, nil
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

	return member, nil
}

func (m *MemberService) GetMember(ctx context.Context, req dto.GetMemberRequest) (models.Member, error) {
	log := m.log.With(zap.Any("request", req))

	member, err := m.storage.GetMember(ctx, req)
	if err != nil {
		log.Error("failed to get member", zap.Error(err))
		return models.Member{}, err
	}

	return member, nil
}

func (m *MemberService) GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersResponse, error) {
	log := m.log.With(zap.Any("request", req))

	members, err := m.storage.GetMembers(ctx, req)
	if err != nil {
		log.Error("failed to add member", zap.Error(err))
		return dto.GetMembersResponse{}, err
	}

	return dto.GetMembersResponse{
		Members:    members,
		TotalCount: int64(len(members)),
	}, nil
}
