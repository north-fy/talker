package dto

import (
	"time"

	"github.com/north-fy/talker/services/message/internal/domain/models"
)

// MessageType represents the type of message
type MessageType int32

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeText
	MessageTypeImage
	MessageTypeFile
	MessageTypeSticker
	MessageTypeSystem
)

type SendMessageRequest struct {
	ChatID      int64  `validate:"required"`
	Content     string `validate:"required,max=255"`
	MessageType MessageType
	ReplyTo     int64
	Attachments []int64
}

type GetMessagesRequest struct {
	ChatID int64  `validate:"required"`
	Limit  int32
	Before int64
	After  int64
}

type GetMessagesResponse struct {
	Messages   []*models.Message
	HasMore    bool
	TotalCount int32
}

type EditMessageRequest struct {
	MessageID int64 `validate:"required"`
	Content   string`validate:"required,max=255"`
}

type DeleteMessageRequest struct {
	MessageID   int64`validate:"required"`
	ForEveryone bool
	ChatID      int64`validate:"required"`
}

type GetMessageRequest struct {
	MessageID int64 `validate:"required"`
}

type AddReactionRequest struct {
	MessageID int64  `validate:"required"`
	Reaction  string `validate:"required"`
	UserID    int64  `validate:"required"`
}

type RemoveReactionRequest struct {
	MessageID int64 `validate:"required"`
	Reaction  string`validate:"required"`
	UserID    int64 `validate:"required"`
}

type SearchMessagesRequest struct {
	ChatID int64 `validate:"required"`
	Query  string
	Limit  int32
	Before int64
}

type SearchMessagesResponse struct {
	Messages []*models.Message
	HasMore  bool
}

type MarkAsReadRequest struct {
	ChatID        int64 `validate:"required"`
	UserID        int64 `validate:"required"`
	UpToMessageID int64
}

type GetUnreadCountRequest struct {
	ChatID int64 `validate:"required"`
	UserID int64 `validate:"required"`
}

type GetUnreadCountResponse struct {
	Count          int32
	LastMessageID  int64
}

type ConnectWebSocketRequest struct {
	ChatID int64 `validate:"required"`
	UserID int64 `validate:"required"`
	Token  string
}

type GetLastMessageRequest struct {
	ChatID int64 `validate:"required"`
}

type DeleteChatMessagesRequest struct {
	ChatID int64 `validate:"required"`
}

type Reaction struct {
	MessageID  int64 `validate:"required"`
	UserID     int64 `validate:"required"`
	Reaction   string `validate:"required"`
	CreatedAt  time.Time
}

