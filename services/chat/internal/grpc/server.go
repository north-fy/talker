package grpc

import (
	messagev1 "github.com/north-fy/talker/pkg/protos/chat"
	"google.golang.org/grpc"
)

type ChatService interface {
	ChatService
	MemberService
	InviteService
	InternalService
}

type serverAPI struct {
	messagev1.UnimplementedChatServiceServer
	serv ChatService
}

func Register(gRPC *grpc.Server, service ChatService) {
	messagev1.RegisterChatServiceServer(gRPC, &serverAPI{serv: service})
}
