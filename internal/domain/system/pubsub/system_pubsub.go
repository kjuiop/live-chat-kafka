package pubsub

import (
	"live-chat-kafka/config"
	"live-chat-kafka/internal/domain/system"
	"live-chat-kafka/internal/message_queue"
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

func (p *PubSub) RegisterSubTopic(topic string) error {
	if err := p.mq.Subscribe(topic); err != nil {
		return err
	}

	return nil
}
