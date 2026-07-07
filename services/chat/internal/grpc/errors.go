package grpc

import (
	"errors"

	"github.com/north-fy/talker/services/chat/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPC(err error) error {
	switch {
	// NotFound
	case errors.Is(err, domain.ErrChatNotFound),
		errors.Is(err, domain.ErrMemberNotInChat),
		errors.Is(err, domain.ErrInviteNotFound):
		return status.Error(codes.NotFound, err.Error())

	// InvalidArgument
	case errors.Is(err, domain.ErrChatNameEmpty),
		errors.Is(err, domain.ErrChatNameTooLong),
		errors.Is(err, domain.ErrInvalidRole),
		errors.Is(err, domain.ErrInvalidPermission),
		errors.Is(err, domain.ErrInvalidStruct):
		return status.Error(codes.InvalidArgument, err.Error())

	// PermissionDenied
	case errors.Is(err, domain.ErrAccessDenied),
		errors.Is(err, domain.ErrNotAuthenticated),
		errors.Is(err, domain.ErrInsufficientRole),
		errors.Is(err, domain.ErrCannotRemoveOwner),
		errors.Is(err, domain.ErrCannotDeletePrivateChat):
		return status.Error(codes.PermissionDenied, err.Error())

	// AlreadyExists
	case errors.Is(err, domain.ErrMemberAlreadyInChat),
		errors.Is(err, domain.ErrChatAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	// FailedPrecondition
	case errors.Is(err, domain.ErrInviteExpired),
		errors.Is(err, domain.ErrInviteMaxUsesReached),
		errors.Is(err, domain.ErrInviteRevoked),
		errors.Is(err, domain.ErrChatIsArchived):
		return status.Error(codes.FailedPrecondition, err.Error())

	// Internal
	case errors.Is(err, domain.ErrInternalStorage):
		return status.Error(codes.Internal, err.Error())

	default:
		return status.Error(codes.Internal, "internal error")
	}
}
