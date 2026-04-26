package model

import (
	"testing"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantErr   bool
		wantTitle string
	}{
		{
			name:      "нормальное название",
			title:     "Купить хлеб",
			wantErr:   false,
			wantTitle: "Купить хлеб",
		},
		{
			name:      "пустое название - ошибка",
			title:     "",
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "название из пробелов",
			title:     "   ",
			wantErr:   false,
			wantTitle: "   ",
		},
		{
			name:      "очень длинное название",
			title:     "А" + string(make([]byte, 1000)),
			wantErr:   false,
			wantTitle: "А" + string(make([]byte, 1000)),
		},
		{
			name:      "название с символами",
			title:     "Задача №42! @#$%",
			wantErr:   false,
			wantTitle: "Задача №42! @#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.title)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTask() хотел ошибку, но получил nil")
				}
				if err != ErrEmptyName {
					t.Errorf("NewTask() ошибка = %v, хотим %v", err, ErrEmptyName)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTask() не ожидал ошибку, получил %v", err)
			}

			if task == nil {
				t.Errorf("NewTask() вернул nil task")
				return
			}

			if task.Title != tt.wantTitle {
				t.Errorf("Title = %v, хотим %v", task.Title, tt.wantTitle)
			}

			if task.ID == "" {
				t.Errorf("ID не должен быть пустым")
			}

			if task.Created.IsZero() {
				t.Errorf("Created должно быть заполнено")
			}

			if task.Tags == nil {
				t.Errorf("Tags = nil, хотим пустой слайс")
			}
			if task.Subtasks == nil {
				t.Errorf("Subtasks = nil, хотим пустой слайс")
			}
		})
	}
}

func TestNewSubtask(t *testing.T) {
	tests := []struct {
		name         string
		subtaskName  string
		needProgress int
		wantErr      bool
		wantNeed     int
	}{
		{
			name:         "нормальный подзадача",
			subtaskName:  "Написать код",
			needProgress: 10,
			wantErr:      false,
			wantNeed:     10,
		},
		{
			name:         "пустое имя - ошибка",
			subtaskName:  "",
			needProgress: 5,
			wantErr:      true,
			wantNeed:     0,
		},
		{
			name:         "прогресс 0 - ошибка",
			subtaskName:  "Тест",
			needProgress: 0,
			wantErr:      true,
			wantNeed:     0,
		},
		{
			name:         "отрицательный прогресс - ошибка",
			subtaskName:  "Отладка",
			needProgress: -5,
			wantErr:      true,
			wantNeed:     0,
		},
		{
			name:         "минимальный прогресс 1",
			subtaskName:  "Минимум",
			needProgress: 1,
			wantErr:      false,
			wantNeed:     1,
		},
		{
			name:         "большой прогресс",
			subtaskName:  "Большая задача",
			needProgress: 1000,
			wantErr:      false,
			wantNeed:     1000,
		},
		{
			name:         "имя из одного символа",
			subtaskName:  "A",
			needProgress: 5,
			wantErr:      false,
			wantNeed:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subtask, err := NewSubtask(tt.subtaskName, tt.needProgress)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSubtask() хотел ошибку, получил nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewSubtask() не ожидал ошибку, получил %v", err)
			}

			if subtask == nil {
				t.Errorf("NewSubtask() вернул nil")
				return
			}

			if subtask.Name != tt.subtaskName {
				t.Errorf("Name = %v, хотим %v", subtask.Name, tt.subtaskName)
			}

			if subtask.NeedProgress != tt.wantNeed {
				t.Errorf("NeedProgress = %v, хотим %v", subtask.NeedProgress, tt.wantNeed)
			}

			if subtask.ID == "" {
				t.Errorf("ID не должен быть пустым")
			}

			if subtask.Progress != 0 {
				t.Errorf("Progress должен быть 0 по умолчанию, получили %v", subtask.Progress)
			}
		})
	}
}

func TestSubtask_ProgressDefault(t *testing.T) {
	subtask, err := NewSubtask("Тест", 10)
	if err != nil {
		t.Fatalf("Не удалось создать подзадачу: %v", err)
	}

	if subtask.Progress != 0 {
		t.Errorf("Progress = %v, хотим 0", subtask.Progress)
	}
}