package sqlite

import (
	"os"
	"testing"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()

	// Создаем временный файл БД
	tmpFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	db, err := NewDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create DB: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpFile.Name())
	})

	return db
}

func TestNewDB(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	db, err := NewDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()

	if db.DB == nil {
		t.Error("DB.DB is nil")
	}
}

func TestNewDBInvalidPath(t *testing.T) {
	_, err := NewDB("/invalid/path/to/db.db")
	if err == nil {
		t.Error("NewDB() with invalid path should return error")
	}
}

func TestClose(t *testing.T) {
	db := setupTestDB(t)
	err := db.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
