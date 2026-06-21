package dto

import (
	"time"

	"github.com/north-fy/talker/services/message/internal/domain/models"
)

type EventRequest struct {
	ChatID   int64
	UserID   int64
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

