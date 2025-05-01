package pubsub

import (
	"live-chat-kafka/config"
	"live-chat-kafka/internal/message_queue"
)

type PubSub struct {
	cfg config.Kafka
	mq  message_queue.Client
}

func NewRoomPubSub(cfg config.Kafka, mq message_queue.Client) *PubSub {
	return &PubSub{
		cfg: cfg,
		mq:  mq,
	}
}
