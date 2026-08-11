package postgres

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

const settingsSelectColumns = `
chat_id, is_private, allow_messages_from_all, 
allow_media, allow_reactions,
message_ttl_seconds, language, 
is_announcement`

type row interface{ Scan(...any) error }

func scanSettings(row row) (int64, dto.ChatSettings, error) {
	var (
		chatID   int64
		settings dto.ChatSettings
	)

	if err := row.Scan(
		&chatID,
		&settings.IsPrivate,
		&settings.AllowMessagesFromAll,
		&settings.AllowMedia,
		&settings.AllowReactions,
		&settings.MessageTTLSeconds,
		&settings.Language,
		&settings.IsAnnouncement,
	); err != nil {
		return 0, dto.ChatSettings{}, err
	}

	return chatID, settings, nil
}

func (s *Storage) SelectSettings(ctx context.Context, chatID int64) (dto.ChatSettings, error) {
	query := `
	SELECT ` + settingsSelectColumns + `
	FROM chat_settings
	WHERE chat_id = $1
	`

	_, settings, err := scanSettings(s.conn.QueryRow(ctx, query, chatID))
	if err != nil {
		return dto.ChatSettings{}, err
	}

	return settings, nil
}

func (s *Storage) SelectSettingsByChatIDs(ctx context.Context, chatIDs []int64) (map[int64]dto.ChatSettings, error) {
	if len(chatIDs) == 0 {
		return map[int64]dto.ChatSettings{}, nil
	}

	query := `
	SELECT ` + settingsSelectColumns + `
	FROM chat_settings
	WHERE chat_id = ANY($1)
	`

	rows, err := s.conn.Query(ctx, query, chatIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]dto.ChatSettings)
	for rows.Next() {
		chatID, settings, err := scanSettings(rows)
		if err != nil {
			return nil, err
		}
		result[chatID] = settings
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
