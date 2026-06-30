package grpc

import (
	"errors"

	"github.com/north-fy/talker/services/message/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPC(err error) error {
	switch {
	case errors.Is(err, domain.ErrMessageNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrChatAccessDenied),
		errors.Is(err, domain.ErrUserNotInChat):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrMessageEmptyContent),
		errors.Is(err, domain.ErrMessageTooLong),
		errors.Is(err, domain.ErrInvalidReaction):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrWebSocketPublish),
		errors.Is(err, domain.ErrEventBusUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
