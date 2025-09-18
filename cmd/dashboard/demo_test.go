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

func setupTestDashboard(t *testing.T) *Dashboard {
	// Create temporary database file
	tmpFile, err := os.CreateTemp("", "libretto_dashboard_test_*.db")
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

	// Initialize GraphWrite service
	graphService := graphwrite.NewService(database)

	return &Dashboard{
		queries:      database.Queries(),
		database:     database,
		graphService: graphService,
	}
}

func TestDashboard_CreateStoryDemo(t *testing.T) {
	dashboard := setupTestDashboard(t)

	// Create request
	req := httptest.NewRequest("POST", "/api/demo/create-story", nil)
	w := httptest.NewRecorder()

	// Handle request
	dashboard.handleCreateStoryDemo(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response
	var result map[string]any
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if result["projectId"] == nil {
		t.Error("Expected projectId in response")
	}

	if result["versionId"] == nil {
		t.Error("Expected versionId in response")
	}

	if result["sceneId"] == nil {
		t.Error("Expected sceneId in response")
	}

	if result["characterId"] == nil {
		t.Error("Expected characterId in response")
	}

	entities, ok := result["entities"].([]any)
	if !ok || len(entities) != 2 {
		t.Errorf("Expected 2 entities, got %v", result["entities"])
	}

	applied, ok := result["applied"].(float64)
	if !ok || applied != 2 {
		t.Errorf("Expected 2 applied deltas, got %v", result["applied"])
	}
}

func TestDashboard_AddCharacterDemo(t *testing.T) {
	dashboard := setupTestDashboard(t)

	// First create a story
	req1 := httptest.NewRequest("POST", "/api/demo/create-story", nil)
	w1 := httptest.NewRecorder()
	dashboard.handleCreateStoryDemo(w1, req1)

	var createResult map[string]any
	if err := json.NewDecoder(w1.Body).Decode(&createResult); err != nil {
		t.Fatalf("Failed to decode create story response: %v", err)
	}

	// Now add character
	requestBody := map[string]any{
		"projectId":       createResult["projectId"],
		"parentVersionId": createResult["versionId"],
		"sceneId":         createResult["sceneId"],
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req2 := httptest.NewRequest("POST", "/api/demo/add-character", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	dashboard.handleAddCharacterDemo(w2, req2)

	// Check response
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(w2.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should now have 4 entities (scene, original character, new character, location)
	entities, ok := result["entities"].([]any)
	if !ok || len(entities) != 4 {
		t.Errorf("Expected 4 entities after adding character, got %v", result["entities"])
	}

	applied, ok := result["applied"].(float64)
	if !ok || applied != 2 {
		t.Errorf("Expected 2 applied deltas, got %v", result["applied"])
	}
}

func TestDashboard_UpdateSceneDemo(t *testing.T) {
	dashboard := setupTestDashboard(t)

	// First create a story
	req1 := httptest.NewRequest("POST", "/api/demo/create-story", nil)
	w1 := httptest.NewRecorder()
	dashboard.handleCreateStoryDemo(w1, req1)

	var createResult map[string]any
	if err := json.NewDecoder(w1.Body).Decode(&createResult); err != nil {
		t.Fatalf("Failed to decode create story response: %v", err)
	}

	// Update scene
	requestBody := map[string]any{
		"projectId":       createResult["projectId"],
		"parentVersionId": createResult["versionId"],
		"sceneId":         createResult["sceneId"],
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req2 := httptest.NewRequest("POST", "/api/demo/update-scene", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	dashboard.handleUpdateSceneDemo(w2, req2)

	// Check response
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(w2.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should still have 2 entities (scene updated, character unchanged)
	entities, ok := result["entities"].([]any)
	if !ok || len(entities) != 2 {
		t.Errorf("Expected 2 entities after updating scene, got %v", result["entities"])
	}

	applied, ok := result["applied"].(float64)
	if !ok || applied != 1 {
		t.Errorf("Expected 1 applied delta, got %v", result["applied"])
	}
}

func TestDashboard_DemoPageRenders(t *testing.T) {
	dashboard := setupTestDashboard(t)

	req := httptest.NewRequest("GET", "/demo", nil)
	w := httptest.NewRecorder()

	dashboard.handleDemo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("GraphWrite Service Demo")) {
		t.Error("Expected demo page to contain title")
	}

	if !bytes.Contains([]byte(body), []byte("Create Story")) {
		t.Error("Expected demo page to contain Create Story button")
	}
}