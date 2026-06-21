package grpc

import (
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
)

type MessageService interface {
	MessageFuncService
	ReactionService
	FeatureService
}

type serverAPI struct {
	messagev1.UnimplementedMessageServiceServer
	serv MessageService
	ws serverWebSocket
}
