package model

import (
	"testing"
	"time"
)

/* Задачи */
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

/* Методы задач */
func TestTask_AddTag(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)

	tests := []struct {
		name    string
		setup   func(*Task)
		tag     string
		wantErr error
	}{
		{"успешное добавление", func(t *Task) { t.DueTime = future }, "urgent", nil},
		{"просроченная задача", func(t *Task) { t.DueTime = past }, "urgent", ErrTaskIsOverdue},
		{"пустой тег", func(t *Task) { t.DueTime = future }, "", ErrEmptyName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			tt.setup(task)
			err := task.AddTag(tt.tag)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_RemoveTag(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)

	tests := []struct {
		name    string
		setup   func(*Task)
		tag     string
		wantErr error
	}{
		{"успешное удаление", func(t *Task) { t.DueTime = future; t.Tags = []string{"urgent"} }, "urgent", nil},
		{"тег не найден", func(t *Task) { t.DueTime = future; t.Tags = []string{"bug"} }, "urgent", ErrNotExisting},
		{"просроченная задача", func(t *Task) { t.DueTime = past }, "urgent", ErrTaskIsOverdue},
		{"пустой тег", func(t *Task) { t.DueTime = future }, "", ErrEmptyName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			tt.setup(task)
			err := task.RemoveTag(tt.tag)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_AddSubtask(t *testing.T) {
	tests := []struct {
		name     string
		subName  string
		progress int
		wantErr  error
	}{
		{"успешное добавление", "Подзадача", 5, nil},
		{"пустое имя", "", 5, ErrEmptyName},
		{"прогресс меньше 1", "Подзадача", 0, ErrInvalidProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			err := task.AddSubtask(tt.subName, tt.progress)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_RemoveSubtask(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)

	t.Run("успешное удаление", func(t *testing.T) {
		task, _ := NewTask("Test")
		task.DueTime = future
		task.AddSubtask("test", 5)
		subtaskID := task.Subtasks[0].ID

		err := task.RemoveSubtask(subtaskID)

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
		if len(task.Subtasks) != 0 {
			t.Errorf("expected 0 subtasks, got %d", len(task.Subtasks))
		}
	})

	t.Run("подзадача не найдена", func(t *testing.T) {
		task, _ := NewTask("Test")
		task.DueTime = future

		err := task.RemoveSubtask("invalid-id")

		if err != ErrNotExisting {
			t.Errorf("got %v, want %v", err, ErrNotExisting)
		}
	})

	t.Run("просроченная задача", func(t *testing.T) {
		task, _ := NewTask("Test")
		task.DueTime = past

		err := task.RemoveSubtask("any-id")

		if err != ErrTaskIsOverdue {
			t.Errorf("got %v, want %v", err, ErrTaskIsOverdue)
		}
	})

	t.Run("пустой id", func(t *testing.T) {
		task, _ := NewTask("Test")
		task.DueTime = future

		err := task.RemoveSubtask("")

		if err != ErrEmptyName {
			t.Errorf("got %v, want %v", err, ErrEmptyName)
		}
	})
}

func TestTask_SetDueTime(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)

	tests := []struct {
		name    string
		dueTime time.Time
		wantErr error
	}{
		{"успешная установка", future, nil},
		{"дедлайн в прошлом", past, ErrPastDeadline},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			err := task.SetDueTime(tt.dueTime)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_IsOverdue(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)
	zero := time.Time{}

	tests := []struct {
		name    string
		dueTime time.Time
		want    bool
	}{
		{"не просрочена", future, false},
		{"просрочена", past, true},
		{"не установлена", zero, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			task.DueTime = tt.dueTime
			if got := task.IsOverdue(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_SetNotifyAt(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)
	futureDue := time.Now().UTC().Add(48 * time.Hour)

	tests := []struct {
		name     string
		setup    func(*Task)
		notifyAt time.Time
		wantErr  error
	}{
		{"успешная установка", func(t *Task) { t.DueTime = futureDue }, future, nil},
		{"напоминание в прошлом", func(t *Task) { t.DueTime = futureDue }, past, ErrPastNotify},
		{"задача просрочена", func(t *Task) { t.DueTime = past }, future, ErrTaskIsOverdue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			tt.setup(task)
			err := task.SetNotifyAt(tt.notifyAt)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_IsNotifiable(t *testing.T) {
	future := time.Now().UTC().Add(24 * time.Hour)
	past := time.Now().UTC().Add(-24 * time.Hour)
	zero := time.Time{}

	tests := []struct {
		name     string
		notifyAt time.Time
		want     bool
	}{
		{"можно уведомить", past, true},
		{"еще не время", future, false},
		{"не установлено", zero, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("Test")
			task.NotifyAt = tt.notifyAt
			if got := task.IsNotifiable(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

/* Субзадачи */
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

func TestSubtask_ChangeMaxProgress(t *testing.T) {
	tests := []struct {
		name        string
		maxProgress int
		wantErr     error
	}{
		{"успешное изменение", 10, nil},
		{"невалидный прогресс", 0, ErrInvalidProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, _ := NewSubtask("Test", 5)
			err := sub.ChangeMaxProgress(tt.maxProgress)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if err == nil && sub.NeedProgress != tt.maxProgress {
				t.Errorf("expected %d, got %d", tt.maxProgress, sub.NeedProgress)
			}
		})
	}
}

func TestSubtask_IncrementProgress(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Subtask)
		increment int
		wantErr   error
		expected  int
	}{
		{"успешное увеличение", func(s *Subtask) { s.NeedProgress = 10; s.Progress = 3 }, 2, nil, 5},
		{"превышение максимума", func(s *Subtask) { s.NeedProgress = 10; s.Progress = 9 }, 2, ErrMaxProgressExceeded, 9},
		{"невалидный инкремент", func(s *Subtask) { s.NeedProgress = 10; s.Progress = 5 }, 0, ErrInvalidProgress, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, _ := NewSubtask("Test", 10)
			tt.setup(sub)
			err := sub.IncrementProgress(tt.increment)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if sub.Progress != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, sub.Progress)
			}
		})
	}
}

func TestSubtask_DecrementProgress(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Subtask)
		decrement int
		wantErr   error
		expected  int
	}{
		{"успешное уменьшение", func(s *Subtask) { s.NeedProgress = 10; s.Progress = 5 }, 2, nil, 3},
		{"ниже нуля", func(s *Subtask) { s.NeedProgress = 10; s.Progress = 1 }, 2, ErrMinProgressExceeded, 1},
		{"невалидный декремент", func(s *Subtask) { s.NeedProgress = 10; s.Progress = 5 }, 0, ErrInvalidProgress, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, _ := NewSubtask("Test", 10)
			tt.setup(sub)
			err := sub.DecrementProgress(tt.decrement)
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if sub.Progress != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, sub.Progress)
			}
		})
	}
}
