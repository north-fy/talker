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
	Attachments []string
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
	MessageID    int64
	ForEveryone  bool
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

// WebSocketEvent represents a WebSocket event
type WebSocketEvent interface {
	EventType() string
}

// NewMessageEvent represents a new message WebSocket event
type NewMessageEvent struct {
	Message models.Message
}

func (e *NewMessageEvent) EventType() string {
	return "new_message"
}

// MessageUpdatedEvent represents a message updated WebSocket event
type MessageUpdatedEvent struct {
	MessageID  int64
	NewContent string
	UpdatedAt  time.Time
}

func (e *MessageUpdatedEvent) EventType() string {
	return "message_updated"
}

// MessageDeletedEvent represents a message deleted WebSocket event
type MessageDeletedEvent struct {
	MessageID   int64
	ForEveryone bool
}

func (e *MessageDeletedEvent) EventType() string {
	return "message_deleted"
}

// ReactionAddedEvent represents a reaction added WebSocket event
type ReactionAddedEvent struct {
	MessageID int64
	UserID    int64
	Reaction  string
	NewCount  int32
}

func (e *ReactionAddedEvent) EventType() string {
	return "reaction_added"
}

// ReactionRemovedEvent represents a reaction removed WebSocket event
type ReactionRemovedEvent struct {
	MessageID int64
	UserID    int64
	Reaction  string
	NewCount  int32
}

func (e *ReactionRemovedEvent) EventType() string {
	return "reaction_removed"
}

// TypingEvent represents a typing WebSocket event
type TypingEvent struct {
	ChatID   int64
	UserID   int64
	Username string
	IsTyping bool
}

func (e *TypingEvent) EventType() string {
	return "typing"
}

// ReadReceiptEvent represents a read receipt WebSocket event
type ReadReceiptEvent struct {
	ChatID             int64
	UserID             int64
	ReadUpToMessageID  int64
	ReadAt             time.Time
}

func (e *ReadReceiptEvent) EventType() string {
	return "read_receipt"
}

