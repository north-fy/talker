package clients

import (
	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func NewUserClient(addr string) (userv1.UserServiceClient, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}

	return userv1.NewUserServiceClient(conn), nil
}

func NewMessageClient(addr string) (messagev1.MessageServiceClient, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}

	return messagev1.NewMessageServiceClient(conn), nil
}

func NewChatClient(addr string) (chatv1.ChatServiceClient, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}

	return chatv1.NewChatServiceClient(conn), nil
}
