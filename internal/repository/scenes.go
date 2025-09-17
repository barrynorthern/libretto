package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/google/uuid"
)

// ScenesRepository provides database operations for scenes
type ScenesRepository struct {
	db *db.Database
}

// NewScenesRepository creates a new scenes repository
func NewScenesRepository(database *db.Database) *ScenesRepository {
	return &ScenesRepository{db: database}
}

// Scene represents a scene in the domain layer
type Scene struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateScene creates a new scene with auto-generated ID
func (r *ScenesRepository) CreateScene(ctx context.Context, title, summary, content string) (*Scene, error) {
	id := uuid.New().String()
	
	dbScene, err := r.db.Queries().CreateScene(ctx, db.CreateSceneParams{
		ID:      id,
		Title:   title,
		Summary: summary,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create scene: %w", err)
	}

	return &Scene{
		ID:        dbScene.ID,
		Title:     dbScene.Title,
		Summary:   dbScene.Summary,
		Content:   dbScene.Content,
		CreatedAt: dbScene.CreatedAt,
		UpdatedAt: dbScene.UpdatedAt,
	}, nil
}

// CreateSceneWithID creates a new scene with a specific ID
func (r *ScenesRepository) CreateSceneWithID(ctx context.Context, id, title, summary, content string) (*Scene, error) {
	dbScene, err := r.db.Queries().CreateScene(ctx, db.CreateSceneParams{
		ID:      id,
		Title:   title,
		Summary: summary,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create scene with ID %s: %w", id, err)
	}

	return &Scene{
		ID:        dbScene.ID,
		Title:     dbScene.Title,
		Summary:   dbScene.Summary,
		Content:   dbScene.Content,
		CreatedAt: dbScene.CreatedAt,
		UpdatedAt: dbScene.UpdatedAt,
	}, nil
}

// GetScene retrieves a scene by ID
func (r *ScenesRepository) GetScene(ctx context.Context, id string) (*Scene, error) {
	dbScene, err := r.db.Queries().GetScene(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scene %s: %w", id, err)
	}

	return &Scene{
		ID:        dbScene.ID,
		Title:     dbScene.Title,
		Summary:   dbScene.Summary,
		Content:   dbScene.Content,
		CreatedAt: dbScene.CreatedAt,
		UpdatedAt: dbScene.UpdatedAt,
	}, nil
}

// ListScenes retrieves all scenes ordered by creation time (newest first)
func (r *ScenesRepository) ListScenes(ctx context.Context) ([]*Scene, error) {
	dbScenes, err := r.db.Queries().ListScenes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list scenes: %w", err)
	}

	scenes := make([]*Scene, len(dbScenes))
	for i, dbScene := range dbScenes {
		scenes[i] = &Scene{
			ID:        dbScene.ID,
			Title:     dbScene.Title,
			Summary:   dbScene.Summary,
			Content:   dbScene.Content,
			CreatedAt: dbScene.CreatedAt,
			UpdatedAt: dbScene.UpdatedAt,
		}
	}

	return scenes, nil
}

// UpdateScene updates an existing scene
func (r *ScenesRepository) UpdateScene(ctx context.Context, id, title, summary, content string) (*Scene, error) {
	dbScene, err := r.db.Queries().UpdateScene(ctx, db.UpdateSceneParams{
		ID:      id,
		Title:   title,
		Summary: summary,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update scene %s: %w", id, err)
	}

	return &Scene{
		ID:        dbScene.ID,
		Title:     dbScene.Title,
		Summary:   dbScene.Summary,
		Content:   dbScene.Content,
		CreatedAt: dbScene.CreatedAt,
		UpdatedAt: dbScene.UpdatedAt,
	}, nil
}

// DeleteScene deletes a scene by ID
func (r *ScenesRepository) DeleteScene(ctx context.Context, id string) error {
	err := r.db.Queries().DeleteScene(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete scene %s: %w", id, err)
	}
	return nil
}
