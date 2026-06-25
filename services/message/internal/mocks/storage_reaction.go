package mocks

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
)

type MockStorageReaction struct {
	InsertReactionFn func(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error)
	DeleteReactionFn func(ctx context.Context, req dto.RemoveReactionRequest) (string, error)
}

func (m *MockStorageReaction) InsertReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error) {
	if m.InsertReactionFn != nil {
		return m.InsertReactionFn(ctx, req)
	}
	return dto.Reaction{}, nil
}

func (m *MockStorageReaction) DeleteReaction(ctx context.Context, req dto.RemoveReactionRequest) (string, error) {
	if m.DeleteReactionFn != nil {
		return m.DeleteReactionFn(ctx, req)
	}
	return "", nil
}
