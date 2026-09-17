package grpc

import (
	"context"

	notificationv1 "github.com/north-fy/talker/pkg/protos/notification"
	"github.com/north-fy/talker/services/notification/internal/domain/dto"
	"github.com/north-fy/talker/services/notification/pkg/convert"
)

type NotificationService interface {
	GetNotification(ctx context.Context, req dto.GetNotificationRequest) (dto.GetNotificationResponse, error)
	GetNotifications(ctx context.Context, req dto.GetNotificationsRequest) (dto.GetNotificationsResponse, error)
	GetUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error)
}

func (s *serverAPI) GetNotification(ctx context.Context, req *notificationv1.GetNotificationRequest) (*notificationv1.Notification, error) {
	notifyReq := dto.GetNotificationRequest{
		ID: req.GetNotificationId(),
	}

	resp, err := s.serv.GetNotification(ctx, notifyReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return convert.ConvertNotificationToProto(&resp.Notification), nil
}

func (s *serverAPI) GetNotifications(ctx context.Context, req *notificationv1.GetNotificationsRequest) (*notificationv1.GetNotificationsResponse, error) {
	notifyReq := dto.GetNotificationsRequest{
		UserID: req.GetUserId(),
		Limit: req.GetLimit(),
		Before: req.GetBefore(),
		OnlyUnread: req.GetOnlyUnread(),
	}

	resp, err := s.serv.GetNotifications(ctx, notifyReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	notifyes := make([]*notificationv1.Notification, len(resp.Notifications))
	for _, notify := range resp.Notifications {
		notifyes = append(notifyes, convert.ConvertNotificationToProto(&notify))
	}

	return &notificationv1.GetNotificationsResponse{
		Notifications: notifyes,
		HasMore: resp.HasMore,
		UnreadCount: resp.UnreadCount,
	}, nil
}

func (s *serverAPI) GetUnreadCount(ctx context.Context, req *notificationv1.GetUnreadCountRequest) (*notificationv1.GetUnreadCountResponse, error) {
	notifyReq := dto.GetUnreadCountRequest{
		UserID: req.GetUserId(),
	}

	resp, err := s.serv.GetUnreadCount(ctx, notifyReq)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &notificationv1.GetUnreadCountResponse{
		Count: resp.Count,
	}, nil
}