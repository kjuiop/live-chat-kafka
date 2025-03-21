package message_queue

type Client interface {
	Subscribe(topic string) error
}
