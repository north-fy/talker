package service

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"go.uber.org/zap"
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
	log     *zap.Logger
	storage StorageMessage
}

func NewMessageFuncService(log *zap.Logger, storage StorageMessage) *MessageFuncService {
	return &MessageFuncService{
		log:     log,
		storage: storage,
	}
}

func (s *MessageFuncService) SendMessage(ctx context.Context, req dto.SendMessageRequest) (models.Message, error) {
	s.log = s.log.With(zap.Any("request", req))

	id := ctx.Value(models.UserIDKey).(int64)

	msg, err := s.storage.CreateMessage(ctx, id, req)
	if err != nil {
		s.log.Error("failed to create message", zap.Error(err))
		return models.Message{}, err
	}

	return msg, nil
}

func (s *MessageFuncService) GetMessages(ctx context.Context, req dto.GetMessagesRequest) (dto.GetMessagesResponse, error) {
	// TODO: добавить проверку на доступ к чат id

	s.log = s.log.With(zap.Any("request", req))

	messages, err := s.storage.SelectMessages(ctx, req)
	if err != nil {
		s.log.Error("failed to select messages", zap.Error(err))
		return dto.GetMessagesResponse{}, err
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

	s.log = s.log.With(zap.Any("request", req))

	msg, err := s.storage.UpdateMessage(ctx, req)
	if err != nil {
		s.log.Error("failed to update message", zap.Error(err))
		return models.Message{}, err
	}

	return msg, nil
}

func (s *MessageFuncService) DeleteMessage(ctx context.Context, req dto.DeleteMessageRequest) (bool, error) {
	// TODO: добавить проверку на доступ к сообщению и чату

	// TODO: короче нужно еще проверку делать на удаление из таблицыи скрытие у одного пользователя
	// TODO: То есть делать 2 метода у storage по ForEveryone
	// DONE

	s.log = s.log.With(zap.Any("request", req))

	var err error
	if req.ForEveryone {
		err = s.storage.DeleteMessage(ctx, req.MessageID)
	} else {
		err = s.storage.DeleteMessageForUser(ctx, req.MessageID)
	}

	if err != nil {
		s.log.Error("failed to delete message", zap.Error(err))
		return false, err
	}

	return true, nil
}

func (s *MessageFuncService) GetMessage(ctx context.Context, req dto.GetMessageRequest) (models.Message, error) {
	// TODO: добавить проверку на доступ к чату??

	s.log = s.log.With(zap.Any("request", req))

	msg, err := s.storage.SelectMessage(ctx, req.MessageID)
	if err != nil {
		s.log.Error("failed to select message", zap.Error(err))
		return models.Message{}, err
	}

	return msg, nil
}
