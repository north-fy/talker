package convert

import (
	notificationv1 "github.com/north-fy/talker/pkg/protos/notification"
	"github.com/north-fy/talker/services/notification/internal/domain/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertNotificationToProto(dtoNotify *models.Notification) *notificationv1.Notification {
	return &notificationv1.Notification{
		Id: dtoNotify.ID,
		UserId: dtoNotify.UserID,
		Type: notificationv1.NotificationType(dtoNotify.Type),
		Title: dtoNotify.Title,
		Body: dtoNotify.Body,
		Data: dtoNotify.Data,
		IsRead: dtoNotify.IsRead,
		CreatedAt: timestamppb.New(dtoNotify.CreatedAt),
	}
}