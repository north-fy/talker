package domain

import "errors"

var (
	ErrMessageNotFound      = errors.New("message not found")
	ErrMessageEmptyContent  = errors.New("message content is empty")
	ErrMessageTooLong       = errors.New("message content exceeds max length")
	ErrInvalidStruct        = errors.New("invalid struct")
	ErrCannotEditDeleted    = errors.New("cannot edit deleted message")
	ErrCannotDeleteDeleted  = errors.New("cannot delete already deleted message")
	ErrEditTimeExpired      = errors.New("edit time limit exceeded")
	ErrChatAccessDenied     = errors.New("access to chat denied")
	ErrChatNotFound         = errors.New("chat not found")
	ErrUserNotInChat        = errors.New("user is not a member of this chat")
	ErrReactionExists       = errors.New("reaction already exists by this user")
	ErrReactionNotFound     = errors.New("reaction not found")
	ErrInvalidReaction      = errors.New("invalid reaction emoji")
	ErrReadReceiptInvalid   = errors.New("invalid read receipt: message not in chat")
	ErrWebSocketSubscribe   = errors.New("failed to subscribe to event bus")
	ErrWebSocketPublish     = errors.New("failed to publish event")
	ErrEventBusUnavailable  = errors.New("event bus unavailable")
	ErrNotAuthenticated     = errors.New("user is not authenticated")
	ErrInternalStorage      = errors.New("internal storage error")
)
