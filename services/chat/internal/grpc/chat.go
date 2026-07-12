package grpc

import (
	"context"

	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/pkg/convert"
)

type ChatFuncService interface {
	CreateChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error)
	GetChat(ctx context.Context, req dto.GetChatRequest) (models.Chat, error)
	GetChats(ctx context.Context, req dto.GetChatsRequest) (dto.GetChatsResponse, error)
	UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error)
}

func (s *serverAPI) CreateChat(ctx context.Context, req *chatv1.CreateChatRequest) (*chatv1.Chat, error) {
	chatType := int32(req.GetType())
	chatReq := dto.CreateChatRequest{
		Name:         req.GetName(),
		Type:         dto.ChatType(chatType),
		MemberIDs:    req.GetMemberIds(),
		AvatarBase64: req.GetAvatarBase64(),
	}

	resp, err := s.serv.CreateChat(ctx, chatReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertChatToProto(&resp), nil
}

func (s *serverAPI) GetChat(ctx context.Context, req *chatv1.GetChatRequest) (*chatv1.Chat, error) {
	chatReq := dto.GetChatRequest{
		ChatID:         req.GetChatId(),
		IncludeMembers: req.GetIncludeMembers(),
	}

	resp, err := s.serv.GetChat(ctx, chatReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertChatToProto(&resp), nil
}

func (s *serverAPI) GetChats(ctx context.Context, req *chatv1.GetChatsRequest) (*chatv1.GetChatsResponse, error) {
	protoFilter := req.GetFilter()
	chatType := int32(protoFilter.GetType())

	chatFilter := dto.ChatFilter{
		Type:            dto.ChatType(chatType),
		Search:          protoFilter.GetSearch(),
		ChatIDs:         protoFilter.GetChatIds(),
		IncludeArchived: protoFilter.GetIncludeArchived(),
	}

	chatReq := dto.GetChatsRequest{Filter: chatFilter}

	resp, err := s.serv.GetChats(ctx, chatReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	// TODO: протестировать по бенчмаркам как влияает на скорость
	var protoChats []*chatv1.Chat
	for _, chat := range resp.Chats {
		protoChats = append(protoChats, convert.ConvertChatToProto(chat))
	}

	return &chatv1.GetChatsResponse{
		Chats:      protoChats,
		TotalCount: resp.TotalCount,
	}, nil
}

func (s *serverAPI) UpdateChat(ctx context.Context, req *chatv1.UpdateChatRequest) (*chatv1.Chat, error) {
	name := req.GetName()
	avatar := req.GetAvatarBase64()

	chatReq := dto.UpdateChatRequest{
		ChatID: req.GetChatId(),
		Name:   &name,
		AvatarBase64: &avatar,
	}

	chat, err := s.serv.UpdateChat(ctx, chatReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertChatToProto(&chat), nil
}
