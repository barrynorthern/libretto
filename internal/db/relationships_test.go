package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestCreateRelationship(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entities
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	sceneID := uuid.New().String()
	characterID := uuid.New().String()

	// Create project
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

	// Create version
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

	// Create entities
	sceneData := map[string]any{"title": "Opening Scene", "summary": "The beginning"}
	characterData := map[string]any{"name": "Hero", "role": "protagonist"}

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

	// Create relationship
	relationshipID := uuid.New().String()
	properties := map[string]any{
		"role":       "protagonist",
		"importance": "high",
	}

	propertiesJSON, _ := json.Marshal(properties)

	relationshipParams := CreateRelationshipParams{
		ID:               relationshipID,
		VersionID:        versionID,
		FromEntityID:     sceneID,
		ToEntityID:       characterID,
		RelationshipType: "features",
		Properties:       propertiesJSON,
	}

	relationship, err := queries.CreateRelationship(ctx, relationshipParams)
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	if relationship.ID != relationshipID {
		t.Errorf("Expected relationship ID %s, got %s", relationshipID, relationship.ID)
	}
	if relationship.VersionID != versionID {
		t.Errorf("Expected version ID %s, got %s", versionID, relationship.VersionID)
	}
	if relationship.FromEntityID != sceneID {
		t.Errorf("Expected from entity ID %s, got %s", sceneID, relationship.FromEntityID)
	}
	if relationship.ToEntityID != characterID {
		t.Errorf("Expected to entity ID %s, got %s", characterID, relationship.ToEntityID)
	}
	if relationship.RelationshipType != "features" {
		t.Errorf("Expected relationship type 'features', got %s", relationship.RelationshipType)
	}

	// Verify properties
	var storedProperties map[string]any
	err = json.Unmarshal(relationship.Properties, &storedProperties)
	if err != nil {
		t.Fatalf("Failed to unmarshal properties: %v", err)
	}

	if storedProperties["role"] != "protagonist" {
		t.Errorf("Expected role 'protagonist', got %v", storedProperties["role"])
	}
}

func TestListRelationshipsByEntity(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entities
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	sceneID := uuid.New().String()
	character1ID := uuid.New().String()
	character2ID := uuid.New().String()

	// Setup project and version
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

	// Create entities
	entities := []CreateEntityParams{
		{
			ID:         sceneID,
			VersionID:  versionID,
			EntityType: "Scene",
			Name:       "Opening Scene",
			Data:       json.RawMessage(`{"title": "Opening Scene"}`),
		},
		{
			ID:         character1ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Hero",
			Data:       json.RawMessage(`{"name": "Hero"}`),
		},
		{
			ID:         character2ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Villain",
			Data:       json.RawMessage(`{"name": "Villain"}`),
		},
	}

	for _, entityParams := range entities {
		_, err = queries.CreateEntity(ctx, entityParams)
		if err != nil {
			t.Fatalf("Failed to create entity %s: %v", entityParams.Name, err)
		}
	}

	// Create relationships
	relationships := []CreateRelationshipParams{
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     sceneID,
			ToEntityID:       character1ID,
			RelationshipType: "features",
			Properties:       json.RawMessage(`{"role": "protagonist"}`),
		},
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     sceneID,
			ToEntityID:       character2ID,
			RelationshipType: "features",
			Properties:       json.RawMessage(`{"role": "antagonist"}`),
		},
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     character1ID,
			ToEntityID:       character2ID,
			RelationshipType: "conflicts",
			Properties:       json.RawMessage(`{"intensity": "high"}`),
		},
	}

	for _, relationshipParams := range relationships {
		_, err = queries.CreateRelationship(ctx, relationshipParams)
		if err != nil {
			t.Fatalf("Failed to create relationship: %v", err)
		}
	}

	// List relationships for scene (should have 2)
	listParams := ListRelationshipsByEntityParams{
		FromEntityID: sceneID,
		ToEntityID:   sceneID,
	}

	sceneRelationships, err := queries.ListRelationshipsByEntity(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list relationships for scene: %v", err)
	}

	if len(sceneRelationships) != 2 {
		t.Errorf("Expected 2 relationships for scene, got %d", len(sceneRelationships))
	}

	// List relationships for character1 (should have 2)
	listParams = ListRelationshipsByEntityParams{
		FromEntityID: character1ID,
		ToEntityID:   character1ID,
	}

	character1Relationships, err := queries.ListRelationshipsByEntity(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list relationships for character1: %v", err)
	}

	if len(character1Relationships) != 2 {
		t.Errorf("Expected 2 relationships for character1, got %d", len(character1Relationships))
	}
}

