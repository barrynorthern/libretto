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

func TestDeleteProject(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_delete_test_*.db")
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
	project, err := database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project to Delete",
		Theme:       sql.NullString{String: "Test", Valid: true},
		Genre:       sql.NullString{String: "Test", Valid: true},
		Description: sql.NullString{String: "This project will be deleted", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Verify project exists
	_, err = database.Queries().GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("Project should exist before deletion: %v", err)
	}

	// Test DELETE request
	req := httptest.NewRequest("DELETE", "/api/project/delete/"+projectID, nil)
	w := httptest.NewRecorder()

	dashboard.handleDeleteProject(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Parse response
	var result map[string]any
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if success, ok := result["success"].(bool); !ok || !success {
		t.Errorf("Expected success=true, got %v", result["success"])
	}

	if projectName, ok := result["projectName"].(string); !ok || projectName != project.Name {
		t.Errorf("Expected projectName='%s', got %v", project.Name, result["projectName"])
	}

	// Verify project is actually deleted
	_, err = database.Queries().GetProject(ctx, projectID)
	if err == nil {
		t.Error("Project should be deleted but still exists")
	}
}

func TestDeleteProject_NotFound(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_delete_test_*.db")
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

	// Test DELETE request for non-existent project
	nonExistentID := uuid.New().String()
	req := httptest.NewRequest("DELETE", "/api/project/delete/"+nonExistentID, nil)
	w := httptest.NewRecorder()

	dashboard.handleDeleteProject(w, req)

	// Should return 404
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteProject_InvalidMethod(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_delete_test_*.db")
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

	// Test GET request (should be rejected)
	req := httptest.NewRequest("GET", "/api/project/delete/some-id", nil)
	w := httptest.NewRecorder()

	dashboard.handleDeleteProject(w, req)

	// Should return 405 Method Not Allowed
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestDeleteProject_WithSharedEntities(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_delete_shared_test_*.db")
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

	// Create two projects
	project1ID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          project1ID,
		Name:        "Book 1: The Lost Artifact",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "First book", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create project 1: %v", err)
	}

	project2ID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          project2ID,
		Name:        "Book 2: The Shadow War",
		Theme:       sql.NullString{String: "War", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Second book", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create project 2: %v", err)
	}

	// Create versions for both projects
	version1ID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           version1ID,
		ProjectID:    project1ID,
		Name:         sql.NullString{String: "Version 1", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	version2ID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           version2ID,
		ProjectID:    project2ID,
		Name:         sql.NullString{String: "Version 2", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Create Elena in project 1
	elenaID := "elena-stormwind-protagonist"
	response1, err := graphService.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: version1ID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name": "Elena Stormwind",
					"role": "protagonist",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create Elena in project 1: %v", err)
	}

	// Update working set for project 1
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        response1.GraphVersionID,
		ProjectID: project1ID,
	})
	if err != nil {
		t.Fatalf("Failed to update working set 1: %v", err)
	}

	// Import Elena into project 2
	_, err = graphService.ImportEntity(ctx, version2ID, project1ID, elenaID)
	if err != nil {
		t.Fatalf("Failed to import Elena to project 2: %v", err)
	}

	// Try to delete project 1 (should fail due to shared entity)
	req := httptest.NewRequest("DELETE", "/api/project/delete/"+project1ID, nil)
	w := httptest.NewRecorder()

	dashboard.handleDeleteProject(w, req)

	// Should return 409 Conflict
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Parse response
	var result map[string]any
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response indicates shared entities conflict
	if success, ok := result["success"].(bool); !ok || success {
		t.Errorf("Expected success=false, got %v", result["success"])
	}

	if sharedEntities, ok := result["sharedEntities"].([]any); !ok || len(sharedEntities) == 0 {
		t.Errorf("Expected sharedEntities list, got %v", result["sharedEntities"])
	}

	// Verify project still exists (wasn't deleted)
	_, err = database.Queries().GetProject(ctx, project1ID)
	if err != nil {
		t.Error("Project should still exist after failed deletion attempt")
	}
}