package grpc

import (
	"context"
	"time"

	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/pkg/convert"
)

type InviteService interface {
	CreateInviteLink(ctx context.Context, req dto.CreateInviteLinkRequest) (models.InviteLink, error)
	JoinChatByInvite(ctx context.Context, req dto.JoinChatByInviteRequest) (models.Chat, error)
	RevokeInviteLink(ctx context.Context, req dto.RevokeInviteLinkRequest) error
}

func (s *serverAPI) CreateInviteLink(ctx context.Context, req *chatv1.CreateInviteLinkRequest) (*chatv1.InviteLink, error) {
	var expires *time.Time
	if ts := req.GetExpiresAt(); ts != nil {
		t := ts.AsTime()
		expires = &t
	}

	invReq := dto.CreateInviteLinkRequest{
		ChatID:    req.GetChatId(),
		MaxUses:   req.GetMaxUses(),
		ExpiresAt: expires,
	}

	invite, err := s.serv.CreateInviteLink(ctx, invReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertInviteToProto(&invite), nil
}

func (s *serverAPI) JoinChatByInvite(ctx context.Context, req *chatv1.JoinChatByInviteRequest) (*chatv1.Chat, error) {
	invReq := dto.JoinChatByInviteRequest{
		InviteCode: req.GetInviteCode(),
		UserID:     req.GetUserId(),
	}

	chat, err := s.serv.JoinChatByInvite(ctx, invReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertChatToProto(&chat), nil
}

func (s *serverAPI) RevokeInviteLink(ctx context.Context, req *chatv1.RevokeInviteLinkRequest) (*chatv1.Empty, error) {
	invReq := dto.RevokeInviteLinkRequest{
		ChatID:   req.GetChatId(),
		InviteID: req.GetInviteId(),
	}

	if err := s.serv.RevokeInviteLink(ctx, invReq); err != nil {
		return nil, toGRPC(err)
	}

	return nil, nil
}
