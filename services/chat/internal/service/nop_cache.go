package service

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

// nopCache — заглушка, используемая в тестах и когда Redis недоступен.
type nopCache struct{}

func (nopCache) GetChat(context.Context, int64) (*models.Chat, error) {
	return nil, nil
}

func (nopCache) SetChat(context.Context, *models.Chat) error {
	return nil
}

func (nopCache) DeleteChat(context.Context, int64) error {
	return nil
}

func (nopCache) GetUserChats(context.Context, int64) (*dto.GetUserChatsResponse, error) {
	return nil, nil
}

func (nopCache) SetUserChats(context.Context, int64, *dto.GetUserChatsResponse) error {
	return nil
}

func (nopCache) DeleteUserChats(context.Context, int64) error {
	return nil
}

func (nopCache) GetMember(context.Context, int64, int64) (*dto.MemberDB, error) {
	return nil, nil
}

func (nopCache) SetMember(context.Context, *dto.MemberDB) error {
	return nil
}

func (nopCache) DeleteMember(context.Context, int64, int64) error {
	return nil
}
