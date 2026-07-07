package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

type InviteService interface {
	CreateInviteLink(ctx context.Context, req dto.CreateInviteLinkRequest, creatorID string) (any, error)
	JoinChatByInvite(ctx context.Context, req dto.JoinChatByInviteRequest) (any, error)
	RevokeInviteLink(ctx context.Context, req dto.RevokeInviteLinkRequest) error
}

func (s *serverAPI) CreateInviteLink(ctx context.Context, req *messagev1.CreateInviteLinkRequest) (*messagev1.InviteLink, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) JoinChatByInvite(ctx context.Context, req *messagev1.JoinChatByInviteRequest) (*messagev1.Chat, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) RevokeInviteLink(ctx context.Context, req *messagev1.RevokeInviteLinkRequest) (*messagev1.Empty, error) {
	// TODO
	return nil, nil
}
