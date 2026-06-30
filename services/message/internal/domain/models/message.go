package models

import "time"

type Message struct {
	ID           int64
	ChatID       int64  `validate:"required"`
	SenderID     int64  `validate:"required"`
	Content      string `validate:"required,max=255"`
	MessageType  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsEdited     bool
	IsDeleted    bool
	ReplyTo      int64
	ReplyInfoMsg ReplyInfo
	Attachments  []int64
	Reactions    map[string]int32
}

type ReplyInfo struct {
	MessageID      int64
	SenderName     string
	ContentPreview string
}
