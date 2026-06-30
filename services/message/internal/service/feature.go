package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
type FeatureService interface {
	SearchMessages(ctx context.Context, req dto.SearchMessagesRequest) (dto.SearchMessagesResponse, error)
	MarkAsRead(ctx context.Context, req dto.MarkAsReadRequest) error
	GetUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error)
	ConnectWebSocket(ctx context.Context, req dto.ConnectWebSocketRequest) error
	GetLastMessage(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error)
	DeleteChatMessages(ctx context.Context, req dto.DeleteChatMessagesRequest) error
}
*/

type StorageFeature interface {
	SearchMessages(ctx context.Context, req dto.SearchMessagesRequest) ([]*models.Message, error) // Реализация не нужна? тогда dto надо исправлять
	SetAsRead(ctx context.Context, req dto.MarkAsReadRequest) error
	SelectUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error)
	SelectLastMessage(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error)
	DeleteChatMessages(ctx context.Context, req dto.DeleteChatMessagesRequest) error
}

type FeatureService struct {
	log      *zap.Logger
	storage  StorageFeature
	eventbus event.EventBus
}

func NewFeatureService(log *zap.Logger, storage StorageFeature, bus event.EventBus) *FeatureService {
	return &FeatureService{
		log:      log,
		storage:  storage,
		eventbus: bus,
	}
}

func (s *FeatureService) SearchMessages(ctx context.Context, req dto.SearchMessagesRequest) (dto.SearchMessagesResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.SearchMessagesResponse{}, domain.ErrInvalidStruct
	}

	messages, err := s.storage.SearchMessages(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("messages not found", zap.Error(err))
			return dto.SearchMessagesResponse{}, domain.ErrMessageNotFound
		}

		log.Error("failed to search messages", zap.Error(err))
		return dto.SearchMessagesResponse{}, domain.ErrInternalStorage
	}

	resp := dto.SearchMessagesResponse{
		Messages: messages,
	}

	if len(messages) > 0 {
		resp.HasMore = true
	}

	return resp, nil
}

func (s *FeatureService) MarkAsRead(ctx context.Context, req dto.MarkAsReadRequest) error {
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return domain.ErrInvalidStruct
	}

	if err := s.storage.SetAsRead(ctx, req); err != nil {
		log.Error("failed to mark messages as read", zap.Error(err))
		return domain.ErrInternalStorage
	}

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_ReadReceipt{
			ReadReceipt: &messagev1.ReadReceiptEvent{
				ChatId:            req.ChatID,
				UserId:            req.UserID,
				ReadUpToMessageId: strconv.Itoa(int(req.UpToMessageID)),
				ReadAt:            timestamppb.New(time.Now()),
			},
		},
	}

	eventData, err := json.Marshal(&wsData)
	if err != nil {
		log.Error("failed to marshal websocket data", zap.Error(err))
		return err
	}

	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_ReadReceipt,
		ChatID:  req.ChatID,
		Payload: eventData,
	}

	if err = s.eventbus.Publish(ctx, &ev); err != nil {
		log.Error("failed to publish msg for stream",
			zap.Int64("chat_id", ev.GetChatID()),
			zap.Error(err))
		return domain.ErrWebSocketPublish
	}

	return nil
}

func (s *FeatureService) GetUnreadCount(ctx context.Context, req dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error) {
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.GetUnreadCountResponse{}, domain.ErrInvalidStruct
	}

	resp, err := s.storage.SelectUnreadCount(ctx, req)
	if err != nil {
		log.Error("failed to get unread count", zap.Error(err))
		return dto.GetUnreadCountResponse{}, domain.ErrInternalStorage
	}

	return resp, nil
}

func (s *FeatureService) GetLastMessage(ctx context.Context, req dto.GetLastMessageRequest) (models.Message, error) {
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Message{}, domain.ErrInvalidStruct
	}

	resp, err := s.storage.SelectLastMessage(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("last message not found", zap.Error(err))
			return models.Message{}, domain.ErrMessageNotFound
		}

		log.Error("failed to select last message", zap.Error(err))
		return models.Message{}, domain.ErrInternalStorage
	}

	return resp, nil
}

func (s *FeatureService) DeleteChatMessages(ctx context.Context, req dto.DeleteChatMessagesRequest) error {
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return domain.ErrInvalidStruct
	}

	if err := s.storage.DeleteChatMessages(ctx, req); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("chat not found", zap.Error(err))
			return domain.ErrChatNotFound
		}

		log.Error("failed to delete chat messages", zap.Error(err))
		return domain.ErrInternalStorage
	}

	return nil
}
