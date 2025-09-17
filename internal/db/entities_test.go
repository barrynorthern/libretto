package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestCreateEntity(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project and version
	projectID := uuid.New().String()
	versionID := uuid.New().String()

	projectParams := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	// Create entity
	entityID := uuid.New().String()
	sceneData := map[string]any{
		"title":          "Opening Scene",
		"summary":        "The hero begins their journey",
		"content":        "It was a dark and stormy night...",
		"act":            "Act1",
		"sequence":       1,
		"emotional_tone": "mysterious",
		"pacing":         "slow",
	}

	dataJSON, err := json.Marshal(sceneData)
	if err != nil {
		t.Fatalf("Failed to marshal scene data: %v", err)
	}

	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       dataJSON,
	}

	entity, err := queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	if entity.ID != entityID {
		t.Errorf("Expected entity ID %s, got %s", entityID, entity.ID)
	}
	if entity.VersionID != versionID {
		t.Errorf("Expected version ID %s, got %s", versionID, entity.VersionID)
	}
	if entity.EntityType != "Scene" {
		t.Errorf("Expected entity type 'Scene', got %s", entity.EntityType)
	}
	if entity.Name != "Opening Scene" {
		t.Errorf("Expected entity name 'Opening Scene', got %s", entity.Name)
	}

	// Verify data was stored correctly
	var storedData map[string]any
	err = json.Unmarshal(entity.Data, &storedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal stored data: %v", err)
	}

	if storedData["title"] != "Opening Scene" {
		t.Errorf("Expected title 'Opening Scene', got %v", storedData["title"])
	}
	if storedData["act"] != "Act1" {
		t.Errorf("Expected act 'Act1', got %v", storedData["act"])
	}
}

func TestListEntitiesByVersion(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project and version
	projectID := uuid.New().String()
	versionID := uuid.New().String()

	projectParams := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	// Create multiple entities
	sceneID := uuid.New().String()
	characterID := uuid.New().String()

	sceneData := map[string]any{
		"title":   "Opening Scene",
		"summary": "The hero begins their journey",
		"content": "It was a dark and stormy night...",
	}

	characterData := map[string]any{
		"name":        "Hero",
		"role":        "protagonist",
		"description": "A brave adventurer",
	}

	sceneDataJSON, _ := json.Marshal(sceneData)
	characterDataJSON, _ := json.Marshal(characterData)

	sceneParams := CreateEntityParams{
		ID:         sceneID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       sceneDataJSON,
	}

	characterParams := CreateEntityParams{
		ID:         characterID,
		VersionID:  versionID,
		EntityType: "Character",
		Name:       "Hero",
		Data:       characterDataJSON,
	}

	_, err = queries.CreateEntity(ctx, sceneParams)
	if err != nil {
		t.Fatalf("Failed to create scene entity: %v", err)
	}

	_, err = queries.CreateEntity(ctx, characterParams)
	if err != nil {
		t.Fatalf("Failed to create character entity: %v", err)
	}

	// List all entities in version
	entities, err := queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		t.Fatalf("Failed to list entities: %v", err)
	}

	if len(entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(entities))
	}

	// Verify entities are present
	entityTypes := make(map[string]bool)
	for _, entity := range entities {
		entityTypes[entity.EntityType] = true
	}

	if !entityTypes["Scene"] {
		t.Error("Expected Scene entity to be present")
	}
	if !entityTypes["Character"] {
		t.Error("Expected Character entity to be present")
	}
}

func TestListEntitiesByType(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project and version
	projectID := uuid.New().String()
	versionID := uuid.New().String()

	projectParams := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	// Create multiple characters and one scene
	character1ID := uuid.New().String()
	character2ID := uuid.New().String()
	sceneID := uuid.New().String()

	character1Data := map[string]any{"name": "Hero", "role": "protagonist"}
	character2Data := map[string]any{"name": "Villain", "role": "antagonist"}
	sceneData := map[string]any{"title": "Opening Scene", "summary": "The beginning"}

	character1DataJSON, _ := json.Marshal(character1Data)
	character2DataJSON, _ := json.Marshal(character2Data)
	sceneDataJSON, _ := json.Marshal(sceneData)

	entities := []CreateEntityParams{
		{
			ID:         character1ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Hero",
			Data:       character1DataJSON,
		},
		{
			ID:         character2ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Villain",
			Data:       character2DataJSON,
		},
		{
			ID:         sceneID,
			VersionID:  versionID,
			EntityType: "Scene",
			Name:       "Opening Scene",
			Data:       sceneDataJSON,
		},
	}

	for _, entityParams := range entities {
		_, err = queries.CreateEntity(ctx, entityParams)
		if err != nil {
			t.Fatalf("Failed to create entity %s: %v", entityParams.Name, err)
		}
	}

	// List only Character entities
	listParams := ListEntitiesByTypeParams{
		VersionID:  versionID,
		EntityType: "Character",
	}

	characters, err := queries.ListEntitiesByType(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list character entities: %v", err)
	}

	if len(characters) != 2 {
		t.Errorf("Expected 2 character entities, got %d", len(characters))
	}

	for _, character := range characters {
		if character.EntityType != "Character" {
			t.Errorf("Expected entity type 'Character', got %s", character.EntityType)
		}
	}
}

