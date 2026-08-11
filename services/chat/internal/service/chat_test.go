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

func newChatService(storage *mocks.MockChatStorage, cache *mocks.MockCache) *ChatFuncService {
	if cache == nil {
		cache = &mocks.MockCache{}
	}
	return NewChatFuncService(zap.NewNop(), storage, cache)
}

func TestChatService_CreateChat(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockChatStorage{
		InsertChatFn: func(_ context.Context, _ dto.CreateChatRequest) (models.Chat, error) {
			return models.Chat{ID: 1, Name: "group", Type: int32(dto.ChatTypeGroup), CreatedBy: 10, CreatedAt: now, UpdatedAt: now}, nil
		},
	}

	svc := newChatService(storage, nil)
	chat, err := svc.CreateChat(context.Background(), dto.CreateChatRequest{
		Name:      "group",
		Type:      dto.ChatTypeGroup,
		MemberIDs: []int64{10, 20},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 1 || chat.CreatedBy != 10 {
		t.Fatalf("unexpected chat: %+v", chat)
	}
}

func TestChatService_CreateChat_EmptyName(t *testing.T) {
	t.Parallel()

	svc := newChatService(&mocks.MockChatStorage{}, nil)
	_, err := svc.CreateChat(context.Background(), dto.CreateChatRequest{
		Type:      dto.ChatTypeGroup,
		MemberIDs: []int64{10},
	})
	if err != domain.ErrChatNameEmpty {
		t.Fatalf("expected ErrChatNameEmpty, got %v", err)
	}
}

func TestChatService_GetChat_FromCache(t *testing.T) {
	t.Parallel()

	cachedChat := &models.Chat{ID: 7, Name: "cached", MembersCount: 3}
	cache := &mocks.MockCache{
		GetChatFn: func(_ context.Context, chatID int64) (*models.Chat, error) {
			return cachedChat, nil
		},
	}
	storage := &mocks.MockChatStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			t.Fatal("storage should not be called on cache hit")
			return models.Chat{}, nil
		},
	}

	svc := newChatService(storage, cache)
	chat, err := svc.GetChat(context.Background(), dto.GetChatRequest{ChatID: 7, IncludeMembers: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 7 || chat.MembersCount != 3 {
		t.Fatalf("unexpected chat: %+v", chat)
	}
}

func TestChatService_GetChat_WithoutMembersCount(t *testing.T) {
	t.Parallel()

	cache := &mocks.MockCache{
		GetChatFn: func(_ context.Context, chatID int64) (*models.Chat, error) {
			return &models.Chat{ID: 7, MembersCount: 3}, nil
		},
	}
	svc := newChatService(&mocks.MockChatStorage{}, cache)

	chat, err := svc.GetChat(context.Background(), dto.GetChatRequest{ChatID: 7, IncludeMembers: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.MembersCount != 0 {
		t.Fatalf("expected MembersCount to be zeroed, got %d", chat.MembersCount)
	}
}

func TestChatService_GetChat_NotFound(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockChatStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{}, pgx.ErrNoRows
		},
	}

	svc := newChatService(storage, nil)
	_, err := svc.GetChat(context.Background(), dto.GetChatRequest{ChatID: 1})
	if err != domain.ErrChatNotFound {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestChatService_UpdateChat_InvalidName(t *testing.T) {
	t.Parallel()

	tooLong := string(make([]byte, 256))
	svc := newChatService(&mocks.MockChatStorage{}, nil)
	_, err := svc.UpdateChat(context.Background(), dto.UpdateChatRequest{
		ChatID: 1,
		Name:   &tooLong,
	})
	if err != domain.ErrChatNameTooLong {
		t.Fatalf("expected ErrChatNameTooLong, got %v", err)
	}
}
