package grpc

import (
	"errors"

	"github.com/north-fy/talker/services/message/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPC(err error) error {
	switch {
	// NotFound
	case errors.Is(err, domain.ErrMessageNotFound),
		errors.Is(err, domain.ErrChatNotFound),
		errors.Is(err, domain.ErrReactionNotFound):
		return status.Error(codes.NotFound, err.Error())

	// InvalidArgument
	case errors.Is(err, domain.ErrMessageEmptyContent),
		errors.Is(err, domain.ErrMessageTooLong),
		errors.Is(err, domain.ErrInvalidStruct),
		errors.Is(err, domain.ErrInvalidReaction),
		errors.Is(err, domain.ErrReadReceiptInvalid):
		return status.Error(codes.InvalidArgument, err.Error())

	// PermissionDenied
	case errors.Is(err, domain.ErrChatAccessDenied),
		errors.Is(err, domain.ErrUserNotInChat),
		errors.Is(err, domain.ErrNotAuthenticated):
		return status.Error(codes.PermissionDenied, err.Error())

	// AlreadyExists
	case errors.Is(err, domain.ErrReactionExists):
		return status.Error(codes.AlreadyExists, err.Error())

	// FailedPrecondition
	case errors.Is(err, domain.ErrCannotEditDeleted),
		errors.Is(err, domain.ErrCannotDeleteDeleted),
		errors.Is(err, domain.ErrEditTimeExpired):
		return status.Error(codes.FailedPrecondition, err.Error())

	// Unavailable
	case errors.Is(err, domain.ErrWebSocketSubscribe),
		errors.Is(err, domain.ErrWebSocketPublish),
		errors.Is(err, domain.ErrEventBusUnavailable):
		return status.Error(codes.Unavailable, err.Error())

	// Internal
	case errors.Is(err, domain.ErrInternalStorage):
		return status.Error(codes.Internal, err.Error())

	default:
		return status.Error(codes.Internal, "internal error")
	}
}
