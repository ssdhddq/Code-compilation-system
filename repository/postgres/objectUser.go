package postgres

import (
	"Code-compilation-system/repository"
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (r *Object) RegisterUser(user *repository.User) error {
	query := `INSER INTO users (id, login, password_hash) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(context.Background(), query, user.Id, user.Login, user.Password)
	return err
}

func (r *Object) AuthUser(s string, s2 string) (bool, *uuid.UUID) {
	var userID uuid.UUID
	var hash string
	query := `SELECT id, password_hash FROM users WHERE login = $1`
	err := r.pool.QueryRow(context.Background(), query, s).Scan(&userID, &hash)
	if err != nil {
		return false, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s2)); err != nil {
		return false, nil
	}
	return true, &userID
}

func (r *Object) DeleteUser(uuid uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.pool.Exec(context.Background(), query, uuid)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return repository.NotFound
	}
	return nil
}

func (r *Object) GetUserByLogin(login string) (*repository.User, error) {
	var user repository.User
	query := `SELECT id, login, password_hash FROM users WHERE login = $1`
	err := r.pool.QueryRow(context.Background(), query, login).Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
