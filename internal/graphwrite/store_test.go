package graphwrite

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) *db.Database {
	// Create temporary database file
	tmpFile, err := os.CreateTemp("", "libretto_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	// Clean up after test
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	// Create database connection
	database, err := db.NewDatabase(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Run migrations
	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return database
}

func createTestProject(t *testing.T, database *db.Database) string {
	ctx := context.Background()
	projectID := uuid.New().String()

	_, err := database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Test project description", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	return projectID
}

func createTestGraphVersion(t *testing.T, database *db.Database, projectID string, isWorkingSet bool) string {
	ctx := context.Background()
	versionID := uuid.New().String()

	_, err := database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           versionID,
		ProjectID:    projectID,
		Name:         sql.NullString{String: "Test Version", Valid: true},
		Description:  sql.NullString{String: "Test version description", Valid: true},
		IsWorkingSet: isWorkingSet,
	})
	if err != nil {
		t.Fatalf("Failed to create test graph version: %v", err)
	}

	return versionID
}

func TestService_Apply_CreateEntity(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Test creating a new entity
	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   uuid.New().String(),
				Fields: map[string]any{
					"name":    "Opening Scene",
					"title":   "The Beginning",
					"summary": "Our hero starts their journey",
					"content": "It was a dark and stormy night...",
				},
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Applied != 1 {
		t.Errorf("Expected 1 delta applied, got %d", response.Applied)
	}

	// Verify the entity was created in the new version
	entities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entities))
	}

	entity := entities[0]
	if entity.EntityType != "Scene" {
		t.Errorf("Expected entity type 'Scene', got '%s'", entity.EntityType)
	}

	if entity.Name != "Opening Scene" {
		t.Errorf("Expected entity name 'Opening Scene', got '%s'", entity.Name)
	}

	// Verify entity data
	if title, ok := entity.Data["title"].(string); !ok || title != "The Beginning" {
		t.Errorf("Expected title 'The Beginning', got %v", entity.Data["title"])
	}
}

func TestService_Apply_UpdateEntity(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Create initial entity
	entityID := uuid.New().String()
	entityData := map[string]any{
		"name":    "Original Scene",
		"title":   "Original Title",
		"summary": "Original summary",
	}
	dataBytes, _ := json.Marshal(entityData)

	_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         entityID,
		VersionID:  parentVersionID,
		EntityType: "Scene",
		Name:       "Original Scene",
		Data:       dataBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create initial entity: %v", err)
	}

	// Test updating the entity
	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "update",
				EntityType: "Scene",
				EntityID:   entityID,
				Fields: map[string]any{
					"name":    "Updated Scene",
					"title":   "Updated Title",
					"summary": "Updated summary",
				},
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Applied != 1 {
		t.Errorf("Expected 1 delta applied, got %d", response.Applied)
	}

	// Verify the entity was updated in the new version
	entities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entities))
	}

	entity := entities[0]
	if entity.Name != "Updated Scene" {
		t.Errorf("Expected entity name 'Updated Scene', got '%s'", entity.Name)
	}

	if title, ok := entity.Data["title"].(string); !ok || title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %v", entity.Data["title"])
	}
}

func TestService_Apply_DeleteEntity(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Create initial entity
	entityID := uuid.New().String()
	entityData := map[string]any{
		"name":  "Scene to Delete",
		"title": "Doomed Scene",
	}
	dataBytes, _ := json.Marshal(entityData)

	_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         entityID,
		VersionID:  parentVersionID,
		EntityType: "Scene",
		Name:       "Scene to Delete",
		Data:       dataBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create initial entity: %v", err)
	}

	// Test deleting the entity
	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "delete",
				EntityType: "Scene",
				EntityID:   entityID,
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Applied != 1 {
		t.Errorf("Expected 1 delta applied, got %d", response.Applied)
	}

	// Verify the entity was deleted in the new version
	entities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(entities) != 0 {
		t.Errorf("Expected 0 entities after deletion, got %d", len(entities))
	}
}

