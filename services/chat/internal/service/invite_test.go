package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/internal/mocks"
	"go.uber.org/zap"
)

func newInviteService(storage *mocks.MockInviteStorage, cache *mocks.MockCache) *InviteService {
	if cache == nil {
		cache = &mocks.MockCache{}
	}
	return NewInviteService(zap.NewNop(), storage, cache)
}

func TestInviteService_CreateInviteLink(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockInviteStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{ID: 1}, nil
		},
		InsertInviteFn: func(_ context.Context, req dto.CreateInviteLinkRequest, _ int64, code string) (models.InviteLink, error) {
			return models.InviteLink{ID: 5, ChatID: req.ChatID, Code: code, MaxUses: req.MaxUses, CreatedAt: now}, nil
		},
	}

	svc := newInviteService(storage, nil)
	invite, err := svc.CreateInviteLink(context.Background(), dto.CreateInviteLinkRequest{ChatID: 1, MaxUses: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invite.ID != 5 || invite.Code == "" {
		t.Fatalf("unexpected invite: %+v", invite)
	}
	if invite.URL != inviteBaseURL+invite.Code {
		t.Fatalf("unexpected url: %q", invite.URL)
	}
}

func TestInviteService_CreateInviteLink_ChatNotFound(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{}, pgx.ErrNoRows
		},
	}

	svc := newInviteService(storage, nil)
	_, err := svc.CreateInviteLink(context.Background(), dto.CreateInviteLinkRequest{ChatID: 1})
	if err != domain.ErrChatNotFound {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestInviteService_JoinChatByInvite(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		SelectInviteByCodeFn: func(_ context.Context, code string) (models.InviteLink, error) {
			return models.InviteLink{ID: 1, ChatID: 1, Code: code, MaxUses: 0, IsActive: true}, nil
		},
		IncrementUsedCountFn: func(_ context.Context, id int64) (models.InviteLink, error) {
			return models.InviteLink{ID: id, ChatID: 1, UsedCount: 1}, nil
		},
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{}, pgx.ErrNoRows
		},
		AddMemberFn: func(_ context.Context, req dto.AddMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: req.ChatID, UserID: req.UserID, Role: dto.RoleMember}, nil
		},
		SelectChatFn: func(_ context.Context, chatID int64) (models.Chat, error) {
			return models.Chat{ID: chatID, Name: "group"}, nil
		},
	}

	svc := newInviteService(storage, nil)
	chat, err := svc.JoinChatByInvite(context.Background(), dto.JoinChatByInviteRequest{
		InviteCode: "abc123",
		UserID:     42,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 1 {
		t.Fatalf("unexpected chat: %+v", chat)
	}
}

func TestInviteService_JoinChatByInvite_NotFound(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		SelectInviteByCodeFn: func(_ context.Context, _ string) (models.InviteLink, error) {
			return models.InviteLink{}, pgx.ErrNoRows
		},
	}

	svc := newInviteService(storage, nil)
	_, err := svc.JoinChatByInvite(context.Background(), dto.JoinChatByInviteRequest{
		InviteCode: "abc",
		UserID:     42,
	})
	if err != domain.ErrInviteNotFound {
		t.Fatalf("expected ErrInviteNotFound, got %v", err)
	}
}

func TestInviteService_JoinChatByInvite_Expired(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		SelectInviteByCodeFn: func(_ context.Context, _ string) (models.InviteLink, error) {
			return models.InviteLink{ID: 1, ChatID: 1, IsActive: true, ExpiresAt: time.Now().Add(-time.Hour)}, nil
		},
	}

	svc := newInviteService(storage, nil)
	_, err := svc.JoinChatByInvite(context.Background(), dto.JoinChatByInviteRequest{
		InviteCode: "abc",
		UserID:     42,
	})
	if err != domain.ErrInviteExpired {
		t.Fatalf("expected ErrInviteExpired, got %v", err)
	}
}

func TestInviteService_JoinChatByInvite_Revoked(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		SelectInviteByCodeFn: func(_ context.Context, _ string) (models.InviteLink, error) {
			return models.InviteLink{ID: 1, ChatID: 1, IsActive: false}, nil
		},
	}

	svc := newInviteService(storage, nil)
	_, err := svc.JoinChatByInvite(context.Background(), dto.JoinChatByInviteRequest{
		InviteCode: "abc",
		UserID:     42,
	})
	if err != domain.ErrInviteRevoked {
		t.Fatalf("expected ErrInviteRevoked, got %v", err)
	}
}

func TestInviteService_JoinChatByInvite_MaxUsesReached(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		SelectInviteByCodeFn: func(_ context.Context, _ string) (models.InviteLink, error) {
			return models.InviteLink{ID: 1, ChatID: 1, IsActive: true, MaxUses: 1, UsedCount: 1}, nil
		},
	}

	svc := newInviteService(storage, nil)
	_, err := svc.JoinChatByInvite(context.Background(), dto.JoinChatByInviteRequest{
		InviteCode: "abc",
		UserID:     42,
	})
	if err != domain.ErrInviteMaxUsesReached {
		t.Fatalf("expected ErrInviteMaxUsesReached, got %v", err)
	}
}

func TestInviteService_RevokeInviteLink(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		DeactivateInviteFn: func(_ context.Context, _, _ int64) error {
			return nil
		},
	}

	svc := newInviteService(storage, nil)
	err := svc.RevokeInviteLink(context.Background(), dto.RevokeInviteLinkRequest{ChatID: 1, InviteID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInviteService_RevokeInviteLink_NotFound(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockInviteStorage{
		DeactivateInviteFn: func(_ context.Context, _, _ int64) error {
			return pgx.ErrNoRows
		},
	}

	svc := newInviteService(storage, nil)
	err := svc.RevokeInviteLink(context.Background(), dto.RevokeInviteLinkRequest{ChatID: 1, InviteID: 1})
	if err != domain.ErrInviteNotFound {
		t.Fatalf("expected ErrInviteNotFound, got %v", err)
	}
}

func TestGenerateInviteCode(t *testing.T) {
	t.Parallel()

	code, err := generateInviteCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != 8 {
		t.Fatalf("expected code length 8, got %d", len(code))
	}

	other, _ := generateInviteCode()
	if code == other {
		t.Fatal("expected generated codes to differ")
	}
}
