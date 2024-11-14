package kafka

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/IBM/sarama"
)

var consumer sarama.Consumer
var partitionConsumer sarama.PartitionConsumer

func InitConsumer() error {
	addr := fmt.Sprintf("kafka:%s", os.Getenv("KAFKA_PORT"))
	var err error
	consumer, err = sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		return fmt.Errorf("consumer init err: %w", err)
	}

	return nil
}

func GetMessage(topic string, partition int32) <-chan *sarama.ConsumerMessage {
	var err error
	partitionConsumer, err = consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		slog.Error("Fail partition consumer", slog.Any("error", err), slog.String("module", "notification-service"))
	}

	return partitionConsumer.Messages()
}

func CloseConsumer() {
	if err := consumer.Close(); err != nil {
		slog.Error("Fail consumer close", slog.Any("error", err), slog.String("module", "notification-service"))
	}
}

func ClosePartitionConsumer() {
	if err := partitionConsumer.Close(); err != nil {
		slog.Error("Fail partition consumer close", slog.Any("error", err), slog.String("module", "notification-service"))
	}
}
