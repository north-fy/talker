package domain

import "errors"

var (
	ErrInternalKafka = errors.New("internal pushment error")
	ErrUserNotFound = errors.New("user not found")
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotificationType = errors.New("invalid notification type")
	ErrInvalidData = errors.New("invalid data")
	ErrInvalidPagination = errors.New("invalid pagination")
)
