package types

import (
	"encoding/json"
	"time"
)

// EntityType represents the different types of narrative entities
type EntityType string

const (
	EntityTypeScene     EntityType = "Scene"
	EntityTypeCharacter EntityType = "Character"
	EntityTypeLocation  EntityType = "Location"
	EntityTypeTheme     EntityType = "Theme"
	EntityTypePlotPoint EntityType = "PlotPoint"
	EntityTypeArc       EntityType = "Arc"
)

// RelationshipType represents the different types of relationships between entities
type RelationshipType string

const (
	RelationshipContains    RelationshipType = "contains"
	RelationshipAdvances    RelationshipType = "advances"
	RelationshipFeatures    RelationshipType = "features"
	RelationshipOccursAt    RelationshipType = "occurs_at"
	RelationshipInfluences  RelationshipType = "influences"
	RelationshipPrecedes    RelationshipType = "precedes"
	RelationshipFollows     RelationshipType = "follows"
	RelationshipConflicts   RelationshipType = "conflicts"
	RelationshipSupports    RelationshipType = "supports"
)

// AnnotationType represents the different types of annotations
type AnnotationType string

const (
	AnnotationEmotionalAnalysis AnnotationType = "emotional_analysis"
	AnnotationThematicScore     AnnotationType = "thematic_score"
	AnnotationContinuityCheck   AnnotationType = "continuity_check"
	AnnotationStructuralNote    AnnotationType = "structural_note"
	AnnotationCharacterVoice    AnnotationType = "character_voice"
	AnnotationPacingAnalysis    AnnotationType = "pacing_analysis"
)

// SceneData represents the data structure for Scene entities
type SceneData struct {
	Title         string            `json:"title"`
	Summary       string            `json:"summary"`
	Content       string            `json:"content"`
	Act           string            `json:"act,omitempty"`
	Sequence      int               `json:"sequence,omitempty"`
	EmotionalTone string            `json:"emotional_tone,omitempty"`
	Pacing        string            `json:"pacing,omitempty"`
	Characters    []string          `json:"characters,omitempty"`    // Entity IDs
	Location      string            `json:"location,omitempty"`      // Entity ID
	Themes        []string          `json:"themes,omitempty"`        // Entity IDs
	Metadata      map[string]any    `json:"metadata,omitempty"`
}

// CharacterData represents the data structure for Character entities
type CharacterData struct {
	Name                string                 `json:"name"`
	Role                string                 `json:"role,omitempty"`
	Description         string                 `json:"description,omitempty"`
	PersonalityTraits   []string               `json:"personality_traits,omitempty"`
	Background          string                 `json:"background,omitempty"`
	VoiceCharacteristics VoiceCharacteristics  `json:"voice_characteristics,omitempty"`
	Relationships       []CharacterRelationship `json:"relationships,omitempty"`
	CharacterArc        CharacterArc           `json:"character_arc,omitempty"`
}

type VoiceCharacteristics struct {
	Tone           string   `json:"tone,omitempty"`
	Vocabulary     string   `json:"vocabulary,omitempty"`
	SpeechPatterns []string `json:"speech_patterns,omitempty"`
}

type CharacterRelationship struct {
	CharacterID      string `json:"character_id"`
	RelationshipType string `json:"relationship_type"`
	Status           string `json:"status,omitempty"`
}

type CharacterArc struct {
	StartingState string `json:"starting_state,omitempty"`
	CurrentState  string `json:"current_state,omitempty"`
	TargetState   string `json:"target_state,omitempty"`
}

// LocationData represents the data structure for Location entities
type LocationData struct {
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	Atmosphere         string            `json:"atmosphere,omitempty"`
	PhysicalDetails    PhysicalDetails   `json:"physical_details,omitempty"`
	Significance       string            `json:"significance,omitempty"`
	ConnectedLocations []string          `json:"connected_locations,omitempty"` // Entity IDs
}

type PhysicalDetails struct {
	Size            string   `json:"size,omitempty"`
	Lighting        string   `json:"lighting,omitempty"`
	NotableFeatures []string `json:"notable_features,omitempty"`
}

