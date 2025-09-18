package server

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/graphwrite"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
	"connectrpc.com/connect"
	"github.com/google/uuid"
)

func setupIntegrationTest(t *testing.T) (*GraphWriteServer, *db.Database, string, string) {
	// Create temporary database file
	tmpFile, err := os.CreateTemp("", "libretto_integration_test_*.db")
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

	// Create project and initial graph version
	projectID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "Integration Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Test project for integration", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	versionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           versionID,
		ProjectID:    projectID,
		Name:         sql.NullString{String: "Initial Version", Valid: true},
		Description:  sql.NullString{String: "Initial version for testing", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create test graph version: %v", err)
	}

	// Create GraphWrite service and server
	service := graphwrite.NewService(database)
	server := NewGraphWriteServer(service)

	return server, database, projectID, versionID
}

func TestGraphWriteServer_Apply_Integration(t *testing.T) {
	server, database, _, versionID := setupIntegrationTest(t)
	defer database.Close()

	ctx := context.Background()

	// Create a request to add a scene
	req := connect.NewRequest(&graphv1.ApplyRequest{
		ParentVersionId: versionID,
		Deltas: []*graphv1.Delta{
			{
				Op:         "create",
				EntityType: "Scene",
				EntityId:   uuid.New().String(),
				Fields: map[string]string{
					"name":    "Opening Scene",
					"title":   "The Beginning",
					"summary": "Our hero starts their journey",
					"content": "It was a dark and stormy night...",
				},
			},
		},
	})

	// Apply the request
	response, err := server.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Msg.Applied != 1 {
		t.Errorf("Expected 1 delta applied, got %d", response.Msg.Applied)
	}

	if response.Msg.GraphVersionId == "" {
		t.Error("Expected non-empty graph version ID")
	}

	// Verify the new version was created
	newVersion, err := database.Queries().GetGraphVersion(ctx, response.Msg.GraphVersionId)
	if err != nil {
		t.Fatalf("Failed to get new version: %v", err)
	}

	if newVersion.ParentVersionID.String != versionID {
		t.Errorf("Expected parent version ID %s, got %s", versionID, newVersion.ParentVersionID.String)
	}

	// Verify the entity was created
	entities, err := database.Queries().ListEntitiesByVersion(ctx, response.Msg.GraphVersionId)
	if err != nil {
		t.Fatalf("Failed to list entities: %v", err)
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
}

func TestGraphWriteServer_Apply_MultipleDeltas_Integration(t *testing.T) {
	server, database, _, versionID := setupIntegrationTest(t)
	defer database.Close()

	ctx := context.Background()

	sceneID := uuid.New().String()
	characterID := uuid.New().String()

	// Create a request with multiple deltas
	req := connect.NewRequest(&graphv1.ApplyRequest{
		ParentVersionId: versionID,
		Deltas: []*graphv1.Delta{
			{
				Op:         "create",
				EntityType: "Scene",
				EntityId:   sceneID,
				Fields: map[string]string{
					"name":    "Opening Scene",
					"title":   "The Beginning",
					"summary": "Our story begins",
				},
			},
			{
				Op:         "create",
				EntityType: "Character",
				EntityId:   characterID,
				Fields: map[string]string{
					"name": "Hero",
					"role": "protagonist",
				},
			},
		},
	})

	// Apply the request
	response, err := server.Apply(ctx, req)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if response.Msg.Applied != 2 {
		t.Errorf("Expected 2 deltas applied, got %d", response.Msg.Applied)
	}

	// Verify both entities were created
	entities, err := database.Queries().ListEntitiesByVersion(ctx, response.Msg.GraphVersionId)
	if err != nil {
		t.Fatalf("Failed to list entities: %v", err)
	}

	if len(entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(entities))
	}

	// Verify entity types
	entityTypes := make(map[string]bool)
	for _, entity := range entities {
		entityTypes[entity.EntityType] = true
	}

	if !entityTypes["Scene"] {
		t.Error("Expected Scene entity to be created")
	}

	if !entityTypes["Character"] {
		t.Error("Expected Character entity to be created")
	}
}

func TestGraphWriteServer_Apply_InvalidRequest_Integration(t *testing.T) {
	server, database, _, _ := setupIntegrationTest(t)
	defer database.Close()

	ctx := context.Background()

	// Test with empty deltas
	req := connect.NewRequest(&graphv1.ApplyRequest{
		ParentVersionId: "some-version",
		Deltas:          []*graphv1.Delta{},
	})

	_, err := server.Apply(ctx, req)
	if err == nil {
		t.Error("Expected error for empty deltas, got nil")
	}

	// Test with non-existent parent version
	req = connect.NewRequest(&graphv1.ApplyRequest{
		ParentVersionId: "non-existent-version",
		Deltas: []*graphv1.Delta{
			{
				Op:         "create",
				EntityType: "Scene",
				EntityId:   uuid.New().String(),
				Fields:     map[string]string{"name": "Test Scene"},
			},
		},
	})

	_, err = server.Apply(ctx, req)
	if err == nil {
		t.Error("Expected error for non-existent parent version, got nil")
	}
}