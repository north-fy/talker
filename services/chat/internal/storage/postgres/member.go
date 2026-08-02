package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

/*

CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'group',
    created_by UUID NOT NULL,
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE chat_members (
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW(),
    last_read_at TIMESTAMP,
    unread_count BIGINT DEFAULT 0,
    PRIMARY KEY (chat_id, user_id)
);

*/

func (s *Storage) AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error) {
	query := `
	INSERT INTO chat_members (chat_id, user_id, role)
	VALUES ($1, $2, $3)
	RETURNING joined_at, last_read_at, unread_count
	`

	member := dto.MemberDB{
		ChatID: req.ChatID,
		UserID: req.UserID,
		Role:   req.Role,
	}

	if err := s.conn.QueryRow(ctx, query, req.ChatID, req.UserID, req.Role).Scan(
		&member.JoinedAt,
		&member.LastReadAt,
		&member.UnreadCount); err != nil {
		return dto.MemberDB{}, err
	}

	return member, nil
}

func (s *Storage) RemoveMember(ctx context.Context, req dto.RemoveMemberRequest) error {
	query := `
	DELETE FROM chat_members
	WHERE chat_id = $1 AND user_id = $2
	`

	stmt, err := s.conn.Exec(ctx, query, req.ChatID, req.UserID)
	if err != nil {
		return err
	}

	if stmt.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *Storage) UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.MemberDB, error) {
	query := `
	UPDATE chat_members
	SET role = $1
	WHERE chat_id = $2 AND user_id = $3
	RETURNING joined_at, last_read_at, unread_count
	`

	member := dto.MemberDB{
		ChatID: req.ChatID,
		UserID: req.UserID,
		Role:   req.Role,
	}

	if err := s.conn.QueryRow(ctx, query, req.Role, req.ChatID, req.UserID).Scan(
		&member.JoinedAt,
		&member.LastReadAt,
		&member.UnreadCount); err != nil {
		return dto.MemberDB{}, err
	}

	return member, nil
}

func (s *Storage) GetMember(ctx context.Context, req dto.GetMemberRequest) (dto.MemberDB, error) {
	query := `
	SELECT chat_id, user_id, role, joined_at, last_read_at, unread_count
	FROM chat_members
	WHERE chat_id = $1 AND user_id = $2
	`

	var member dto.MemberDB
	if err := s.conn.QueryRow(ctx, query, req.ChatID, req.UserID).Scan(
		&member.ChatID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
		&member.LastReadAt,
		&member.UnreadCount); err != nil {
		return dto.MemberDB{}, err
	}

	return member, nil
}

func (s *Storage) GetMembers(ctx context.Context, req dto.GetMembersRequest) ([]*dto.MemberDB, error) {
	query := `
	SELECT chat_id, user_id, role, joined_at, last_read_at, unread_count
	FROM chat_members
	WHERE chat_id = $1 AND user_id = $2
	`

	rows, err := s.conn.Query(ctx, query, req.ChatID)
	if err != nil {
		return nil, err
	}

	var members []*dto.MemberDB
	for rows.Next() {
		var member dto.MemberDB
		if err := rows.Scan(
			&member.ChatID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&member.LastReadAt,
			&member.UnreadCount); err != nil {
			return nil, err
		}

		members = append(members, &member)
	}

	return members, nil
}
