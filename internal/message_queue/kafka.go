package message_queue

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"live-chat-kafka/config"
)

type kafkaClient struct {
	cfg      config.Kafka
	consumer *kafka.Consumer
}

func NewKafkaConsumerClient(cfg config.Kafka) (Client, error) {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.URL,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		return nil, err
	}

	return &kafkaClient{
		cfg:      cfg,
		consumer: consumer,
	}, nil
}

func (k *kafkaClient) Subscribe(topic string) error {
	if err := k.consumer.Subscribe(topic, nil); err != nil {
		return err
	}
	return nil
}
