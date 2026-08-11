package models

import "time"

type Chat struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Type         int32     `json:"type"`
	CreatedBy    int64     `json:"created_by"`
	AvatarURL    string    `json:"avatar_url"`
	MembersCount int32     `json:"members_count"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Member struct {
	UserID      int64     `json:"user_id"`
	ChatID      int64     `json:"chat_id"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
	LastReadAt  time.Time `json:"last_read_at"`
	UnreadCount int64     `json:"unread_count"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	AvatarURL   string    `json:"avatar_url"`
}

type InviteLink struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Code      string    `json:"code"`
	URL       string    `json:"url"`
	MaxUses   int32     `json:"max_uses"`
	UsedCount int32     `json:"used_count"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int64     `json:"created_by"`
	IsActive  bool      `json:"is_active"`
}
