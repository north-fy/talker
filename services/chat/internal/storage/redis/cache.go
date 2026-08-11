package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
	"github.com/redis/go-redis/v9"
)

const (
	chatKeyPrefix      = "chat:%d"
	userChatsKeyPrefix = "user_chats:%d"
	memberKeyPrefix    = "member:%d:%d"
	cacheTTL           = 5 * time.Minute
)

func chatKey(id int64) string      { return fmt.Sprintf(chatKeyPrefix, id) }
func userChatsKey(id int64) string { return fmt.Sprintf(userChatsKeyPrefix, id) }
func memberKey(chatID, userID int64) string {
	return fmt.Sprintf(memberKeyPrefix, chatID, userID)
}

func (s *Storage) GetChat(ctx context.Context, chatID int64) (*models.Chat, error) {
	data, err := s.client.Get(ctx, chatKey(chatID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var chat models.Chat
	if err := json.Unmarshal(data, &chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (s *Storage) SetChat(ctx context.Context, chat *models.Chat) error {
	data, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, chatKey(chat.ID), data, cacheTTL).Err()
}

func (s *Storage) DeleteChat(ctx context.Context, chatID int64) error {
	return s.client.Del(ctx, chatKey(chatID)).Err()
}

func (s *Storage) GetUserChats(ctx context.Context, userID int64) (*dto.GetUserChatsResponse, error) {
	data, err := s.client.Get(ctx, userChatsKey(userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var resp dto.GetUserChatsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *Storage) SetUserChats(ctx context.Context, userID int64, resp *dto.GetUserChatsResponse) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, userChatsKey(userID), data, cacheTTL).Err()
}

func (s *Storage) DeleteUserChats(ctx context.Context, userID int64) error {
	return s.client.Del(ctx, userChatsKey(userID)).Err()
}

func (s *Storage) GetMember(ctx context.Context, chatID, userID int64) (*dto.MemberDB, error) {
	data, err := s.client.Get(ctx, memberKey(chatID, userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var member dto.MemberDB
	if err := json.Unmarshal(data, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func (s *Storage) SetMember(ctx context.Context, member *dto.MemberDB) error {
	data, err := json.Marshal(member)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, memberKey(member.ChatID, member.UserID), data, cacheTTL).Err()
}

func (s *Storage) DeleteMember(ctx context.Context, chatID, userID int64) error {
	return s.client.Del(ctx, memberKey(chatID, userID)).Err()
}
