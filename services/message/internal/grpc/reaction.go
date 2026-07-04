package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReactionService interface {
	AddReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error)
	RemoveReaction(ctx context.Context, req dto.RemoveReactionRequest) error
}

func (s *serverAPI) AddReaction(ctx context.Context, req *messagev1.AddReactionRequest) (*messagev1.Reaction, error) {
	rectReg := dto.AddReactionRequest{
		MessageID: req.GetMessageId(),
		Reaction:  req.GetReaction(),
	}

	reaction, err := s.serv.AddReaction(ctx, rectReg)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &messagev1.Reaction{
		MessageId: reaction.MessageID,
		UserId:    reaction.UserID,
		Reaction:  reaction.Reaction,
		CreatedAt: timestamppb.New(reaction.CreatedAt),
	}, nil
}

func (s *serverAPI) RemoveReaction(ctx context.Context, req *messagev1.RemoveReactionRequest) (*messagev1.Empty, error) {
	rectReg := dto.RemoveReactionRequest{
		MessageID: req.GetMessageId(),
		Reaction:  req.GetReaction(),
	}

	if err := s.serv.RemoveReaction(ctx, rectReg); err != nil {
		return nil, toGRPC(err)
	}

	return nil, nil
}
