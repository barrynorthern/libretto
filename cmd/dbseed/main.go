package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/types"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dbPath = flag.String("db", "libretto.db", "Path to SQLite database")
		preset = flag.String("preset", "fantasy", "Preset to load: fantasy, scifi, mystery")
		clean  = flag.Bool("clean", false, "Clean database before seeding")
	)
	flag.Parse()

	database, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Apply migrations if needed
	if err := applyMigrations(database); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	queries := db.New(database)
	ctx := context.Background()

	if *clean {
		if err := cleanDatabase(database); err != nil {
			log.Fatalf("Failed to clean database: %v", err)
		}
		fmt.Println("Database cleaned.")
	}

	switch *preset {
	case "fantasy":
		seedFantasyStory(ctx, queries)
	case "scifi":
		seedSciFiStory(ctx, queries)
	case "mystery":
		seedMysteryStory(ctx, queries)
	default:
		log.Fatalf("Unknown preset: %s", *preset)
	}

	fmt.Printf("Database seeded with %s preset.\n", *preset)
}

func applyMigrations(database *sql.DB) error {
	migrations := []string{
		// Initial schema
		`CREATE TABLE IF NOT EXISTS scenes (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			summary TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		// Living Narrative schema
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			theme TEXT,
			genre TEXT,
			description TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS graph_versions (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			parent_version_id TEXT,
			name TEXT DEFAULT '',
			description TEXT DEFAULT '',
			is_working_set BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_version_id) REFERENCES graph_versions(id)
		)`,
		`CREATE TABLE IF NOT EXISTS entities (
			id TEXT PRIMARY KEY,
			version_id TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			data JSON NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (version_id) REFERENCES graph_versions(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS relationships (
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
		)`,
		`CREATE TABLE IF NOT EXISTS annotations (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL,
			annotation_type TEXT NOT NULL,
			content TEXT NOT NULL,
			metadata JSON,
			agent_name TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE
		)`,
	}

	for _, migration := range migrations {
		if _, err := database.Exec(migration); err != nil {
			return fmt.Errorf("failed to apply migration: %v", err)
		}
	}

	return nil
}

func cleanDatabase(database *sql.DB) error {
	tables := []string{"annotations", "relationships", "entities", "graph_versions", "projects", "scenes"}
	
	for _, table := range tables {
		if _, err := database.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to clean table %s: %v", table, err)
		}
	}
	
	return nil
}

func seedFantasyStory(ctx context.Context, queries *db.Queries) {
	// Create project
	projectID := uuid.New().String()
	project := db.CreateProjectParams{
		ID:          projectID,
		Name:        "The Crystal of Eternal Light",
		Theme:       sql.NullString{String: "Good vs Evil", Valid: true},
		Genre:       sql.NullString{String: "Epic Fantasy", Valid: true},
		Description: sql.NullString{String: "A tale of heroes seeking an ancient crystal to save their realm", Valid: true},
	}

	_, err := queries.CreateProject(ctx, project)
	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}

	// Create working set version
	versionID := uuid.New().String()
	version := db.CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "First Draft", Valid: true},
		Description:   sql.NullString{String: "Initial version of the fantasy epic", Valid: true},
		IsWorkingSet:  true,
	}

	_, err = queries.CreateGraphVersion(ctx, version)
	if err != nil {
		log.Fatalf("Failed to create version: %v", err)
	}

	// Create entities
	entities := createFantasyEntities(versionID)
	entityIDs := make(map[string]string) // name -> id mapping

	for _, entity := range entities {
		created, err := queries.CreateEntity(ctx, entity)
		if err != nil {
			log.Fatalf("Failed to create entity %s: %v", entity.Name, err)
		}
		entityIDs[entity.Name] = created.ID
	}

	// Create relationships
	relationships := createFantasyRelationships(versionID, entityIDs)
	for _, rel := range relationships {
		_, err := queries.CreateRelationship(ctx, rel)
		if err != nil {
			log.Fatalf("Failed to create relationship: %v", err)
		}
	}

	// Create annotations
	annotations := createFantasyAnnotations(entityIDs)
	for _, annotation := range annotations {
		_, err := queries.CreateAnnotation(ctx, annotation)
		if err != nil {
			log.Fatalf("Failed to create annotation: %v", err)
		}
	}

	fmt.Printf("Created fantasy story with project ID: %s\n", projectID)
}

