package mocks

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

// MockChatStorage — мок ChatStorage.
type MockChatStorage struct {
	InsertChatFn      func(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error)
	SelectChatFn      func(ctx context.Context, chatID int64) (models.Chat, error)
	SelectChatsFn     func(ctx context.Context, filter dto.ChatFilter) (dto.GetChatsResponse, error)
	UpdateChatFn      func(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error)
	SelectMemberIDsFn func(ctx context.Context, chatID int64) ([]int64, error)
}

func (m *MockChatStorage) InsertChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error) {
	if m.InsertChatFn != nil {
		return m.InsertChatFn(ctx, req)
	}
	return models.Chat{}, nil
}

func (m *MockChatStorage) SelectChat(ctx context.Context, chatID int64) (models.Chat, error) {
	if m.SelectChatFn != nil {
		return m.SelectChatFn(ctx, chatID)
	}
	return models.Chat{}, nil
}

func (m *MockChatStorage) SelectChats(ctx context.Context, filter dto.ChatFilter) (dto.GetChatsResponse, error) {
	if m.SelectChatsFn != nil {
		return m.SelectChatsFn(ctx, filter)
	}
	return dto.GetChatsResponse{}, nil
}

func (m *MockChatStorage) UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error) {
	if m.UpdateChatFn != nil {
		return m.UpdateChatFn(ctx, req)
	}
	return models.Chat{}, nil
}

func (m *MockChatStorage) SelectMemberIDs(ctx context.Context, chatID int64) ([]int64, error) {
	if m.SelectMemberIDsFn != nil {
		return m.SelectMemberIDsFn(ctx, chatID)
	}
	return nil, nil
}

