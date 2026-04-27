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

/*
	Set DueTime
	Set NotifyAt
	
	Get Overdue
	Get Notifiable
*/


var (
	/* Общие */
	ErrEmptyName       = errors.New("Название пустое")
	
	/* Задача */
	ErrTaskIsOverdue = errors.New("Задача в дедлайне")
	ErrPastDeadline = errors.New("Дедлайн в прошлом")
	ErrPastNotify = errors.New("Напоминание в прошлом")
	ErrNotExisting = errors.New("Тег не существует")
	
	/* Субзадача */
	ErrInvalidProgress = errors.New("Прогресс  меньше единицы")
	ErrMaxProgressExceeded = errors.New("Прогресс выше максимально")
	ErrMinProgressExceeded = errors.New("Прогресс отрицателен")
)

/* ЗАДАЧИ */
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

func (t *Task) AddTag(tag string) error {
	if t.IsOverdue() {
		return ErrTaskIsOverdue
	}
	
	if tag == "" {
		return ErrEmptyName
	}
	
	t.Tags = append(t.Tags, tag)
	
	return nil
}

func (t *Task) RemoveTag(tag string) error {
	if t.IsOverdue() {
		return ErrTaskIsOverdue
	}
	
	if tag == "" {
		return ErrEmptyName
	}
	
	for i, val := range t.Tags {
		if tag == val {
			t.Tags = append(t.Tags[:i], t.Tags[i+1:]...)
			return nil
		}
	}
	
	return ErrNotExisting
}

func (t *Task) AddSubtask(name string, progress int) error {
	subtask, err := NewSubtask(name, progress)
	if err != nil {
		return err
	}
	
	t.Subtasks = append(t.Subtasks, subtask)
	return nil
}

func (t *Task) RemoveSubtask(id string) error {
	if t.IsOverdue() {
		return ErrTaskIsOverdue
	}
	
	if id == "" {
		return ErrEmptyName
	}
	
	for i, val := range t.Subtasks {
		if id == val.ID {
			t.Subtasks = append(t.Subtasks[:i], t.Subtasks[i+1:]...)
			return nil
		}
	}
	
	return ErrNotExisting
}

func (t *Task) SetDueTime(duetime time.Time) error {
	if duetime.Before(time.Now().UTC()) {
		return ErrPastDeadline
	}
	
	t.DueTime = duetime
	
	return nil
}

func (t *Task) IsOverdue() bool {
	return (!t.DueTime.IsZero() && t.DueTime.Before(time.Now().UTC()))
}

func (t *Task) SetNotifyAt(notifyat time.Time) error {
	if t.IsOverdue() {
		return ErrTaskIsOverdue
	}

	if notifyat.Before(time.Now().UTC()) {
		return ErrPastNotify
	}
	
	t.NotifyAt = notifyat
	
	return nil
}

func (t *Task) IsNotifiable() bool {
	return (!t.NotifyAt.IsZero() && t.NotifyAt.Before(time.Now().UTC()))
}

/* СУБЗАДАЧИ */
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

func (t *Subtask) ChangeMaxProgress(MaxProgress int) error {
	if MaxProgress < 1 {
		return ErrInvalidProgress
	}
	
	t.NeedProgress = MaxProgress
	return nil
}

func (t *Subtask) IncrementProgress(progress int) error {
	if progress < 1 {
		return ErrInvalidProgress
	}

	if t.Progress + progress > t.NeedProgress {
		return ErrMaxProgressExceeded
	}
	
	t.Progress += progress
	return nil
}

func (t *Subtask) DecrementProgress(progress int) error {
	if progress < 1 {
		return ErrInvalidProgress
	}

	if t.Progress - progress < 0 {
		return ErrMinProgressExceeded
	}
	
	t.Progress -= progress
	return nil
}