package service

import (
	"context"
	"os"
	"testing"

	"github.com/absdekty/taskmanager/internal/model"
	"github.com/absdekty/taskmanager/internal/repository/sqlite"
)

func setupTestService(t *testing.T) (*Service, *sqlite.DB) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	db, err := sqlite.NewDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create DB: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpFile.Name())
	})

	return NewService(db), db
}

func TestCreateTask(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, "Test Task")
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	if task.Title != "Test Task" {
		t.Errorf("Title = %v, want 'Test Task'", task.Title)
	}

	if task.ID == "" {
		t.Error("ID is empty")
	}
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.CreateTask(ctx, "")
	if err == nil {
		t.Error("CreateTask() with empty title should return error")
	}

	if err != model.ErrEmptyName {
		t.Errorf("Expected ErrEmptyName, got %v", err)
	}
}

func TestGetTask(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateTask(ctx, "Get Test")
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	task, err := svc.GetTask(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if task.Title != "Get Test" {
		t.Errorf("Title = %v, want 'Get Test'", task.Title)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.GetTask(ctx, "non-existent-id")
	if err != model.ErrTaskNotFound {
		t.Errorf("Expected ErrTaskNotFound, got %v", err)
	}
}

func TestUpdateTask(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, "Original")
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	task.Title = "Updated"
	err = svc.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("UpdateTask() error = %v", err)
	}

	updated, err := svc.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if updated.Title != "Updated" {
		t.Errorf("Title = %v, want 'Updated'", updated.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, "To Delete")
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	err = svc.DeleteTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	_, err = svc.GetTask(ctx, task.ID)
	if err != model.ErrTaskNotFound {
		t.Error("Task still exists after delete")
	}
}

func TestListTasks(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	titles := []string{"Task 1", "Task 2", "Task 3"}
	for _, title := range titles {
		_, err := svc.CreateTask(ctx, title)
		if err != nil {
			t.Fatalf("CreateTask() error = %v", err)
		}
	}

	tasks, err := svc.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestAddSubtask(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, "Parent")
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	subtasks, err := svc.AddSubtask(ctx, task.ID, "Sub 1", 10)
	if err != nil {
		t.Fatalf("AddSubtask() error = %v", err)
	}

	if len(subtasks) != 1 {
		t.Errorf("Expected 1 subtask, got %d", len(subtasks))
	}

	if subtasks[0].Name != "Sub 1" {
		t.Errorf("Name = %v, want 'Sub 1'", subtasks[0].Name)
	}
}

func TestUpdateSubtaskProgress(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, "Task")
	subtasks, _ := svc.AddSubtask(ctx, task.ID, "Sub", 10)

	err := svc.UpdateSubtaskProgress(ctx, subtasks[0].ID, 7)
	if err != nil {
		t.Fatalf("UpdateSubtaskProgress() error = %v", err)
	}

	updated, err := svc.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if updated.Subtasks[0].Progress != 7 {
		t.Errorf("Progress = %d, want 7", updated.Subtasks[0].Progress)
	}
}

func TestUpdateSubtaskProgressOverLimit(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, "Task")
	subtasks, _ := svc.AddSubtask(ctx, task.ID, "Sub", 10)

	err := svc.UpdateSubtaskProgress(ctx, subtasks[0].ID, 99)
	if err != nil {
		t.Fatalf("UpdateSubtaskProgress() error = %v", err)
	}

	updated, _ := svc.GetTask(ctx, task.ID)
	if updated.Subtasks[0].Progress != 10 {
		t.Errorf("Progress = %d, want 10 (capped)", updated.Subtasks[0].Progress)
	}
}

func TestDeleteSubtask(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, "Task")
	subtasks, _ := svc.AddSubtask(ctx, task.ID, "To Delete", 5)

	err := svc.DeleteSubtask(ctx, subtasks[0].ID)
	if err != nil {
		t.Fatalf("DeleteSubtask() error = %v", err)
	}

	updated, _ := svc.GetTask(ctx, task.ID)
	if len(updated.Subtasks) != 0 {
		t.Error("Subtask still exists after delete")
	}
}

func TestAddTag(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, "Tagged Task")

	err := svc.AddTag(ctx, task.ID, "work")
	if err != nil {
		t.Fatalf("AddTag() error = %v", err)
	}

	updated, _ := svc.GetTask(ctx, task.ID)
	if len(updated.Tags) != 1 || updated.Tags[0] != "work" {
		t.Errorf("Tags = %v, want [work]", updated.Tags)
	}
}

func TestRemoveTag(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, "Task")
	svc.AddTag(ctx, task.ID, "test")
	svc.AddTag(ctx, task.ID, "remove")

	err := svc.RemoveTag(ctx, task.ID, "remove")
	if err != nil {
		t.Fatalf("RemoveTag() error = %v", err)
	}

	updated, _ := svc.GetTask(ctx, task.ID)
	if len(updated.Tags) != 1 || updated.Tags[0] != "test" {
		t.Errorf("Tags = %v, want [test]", updated.Tags)
	}
}

func TestMultipleTags(t *testing.T) {
	svc, _ := setupTestService(t)
	ctx := context.Background()

	task, _ := svc.CreateTask(ctx, "Multi Tag")

	tags := []string{"work", "urgent", "important"}
	for _, tag := range tags {
		svc.AddTag(ctx, task.ID, tag)
	}

	updated, _ := svc.GetTask(ctx, task.ID)
	if len(updated.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(updated.Tags))
	}
}
