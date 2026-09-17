package grpc

import (
	"errors"

	"github.com/north-fy/talker/services/notification/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPC(err error) error {
	switch {

	case errors.Is(err, domain.ErrInternalKafka):
			return status.Error(codes.Internal, err.Error())
	case errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrNotificationNotFound):
			return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidData),
		errors.Is(err, domain.ErrInvalidNotificationType),
		errors.Is(err, domain.ErrInvalidPagination):
			return status.Error(codes.InvalidArgument, err.Error())
	}

}
