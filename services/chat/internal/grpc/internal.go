package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

type InternalService interface {
	GetChatInternal(ctx context.Context, req dto.GetChatInternalRequest) (any, error)
	GetChatsInternal(ctx context.Context, req dto.GetChatsInternalRequest) (dto.GetChatsInternalResponse, error)
	ValidateMemberAccess(ctx context.Context, req dto.ValidateMemberAccessRequest) (dto.ValidateMemberAccessResponse, error)
}

func (s *serverAPI) GetChatInternal(ctx context.Context, req *messagev1.GetChatInternalRequest) (*messagev1.ChatInternal, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) GetChatsInternal(ctx context.Context, req *messagev1.GetChatsInternalRequest) (*messagev1.GetChatsInternalResponse, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) ValidateMemberAccess(ctx context.Context, req *messagev1.ValidateMemberAccessRequest) (*messagev1.ValidateMemberAccessResponse, error) {
	// TODO
	return nil, nil
}
