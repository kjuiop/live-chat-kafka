package message_queue

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/message_queue/types"
)

type kafkaClient struct {
	cfg      config.Kafka
	consumer *kafka.Consumer
	producer *kafka.Producer
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

func NewKafkaProducerClient(cfg config.Kafka) (Client, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.URL,
		"client.id":         cfg.ClientID,
		"acks":              "all",
	})
	if err != nil {
		return nil, err
	}

	return &kafkaClient{
		cfg:      cfg,
		producer: producer,
	}, nil
}

func (k *kafkaClient) Subscribe(topic string) error {
	if err := k.consumer.Subscribe(topic, nil); err != nil {
		return err
	}
	return nil
}

func (k *kafkaClient) Poll(timeoutMs int) types.Event {
	ev := k.consumer.Poll(timeoutMs)
	switch event := ev.(type) {
	case *kafka.Message:
		return &types.Message{Value: event.Value}
	case *kafka.Error:
		return &types.Error{Error: event}
	default:
		return nil
	}
}

func (k *kafkaClient) PublishEvent(topic string, data []byte) (types.Event, error) {

	ch := make(chan kafka.Event)

	if err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, ch); err != nil {
		return nil, err
	}

	event := <-ch
	switch e := event.(type) {
	case *kafka.Message:
		return &types.Message{Value: e.Value}, nil
	case *kafka.Error:
		return &types.Error{Error: e}, nil
	default:
		return nil, errors.New("unexpected event type")
	}
}
