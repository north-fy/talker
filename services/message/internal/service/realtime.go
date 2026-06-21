package service

import (
	"context"
	"encoding/json"
	"fmt"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"go.uber.org/zap"
)

type Storage interface {
	PubMessage(ctx context.Context, channel string, data interface{}) error
	SubMessage(ctx context.Context, channels... string) dto.Subscribe
}

type WebSocketService struct {
	log     *zap.Logger
	storage Storage
}

func NewWebSocketService(log *zap.Logger, storage Storage) *WebSocketService {
	return &WebSocketService{
		log:     log,
		storage: storage,
	}
}

func (serv *WebSocketService) HandleClientMessage(ctx context.Context, req dto.EventRequest, sendChan chan *messagev1.WebSocketMessage) error {
	serv.log = serv.log.With(zap.Int64("chat_id", req.ChatID), zap.Int64("user_id", req.UserID))

	sub := serv.storage.SubMessage(ctx, fmt.Sprintf("%d:%d", req.ChatID, req.UserID))
	errChan := make(chan error, 1)

	defer func() {
		serv.log.Info("closing subscription")
		if err := sub.Close(); err != nil {
			serv.log.Error("failed to close subscription", zap.Error(err))
			errChan <- err
		}
		close(errChan)
	}()

	for {
		select {
		case <-ctx.Done():
			serv.log.Warn("context is done:", zap.Error(ctx.Err()))
			return ctx.Err()

		case err := <-errChan:
			return err

		default:
			ch := sub.GetData()
			for subMsg := range ch {
				var wsMsg messagev1.WebSocketMessage
				if err := json.Unmarshal([]byte(subMsg.Payload), &wsMsg); err != nil {
					serv.log.Error("failed to unmarshal message", zap.Error(err))
					return err
				}

				sendChan <- &wsMsg
			}
		}
	}
}
