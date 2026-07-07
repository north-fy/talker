package service

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/internal/storage"
	"go.uber.org/zap"
)

type ChatService struct {
	log    *zap.Logger
	chat   storage.ChatStorage
	member storage.MemberStorage
	invite storage.InviteStorage
}

func NewChatService(log *zap.Logger, chat storage.ChatStorage, member storage.MemberStorage, invite storage.InviteStorage) *ChatService {
	return &ChatService{
		log:    log,
		chat:   chat,
		member: member,
		invite: invite,
	}
}

// ==================== Chat ====================

func (s *ChatService) CreateChat(ctx context.Context, req dto.CreateChatRequest, creatorID string) (models.Chat, error) {
	// TODO: validate, create chat, add creator as OWNER
	return models.Chat{}, nil
}

func (s *ChatService) GetChat(ctx context.Context, req dto.GetChatRequest) (models.Chat, error) {
	// TODO: check access, get chat
	return models.Chat{}, nil
}

func (s *ChatService) GetChats(ctx context.Context, req dto.GetChatsRequest) (dto.GetChatsResponse, error) {
	// TODO: list chats with filters
	return dto.GetChatsResponse{}, nil
}

func (s *ChatService) UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error) {
	// TODO: check admin/owner role, update chat
	return models.Chat{}, nil
}

func (s *ChatService) DeleteChat(ctx context.Context, req dto.DeleteChatRequest) error {
	// TODO: check owner role, delete chat and all members
	return nil
}

// ==================== Member ====================

func (s *ChatService) AddMember(ctx context.Context, req dto.AddMemberRequest) (models.Member, error) {
	// TODO: check admin role, check not already member, add member
	return models.Member{}, nil
}

func (s *ChatService) RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error {
	// TODO: check admin role, cannot remove owner, remove member
	return nil
}

func (s *ChatService) GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersResponse, error) {
	// TODO: list members with filters
	return dto.GetMembersResponse{}, nil
}

func (s *ChatService) GetMember(ctx context.Context, req dto.GetMemberRequest) (models.Member, error) {
	// TODO: get member by chat_id + user_id
	return models.Member{}, nil
}

func (s *ChatService) UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (models.Member, error) {
	// TODO: check admin/owner role, update role
	return models.Member{}, nil
}

func (s *ChatService) IsMember(ctx context.Context, req dto.GetMemberRequest) (bool, string, error) {
	// TODO: check membership
	return false, "", nil
}

func (s *ChatService) GetUserChats(ctx context.Context, req dto.GetUserChatsRequest) (dto.GetUserChatsResponse, error) {
	// TODO: get all chats for user
	return dto.GetUserChatsResponse{}, nil
}

func (s *ChatService) LeaveChat(ctx context.Context, req dto.RemoveMemberRequest) error {
	// TODO: user leaves chat (cannot leave if owner)
	return nil
}

// ==================== Invite ====================

func (s *ChatService) CreateInviteLink(ctx context.Context, req dto.CreateInviteLinkRequest, creatorID string) (models.InviteLink, error) {
	// TODO: check admin role, generate code, create link
	return models.InviteLink{}, nil
}

func (s *ChatService) JoinChatByInvite(ctx context.Context, req dto.JoinChatByInviteRequest) (models.Chat, error) {
	// TODO: validate invite, check not expired/max uses, add member
	return models.Chat{}, nil
}

func (s *ChatService) RevokeInviteLink(ctx context.Context, req dto.RevokeInviteLinkRequest) error {
	// TODO: check admin role, revoke invite
	return nil
}

// ==================== Internal ====================

func (s *ChatService) GetChatInternal(ctx context.Context, req dto.GetChatInternalRequest) (models.Chat, error) {
	// TODO: internal method for other services
	return models.Chat{}, nil
}

func (s *ChatService) GetChatsInternal(ctx context.Context, req dto.GetChatsInternalRequest) (dto.GetChatsInternalResponse, error) {
	// TODO: batch get chats for other services
	return dto.GetChatsInternalResponse{}, nil
}

func (s *ChatService) ValidateMemberAccess(ctx context.Context, req dto.ValidateMemberAccessRequest) (dto.ValidateMemberAccessResponse, error) {
	// TODO: validate access for message service
	return dto.ValidateMemberAccessResponse{}, nil
}