func createFantasyEntities(versionID string) []db.CreateEntityParams {
	var entities []db.CreateEntityParams

	// Scenes
	scenes := []struct {
		name string
		data *types.SceneData
	}{
		{
			"The Call to Adventure",
			&types.SceneData{
				Title:         "The Call to Adventure",
				Summary:       "Elara receives the quest to find the Crystal of Eternal Light",
				Content:       "The ancient wizard Gandor appeared at Elara's door with urgent news...",
				Act:           "Act1",
				Sequence:      1,
				EmotionalTone: "mysterious",
				Pacing:        "slow",
				Metadata: map[string]any{
					"word_count":         850,
					"emotional_score":    0.7,
					"thematic_relevance": 0.9,
				},
			},
		},
		{
			"The Dark Forest",
			&types.SceneData{
				Title:         "The Dark Forest",
				Summary:       "Heroes enter the perilous Shadowwood",
				Content:       "Ancient trees loomed overhead as our heroes stepped into darkness...",
				Act:           "Act2",
				Sequence:      8,
				EmotionalTone: "ominous",
				Pacing:        "medium",
				Metadata: map[string]any{
					"word_count":         1200,
					"emotional_score":    0.4,
					"thematic_relevance": 0.8,
				},
			},
		},
		{
			"The Final Battle",
			&types.SceneData{
				Title:         "The Final Battle",
				Summary:       "Epic confrontation with the Shadow Lord",
				Content:       "Lightning crackled as Elara raised the Crystal high above her head...",
				Act:           "Act3",
				Sequence:      25,
				EmotionalTone: "triumphant",
				Pacing:        "fast",
				Metadata: map[string]any{
					"word_count":         2100,
					"emotional_score":    0.95,
					"thematic_relevance": 1.0,
				},
			},
		},
	}

	for _, scene := range scenes {
		data, _ := types.MarshalEntityData(scene.data)
		entities = append(entities, db.CreateEntityParams{
			ID:         uuid.New().String(),
			VersionID:  versionID,
			EntityType: string(types.EntityTypeScene),
			Name:       scene.name,
			Data:       data,
		})
	}

	// Characters
	characters := []struct {
		name string
		data *types.CharacterData
	}{
		{
			"Elara the Brave",
			&types.CharacterData{
				Name:        "Elara the Brave",
				Role:        "protagonist",
				Description: "A young warrior destined to save the realm",
				PersonalityTraits: []string{"brave", "determined", "compassionate"},
				Background:        "Born in the northern kingdoms, trained by ancient masters",
				VoiceCharacteristics: types.VoiceCharacteristics{
					Tone:           "confident",
					Vocabulary:     "formal",
					SpeechPatterns: []string{"uses archaic terms", "speaks with authority"},
				},
				CharacterArc: types.CharacterArc{
					StartingState: "naive_hero",
					CurrentState:  "growing_wisdom",
					TargetState:   "wise_leader",
				},
			},
		},
		{
			"Shadow Lord Malachar",
			&types.CharacterData{
				Name:        "Shadow Lord Malachar",
				Role:        "antagonist",
				Description: "Ancient evil seeking to plunge the world into darkness",
				PersonalityTraits: []string{"cunning", "ruthless", "charismatic"},
				Background:        "Once a noble king, corrupted by dark magic centuries ago",
				VoiceCharacteristics: types.VoiceCharacteristics{
					Tone:           "menacing",
					Vocabulary:     "archaic",
					SpeechPatterns: []string{"speaks in riddles", "uses dark metaphors"},
				},
				CharacterArc: types.CharacterArc{
					StartingState: "supreme_confidence",
					CurrentState:  "growing_desperation",
					TargetState:   "ultimate_defeat",
				},
			},
		},
	}

	for _, char := range characters {
		data, _ := types.MarshalEntityData(char.data)
		entities = append(entities, db.CreateEntityParams{
			ID:         uuid.New().String(),
			VersionID:  versionID,
			EntityType: string(types.EntityTypeCharacter),
			Name:       char.name,
			Data:       data,
		})
	}

	// Locations
	locations := []struct {
		name string
		data *types.LocationData
	}{
		{
			"Shadowwood Forest",
			&types.LocationData{
				Name:        "Shadowwood Forest",
				Description: "Ancient woodland shrouded in perpetual twilight",
				Atmosphere:  "ominous",
				PhysicalDetails: types.PhysicalDetails{
					Size:     "vast",
					Lighting: "dim twilight",
					NotableFeatures: []string{"twisted ancient trees", "glowing mushrooms", "hidden paths"},
				},
				Significance: "Gateway to the Shadow Realm",
			},
		},
		{
			"Crystal Caverns",
			&types.LocationData{
				Name:        "Crystal Caverns",
				Description: "Mystical caves where the Crystal of Eternal Light rests",
				Atmosphere:  "mystical",
				PhysicalDetails: types.PhysicalDetails{
					Size:     "cathedral-like",
					Lighting: "ethereal crystal glow",
					NotableFeatures: []string{"floating crystals", "ancient runes", "pools of starlight"},
				},
				Significance: "Final destination of the quest",
			},
		},
	}

	for _, loc := range locations {
		data, _ := types.MarshalEntityData(loc.data)
		entities = append(entities, db.CreateEntityParams{
			ID:         uuid.New().String(),
			VersionID:  versionID,
			EntityType: string(types.EntityTypeLocation),
			Name:       loc.name,
			Data:       data,
		})
	}

	// Themes
	themes := []struct {
		name string
		data *types.ThemeData
	}{
		{
			"Good vs Evil",
			&types.ThemeData{
				Name:        "Good vs Evil",
				Description: "The eternal struggle between light and darkness",
				Questions:   []string{"What defines true good?", "Can evil ever be redeemed?"},
				Symbols:     []string{"light", "darkness", "crystal", "shadow"},
				Relevance:   0.95,
			},
		},
		{
			"Courage and Sacrifice",
			&types.ThemeData{
				Name:        "Courage and Sacrifice",
				Description: "The price of heroism and the courage to pay it",
				Questions:   []string{"What would you sacrifice for others?", "Where does true courage come from?"},
				Symbols:     []string{"sword", "shield", "flame", "mountain"},
				Relevance:   0.85,
			},
		},
	}

	for _, theme := range themes {
		data, _ := types.MarshalEntityData(theme.data)
		entities = append(entities, db.CreateEntityParams{
			ID:         uuid.New().String(),
			VersionID:  versionID,
			EntityType: string(types.EntityTypeTheme),
			Name:       theme.name,
			Data:       data,
		})
	}

	return entities
}

