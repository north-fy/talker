package models

import "time"

type Chat struct {
	ID          int64
	Name        string
	Type        int32
	CreatedBy   int64
	AvatarURL   string
	MembersCount int32
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Member struct {
	UserID     int64
	ChatID     int64
	Role       string
	JoinedAt   time.Time
	LastReadAt time.Time
	UnreadCount int64
	Username   string
	FullName   string
	AvatarURL  string
}

type InviteLink struct {
	ID        int64
	ChatID    int64
	Code      string
	URL       string
	MaxUses   int32
	UsedCount int32
	ExpiresAt time.Time
	CreatedAt time.Time
	CreatedBy string
	IsActive  bool
}
