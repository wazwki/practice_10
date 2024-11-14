package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/db"
	"user-service/internal/handlers"
	"user-service/internal/middlewares"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	logger.LogInit()

	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Fail load env", slog.Any("error", err), slog.String("module", "user-service"))
		os.Exit(1)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("US_PORT")
	if port == "" {
		port = "8080"
	}

	pool, err := db.ConnectPool()
	if err != nil {
		slog.Error("Fail connect pool", slog.Any("error", err), slog.String("module", "user-service"))
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("Connect pool success", slog.String("module", "user-service"))

	if err := db.RunMigrations(); err != nil {
		slog.Error("Fail migrate", slog.Any("error", err), slog.String("module", "user-service"))
		os.Exit(1)
	}
	slog.Info("Migrate success", slog.String("module", "user-service"))

	mux := http.NewServeMux()
	corsMux := middlewares.CorsMiddleware(mux)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%v:%v", host, port),
		ReadHeaderTimeout: 800 * time.Millisecond,
		ReadTimeout:       800 * time.Millisecond,
		Handler:           corsMux,
	}

	repository := repository.NewRepository(pool)
	service := service.NewService(repository)
	handlers := handlers.NewHandler(service)

	mux.HandleFunc("GET /users", handlers.GetHandler)
	mux.HandleFunc("PUT /users/{id}", handlers.UpdateHandler)
	mux.HandleFunc("DELETE /users/{id}", handlers.DeleteHandler)
	mux.HandleFunc("POST /users", handlers.CreateHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info(fmt.Sprintf("Server up with address: %v:%v", host, port), slog.String("module", "user-service"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Fail start server", slog.Any("error", err), slog.String("module", "user-service"))
		}
	}()

	<-quit
	slog.Info("Shutting down server...", slog.String("module", "user-service"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Fail shutdown", slog.Any("error", err), slog.String("module", "user-service"))
	} else {
		slog.Info("Server gracefully stopped", slog.String("module", "user-service"))
	}
}
