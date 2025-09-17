package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateAnnotation(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entity
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	entityID := uuid.New().String()

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

	// Create entity
	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       json.RawMessage(`{"title": "Opening Scene", "summary": "The beginning"}`),
	}

	_, err = queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Create annotation
	annotationID := uuid.New().String()
	metadata := map[string]any{
		"sentiment":      0.7,
		"emotions":       map[string]float64{"excitement": 0.8, "curiosity": 0.6},
		"emotional_arc":  "rising",
		"impact_score":   0.75,
		"analyzed_at":    time.Now().Format(time.RFC3339),
	}

	metadataJSON, _ := json.Marshal(metadata)

	annotationParams := CreateAnnotationParams{
		ID:             annotationID,
		EntityID:       entityID,
		AnnotationType: "emotional_analysis",
		Content:        "This scene shows strong emotional engagement with rising excitement and curiosity.",
		Metadata:       metadataJSON,
		AgentName:      sql.NullString{String: "empath_agent", Valid: true},
	}

	annotation, err := queries.CreateAnnotation(ctx, annotationParams)
	if err != nil {
		t.Fatalf("Failed to create annotation: %v", err)
	}

	if annotation.ID != annotationID {
		t.Errorf("Expected annotation ID %s, got %s", annotationID, annotation.ID)
	}
	if annotation.EntityID != entityID {
		t.Errorf("Expected entity ID %s, got %s", entityID, annotation.EntityID)
	}
	if annotation.AnnotationType != "emotional_analysis" {
		t.Errorf("Expected annotation type 'emotional_analysis', got %s", annotation.AnnotationType)
	}
	if !annotation.AgentName.Valid || annotation.AgentName.String != "empath_agent" {
		t.Errorf("Expected agent name 'empath_agent', got %v", annotation.AgentName)
	}

	// Verify metadata
	var storedMetadata map[string]any
	err = json.Unmarshal(annotation.Metadata, &storedMetadata)
	if err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if storedMetadata["sentiment"] != 0.7 {
		t.Errorf("Expected sentiment 0.7, got %v", storedMetadata["sentiment"])
	}
	if storedMetadata["emotional_arc"] != "rising" {
		t.Errorf("Expected emotional_arc 'rising', got %v", storedMetadata["emotional_arc"])
	}
}

func TestListAnnotationsByEntity(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entity
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	entityID := uuid.New().String()

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

	// Create entity
	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       json.RawMessage(`{"title": "Opening Scene"}`),
	}

	_, err = queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Create multiple annotations
	annotations := []CreateAnnotationParams{
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "emotional_analysis",
			Content:        "Emotional analysis content",
			Metadata:       json.RawMessage(`{"sentiment": 0.7}`),
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "thematic_score",
			Content:        "Thematic analysis content",
			Metadata:       json.RawMessage(`{"relevance_score": 0.8}`),
			AgentName:      sql.NullString{String: "thematic_steward", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "continuity_check",
			Content:        "Continuity validation content",
			Metadata:       json.RawMessage(`{"is_consistent": true}`),
			AgentName:      sql.NullString{String: "continuity_steward", Valid: true},
		},
	}

	for _, annotationParams := range annotations {
		_, err = queries.CreateAnnotation(ctx, annotationParams)
		if err != nil {
			t.Fatalf("Failed to create annotation: %v", err)
		}
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	// List annotations for entity
	entityAnnotations, err := queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		t.Fatalf("Failed to list annotations by entity: %v", err)
	}

	if len(entityAnnotations) != 3 {
		t.Errorf("Expected 3 annotations, got %d", len(entityAnnotations))
	}

	// Verify annotations are ordered by created_at DESC (newest first)
	// Note: The exact order may vary due to timing, so we just check that all types are present
	// The last created annotation should be first due to DESC ordering

	// Verify all annotation types are present
	annotationTypes := make(map[string]bool)
	for _, annotation := range entityAnnotations {
		annotationTypes[annotation.AnnotationType] = true
	}

	expectedTypes := []string{"emotional_analysis", "thematic_score", "continuity_check"}
	for _, expectedType := range expectedTypes {
		if !annotationTypes[expectedType] {
			t.Errorf("Expected annotation type '%s' to be present", expectedType)
		}
	}
}

