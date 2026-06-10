package service

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"go.uber.org/zap"
)

type StorageReaction interface {
	InsertReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error)
	DeleteReaction(ctx context.Context, req dto.RemoveReactionRequest) error
}

type ReactionService struct {
	log     *zap.Logger
	storage StorageReaction
}

func NewReactionService(log *zap.Logger, storage StorageReaction) *ReactionService {
	return &ReactionService{
		log:     log,
		storage: storage,
	}
}

func (s *ReactionService) AddReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error) {
	s.log = s.log.With(zap.Any("request", req))

	req.UserID = ctx.Value(models.UserIDKey).(int64)

	react, err := s.storage.InsertReaction(ctx, req)
	if err != nil {
		s.log.Error("failed to create reaction", zap.Error(err))
		return dto.Reaction{}, err
	}

	return react, nil
}

func (s *ReactionService) RemoveReaction(ctx context.Context, req dto.RemoveReactionRequest) error {
	s.log = s.log.With(zap.Any("request", req))

	req.UserID = ctx.Value(models.UserIDKey).(int64)

	if err := s.storage.DeleteReaction(ctx, req); err != nil {
		s.log.Error("failed to delete reaction", zap.Error(err))
		return err
	}

	return nil
}
