package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/services/chat/internal/domain"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/north-fy/talker/services/chat/internal/mocks"
	"go.uber.org/zap"
)

func newMemberService(storage *mocks.MockMemberStorage, userClient *mocks.MockUserClient, cache *mocks.MockCache) *MemberService {
	if userClient == nil {
		userClient = &mocks.MockUserClient{}
	}
	if cache == nil {
		cache = &mocks.MockCache{}
	}
	return NewMemberService(zap.NewNop(), storage, userClient, cache)
}

func TestMemberService_AddMember(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockMemberStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{ID: 1}, nil
		},
		AddMemberFn: func(_ context.Context, req dto.AddMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: req.ChatID, UserID: req.UserID, Role: dto.RoleMember, JoinedAt: now}, nil
		},
	}
	userClient := &mocks.MockUserClient{
		GetUsersFn: func(_ context.Context, in *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
			return &userv1.GetUsersResponse{Users: []*userv1.User{
				{UserId: in.UserIds[0], FirstName: "John", LastName: "Doe"},
			}}, nil
		},
	}

	svc := newMemberService(storage, userClient, nil)
	member, err := svc.AddMember(context.Background(), dto.AddMemberRequest{
		ChatID:    1,
		UserID:    42,
		Role:      dto.RoleMember,
		InvitedBy: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.UserID != 42 || member.ChatID != 1 || member.FullName != "John Doe" {
		t.Fatalf("unexpected member: %+v", member)
	}
}

func TestMemberService_AddMember_ChatNotFound(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{}, pgx.ErrNoRows
		},
	}

	svc := newMemberService(storage, nil, nil)
	_, err := svc.AddMember(context.Background(), dto.AddMemberRequest{
		ChatID:    1,
		UserID:    42,
		Role:      dto.RoleMember,
		InvitedBy: 10,
	})
	if err != domain.ErrChatNotFound {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestMemberService_AddMember_AlreadyMember(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		SelectChatFn: func(_ context.Context, _ int64) (models.Chat, error) {
			return models.Chat{ID: 1}, nil
		},
		AddMemberFn: func(_ context.Context, _ dto.AddMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{}, &pgconn.PgError{Code: "23505"}
		},
	}

	svc := newMemberService(storage, nil, nil)
	_, err := svc.AddMember(context.Background(), dto.AddMemberRequest{
		ChatID:    1,
		UserID:    42,
		Role:      dto.RoleMember,
		InvitedBy: 10,
	})
	if err != domain.ErrMemberAlreadyInChat {
		t.Fatalf("expected ErrMemberAlreadyInChat, got %v", err)
	}
}

func TestMemberService_AddMember_InvalidRole(t *testing.T) {
	t.Parallel()

	svc := newMemberService(&mocks.MockMemberStorage{}, nil, nil)
	_, err := svc.AddMember(context.Background(), dto.AddMemberRequest{
		ChatID: 1,
		UserID: 42,
		Role:   dto.RoleUnknown,
	})
	if err != domain.ErrInvalidRole {
		t.Fatalf("expected ErrInvalidRole, got %v", err)
	}
}

func TestMemberService_RemoveMember(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleMember}, nil
		},
		RemoveMemberFn: func(_ context.Context, _ dto.RemoveMemberRequest) error {
			return nil
		},
	}

	svc := newMemberService(storage, nil, nil)
	err := svc.RemoveMember(context.Background(), dto.RemoveMemberRequest{ChatID: 1, UserID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMemberService_RemoveMember_Owner(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleOwner}, nil
		},
	}

	svc := newMemberService(storage, nil, nil)
	err := svc.RemoveMember(context.Background(), dto.RemoveMemberRequest{ChatID: 1, UserID: 42})
	if err != domain.ErrCannotRemoveOwner {
		t.Fatalf("expected ErrCannotRemoveOwner, got %v", err)
	}
}

func TestMemberService_RemoveMember_NotInChat(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		RemoveMemberFn: func(_ context.Context, _ dto.RemoveMemberRequest) error {
			return pgx.ErrNoRows
		},
	}

	svc := newMemberService(storage, nil, nil)
	err := svc.RemoveMember(context.Background(), dto.RemoveMemberRequest{ChatID: 1, UserID: 42})
	if err != domain.ErrMemberNotInChat {
		t.Fatalf("expected ErrMemberNotInChat, got %v", err)
	}
}

func TestMemberService_UpdateMemberRole(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleMember}, nil
		},
		UpdateMemberRoleFn: func(_ context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: req.ChatID, UserID: req.UserID, Role: req.Role, JoinedAt: now}, nil
		},
	}
	userClient := &mocks.MockUserClient{
		GetUsersFn: func(_ context.Context, in *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
			return &userv1.GetUsersResponse{Users: []*userv1.User{
				{UserId: in.UserIds[0], FirstName: "Jane", LastName: "Doe"},
			}}, nil
		},
	}

	svc := newMemberService(storage, userClient, nil)
	member, err := svc.UpdateMemberRole(context.Background(), dto.UpdateMemberRoleRequest{
		ChatID: 1,
		UserID: 42,
		Role:   dto.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.Role != "3" {
		t.Fatalf("expected role admin (3), got %q", member.Role)
	}
}

func TestMemberService_UpdateMemberRole_Owner(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleOwner}, nil
		},
	}

	svc := newMemberService(storage, nil, nil)
	_, err := svc.UpdateMemberRole(context.Background(), dto.UpdateMemberRoleRequest{
		ChatID: 1,
		UserID: 42,
		Role:   dto.RoleAdmin,
	})
	if err != domain.ErrCannotModifyOwner {
		t.Fatalf("expected ErrCannotModifyOwner, got %v", err)
	}
}

