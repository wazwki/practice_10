package kafka

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitProducer() error {
	addr := fmt.Sprintf("kafka:%s", os.Getenv("KAFKA_PORT"))
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		return fmt.Errorf("producer init err: %w", err)
	}

	return nil
}

func SendMessage(topic, message string, partition int32) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Partition: partition,
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("msg err: %w", err)
	}

	slog.Info(fmt.Sprintf("msg send to partition %d offset %d\n", partition, offset), slog.String("module", "user-service"))
	return nil
}

func CloseProducer() {
	if err := producer.Close(); err != nil {
		slog.Error("Fail producer close", slog.Any("error", err), slog.String("module", "user-service"))
	}
}
