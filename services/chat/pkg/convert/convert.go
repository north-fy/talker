package convert

import (
	"fmt"
	"strconv"

	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
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

func ConvertMemberToProto(m *models.Member) *chatv1.Member {
	r, _ := strconv.Atoi(m.Role)

	return &chatv1.Member{
		UserId: m.UserID,
		ChatId: m.ChatID,
		Role: chatv1.Role(r),
		JoinedAt: timestamppb.New(m.JoinedAt),
		LastReadAt: timestamppb.New(m.LastReadAt),
		UnreadCount: m.UnreadCount,
		Username: m.Username,
		FullName: m.FullName,
		AvatarUrl: m.AvatarURL,
	}
}

func ConvertMemberToDTO(memberDB *dto.MemberDB, memberUser *userv1.User) *models.Member {
	return &models.Member{
		UserID:      memberDB.UserID,
		ChatID:      memberDB.ChatID,
		Role:        string(memberDB.Role),
		JoinedAt:    memberDB.JoinedAt,
		LastReadAt:  memberDB.LastReadAt,
		UnreadCount: memberDB.UnreadCount,
		Username:    memberUser.Username,
		FullName:    fmt.Sprintf("%s %s", memberUser.FirstName, memberUser.LastName),
	}
}