package postgres

import (
	"context"

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

func (s *Storage) CreateMessage(ctx context.Context,  req dto.SendMessageRequest) (models.Message, error) {
	query := `
	INSERT INTO messages(chat_id, )

	`
}
