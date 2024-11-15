package service

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type ServiceInterface interface {
	HandleConnection(conn *websocket.Conn)
}

type Service struct{}

func NewService() ServiceInterface {
	return &Service{}
}

func (s *Service) HandleConnection(conn *websocket.Conn) {
	msgTextChan := make(chan []byte)
	msgBinChan := make(chan []byte)

	go func() {
		for {
			msgType, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Error reading message", slog.Any("error", err), slog.String("module", "live-service"))
				return
			}

			if msgType == websocket.TextMessage {
				msgTextChan <- message
				slog.Info("Get text msg")
			} else if msgType == websocket.BinaryMessage {
				msgBinChan <- message
				slog.Info("Get binary msg")
			} else {
				slog.Error("Unsupported message type", slog.String("module", "live-service"))
				return
			}
		}
	}()

	for {
		select {
		case textMsg := <-msgTextChan:
			go func(message []byte) {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					slog.Error("Error writing message", slog.Any("error", err), slog.String("module", "live-service"))
					return
				}
				slog.Info("Send text msg")
			}(textMsg)
		case binaryMsg := <-msgBinChan:
			go func(message []byte) {
				err := conn.WriteMessage(websocket.BinaryMessage, message)
				if err != nil {
					slog.Error("Error writing message", slog.Any("error", err), slog.String("module", "live-service"))
					return
				}
				slog.Info("Send binary msg")
			}(binaryMsg)
		}
	}
}
