package producer

import (
	"github.com/IBM/sarama"
	"log/slog"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	logger   *slog.Logger
}

// NewKafkaProducer creates a new KafkaProducer
func NewKafkaProducer(hosts []string, logger *slog.Logger) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(hosts, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		producer: producer,
		logger:   logger.With("kafka", "producer"),
	}, nil
}

func (k *KafkaProducer) SendMessage(topic string, message string) error {
	partition, offset, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	if err != nil {
		k.logger.Error("Failed to send message", "err", err, "topic", topic, "message", message)
	} else {
		k.logger.Info("Message sent", "topic", topic, "message", message, "partition", partition, "offset", offset)
	}
	return err
}

func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}
