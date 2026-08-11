package service

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

// Cache предоставляет cache-aside хранилище для горячих данных чатов.
// Методы Get* возвращают nil при промахе.
type Cache interface {
	GetChat(ctx context.Context, chatID int64) (*models.Chat, error)
	SetChat(ctx context.Context, chat *models.Chat) error
	DeleteChat(ctx context.Context, chatID int64) error

	GetUserChats(ctx context.Context, userID int64) (*dto.GetUserChatsResponse, error)
	SetUserChats(ctx context.Context, userID int64, resp *dto.GetUserChatsResponse) error
	DeleteUserChats(ctx context.Context, userID int64) error

	GetMember(ctx context.Context, chatID, userID int64) (*dto.MemberDB, error)
	SetMember(ctx context.Context, member *dto.MemberDB) error
	DeleteMember(ctx context.Context, chatID, userID int64) error
}
