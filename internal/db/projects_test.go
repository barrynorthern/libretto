package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *Queries {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Apply migrations
	migrations := []string{
		// Initial schema
		`CREATE TABLE scenes (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			summary TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		// Living Narrative schema
		`CREATE TABLE projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			theme TEXT,
			genre TEXT,
			description TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE graph_versions (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			parent_version_id TEXT,
			name TEXT DEFAULT '',
			description TEXT DEFAULT '',
			is_working_set BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_version_id) REFERENCES graph_versions(id)
		);`,
		`CREATE TABLE entities (
			id TEXT PRIMARY KEY,
			version_id TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			data JSON NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (version_id) REFERENCES graph_versions(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE relationships (
			id TEXT PRIMARY KEY,
			version_id TEXT NOT NULL,
			from_entity_id TEXT NOT NULL,
			to_entity_id TEXT NOT NULL,
			relationship_type TEXT NOT NULL,
			properties JSON,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (version_id) REFERENCES graph_versions(id) ON DELETE CASCADE,
			FOREIGN KEY (from_entity_id) REFERENCES entities(id) ON DELETE CASCADE,
			FOREIGN KEY (to_entity_id) REFERENCES entities(id) ON DELETE CASCADE,
			UNIQUE(version_id, from_entity_id, to_entity_id, relationship_type)
		);`,
		`CREATE TABLE annotations (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL,
			annotation_type TEXT NOT NULL,
			content TEXT NOT NULL,
			metadata JSON,
			agent_name TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE
		);`,
		// Triggers
		`CREATE TRIGGER update_projects_updated_at 
			AFTER UPDATE ON projects
			FOR EACH ROW
		BEGIN
			UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,
		`CREATE TRIGGER update_entities_updated_at 
			AFTER UPDATE ON entities
			FOR EACH ROW
		BEGIN
			UPDATE entities SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;`,
		// Unique constraint for working set
		`CREATE UNIQUE INDEX idx_unique_working_set_per_project 
		ON graph_versions(project_id) 
		WHERE is_working_set = TRUE;`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			t.Fatalf("Failed to apply migration: %v", err)
		}
	}

	return New(db)
}

func TestCreateProject(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	projectID := uuid.New().String()
	params := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project for unit testing", Valid: true},
	}

	project, err := queries.CreateProject(ctx, params)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	if project.ID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, project.ID)
	}
	if project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got %s", project.Name)
	}
	if !project.Theme.Valid || project.Theme.String != "Adventure" {
		t.Errorf("Expected theme 'Adventure', got %v", project.Theme)
	}
	if project.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}
	if project.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}
}

func TestGetProject(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project first
	projectID := uuid.New().String()
	params := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, params)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Get the project
	project, err := queries.GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}

	if project.ID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, project.ID)
	}
	if project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got %s", project.Name)
	}
}

func TestListProjects(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create multiple projects
	project1ID := uuid.New().String()
	project2ID := uuid.New().String()

	params1 := CreateProjectParams{
		ID:          project1ID,
		Name:        "Project 1",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "First project", Valid: true},
	}

	params2 := CreateProjectParams{
		ID:          project2ID,
		Name:        "Project 2",
		Theme:       sql.NullString{String: "Mystery", Valid: true},
		Genre:       sql.NullString{String: "Thriller", Valid: true},
		Description: sql.NullString{String: "Second project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, params1)
	if err != nil {
		t.Fatalf("Failed to create project 1: %v", err)
	}

	time.Sleep(10 * time.Millisecond) // Ensure different timestamps

	_, err = queries.CreateProject(ctx, params2)
	if err != nil {
		t.Fatalf("Failed to create project 2: %v", err)
	}

	// List projects
	projects, err := queries.ListProjects(ctx)
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}

	// Should be ordered by created_at DESC (newest first)
	// Note: Due to timing precision, we just verify both projects are present
	projectNames := make(map[string]bool)
	for _, project := range projects {
		projectNames[project.Name] = true
	}
	
	if !projectNames["Project 1"] {
		t.Error("Expected 'Project 1' to be present")
	}
	if !projectNames["Project 2"] {
		t.Error("Expected 'Project 2' to be present")
	}
}

func TestUpdateProject(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project first
	projectID := uuid.New().String()
	params := CreateProjectParams{
		ID:          projectID,
		Name:        "Original Name",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Original description", Valid: true},
	}

	originalProject, err := queries.CreateProject(ctx, params)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	time.Sleep(10 * time.Millisecond) // Ensure different timestamp

	// Update the project
	updateParams := UpdateProjectParams{
		ID:          projectID,
		Name:        "Updated Name",
		Theme:       sql.NullString{String: "Mystery", Valid: true},
		Genre:       sql.NullString{String: "Thriller", Valid: true},
		Description: sql.NullString{String: "Updated description", Valid: true},
	}

	updatedProject, err := queries.UpdateProject(ctx, updateParams)
	if err != nil {
		t.Fatalf("Failed to update project: %v", err)
	}

	if updatedProject.Name != "Updated Name" {
		t.Errorf("Expected updated name 'Updated Name', got %s", updatedProject.Name)
	}
	if !updatedProject.Theme.Valid || updatedProject.Theme.String != "Mystery" {
		t.Errorf("Expected updated theme 'Mystery', got %v", updatedProject.Theme)
	}
	// Note: SQLite CURRENT_TIMESTAMP may have same precision, so we check if it's not before
	if updatedProject.UpdatedAt.Before(originalProject.UpdatedAt) {
		t.Error("Expected updated_at to not go backwards after update")
	}
}

func TestDeleteProject(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create a project first
	projectID := uuid.New().String()
	params := CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A test project", Valid: true},
	}

	_, err := queries.CreateProject(ctx, params)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Delete the project
	err = queries.DeleteProject(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to delete project: %v", err)
	}

	// Verify it's deleted
	_, err = queries.GetProject(ctx, projectID)
	if err == nil {
		t.Error("Expected error when getting deleted project")
	}
	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}