package main

import (
	"context"
	"fmt"

	"live-service/internal/handlers"
	"live-service/internal/middlewares"
	"live-service/internal/service"
	"live-service/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	logger.LogInit()

	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Fail load env", slog.Any("error", err), slog.String("module", "live-service"))
		os.Exit(1)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	serv := service.NewService()
	hand := handlers.NewHandler(serv, upgrader)

	mux := http.NewServeMux()
	corsMux := middlewares.CorsMiddleware(mux)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%v:%v", host, port),
		ReadHeaderTimeout: 800 * time.Millisecond,
		ReadTimeout:       800 * time.Millisecond,
		Handler:           corsMux,
	}

	mux.HandleFunc("/ws", hand.Handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info(fmt.Sprintf("Server up with address: %v:%v", host, port), slog.String("module", "live-service"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Fail start server", slog.Any("error", err), slog.String("module", "live-service"))
		}
	}()

	<-quit
	slog.Info("Shutting down server...", slog.String("module", "live-service"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Fail shutdown", slog.Any("error", err), slog.String("module", "live-service"))
	} else {
		slog.Info("Server gracefully stopped", slog.String("module", "live-service"))
	}
}
