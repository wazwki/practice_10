package main

import (
	"context"
	"fmt"
	"log/slog"
	"notification-service/pkg/kafka"
	"notification-service/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	logger.LogInit()

	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Fail load env", slog.Any("error", err), slog.String("module", "notification-service"))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	if err := kafka.InitConsumer(); err != nil {
		slog.Error("Fail init consumer", slog.Any("error", err), slog.String("module", "notification-service"))
	}
	defer kafka.CloseConsumer()

	msgs := kafka.GetMessage("registration-topic", 0)
	defer kafka.ClosePartitionConsumer()

	wg := sync.WaitGroup{}
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Message channel closed")
				return
			}
			wg.Add(1)
			go func(value []byte) {
				defer wg.Done()
				if err := SendNotification(string(value)); err != nil {
					slog.Error("Fail send notification", slog.Any("error", err), slog.String("module", "notification-service"))
				}
			}(msg.Value)
		case <-quit:
			slog.Info("Shutting down server...", slog.String("module", "notification-service"))
			cancel()
			wg.Wait()
			return
		case <-ctx.Done():
			slog.Info("Context canceled", slog.String("module", "notification-service"))
			wg.Wait()
			return
		}
	}
}

func SendNotification(email string) error {
	fmt.Printf("Send notification to user: %s\n", email)
	return nil
}