func TestListRelationshipsByType(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entities
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	sceneID := uuid.New().String()
	character1ID := uuid.New().String()
	character2ID := uuid.New().String()

	// Setup project and version
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

	// Create entities
	entities := []CreateEntityParams{
		{
			ID:         sceneID,
			VersionID:  versionID,
			EntityType: "Scene",
			Name:       "Opening Scene",
			Data:       json.RawMessage(`{"title": "Opening Scene"}`),
		},
		{
			ID:         character1ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Hero",
			Data:       json.RawMessage(`{"name": "Hero"}`),
		},
		{
			ID:         character2ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Villain",
			Data:       json.RawMessage(`{"name": "Villain"}`),
		},
	}

	for _, entityParams := range entities {
		_, err = queries.CreateEntity(ctx, entityParams)
		if err != nil {
			t.Fatalf("Failed to create entity %s: %v", entityParams.Name, err)
		}
	}

	// Create relationships of different types
	relationships := []CreateRelationshipParams{
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     sceneID,
			ToEntityID:       character1ID,
			RelationshipType: "features",
			Properties:       json.RawMessage(`{"role": "protagonist"}`),
		},
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     sceneID,
			ToEntityID:       character2ID,
			RelationshipType: "features",
			Properties:       json.RawMessage(`{"role": "antagonist"}`),
		},
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     character1ID,
			ToEntityID:       character2ID,
			RelationshipType: "conflicts",
			Properties:       json.RawMessage(`{"intensity": "high"}`),
		},
	}

	for _, relationshipParams := range relationships {
		_, err = queries.CreateRelationship(ctx, relationshipParams)
		if err != nil {
			t.Fatalf("Failed to create relationship: %v", err)
		}
	}

	// List "features" relationships
	listParams := ListRelationshipsByTypeParams{
		VersionID:        versionID,
		RelationshipType: "features",
	}

	featuresRelationships, err := queries.ListRelationshipsByType(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list features relationships: %v", err)
	}

	if len(featuresRelationships) != 2 {
		t.Errorf("Expected 2 features relationships, got %d", len(featuresRelationships))
	}

	for _, rel := range featuresRelationships {
		if rel.RelationshipType != "features" {
			t.Errorf("Expected relationship type 'features', got %s", rel.RelationshipType)
		}
	}

	// List "conflicts" relationships
	listParams.RelationshipType = "conflicts"
	conflictsRelationships, err := queries.ListRelationshipsByType(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list conflicts relationships: %v", err)
	}

	if len(conflictsRelationships) != 1 {
		t.Errorf("Expected 1 conflicts relationship, got %d", len(conflictsRelationships))
	}
}

func TestGetRelationshipsBetweenEntities(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entities
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	character1ID := uuid.New().String()
	character2ID := uuid.New().String()

	// Setup project and version
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

	// Create entities
	entities := []CreateEntityParams{
		{
			ID:         character1ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Hero",
			Data:       json.RawMessage(`{"name": "Hero"}`),
		},
		{
			ID:         character2ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Villain",
			Data:       json.RawMessage(`{"name": "Villain"}`),
		},
	}

	for _, entityParams := range entities {
		_, err = queries.CreateEntity(ctx, entityParams)
		if err != nil {
			t.Fatalf("Failed to create entity %s: %v", entityParams.Name, err)
		}
	}

	// Create multiple relationships between the same entities
	relationships := []CreateRelationshipParams{
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     character1ID,
			ToEntityID:       character2ID,
			RelationshipType: "conflicts",
			Properties:       json.RawMessage(`{"intensity": "high"}`),
		},
		{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     character1ID,
			ToEntityID:       character2ID,
			RelationshipType: "knows",
			Properties:       json.RawMessage(`{"since": "childhood"}`),
		},
	}

	for _, relationshipParams := range relationships {
		_, err = queries.CreateRelationship(ctx, relationshipParams)
		if err != nil {
			t.Fatalf("Failed to create relationship: %v", err)
		}
	}

	// Get relationships between the two characters
	betweenParams := GetRelationshipsBetweenEntitiesParams{
		FromEntityID: character1ID,
		ToEntityID:   character2ID,
	}

	betweenRelationships, err := queries.GetRelationshipsBetweenEntities(ctx, betweenParams)
	if err != nil {
		t.Fatalf("Failed to get relationships between entities: %v", err)
	}

	if len(betweenRelationships) != 2 {
		t.Errorf("Expected 2 relationships between entities, got %d", len(betweenRelationships))
	}

	// Verify relationship types
	relationshipTypes := make(map[string]bool)
	for _, rel := range betweenRelationships {
		relationshipTypes[rel.RelationshipType] = true
	}

	if !relationshipTypes["conflicts"] {
		t.Error("Expected 'conflicts' relationship to be present")
	}
	if !relationshipTypes["knows"] {
		t.Error("Expected 'knows' relationship to be present")
	}
}

func TestUniqueRelationshipConstraint(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entities
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	character1ID := uuid.New().String()
	character2ID := uuid.New().String()

	// Setup project and version
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

	// Create entities
	entities := []CreateEntityParams{
		{
			ID:         character1ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Hero",
			Data:       json.RawMessage(`{"name": "Hero"}`),
		},
		{
			ID:         character2ID,
			VersionID:  versionID,
			EntityType: "Character",
			Name:       "Villain",
			Data:       json.RawMessage(`{"name": "Villain"}`),
		},
	}

	for _, entityParams := range entities {
		_, err = queries.CreateEntity(ctx, entityParams)
		if err != nil {
			t.Fatalf("Failed to create entity %s: %v", entityParams.Name, err)
		}
	}

	// Create first relationship
	relationship1Params := CreateRelationshipParams{
		ID:               uuid.New().String(),
		VersionID:        versionID,
		FromEntityID:     character1ID,
		ToEntityID:       character2ID,
		RelationshipType: "conflicts",
		Properties:       json.RawMessage(`{"intensity": "high"}`),
	}

	_, err = queries.CreateRelationship(ctx, relationship1Params)
	if err != nil {
		t.Fatalf("Failed to create first relationship: %v", err)
	}

	// Try to create duplicate relationship - should fail
	relationship2Params := CreateRelationshipParams{
		ID:               uuid.New().String(),
		VersionID:        versionID,
		FromEntityID:     character1ID,
		ToEntityID:       character2ID,
		RelationshipType: "conflicts", // Same type as first relationship
		Properties:       json.RawMessage(`{"intensity": "low"}`),
	}

	_, err = queries.CreateRelationship(ctx, relationship2Params)
	if err == nil {
		t.Error("Expected error when creating duplicate relationship")
	}
}