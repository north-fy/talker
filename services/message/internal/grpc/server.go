package grpc

import (
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"google.golang.org/grpc"
)

type MessageService interface {
	MessageFuncService
	ReactionService
	FeatureService
}

type serverAPI struct {
	messagev1.UnimplementedMessageServiceServer
	serv MessageService
	ws   ServerWebSocket
}

func Register(gRPC *grpc.Server, service MessageService, ws ServerWebSocket) {
	messagev1.RegisterMessageServiceServer(gRPC, &serverAPI{serv: service, ws: ws})
}
