package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/domain/chat"
	"live-chat-kafka/internal/message_queue"
	"live-chat-kafka/internal/message_queue/types"
	"log/slog"
	"sync"
)

const (
	prefix = "chat"
)

type chatPubSub struct {
	cfg        config.Kafka
	mq         message_queue.Client
	subscribed map[string]struct{}
	mu         sync.Mutex
}

func NewChatPubSub(cfg config.Kafka, mq message_queue.Client) chat.PubSub {
	return &chatPubSub{
		cfg:        cfg,
		mq:         mq,
		subscribed: make(map[string]struct{}),
		mu:         sync.Mutex{},
	}
}

func (c *chatPubSub) SubscribeTopic(ctx context.Context, topic string, handler func(msg *chat.Message)) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.subscribed[topic]; exists {
		return nil // 이미 구독 중
	}

	kafkaTopic := fmt.Sprintf("%s_%s", prefix, topic)

	if err := c.mq.Subscribe(kafkaTopic); err != nil {
		return fmt.Errorf("failed to subscribe topic %s: %w", kafkaTopic, err)
	}

	c.subscribed[topic] = struct{}{}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Debug("close Loop SubscribeTopic goroutine")
				return
			default:
				ev := c.mq.Poll(100) // 100ms 동안 폴링

				if ev == nil {
					continue
				}

				if ev.IsError() {
					errorEvent := ev.(*types.Error)
					slog.Error("Failed to Polling event", "topic", topic, "error", errorEvent.Error)
					continue
				}

				if !ev.IsMessage() {
					slog.Error("is not expected message", "topic", topic, "event", ev)
					continue
				}

				message, ok := ev.(*types.Message)
				if !ok {
					slog.Error("event is message type but cast failed", "topic", topic, "event", ev)
					continue
				}

				var decoder chat.Message
				if err := json.Unmarshal(message.Value, &decoder); err != nil {
					slog.Error("Failed to Unmarshal event", "topic", topic, "error", err, "event", ev)
					continue
				}

				slog.Debug("kafka message", "topic", topic, "message", string(message.Value))
				handler(&decoder)
			}
		}
	}()

	return nil
}

func (c *chatPubSub) PublishMessage(roomId string, message *chat.Message) error {
	// topic 이름 구성
	kafkaTopic := fmt.Sprintf("%s_%s", prefix, roomId)

	// 메시지 직렬화
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	event, err := c.mq.PublishEvent(kafkaTopic, data)
	if err != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", kafkaTopic, err)
	}

	slog.Debug("published kafka message", "topic", kafkaTopic, "message", string(data), "event", event)
	return nil
}
