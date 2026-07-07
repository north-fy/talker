package models

import "time"

type Chat struct {
	ID          string
	Name        string
	Type        string
	CreatedBy   string
	AvatarURL   string
	MembersCount int32
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Member struct {
	UserID     string
	ChatID     string
	Role       string
	JoinedAt   time.Time
	LastReadAt time.Time
	UnreadCount int64
	Username   string
	FullName   string
	AvatarURL  string
}

type InviteLink struct {
	ID        string
	ChatID    string
	Code      string
	URL       string
	MaxUses   int32
	UsedCount int32
	ExpiresAt time.Time
	CreatedAt time.Time
	CreatedBy string
	IsActive  bool
}

type UserChat struct {
	Chat        Chat
	MemberInfo  Member
	UnreadCount int64
}
