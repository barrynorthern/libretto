package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/barrynorthern/libretto/internal/types"
	"github.com/google/uuid"
)

// TestFullWorkflow tests a complete workflow from project creation to entity relationships
func TestFullWorkflow(t *testing.T) {
	queries := setupTestDB(t)
	ctx := context.Background()

	// 1. Create a project
	projectID := uuid.New().String()
	projectParams := CreateProjectParams{
		ID:          projectID,
		Name:        "Epic Fantasy Adventure",
		Theme:       sql.NullString{String: "Good vs Evil", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A tale of heroes and villains", Valid: true},
	}

	project, err := queries.CreateProject(ctx, projectParams)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// 2. Create a working set version
	versionID := uuid.New().String()
	versionParams := CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Draft", Valid: true},
		Description:   sql.NullString{String: "First version of the story", Valid: true},
		IsWorkingSet:  true,
	}

	version, err := queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		t.Fatalf("Failed to create graph version: %v", err)
	}

	// 3. Create entities using typed data structures
	
	// Create a scene
	sceneID := uuid.New().String()
	sceneData := &types.SceneData{
		Title:         "The Dark Forest",
		Summary:       "Heroes enter the mysterious dark forest",
		Content:       "The ancient trees loomed overhead as our heroes stepped into the shadows...",
		Act:           "Act1",
		Sequence:      1,
		EmotionalTone: "mysterious",
		Pacing:        "slow",
		Characters:    []string{}, // Will be filled after creating characters
		Location:      "",         // Will be filled after creating location
		Themes:        []string{}, // Will be filled after creating themes
		Metadata: map[string]any{
			"word_count":         250,
			"emotional_score":    0.6,
			"thematic_relevance": 0.8,
			"generated_by":       "plot_weaver",
		},
	}

	sceneDataJSON, err := types.MarshalEntityData(sceneData)
	if err != nil {
		t.Fatalf("Failed to marshal scene data: %v", err)
	}

	sceneParams := CreateEntityParams{
		ID:         sceneID,
		VersionID:  versionID,
		EntityType: string(types.EntityTypeScene),
		Name:       sceneData.Title,
		Data:       sceneDataJSON,
	}

	scene, err := queries.CreateEntity(ctx, sceneParams)
	if err != nil {
		t.Fatalf("Failed to create scene entity: %v", err)
	}

	// Create a character
	characterID := uuid.New().String()
	characterData := &types.CharacterData{
		Name:        "Elara the Brave",
		Role:        "protagonist",
		Description: "A courageous warrior with a mysterious past",
		PersonalityTraits: []string{"brave", "determined", "compassionate"},
		Background:        "Born in the northern kingdoms, trained by the ancient order...",
		VoiceCharacteristics: types.VoiceCharacteristics{
			Tone:           "confident",
			Vocabulary:     "formal",
			SpeechPatterns: []string{"uses archaic terms", "speaks with authority"},
		},
		CharacterArc: types.CharacterArc{
			StartingState: "naive_hero",
			CurrentState:  "questioning_beliefs",
			TargetState:   "wise_leader",
		},
	}

	characterDataJSON, err := types.MarshalEntityData(characterData)
	if err != nil {
		t.Fatalf("Failed to marshal character data: %v", err)
	}

	characterParams := CreateEntityParams{
		ID:         characterID,
		VersionID:  versionID,
		EntityType: string(types.EntityTypeCharacter),
		Name:       characterData.Name,
		Data:       characterDataJSON,
	}

	_, err = queries.CreateEntity(ctx, characterParams)
	if err != nil {
		t.Fatalf("Failed to create character entity: %v", err)
	}

	// Create a location
	locationID := uuid.New().String()
	locationData := &types.LocationData{
		Name:        "The Dark Forest",
		Description: "An ancient woodland shrouded in mystery and danger",
		Atmosphere:  "ominous",
		PhysicalDetails: types.PhysicalDetails{
			Size:     "vast",
			Lighting: "dim filtered sunlight",
			NotableFeatures: []string{"towering ancient trees", "twisted paths", "strange sounds"},
		},
		Significance:       "Gateway to the forbidden realm",
		ConnectedLocations: []string{}, // Could connect to other locations
	}

	locationDataJSON, err := types.MarshalEntityData(locationData)
	if err != nil {
		t.Fatalf("Failed to marshal location data: %v", err)
	}

	locationParams := CreateEntityParams{
		ID:         locationID,
		VersionID:  versionID,
		EntityType: string(types.EntityTypeLocation),
		Name:       locationData.Name,
		Data:       locationDataJSON,
	}

	_, err = queries.CreateEntity(ctx, locationParams)
	if err != nil {
		t.Fatalf("Failed to create location entity: %v", err)
	}

	// Create a theme
	themeID := uuid.New().String()
	themeData := &types.ThemeData{
		Name:        "Courage in Darkness",
		Description: "The theme of finding inner strength when facing the unknown",
		Questions:   []string{"What defines true courage?", "How do we face our fears?"},
		Symbols:     []string{"light", "shadows", "ancient trees"},
		Relevance:   0.9,
	}

	themeDataJSON, err := types.MarshalEntityData(themeData)
	if err != nil {
		t.Fatalf("Failed to marshal theme data: %v", err)
	}

	themeParams := CreateEntityParams{
		ID:         themeID,
		VersionID:  versionID,
		EntityType: string(types.EntityTypeTheme),
		Name:       themeData.Name,
		Data:       themeDataJSON,
	}

	_, err = queries.CreateEntity(ctx, themeParams)
	if err != nil {
		t.Fatalf("Failed to create theme entity: %v", err)
	}

	// 4. Create relationships between entities
	
	// Scene features Character
	sceneCharacterRelID := uuid.New().String()
	sceneCharacterProps := map[string]any{
		"role":       "protagonist",
		"importance": "primary",
	}
	sceneCharacterPropsJSON, _ := json.Marshal(sceneCharacterProps)

	sceneCharacterRelParams := CreateRelationshipParams{
		ID:               sceneCharacterRelID,
		VersionID:        versionID,
		FromEntityID:     sceneID,
		ToEntityID:       characterID,
		RelationshipType: string(types.RelationshipFeatures),
		Properties:       sceneCharacterPropsJSON,
	}

	_, err = queries.CreateRelationship(ctx, sceneCharacterRelParams)
	if err != nil {
		t.Fatalf("Failed to create scene-character relationship: %v", err)
	}

	// Scene occurs at Location
	sceneLocationRelID := uuid.New().String()
	sceneLocationRelParams := CreateRelationshipParams{
		ID:               sceneLocationRelID,
		VersionID:        versionID,
		FromEntityID:     sceneID,
		ToEntityID:       locationID,
		RelationshipType: string(types.RelationshipOccursAt),
		Properties:       json.RawMessage(`{}`),
	}

	_, err = queries.CreateRelationship(ctx, sceneLocationRelParams)
	if err != nil {
		t.Fatalf("Failed to create scene-location relationship: %v", err)
	}

	// Scene explores Theme
	sceneThemeRelID := uuid.New().String()
	sceneThemeRelParams := CreateRelationshipParams{
		ID:               sceneThemeRelID,
		VersionID:        versionID,
		FromEntityID:     sceneID,
		ToEntityID:       themeID,
		RelationshipType: string(types.RelationshipInfluences),
		Properties:       json.RawMessage(`{"strength": "strong"}`),
	}

	_, err = queries.CreateRelationship(ctx, sceneThemeRelParams)
	if err != nil {
		t.Fatalf("Failed to create scene-theme relationship: %v", err)
	}

	// 5. Create annotations
	
	// Emotional analysis annotation
	emotionalAnnotationID := uuid.New().String()
	emotionalData := &types.EmotionalAnalysisData{
		Sentiment:    0.6,
		Emotions:     map[string]float64{"mystery": 0.8, "anticipation": 0.7},
		EmotionalArc: "building",
		ImpactScore:  0.75,
		Suggestions:  []string{"Consider adding more sensory details"},
		AnalyzedAt:   scene.CreatedAt,
	}

	emotionalDataJSON, _ := json.Marshal(emotionalData)

	emotionalAnnotationParams := CreateAnnotationParams{
		ID:             emotionalAnnotationID,
		EntityID:       sceneID,
		AnnotationType: string(types.AnnotationEmotionalAnalysis),
		Content:        "Scene establishes mysterious atmosphere effectively",
		Metadata:       emotionalDataJSON,
		AgentName:      sql.NullString{String: "empath_agent", Valid: true},
	}

	_, err = queries.CreateAnnotation(ctx, emotionalAnnotationParams)
	if err != nil {
		t.Fatalf("Failed to create emotional annotation: %v", err)
	}

	// Thematic score annotation
	thematicAnnotationID := uuid.New().String()
	thematicData := &types.ThematicScoreData{
		RelevanceScore: 0.85,
		ThemeAlignment: map[string]float64{themeID: 0.9},
		Contributions:  []string{"Establishes the courage theme", "Sets up character growth"},
		Concerns:       []string{},
		AnalyzedAt:     scene.CreatedAt,
	}

	thematicDataJSON, _ := json.Marshal(thematicData)

	thematicAnnotationParams := CreateAnnotationParams{
		ID:             thematicAnnotationID,
		EntityID:       sceneID,
		AnnotationType: string(types.AnnotationThematicScore),
		Content:        "Strong thematic alignment with courage theme",
		Metadata:       thematicDataJSON,
		AgentName:      sql.NullString{String: "thematic_steward", Valid: true},
	}

	_, err = queries.CreateAnnotation(ctx, thematicAnnotationParams)
	if err != nil {
		t.Fatalf("Failed to create thematic annotation: %v", err)
	}

	// 6. Verify the complete narrative graph
	
	// Check project exists
	retrievedProject, err := queries.GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to retrieve project: %v", err)
	}
	if retrievedProject.Name != project.Name {
		t.Errorf("Expected project name %s, got %s", project.Name, retrievedProject.Name)
	}

	// Check working set version
	workingSet, err := queries.GetWorkingSetVersion(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to get working set version: %v", err)
	}
	if workingSet.ID != versionID {
		t.Errorf("Expected working set ID %s, got %s", versionID, workingSet.ID)
	}

	// Check entities count
	entities, err := queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		t.Fatalf("Failed to list entities: %v", err)
	}
	if len(entities) != 4 {
		t.Errorf("Expected 4 entities, got %d", len(entities))
	}

	// Check relationships count
	relationships, err := queries.ListRelationshipsByVersion(ctx, versionID)
	if err != nil {
		t.Fatalf("Failed to list relationships: %v", err)
	}
	if len(relationships) != 3 {
		t.Errorf("Expected 3 relationships, got %d", len(relationships))
	}

	// Check annotations count
	annotations, err := queries.ListAnnotationsByEntity(ctx, sceneID)
	if err != nil {
		t.Fatalf("Failed to list annotations: %v", err)
	}
	if len(annotations) != 2 {
		t.Errorf("Expected 2 annotations, got %d", len(annotations))
	}

	// Verify data integrity by unmarshaling and checking content
	retrievedScene, err := queries.GetEntity(ctx, sceneID)
	if err != nil {
		t.Fatalf("Failed to retrieve scene: %v", err)
	}

	retrievedSceneData, err := types.UnmarshalSceneData(retrievedScene.Data)
	if err != nil {
		t.Fatalf("Failed to unmarshal scene data: %v", err)
	}

	if retrievedSceneData.Title != sceneData.Title {
		t.Errorf("Expected scene title %s, got %s", sceneData.Title, retrievedSceneData.Title)
	}
	if retrievedSceneData.Act != sceneData.Act {
		t.Errorf("Expected scene act %s, got %s", sceneData.Act, retrievedSceneData.Act)
	}

	t.Logf("Successfully created complete narrative graph with:")
	t.Logf("- Project: %s", project.Name)
	t.Logf("- Version: %s", version.Name.String)
	t.Logf("- Entities: %d", len(entities))
	t.Logf("- Relationships: %d", len(relationships))
	t.Logf("- Annotations: %d", len(annotations))
}