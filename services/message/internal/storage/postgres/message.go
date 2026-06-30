package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
)

/*
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    content TEXT,
    type VARCHAR(50) DEFAULT 'text',
    reply_to UUID,
    attachments TEXT[],
    reactions JSONB DEFAULT '{}',
    is_edited BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
) PARTITION BY RANGE (created_at);

*/

func (s *Storage) CreateMessage(ctx context.Context, senderID int64, req dto.SendMessageRequest) (models.Message, error) {
	// TODO: add type conversion
	query := `
	INSERT INTO messages(chat_id, sender_id, content, reply_to, attachments, reactions)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, type, reactions, is_edited, is_deleted, created_at
	`

	msg := models.Message{
		ChatID:      req.ChatID,
		SenderID:    senderID,
		Content:     req.Content,
		ReplyTo:     req.ReplyTo,
		Attachments: req.Attachments,
		// UpdatedAt в данном случаем мы берем как sql.NullTime.
		// Мы не можем записать nil значение в тип time.Time
		UpdatedAt: sql.NullTime{}.Time,
	}

	reactions := map[string]int32{}

	row := s.conn.QueryRow(ctx, query, req.ChatID, senderID, req.Content, req.ReplyTo, req.Attachments, reactions)
	if err := row.Scan(&msg.ID, &msg.MessageType, &msg.Reactions, &msg.IsEdited, &msg.IsDeleted, &msg.CreatedAt); err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (s *Storage) SelectMessages(ctx context.Context, req dto.GetMessagesRequest) ([]*models.Message, error) {
	query := `
	SELECT id, chat_id, sender_id, content, type, reply_to, attachments, reactions, 
	       is_edited, is_deleted, created_at, updated_at FROM messages
	WHERE chat_id = $1 AND id > $2 AND id < $3
	LIMIT $4
	`

	rows, err := s.conn.Query(ctx, query, req.ChatID, req.After, req.Before, req.Limit)
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

func (s *Storage) UpdateMessage(ctx context.Context, req dto.EditMessageRequest) (models.Message, error) {
	query := `
	UPDATE messages
	SET content=$1, 
	    is_edited=true, 
	    updated_at=$2
	WHERE id=$3
	RETURNING id, chat_id, sender_id, content, type, reply_to, attachments, reactions, 
	          is_edited, is_deleted, created_at, updated_at
	`

	row := s.conn.QueryRow(ctx, query, req.Content, time.Now(), req.MessageID)

	var msg models.Message
	if err := row.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content,
		&msg.MessageType, &msg.ReplyTo, &msg.Attachments, &msg.Reactions,
		&msg.IsEdited, &msg.IsDeleted, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (s *Storage) DeleteMessageForUser(ctx context.Context, id int64) error {
	query := `
	UPDATE messages
	SET is_deleted = true
	WHERE id = $1
	`

	ct, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Storage) DeleteMessage(ctx context.Context, id int64) error {
	query := `
	DELETE FROM messages
	WHERE id = $1
	`

	ct, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Storage) SelectMessage(ctx context.Context, id int64) (models.Message, error) {
	query := `
	SELECT id, chat_id, sender_id, content, type, reply_to, attachments, reactions, 
	       is_edited, is_deleted, created_at, updated_at FROM messages
	WHERE id = $1
	`

	row := s.conn.QueryRow(ctx, query, id)

	var msg models.Message
	if err := row.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content,
		&msg.MessageType, &msg.ReplyTo, &msg.Attachments, &msg.Reactions,
		&msg.IsEdited, &msg.IsDeleted, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
		return models.Message{}, err
	}

	return msg, nil
}
