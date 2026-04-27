package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/absdekty/taskmanager/internal/model"
	"github.com/google/uuid"
)

func createTestTask(t *testing.T, db *DB, title string) *model.Task {
	t.Helper()

	task := &model.Task{
		ID:      uuid.New().String(),
		Title:   title,
		Created: time.Now(),
	}

	err := db.CreateTask(context.Background(), task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	return task
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := &model.Task{
		ID:       uuid.New().String(),
		Title:    "Test Task",
		Created:  time.Now(),
		DueTime:  time.Now().Add(24 * time.Hour),
		NotifyAt: time.Now().Add(1 * time.Hour),
	}

	err := db.CreateTask(ctx, task)
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	// Verify task was created
	saved, err := db.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if saved.Title != task.Title {
		t.Errorf("Title = %v, want %v", saved.Title, task.Title)
	}
}

func TestGetTask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	created := createTestTask(t, db, "Get Task Test")

	task, err := db.GetTask(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if task.ID != created.ID {
		t.Errorf("ID = %v, want %v", task.ID, created.ID)
	}
	if task.Title != "Get Task Test" {
		t.Errorf("Title = %v, want 'Get Task Test'", task.Title)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task, err := db.GetTask(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if task != nil {
		t.Errorf("Expected nil, got %v", task)
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Original Title")

	task.Title = "Updated Title"
	task.DueTime = time.Now().Add(48 * time.Hour)

	err := db.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("UpdateTask() error = %v", err)
	}

	updated, err := db.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("Title = %v, want 'Updated Title'", updated.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	task := createTestTask(t, db, "Task to Delete")

	err := db.DeleteTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	deleted, err := db.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if deleted != nil {
		t.Error("Task was not deleted")
	}
}
