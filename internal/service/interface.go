package service

import (
	"context"
	"github.com/absdekty/taskmanager/internal/model"
)

type TaskServiceI interface {
	CreateTask(ctx context.Context, title string) (*model.Task, error)
	GetTask(ctx context.Context, id string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context) ([]*model.Task, error)
}

type SubtaskServiceI interface {
	AddSubtask(ctx context.Context, taskID, name string, needProgress int) ([]*model.Subtask, error)
	UpdateSubtaskProgress(ctx context.Context, subtaskID string, progress int) error
	DeleteSubtask(ctx context.Context, subtaskID string) error
}

type TagServiceI interface {
	AddTag(ctx context.Context, taskID, tag string) error
	RemoveTag(ctx context.Context, taskID, tag string) error
}

type ServiceI interface {
	TaskServiceI
	SubtaskServiceI
	TagServiceI
}
