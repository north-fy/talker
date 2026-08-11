package grpc

import (
	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"google.golang.org/grpc"
)

type ChatService interface {
	ChatFuncService
	MemberService
	FeatureService
	InternalService
	InviteService
}

type serverAPI struct {
	chatv1.UnimplementedChatServiceServer
	serv ChatService
}

func Register(gRPC *grpc.Server, service ChatService) {
	chatv1.RegisterChatServiceServer(gRPC, &serverAPI{serv: service})
}