func TestService_Apply_WithRelationships(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Create entities with relationships
	sceneID := uuid.New().String()
	characterID := uuid.New().String()
	relationshipID := uuid.New().String()

	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   sceneID,
				Fields: map[string]any{
					"name":  "Scene with Character",
					"title": "Character Introduction",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   characterID,
				Fields: map[string]any{
					"name": "Hero",
					"role": "protagonist",
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						RelationshipID:   relationshipID,
						FromEntityID:     sceneID,
						ToEntityID:       characterID,
						RelationshipType: "features",
						Properties: map[string]any{
							"importance": "primary",
						},
					},
				},
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Applied != 2 {
		t.Errorf("Expected 2 deltas applied, got %d", response.Applied)
	}

	// Verify entities were created
	entities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(entities))
	}

	// Verify relationship was created using version-aware method
	neighbors, err := service.GetNeighborsInVersion(ctx, response.GraphVersionID, sceneID, "features")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed: %v", err)
	}

	if len(neighbors) != 1 {
		t.Errorf("Expected 1 neighbor, got %d", len(neighbors))
	}

	if neighbors[0].ID != characterID {
		t.Errorf("Expected neighbor ID %s, got %s", characterID, neighbors[0].ID)
	}
}

func TestService_GetVersion(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	versionID := createTestGraphVersion(t, database, projectID, true)

	// Test getting version
	version, err := service.GetVersion(ctx, versionID)
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if version.ID != versionID {
		t.Errorf("Expected version ID %s, got %s", versionID, version.ID)
	}

	if version.ProjectID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, version.ProjectID)
	}

	if version.IsWorkingSet != true {
		t.Errorf("Expected IsWorkingSet to be true, got %v", version.IsWorkingSet)
	}
}

func TestService_ListEntities_WithFilter(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	versionID := createTestGraphVersion(t, database, projectID, true)

	// Create entities of different types
	sceneData, _ := json.Marshal(map[string]any{"name": "Test Scene"})
	characterData, _ := json.Marshal(map[string]any{"name": "Test Character"})

	_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         uuid.New().String(),
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Test Scene",
		Data:       sceneData,
	})
	if err != nil {
		t.Fatalf("Failed to create scene entity: %v", err)
	}

	_, err = database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         uuid.New().String(),
		VersionID:  versionID,
		EntityType: "Character",
		Name:       "Test Character",
		Data:       characterData,
	})
	if err != nil {
		t.Fatalf("Failed to create character entity: %v", err)
	}

	// Test filtering by entity type
	entityType := "Scene"
	entities, err := service.ListEntities(ctx, versionID, EntityFilter{
		EntityType: &entityType,
	})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(entities) != 1 {
		t.Errorf("Expected 1 scene entity, got %d", len(entities))
	}

	if entities[0].EntityType != "Scene" {
		t.Errorf("Expected entity type 'Scene', got '%s'", entities[0].EntityType)
	}

	// Test listing all entities
	allEntities, err := service.ListEntities(ctx, versionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(allEntities) != 2 {
		t.Errorf("Expected 2 total entities, got %d", len(allEntities))
	}
}

func TestService_Apply_InvalidParentVersion(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Test with non-existent parent version
	req := &ApplyRequest{
		ParentVersionID: "non-existent-version",
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   uuid.New().String(),
				Fields:     map[string]any{"name": "Test Scene"},
			},
		},
	}

	_, err := service.Apply(ctx, req)
	if err == nil {
		t.Error("Expected error for non-existent parent version, got nil")
	}
}

func TestService_Apply_EmptyDeltas(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Test with empty deltas
	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas:          []*Delta{},
	}

	_, err := service.Apply(ctx, req)
	if err == nil {
		t.Error("Expected error for empty deltas, got nil")
	}
}

func TestService_GetNeighbors_NoRelationships(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	versionID := createTestGraphVersion(t, database, projectID, true)

	// Create entity without relationships
	entityID := uuid.New().String()
	entityData, _ := json.Marshal(map[string]any{"name": "Lonely Entity"})

	_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Lonely Entity",
		Data:       entityData,
	})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Test getting neighbors
	neighbors, err := service.GetNeighborsInVersion(ctx, versionID, entityID, "")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed: %v", err)
	}

	if len(neighbors) != 0 {
		t.Errorf("Expected 0 neighbors, got %d", len(neighbors))
	}
}

