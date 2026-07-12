package dto

import (
	"time"

	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

type ChatType int32

const (
	ChatTypeUnknown ChatType = iota
	ChatTypePrivate
	ChatTypeGroup
	ChatTypeChannel
)

type Role int32

const (
	RoleUnknown Role = iota
	RoleMember
	RoleModerator
	RoleAdmin
	RoleOwner
)

type PermissionType int32

const (
	PermissionUnknown PermissionType = iota
	PermissionRead
	PermissionWrite
	PermissionDelete
	PermissionManageMembers
	PermissionUpdateSettings
)

type CreateChatRequest struct {
	Name         string   `validate:"required,max=255"`
	Type         ChatType `validate:"required"`
	MemberIDs    []int64  `validate:"required,min=1"`
	AvatarBase64 string
}

type GetChatRequest struct {
	ChatID         int64 `validate:"required"`
	IncludeMembers bool
}

type GetChatsRequest struct {
	Filter ChatFilter
}

type ChatFilter struct {
	Type            ChatType
	Search          string
	ChatIDs         []int64
	IncludeArchived bool
}

type GetChatsResponse struct {
	Chats      []*models.Chat
	TotalCount int64
}

type UpdateChatRequest struct {
	ChatID       int64   `validate:"required"`
	Name         *string `validate:"omitempty,max=255"`
	AvatarBase64 *string
}

type AddMemberRequest struct {
	ChatID    int64 `validate:"required"`
	UserID    int64 `validate:"required"`
	Role      Role  `validate:"required,oneof=1 2 3 4"`
	InvitedBy int64 `validate:"required"`
}

type RemoveMemberRequest struct {
	ChatID int64 `validate:"required"`
	UserID int64 `validate:"required"`
}

type GetMembersRequest struct {
	ChatID int64 `validate:"required"`
	Filter MemberFilter
}

type MemberFilter struct {
	Role   Role
	Search string
}

type GetMembersResponse struct {
	Members    []*models.Member
	TotalCount int64
}

type UpdateMemberRoleRequest struct {
	ChatID int64 `validate:"required"`
	UserID int64 `validate:"required"`
	Role   Role  `validate:"required"`
}

type GetMemberRequest struct {
	ChatID int64 `validate:"required"`
	UserID int64 `validate:"required"`
}

type IsMemberRequest struct {
	ChatID int64 `validate:"required"`
	UserID int64 `validate:"required"`
}

type IsMemberResponse struct {
	IsMember bool
	Role     Role
}

type GetUserChatsRequest struct {
	UserID             int64 `validate:"required"`
	IncludeLastMessage bool
}

type GetUserChatsResponse struct {
	UserChats  []*UserChatResponse
	TotalCount int64
}

type UserChatResponse struct {
	Chat        *models.Chat
	MemberInfo  *models.Member
	LastMessage *MessageResponse
	UnreadCount int64
}

type CreateInviteLinkRequest struct {
	ChatID    int64 `validate:"required"`
	MaxUses   int32 `validate:"min=0"`
	ExpiresAt *time.Time
}

type JoinChatByInviteRequest struct {
	InviteCode string `validate:"required"`
	UserID     int64  `validate:"required"`
}

type RevokeInviteLinkRequest struct {
	ChatID   int64 `validate:"required"`
	InviteID int64 `validate:"required"`
}

type UpdateChatSettingsRequest struct {
	ChatID   int64 `validate:"required"`
	Settings ChatSettings
}

type GetChatSettingsRequest struct {
	ChatID int64 `validate:"required"`
}

type ChatSettings struct {
	IsPrivate            bool
	AllowMessagesFromAll bool
	AllowMedia           bool
	AllowReactions       bool
	MessageTTLSeconds    int32  `validate:"min=0"`
	Language             string `validate:"max=10"`
	IsAnnouncement       bool
}

type GetChatInternalRequest struct {
	ChatID         int64 `validate:"required"`
	IncludeMembers bool
}

type ChatInternalResponse struct {
	ID        int64
	Name      string
	Type      ChatType
	IsActive  bool
	MemberIDs []int64
	Settings  ChatSettings
}

type GetChatsInternalRequest struct {
	ChatIDs []int64 `validate:"required,min=1"`
}

type GetChatsInternalResponse struct {
	Chats map[int64]*ChatInternalResponse
}

type ValidateMemberAccessRequest struct {
	ChatID             int64          `validate:"required"`
	UserID             int64          `validate:"required"`
	RequiredPermission PermissionType `validate:"required"`
}

type ValidateMemberAccessResponse struct {
	HasAccess bool
	Role      Role
	Reason    string
}

type UpdateLastReadRequest struct {
	ChatID    int64 `validate:"required"`
	UserID    int64 `validate:"required"`
	MessageID int64 `validate:"required"`
}

type MessageResponse struct {
	ID        int64
	ChatID    int64
	SenderID  int64
	Content   string
	CreatedAt time.Time
}
