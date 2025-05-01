package pubsub

import (
	"context"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/message_queue"
	"live-chat-kafka/internal/message_queue/types"
)

const (
	serverTopic = "live-chat-server-info"
)

type PubSub struct {
	cfg config.Kafka
	mq  message_queue.Client
}

func NewSystemPubSub(cfg config.Kafka, mq message_queue.Client) system.PubSub {
	return &PubSub{
		cfg: cfg,
		mq:  mq,
	}
}

func (p *PubSub) CreateChatServerTopic() error {

	if err := p.mq.CreateTopic(context.Background(), serverTopic); err != nil {
		return err
	}

	return nil
}

func (p *PubSub) SubscribeTopic(topic string) error {
	if err := p.mq.Subscribe(topic); err != nil {
		return err
	}
	return nil
}

func (p *PubSub) Poll(timeoutMs int) types.Event {
	return p.mq.Poll(timeoutMs)
}

func (p *PubSub) PublishEvent(topic string, data []byte) (types.Event, error) {
	return p.mq.PublishEvent(topic, data)
}
