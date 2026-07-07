package dto

import (
	"time"

	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

type CreateChatRequest struct {
	Name            string   `validate:"required"`
	Type            string   `validate:"required"`
	MemberIDs       []string `validate:"required"`
	AvatarBase64    string
}

type GetChatRequest struct {
	ChatID          string `validate:"required"`
	IncludeMembers  bool
}

type GetChatsRequest struct {
	Type            string
	Search          string
	ChatIDs         []string
	IncludeArchived bool
}

type GetChatsResponse struct {
	Chats       []models.Chat
	TotalCount  int64
}

type UpdateChatRequest struct {
	ChatID      string  `validate:"required"`
	Name        *string
	AvatarBase64 *string
}

type DeleteChatRequest struct {
	ChatID string `validate:"required"`
}

// ==================== Member ====================

type AddMemberRequest struct {
	ChatID    string `validate:"required"`
	UserID    string `validate:"required"`
	Role      string
	InvitedBy string
}

type RemoveMemberRequest struct {
	ChatID string `validate:"required"`
	UserID string `validate:"required"`
}

type GetMembersRequest struct {
	ChatID     string `validate:"required"`
	Role       string
	Search     string
}

type GetMembersResponse struct {
	Members     []models.Member
	TotalCount  int64
}

type UpdateMemberRoleRequest struct {
	ChatID string `validate:"required"`
	UserID string `validate:"required"`
	Role   string `validate:"required"`
}

type GetMemberRequest struct {
	ChatID string `validate:"required"`
	UserID string `validate:"required"`
}

type GetUserChatsRequest struct {
	UserID              string `validate:"required"`
	IncludeLastMessage  bool
}

type GetUserChatsResponse struct {
	UserChats   []models.UserChat
	TotalCount  int64
}

// ==================== Invite ====================

type CreateInviteLinkRequest struct {
	ChatID    string `validate:"required"`
	MaxUses   int32
	ExpiresAt time.Time
}

type JoinChatByInviteRequest struct {
	InviteCode string `validate:"required"`
	UserID     string
}

type RevokeInviteLinkRequest struct {
	ChatID   string `validate:"required"`
	InviteID string `validate:"required"`
}

// ==================== Internal ====================

type GetChatInternalRequest struct {
	ChatID         string `validate:"required"`
	IncludeMembers bool
}

type GetChatsInternalRequest struct {
	ChatIDs []string `validate:"required"`
}

type GetChatsInternalResponse struct {
	Chats map[string]models.Chat
}

type ValidateMemberAccessRequest struct {
	ChatID   string `validate:"required"`
	UserID   string `validate:"required"`
	Permission string
}

type ValidateMemberAccessResponse struct {
	HasAccess bool
	Role      string
	Reason    string
}
