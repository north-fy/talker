package service

import (
	"context"
	"encoding/json"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"go.uber.org/zap"
)

type StorageReaction interface {
	InsertReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error)
	DeleteReaction(ctx context.Context, req dto.RemoveReactionRequest) (string, error)
}

type ReactionService struct {
	log      *zap.Logger
	storage  StorageReaction
	eventbus event.EventBus
}

func NewReactionService(log *zap.Logger, storage StorageReaction, ev event.EventBus) *ReactionService {
	return &ReactionService{
		log:      log,
		storage:  storage,
		eventbus: ev,
	}
}

func (s *ReactionService) AddReaction(ctx context.Context, req dto.AddReactionRequest) (dto.Reaction, error) {
	log := s.log.With(zap.Any("request", req))

	req.UserID = ctx.Value(models.UserIDKey).(int64)

	react, err := s.storage.InsertReaction(ctx, req)
	if err != nil {
		log.Error("failed to create reaction", zap.Error(err))
		return dto.Reaction{}, err
	}

	var mapReact map[string]any
	if err = json.Unmarshal([]byte(react.Reaction), &mapReact); err != nil {
		log.Error("failed to unmarshal reactions", zap.Error(err))
		return dto.Reaction{}, err
	}

	count, ok := mapReact[req.Reaction].(int32)
	if !ok {
		log.Error("failed to get reaction count", zap.Any("reaction", mapReact))
		return dto.Reaction{}, err
	}

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_ReactionAdded{
			ReactionAdded: &messagev1.ReactionAddedEvent{
				MessageId: react.MessageID,
				UserId:    react.UserID,
				Reaction:  react.Reaction,
				NewCount:  count,
			},
		},
	}

	eventData, err := json.Marshal(&wsData)
	if err != nil {
		log.Error("failed to marshal websocket data", zap.Error(err))
		return dto.Reaction{}, err
	}

	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_ReactionAdded,
		ChatID:  0, // TODO: get chat id from ctx
		Payload: eventData,
	}

	if err = s.eventbus.Publish(ctx, &ev); err != nil {
		log.Error("failed to publish msg for stream",
			zap.Int64("chat_id", ev.GetChatID()),
			zap.Error(err))
		return dto.Reaction{}, err
	}

	return react, nil
}

func (s *ReactionService) RemoveReaction(ctx context.Context, req dto.RemoveReactionRequest) error {
	log := s.log.With(zap.Any("request", req))

	req.UserID = ctx.Value(models.UserIDKey).(int64)

	react, err := s.storage.DeleteReaction(ctx, req)
	if err != nil {
		log.Error("failed to delete reaction", zap.Error(err))
		return err
	}

	var mapReact map[string]any
	if err = json.Unmarshal([]byte(react), &mapReact); err != nil {
		log.Error("failed to unmarshal reactions", zap.Error(err))
		return err
	}

	count, ok := mapReact[req.Reaction].(int32)
	if !ok {
		log.Error("failed to get reaction count", zap.Any("reaction", mapReact))
		return err
	}

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_ReactionRemoved{
			ReactionRemoved: &messagev1.ReactionRemovedEvent{
				MessageId: req.MessageID,
				UserId:    req.UserID,
				Reaction:  req.Reaction,
				NewCount:  count,
			},
		},
	}

	eventData, err := json.Marshal(&wsData)
	if err != nil {
		log.Error("failed to marshal websocket data", zap.Error(err))
		return err
	}

	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_ReactionRemoved,
		ChatID:  0, // TODO: get chat id from ctx
		Payload: eventData,
	}

	if err = s.eventbus.Publish(ctx, &ev); err != nil {
		log.Error("failed to publish msg for stream",
			zap.Int64("chat_id", ev.GetChatID()),
			zap.Error(err))
		return err
	}

	return nil
}
