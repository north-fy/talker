package grpc

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/event"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"github.com/north-fy/talker/services/message/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// inMemoryEventBus — in-memory pub/sub для тестов
type inMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[int64][]chan event.Event
}

func newInMemoryEventBus() *inMemoryEventBus {
	return &inMemoryEventBus{
		subscribers: make(map[int64][]chan event.Event),
	}
}

func (b *inMemoryEventBus) Publish(_ context.Context, ev event.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers[ev.GetChatID()] {
		select {
		case ch <- ev:
		default:
		}
	}
	return nil
}

func (b *inMemoryEventBus) Subscribe(_ context.Context, chatID int64) (<-chan event.Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan event.Event, 100)
	b.subscribers[chatID] = append(b.subscribers[chatID], ch)
	return ch, nil
}

func (b *inMemoryEventBus) Close() error { return nil }

type stubMessageService struct{}

func (s *stubMessageService) SendMessage(context.Context, dto.SendMessageRequest) (models.Message, error) {
	return models.Message{}, nil
}
func (s *stubMessageService) GetMessages(context.Context, dto.GetMessagesRequest) (dto.GetMessagesResponse, error) {
	return dto.GetMessagesResponse{}, nil
}
func (s *stubMessageService) EditMessage(context.Context, dto.EditMessageRequest) (models.Message, error) {
	return models.Message{}, nil
}
func (s *stubMessageService) DeleteMessage(context.Context, dto.DeleteMessageRequest) (bool, error) {
	return true, nil
}
func (s *stubMessageService) GetMessage(context.Context, dto.GetMessageRequest) (models.Message, error) {
	return models.Message{}, nil
}
func (s *stubMessageService) AddReaction(context.Context, dto.AddReactionRequest) (dto.Reaction, error) {
	return dto.Reaction{}, nil
}
func (s *stubMessageService) RemoveReaction(context.Context, dto.RemoveReactionRequest) error {
	return nil
}
func (s *stubMessageService) SearchMessages(context.Context, dto.SearchMessagesRequest) (dto.SearchMessagesResponse, error) {
	return dto.SearchMessagesResponse{}, nil
}
func (s *stubMessageService) MarkAsRead(context.Context, dto.MarkAsReadRequest) error {
	return nil
}
func (s *stubMessageService) GetUnreadCount(context.Context, dto.GetUnreadCountRequest) (dto.GetUnreadCountResponse, error) {
	return dto.GetUnreadCountResponse{}, nil
}
func (s *stubMessageService) GetLastMessage(context.Context, dto.GetLastMessageRequest) (models.Message, error) {
	return models.Message{}, nil
}
func (s *stubMessageService) DeleteChatMessages(context.Context, dto.DeleteChatMessagesRequest) error {
	return nil
}

func setupTestServer(t *testing.T) (*inMemoryEventBus, string, func()) {
	t.Helper()

	bus := newInMemoryEventBus()
	wsSvc := service.NewWebSocketService(zap.NewNop(), bus)
	wsServer := NewServerWebSocket(wsSvc)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	Register(srv, &stubMessageService{}, *wsServer)

	go func() {
		_ = srv.Serve(lis)
	}()

	return bus, lis.Addr().String(), func() {
		srv.GracefulStop()
		bus.Close()
	}
}

func dial(t *testing.T, addr string) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	return conn
}

