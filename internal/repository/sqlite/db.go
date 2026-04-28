package sqlite

import (
	"database/sql"
	"github.com/absdekty/taskmanager/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

var _ repository.RepositoryI = (*DB)(nil)

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func initSchema(db *sql.DB) error {
	queries := []string{
		// Таблица задач
		`CREATE TABLE IF NOT EXISTS tasks (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            created DATETIME NOT NULL,
            due_time DATETIME,
            notify_at DATETIME
        )`,

		// Таблица подзадач
		`CREATE TABLE IF NOT EXISTS subtasks (
            id TEXT PRIMARY KEY,
            task_id TEXT NOT NULL,
            name TEXT NOT NULL,
            need_progress INTEGER NOT NULL,
            progress INTEGER NOT NULL DEFAULT 0,
            FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
        )`,

		`CREATE TABLE IF NOT EXISTS tags (
			task_id TEXT NOT NULL,
			tag TEXT NOT NULL,
			PRIMARY KEY (task_id, tag),
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
