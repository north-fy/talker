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
		Id:           chat.ID,
		Name:         chat.Name,
		Type:         chatv1.ChatType(chat.Type),
		CreatedBy:    chat.CreatedBy,
		AvatarUrl:    chat.AvatarURL,
		MembersCount: chat.MembersCount,
		LastMessage:  nil,
		CreatedAt:    timestamppb.New(chat.CreatedAt),
		UpdatedAt:    timestamppb.New(chat.UpdatedAt),
	}
}

func ConvertMemberToProto(m *models.Member) *chatv1.Member {
	r, _ := strconv.Atoi(m.Role)

	return &chatv1.Member{
		UserId:      m.UserID,
		ChatId:      m.ChatID,
		Role:        chatv1.Role(r),
		JoinedAt:    timestamppb.New(m.JoinedAt),
		LastReadAt:  timestamppb.New(m.LastReadAt),
		UnreadCount: m.UnreadCount,
		Username:    m.Username,
		FullName:    m.FullName,
		AvatarUrl:   m.AvatarURL,
	}
}

func ConvertMemberToDTO(memberDB *dto.MemberDB, memberUser *userv1.User) *models.Member {
	return &models.Member{
		UserID:      memberDB.UserID,
		ChatID:      memberDB.ChatID,
		Role:        strconv.Itoa(int(memberDB.Role)),
		JoinedAt:    memberDB.JoinedAt,
		LastReadAt:  memberDB.LastReadAt,
		UnreadCount: memberDB.UnreadCount,
		Username:    memberUser.Username,
		FullName:    fmt.Sprintf("%s %s", memberUser.FirstName, memberUser.LastName),
	}
}

func ConvertMemberDBToModel(memberDB *dto.MemberDB) *models.Member {
	return &models.Member{
		UserID:      memberDB.UserID,
		ChatID:      memberDB.ChatID,
		Role:        strconv.Itoa(int(memberDB.Role)),
		JoinedAt:    memberDB.JoinedAt,
		LastReadAt:  memberDB.LastReadAt,
		UnreadCount: memberDB.UnreadCount,
	}
}

func ConvertMessageToProto(m *dto.MessageResponse) *chatv1.Message {
	return &chatv1.Message{
		Id:        m.ID,
		ChatId:    m.ChatID,
		SenderId:  m.SenderID,
		Content:   m.Content,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}
}

func ConvertInviteToProto(invite *models.InviteLink) *chatv1.InviteLink {
	return &chatv1.InviteLink{
		Id:        invite.ID,
		ChatId:    invite.ChatID,
		Code:      invite.Code,
		Url:       invite.URL,
		MaxUses:   invite.MaxUses,
		UsedCount: invite.UsedCount,
		ExpiresAt: timestamppb.New(invite.ExpiresAt),
		CreatedAt: timestamppb.New(invite.CreatedAt),
		CreatedBy: invite.CreatedBy,
		IsActive:  invite.IsActive,
	}
}

func ConvertSettingsToProto(settings *dto.ChatSettings) *chatv1.ChatSettings {
	return &chatv1.ChatSettings{
		IsPrivate:            settings.IsPrivate,
		AllowMessagesFromAll: settings.AllowMessagesFromAll,
		AllowMedia:           settings.AllowMedia,
		AllowReactions:       settings.AllowReactions,
		MessageTtlSeconds:    settings.MessageTTLSeconds,
		Language:             settings.Language,
		IsAnnouncement:       settings.IsAnnouncement,
	}
}

func ConvertChatInternalToProto(in *dto.ChatInternalResponse) *chatv1.ChatInternal {
	return &chatv1.ChatInternal{
		Id:        in.ID,
		Name:      in.Name,
		Type:      chatv1.ChatType(in.Type),
		IsActive:  in.IsActive,
		MemberIds: in.MemberIDs,
		Settings:  ConvertSettingsToProto(&in.Settings),
	}
}
