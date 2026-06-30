package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"github.com/north-fy/talker/services/message/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StorageMessage interface {
	CreateMessage(ctx context.Context, senderID int64, req dto.SendMessageRequest) (models.Message, error)
	SelectMessages(ctx context.Context, req dto.GetMessagesRequest) ([]*models.Message, error)
	UpdateMessage(ctx context.Context, req dto.EditMessageRequest) (models.Message, error)
	DeleteMessageForUser(ctx context.Context, id int64) error
	DeleteMessage(ctx context.Context, id int64) error
	SelectMessage(ctx context.Context, id int64) (models.Message, error)
}

type MessageFuncService struct {
	log      *zap.Logger
	storage  StorageMessage
	eventbus event.EventBus
}

func NewMessageFuncService(log *zap.Logger, storage StorageMessage, ev event.EventBus) *MessageFuncService {
	return &MessageFuncService{
		log:      log,
		storage:  storage,
		eventbus: ev,
	}
}

func (s *MessageFuncService) SendMessage(ctx context.Context, req dto.SendMessageRequest) (models.Message, error) {
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fe := range validationErrors {
				switch fe.Field() {
				case "Content":
					if fe.Tag() == "required" {
						return models.Message{}, domain.ErrMessageEmptyContent
					}

					return models.Message{}, domain.ErrMessageTooLong
				default:
					return models.Message{}, domain.ErrInvalidStruct
				}
			}
		}

		return models.Message{}, err
	}

	id, ok := ctx.Value(models.UserIDKey).(int64)
	if !ok {
		return models.Message{}, domain.ErrNotAuthenticated
	}

	msg, err := s.storage.CreateMessage(ctx, id, req)
	if err != nil {
		log.Error("failed to create message", zap.Error(err))
		return models.Message{}, domain.ErrInternalStorage
	}

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_NewMessage{
			NewMessage: &messagev1.NewMessageEvent{
				Message: utils.ConvertToProtoMessage(&msg),
			},
		},
	}
	eventData, err := json.Marshal(&wsData)
	if err != nil {
		log.Error("failed to marshal websocket data", zap.Error(err))
		return msg, err
	}

	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_NewMessage,
		ChatID:  req.ChatID,
		Payload: eventData,
	}

	if err := s.eventbus.Publish(ctx, &ev); err != nil {
		log.Error("failed to publish msg for stream",
			zap.Int64("chat_id", ev.GetChatID()),
			zap.Error(err))
		return msg, domain.ErrWebSocketPublish
	}

	return msg, nil
}

func (s *MessageFuncService) GetMessages(ctx context.Context, req dto.GetMessagesRequest) (dto.GetMessagesResponse, error) {
	// TODO: добавить проверку на доступ к чат id
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return dto.GetMessagesResponse{}, domain.ErrInvalidStruct
	}

	messages, err := s.storage.SelectMessages(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("message not found", zap.Error(err))
			return dto.GetMessagesResponse{}, domain.ErrMessageNotFound
		}

		log.Error("failed to select messages", zap.Error(err))
		return dto.GetMessagesResponse{}, domain.ErrInternalStorage
	}

	var (
		isMore bool
		count  int32
	)

	if len(messages) > 1 {
		isMore = true
		count = int32(len(messages))
	}

	return dto.GetMessagesResponse{
		Messages:   messages,
		HasMore:    isMore,
		TotalCount: count,
	}, nil
}

func (s *MessageFuncService) EditMessage(ctx context.Context, req dto.EditMessageRequest) (models.Message, error) {
	// TODO: добавить проверку на доступ к сообщению и чату
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fe := range validationErrors {
				switch fe.Field() {
				case "Content":
					if fe.Tag() == "required" {
						return models.Message{}, domain.ErrMessageEmptyContent
					}

					return models.Message{}, domain.ErrMessageTooLong
				default:
					return models.Message{}, domain.ErrInvalidStruct
				}
			}
		}

		return models.Message{}, err
	}

	msg, err := s.storage.UpdateMessage(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("message not found", zap.Error(err))
			return models.Message{}, domain.ErrMessageNotFound
		}

		log.Error("failed to update message", zap.Error(err))
		return models.Message{}, domain.ErrInternalStorage
	}

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_MessageUpdated{
			MessageUpdated: &messagev1.MessageUpdatedEvent{
				MessageId:  msg.ID,
				NewContent: msg.Content,
				UpdatedAt:  timestamppb.New(msg.UpdatedAt),
			},
		},
	}
	eventData, err := json.Marshal(&wsData)
	if err != nil {
		log.Error("failed to marshal websocket data", zap.Error(err))
		return msg, err
	}

	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_MessageUpdated,
		ChatID:  msg.ChatID,
		Payload: eventData,
	}

	if err := s.eventbus.Publish(ctx, &ev); err != nil {
		log.Error("failed to publish msg for stream",
			zap.Int64("chat_id", ev.GetChatID()),
			zap.Error(err))
		return msg, domain.ErrWebSocketPublish
	}

	return msg, nil
}

func (s *MessageFuncService) DeleteMessage(ctx context.Context, req dto.DeleteMessageRequest) (bool, error) {
	// TODO: добавить проверку на доступ к сообщению и чату

	// TODO: короче нужно еще проверку делать на удаление из таблицыи скрытие у одного пользователя
	// TODO: То есть делать 2 метода у storage по ForEveryone
	// DONE
	log := s.log.With(zap.Any("request", req))

	var err error
	if err = Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return false, domain.ErrInvalidStruct
	}

	if req.ForEveryone {
		err = s.storage.DeleteMessage(ctx, req.MessageID)
	} else {
		err = s.storage.DeleteMessageForUser(ctx, req.MessageID)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("message not found", zap.Error(err))
			return false, domain.ErrMessageNotFound
		}

		log.Error("failed to delete message", zap.Error(err))
		return false, domain.ErrInternalStorage
	}

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_MessageDeleted{
			MessageDeleted: &messagev1.MessageDeletedEvent{
				MessageId:   req.MessageID,
				ForEveryone: req.ForEveryone,
			},
		},
	}
	eventData, err := json.Marshal(&wsData)
	if err != nil {
		log.Error("failed to marshal websocket data", zap.Error(err))
		return true, err
	}

	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_MessageDeleted,
		ChatID:  req.ChatID,
		Payload: eventData,
	}

	if err := s.eventbus.Publish(ctx, &ev); err != nil {
		log.Error("failed to publish msg for stream",
			zap.Int64("chat_id", ev.GetChatID()),
			zap.Error(err))
		return false, domain.ErrWebSocketPublish
	}

	return true, nil
}

func (s *MessageFuncService) GetMessage(ctx context.Context, req dto.GetMessageRequest) (models.Message, error) {
	// TODO: добавить проверку на доступ к чату??
	log := s.log.With(zap.Any("request", req))

	if err := Validator.StructCtx(ctx, &req); err != nil {
		log.Error("failed to validate request", zap.Error(err))
		return models.Message{}, domain.ErrInvalidStruct
	}

	msg, err := s.storage.SelectMessage(ctx, req.MessageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("message not found", zap.Error(err))
			return models.Message{}, domain.ErrMessageNotFound
		}

		log.Error("failed to select message", zap.Error(err))
		return models.Message{}, domain.ErrInternalStorage
	}

	return msg, nil
}
