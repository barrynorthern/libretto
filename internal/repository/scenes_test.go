package repository

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/barrynorthern/libretto/internal/db"
)

func setupTestDB(t *testing.T) *db.Database {
	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	database, err := db.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return database
}

func TestScenesRepository_CreateAndList(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewScenesRepository(database)
	ctx := context.Background()

	// Test creating a scene
	scene, err := repo.CreateScene(ctx, "Test Scene", "A test scene summary", "This is the content of the test scene.")
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	if scene.ID == "" {
		t.Error("Expected scene ID to be generated")
	}
	if scene.Title != "Test Scene" {
		t.Errorf("Expected title 'Test Scene', got '%s'", scene.Title)
	}
	if scene.Summary != "A test scene summary" {
		t.Errorf("Expected summary 'A test scene summary', got '%s'", scene.Summary)
	}

	// Test listing scenes
	scenes, err := repo.ListScenes(ctx)
	if err != nil {
		t.Fatalf("Failed to list scenes: %v", err)
	}

	if len(scenes) != 1 {
		t.Errorf("Expected 1 scene, got %d", len(scenes))
	}

	if scenes[0].ID != scene.ID {
		t.Errorf("Expected scene ID %s, got %s", scene.ID, scenes[0].ID)
	}
}

func TestScenesRepository_GetScene(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewScenesRepository(database)
	ctx := context.Background()

	// Create a scene
	created, err := repo.CreateScene(ctx, "Get Test", "Summary", "Content")
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// Get the scene
	retrieved, err := repo.GetScene(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to get scene: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, retrieved.ID)
	}
	if retrieved.Title != created.Title {
		t.Errorf("Expected title %s, got %s", created.Title, retrieved.Title)
	}
}

func TestScenesRepository_UpdateScene(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewScenesRepository(database)
	ctx := context.Background()

	// Create a scene
	created, err := repo.CreateScene(ctx, "Original Title", "Original Summary", "Original Content")
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// Update the scene
	updated, err := repo.UpdateScene(ctx, created.ID, "Updated Title", "Updated Summary", "Updated Content")
	if err != nil {
		t.Fatalf("Failed to update scene: %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
	}
	if updated.Summary != "Updated Summary" {
		t.Errorf("Expected summary 'Updated Summary', got '%s'", updated.Summary)
	}
	if updated.Content != "Updated Content" {
		t.Errorf("Expected content 'Updated Content', got '%s'", updated.Content)
	}

	// Verify the update persisted
	retrieved, err := repo.GetScene(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to get updated scene: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("Expected persisted title 'Updated Title', got '%s'", retrieved.Title)
	}
}

func TestScenesRepository_DeleteScene(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewScenesRepository(database)
	ctx := context.Background()

	// Create a scene
	created, err := repo.CreateScene(ctx, "To Delete", "Summary", "Content")
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// Delete the scene
	err = repo.DeleteScene(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to delete scene: %v", err)
	}

	// Verify it's gone
	scenes, err := repo.ListScenes(ctx)
	if err != nil {
		t.Fatalf("Failed to list scenes: %v", err)
	}

	if len(scenes) != 0 {
		t.Errorf("Expected 0 scenes after deletion, got %d", len(scenes))
	}
}
