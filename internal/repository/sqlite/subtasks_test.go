package sqlite

import (
	"context"
	"testing"

	"github.com/absdekty/taskmanager/internal/model"
	"github.com/google/uuid"
)

func createTestSubtask(t *testing.T, db *DB, taskID string, name string, needProgress int) *model.Subtask {
	t.Helper()

	subtask := &model.Subtask{
		ID:           uuid.New().String(),
		TaskID:       taskID,
		Name:         name,
		NeedProgress: needProgress,
		Progress:     0,
	}

	err := db.CreateSubtask(context.Background(), subtask)
	if err != nil {
		t.Fatalf("Failed to create test subtask: %v", err)
	}

	return subtask
}

func TestCreateSubtask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Parent Task")

	subtask := &model.Subtask{
		ID:           uuid.New().String(),
		TaskID:       task.ID,
		Name:         "Test Subtask",
		NeedProgress: 10,
		Progress:     0,
	}

	err := db.CreateSubtask(ctx, subtask)
	if err != nil {
		t.Fatalf("CreateSubtask() error = %v", err)
	}

	// Verify subtask was created
	subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}

	if len(subtasks) != 1 {
		t.Fatalf("Expected 1 subtask, got %d", len(subtasks))
	}

	if subtasks[0].Name != "Test Subtask" {
		t.Errorf("Name = %v, want 'Test Subtask'", subtasks[0].Name)
	}
	if subtasks[0].NeedProgress != 10 {
		t.Errorf("NeedProgress = %v, want 10", subtasks[0].NeedProgress)
	}
	if subtasks[0].Progress != 0 {
		t.Errorf("Progress = %v, want 0", subtasks[0].Progress)
	}
}

func TestGetSubtasksByTask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Task with Subtasks")

	// Create multiple subtasks
	names := []string{"Sub 1", "Sub 2", "Sub 3"}
	needs := []int{5, 10, 15}

	for i := 0; i < 3; i++ {
		createTestSubtask(t, db, task.ID, names[i], needs[i])
	}

	subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}

	if len(subtasks) != 3 {
		t.Fatalf("Expected 3 subtasks, got %d", len(subtasks))
	}

	for i, subtask := range subtasks {
		if subtask.Name != names[i] {
			t.Errorf("Subtask %d: Name = %v, want %v", i, subtask.Name, names[i])
		}
		if subtask.NeedProgress != needs[i] {
			t.Errorf("Subtask %d: NeedProgress = %v, want %v", i, subtask.NeedProgress, needs[i])
		}
	}
}

func TestGetSubtask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Parent Task")
	created := createTestSubtask(t, db, task.ID, "Test Subtask", 10)

	subtask, err := db.GetSubtask(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSubtask() error = %v", err)
	}

	if subtask == nil {
		t.Fatal("GetSubtask() returned nil")
	}

	if subtask.ID != created.ID {
		t.Errorf("ID = %v, want %v", subtask.ID, created.ID)
	}
	if subtask.TaskID != task.ID {
		t.Errorf("TaskID = %v, want %v", subtask.TaskID, task.ID)
	}
	if subtask.Name != "Test Subtask" {
		t.Errorf("Name = %v, want 'Test Subtask'", subtask.Name)
	}
	if subtask.NeedProgress != 10 {
		t.Errorf("NeedProgress = %v, want 10", subtask.NeedProgress)
	}
	if subtask.Progress != 0 {
		t.Errorf("Progress = %v, want 0", subtask.Progress)
	}
}

func TestGetSubtaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	subtask, err := db.GetSubtask(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("GetSubtask() error = %v", err)
	}

	if subtask != nil {
		t.Errorf("Expected nil, got %v", subtask)
	}
}

func TestGetSubtaskAfterUpdate(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Parent Task")
	subtask := createTestSubtask(t, db, task.ID, "Original", 10)

	// Обновляем прогресс
	subtask.Progress = 7
	err := db.UpdateSubtask(ctx, subtask)
	if err != nil {
		t.Fatalf("UpdateSubtask() error = %v", err)
	}

	// Получаем и проверяем
	updated, err := db.GetSubtask(ctx, subtask.ID)
	if err != nil {
		t.Fatalf("GetSubtask() error = %v", err)
	}

	if updated.Progress != 7 {
		t.Errorf("Progress = %v, want 7", updated.Progress)
	}
}

func TestGetSubtasksByTaskEmpty(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Task without Subtasks")

	subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}

	if len(subtasks) != 0 {
		t.Errorf("Expected 0 subtasks, got %d", len(subtasks))
	}
}

func TestUpdateSubtask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Parent Task")
	subtask := createTestSubtask(t, db, task.ID, "Original", 10)

	// Update subtask
	subtask.Name = "Updated"
	subtask.NeedProgress = 20
	subtask.Progress = 15

	err := db.UpdateSubtask(ctx, subtask)
	if err != nil {
		t.Fatalf("UpdateSubtask() error = %v", err)
	}

	// Verify update
	subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}

	if len(subtasks) != 1 {
		t.Fatalf("Expected 1 subtask, got %d", len(subtasks))
	}

	updated := subtasks[0]
	if updated.Name != "Updated" {
		t.Errorf("Name = %v, want 'Updated'", updated.Name)
	}
	if updated.NeedProgress != 20 {
		t.Errorf("NeedProgress = %v, want 20", updated.NeedProgress)
	}
	if updated.Progress != 15 {
		t.Errorf("Progress = %v, want 15", updated.Progress)
	}
}

func TestDeleteSubtask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Parent Task")
	subtask := createTestSubtask(t, db, task.ID, "To Delete", 5)

	// Verify subtask exists
	subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}
	if len(subtasks) != 1 {
		t.Fatalf("Expected 1 subtask before delete, got %d", len(subtasks))
	}

	// Delete subtask
	err = db.DeleteSubtask(ctx, subtask.ID)
	if err != nil {
		t.Fatalf("DeleteSubtask() error = %v", err)
	}

	// Verify deletion
	subtasks, err = db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}

	if len(subtasks) != 0 {
		t.Errorf("Expected 0 subtasks after delete, got %d", len(subtasks))
	}

	// Verify parent task still exists
	parentTask, err := db.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if parentTask == nil {
		t.Error("Parent task was deleted when subtask was deleted")
	}
}

func TestCascadeDeleteSubtasks(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Task to Delete")

	// Create subtasks
	createTestSubtask(t, db, task.ID, "Subtask 1", 5)
	createTestSubtask(t, db, task.ID, "Subtask 2", 10)

	// Verify subtasks exist
	subtasks, err := db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}
	if len(subtasks) != 2 {
		t.Fatalf("Expected 2 subtasks, got %d", len(subtasks))
	}

	// Delete task (should cascade delete subtasks)
	err = db.DeleteTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	// Verify subtasks are gone
	subtasks, err = db.GetSubtasksByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetSubtasksByTask() error = %v", err)
	}

	if len(subtasks) != 0 {
		t.Errorf("Expected 0 subtasks after task deletion, got %d", len(subtasks))
	}
}
