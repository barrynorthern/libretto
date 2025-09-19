package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/graphwrite"
	"github.com/google/uuid"
)

func TestGraphAPI_WithRelationships(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_graph_api_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Initialize database
	database, err := db.NewDatabase(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Create GraphWrite service
	graphService := graphwrite.NewService(database)

	// Create dashboard
	dashboard := &Dashboard{
		queries:      database.Queries(),
		database:     database,
		graphService: graphService,
	}

	// Create test project
	projectID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Test", Valid: true},
		Genre:       sql.NullString{String: "Test", Valid: true},
		Description: sql.NullString{String: "Test project", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Create initial version
	versionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           versionID,
		ProjectID:    projectID,
		Name:         sql.NullString{String: "Test Version", Valid: true},
		Description:  sql.NullString{String: "Test version", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Create entities with relationships using GraphWrite service
	sceneID := "test-scene-001"
	characterID := "test-character-001"

	response, err := graphService.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: versionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   sceneID,
				Fields: map[string]any{
					"name":  "Test Scene",
					"title": "A Test Scene",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   characterID,
				Fields: map[string]any{
					"name": "Test Character",
					"role": "protagonist",
				},
				Relationships: []*graphwrite.RelationshipDelta{
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
		},
	})
	if err != nil {
		t.Fatalf("Failed to create entities: %v", err)
	}

	// Update the working set to point to the new version with entities
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        response.GraphVersionID,
		ProjectID: projectID,
	})
	if err != nil {
		t.Fatalf("Failed to update working set: %v", err)
	}

	// Test the graph API
	req := httptest.NewRequest("GET", "/api/graph/"+projectID, nil)
	w := httptest.NewRecorder()

	dashboard.handleGraphAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Parse response
	var graph GraphVisualization
	if err := json.NewDecoder(w.Body).Decode(&graph); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify we have 2 nodes
	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	// Verify we have 1 relationship
	if len(graph.Links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(graph.Links))
		t.Logf("Nodes: %+v", graph.Nodes)
		t.Logf("Links: %+v", graph.Links)
	}

	// Verify the relationship connects the right entities
	if len(graph.Links) > 0 {
		link := graph.Links[0]
		if link.Source != sceneID || link.Target != characterID {
			t.Errorf("Expected link from %s to %s, got from %s to %s", 
				sceneID, characterID, link.Source, link.Target)
		}
		
		if link.Type != "features" {
			t.Errorf("Expected relationship type 'features', got '%s'", link.Type)
		}
	}

	// Verify nodes have correct connection counts
	sceneNode := findNodeByID(graph.Nodes, sceneID)
	characterNode := findNodeByID(graph.Nodes, characterID)

	if sceneNode == nil {
		t.Error("Scene node not found")
	} else if sceneNode.Size != 1 {
		t.Errorf("Expected scene to have 1 connection, got %d", sceneNode.Size)
	}

	if characterNode == nil {
		t.Error("Character node not found")
	} else if characterNode.Size != 1 {
		t.Errorf("Expected character to have 1 connection, got %d", characterNode.Size)
	}
}

func findNodeByID(nodes []Node, id string) *Node {
	for _, node := range nodes {
		if node.ID == id {
			return &node
		}
	}
	return nil
}