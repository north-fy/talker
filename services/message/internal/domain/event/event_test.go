package event_test

import (
	"encoding/json"
	"testing"

	"github.com/north-fy/talker/services/message/internal/domain/event"
)

func TestMessageEvent_Getters(t *testing.T) {
	payload := json.RawMessage(`{"test":"data"}`)
	ev := event.MessageEvent{
		Type:    event.WebSocketMessage_NewMessage,
		ChatID:  10,
		Payload: payload,
	}

	if ev.GetType() != event.WebSocketMessage_NewMessage {
		t.Errorf("GetType() = %v, want %v", ev.GetType(), event.WebSocketMessage_NewMessage)
	}

	if ev.GetChatID() != 10 {
		t.Errorf("GetChatID() = %v, want %v", ev.GetChatID(), 10)
	}

	if string(ev.GetData()) != string(payload) {
		t.Errorf("GetData() = %v, want %v", ev.GetData(), payload)
	}
}

