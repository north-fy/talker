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
	ChatID      string
	Content     string
	MessageType MessageType
	ReplyTo     string
	Attachments []string
}

type GetMessagesRequest struct {
	ChatID string
	Limit  int32
	Before string
	After  string
}

type GetMessagesResponse struct {
	Messages   []*models.Message
	HasMore    bool
	TotalCount int32
}

type EditMessageRequest struct {
	MessageID string
	Content   string
}

type DeleteMessageRequest struct {
	MessageID    string
	ForEveryone  bool
}

type GetMessageRequest struct {
	MessageID string
}

type AddReactionRequest struct {
	MessageID string
	Reaction  string
	UserID    string
}

type RemoveReactionRequest struct {
	MessageID string
	Reaction  string
	UserID    string
}

type SearchMessagesRequest struct {
	ChatID string
	Query  string
	Limit  int32
	Before string
}

type SearchMessagesResponse struct {
	Messages []*models.Message
	HasMore  bool
}

type MarkAsReadRequest struct {
	ChatID        string
	UserID        string
	UpToMessageID string
}

type GetUnreadCountRequest struct {
	ChatID string
	UserID string
}

type GetUnreadCountResponse struct {
	Count          int32
	LastMessageID  string
}

type ConnectWebSocketRequest struct {
	ChatID string
	UserID string
	Token  string
}

type GetLastMessageRequest struct {
	ChatID string
}

type DeleteChatMessagesRequest struct {
	ChatID string
}

type Reaction struct {
	MessageID  string
	UserID     string
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
	MessageID  string
	NewContent string
	UpdatedAt  time.Time
}

func (e *MessageUpdatedEvent) EventType() string {
	return "message_updated"
}

// MessageDeletedEvent represents a message deleted WebSocket event
type MessageDeletedEvent struct {
	MessageID   string
	ForEveryone bool
}

func (e *MessageDeletedEvent) EventType() string {
	return "message_deleted"
}

// ReactionAddedEvent represents a reaction added WebSocket event
type ReactionAddedEvent struct {
	MessageID string
	UserID    string
	Reaction  string
	NewCount  int32
}

func (e *ReactionAddedEvent) EventType() string {
	return "reaction_added"
}

// ReactionRemovedEvent represents a reaction removed WebSocket event
type ReactionRemovedEvent struct {
	MessageID string
	UserID    string
	Reaction  string
	NewCount  int32
}

func (e *ReactionRemovedEvent) EventType() string {
	return "reaction_removed"
}

// TypingEvent represents a typing WebSocket event
type TypingEvent struct {
	ChatID   string
	UserID   string
	Username string
	IsTyping bool
}

func (e *TypingEvent) EventType() string {
	return "typing"
}

// ReadReceiptEvent represents a read receipt WebSocket event
type ReadReceiptEvent struct {
	ChatID             string
	UserID             string
	ReadUpToMessageID  string
	ReadAt             time.Time
}

func (e *ReadReceiptEvent) EventType() string {
	return "read_receipt"
}

