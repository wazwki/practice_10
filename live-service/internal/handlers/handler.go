package handlers

import (
	"live-service/internal/service"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HandlerInterface interface {
	Handler(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service  service.ServiceInterface
	upgrader websocket.Upgrader
}

func NewHandler(s service.ServiceInterface, u websocket.Upgrader) HandlerInterface {
	return &Handler{service: s, upgrader: u}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.service.HandleConnection(conn)
		}()
	}

	wg.Wait()
}
