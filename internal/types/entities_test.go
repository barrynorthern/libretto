package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSceneDataMarshalUnmarshal(t *testing.T) {
	original := &SceneData{
		Title:         "The Great Confrontation",
		Summary:       "Hero faces the villain in final battle",
		Content:       "The wind howled as they faced each other...",
		Act:           "Act3",
		Sequence:      25,
		EmotionalTone: "intense",
		Pacing:        "fast",
		Characters:    []string{"char_001", "char_002"},
		Location:      "loc_001",
		Themes:        []string{"theme_001", "theme_002"},
		Metadata: map[string]any{
			"word_count":         1500,
			"emotional_score":    0.9,
			"thematic_relevance": 0.85,
			"generated_by":       "plot_weaver",
		},
	}

	// Marshal to JSON
	data, err := MarshalEntityData(original)
	if err != nil {
		t.Fatalf("Failed to marshal scene data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalSceneData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal scene data: %v", err)
	}

	// Verify fields
	if unmarshaled.Title != original.Title {
		t.Errorf("Expected title %s, got %s", original.Title, unmarshaled.Title)
	}
	if unmarshaled.Act != original.Act {
		t.Errorf("Expected act %s, got %s", original.Act, unmarshaled.Act)
	}
	if unmarshaled.Sequence != original.Sequence {
		t.Errorf("Expected sequence %d, got %d", original.Sequence, unmarshaled.Sequence)
	}
	if len(unmarshaled.Characters) != len(original.Characters) {
		t.Errorf("Expected %d characters, got %d", len(original.Characters), len(unmarshaled.Characters))
	}
	if unmarshaled.Metadata["word_count"] != float64(1500) { // JSON numbers become float64
		t.Errorf("Expected word_count 1500, got %v", unmarshaled.Metadata["word_count"])
	}
}