// ThemeData represents the data structure for Theme entities
type ThemeData struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Questions   []string `json:"questions,omitempty"`
	Symbols     []string `json:"symbols,omitempty"`
	Relevance   float64  `json:"relevance,omitempty"` // 0.0 to 1.0
}

// PlotPointData represents the data structure for PlotPoint entities
type PlotPointData struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Type        string   `json:"type,omitempty"` // inciting_incident, plot_twist, climax, etc.
	Act         string   `json:"act,omitempty"`
	Sequence    int      `json:"sequence,omitempty"`
	Characters  []string `json:"characters,omitempty"` // Entity IDs
	Themes      []string `json:"themes,omitempty"`     // Entity IDs
}

// ArcData represents the data structure for Arc entities
type ArcData struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Type        string   `json:"type,omitempty"` // character_arc, plot_arc, thematic_arc
	StartAct    string   `json:"start_act,omitempty"`
	EndAct      string   `json:"end_act,omitempty"`
	Characters  []string `json:"characters,omitempty"` // Entity IDs
	PlotPoints  []string `json:"plot_points,omitempty"` // Entity IDs
	Status      string   `json:"status,omitempty"` // planned, active, completed
}

// EmotionalAnalysisData represents emotional analysis annotation data
type EmotionalAnalysisData struct {
	Sentiment      float64            `json:"sentiment"`       // -1.0 to 1.0
	Emotions       map[string]float64 `json:"emotions"`        // emotion -> intensity
	EmotionalArc   string             `json:"emotional_arc"`   // rising, falling, stable
	ImpactScore    float64            `json:"impact_score"`    // 0.0 to 1.0
	Suggestions    []string           `json:"suggestions,omitempty"`
	AnalyzedAt     time.Time          `json:"analyzed_at"`
}

// ThematicScoreData represents thematic relevance annotation data
type ThematicScoreData struct {
	RelevanceScore float64           `json:"relevance_score"` // 0.0 to 1.0
	ThemeAlignment map[string]float64 `json:"theme_alignment"` // theme_id -> alignment score
	Contributions  []string          `json:"contributions,omitempty"`
	Concerns       []string          `json:"concerns,omitempty"`
	AnalyzedAt     time.Time         `json:"analyzed_at"`
}

// ContinuityCheckData represents continuity validation annotation data
type ContinuityCheckData struct {
	IsConsistent   bool                   `json:"is_consistent"`
	Violations     []ContinuityViolation  `json:"violations,omitempty"`
	Validations    []ContinuityValidation `json:"validations,omitempty"`
	CheckedAt      time.Time              `json:"checked_at"`
}

type ContinuityViolation struct {
	Type        string `json:"type"`        // timeline, character_knowledge, physical, etc.
	Description string `json:"description"`
	Severity    string `json:"severity"`    // low, medium, high, critical
	EntityIDs   []string `json:"entity_ids,omitempty"` // Related entities
}

type ContinuityValidation struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	EntityIDs   []string `json:"entity_ids,omitempty"`
}

// Helper functions to marshal/unmarshal entity data

func MarshalEntityData(data any) (json.RawMessage, error) {
	return json.Marshal(data)
}

func UnmarshalSceneData(raw json.RawMessage) (*SceneData, error) {
	var data SceneData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalCharacterData(raw json.RawMessage) (*CharacterData, error) {
	var data CharacterData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalLocationData(raw json.RawMessage) (*LocationData, error) {
	var data LocationData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalThemeData(raw json.RawMessage) (*ThemeData, error) {
	var data ThemeData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalPlotPointData(raw json.RawMessage) (*PlotPointData, error) {
	var data PlotPointData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalArcData(raw json.RawMessage) (*ArcData, error) {
	var data ArcData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalEmotionalAnalysisData(raw json.RawMessage) (*EmotionalAnalysisData, error) {
	var data EmotionalAnalysisData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalThematicScoreData(raw json.RawMessage) (*ThematicScoreData, error) {
	var data ThematicScoreData
	err := json.Unmarshal(raw, &data)
	return &data, err
}

func UnmarshalContinuityCheckData(raw json.RawMessage) (*ContinuityCheckData, error) {
	var data ContinuityCheckData
	err := json.Unmarshal(raw, &data)
	return &data, err
}