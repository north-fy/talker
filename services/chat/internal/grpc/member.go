package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

type MemberService interface {
	AddMember(ctx context.Context, req dto.AddMemberRequest) (any, error)
	RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error
	GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersResponse, error)
	GetMember(ctx context.Context, req dto.GetMemberRequest) (any, error)
	UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (any, error)
	IsMember(ctx context.Context, req dto.GetMemberRequest) (bool, string, error)
	GetUserChats(ctx context.Context, req dto.GetUserChatsRequest) (dto.GetUserChatsResponse, error)
	LeaveChat(ctx context.Context, req dto.RemoveMemberRequest) error
}

func (s *serverAPI) AddMember(ctx context.Context, req *messagev1.AddMemberRequest) (*messagev1.Member, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) RemoveMember(ctx context.Context, req *messagev1.RemoveMemberRequest) (*messagev1.Empty, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) GetMembers(ctx context.Context, req *messagev1.GetMembersRequest) (*messagev1.GetMembersResponse, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) GetMember(ctx context.Context, req *messagev1.GetMemberRequest) (*messagev1.Member, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) UpdateMemberRole(ctx context.Context, req *messagev1.UpdateMemberRoleRequest) (*messagev1.Member, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) IsMember(ctx context.Context, req *messagev1.IsMemberRequest) (*messagev1.IsMemberResponse, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) GetUserChats(ctx context.Context, req *messagev1.GetUserChatsRequest) (*messagev1.GetUserChatsResponse, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) LeaveChat(ctx context.Context, req *messagev1.LeaveChatRequest) (*messagev1.Empty, error) {
	// TODO
	return nil, nil
}
