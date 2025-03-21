package message_queue

import "live-chat-kafka/internal/message_queue/types"

type Client interface {
	Subscribe(topic string) error
	Poll(timeoutMs int) types.Event
	PublishEvent(topic string, data []byte) (types.Event, error)
}