func createFantasyRelationships(versionID string, entityIDs map[string]string) []db.CreateRelationshipParams {
	var relationships []db.CreateRelationshipParams

	// Scene-Character relationships
	sceneCharRels := []struct {
		scene      string
		character  string
		relType    types.RelationshipType
		properties map[string]any
	}{
		{"The Call to Adventure", "Elara the Brave", types.RelationshipFeatures, map[string]any{"role": "protagonist", "importance": "primary"}},
		{"The Dark Forest", "Elara the Brave", types.RelationshipFeatures, map[string]any{"role": "protagonist", "importance": "primary"}},
		{"The Final Battle", "Elara the Brave", types.RelationshipFeatures, map[string]any{"role": "protagonist", "importance": "primary"}},
		{"The Final Battle", "Shadow Lord Malachar", types.RelationshipFeatures, map[string]any{"role": "antagonist", "importance": "primary"}},
	}

	for _, rel := range sceneCharRels {
		props, _ := json.Marshal(rel.properties)
		relationships = append(relationships, db.CreateRelationshipParams{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     entityIDs[rel.scene],
			ToEntityID:       entityIDs[rel.character],
			RelationshipType: string(rel.relType),
			Properties:       props,
		})
	}

	// Scene-Location relationships
	sceneLocRels := []struct {
		scene    string
		location string
		relType  types.RelationshipType
	}{
		{"The Dark Forest", "Shadowwood Forest", types.RelationshipOccursAt},
		{"The Final Battle", "Crystal Caverns", types.RelationshipOccursAt},
	}

	for _, rel := range sceneLocRels {
		relationships = append(relationships, db.CreateRelationshipParams{
			ID:               uuid.New().String(),
			VersionID:        versionID,
			FromEntityID:     entityIDs[rel.scene],
			ToEntityID:       entityIDs[rel.location],
			RelationshipType: string(rel.relType),
			Properties:       json.RawMessage(`{}`),
		})
	}

	// Character conflicts
	relationships = append(relationships, db.CreateRelationshipParams{
		ID:               uuid.New().String(),
		VersionID:        versionID,
		FromEntityID:     entityIDs["Elara the Brave"],
		ToEntityID:       entityIDs["Shadow Lord Malachar"],
		RelationshipType: string(types.RelationshipConflicts),
		Properties:       json.RawMessage(`{"intensity": "ultimate", "type": "good_vs_evil"}`),
	})

	return relationships
}

