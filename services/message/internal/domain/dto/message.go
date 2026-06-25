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
	ChatID      int64
	Content     string
	MessageType MessageType
	ReplyTo     int64
	Attachments []int64
}

type GetMessagesRequest struct {
	ChatID int64
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
	MessageID int64
	Content   string
}

type DeleteMessageRequest struct {
	MessageID   int64
	ForEveryone bool
	ChatID      int64
}

type GetMessageRequest struct {
	MessageID int64
}

type AddReactionRequest struct {
	MessageID int64
	Reaction  string
	UserID    int64
}

type RemoveReactionRequest struct {
	MessageID int64
	Reaction  string
	UserID    int64
}

type SearchMessagesRequest struct {
	ChatID int64
	Query  string
	Limit  int32
	Before int64
}

type SearchMessagesResponse struct {
	Messages []*models.Message
	HasMore  bool
}

type MarkAsReadRequest struct {
	ChatID        int64
	UserID        int64
	UpToMessageID int64
}

type GetUnreadCountRequest struct {
	ChatID int64
	UserID int64
}

type GetUnreadCountResponse struct {
	Count          int32
	LastMessageID  int64
}

type ConnectWebSocketRequest struct {
	ChatID int64
	UserID int64
	Token  string
}

type GetLastMessageRequest struct {
	ChatID int64
}

type DeleteChatMessagesRequest struct {
	ChatID int64
}

type Reaction struct {
	MessageID  int64
	UserID     int64
	Reaction   string
	CreatedAt  time.Time
}