func TestUpdateEntity(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project and version
	projectID := uuid.New().String()
	versionID := uuid.New().String()

	projectParams := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	// Create entity
	entityID := uuid.New().String()
	originalData := map[string]any{
		"title":   "Original Title",
		"summary": "Original summary",
		"content": "Original content",
	}

	originalDataJSON, _ := json.Marshal(originalData)

	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Original Scene",
		Data:       originalDataJSON,
	}

	originalEntity, err := queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Update entity
	updatedData := map[string]any{
		"title":   "Updated Title",
		"summary": "Updated summary",
		"content": "Updated content",
		"act":     "Act1",
	}

	updatedDataJSON, _ := json.Marshal(updatedData)

	updateParams := UpdateEntityParams{
		ID:   entityID,
		Name: "Updated Scene",
		Data: updatedDataJSON,
	}

	updatedEntity, err := queries.UpdateEntity(ctx, updateParams)
	if err != nil {
		t.Fatalf("Failed to update entity: %v", err)
	}

	if updatedEntity.Name != "Updated Scene" {
		t.Errorf("Expected updated name 'Updated Scene', got %s", updatedEntity.Name)
	}

	// Verify data was updated
	var storedData map[string]any
	err = json.Unmarshal(updatedEntity.Data, &storedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal updated data: %v", err)
	}

	if storedData["title"] != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %v", storedData["title"])
	}
	if storedData["act"] != "Act1" {
		t.Errorf("Expected act 'Act1', got %v", storedData["act"])
	}

	// Verify updated_at changed (or at least didn't go backwards)
	if updatedEntity.UpdatedAt.Before(originalEntity.UpdatedAt) {
		t.Error("Expected updated_at to not go backwards after update")
	}
}

func TestCountEntitiesByType(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project and version
	projectID := uuid.New().String()
	versionID := uuid.New().String()

	projectParams := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	// Create multiple entities of different types
	entities := []struct {
		entityType string
		name       string
	}{
		{"Character", "Hero"},
		{"Character", "Villain"},
		{"Character", "Sidekick"},
		{"Scene", "Opening"},
		{"Scene", "Climax"},
		{"Location", "Castle"},
	}

	for i, entity := range entities {
		entityID := uuid.New().String()
		data := map[string]any{"name": entity.name}
		dataJSON, _ := json.Marshal(data)

		entityParams := CreateEntityParams{
			ID:         entityID,
			VersionID:  versionID,
			EntityType: entity.entityType,
			Name:       entity.name,
			Data:       dataJSON,
		}

		_, err = queries.CreateEntity(ctx, entityParams)
		if err != nil {
			t.Fatalf("Failed to create entity %d: %v", i, err)
		}
	}

	// Count characters
	countParams := CountEntitiesByTypeParams{
		VersionID:  versionID,
		EntityType: "Character",
	}

	characterCount, err := queries.CountEntitiesByType(ctx, countParams)
	if err != nil {
		t.Fatalf("Failed to count character entities: %v", err)
	}

	if characterCount != 3 {
		t.Errorf("Expected 3 character entities, got %d", characterCount)
	}

	// Count scenes
	countParams.EntityType = "Scene"
	sceneCount, err := queries.CountEntitiesByType(ctx, countParams)
	if err != nil {
		t.Fatalf("Failed to count scene entities: %v", err)
	}

	if sceneCount != 2 {
		t.Errorf("Expected 2 scene entities, got %d", sceneCount)
	}

	// Count locations
	countParams.EntityType = "Location"
	locationCount, err := queries.CountEntitiesByType(ctx, countParams)
	if err != nil {
		t.Fatalf("Failed to count location entities: %v", err)
	}

	if locationCount != 1 {
		t.Errorf("Expected 1 location entity, got %d", locationCount)
	}
}