package message_queue

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"live-chat-kafka/config"
	"live-chat-kafka/internal/message_queue/types"
	"log/slog"
)

type kafkaClient struct {
	cfg      config.Kafka
	consumer *kafka.Consumer
	producer *kafka.Producer
	admin    *kafka.AdminClient
}

func NewKafkaClient(cfg config.Kafka, withProducer, withConsumer bool) (Client, error) {

	var producer *kafka.Producer
	var consumer *kafka.Consumer
	var err error

	if withProducer {
		producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": cfg.URL,
			"client.id":         cfg.ClientID,
			"acks":              "all",
		})
		if err != nil {
			return nil, err
		}
	}

	if withConsumer {
		consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": cfg.URL,
			"group.id":          cfg.GroupID,
			"auto.offset.reset": "latest",
		})
		if err != nil {
			return nil, err
		}
	}

	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": cfg.URL,
	})
	if err != nil {
		return nil, err
	}

	return &kafkaClient{
		cfg:      cfg,
		producer: producer,
		consumer: consumer,
		admin:    admin,
	}, nil
}

func (k *kafkaClient) CreateTopic(ctx context.Context, topic string) error {

	metadata, err := k.admin.GetMetadata(&topic, false, 5000)
	if err != nil {
		return err
	}

	if _, exists := metadata.Topics[topic]; exists {
		return nil
	}

	_, err = k.admin.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (k *kafkaClient) DeleteTopic(ctx context.Context, topic string) error {
	_, err := k.admin.DeleteTopics(ctx, []string{topic})
	return err
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

func (k *kafkaClient) Close() {
	k.admin.Close()
	if k.consumer != nil {
		if err := k.consumer.Close(); err != nil {
			slog.Error("failed to close kafka consumer", "error", err)
		}
	}
	if k.producer != nil {
		k.producer.Close()
	}
}
