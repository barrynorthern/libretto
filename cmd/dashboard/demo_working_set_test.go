package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/graphwrite"
)

func TestDemoWorkingSetUpdate(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_demo_working_set_test_*.db")
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

	// Step 1: Create story via demo
	req1 := httptest.NewRequest("POST", "/api/demo/create-story", nil)
	w1 := httptest.NewRecorder()

	dashboard.handleCreateStoryDemo(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200 for create story, got %d", w1.Code)
		t.Logf("Response body: %s", w1.Body.String())
	}

	// Parse create story response
	var createResult map[string]any
	if err := json.NewDecoder(w1.Body).Decode(&createResult); err != nil {
		t.Fatalf("Failed to decode create story response: %v", err)
	}

	projectID, ok := createResult["projectId"].(string)
	if !ok {
		t.Fatalf("Expected projectId in response, got %v", createResult["projectId"])
	}

	// Step 2: Check that the graph API now shows entities and relationships
	req2 := httptest.NewRequest("GET", "/api/graph/"+projectID, nil)
	w2 := httptest.NewRecorder()

	dashboard.handleGraphAPI(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200 for graph API, got %d", w2.Code)
		t.Logf("Response body: %s", w2.Body.String())
	}

	// Parse graph response
	var graph GraphVisualization
	if err := json.NewDecoder(w2.Body).Decode(&graph); err != nil {
		t.Fatalf("Failed to decode graph response: %v", err)
	}

	// Verify we have entities
	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes after create story, got %d", len(graph.Nodes))
		t.Logf("Nodes: %+v", graph.Nodes)
	}

	// Verify we have relationships
	if len(graph.Links) != 1 {
		t.Errorf("Expected 1 link after create story, got %d", len(graph.Links))
		t.Logf("Links: %+v", graph.Links)
	}

	// Step 3: Add character via demo
	addCharReq := map[string]any{
		"projectId":       projectID,
		"parentVersionId": createResult["versionId"],
		"sceneId":         createResult["sceneId"],
	}
	addCharBody, _ := json.Marshal(addCharReq)

	req3 := httptest.NewRequest("POST", "/api/demo/add-character", bytes.NewReader(addCharBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()

	dashboard.handleAddCharacterDemo(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200 for add character, got %d", w3.Code)
		t.Logf("Response body: %s", w3.Body.String())
	}

	// Step 4: Check graph API again - should now have more entities
	req4 := httptest.NewRequest("GET", "/api/graph/"+projectID, nil)
	w4 := httptest.NewRecorder()

	dashboard.handleGraphAPI(w4, req4)

	if w4.Code != http.StatusOK {
		t.Errorf("Expected status 200 for graph API after add character, got %d", w4.Code)
		t.Logf("Response body: %s", w4.Body.String())
	}

	// Parse updated graph response
	var updatedGraph GraphVisualization
	if err := json.NewDecoder(w4.Body).Decode(&updatedGraph); err != nil {
		t.Fatalf("Failed to decode updated graph response: %v", err)
	}

	// Verify we now have more entities (scene, Elena, Mordak, tavern)
	if len(updatedGraph.Nodes) != 4 {
		t.Errorf("Expected 4 nodes after add character, got %d", len(updatedGraph.Nodes))
		t.Logf("Updated Nodes: %+v", updatedGraph.Nodes)
	}

	// Verify we have more relationships
	if len(updatedGraph.Links) != 3 {
		t.Errorf("Expected 3 links after add character, got %d", len(updatedGraph.Links))
		t.Logf("Updated Links: %+v", updatedGraph.Links)
	}

	// Verify entity types are correct
	entityTypes := make(map[string]int)
	for _, node := range updatedGraph.Nodes {
		entityTypes[node.Type]++
	}

	expectedTypes := map[string]int{
		"Scene":     1,
		"Character": 2, // Elena + Mordak
		"Location":  1, // Tavern
	}

	for expectedType, expectedCount := range expectedTypes {
		if actualCount := entityTypes[expectedType]; actualCount != expectedCount {
			t.Errorf("Expected %d %s entities, got %d", expectedCount, expectedType, actualCount)
		}
	}

	t.Logf("✅ Demo working set update test passed!")
	t.Logf("   - Created story with 2 entities and 1 relationship")
	t.Logf("   - Added characters to get 4 entities and 3 relationships")
	t.Logf("   - Graph API correctly shows all entities and relationships")
}