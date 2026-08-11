package grpc

import (
	"context"

	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/pkg/convert"
)

type FeatureService interface {
	IsMember(ctx context.Context, req dto.IsMemberRequest) (dto.IsMemberResponse, error)
	GetUserChats(ctx context.Context, req dto.GetUserChatsRequest) (dto.GetUserChatsResponse, error)
}

func (s *serverAPI) IsMember(ctx context.Context, req *chatv1.IsMemberRequest) (*chatv1.IsMemberResponse, error) {
	mReq := dto.IsMemberRequest{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
	}

	resp, err := s.serv.IsMember(ctx, mReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &chatv1.IsMemberResponse{
		IsMember: resp.IsMember,
		Role:     chatv1.Role(resp.Role),
	}, nil
}

func (s *serverAPI) GetUserChats(ctx context.Context, req *chatv1.GetUserChatsRequest) (*chatv1.GetUserChatsResponse, error) {
	mReq := dto.GetUserChatsRequest{
		UserID:             req.GetUserId(),
		IncludeLastMessage: req.GetIncludeLastMessage(),
	}

	resp, err := s.serv.GetUserChats(ctx, mReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	userChats := make([]*chatv1.UserChat, 0, len(resp.UserChats))
	for _, val := range resp.UserChats {
		chat := &chatv1.UserChat{}

		chat.Chat = convert.ConvertChatToProto(val.Chat)
		chat.MemberInfo = convert.ConvertMemberToProto(val.MemberInfo)
		chat.LastMessage = convert.ConvertMessageToProto(val.LastMessage)
		chat.UnreadCount = val.UnreadCount

		userChats = append(userChats, chat)
	}

	return &chatv1.GetUserChatsResponse{
		UserChats:  userChats,
		TotalCount: resp.TotalCount,
	}, nil
}
