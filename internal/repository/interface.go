package repository

import (
	"context"
	"github.com/absdekty/taskmanager/internal/model"
)

/* Общее */
type TaskMainRepositoryI interface {
	ListTasks(ctx context.Context) ([]*model.Task, error)
}

/* Задача */
type TaskRepositoryI interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTask(ctx context.Context, id string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, id string) error
}

/* Субзадача */
type SubtaskRepositoryI interface {
	CreateSubtask(ctx context.Context, subtask *model.Subtask) error
	GetSubtasksByTask(ctx context.Context, taskID string) ([]*model.Subtask, error)
	UpdateSubtask(ctx context.Context, subtask *model.Subtask) error
	DeleteSubtask(ctx context.Context, id string) error
}

type TagRepositoryI interface {
	AddTag(ctx context.Context, taskID, tag string) error
	RemoveTag(ctx context.Context, taskID, tag string) error
	GetTagsByTask(ctx context.Context, taskID string) ([]string, error)
}

/* Основной интерфейс, реализующий CRUD-методы для tasks, subtasks */
type RepositoryI interface {
	TaskMainRepositoryI
	TaskRepositoryI
	SubtaskRepositoryI
	TagRepositoryI
	Close() error
}
