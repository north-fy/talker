package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
)

func (s *serverAPI) SendMessage(ctx context.Context, req *messagev1.SendMessageRequest) (*messagev1.Message, error) {
	return nil, nil
}

func (s *serverAPI) GetMessages(ctx context.Context, req *messagev1.GetMessagesRequest) (*messagev1.GetMessagesResponse, error) {
	return nil, nil
}

func (s *serverAPI) EditMessage(ctx context.Context, req *messagev1.EditMessageRequest) (*messagev1.Message, error) {
	return nil, nil
}

func (s *serverAPI) DeleteMessage(ctx context.Context, req *messagev1.DeleteMessageRequest) (*messagev1.Empty, error) {
	return nil, nil
}

func (s *serverAPI) GetMessage(ctx context.Context, req *messagev1.GetMessageRequest) (*messagev1.Message, error) {
	return nil, nil
}
