package service

import (
	"log/slog"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/pkg/kafka"
)

type ServiceInterface interface {
	GetUsers() ([]*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(id string, user *models.User) error
	DeleteUser(id string) error
}

type Service struct {
	repo repository.StorageInterface
}

func NewService(repo repository.StorageInterface) ServiceInterface {
	return &Service{repo: repo}
}

func (s *Service) GetUsers() ([]*models.User, error) {
	users, err := s.repo.Get()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) CreateUser(user *models.User) error {
	err := s.repo.Create(user)
	if err != nil {
		return err
	}

	if err := kafka.InitProducer(); err != nil {
		slog.Error("Fail init produser", slog.Any("error", err), slog.String("module", "user-service"))
	}
	defer kafka.CloseProducer()
	if err := kafka.SendMessage("registration-topic", user.Email, 0); err != nil {
		slog.Error("Fail send message to kafka", slog.Any("error", err), slog.String("module", "user-service"))
	}

	return nil
}

func (s *Service) UpdateUser(id string, user *models.User) error {
	err := s.repo.Update(user, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteUser(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