func createFantasyAnnotations(entityIDs map[string]string) []db.CreateAnnotationParams {
	var annotations []db.CreateAnnotationParams

	// Emotional analysis for scenes
	emotionalAnnotations := []struct {
		entityName string
		data       *types.EmotionalAnalysisData
		content    string
	}{
		{
			"The Call to Adventure",
			&types.EmotionalAnalysisData{
				Sentiment:    0.7,
				Emotions:     map[string]float64{"mystery": 0.8, "anticipation": 0.9, "hope": 0.7},
				EmotionalArc: "rising",
				ImpactScore:  0.8,
				Suggestions:  []string{"Consider adding more personal stakes", "Enhance the sense of urgency"},
				AnalyzedAt:   time.Now(),
			},
			"Strong opening with good emotional engagement and mystery setup",
		},
		{
			"The Final Battle",
			&types.EmotionalAnalysisData{
				Sentiment:    0.95,
				Emotions:     map[string]float64{"triumph": 0.95, "relief": 0.8, "satisfaction": 0.9},
				EmotionalArc: "climactic",
				ImpactScore:  0.98,
				Suggestions:  []string{"Perfect emotional climax", "Consider brief moment of doubt before victory"},
				AnalyzedAt:   time.Now(),
			},
			"Excellent climactic scene with maximum emotional impact and satisfying resolution",
		},
	}

	for _, ea := range emotionalAnnotations {
		data, _ := json.Marshal(ea.data)
		annotations = append(annotations, db.CreateAnnotationParams{
			ID:             uuid.New().String(),
			EntityID:       entityIDs[ea.entityName],
			AnnotationType: string(types.AnnotationEmotionalAnalysis),
			Content:        ea.content,
			Metadata:       data,
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		})
	}

	// Thematic analysis
	thematicAnnotations := []struct {
		entityName string
		data       *types.ThematicScoreData
		content    string
	}{
		{
			"The Final Battle",
			&types.ThematicScoreData{
				RelevanceScore: 0.98,
				ThemeAlignment: map[string]float64{
					entityIDs["Good vs Evil"]:         0.98,
					entityIDs["Courage and Sacrifice"]: 0.92,
				},
				Contributions: []string{"Ultimate expression of good vs evil theme", "Demonstrates ultimate sacrifice and courage"},
				Concerns:      []string{},
				AnalyzedAt:    time.Now(),
			},
			"Perfect thematic culmination bringing together all major themes of the story",
		},
	}

	for _, ta := range thematicAnnotations {
		data, _ := json.Marshal(ta.data)
		annotations = append(annotations, db.CreateAnnotationParams{
			ID:             uuid.New().String(),
			EntityID:       entityIDs[ta.entityName],
			AnnotationType: string(types.AnnotationThematicScore),
			Content:        ta.content,
			Metadata:       data,
			AgentName:      sql.NullString{String: "thematic_steward", Valid: true},
		})
	}

	return annotations
}

func seedSciFiStory(ctx context.Context, queries *db.Queries) {
	// Similar structure but with sci-fi content
	fmt.Println("Sci-fi seeding not yet implemented")
}

func seedMysteryStory(ctx context.Context, queries *db.Queries) {
	// Similar structure but with mystery content
	fmt.Println("Mystery seeding not yet implemented")
}