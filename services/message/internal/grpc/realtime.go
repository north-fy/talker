package grpc

import (
	"context"
	"sync"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
)

// TODO: переписать в 1 сервис потом через него взаимодействовать для разных целей в бд
//type WebSocketService interface {
//	NewMessageEvent(ctx context.Context, req dto.EventRequest) (dto.NewMessageEvent, error)
//	IsUpdatedMessageEvent(ctx context.Context, req dto.EventRequest) (dto.MessageUpdatedEvent, error)
//	IsDeletedMessageEvent(ctx context.Context, req dto.EventRequest) (dto.MessageDeletedEvent, error)
//	NewReactionEvent(ctx context.Context, req dto.EventRequest) (dto.ReactionAddedEvent, error)
//	IsDeletedReactionEvent(ctx context.Context, req dto.EventRequest) (dto.ReactionRemovedEvent, error)
//	TypingEvent(ctx context.Context, req dto.EventRequest) (dto.TypingEvent, error)
//	ReadReceiptEvent(ctx context.Context, req dto.EventRequest) (dto.ReadReceiptEvent, error)
//}

type WebSocketService interface {
	HandleClientMessage(ctx context.Context, req dto.EventRequest, sendChan chan *messagev1.WebSocketMessage) error
}

type ServerWebSocket struct {
	serv    WebSocketService
	mu      *sync.RWMutex
	clients map[int64]chan *messagev1.WebSocketMessage
}

func NewServerWebSocket(serv WebSocketService) *ServerWebSocket {
	return &ServerWebSocket{
		serv:    serv,
		mu:      &sync.RWMutex{},
		clients: make(map[int64]chan *messagev1.WebSocketMessage),
	}
}

func (s *serverAPI) ConnectWebSocket(req *messagev1.ConnectWebSocketRequest, stream messagev1.MessageService_ConnectWebSocketServer) error {
	clientID, chatID := req.GetUserId(), req.GetChatId()

	msgChan := make(chan *messagev1.WebSocketMessage, 100)
	errChan := make(chan error, 1)

	s.ws.mu.Lock()
	s.ws.clients[clientID] = msgChan
	s.ws.mu.Unlock()

	defer func() {
		close(msgChan)

		s.ws.mu.Lock()
		delete(s.ws.clients, clientID)
		s.ws.mu.Unlock()
	}()

	// функция для перехавата сообщений чата
	go func() {
		err := s.ws.serv.HandleClientMessage(stream.Context(), dto.EventRequest{ChatID: chatID, UserID: clientID}, msgChan)
		if err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	// отправление в стрим поток
	for {
		select {
		case err := <-errChan:
			return err

		case msg := <-msgChan:
			if err := stream.Send(msg); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
