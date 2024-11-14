package repository

import (
	"context"
	"fmt"
	"user-service/internal/models"

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

	return users, nil
}

func (repo *Repository) Create(user *models.User) error {
	query := `INSERT INTO users_table (name, email) VALUES ($1, $2)`

	_, err := repo.DataBase.Exec(context.Background(), query, user.Name, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Update(user *models.User, id string) error {
	query := `UPDATE users_table SET name = $1, email = $2 WHERE id = $3`
	commandTag, err := repo.DataBase.Exec(context.Background(), query, user.Name, user.Email, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

func (repo *Repository) Delete(id string) error {
	query := `DELETE FROM users_table WHERE id = $1`
	commandTag, err := repo.DataBase.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}
