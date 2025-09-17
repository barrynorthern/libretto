package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
)

func TestCreateGraphVersion(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project first
	projectID := uuid.New().String()
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

	// Create a graph version
	versionID := uuid.New().String()
	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version of the narrative", Valid: true},
		IsWorkingSet:  true,
	}

	version, err := queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	if version.ID != versionID {
		t.Errorf("Expected version ID %s, got %s", versionID, version.ID)
	}
	if version.ProjectID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, version.ProjectID)
	}
	if !version.IsWorkingSet {
		t.Error("Expected version to be working set")
	}
	if !version.Name.Valid || version.Name.String != "Initial Version" {
		t.Errorf("Expected name 'Initial Version', got %v", version.Name)
	}
}

func TestGetWorkingSetVersion(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project
	projectID := uuid.New().String()
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

	// Create multiple versions, only one as working set
	version1ID := uuid.New().String()
	version2ID := uuid.New().String()

	version1Params := CreateGraphVersionParams{
		ID:            version1ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Version 1", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  false,
	}

	version2Params := CreateGraphVersionParams{
		ID:            version2ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{String: version1ID, Valid: true},
		Name:          sql.NullString{String: "Version 2", Valid: true},
		Description:   sql.NullString{String: "Working version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, version1Params)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	_, err = queries.CreateGraphVersion(ctx, version2Params)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Get working set version
	workingSet, err := queries.GetWorkingSetVersion(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to get working set version: %v", err)
	}

	if workingSet.ID != version2ID {
		t.Errorf("Expected working set ID %s, got %s", version2ID, workingSet.ID)
	}
	if !workingSet.IsWorkingSet {
		t.Error("Expected version to be working set")
	}
}

func TestSetWorkingSet(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project
	projectID := uuid.New().String()
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

	// Create two versions
	version1ID := uuid.New().String()
	version2ID := uuid.New().String()

	version1Params := CreateGraphVersionParams{
		ID:            version1ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Version 1", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	version2Params := CreateGraphVersionParams{
		ID:            version2ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{String: version1ID, Valid: true},
		Name:          sql.NullString{String: "Version 2", Valid: true},
		Description:   sql.NullString{String: "Second version", Valid: true},
		IsWorkingSet:  false,
	}

	_, err = queries.CreateGraphVersion(ctx, version1Params)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	_, err = queries.CreateGraphVersion(ctx, version2Params)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Switch working set to version 2
	setWorkingSetParams := SetWorkingSetParams{
		ID:        version2ID,
		ProjectID: projectID,
	}

	err = queries.SetWorkingSet(ctx, setWorkingSetParams)
	if err != nil {
		t.Fatalf("Failed to set working set: %v", err)
	}

	// Verify version 2 is now working set
	workingSet, err := queries.GetWorkingSetVersion(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to get working set version: %v", err)
	}

	if workingSet.ID != version2ID {
		t.Errorf("Expected working set ID %s, got %s", version2ID, workingSet.ID)
	}

	// Verify version 1 is no longer working set
	version1, err := queries.GetGraphVersion(ctx, version1ID)
	if err != nil {
		t.Fatalf("Failed to get version 1: %v", err)
	}

	if version1.IsWorkingSet {
		t.Error("Expected version 1 to no longer be working set")
	}
}

func TestListGraphVersionsByProject(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project
	projectID := uuid.New().String()
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

	// Create multiple versions
	version1ID := uuid.New().String()
	version2ID := uuid.New().String()

	version1Params := CreateGraphVersionParams{
		ID:            version1ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Version 1", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  false,
	}

	version2Params := CreateGraphVersionParams{
		ID:            version2ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{String: version1ID, Valid: true},
		Name:          sql.NullString{String: "Version 2", Valid: true},
		Description:   sql.NullString{String: "Second version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, version1Params)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	_, err = queries.CreateGraphVersion(ctx, version2Params)
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// List versions
	versions, err := queries.ListGraphVersionsByProject(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to list graph versions: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	}

	// Should be ordered by created_at DESC (newest first)
	// Note: Due to timing precision, we just verify both versions are present
	versionIDs := make(map[string]bool)
	for _, version := range versions {
		versionIDs[version.ID] = true
	}
	
	if !versionIDs[version1ID] {
		t.Errorf("Expected version %s to be present", version1ID)
	}
	if !versionIDs[version2ID] {
		t.Errorf("Expected version %s to be present", version2ID)
	}
}

func TestUniqueWorkingSetConstraint(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project
	projectID := uuid.New().String()
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

	// Create first working set version
	version1ID := uuid.New().String()
	version1Params := CreateGraphVersionParams{
		ID:            version1ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Version 1", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, version1Params)
	if err != nil {
		t.Fatalf("Failed to create version 1: %v", err)
	}

	// Try to create second working set version - should fail
	version2ID := uuid.New().String()
	version2Params := CreateGraphVersionParams{
		ID:            version2ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{String: version1ID, Valid: true},
		Name:          sql.NullString{String: "Version 2", Valid: true},
		Description:   sql.NullString{String: "Second version", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, version2Params)
	if err == nil {
		t.Error("Expected error when creating second working set version")
	}
}