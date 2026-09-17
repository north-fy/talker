package grpc

import (
	notificationv1 "github.com/north-fy/talker/pkg/protos/notification"
	"google.golang.org/grpc"
)

type serverAPI struct {
	notificationv1.UnimplementedNotificationServiceServer
	serv NotificationService // implement
}

func Register(gRPC *grpc.Server, service NotificationService) {
	notificationv1.RegisterNotificationServiceServer(gRPC, &serverAPI{serv: service})
}
