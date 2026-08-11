package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/internal/mocks"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newFeatureService(storage *mocks.MockFeatStorage, messageClient *mocks.MockMessageClient, cache *mocks.MockCache) *FeatureService {
	if messageClient == nil {
		messageClient = &mocks.MockMessageClient{}
	}
	if cache == nil {
		cache = &mocks.MockCache{}
	}
	return NewFeatureService(zap.NewNop(), storage, messageClient, cache)
}

func TestFeatureService_IsMember_Member(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockFeatStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleMember}, nil
		},
	}

	svc := newFeatureService(storage, nil, nil)
	resp, err := svc.IsMember(context.Background(), dto.IsMemberRequest{ChatID: 1, UserID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsMember || resp.Role != dto.RoleMember {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestFeatureService_IsMember_NotMember(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockFeatStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{}, pgx.ErrNoRows
		},
	}

	svc := newFeatureService(storage, nil, nil)
	resp, err := svc.IsMember(context.Background(), dto.IsMemberRequest{ChatID: 1, UserID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.IsMember {
		t.Fatal("expected not a member")
	}
}

func TestFeatureService_IsMember_FromCache(t *testing.T) {
	t.Parallel()

	cache := &mocks.MockCache{
		GetMemberFn: func(_ context.Context, _, _ int64) (*dto.MemberDB, error) {
			return &dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleAdmin}, nil
		},
	}
	storage := &mocks.MockFeatStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			t.Fatal("storage should not be called on cache hit")
			return dto.MemberDB{}, nil
		},
	}

	svc := newFeatureService(storage, nil, cache)
	resp, err := svc.IsMember(context.Background(), dto.IsMemberRequest{ChatID: 1, UserID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Role != dto.RoleAdmin {
		t.Fatalf("unexpected role: %v", resp.Role)
	}
}

func TestFeatureService_GetUserChats(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockFeatStorage{
		GetChatsByUserFn: func(_ context.Context, userID int64) ([]*dto.UserChatDB, error) {
			return []*dto.UserChatDB{
				{
					Chat:   models.Chat{ID: 1, Name: "one"},
					Member: dto.MemberDB{ChatID: 1, UserID: userID, Role: dto.RoleOwner, UnreadCount: 5},
				},
			}, nil
		},
	}
	messageClient := &mocks.MockMessageClient{
		GetLastMessageFn: func(_ context.Context, in *messagev1.GetLastMessageRequest) (*messagev1.Message, error) {
			return &messagev1.Message{
				Id:        100,
				ChatId:    in.ChatId,
				SenderId:  9,
				Content:   "hello",
				CreatedAt: timestamppb.New(now),
			}, nil
		},
	}

	svc := newFeatureService(storage, messageClient, nil)
	resp, err := svc.GetUserChats(context.Background(), dto.GetUserChatsRequest{UserID: 42, IncludeLastMessage: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 1 {
		t.Fatalf("unexpected total count: %d", resp.TotalCount)
	}
	item := resp.UserChats[0]
	if item.Chat.ID != 1 || item.UnreadCount != 5 {
		t.Fatalf("unexpected user chat: %+v", item)
	}
	if item.LastMessage == nil || item.LastMessage.ID != 100 || item.LastMessage.Content != "hello" {
		t.Fatalf("unexpected last message: %+v", item.LastMessage)
	}
	if item.MemberInfo.Role != "4" {
		t.Fatalf("unexpected role: %q", item.MemberInfo.Role)
	}
}

func TestFeatureService_GetUserChats_FromCache(t *testing.T) {
	t.Parallel()

	cached := &dto.GetUserChatsResponse{
		UserChats:  []*dto.UserChatResponse{{Chat: &models.Chat{ID: 1}}},
		TotalCount: 1,
	}
	cache := &mocks.MockCache{
		GetUserChatsFn: func(_ context.Context, _ int64) (*dto.GetUserChatsResponse, error) {
			return cached, nil
		},
	}
	storage := &mocks.MockFeatStorage{
		GetChatsByUserFn: func(_ context.Context, _ int64) ([]*dto.UserChatDB, error) {
			t.Fatal("storage should not be called on cache hit")
			return nil, nil
		},
	}

	svc := newFeatureService(storage, nil, cache)
	resp, err := svc.GetUserChats(context.Background(), dto.GetUserChatsRequest{UserID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
