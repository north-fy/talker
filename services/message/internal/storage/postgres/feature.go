package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
)

/*
-- read_receipts
CREATE TABLE read_receipts (
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    last_read_message_id UUID,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);
*/

func (s *Storage) SearchMessages(ctx context.Context, req dto.SearchMessagesRequest) ([]*models.Message, error) {
	query := `
	SELECT id, chat_id, sender_id, content, type, reply_to, attachments, reactions, 
	       is_edited, is_deleted, created_at, updated_at FROM messages
	WHERE chat_id = $1 
	  AND content @@ plainto_tsquery('russian', $2)
	  AND id > $3
	LIMIT $4
	`

	rows, err := s.conn.Query(ctx, query, req.ChatID, req.Query, req.Before, req.Limit)
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, 0, req.Limit)
	for rows.Next() {
		msg := models.Message{}
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.MessageType, &msg.ReplyTo,
			&msg.Attachments, &msg.Reactions, &msg.IsEdited, &msg.IsDeleted, &msg.CreatedAt,
			&msg.UpdatedAt); err != nil {
			return nil, err
		}

		messages = append(messages, &msg)
	}

	return messages, nil
}

func (s *Storage) SetAsRead(ctx context.Context, req dto.MarkAsReadRequest) error {
	query := `
	INSERT INTO read_receipts (chat_id, user_id, last_read_message_id)
	VALUES ($1, $2, $3)
	`

	ct, err := s.conn.Exec(ctx, query, req.ChatID, req.UserID, req.UpToMessageID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Storage) SelectUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error) {
	query := `
	SELECT count(*), MAX(last_read_message_id) FROM read_receipts
	WHERE chat_id = $1 AND user_id = $2
	`

	resp := dto.GetUnreadCountResponse{}
	if err := s.conn.QueryRow(ctx, query, req.ChatID, req.UserID).Scan(&resp.Count, &resp.LastMessageID); err != nil {
		return dto.GetUnreadCountResponse{}, err
	}

	return resp, nil
}

func (s *Storage) SelectLastMessage(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error) {
	query := `
	SELECT id, chat_id, sender_id, content, type, reply_to, attachments, reactions, 
	       is_edited, is_deleted, created_at, updated_at 
	FROM messages
	WHERE chat_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	msg := models.Message{}

	row := s.conn.QueryRow(ctx, query, req.ChatID)
	if err := row.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.MessageType, &msg.ReplyTo,
		&msg.Attachments, &msg.Reactions, &msg.IsEdited, &msg.IsDeleted, &msg.CreatedAt,
		&msg.UpdatedAt); err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (s *Storage) DeleteChatMessages(ctx context.Context, req dto.DeleteChatMessagesRequest) error {
	query := `
	DELETE FROM messages
	WHERE chat_id = $1
	`

	ct, err := s.conn.Exec(ctx, query, req.ChatID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
