package sqlite

import (
    "context"
    "testing"
)

func TestAddTag(t *testing.T) {
    db := setupTestDB(t)
    ctx := context.Background()
    
    task := createTestTask(t, db, "Tag Task")
    
    err := db.AddTag(ctx, task.ID, "work")
    if err != nil {
        t.Fatalf("AddTag() error = %v", err)
    }
    
    tags, err := db.GetTagsByTask(ctx, task.ID)
    if err != nil {
        t.Fatalf("GetTagsByTask() error = %v", err)
    }
    
    if len(tags) != 1 || tags[0] != "work" {
        t.Errorf("Tags = %v, want [work]", tags)
    }
}

func TestAddMultipleTags(t *testing.T) {
    db := setupTestDB(t)
    ctx := context.Background()
    
    task := createTestTask(t, db, "Multi Tag Task")
    
    tagsToAdd := []string{"work", "urgent", "important"}
    for _, tag := range tagsToAdd {
        err := db.AddTag(ctx, task.ID, tag)
        if err != nil {
            t.Fatalf("AddTag(%s) error = %v", tag, err)
        }
    }
    
    tags, err := db.GetTagsByTask(ctx, task.ID)
    if err != nil {
        t.Fatalf("GetTagsByTask() error = %v", err)
    }
    
    if len(tags) != 3 {
        t.Errorf("Expected 3 tags, got %d", len(tags))
    }
}

func TestRemoveTag(t *testing.T) {
    db := setupTestDB(t)
    ctx := context.Background()
    
    task := createTestTask(t, db, "Remove Tag Task")
    
    db.AddTag(ctx, task.ID, "test")
    db.AddTag(ctx, task.ID, "remove")
    
    err := db.RemoveTag(ctx, task.ID, "remove")
    if err != nil {
        t.Fatalf("RemoveTag() error = %v", err)
    }
    
    tags, err := db.GetTagsByTask(ctx, task.ID)
    if err != nil {
        t.Fatalf("GetTagsByTask() error = %v", err)
    }
    
    if len(tags) != 1 || tags[0] != "test" {
        t.Errorf("Tags = %v, want [test]", tags)
    }
}

func TestGetTagsByTaskEmpty(t *testing.T) {
    db := setupTestDB(t)
    ctx := context.Background()
    
    task := createTestTask(t, db, "No Tags Task")
    
    tags, err := db.GetTagsByTask(ctx, task.ID)
    if err != nil {
        t.Fatalf("GetTagsByTask() error = %v", err)
    }
    
    if len(tags) != 0 {
        t.Errorf("Expected 0 tags, got %d", len(tags))
    }
}

func TestTagsCascadeDelete(t *testing.T) {
    db := setupTestDB(t)
    ctx := context.Background()
    
    task := createTestTask(t, db, "Delete Me")
    db.AddTag(ctx, task.ID, "test")
    
    tags, _ := db.GetTagsByTask(ctx, task.ID)
    if len(tags) != 1 {
        t.Fatal("Tag not added")
    }
    
    err := db.DeleteTask(ctx, task.ID)
    if err != nil {
        t.Fatalf("DeleteTask() error = %v", err)
    }
    
    tags, err = db.GetTagsByTask(ctx, task.ID)
    if err != nil {
        t.Fatalf("GetTagsByTask() error = %v", err)
    }
    
    if len(tags) != 0 {
        t.Errorf("Expected 0 tags after task deletion, got %d", len(tags))
    }
}