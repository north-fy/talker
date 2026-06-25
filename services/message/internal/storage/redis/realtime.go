package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/north-fy/talker/services/message/internal/domain/event"
)

func (s *Storage) Publish(ctx context.Context, ev event.Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	return s.client.Publish(ctx, fmt.Sprintf("chat:%d", ev.GetChatID()), data).Err()
}

func (s *Storage) Subscribe(ctx context.Context, chatID int64) (<-chan event.Event, error) {
	sub := s.client.Subscribe(ctx, fmt.Sprintf("chat:%d", chatID))

	_, err := sub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.subs[chatID] = sub
	s.mu.Unlock()

	ch := make(chan event.Event, 100)
	errCh := make(chan error, 1)

	go func() {
		defer func() {
			close(ch)
			close(errCh)
			sub.Close()

			s.mu.Lock()
			delete(s.subs, chatID)
			s.mu.Unlock()
		}()

		for {
			msg, err := sub.ReceiveMessage(ctx)
			if err != nil {
				errCh <- err
				return
			}

			var msgEvent event.MessageEvent
			if err = json.Unmarshal([]byte(msg.Payload), &msgEvent); err != nil {
				errCh <- err
				return
			}

			ch <- &msgEvent
		}
	}()

	return ch, nil
}
