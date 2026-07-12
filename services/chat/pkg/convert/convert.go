package convert

import (
	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertChatToProto(chat *models.Chat) *chatv1.Chat {
	return &chatv1.Chat{
		Id: chat.ID,
		Name: chat.Name,
		Type: chatv1.ChatType(chat.Type),
		CreatedBy: chat.CreatedBy,
		AvatarUrl: chat.AvatarURL,
		MembersCount: chat.MembersCount,
		LastMessage: nil,
		CreatedAt: timestamppb.New(chat.CreatedAt),
		UpdatedAt: timestamppb.New(chat.UpdatedAt),
	}
}
