package grpc

import (
	"context"

	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/pkg/convert"
)

type MemberService interface {
	AddMember(ctx context.Context, req dto.AddMemberRequest) (models.Member, error)
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error
	UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (models.Member, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (models.Member, error)
	GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersResponse, error)
}

func (s *serverAPI) AddMember(ctx context.Context, req *chatv1.AddMemberRequest) (*chatv1.Member, error) {
	mReq := dto.AddMemberRequest{
		ChatID:    req.GetChatId(),
		UserID:    req.GetUserId(),
		Role:      dto.Role(req.GetRole()),
		InvitedBy: req.GetInvitedBy(),
	}

	resp, err := s.serv.AddMember(ctx, mReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertMemberToProto(&resp), nil
}

func (s *serverAPI) RemoveMember(ctx context.Context, req *chatv1.RemoveMemberRequest) (*chatv1.Empty, error) {
	mReq := dto.RemoveMemberRequest{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
	}

	if err := s.serv.RemoveMember(ctx, mReq); err != nil {
		return nil, toGRPC(err)
	}

	return nil, nil
}

func (s *serverAPI) UpdateMemberRole(ctx context.Context, req *chatv1.UpdateMemberRoleRequest) (*chatv1.Member, error) {
	mReq := dto.UpdateMemberRoleRequest{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
		Role:   dto.Role(req.GetRole()),
	}

	resp, err := s.serv.UpdateMemberRole(ctx, mReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertMemberToProto(&resp), nil
}

func (s *serverAPI) GetMember(ctx context.Context, req *chatv1.GetMemberRequest) (*chatv1.Member, error) {
	mReq := dto.GetMemberRequest{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
	}

	resp, err := s.serv.GetMember(ctx, mReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertMemberToProto(&resp), nil
}

func (s *serverAPI) GetMembers(ctx context.Context, req *chatv1.GetMembersRequest) (*chatv1.GetMembersResponse, error) {
	filter := dto.MemberFilter{
		Role:   dto.Role(req.GetFilter().GetRole()),
		Search: req.GetFilter().GetSearch(),
	}

	mReq := dto.GetMembersRequest{
		ChatID: req.GetChatId(),
		Filter: filter,
	}

	resp, err := s.serv.GetMembers(ctx, mReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	members := make([]*chatv1.Member, 0, len(resp.Members))
	for _, member := range resp.Members {
		members = append(members, convert.ConvertMemberToProto(member))
	}

	return &chatv1.GetMembersResponse{
		Members:    members,
		TotalCount: resp.TotalCount,
	}, nil
}
