package event

import "encoding/json"

type MessageEvent struct {
	Type    int32         `json:"type"`
	ChatID  int64           `json:"chat_id"`
	Payload json.RawMessage `json:"payload"`
}

func (m *MessageEvent) GetType() int32 {
	return m.Type
}

func (m *MessageEvent) GetChatID() int64 {
	return m.ChatID
}

func (m *MessageEvent) GetData() []byte {
	return m.Payload
}