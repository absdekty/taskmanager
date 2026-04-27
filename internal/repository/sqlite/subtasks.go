package sqlite

import (
	"context"
	"fmt"
	"github.com/absdekty/taskmanager/internal/model"
)

func (db *DB) CreateSubtask(ctx context.Context, subtask *model.Subtask) error {
	query := `INSERT INTO subtasks (id, task_id, name, need_progress, progress)
              VALUES (?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query,
		subtask.ID, subtask.TaskID, subtask.Name, subtask.NeedProgress, subtask.Progress)
	if err != nil {
		return fmt.Errorf("create subtask: %w", err)
	}
	return nil
}

func (db *DB) GetSubtasksByTask(ctx context.Context, taskID string) ([]*model.Subtask, error) {
	query := `SELECT id, task_id, name, need_progress, progress
              FROM subtasks WHERE task_id = ?`

	rows, err := db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []*model.Subtask
	for rows.Next() {
		var subtask model.Subtask
		err := rows.Scan(&subtask.ID, &subtask.TaskID, &subtask.Name, &subtask.NeedProgress, &subtask.Progress)
		if err != nil {
			return nil, err
		}
		subtasks = append(subtasks, &subtask)
	}

	return subtasks, rows.Err()
}

func (db *DB) UpdateSubtask(ctx context.Context, subtask *model.Subtask) error {
	query := `UPDATE subtasks SET name=?, need_progress=?, progress=? WHERE id=?`
	_, err := db.ExecContext(ctx, query, subtask.Name, subtask.NeedProgress, subtask.Progress, subtask.ID)
	return err
}

func (db *DB) DeleteSubtask(ctx context.Context, id string) error {
	query := `DELETE FROM subtasks WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}
