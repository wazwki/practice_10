package repository

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/models"
	"user-service/pkg/metrics"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageInterface interface {
	Get() ([]*models.User, error)
	Create(i *models.User) error
	Update(i *models.User, id string) error
	Delete(id string) error
}

type Repository struct {
	DataBase *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) StorageInterface {
	return &Repository{DataBase: db}
}

func (repo *Repository) Get() ([]*models.User, error) {
	start := time.Now()

	query := `SELECT name, email FROM users_table`
	rows, err := repo.DataBase.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	metrics.RepositoryDuration.WithLabelValues("Get").Observe(time.Since(start).Seconds())
	return users, nil
}

func (repo *Repository) Create(user *models.User) error {
	start := time.Now()

	query := `INSERT INTO users_table (name, email) VALUES ($1, $2)`
	repo.DataBase.QueryRow(context.Background(), query, user.Name, user.Email)

	metrics.RepositoryDuration.WithLabelValues("Create").Observe(time.Since(start).Seconds())
	return nil
}

func (repo *Repository) Update(user *models.User, id string) error {
	start := time.Now()

	query := `UPDATE users_table SET name = $1, email = $2 WHERE id = $3`
	commandTag, err := repo.DataBase.Exec(context.Background(), query, user.Name, user.Email, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	metrics.RepositoryDuration.WithLabelValues("Update").Observe(time.Since(start).Seconds())
	return nil
}

func (repo *Repository) Delete(id string) error {
	start := time.Now()

	query := `DELETE FROM users_table WHERE id = $1`
	commandTag, err := repo.DataBase.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	metrics.RepositoryDuration.WithLabelValues("Delete").Observe(time.Since(start).Seconds())
	return nil
}
