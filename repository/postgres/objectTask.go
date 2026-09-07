package postgres

import (
	"Code-compilation-system/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Object) GetTask(uuid uuid.UUID) (*repository.Task, error) {
	var task repository.Task
	query := `SELECT id, translator, code, result, status FROM tasks WHERE id = $1`
	err := r.pool.QueryRow(context.Background(), query, uuid).Scan(&task.ID, &task.Translator, &task.Code, &task.Result, &task.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *Object) SaveTask(uuid uuid.UUID, task *repository.Task) error {
	query := `UPDATE tasks SET translator = $1, code = $2, result = $3, status = $4 WHERE id = $5`
	_, err := r.pool.Exec(context.Background(), query, task.Translator, task.Code, task.Result, task.Status, uuid)
	return err
}

func (r *Object) CreateTask(task *repository.Task) error {
	query := `INSERT INTO tasks (id, translator, code, result, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(context.Background(), query, task.ID, task.Translator, task.Code, task.Result, task.Status)
	return err
}

func (r *Object) DeleteTask(uuid uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.pool.Exec(context.Background(), query, uuid)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return repository.NotFound
	}
	return nil
}

func (r *Object) UpdateStatus(uuid uuid.UUID, s string) error {
	query := `UPDATE tasks SET status = $1 WHERE id = $2`
	result, err := r.pool.Exec(context.Background(), query, s, uuid)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return repository.NotFound
	}
	return nil
}

func (r *Object) UpdateResult(uuid uuid.UUID, s string) error {
	query := `UPDATE tasks SET result = $1 WHERE id = $2`
	_, err := r.pool.Exec(context.Background(), query, s, uuid)
	return err
}
