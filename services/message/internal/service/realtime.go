package service

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type WebSocketService struct {
	log      *zap.Logger
	eventbus event.EventBus
}

func NewWebSocketService(log *zap.Logger, bus event.EventBus) *WebSocketService {
	return &WebSocketService{
		log:      log,
		eventbus: bus,
	}
}

func (serv *WebSocketService) HandleClientMessage(ctx context.Context, req dto.EventRequest, sendChan chan *messagev1.WebSocketMessage) error {
	log := serv.log.With(zap.Int64("chat_id", req.ChatID), zap.Int64("user_id", req.UserID))

	ch, err := serv.eventbus.Subscribe(ctx, req.ChatID)
	if err != nil {
		log.Error("failed to subscribe eventbus", zap.Error(err))
		return domain.ErrWebSocketSubscribe
	}

	for {
		select {
		case <-ctx.Done():
			log.Warn("context is done:", zap.Error(ctx.Err()))
			return ctx.Err()

		case subMsg := <-ch:
			var wsMsg messagev1.WebSocketMessage
			if err := protojson.Unmarshal(subMsg.GetData(), &wsMsg); err != nil {
				log.Error("failed to unmarshal message",
					zap.Int32("type", subMsg.GetType()),
					zap.Error(err))
				return err
			}

			sendChan <- &wsMsg
		}
	}
}
