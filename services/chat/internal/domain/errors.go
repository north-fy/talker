package domain

import "errors"

var (
	// Chat
	ErrChatNotFound         = errors.New("chat not found")
	ErrChatAlreadyExists    = errors.New("chat already exists")
	ErrChatNameEmpty        = errors.New("chat name is empty")
	ErrChatNameTooLong      = errors.New("chat name exceeds max length")
	ErrChatIsArchived       = errors.New("chat is archived")
	ErrCannotDeletePrivateChat = errors.New("cannot delete private chat")

	// Members
	ErrMemberAlreadyInChat  = errors.New("user is already a member of this chat")
	ErrMemberNotInChat      = errors.New("user is not a member of this chat")
	ErrCannotRemoveOwner    = errors.New("cannot remove chat owner")
	ErrCannotRemoveSelf     = errors.New("cannot remove yourself, use leave chat")
	ErrInsufficientRole     = errors.New("insufficient role to perform this action")
	ErrInvalidRole          = errors.New("invalid role")

	// Invite links
	ErrInviteNotFound       = errors.New("invite link not found")
	ErrInviteExpired        = errors.New("invite link has expired")
	ErrInviteMaxUsesReached = errors.New("invite link max uses reached")
	ErrInviteRevoked        = errors.New("invite link has been revoked")

	// Permissions
	ErrAccessDenied         = errors.New("access denied")
	ErrNotAuthenticated     = errors.New("user is not authenticated")
	ErrInvalidPermission    = errors.New("invalid permission type")

	// Storage
	ErrInternalStorage      = errors.New("internal storage error")
	ErrInvalidStruct        = errors.New("invalid struct")
)
