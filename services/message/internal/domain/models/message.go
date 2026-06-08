package models

import "time"

type Message struct {
	ID           string
	ChatID       string
	SenderID     string
	Content      string
	MessageType  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsEdited     bool
	IsDeleted    bool
	ReplyTo      string
	ReplyInfoMsg ReplyInfo
	Attachments  []string
	Reactions    map[string]int32
}

type ReplyInfo struct {
	MessageID      string
	SenderName     string
	ContentPreview string
}
