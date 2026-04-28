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
            return nil, fmt.Errorf("get subtasks: %w", err)
        }
        task.Subtasks = subtasks
        
        tags, err := db.GetTagsByTask(ctx, task.ID)
        if err != nil {
            return nil, fmt.Errorf("get tags: %w", err)
        }
        task.Tags = tags
        
        tasks = append(tasks, &task)
    }
    
    return tasks, rows.Err()
}