package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/internal/mocks"
	"go.uber.org/zap"
)

func newInternalService(storage *mocks.MockInternalStorage) *InternalService {
	return NewInternalService(zap.NewNop(), storage)
}

func TestInternalService_GetChatInternal(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{ID: 1, Name: "group", Type: int32(dto.ChatTypeGroup), IsActive: true}, nil
		},
		SelectMemberIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return []int64{10, 20}, nil
		},
		SelectSettingsFn: func(_ context.Context, _ int64) (dto.ChatSettings, error) {
			return dto.ChatSettings{Language: "ru", AllowMedia: true}, nil
		},
	}

	svc := newInternalService(storage)
	resp, err := svc.GetChatInternal(context.Background(), dto.GetChatInternalRequest{ChatID: 1, IncludeMembers: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != 1 || len(resp.MemberIDs) != 2 || resp.Settings.Language != "ru" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestInternalService_GetChatInternal_NotFound(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{}, pgx.ErrNoRows
		},
	}

	svc := newInternalService(storage)
	_, err := svc.GetChatInternal(context.Background(), dto.GetChatInternalRequest{ChatID: 1})
	if err != domain.ErrChatNotFound {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestInternalService_GetChatsInternal(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		SelectChatsByIDsFn: func(_ context.Context, ids []int64) ([]*models.Chat, error) {
			return []*models.Chat{
				{ID: 1, Name: "one", Type: int32(dto.ChatTypeGroup), IsActive: true},
				{ID: 2, Name: "two", Type: int32(dto.ChatTypePrivate), IsActive: true},
			}, nil
		},
		SelectMembersByChatIDsFn: func(_ context.Context, _ []int64) (map[int64][]int64, error) {
			return map[int64][]int64{1: {10}, 2: {20, 30}}, nil
		},
		SelectSettingsByChatIDsFn: func(_ context.Context, _ []int64) (map[int64]dto.ChatSettings, error) {
			return map[int64]dto.ChatSettings{1: {Language: "ru"}}, nil
		},
	}

	svc := newInternalService(storage)
	resp, err := svc.GetChatsInternal(context.Background(), dto.GetChatsInternalRequest{ChatIDs: []int64{1, 2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Chats) != 2 {
		t.Fatalf("unexpected count: %d", len(resp.Chats))
	}
	if len(resp.Chats[2].MemberIDs) != 2 {
		t.Fatalf("unexpected member ids for chat 2: %+v", resp.Chats[2].MemberIDs)
	}
}

func TestInternalService_ValidateMemberAccess_Read(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleMember}, nil
		},
	}

	svc := newInternalService(storage)
	resp, err := svc.ValidateMemberAccess(context.Background(), dto.ValidateMemberAccessRequest{
		ChatID:             1,
		UserID:             42,
		RequiredPermission: dto.PermissionRead,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.HasAccess {
		t.Fatal("member should have read access")
	}
}

func TestInternalService_ValidateMemberAccess_DeleteDenied(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleMember}, nil
		},
	}

	svc := newInternalService(storage)
	resp, err := svc.ValidateMemberAccess(context.Background(), dto.ValidateMemberAccessRequest{
		ChatID:             1,
		UserID:             42,
		RequiredPermission: dto.PermissionDelete,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.HasAccess {
		t.Fatal("member should not have delete access")
	}
	if resp.Reason == "" {
		t.Fatal("expected reason to be set")
	}
}

func TestInternalService_ValidateMemberAccess_ManageMembers(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleAdmin}, nil
		},
	}

	svc := newInternalService(storage)
	resp, err := svc.ValidateMemberAccess(context.Background(), dto.ValidateMemberAccessRequest{
		ChatID:             1,
		UserID:             42,
		RequiredPermission: dto.PermissionManageMembers,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.HasAccess {
		t.Fatal("admin should have manage members access")
	}
}

func TestInternalService_ValidateMemberAccess_NotMember(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInternalStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{}, pgx.ErrNoRows
		},
	}

	svc := newInternalService(storage)
	resp, err := svc.ValidateMemberAccess(context.Background(), dto.ValidateMemberAccessRequest{
		ChatID:             1,
		UserID:             42,
		RequiredPermission: dto.PermissionRead,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.HasAccess {
		t.Fatal("expected no access")
	}
}
