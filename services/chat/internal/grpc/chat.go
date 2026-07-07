package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

type ChatService interface {
	CreateChat(ctx context.Context, req dto.CreateChatRequest, creatorID string) (any, error)
	GetChat(ctx context.Context, req dto.GetChatRequest) (any, error)
	GetChats(ctx context.Context, req dto.GetChatsRequest) (dto.GetChatsResponse, error)
	UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (any, error)
	DeleteChat(ctx context.Context, req dto.DeleteChatRequest) error
}

func (s *serverAPI) CreateChat(ctx context.Context, req *messagev1.CreateChatRequest) (*messagev1.Chat, error) {
	// TODO: convert proto -> dto, call service, convert result -> proto
	return nil, nil
}

func (s *serverAPI) GetChat(ctx context.Context, req *messagev1.GetChatRequest) (*messagev1.Chat, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) GetChats(ctx context.Context, req *messagev1.GetChatsRequest) (*messagev1.GetChatsResponse, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) UpdateChat(ctx context.Context, req *messagev1.UpdateChatRequest) (*messagev1.Chat, error) {
	// TODO
	return nil, nil
}

func (s *serverAPI) DeleteChat(ctx context.Context, req *messagev1.DeleteChatRequest) (*messagev1.Empty, error) {
	// TODO
	return nil, nil
}