// MockMemberStorage — мок MemberStorage.
type MockMemberStorage struct {
	AddMemberFn        func(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error)
	RemoveMemberFn     func(ctx context.Context, req dto.RemoveMemberRequest) error
	UpdateMemberRoleFn func(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error)
	GetMemberFn        func(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
	GetMembersFn       func(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersDBResponse, error)
	SelectChatFn       func(ctx context.Context, chatID int64) (models.Chat, error)
}

func (m *MockMemberStorage) AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error) {
	if m.AddMemberFn != nil {
		return m.AddMemberFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

func (m *MockMemberStorage) RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error {
	if m.RemoveMemberFn != nil {
		return m.RemoveMemberFn(ctx, req)
	}
	return nil
}

func (m *MockMemberStorage) UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error) {
	if m.UpdateMemberRoleFn != nil {
		return m.UpdateMemberRoleFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

func (m *MockMemberStorage) GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error) {
	if m.GetMemberFn != nil {
		return m.GetMemberFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

func (m *MockMemberStorage) GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersDBResponse, error) {
	if m.GetMembersFn != nil {
		return m.GetMembersFn(ctx, req)
	}
	return dto.GetMembersDBResponse{}, nil
}

func (m *MockMemberStorage) SelectChat(ctx context.Context, chatID int64) (models.Chat, error) {
	if m.SelectChatFn != nil {
		return m.SelectChatFn(ctx, chatID)
	}
	return models.Chat{}, nil
}

// MockFeatStorage — мок FeatStorage.
type MockFeatStorage struct {
	GetMemberFn      func(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
	GetChatsByUserFn func(ctx context.Context, userID int64) ([]*dto.UserChatDB, error)
}

func (m *MockFeatStorage) GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error) {
	if m.GetMemberFn != nil {
		return m.GetMemberFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

func (m *MockFeatStorage) GetChatsByUser(ctx context.Context, userID int64) ([]*dto.UserChatDB, error) {
	if m.GetChatsByUserFn != nil {
		return m.GetChatsByUserFn(ctx, userID)
	}
	return nil, nil
}

// MockInternalStorage — мок InternalStorage.
type MockInternalStorage struct {
	SelectChatFn              func(ctx context.Context, chatID int64) (models.Chat, error)
	SelectChatsByIDsFn        func(ctx context.Context, ids []int64) ([]*models.Chat, error)
	SelectMemberIDsFn         func(ctx context.Context, chatID int64) ([]int64, error)
	SelectMembersByChatIDsFn  func(ctx context.Context, chatIDs []int64) (map[int64][]int64, error)
	SelectSettingsFn          func(ctx context.Context, chatID int64) (dto.ChatSettings, error)
	SelectSettingsByChatIDsFn func(ctx context.Context, chatIDs []int64) (map[int64]dto.ChatSettings, error)
	GetMemberFn               func(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
}

func (m *MockInternalStorage) SelectChat(ctx context.Context, chatID int64) (models.Chat, error) {
	if m.SelectChatFn != nil {
		return m.SelectChatFn(ctx, chatID)
	}
	return models.Chat{}, nil
}

func (m *MockInternalStorage) SelectChatsByIDs(ctx context.Context, ids []int64) ([]*models.Chat, error) {
	if m.SelectChatsByIDsFn != nil {
		return m.SelectChatsByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (m *MockInternalStorage) SelectMemberIDs(ctx context.Context, chatID int64) ([]int64, error) {
	if m.SelectMemberIDsFn != nil {
		return m.SelectMemberIDsFn(ctx, chatID)
	}
	return nil, nil
}

func (m *MockInternalStorage) SelectMembersByChatIDs(ctx context.Context, chatIDs []int64) (map[int64][]int64, error) {
	if m.SelectMembersByChatIDsFn != nil {
		return m.SelectMembersByChatIDsFn(ctx, chatIDs)
	}
	return nil, nil
}

func (m *MockInternalStorage) SelectSettings(ctx context.Context, chatID int64) (dto.ChatSettings, error) {
	if m.SelectSettingsFn != nil {
		return m.SelectSettingsFn(ctx, chatID)
	}
	return dto.ChatSettings{}, nil
}

func (m *MockInternalStorage) SelectSettingsByChatIDs(ctx context.Context, chatIDs []int64) (map[int64]dto.ChatSettings, error) {
	if m.SelectSettingsByChatIDsFn != nil {
		return m.SelectSettingsByChatIDsFn(ctx, chatIDs)
	}
	return nil, nil
}

func (m *MockInternalStorage) GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error) {
	if m.GetMemberFn != nil {
		return m.GetMemberFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

// MockInviteStorage — мок InviteStorage.
type MockInviteStorage struct {
	SelectChatFn         func(ctx context.Context, chatID int64) (models.Chat, error)
	InsertInviteFn       func(ctx context.Context, req dto.CreateInviteLinkRequest, createdBy int64, code string) (models.InviteLink, error)
	SelectInviteFn       func(ctx context.Context, id int64) (models.InviteLink, error)
	SelectInviteByCodeFn func(ctx context.Context, code string) (models.InviteLink, error)
	IncrementUsedCountFn func(ctx context.Context, id int64) (models.InviteLink, error)
	DeactivateInviteFn   func(ctx context.Context, chatID, inviteID int64) error
	AddMemberFn          func(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error)
	GetMemberFn          func(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error)
}

func (m *MockInviteStorage) SelectChat(ctx context.Context, chatID int64) (models.Chat, error) {
	if m.SelectChatFn != nil {
		return m.SelectChatFn(ctx, chatID)
	}
	return models.Chat{}, nil
}

func (m *MockInviteStorage) InsertInvite(ctx context.Context, req dto.CreateInviteLinkRequest, createdBy int64, code string) (models.InviteLink, error) {
	if m.InsertInviteFn != nil {
		return m.InsertInviteFn(ctx, req, createdBy, code)
	}
	return models.InviteLink{}, nil
}

func (m *MockInviteStorage) SelectInvite(ctx context.Context, id int64) (models.InviteLink, error) {
	if m.SelectInviteFn != nil {
		return m.SelectInviteFn(ctx, id)
	}
	return models.InviteLink{}, nil
}

func (m *MockInviteStorage) SelectInviteByCode(ctx context.Context, code string) (models.InviteLink, error) {
	if m.SelectInviteByCodeFn != nil {
		return m.SelectInviteByCodeFn(ctx, code)
	}
	return models.InviteLink{}, nil
}

func (m *MockInviteStorage) IncrementUsedCount(ctx context.Context, id int64) (models.InviteLink, error) {
	if m.IncrementUsedCountFn != nil {
		return m.IncrementUsedCountFn(ctx, id)
	}
	return models.InviteLink{}, nil
}

func (m *MockInviteStorage) DeactivateInvite(ctx context.Context, chatID, inviteID int64) error {
	if m.DeactivateInviteFn != nil {
		return m.DeactivateInviteFn(ctx, chatID, inviteID)
	}
	return nil
}

func (m *MockInviteStorage) AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error) {
	if m.AddMemberFn != nil {
		return m.AddMemberFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

func (m *MockInviteStorage) GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error) {
	if m.GetMemberFn != nil {
		return m.GetMemberFn(ctx, req)
	}
	return dto.MemberDB{}, nil
}

// MockCache — мок Cache.
type MockCache struct {
	GetChatFn         func(ctx context.Context, chatID int64) (*models.Chat, error)
	SetChatFn         func(ctx context.Context, chat *models.Chat) error
	DeleteChatFn      func(ctx context.Context, chatID int64) error
	GetUserChatsFn    func(ctx context.Context, userID int64) (*dto.GetUserChatsResponse, error)
	SetUserChatsFn    func(ctx context.Context, userID int64, resp *dto.GetUserChatsResponse) error
	DeleteUserChatsFn func(ctx context.Context, userID int64) error
	GetMemberFn       func(ctx context.Context, chatID, userID int64) (*dto.MemberDB, error)
	SetMemberFn       func(ctx context.Context, member *dto.MemberDB) error
	DeleteMemberFn    func(ctx context.Context, chatID, userID int64) error
}

func (m *MockCache) GetChat(ctx context.Context, chatID int64) (*models.Chat, error) {
	if m.GetChatFn != nil {
		return m.GetChatFn(ctx, chatID)
	}
	return nil, nil
}

func (m *MockCache) SetChat(ctx context.Context, chat *models.Chat) error {
	if m.SetChatFn != nil {
		return m.SetChatFn(ctx, chat)
	}
	return nil
}

func (m *MockCache) DeleteChat(ctx context.Context, chatID int64) error {
	if m.DeleteChatFn != nil {
		return m.DeleteChatFn(ctx, chatID)
	}
	return nil
}

func (m *MockCache) GetUserChats(ctx context.Context, userID int64) (*dto.GetUserChatsResponse, error) {
	if m.GetUserChatsFn != nil {
		return m.GetUserChatsFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockCache) SetUserChats(ctx context.Context, userID int64, resp *dto.GetUserChatsResponse) error {
	if m.SetUserChatsFn != nil {
		return m.SetUserChatsFn(ctx, userID, resp)
	}
	return nil
}

func (m *MockCache) DeleteUserChats(ctx context.Context, userID int64) error {
	if m.DeleteUserChatsFn != nil {
		return m.DeleteUserChatsFn(ctx, userID)
	}
	return nil
}

func (m *MockCache) GetMember(ctx context.Context, chatID, userID int64) (*dto.MemberDB, error) {
	if m.GetMemberFn != nil {
		return m.GetMemberFn(ctx, chatID, userID)
	}
	return nil, nil
}

func (m *MockCache) SetMember(ctx context.Context, member *dto.MemberDB) error {
	if m.SetMemberFn != nil {
		return m.SetMemberFn(ctx, member)
	}
	return nil
}

func (m *MockCache) DeleteMember(ctx context.Context, chatID, userID int64) error {
	if m.DeleteMemberFn != nil {
		return m.DeleteMemberFn(ctx, chatID, userID)
	}
	return nil
}
