package grpc

import (
	"context"
	"strconv"

	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/pkg/convert"
)

type InternalService interface {
	GetChatInternal(ctx context.Context, req dto.GetChatInternalRequest) (dto.ChatInternalResponse, error)
	GetChatsInternal(ctx context.Context, req dto.GetChatsInternalRequest) (dto.GetChatsInternalResponse, error)
	ValidateMemberAccess(ctx context.Context, req dto.ValidateMemberAccessRequest) (dto.ValidateMemberAccessResponse, error)
}

func (s *serverAPI) GetChatInternal(ctx context.Context, req *chatv1.GetChatInternalRequest) (*chatv1.ChatInternal, error) {
	chatReq := dto.GetChatInternalRequest{
		ChatID:         req.GetChatId(),
		IncludeMembers: req.GetIncludeMembers(),
	}

	resp, err := s.serv.GetChatInternal(ctx, chatReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertChatInternalToProto(&resp), nil
}

func (s *serverAPI) GetChatsInternal(ctx context.Context, req *chatv1.GetChatsInternalRequest) (*chatv1.GetChatsInternalResponse, error) {
	chatReq := dto.GetChatsInternalRequest{
		ChatIDs: req.GetChatIds(),
	}

	resp, err := s.serv.GetChatsInternal(ctx, chatReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	chats := make(map[string]*chatv1.ChatInternal, len(resp.Chats))
	for id, chat := range resp.Chats {
		chats[strconv.FormatInt(id, 10)] = convert.ConvertChatInternalToProto(chat)
	}

	return &chatv1.GetChatsInternalResponse{
		Chats: chats,
	}, nil
}

func (s *serverAPI) ValidateMemberAccess(ctx context.Context, req *chatv1.ValidateMemberAccessRequest) (*chatv1.ValidateMemberAccessResponse, error) {
	accReq := dto.ValidateMemberAccessRequest{
		ChatID:             req.GetChatId(),
		UserID:             req.GetUserId(),
		RequiredPermission: dto.PermissionType(req.GetRequiredPermission()),
	}

	resp, err := s.serv.ValidateMemberAccess(ctx, accReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &chatv1.ValidateMemberAccessResponse{
		HasAccess: resp.HasAccess,
		Role:      chatv1.Role(resp.Role),
		Reason:    resp.Reason,
	}, nil
}
