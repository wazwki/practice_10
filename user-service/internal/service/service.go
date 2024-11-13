package service

import (
	"time"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/pkg/metrics"
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
	start := time.Now()
	users, err := s.repo.Get()
	if err != nil {
		return nil, err
	}
	metrics.ServiceDuration.WithLabelValues("GetUser").Observe(time.Since(start).Seconds())
	return users, nil
}

func (s *Service) CreateUser(user *models.User) error {
	start := time.Now()
	err := s.repo.Create(user)
	if err != nil {
		return err
	}
	metrics.ServiceDuration.WithLabelValues("CreateUser").Observe(time.Since(start).Seconds())
	return nil
}

func (s *Service) UpdateUser(id string, user *models.User) error {
	start := time.Now()
	err := s.repo.Update(user, id)
	if err != nil {
		return err
	}
	metrics.ServiceDuration.WithLabelValues("UpdateUser").Observe(time.Since(start).Seconds())
	return nil
}

func (s *Service) DeleteUser(id string) error {
	start := time.Now()
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	metrics.ServiceDuration.WithLabelValues("DeleteUser").Observe(time.Since(start).Seconds())
	return nil
}
