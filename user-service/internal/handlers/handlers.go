package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"user-service/internal/models"
	"user-service/internal/service"
	"user-service/pkg/metrics"
)

type HandlerInterface interface {
	GetHandler(http.ResponseWriter, *http.Request)
	UpdateHandler(http.ResponseWriter, *http.Request)
	DeleteHandler(http.ResponseWriter, *http.Request)
	CreateHandler(http.ResponseWriter, *http.Request)
}

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(s service.ServiceInterface) HandlerInterface {
	return &Handler{service: s}
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	users, err := h.service.GetUsers()
	if err != nil {
		slog.Error("Fail in get path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		slog.Error("Fail in get path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	metrics.HandlersDuration.WithLabelValues(r.Method, "/users").Observe(time.Since(start).Seconds())
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := r.PathValue("id")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("Fail in put path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateUser(id, &user); err != nil {
		slog.Error("Fail in put path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	metrics.HandlersDuration.WithLabelValues(r.Method, "/users/{id}").Observe(time.Since(start).Seconds())
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := r.PathValue("id")

	if err := h.service.DeleteUser(id); err != nil {
		slog.Error("Fail in delete path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	metrics.HandlersDuration.WithLabelValues(r.Method, "/users/{id}").Observe(time.Since(start).Seconds())
}

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("Fail in post path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		slog.Error("Fail in post path", slog.Any("error", err), slog.String("module", "user-service"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	metrics.HandlersDuration.WithLabelValues(r.Method, "/users").Observe(time.Since(start).Seconds())
}