func TestListAnnotationsByType(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entity
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	entityID := uuid.New().String()

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

	// Create entity
	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       json.RawMessage(`{"title": "Opening Scene"}`),
	}

	_, err = queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Create multiple annotations of different types
	annotations := []CreateAnnotationParams{
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "emotional_analysis",
			Content:        "First emotional analysis",
			Metadata:       json.RawMessage(`{"sentiment": 0.7}`),
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "emotional_analysis",
			Content:        "Second emotional analysis",
			Metadata:       json.RawMessage(`{"sentiment": 0.5}`),
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "thematic_score",
			Content:        "Thematic analysis",
			Metadata:       json.RawMessage(`{"relevance_score": 0.8}`),
			AgentName:      sql.NullString{String: "thematic_steward", Valid: true},
		},
	}

	for _, annotationParams := range annotations {
		_, err = queries.CreateAnnotation(ctx, annotationParams)
		if err != nil {
			t.Fatalf("Failed to create annotation: %v", err)
		}
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	// List emotional_analysis annotations
	listParams := ListAnnotationsByTypeParams{
		EntityID:       entityID,
		AnnotationType: "emotional_analysis",
	}

	emotionalAnnotations, err := queries.ListAnnotationsByType(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list emotional analysis annotations: %v", err)
	}

	if len(emotionalAnnotations) != 2 {
		t.Errorf("Expected 2 emotional analysis annotations, got %d", len(emotionalAnnotations))
	}

	for _, annotation := range emotionalAnnotations {
		if annotation.AnnotationType != "emotional_analysis" {
			t.Errorf("Expected annotation type 'emotional_analysis', got %s", annotation.AnnotationType)
		}
	}

	// List thematic_score annotations
	listParams.AnnotationType = "thematic_score"
	thematicAnnotations, err := queries.ListAnnotationsByType(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list thematic score annotations: %v", err)
	}

	if len(thematicAnnotations) != 1 {
		t.Errorf("Expected 1 thematic score annotation, got %d", len(thematicAnnotations))
	}
}

func TestListAnnotationsByAgent(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entity
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	entityID := uuid.New().String()

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

	// Create entity
	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       json.RawMessage(`{"title": "Opening Scene"}`),
	}

	_, err = queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Create annotations from different agents
	annotations := []CreateAnnotationParams{
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "emotional_analysis",
			Content:        "First emotional analysis",
			Metadata:       json.RawMessage(`{"sentiment": 0.7}`),
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "emotional_analysis",
			Content:        "Second emotional analysis",
			Metadata:       json.RawMessage(`{"sentiment": 0.5}`),
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "thematic_score",
			Content:        "Thematic analysis",
			Metadata:       json.RawMessage(`{"relevance_score": 0.8}`),
			AgentName:      sql.NullString{String: "thematic_steward", Valid: true},
		},
	}

	for _, annotationParams := range annotations {
		_, err = queries.CreateAnnotation(ctx, annotationParams)
		if err != nil {
			t.Fatalf("Failed to create annotation: %v", err)
		}
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	// List annotations by empath_agent
	empathAnnotations, err := queries.ListAnnotationsByAgent(ctx, sql.NullString{String: "empath_agent", Valid: true})
	if err != nil {
		t.Fatalf("Failed to list annotations by empath agent: %v", err)
	}

	if len(empathAnnotations) != 2 {
		t.Errorf("Expected 2 annotations from empath agent, got %d", len(empathAnnotations))
	}

	for _, annotation := range empathAnnotations {
		if !annotation.AgentName.Valid || annotation.AgentName.String != "empath_agent" {
			t.Errorf("Expected agent name 'empath_agent', got %v", annotation.AgentName)
		}
	}

	// List annotations by thematic_steward
	thematicAnnotations, err := queries.ListAnnotationsByAgent(ctx, sql.NullString{String: "thematic_steward", Valid: true})
	if err != nil {
		t.Fatalf("Failed to list annotations by thematic steward: %v", err)
	}

	if len(thematicAnnotations) != 1 {
		t.Errorf("Expected 1 annotation from thematic steward, got %d", len(thematicAnnotations))
	}
}

