package models

import "time"

type Notification struct {
	ID int64
	UserID int64
	Type NotificationType
	Title string
	Body string
	Data map[string]string
	IsRead bool
	CreatedAt time.Time
}