func TestMemberService_UpdateMemberRole_NotInChat(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{}, pgx.ErrNoRows
		},
	}

	svc := newMemberService(storage, nil, nil)
	_, err := svc.UpdateMemberRole(context.Background(), dto.UpdateMemberRoleRequest{
		ChatID: 1,
		UserID: 42,
		Role:   dto.RoleAdmin,
	})
	if err != domain.ErrMemberNotInChat {
		t.Fatalf("expected ErrMemberNotInChat, got %v", err)
	}
}

func TestMemberService_GetMember_FromCache(t *testing.T) {
	t.Parallel()

	cache := &mocks.MockCache{
		GetMemberFn: func(_ context.Context, _, _ int64) (*dto.MemberDB, error) {
			return &dto.MemberDB{ChatID: 1, UserID: 42, Role: dto.RoleAdmin}, nil
		},
	}
	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			t.Fatal("storage should not be called on cache hit")
			return dto.MemberDB{}, nil
		},
	}
	userClient := &mocks.MockUserClient{
		GetUsersFn: func(_ context.Context, in *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
			return &userv1.GetUsersResponse{Users: []*userv1.User{{UserId: in.UserIds[0]}}}, nil
		},
	}

	svc := newMemberService(storage, userClient, cache)
	member, err := svc.GetMember(context.Background(), dto.GetMemberRequest{ChatID: 1, UserID: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.Role != "3" {
		t.Fatalf("unexpected role: %q", member.Role)
	}
}

func TestMemberService_GetMember_NotInChat(t *testing.T) {
	t.Parallel()

	storage := &mocks.MockMemberStorage{
		GetMemberFn: func(_ context.Context, _ dto.GetMemberRequest) (dto.MemberDB, error) {
			return dto.MemberDB{}, pgx.ErrNoRows
		},
	}

	svc := newMemberService(storage, nil, nil)
	_, err := svc.GetMember(context.Background(), dto.GetMemberRequest{ChatID: 1, UserID: 42})
	if err != domain.ErrMemberNotInChat {
		t.Fatalf("expected ErrMemberNotInChat, got %v", err)
	}
}

func TestMemberService_GetMembers(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockMemberStorage{
		GetMembersFn: func(_ context.Context, _ dto.GetMembersRequest) (dto.GetMembersDBResponse, error) {
			return dto.GetMembersDBResponse{
				Members: []*dto.MemberDB{
					{ChatID: 1, UserID: 42, Role: dto.RoleAdmin, JoinedAt: now},
					{ChatID: 1, UserID: 43, Role: dto.RoleMember, JoinedAt: now},
				},
				IDs: []int64{42, 43},
			}, nil
		},
	}
	userClient := &mocks.MockUserClient{
		GetUsersFn: func(_ context.Context, _ *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
			return &userv1.GetUsersResponse{Users: []*userv1.User{
				{UserId: 42, FirstName: "John", LastName: "Doe"},
				{UserId: 43, FirstName: "Jane", LastName: "Doe"},
			}}, nil
		},
	}

	svc := newMemberService(storage, userClient, nil)
	resp, err := svc.GetMembers(context.Background(), dto.GetMembersRequest{ChatID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 2 {
		t.Fatalf("unexpected total count: %d", resp.TotalCount)
	}
	if resp.Members[0].FullName != "John Doe" {
		t.Fatalf("unexpected member: %+v", resp.Members[0])
	}
}

func TestMemberService_GetMembers_Search(t *testing.T) {
	t.Parallel()

	now := time.Now()
	storage := &mocks.MockMemberStorage{
		GetMembersFn: func(_ context.Context, _ dto.GetMembersRequest) (dto.GetMembersDBResponse, error) {
			return dto.GetMembersDBResponse{
				Members: []*dto.MemberDB{
					{ChatID: 1, UserID: 42, Role: dto.RoleMember, JoinedAt: now},
					{ChatID: 1, UserID: 43, Role: dto.RoleMember, JoinedAt: now},
				},
				IDs: []int64{42, 43},
			}, nil
		},
	}
	userClient := &mocks.MockUserClient{
		GetUsersFn: func(_ context.Context, _ *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
			return &userv1.GetUsersResponse{Users: []*userv1.User{
				{UserId: 42, FirstName: "John", LastName: "Doe"},
				{UserId: 43, FirstName: "Alice", LastName: "Smith"},
			}}, nil
		},
	}

	svc := newMemberService(storage, userClient, nil)
	resp, err := svc.GetMembers(context.Background(), dto.GetMembersRequest{
		ChatID: 1,
		Filter: dto.MemberFilter{Search: "alice"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 1 || resp.Members[0].UserID != 43 {
		t.Fatalf("unexpected members: %+v", resp.Members)
	}
}
