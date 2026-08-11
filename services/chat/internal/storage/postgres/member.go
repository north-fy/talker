package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
)

func (s *Storage) AddMember(ctx context.Context, req dto.AddMemberRequest) (dto.MemberDB, error) {
	query := `
	INSERT INTO chat_members (chat_id, user_id, role)
	VALUES ($1, $2, $3)
	RETURNING chat_id, user_id, role, joined_at, last_read_at, unread_count
	`

	member := dto.MemberDB{
		ChatID: req.ChatID,
		UserID: req.UserID,
		Role:   req.Role,
	}

	if err := s.conn.QueryRow(ctx, query, req.ChatID, req.UserID, req.Role).Scan(
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
	RETURNING chat_id, user_id, role, joined_at, last_read_at, unread_count
	`

	member := dto.MemberDB{
		ChatID: req.ChatID,
		UserID: req.UserID,
		Role:   req.Role,
	}

	if err := s.conn.QueryRow(ctx, query, req.Role, req.ChatID, req.UserID).Scan(
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

func (s *Storage) GetMembers(ctx context.Context, req dto.GetMembersRequest) (dto.GetMembersDBResponse, error) {
	query := `
	SELECT chat_id, user_id, role, joined_at, last_read_at, unread_count
	FROM chat_members
	WHERE chat_id = $1
	`

	args := []any{req.ChatID}
	argIdx := 2

	if req.Filter.Role != dto.RoleUnknown {
		query += fmt.Sprintf(` AND role = $%d`, argIdx)
		args = append(args, req.Filter.Role)
		argIdx++
	}

	query += ` ORDER BY joined_at`

	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return dto.GetMembersDBResponse{}, err
	}
	defer rows.Close()

	var (
		members []*dto.MemberDB
		ids     []int64
	)

	for rows.Next() {
		var member dto.MemberDB
		if err := rows.Scan(
			&member.ChatID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&member.LastReadAt,
			&member.UnreadCount); err != nil {
			return dto.GetMembersDBResponse{}, err
		}

		members = append(members, &member)
		ids = append(ids, member.UserID)
	}

	if err := rows.Err(); err != nil {
		return dto.GetMembersDBResponse{}, err
	}

	return dto.GetMembersDBResponse{
		Members: members,
		IDs:     ids,
	}, nil
}

func (s *Storage) GetChatsByUser(ctx context.Context, userID int64) ([]*dto.UserChatDB, error) {
	query := `
	SELECT c.id, c.name, c.type, c.created_by, c.avatar_url, c.is_active, c.created_at, c.updated_at,
	       (SELECT COUNT(*) FROM chat_members WHERE chat_id = c.id),
	       cm.role, cm.joined_at, cm.last_read_at, cm.unread_count
	FROM chats c
	JOIN chat_members cm ON cm.chat_id = c.id
	WHERE cm.user_id = $1
	ORDER BY c.updated_at DESC
	`

	rows, err := s.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userChats []*dto.UserChatDB
	for rows.Next() {
		uc := &dto.UserChatDB{}
		if err := rows.Scan(
			&uc.Chat.ID,
			&uc.Chat.Name,
			&uc.Chat.Type,
			&uc.Chat.CreatedBy,
			&uc.Chat.AvatarURL,
			&uc.Chat.IsActive,
			&uc.Chat.CreatedAt,
			&uc.Chat.UpdatedAt,
			&uc.Chat.MembersCount,
			&uc.Member.Role,
			&uc.Member.JoinedAt,
			&uc.Member.LastReadAt,
			&uc.Member.UnreadCount,
		); err != nil {
			return nil, err
		}
		uc.Member.ChatID = uc.Chat.ID
		uc.Member.UserID = userID
		userChats = append(userChats, uc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userChats, nil
}

func (s *Storage) SelectMemberIDs(ctx context.Context, chatID int64) ([]int64, error) {
	query := `
	SELECT user_id
	FROM chat_members
	WHERE chat_id = $1
	ORDER BY joined_at
	`

	rows, err := s.conn.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Storage) SelectMembersByChatIDs(ctx context.Context, chatIDs []int64) (map[int64][]int64, error) {
	if len(chatIDs) == 0 {
		return map[int64][]int64{}, nil
	}

	query := `
	SELECT chat_id, user_id
	FROM chat_members
	WHERE chat_id = ANY($1)
	`

	rows, err := s.conn.Query(ctx, query, chatIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]int64)
	for rows.Next() {
		var (
			chatID int64
			userID int64
		)
		if err := rows.Scan(&chatID, &userID); err != nil {
			return nil, err
		}
		result[chatID] = append(result[chatID], userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
