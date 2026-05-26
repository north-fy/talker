package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
)

func (s *serverAPI) SearchMessages(ctx context.Context, req *messagev1.SearchMessagesRequest) (*messagev1.SearchMessagesResponse, error) {
	return nil, nil
}

func (s *serverAPI) MarkAsRead(ctx context.Context, req *messagev1.MarkAsReadRequest) (*messagev1.Empty, error) {
	return nil, nil
}

func (s *serverAPI) GetUnreadCount(ctx context.Context, req *messagev1.GetUnreadCountRequest) (*messagev1.GetUnreadCountResponse, error) {
	return nil, nil
}

func (s *serverAPI) ConnectWebSocket(req *messagev1.ConnectWebSocketRequest, stream messagev1.MessageService_ConnectWebSocketServer) error {
	return nil
}

func (s *serverAPI) GetLastMessage(ctx context.Context, req *messagev1.GetLastMessageRequest) (*messagev1.Message, error) {
	return nil, nil
}

func (s *serverAPI) DeleteChatMessages(ctx context.Context, req *messagev1.DeleteChatMessagesRequest) (*messagev1.Empty, error) {
	return nil, nil
}
