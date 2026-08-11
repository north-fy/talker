package mocks

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"google.golang.org/grpc"
)

// MockUserClient — мок userv1.UserServiceClient.
type MockUserClient struct {
	GetUsersFn func(ctx context.Context, in *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error)
}

func (m *MockUserClient) Register(ctx context.Context, in *userv1.RegisterRequest, _ ...grpc.CallOption) (*userv1.RegisterResponse, error) {
	return nil, nil
}

func (m *MockUserClient) Login(ctx context.Context, in *userv1.LoginRequest, _ ...grpc.CallOption) (*userv1.LoginResponse, error) {
	return nil, nil
}

func (m *MockUserClient) GetMe(ctx context.Context, in *userv1.GetMeRequest, _ ...grpc.CallOption) (*userv1.GetMeResponse, error) {
	return nil, nil
}

func (m *MockUserClient) GetUsers(ctx context.Context, in *userv1.GetUsersRequest, _ ...grpc.CallOption) (*userv1.GetUsersResponse, error) {
	if m.GetUsersFn != nil {
		return m.GetUsersFn(ctx, in)
	}
	return &userv1.GetUsersResponse{}, nil
}

func (m *MockUserClient) ValidateToken(ctx context.Context, in *userv1.ValidateTokenRequest, _ ...grpc.CallOption) (*userv1.ValidateTokenResponse, error) {
	return nil, nil
}

// MockMessageClient — мок messagev1.MessageServiceClient.
type MockMessageClient struct {
	GetLastMessageFn func(ctx context.Context, in *messagev1.GetLastMessageRequest) (*messagev1.Message, error)
}

func (m *MockMessageClient) SendMessage(ctx context.Context, in *messagev1.SendMessageRequest, _ ...grpc.CallOption) (*messagev1.Message, error) {
	return nil, nil
}

func (m *MockMessageClient) GetMessages(ctx context.Context, in *messagev1.GetMessagesRequest, _ ...grpc.CallOption) (*messagev1.GetMessagesResponse, error) {
	return nil, nil
}

func (m *MockMessageClient) EditMessage(ctx context.Context, in *messagev1.EditMessageRequest, _ ...grpc.CallOption) (*messagev1.Message, error) {
	return nil, nil
}

func (m *MockMessageClient) DeleteMessage(ctx context.Context, in *messagev1.DeleteMessageRequest, _ ...grpc.CallOption) (*messagev1.Empty, error) {
	return nil, nil
}

func (m *MockMessageClient) GetMessage(ctx context.Context, in *messagev1.GetMessageRequest, _ ...grpc.CallOption) (*messagev1.Message, error) {
	return nil, nil
}

func (m *MockMessageClient) AddReaction(ctx context.Context, in *messagev1.AddReactionRequest, _ ...grpc.CallOption) (*messagev1.Reaction, error) {
	return nil, nil
}

func (m *MockMessageClient) RemoveReaction(ctx context.Context, in *messagev1.RemoveReactionRequest, _ ...grpc.CallOption) (*messagev1.Empty, error) {
	return nil, nil
}

func (m *MockMessageClient) SearchMessages(ctx context.Context, in *messagev1.SearchMessagesRequest, _ ...grpc.CallOption) (*messagev1.SearchMessagesResponse, error) {
	return nil, nil
}

func (m *MockMessageClient) MarkAsRead(ctx context.Context, in *messagev1.MarkAsReadRequest, _ ...grpc.CallOption) (*messagev1.Empty, error) {
	return nil, nil
}

func (m *MockMessageClient) GetUnreadCount(ctx context.Context, in *messagev1.GetUnreadCountRequest, _ ...grpc.CallOption) (*messagev1.GetUnreadCountResponse, error) {
	return nil, nil
}

func (m *MockMessageClient) ConnectWebSocket(ctx context.Context, in *messagev1.ConnectWebSocketRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[messagev1.WebSocketMessage], error) {
	return nil, nil
}

func (m *MockMessageClient) GetLastMessage(ctx context.Context, in *messagev1.GetLastMessageRequest, _ ...grpc.CallOption) (*messagev1.Message, error) {
	if m.GetLastMessageFn != nil {
		return m.GetLastMessageFn(ctx, in)
	}
	return nil, nil
}

func (m *MockMessageClient) DeleteChatMessages(ctx context.Context, in *messagev1.DeleteChatMessagesRequest, _ ...grpc.CallOption) (*messagev1.Empty, error) {
	return nil, nil
}