func TestConnectWebSocket_ReceiveMessageEvent(t *testing.T) {
	bus, addr, teardown := setupTestServer(t)
	defer teardown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn := dial(t, addr)
	defer conn.Close()

	client := messagev1.NewMessageServiceClient(conn)
	stream, err := client.ConnectWebSocket(ctx, &messagev1.ConnectWebSocketRequest{
		ChatId: 100,
		UserId: 1,
	})
	if err != nil {
		t.Fatalf("ConnectWebSocket: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Публикуем событие через protojson
	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_NewMessage{
			NewMessage: &messagev1.NewMessageEvent{
				Message: &messagev1.Message{
					Id:       1,
					ChatId:   100,
					SenderId: 2,
					Content:  "hello world",
				},
			},
		},
	}
	eventData, _ := protojson.Marshal(&wsData)
	if err := bus.Publish(ctx, &event.MessageEvent{
		Type:    event.WebSocketMessage_NewMessage,
		ChatID:  100,
		Payload: eventData,
	}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}

	newMsg := msg.GetNewMessage()
	if newMsg == nil {
		t.Fatal("expected NewMessage, got nil")
	}
	if newMsg.Message.Content != "hello world" {
		t.Errorf("content = %q, want %q", newMsg.Message.Content, "hello world")
	}
	if newMsg.Message.ChatId != 100 {
		t.Errorf("chat_id = %d, want 100", newMsg.Message.ChatId)
	}
}

func TestConnectWebSocket_ReceiveMultipleEvents(t *testing.T) {
	bus, addr, teardown := setupTestServer(t)
	defer teardown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn := dial(t, addr)
	defer conn.Close()

	client := messagev1.NewMessageServiceClient(conn)
	stream, err := client.ConnectWebSocket(ctx, &messagev1.ConnectWebSocketRequest{
		ChatId: 200,
		UserId: 1,
	})
	if err != nil {
		t.Fatalf("ConnectWebSocket: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	events := []messagev1.WebSocketMessage{
		{Event: &messagev1.WebSocketMessage_NewMessage{
			NewMessage: &messagev1.NewMessageEvent{
				Message: &messagev1.Message{Id: 1, ChatId: 200, Content: "msg1"},
			},
		}},
		{Event: &messagev1.WebSocketMessage_MessageUpdated{
			MessageUpdated: &messagev1.MessageUpdatedEvent{MessageId: 1, NewContent: "edited"},
		}},
		{Event: &messagev1.WebSocketMessage_MessageDeleted{
			MessageDeleted: &messagev1.MessageDeletedEvent{MessageId: 1, ForEveryone: true},
		}},
	}

	for _, ev := range events {
		data, _ := protojson.Marshal(&ev)
		_ = bus.Publish(ctx, &event.MessageEvent{
			Type:    event.WebSocketMessage_NewMessage,
			ChatID:  200,
			Payload: data,
		})
	}

	for i := 0; i < 3; i++ {
		msg, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv %d: %v", i, err)
		}
		if msg.GetEvent() == nil {
			t.Errorf("message %d: event is nil", i)
		}
	}
}

func TestConnectWebSocket_ContextCancellation(t *testing.T) {
	_, addr, teardown := setupTestServer(t)
	defer teardown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn := dial(t, addr)
	defer conn.Close()

	client := messagev1.NewMessageServiceClient(conn)

	streamCtx, streamCancel := context.WithCancel(ctx)
	stream, err := client.ConnectWebSocket(streamCtx, &messagev1.ConnectWebSocketRequest{
		ChatId: 300,
		UserId: 1,
	})
	if err != nil {
		t.Fatalf("ConnectWebSocket: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	streamCancel()

	_, err = stream.Recv()
	if err == nil {
		t.Fatal("expected error after cancel, got nil")
	}
	t.Logf("stream closed: %v", err)
}

func TestConnectWebSocket_EventBusToStream(t *testing.T) {
	bus, addr, teardown := setupTestServer(t)
	defer teardown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn := dial(t, addr)
	defer conn.Close()

	client := messagev1.NewMessageServiceClient(conn)
	stream, err := client.ConnectWebSocket(ctx, &messagev1.ConnectWebSocketRequest{
		ChatId: 400,
		UserId: 1,
	})
	if err != nil {
		t.Fatalf("ConnectWebSocket: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	reactionEvent := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_ReactionAdded{
			ReactionAdded: &messagev1.ReactionAddedEvent{
				MessageId: 10,
				UserId:    1,
				Reaction:  "👍",
				NewCount:  3,
			},
		},
	}
	data, _ := protojson.Marshal(&reactionEvent)
	_ = bus.Publish(ctx, &event.MessageEvent{
		Type:    event.WebSocketMessage_ReactionAdded,
		ChatID:  400,
		Payload: data,
	})

	msg, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}
	reaction := msg.GetReactionAdded()
	if reaction == nil {
		t.Fatal("expected ReactionAdded, got nil")
	}
	if reaction.Reaction != "👍" {
		t.Errorf("reaction = %q, want %q", reaction.Reaction, "👍")
	}
	if reaction.NewCount != 3 {
		t.Errorf("count = %d, want 3", reaction.NewCount)
	}
}

func TestConnectWebSocket_HandleClientMessage_SendChan(t *testing.T) {
	bus := newInMemoryEventBus()
	wsSvc := service.NewWebSocketService(zap.NewNop(), bus)
	sendChan := make(chan *messagev1.WebSocketMessage, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = wsSvc.HandleClientMessage(ctx, dto.EventRequest{ChatID: 500, UserID: 1}, sendChan)
	}()

	time.Sleep(50 * time.Millisecond)

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_Typing{
			Typing: &messagev1.TypingEvent{
				ChatId:   500,
				UserId:   2,
				Username: "testuser",
				IsTyping: true,
			},
		},
	}
	data, _ := protojson.Marshal(&wsData)
	_ = bus.Publish(ctx, &event.MessageEvent{
		Type:    event.WebSocketMessage_NewMessage,
		ChatID:  500,
		Payload: data,
	})

	select {
	case msg := <-sendChan:
		typing := msg.GetTyping()
		if typing == nil {
			t.Fatal("expected Typing event")
		}
		if typing.Username != "testuser" {
			t.Errorf("username = %q, want %q", typing.Username, "testuser")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestConnectWebSocket_HandleClientMessage_JSON_Bug(t *testing.T) {
	// Демонстрирует баг: encoding/json не умеет десериализовать protobuf oneof.
	// В проде используется json.Marshal/Unmarshal — это ломает oneof.
	// Нужно заменить на protojson.

	bus := newInMemoryEventBus()
	wsSvc := service.NewWebSocketService(zap.NewNop(), bus)
	sendChan := make(chan *messagev1.WebSocketMessage, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := wsSvc.HandleClientMessage(ctx, dto.EventRequest{ChatID: 600, UserID: 1}, sendChan)
		if err != nil {
			t.Logf("HandleClientMessage returned error (expected): %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_NewMessage{
			NewMessage: &messagev1.NewMessageEvent{
				Message: &messagev1.Message{Id: 1, ChatId: 600, Content: "test"},
			},
		},
	}

	// Используем encoding/json как в проде — должен упасть
	data, _ := json.Marshal(&wsData)
	_ = bus.Publish(ctx, &event.MessageEvent{
		Type:    event.WebSocketMessage_NewMessage,
		ChatID:  600,
		Payload: data,
	})

	select {
	case <-sendChan:
		t.Error("expected HandleClientMessage to fail on json.Unmarshal, but it succeeded")
	case <-time.After(500 * time.Millisecond):
		t.Logf("confirmed: encoding/json cannot unmarshal protobuf oneof (bug in service/realtime.go:41)")
	}

	cancel()
}

func TestConnectWebSocket_Stress(t *testing.T) {
	bus := newInMemoryEventBus()
	wsSvc := service.NewWebSocketService(zap.NewNop(), bus)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const numMessages = 100
	const numClients = 5

	var wg sync.WaitGroup
	var mu sync.Mutex
	received := make(map[int]int)

	for clientID := 0; clientID < numClients; clientID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sendChan := make(chan *messagev1.WebSocketMessage, numMessages)

			go func() {
				_ = wsSvc.HandleClientMessage(ctx, dto.EventRequest{ChatID: 700, UserID: int64(id)}, sendChan)
			}()

			count := 0
			for count < numMessages {
				select {
				case <-sendChan:
					count++
				case <-time.After(5 * time.Second):
					t.Errorf("client %d: timeout after %d msgs", id, count)
					return
				}
			}

			mu.Lock()
			received[id] = count
			mu.Unlock()
		}(clientID)
	}

	time.Sleep(50 * time.Millisecond)

	wsData := messagev1.WebSocketMessage{
		Event: &messagev1.WebSocketMessage_NewMessage{
			NewMessage: &messagev1.NewMessageEvent{
				Message: &messagev1.Message{Id: 1, ChatId: 700, Content: "stress"},
			},
		},
	}
	data, _ := protojson.Marshal(&wsData)

	for i := 0; i < numMessages; i++ {
		_ = bus.Publish(ctx, &event.MessageEvent{
			Type:    event.WebSocketMessage_NewMessage,
			ChatID:  700,
			Payload: data,
		})
	}

	wg.Wait()

	for id, count := range received {
		if count != numMessages {
			t.Errorf("client %d: received %d, want %d", id, count, numMessages)
		}
	}
	t.Logf("stress: %d clients x %d messages OK", numClients, numMessages)
}
