package sqlite

import (
	"context"
	"fmt"
)

func (db *DB) AddTag(ctx context.Context, taskID, tag string) error {
	query := `INSERT INTO tags (task_id, tag) VALUES (?, ?)`
	_, err := db.ExecContext(ctx, query, taskID, tag)
	if err != nil {
		return fmt.Errorf("add tag: %w", err)
	}
	return nil
}

func (db *DB) RemoveTag(ctx context.Context, taskID, tag string) error {
	query := `DELETE FROM tags WHERE task_id = ? AND tag = ?`
	_, err := db.ExecContext(ctx, query, taskID, tag)
	if err != nil {
		return fmt.Errorf("remove tag: %w", err)
	}
	return nil
}

func (db *DB) GetTagsByTask(ctx context.Context, taskID string) ([]string, error) {
	query := `SELECT tag FROM tags WHERE task_id = ? ORDER BY tag`

	rows, err := db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}