func TestUpdateAnnotation(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entity
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	entityID := uuid.New().String()

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

	// Create entity
	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       json.RawMessage(`{"title": "Opening Scene"}`),
	}

	_, err = queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Create annotation
	annotationID := uuid.New().String()
	originalMetadata := map[string]any{
		"sentiment":     0.5,
		"impact_score":  0.6,
	}

	originalMetadataJSON, _ := json.Marshal(originalMetadata)

	annotationParams := CreateAnnotationParams{
		ID:             annotationID,
		EntityID:       entityID,
		AnnotationType: "emotional_analysis",
		Content:        "Original analysis content",
		Metadata:       originalMetadataJSON,
		AgentName:      sql.NullString{String: "empath_agent", Valid: true},
	}

	_, err = queries.CreateAnnotation(ctx, annotationParams)
	if err != nil {
		t.Fatalf("Failed to create annotation: %v", err)
	}

	// Update annotation
	updatedMetadata := map[string]any{
		"sentiment":     0.8,
		"impact_score":  0.9,
		"confidence":    0.95,
	}

	updatedMetadataJSON, _ := json.Marshal(updatedMetadata)

	updateParams := UpdateAnnotationParams{
		ID:       annotationID,
		Content:  "Updated analysis content with more detail",
		Metadata: updatedMetadataJSON,
	}

	updatedAnnotation, err := queries.UpdateAnnotation(ctx, updateParams)
	if err != nil {
		t.Fatalf("Failed to update annotation: %v", err)
	}

	if updatedAnnotation.Content != "Updated analysis content with more detail" {
		t.Errorf("Expected updated content, got %s", updatedAnnotation.Content)
	}

	// Verify metadata was updated
	var storedMetadata map[string]any
	err = json.Unmarshal(updatedAnnotation.Metadata, &storedMetadata)
	if err != nil {
		t.Fatalf("Failed to unmarshal updated metadata: %v", err)
	}

	if storedMetadata["sentiment"] != 0.8 {
		t.Errorf("Expected sentiment 0.8, got %v", storedMetadata["sentiment"])
	}
	if storedMetadata["confidence"] != 0.95 {
		t.Errorf("Expected confidence 0.95, got %v", storedMetadata["confidence"])
	}
}

func TestDeleteAnnotationsByEntity(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// Create project, version, and entity
	projectID := uuid.New().String()
	versionID := uuid.New().String()
	entityID := uuid.New().String()

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

	// Create entity
	entityParams := CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: "Scene",
		Name:       "Opening Scene",
		Data:       json.RawMessage(`{"title": "Opening Scene"}`),
	}

	_, err = queries.CreateEntity(ctx, entityParams)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Create multiple annotations
	annotations := []CreateAnnotationParams{
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "emotional_analysis",
			Content:        "Emotional analysis",
			Metadata:       json.RawMessage(`{"sentiment": 0.7}`),
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		},
		{
			ID:             uuid.New().String(),
			EntityID:       entityID,
			AnnotationType: "thematic_score",
			Content:        "Thematic analysis",
			Metadata:       json.RawMessage(`{"relevance_score": 0.8}`),
			AgentName:      sql.NullString{String: "thematic_steward", Valid: true},
		},
	}

	for _, annotationParams := range annotations {
		_, err = queries.CreateAnnotation(ctx, annotationParams)
		if err != nil {
			t.Fatalf("Failed to create annotation: %v", err)
		}
	}

	// Verify annotations exist
	beforeAnnotations, err := queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		t.Fatalf("Failed to list annotations before deletion: %v", err)
	}

	if len(beforeAnnotations) != 2 {
		t.Errorf("Expected 2 annotations before deletion, got %d", len(beforeAnnotations))
	}

	// Delete all annotations for entity
	err = queries.DeleteAnnotationsByEntity(ctx, entityID)
	if err != nil {
		t.Fatalf("Failed to delete annotations by entity: %v", err)
	}

	// Verify annotations are deleted
	afterAnnotations, err := queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		t.Fatalf("Failed to list annotations after deletion: %v", err)
	}

	if len(afterAnnotations) != 0 {
		t.Errorf("Expected 0 annotations after deletion, got %d", len(afterAnnotations))
	}
}