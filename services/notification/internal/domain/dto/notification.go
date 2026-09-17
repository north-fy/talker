package dto

import "github.com/north-fy/talker/services/notification/internal/domain/models"

// Get one notification
type GetNotificationRequest struct {
    ID int64
}

type GetNotificationResponse struct {
    models.Notification
}

// Get notifications list
type GetNotificationsRequest struct {
	UserID int64
	Limit int32
	Before int64
	OnlyUnread bool
}

type GetNotificationsResponse struct {
	Notifications []models.Notification
	HasMore bool
	UnreadCount int32
}

// Get count notifications
type GetUnreadCountRequest struct {
	UserID int64
}

type GetUnreadCountResponse struct {
	Count int32
}