package consumer

import (
	"github.com/IBM/sarama"
	"log/slog"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	logger   *slog.Logger
}

// NewKafkaConsumer creates a new KafkaConsumer
func NewKafkaConsumer(hosts []string, logger *slog.Logger) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(hosts, config)

	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{
		consumer: consumer,
		logger:   logger.With("kafka", "consumer"),
	}, nil
}

// MessageChannel reads a message from the given topic
func (k *KafkaConsumer) MessageChannel(topic string) (chan string, error) {
	messages := make(chan string)
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		return nil, err
	}
	for _, partition := range partitionList {
		go func(partition int32) {
			pc, err := k.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				k.logger.Error("Failed to consume partition", "err", err)
				return
			}
			defer pc.Close()

			for message := range pc.Messages() {
				messages <- string(message.Value)
			}
		}(partition)
	}
	return messages, nil
}

// Close closes the consumer
func (k *KafkaConsumer) Close() error {
	return k.consumer.Close()
}
