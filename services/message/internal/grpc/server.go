package grpc

import (
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
)

type MessageService interface {
}

type serverAPI struct {
	messagev1.UnimplementedMessageServiceServer
	msg MessageService
}
