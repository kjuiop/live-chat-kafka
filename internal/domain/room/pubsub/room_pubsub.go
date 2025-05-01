package pubsub

import (
	"context"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/domain/room"
	"live-chat-kafka/internal/message_queue"
)

type roomPubSub struct {
	cfg config.Kafka
	mq  message_queue.Client
}

func NewRoomPubSub(cfg config.Kafka, mq message_queue.Client) room.PubSub {
	return &roomPubSub{
		cfg: cfg,
		mq:  mq,
	}
}

func (r *roomPubSub) CreateChatRoom(ctx context.Context, roomId string) error {

	if err := r.mq.CreateTopic(ctx, roomId); err != nil {
		return err
	}

	return nil
}

func (r *roomPubSub) DeleteChatRoom(ctx context.Context, roomId string) error {

	if err := r.mq.DeleteTopic(ctx, roomId); err != nil {
		return err
	}

	return nil
}