func TestCharacterDataMarshalUnmarshal(t *testing.T) {
	original := &CharacterData{
		Name:        "Elena Vasquez",
		Role:        "protagonist",
		Description: "Determined archaeologist seeking truth",
		PersonalityTraits: []string{"curious", "stubborn", "loyal"},
		Background:        "Born in Madrid, studied at Oxford...",
		VoiceCharacteristics: VoiceCharacteristics{
			Tone:           "direct",
			Vocabulary:     "academic",
			SpeechPatterns: []string{"uses metaphors", "asks probing questions"},
		},
		Relationships: []CharacterRelationship{
			{
				CharacterID:      "char_002",
				RelationshipType: "mentor_student",
				Status:           "strained",
			},
		},
		CharacterArc: CharacterArc{
			StartingState: "naive_trust",
			CurrentState:  "growing_suspicion",
			TargetState:   "hardened_wisdom",
		},
	}

	// Marshal to JSON
	data, err := MarshalEntityData(original)
	if err != nil {
		t.Fatalf("Failed to marshal character data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalCharacterData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal character data: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, unmarshaled.Name)
	}
	if unmarshaled.Role != original.Role {
		t.Errorf("Expected role %s, got %s", original.Role, unmarshaled.Role)
	}
	if len(unmarshaled.PersonalityTraits) != len(original.PersonalityTraits) {
		t.Errorf("Expected %d personality traits, got %d", len(original.PersonalityTraits), len(unmarshaled.PersonalityTraits))
	}
	if unmarshaled.VoiceCharacteristics.Tone != original.VoiceCharacteristics.Tone {
		t.Errorf("Expected tone %s, got %s", original.VoiceCharacteristics.Tone, unmarshaled.VoiceCharacteristics.Tone)
	}
	if len(unmarshaled.Relationships) != len(original.Relationships) {
		t.Errorf("Expected %d relationships, got %d", len(original.Relationships), len(unmarshaled.Relationships))
	}
	if unmarshaled.CharacterArc.StartingState != original.CharacterArc.StartingState {
		t.Errorf("Expected starting state %s, got %s", original.CharacterArc.StartingState, unmarshaled.CharacterArc.StartingState)
	}
}

func TestLocationDataMarshalUnmarshal(t *testing.T) {
	original := &LocationData{
		Name:        "Ancient Library Vault",
		Description: "Hidden chamber beneath university library",
		Atmosphere:  "mysterious",
		PhysicalDetails: PhysicalDetails{
			Size:     "small chamber",
			Lighting: "dim candlelight",
			NotableFeatures: []string{"stone walls", "ancient texts", "hidden passages"},
		},
		Significance:       "Contains evidence of mentor's deception",
		ConnectedLocations: []string{"loc_002", "loc_003"},
	}

	// Marshal to JSON
	data, err := MarshalEntityData(original)
	if err != nil {
		t.Fatalf("Failed to marshal location data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalLocationData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal location data: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, unmarshaled.Name)
	}
	if unmarshaled.Atmosphere != original.Atmosphere {
		t.Errorf("Expected atmosphere %s, got %s", original.Atmosphere, unmarshaled.Atmosphere)
	}
	if unmarshaled.PhysicalDetails.Size != original.PhysicalDetails.Size {
		t.Errorf("Expected size %s, got %s", original.PhysicalDetails.Size, unmarshaled.PhysicalDetails.Size)
	}
	if len(unmarshaled.PhysicalDetails.NotableFeatures) != len(original.PhysicalDetails.NotableFeatures) {
		t.Errorf("Expected %d notable features, got %d", len(original.PhysicalDetails.NotableFeatures), len(unmarshaled.PhysicalDetails.NotableFeatures))
	}
	if len(unmarshaled.ConnectedLocations) != len(original.ConnectedLocations) {
		t.Errorf("Expected %d connected locations, got %d", len(original.ConnectedLocations), len(unmarshaled.ConnectedLocations))
	}
}

func TestThemeDataMarshalUnmarshal(t *testing.T) {
	original := &ThemeData{
		Name:        "Truth vs Deception",
		Description: "The struggle between revealing and concealing truth",
		Questions:   []string{"What is the cost of truth?", "Can deception ever be justified?"},
		Symbols:     []string{"mirrors", "masks", "light and shadow"},
		Relevance:   0.85,
	}

	// Marshal to JSON
	data, err := MarshalEntityData(original)
	if err != nil {
		t.Fatalf("Failed to marshal theme data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalThemeData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal theme data: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, unmarshaled.Name)
	}
	if unmarshaled.Description != original.Description {
		t.Errorf("Expected description %s, got %s", original.Description, unmarshaled.Description)
	}
	if len(unmarshaled.Questions) != len(original.Questions) {
		t.Errorf("Expected %d questions, got %d", len(original.Questions), len(unmarshaled.Questions))
	}
	if len(unmarshaled.Symbols) != len(original.Symbols) {
		t.Errorf("Expected %d symbols, got %d", len(original.Symbols), len(unmarshaled.Symbols))
	}
	if unmarshaled.Relevance != original.Relevance {
		t.Errorf("Expected relevance %f, got %f", original.Relevance, unmarshaled.Relevance)
	}
}

func TestPlotPointDataMarshalUnmarshal(t *testing.T) {
	original := &PlotPointData{
		Name:        "The Revelation",
		Description: "Hero discovers mentor's true identity",
		Type:        "plot_twist",
		Act:         "Act2",
		Sequence:    15,
		Characters:  []string{"char_001", "char_002"},
		Themes:      []string{"theme_001"},
	}

	// Marshal to JSON
	data, err := MarshalEntityData(original)
	if err != nil {
		t.Fatalf("Failed to marshal plot point data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalPlotPointData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal plot point data: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, unmarshaled.Name)
	}
	if unmarshaled.Type != original.Type {
		t.Errorf("Expected type %s, got %s", original.Type, unmarshaled.Type)
	}
	if unmarshaled.Sequence != original.Sequence {
		t.Errorf("Expected sequence %d, got %d", original.Sequence, unmarshaled.Sequence)
	}
	if len(unmarshaled.Characters) != len(original.Characters) {
		t.Errorf("Expected %d characters, got %d", len(original.Characters), len(unmarshaled.Characters))
	}
}

func TestArcDataMarshalUnmarshal(t *testing.T) {
	original := &ArcData{
		Name:        "Hero's Journey",
		Description: "The protagonist's transformation from naive to wise",
		Type:        "character_arc",
		StartAct:    "Act1",
		EndAct:      "Act3",
		Characters:  []string{"char_001"},
		PlotPoints:  []string{"plot_001", "plot_002", "plot_003"},
		Status:      "active",
	}

	// Marshal to JSON
	data, err := MarshalEntityData(original)
	if err != nil {
		t.Fatalf("Failed to marshal arc data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalArcData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal arc data: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, unmarshaled.Name)
	}
	if unmarshaled.Type != original.Type {
		t.Errorf("Expected type %s, got %s", original.Type, unmarshaled.Type)
	}
	if unmarshaled.Status != original.Status {
		t.Errorf("Expected status %s, got %s", original.Status, unmarshaled.Status)
	}
	if len(unmarshaled.PlotPoints) != len(original.PlotPoints) {
		t.Errorf("Expected %d plot points, got %d", len(original.PlotPoints), len(unmarshaled.PlotPoints))
	}
}

func TestEmotionalAnalysisDataMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	original := &EmotionalAnalysisData{
		Sentiment: 0.7,
		Emotions: map[string]float64{
			"excitement": 0.8,
			"curiosity":  0.6,
			"tension":    0.4,
		},
		EmotionalArc: "rising",
		ImpactScore:  0.75,
		Suggestions:  []string{"Consider adding more tension", "Enhance emotional contrast"},
		AnalyzedAt:   now,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal emotional analysis data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalEmotionalAnalysisData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal emotional analysis data: %v", err)
	}

	// Verify fields
	if unmarshaled.Sentiment != original.Sentiment {
		t.Errorf("Expected sentiment %f, got %f", original.Sentiment, unmarshaled.Sentiment)
	}
	if unmarshaled.EmotionalArc != original.EmotionalArc {
		t.Errorf("Expected emotional arc %s, got %s", original.EmotionalArc, unmarshaled.EmotionalArc)
	}
	if unmarshaled.ImpactScore != original.ImpactScore {
		t.Errorf("Expected impact score %f, got %f", original.ImpactScore, unmarshaled.ImpactScore)
	}
	if len(unmarshaled.Emotions) != len(original.Emotions) {
		t.Errorf("Expected %d emotions, got %d", len(original.Emotions), len(unmarshaled.Emotions))
	}
	if unmarshaled.Emotions["excitement"] != original.Emotions["excitement"] {
		t.Errorf("Expected excitement %f, got %f", original.Emotions["excitement"], unmarshaled.Emotions["excitement"])
	}
	if len(unmarshaled.Suggestions) != len(original.Suggestions) {
		t.Errorf("Expected %d suggestions, got %d", len(original.Suggestions), len(unmarshaled.Suggestions))
	}
}

func TestThematicScoreDataMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	original := &ThematicScoreData{
		RelevanceScore: 0.85,
		ThemeAlignment: map[string]float64{
			"theme_001": 0.9,
			"theme_002": 0.7,
		},
		Contributions: []string{"Reinforces central conflict", "Develops character motivation"},
		Concerns:      []string{"May be too heavy-handed"},
		AnalyzedAt:    now,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal thematic score data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalThematicScoreData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal thematic score data: %v", err)
	}

	// Verify fields
	if unmarshaled.RelevanceScore != original.RelevanceScore {
		t.Errorf("Expected relevance score %f, got %f", original.RelevanceScore, unmarshaled.RelevanceScore)
	}
	if len(unmarshaled.ThemeAlignment) != len(original.ThemeAlignment) {
		t.Errorf("Expected %d theme alignments, got %d", len(original.ThemeAlignment), len(unmarshaled.ThemeAlignment))
	}
	if unmarshaled.ThemeAlignment["theme_001"] != original.ThemeAlignment["theme_001"] {
		t.Errorf("Expected theme_001 alignment %f, got %f", original.ThemeAlignment["theme_001"], unmarshaled.ThemeAlignment["theme_001"])
	}
	if len(unmarshaled.Contributions) != len(original.Contributions) {
		t.Errorf("Expected %d contributions, got %d", len(original.Contributions), len(unmarshaled.Contributions))
	}
	if len(unmarshaled.Concerns) != len(original.Concerns) {
		t.Errorf("Expected %d concerns, got %d", len(original.Concerns), len(unmarshaled.Concerns))
	}
}

func TestContinuityCheckDataMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	original := &ContinuityCheckData{
		IsConsistent: false,
		Violations: []ContinuityViolation{
			{
				Type:        "timeline",
				Description: "Character appears in two places at once",
				Severity:    "high",
				EntityIDs:   []string{"scene_001", "scene_002"},
			},
			{
				Type:        "character_knowledge",
				Description: "Character knows information they shouldn't have",
				Severity:    "medium",
				EntityIDs:   []string{"char_001", "scene_003"},
			},
		},
		Validations: []ContinuityValidation{
			{
				Type:        "physical",
				Description: "Location descriptions are consistent",
				EntityIDs:   []string{"loc_001"},
			},
		},
		CheckedAt: now,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal continuity check data: %v", err)
	}

	// Unmarshal back
	unmarshaled, err := UnmarshalContinuityCheckData(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal continuity check data: %v", err)
	}

	// Verify fields
	if unmarshaled.IsConsistent != original.IsConsistent {
		t.Errorf("Expected is_consistent %t, got %t", original.IsConsistent, unmarshaled.IsConsistent)
	}
	if len(unmarshaled.Violations) != len(original.Violations) {
		t.Errorf("Expected %d violations, got %d", len(original.Violations), len(unmarshaled.Violations))
	}
	if len(unmarshaled.Validations) != len(original.Validations) {
		t.Errorf("Expected %d validations, got %d", len(original.Validations), len(unmarshaled.Validations))
	}

	// Check first violation
	if len(unmarshaled.Violations) > 0 {
		violation := unmarshaled.Violations[0]
		originalViolation := original.Violations[0]
		if violation.Type != originalViolation.Type {
			t.Errorf("Expected violation type %s, got %s", originalViolation.Type, violation.Type)
		}
		if violation.Severity != originalViolation.Severity {
			t.Errorf("Expected violation severity %s, got %s", originalViolation.Severity, violation.Severity)
		}
		if len(violation.EntityIDs) != len(originalViolation.EntityIDs) {
			t.Errorf("Expected %d entity IDs, got %d", len(originalViolation.EntityIDs), len(violation.EntityIDs))
		}
	}
}

func TestEntityTypeConstants(t *testing.T) {
	expectedTypes := []EntityType{
		EntityTypeScene,
		EntityTypeCharacter,
		EntityTypeLocation,
		EntityTypeTheme,
		EntityTypePlotPoint,
		EntityTypeArc,
	}

	expectedValues := []string{
		"Scene",
		"Character",
		"Location",
		"Theme",
		"PlotPoint",
		"Arc",
	}

	for i, entityType := range expectedTypes {
		if string(entityType) != expectedValues[i] {
			t.Errorf("Expected entity type %s, got %s", expectedValues[i], string(entityType))
		}
	}
}

func TestRelationshipTypeConstants(t *testing.T) {
	expectedTypes := []RelationshipType{
		RelationshipContains,
		RelationshipAdvances,
		RelationshipFeatures,
		RelationshipOccursAt,
		RelationshipInfluences,
		RelationshipPrecedes,
		RelationshipFollows,
		RelationshipConflicts,
		RelationshipSupports,
	}

	expectedValues := []string{
		"contains",
		"advances",
		"features",
		"occurs_at",
		"influences",
		"precedes",
		"follows",
		"conflicts",
		"supports",
	}

	for i, relType := range expectedTypes {
		if string(relType) != expectedValues[i] {
			t.Errorf("Expected relationship type %s, got %s", expectedValues[i], string(relType))
		}
	}
}

func TestAnnotationTypeConstants(t *testing.T) {
	expectedTypes := []AnnotationType{
		AnnotationEmotionalAnalysis,
		AnnotationThematicScore,
		AnnotationContinuityCheck,
		AnnotationStructuralNote,
		AnnotationCharacterVoice,
		AnnotationPacingAnalysis,
	}

	expectedValues := []string{
		"emotional_analysis",
		"thematic_score",
		"continuity_check",
		"structural_note",
		"character_voice",
		"pacing_analysis",
	}

	for i, annotationType := range expectedTypes {
		if string(annotationType) != expectedValues[i] {
			t.Errorf("Expected annotation type %s, got %s", expectedValues[i], string(annotationType))
		}
	}
}