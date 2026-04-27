package sqlite

import (
	"context"
	"testing"
)

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	titles := []string{"Task 1", "Task 2", "Task 3"}
	for _, title := range titles {
		createTestTask(t, db, title)
	}

	tasks, err := db.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != len(titles) {
		t.Errorf("ListTasks() returned %d tasks, want %d", len(tasks), len(titles))
	}
}
