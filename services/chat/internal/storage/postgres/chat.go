package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

func (s *Storage) InsertChat(ctx context.Context, req dto.CreateChatRequest) (models.Chat, error) {
	if len(req.MemberIDs) == 0 {
		return models.Chat{}, fmt.Errorf("chat must have at least one member")
	}

	createdBy := req.MemberIDs[0]

	query := `
	WITH new_chat AS (
		INSERT INTO chats (name, type, created_by, avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, type, created_by, avatar_url, is_active, created_at, updated_at
	),
	new_members AS (
		INSERT INTO chat_members (chat_id, user_id, role)
		SELECT nc.id, m.uid, CASE WHEN m.uid = $3 THEN 4 ELSE 1 END
		FROM new_chat nc, unnest($5::bigint[]) AS m(uid)
	),
	new_settings AS (
		INSERT INTO chat_settings (chat_id)
		SELECT nc.id FROM new_chat nc
	)
	SELECT nc.id, nc.name, nc.type, nc.created_by, nc.avatar_url, nc.is_active, nc.created_at, nc.updated_at,
	       (SELECT COUNT(*) FROM chat_members WHERE chat_id = nc.id)
	FROM new_chat nc
	`

	var chat models.Chat
	err := s.conn.QueryRow(ctx, query, req.Name, req.Type, createdBy, req.AvatarBase64, req.MemberIDs).Scan(
		&chat.ID,
		&chat.Name,
		&chat.Type,
		&chat.CreatedBy,
		&chat.AvatarURL,
		&chat.IsActive,
		&chat.CreatedAt,
		&chat.UpdatedAt,
		&chat.MembersCount,
	)
	if err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (s *Storage) SelectChat(ctx context.Context, chatID int64) (models.Chat, error) {
	query := `
	SELECT id, name, type, created_by, avatar_url, (
		SELECT COUNT(*) 
		FROM chat_members 
		WHERE chat_id = c.id	
		), is_active, created_at, updated_at
	FROM chats AS c
	WHERE id = $1
	`

	var chat models.Chat
	if err := s.conn.QueryRow(ctx, query, chatID).Scan(
		&chat.ID,
		&chat.Name,
		&chat.Type,
		&chat.CreatedBy,
		&chat.AvatarURL,
		&chat.MembersCount,
		&chat.IsActive,
		&chat.CreatedAt,
		&chat.UpdatedAt); err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (s *Storage) SelectChats(ctx context.Context, filter dto.ChatFilter) (dto.GetChatsResponse, error) {
	baseQuery := `
	SELECT c.id, c.name, c.type, c.created_by, c.avatar_url,
	       (SELECT COUNT(*) FROM chat_members WHERE chat_id = c.id),
	       c.is_active, c.created_at, c.updated_at
	FROM chats c
	WHERE 1=1
	`

	args := []any{}
	argIdx := 1

	if filter.Type != 0 {
		baseQuery += fmt.Sprintf(` AND c.type = $%d`, argIdx)
		args = append(args, filter.Type)
		argIdx++
	}

	if filter.Search != "" {
		baseQuery += fmt.Sprintf(` AND to_tsvector('russian', c.name) @@ plainto_tsquery('russian', $%d)`, argIdx)
		args = append(args, filter.Search)
		argIdx++
	}

	if len(filter.ChatIDs) > 0 {
		baseQuery += fmt.Sprintf(` AND c.id = ANY($%d)`, argIdx)
		args = append(args, filter.ChatIDs)
		argIdx++
	}

	if !filter.IncludeArchived {
		baseQuery += ` AND c.is_active = true`
	}

	baseQuery += ` ORDER BY c.created_at DESC`

	rows, err := s.conn.Query(ctx, baseQuery, args...)
	if err != nil {
		return dto.GetChatsResponse{}, err
	}
	defer rows.Close()

	var chats []*models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.Type,
			&chat.CreatedBy,
			&chat.AvatarURL,
			&chat.MembersCount,
			&chat.IsActive,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		); err != nil {
			return dto.GetChatsResponse{}, err
		}
		chats = append(chats, &chat)
	}

	if err := rows.Err(); err != nil {
		return dto.GetChatsResponse{}, err
	}

	return dto.GetChatsResponse{
		Chats:      chats,
		TotalCount: int64(len(chats)),
	}, nil
}

func (s *Storage) UpdateChat(ctx context.Context, req dto.UpdateChatRequest) (models.Chat, error) {
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf(`name = $%d`, argIdx))
		args = append(args, *req.Name)
		argIdx++
	}

	if req.AvatarBase64 != nil {
		setClauses = append(setClauses, fmt.Sprintf(`avatar_url = $%d`, argIdx))
		args = append(args, *req.AvatarBase64)
		argIdx++
	}

	if len(setClauses) == 0 {
		return s.SelectChat(ctx, req.ChatID)
	}

	setClauses = append(setClauses, `updated_at = NOW()`)

	query := fmt.Sprintf(`
	UPDATE chats SET %s
	WHERE id = $%d
	RETURNING id, name, type, created_by, avatar_url,
	          (SELECT COUNT(*) FROM chat_members WHERE chat_id = chats.id),
	          is_active, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIdx)
	args = append(args, req.ChatID)

	var chat models.Chat
	err := s.conn.QueryRow(ctx, query, args...).Scan(
		&chat.ID,
		&chat.Name,
		&chat.Type,
		&chat.CreatedBy,
		&chat.AvatarURL,
		&chat.MembersCount,
		&chat.IsActive,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (s *Storage) SelectChatsByIDs(ctx context.Context, ids []int64) ([]*models.Chat, error) {
	if len(ids) == 0 {
		return []*models.Chat{}, nil
	}

	query := `
	SELECT c.id, c.name, c.type, c.created_by, c.avatar_url,
	       (SELECT COUNT(*) FROM chat_members WHERE chat_id = c.id),
	       c.is_active, c.created_at, c.updated_at
	FROM chats c
	WHERE c.id = ANY($1)
	ORDER BY c.created_at DESC
	`

	rows, err := s.conn.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.Type,
			&chat.CreatedBy,
			&chat.AvatarURL,
			&chat.MembersCount,
			&chat.IsActive,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		); err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
