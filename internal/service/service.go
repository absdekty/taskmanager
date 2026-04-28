package service

import (
	"context"
	"github.com/absdekty/taskmanager/internal/model"
	"github.com/absdekty/taskmanager/internal/repository"
)

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) *Service {
	return &Service{repo: repo}
}

var _ ServiceI = (*Service)(nil)

func (s *Service) CreateTask(ctx context.Context, title string) (*model.Task, error) {
	task, err := model.NewTask(title)
	if err != nil {
		return nil, err
	}

	if err = s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) GetTask(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, model.ErrTaskNotFound
	}

	subtasks, err := s.repo.GetSubtasksByTask(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Subtasks = subtasks

	tags, err := s.repo.GetTagsByTask(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Tags = tags

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, task *model.Task) error {
	return s.repo.UpdateTask(ctx, task)
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	return s.repo.DeleteTask(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context) ([]*model.Task, error) {
	return s.repo.ListTasks(ctx)
}

func (s *Service) AddSubtask(ctx context.Context, taskID, name string, needProgress int) ([]*model.Subtask, error) {
	subtask, err := model.NewSubtask(name, needProgress)
	if err != nil {
		return nil, err
	}

	subtask.TaskID = taskID
	if err = s.repo.CreateSubtask(ctx, subtask); err != nil {
		return nil, err
	}

	return s.repo.GetSubtasksByTask(ctx, taskID)
}

func (s *Service) UpdateSubtaskProgress(ctx context.Context, subtaskID string, progress int) error {
	subtask, err := s.repo.GetSubtask(ctx, subtaskID)
	if err != nil {
		return err
	}
	if subtask == nil {
		return model.ErrSubtaskNotFound
	}

	if progress > subtask.NeedProgress {
		progress = subtask.NeedProgress
	}
	if progress < 0 {
		progress = 0
	}

	subtask.Progress = progress
	return s.repo.UpdateSubtask(ctx, subtask)
}

func (s *Service) DeleteSubtask(ctx context.Context, subtaskID string) error {
	return s.repo.DeleteSubtask(ctx, subtaskID)
}

func (s *Service) AddTag(ctx context.Context, taskID, tag string) error {
	return s.repo.AddTag(ctx, taskID, tag)
}

func (s *Service) RemoveTag(ctx context.Context, taskID, tag string) error {
	return s.repo.RemoveTag(ctx, taskID, tag)
}
