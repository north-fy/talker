package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"github.com/north-fy/talker/services/message/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FeatureService interface {
	SearchMessages(ctx context.Context, req dto.SearchMessagesRequest) (dto.SearchMessagesResponse, error)
	MarkAsRead(ctx context.Context, req dto.MarkAsReadRequest) error
	GetUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error)
	ConnectWebSocket(ctx context.Context, req dto.ConnectWebSocketRequest) error
	GetLastMessage(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error)
	DeleteChatMessages(ctx context.Context, req dto.DeleteChatMessagesRequest) error
}

func (s *serverAPI) SearchMessages(ctx context.Context, req *messagev1.SearchMessagesRequest) (*messagev1.SearchMessagesResponse, error) {
	msgReq := dto.SearchMessagesRequest{
		ChatID: req.GetChatId(),
		Query:  req.GetQuery(),
		Limit:  req.GetLimit(),
		Before: req.GetBefore(),
	}

	msgResp, err := s.serv.SearchMessages(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	messages := utils.ParallelConvert(msgResp.Messages, workers)

	return &messagev1.SearchMessagesResponse{
		Messages: messages,
		HasMore:  msgResp.HasMore,
	}, nil
}

func (s *serverAPI) MarkAsRead(ctx context.Context, req *messagev1.MarkAsReadRequest) (*messagev1.Empty, error) {
	msgReq := dto.MarkAsReadRequest{
		ChatID:        req.GetChatId(),
		UserID:        req.GetUserId(),
		UpToMessageID: req.GetUpToMessageId(),
	}

	err := s.serv.MarkAsRead(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return &messagev1.Empty{}, nil
}

func (s *serverAPI) GetUnreadCount(ctx context.Context, req *messagev1.GetUnreadCountRequest) (*messagev1.GetUnreadCountResponse, error) {
	msgReq := dto.GetUnreadCountRequest{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
	}

	msgResp, err := s.serv.GetUnreadCount(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return &messagev1.GetUnreadCountResponse{
		Count:         msgResp.Count,
		LastMessageId: msgResp.LastMessageID,
	}, nil
}

func (s *serverAPI) GetLastMessage(ctx context.Context, req *messagev1.GetLastMessageRequest) (*messagev1.Message, error) {
	msgReq := dto.GetLastMessageRequest{
		ChatID: req.GetChatId(),
	}

	msgResp, err := s.serv.GetLastMessage(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return utils.ConvertToProtoMessage(&msgResp), nil
}

func (s *serverAPI) DeleteChatMessages(ctx context.Context, req *messagev1.DeleteChatMessagesRequest) (*messagev1.Empty, error) {
	msgReq := dto.DeleteChatMessagesRequest{
		ChatID: req.GetChatId(),
	}

	if err := s.serv.DeleteChatMessages(ctx, msgReq); err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return nil, nil
}
