package sqlite

import (
	"context"
	"fmt"

	"github.com/absdekty/taskmanager/internal/model"
)

func (db *DB) ListTasks(ctx context.Context) ([]*model.Task, error) {
	query := `SELECT id, title, created, due_time, notify_at FROM tasks ORDER BY created DESC`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Created, &task.DueTime, &task.NotifyAt)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
		if err != nil {
			return nil, fmt.Errorf("get subtasks for task %s: %w", task.ID, err)
		}
		task.Subtasks = subtasks

		task.Tags = []string{}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return tasks, nil
}
