package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/chat/internal/domain/dto"
	"github.com/north-fy/talker/services/chat/internal/domain/models"
)

const inviteSelectColumns = `id, chat_id, code, max_uses, used_count, expires_at, created_at, created_by, is_active`

func scanInvite(row pgx.Row) (models.InviteLink, error) {
	var invite models.InviteLink
	var expires sql.NullTime

	if err := row.Scan(
		&invite.ID,
		&invite.ChatID,
		&invite.Code,
		&invite.MaxUses,
		&invite.UsedCount,
		&expires,
		&invite.CreatedAt,
		&invite.CreatedBy,
		&invite.IsActive,
	); err != nil {
		return models.InviteLink{}, err
	}

	if expires.Valid {
		invite.ExpiresAt = expires.Time
	}

	return invite, nil
}

func (s *Storage) InsertInvite(ctx context.Context, req dto.CreateInviteLinkRequest, createdBy int64, code string) (models.InviteLink, error) {
	query := `
	INSERT INTO invite_links (chat_id, code, max_uses, expires_at, created_by)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING ` + inviteSelectColumns + `
	`

	row := s.conn.QueryRow(ctx, query, req.ChatID, code, req.MaxUses, req.ExpiresAt, createdBy)

	return scanInvite(row)
}

func (s *Storage) SelectInvite(ctx context.Context, id int64) (models.InviteLink, error) {
	query := `
	SELECT ` + inviteSelectColumns + `
	FROM invite_links
	WHERE id = $1
	`

	return scanInvite(s.conn.QueryRow(ctx, query, id))
}

func (s *Storage) SelectInviteByCode(ctx context.Context, code string) (models.InviteLink, error) {
	query := `
	SELECT ` + inviteSelectColumns + `
	FROM invite_links
	WHERE code = $1
	`

	return scanInvite(s.conn.QueryRow(ctx, query, code))
}

func (s *Storage) IncrementUsedCount(ctx context.Context, id int64) (models.InviteLink, error) {
	query := `
	UPDATE invite_links
	SET used_count = used_count + 1
	WHERE id = $1
	  AND is_active = true
	  AND (max_uses = 0 OR used_count < max_uses)
	  AND (expires_at IS NULL OR expires_at > NOW())
	RETURNING ` + inviteSelectColumns + `
	`

	return scanInvite(s.conn.QueryRow(ctx, query, id))
}

func (s *Storage) DeactivateInvite(ctx context.Context, chatID, inviteID int64) error {
	query := `
	UPDATE invite_links
	SET is_active = false
	WHERE id = $1 AND chat_id = $2
	`

	ct, err := s.conn.Exec(ctx, query, inviteID, chatID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