func TestService_Apply_VersioningAndParentChild(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Create initial entity in parent version
	entityID := uuid.New().String()
	entityData := map[string]any{
		"name":  "Original Entity",
		"value": "original",
	}
	dataBytes, _ := json.Marshal(entityData)

	_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         entityID,
		VersionID:  parentVersionID,
		EntityType: "TestEntity",
		Name:       "Original Entity",
		Data:       dataBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create initial entity: %v", err)
	}

	// Apply changes to create child version
	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "update",
				EntityType: "TestEntity",
				EntityID:   entityID,
				Fields: map[string]any{
					"name":  "Updated Entity",
					"value": "updated",
				},
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Verify parent-child relationship
	childVersion, err := service.GetVersion(ctx, response.GraphVersionID)
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if childVersion.ParentVersionID == nil || *childVersion.ParentVersionID != parentVersionID {
		t.Errorf("Expected parent version ID %s, got %v", parentVersionID, childVersion.ParentVersionID)
	}

	// Verify parent version is unchanged
	parentEntities, err := service.ListEntities(ctx, parentVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed for parent: %v", err)
	}

	if len(parentEntities) != 1 {
		t.Errorf("Expected 1 entity in parent version, got %d", len(parentEntities))
	}

	if parentEntities[0].Name != "Original Entity" {
		t.Errorf("Expected parent entity name 'Original Entity', got '%s'", parentEntities[0].Name)
	}

	// Verify child version has updated entity
	childEntities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed for child: %v", err)
	}

	if len(childEntities) != 1 {
		t.Errorf("Expected 1 entity in child version, got %d", len(childEntities))
	}

	if childEntities[0].Name != "Updated Entity" {
		t.Errorf("Expected child entity name 'Updated Entity', got '%s'", childEntities[0].Name)
	}
}

func TestService_Apply_ReferentialIntegrity(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Create entities with relationships in parent version
	sceneID := uuid.New().String()
	characterID := uuid.New().String()
	relationshipID := uuid.New().String()

	sceneData, _ := json.Marshal(map[string]any{"name": "Test Scene"})
	characterData, _ := json.Marshal(map[string]any{"name": "Test Character"})

	_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         sceneID,
		VersionID:  parentVersionID,
		EntityType: "Scene",
		Name:       "Test Scene",
		Data:       sceneData,
	})
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	_, err = database.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         characterID,
		VersionID:  parentVersionID,
		EntityType: "Character",
		Name:       "Test Character",
		Data:       characterData,
	})
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	_, err = database.Queries().CreateRelationship(ctx, db.CreateRelationshipParams{
		ID:               relationshipID,
		VersionID:        parentVersionID,
		FromEntityID:     sceneID,
		ToEntityID:       characterID,
		RelationshipType: "features",
		Properties:       []byte("{}"),
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	// Delete character, which should also delete its relationships
	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "delete",
				EntityType: "Character",
				EntityID:   characterID,
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Verify character is deleted
	entities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	// Should only have the scene left
	if len(entities) != 1 {
		t.Errorf("Expected 1 entity after deletion, got %d", len(entities))
	}

	if entities[0].EntityType != "Scene" {
		t.Errorf("Expected remaining entity to be Scene, got %s", entities[0].EntityType)
	}

	// Scene ID should remain the same across versions (stable entity IDs)
	// Verify relationships involving the deleted character are also gone
	neighbors, err := service.GetNeighborsInVersion(ctx, response.GraphVersionID, sceneID, "features")
	if err != nil {
		t.Fatalf("GetNeighbors failed: %v", err)
	}

	if len(neighbors) != 0 {
		t.Errorf("Expected 0 neighbors after character deletion, got %d", len(neighbors))
	}
}

func TestService_Apply_ComplexWorkflow(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	parentVersionID := createTestGraphVersion(t, database, projectID, true)

	// Apply multiple operations in a single request
	sceneID := uuid.New().String()
	characterID := uuid.New().String()
	locationID := uuid.New().String()

	req := &ApplyRequest{
		ParentVersionID: parentVersionID,
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   sceneID,
				Fields: map[string]any{
					"name":    "Opening Scene",
					"title":   "The Beginning",
					"summary": "Our story begins",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   characterID,
				Fields: map[string]any{
					"name": "Hero",
					"role": "protagonist",
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     sceneID,
						ToEntityID:       characterID,
						RelationshipType: "features",
						Properties: map[string]any{
							"importance": "primary",
						},
					},
				},
			},
			{
				Operation:  "create",
				EntityType: "Location",
				EntityID:   locationID,
				Fields: map[string]any{
					"name":        "Tavern",
					"description": "A cozy tavern",
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     sceneID,
						ToEntityID:       locationID,
						RelationshipType: "occurs_at",
						Properties: map[string]any{
							"time": "evening",
						},
					},
				},
			},
		},
	}

	response, err := service.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Applied != 3 {
		t.Errorf("Expected 3 deltas applied, got %d", response.Applied)
	}

	// Verify all entities were created
	entities, err := service.ListEntities(ctx, response.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("ListEntities failed: %v", err)
	}

	if len(entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(entities))
	}

	// Verify relationships were created correctly using stable entity IDs
	sceneNeighbors, err := service.GetNeighborsInVersion(ctx, response.GraphVersionID, sceneID, "")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed: %v", err)
	}

	if len(sceneNeighbors) != 2 {
		t.Errorf("Expected scene to have 2 neighbors, got %d", len(sceneNeighbors))
	}

	// Verify specific relationship types using original entity IDs
	featuredCharacters, err := service.GetNeighborsInVersion(ctx, response.GraphVersionID, sceneID, "features")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed for features: %v", err)
	}

	if len(featuredCharacters) != 1 {
		t.Errorf("Expected 1 featured character, got %d", len(featuredCharacters))
	}

	locations, err := service.GetNeighborsInVersion(ctx, response.GraphVersionID, sceneID, "occurs_at")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed for occurs_at: %v", err)
	}

	if len(locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(locations))
	}
}

