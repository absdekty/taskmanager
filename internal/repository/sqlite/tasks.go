package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/absdekty/taskmanager/internal/model"
)

func (db *DB) CreateTask(ctx context.Context, task *model.Task) error {
	query := `INSERT INTO tasks (id, title, created, due_time, notify_at)
              VALUES (?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query,
		task.ID, task.Title, task.Created, task.DueTime, task.NotifyAt)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (db *DB) GetTask(ctx context.Context, id string) (*model.Task, error) {
	query := `SELECT id, title, created, due_time, notify_at
              FROM tasks WHERE id = ?`

	var task model.Task
	task.Subtasks = []*model.Subtask{}
	task.Tags = []string{}

	err := db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Created, &task.DueTime, &task.NotifyAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	return &task, nil
}

func (db *DB) UpdateTask(ctx context.Context, task *model.Task) error {
	query := `UPDATE tasks SET title=?, due_time=?, notify_at=? WHERE id=?`

	_, err := db.ExecContext(ctx, query, task.Title, task.DueTime, task.NotifyAt, task.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func (db *DB) DeleteTask(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}
