package models

type NotificationType int32

const (
	NotificationTypeUnknown NotificationType = iota
	NotificationTypeMessage
	NotificationTypeReaction
	NotificationTypeMention
	NotificationTypeChatInvite
	NotificationTypeSystem
)