func TestService_GetNeighbors_FilterByRelationshipType(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	service := NewService(database)
	ctx := context.Background()

	// Setup test data
	projectID := createTestProject(t, database)
	versionID := createTestGraphVersion(t, database, projectID, true)

	// Create entities
	sceneID := uuid.New().String()
	character1ID := uuid.New().String()
	character2ID := uuid.New().String()
	locationID := uuid.New().String()

	entities := []struct {
		id         string
		entityType string
		name       string
	}{
		{sceneID, "Scene", "Test Scene"},
		{character1ID, "Character", "Hero"},
		{character2ID, "Character", "Villain"},
		{locationID, "Location", "Castle"},
	}

	for _, entity := range entities {
		data, _ := json.Marshal(map[string]any{"name": entity.name})
		_, err := database.Queries().CreateEntity(ctx, db.CreateEntityParams{
			ID:         entity.id,
			VersionID:  versionID,
			EntityType: entity.entityType,
			Name:       entity.name,
			Data:       data,
		})
		if err != nil {
			t.Fatalf("Failed to create entity %s: %v", entity.name, err)
		}
	}

	// Create relationships
	relationships := []struct {
		fromID string
		toID   string
		relType string
	}{
		{sceneID, character1ID, "features"},
		{sceneID, character2ID, "features"},
		{sceneID, locationID, "occurs_at"},
	}

	for i, rel := range relationships {
		_, err := database.Queries().CreateRelationship(ctx, db.CreateRelationshipParams{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     rel.fromID,
			ToEntityID:       rel.toID,
			RelationshipType: rel.relType,
			Properties:       []byte("{}"),
		})
		if err != nil {
			t.Fatalf("Failed to create relationship %d: %v", i, err)
		}
	}

	// Test filtering by relationship type
	featuredEntities, err := service.GetNeighborsInVersion(ctx, versionID, sceneID, "features")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed for features: %v", err)
	}

	if len(featuredEntities) != 2 {
		t.Errorf("Expected 2 featured entities, got %d", len(featuredEntities))
	}

	locationEntities, err := service.GetNeighborsInVersion(ctx, versionID, sceneID, "occurs_at")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed for occurs_at: %v", err)
	}

	if len(locationEntities) != 1 {
		t.Errorf("Expected 1 location entity, got %d", len(locationEntities))
	}

	// Test getting all neighbors (no filter)
	allNeighbors, err := service.GetNeighborsInVersion(ctx, versionID, sceneID, "")
	if err != nil {
		t.Fatalf("GetNeighborsInVersion failed for all: %v", err)
	}

	if len(allNeighbors) != 3 {
		t.Errorf("Expected 3 total neighbors, got %d", len(allNeighbors))
	}
}