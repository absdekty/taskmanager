package model

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Subtask struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	NeedProgress int    `json:"need"`
	Progress     int    `json:"progress"`
}

type Task struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Tags     []string   `json:"tags"`
	Subtasks []*Subtask `json:"subtasks"`
	Created  time.Time  `json:"created"`
	DueTime  time.Time  `json:"overtime"`
	NotifyAt time.Time  `json:"notify"`
}

var (
	ErrEmptyName       = errors.New("Название не может быть пустым")
	ErrInvalidProgress = errors.New("Прогресс не может быть меньше единицы")
)

func NewTask(title string) (*Task, error) {
	if title == "" {
		return nil, ErrEmptyName
	}

	return &Task{
		ID:       uuid.New().String(),
		Title:    title,
		Tags:     []string{},
		Subtasks: []*Subtask{},
		Created:  time.Now().UTC(),
		DueTime:  time.Time{},
		NotifyAt: time.Time{},
	}, nil
}

func NewSubtask(name string, progress int) (*Subtask, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	if progress < 1 {
		return nil, ErrInvalidProgress
	}

	return &Subtask{
		ID:           uuid.New().String(),
		Name:         name,
		NeedProgress: progress,
	}, nil
}
