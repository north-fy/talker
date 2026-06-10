package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
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

func (s *Storage) InsertReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error) {
	query := `
	UPDATE messages 
	SET reactions = jsonb_set(
		reactions,
		ARRAY[$1],
		COALESCE((reactions->>$1)::int, 0)::int + 1,
		true
	)
	WHERE id = $2
	`

	ct, err := s.conn.Exec(ctx, query, req.Reaction, req.MessageID)
	if err != nil {
		return dto.Reaction{}, err
	}

	if ct.RowsAffected() == 0 {
		return dto.Reaction{}, sql.ErrNoRows
	}

	return dto.Reaction{
		MessageID: req.MessageID,
		UserID:    req.UserID,
		Reaction:  req.Reaction,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Storage) DeleteReaction(ctx context.Context, req dto.RemoveReactionRequest) error {
	query := `
	UPDATE messages
	SET reactions = CASE 
		WHEN COALESCE((reactions->>$1)::int, 0) <= 1 THEN reactions - $1
		ELSE reactions || jsonb_build_object($1, (reactions->>$1)::int - 1)
	END
	WHERE id = $2
    `

	ct, err := s.conn.Exec(ctx, query, req.Reaction, req.MessageID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}
