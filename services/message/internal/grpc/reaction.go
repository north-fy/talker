package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
)

func (s *serverAPI) AddReaction(ctx context.Context, request *messagev1.AddReactionRequest) (*messagev1.Reaction, error) {
	return nil, nil
}

func (s *serverAPI) RemoveReaction(ctx context.Context, request *messagev1.RemoveReactionRequest) (*messagev1.Empty, error) {
	return nil, nil
}
