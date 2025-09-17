package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/monitoring"
	"github.com/barrynorthern/libretto/internal/types"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type TestSuite struct {
	queries   *db.Queries
	logger    *monitoring.Logger
	metrics   *monitoring.DatabaseMetrics
	database  *sql.DB
}

type TestResult struct {
	Name     string        `json:"name"`
	Passed   bool          `json:"passed"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
	Details  interface{}   `json:"details,omitempty"`
}

type TestReport struct {
	Timestamp    time.Time    `json:"timestamp"`
	TotalTests   int          `json:"total_tests"`
	PassedTests  int          `json:"passed_tests"`
	FailedTests  int          `json:"failed_tests"`
	TotalTime    time.Duration `json:"total_time"`
	Results      []TestResult `json:"results"`
}

func main() {
	var (
		dbPath     = flag.String("db", ":memory:", "Path to SQLite database (use :memory: for in-memory)")
		outputFile = flag.String("output", "", "Output file for test results (JSON)")
		verbose    = flag.Bool("v", false, "Verbose output")
	)
	flag.Parse()

	// Setup database
	database, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Apply migrations
	if err := applyMigrations(database); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Setup monitoring
	logger := monitoring.NewLogger("integration-test")
	metrics := monitoring.NewDatabaseMetrics(logger)

	suite := &TestSuite{
		queries:  db.New(database),
		logger:   logger,
		metrics:  metrics,
		database: database,
	}

	ctx := context.Background()
	
	// Run test suite
	report := suite.RunAllTests(ctx, *verbose)
	
	// Output results
	if *outputFile != "" {
		if err := suite.SaveReport(report, *outputFile); err != nil {
			log.Printf("Failed to save report: %v", err)
		}
	}
	
	suite.PrintSummary(report)
	
	if report.FailedTests > 0 {
		os.Exit(1)
	}
}

func (ts *TestSuite) RunAllTests(ctx context.Context, verbose bool) *TestReport {
	startTime := time.Now()
	
	tests := []struct {
		name string
		fn   func(context.Context) TestResult
	}{
		{"Database Schema Validation", ts.testDatabaseSchema},
		{"Project CRUD Operations", ts.testProjectCRUD},
		{"Graph Version Management", ts.testGraphVersions},
		{"Entity Management", ts.testEntityManagement},
		{"Relationship Management", ts.testRelationshipManagement},
		{"Annotation System", ts.testAnnotationSystem},
		{"Data Model Validation", ts.testDataModelValidation},
		{"Complex Narrative Graph", ts.testComplexNarrativeGraph},
		{"Performance Benchmarks", ts.testPerformance},
		{"Data Integrity", ts.testDataIntegrity},
		{"Concurrent Operations", ts.testConcurrentOperations},
	}
	
	var results []TestResult
	passed := 0
	
	for _, test := range tests {
		if verbose {
			fmt.Printf("Running test: %s...\n", test.name)
		}
		
		result := test.fn(ctx)
		results = append(results, result)
		
		if result.Passed {
			passed++
			if verbose {
				fmt.Printf("✓ %s (%.2fms)\n", test.name, float64(result.Duration.Nanoseconds())/1e6)
			}
		} else {
			if verbose {
				fmt.Printf("✗ %s (%.2fms): %s\n", test.name, float64(result.Duration.Nanoseconds())/1e6, result.Error)
			}
		}
	}
	
	return &TestReport{
		Timestamp:   startTime,
		TotalTests:  len(tests),
		PassedTests: passed,
		FailedTests: len(tests) - passed,
		TotalTime:   time.Since(startTime),
		Results:     results,
	}
}

func (ts *TestSuite) testDatabaseSchema(ctx context.Context) TestResult {
	start := time.Now()
	
	// Check that all expected tables exist
	expectedTables := []string{"projects", "graph_versions", "entities", "relationships", "annotations", "scenes"}
	
	for _, table := range expectedTables {
		var count int
		err := ts.database.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)).Scan(&count)
		if err != nil {
			return TestResult{
				Name:     "Database Schema Validation",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to check table %s: %v", table, err),
			}
		}
		if count != 1 {
			return TestResult{
				Name:     "Database Schema Validation",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Table %s not found", table),
			}
		}
	}
	
	return TestResult{
		Name:     "Database Schema Validation",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"tables_validated": len(expectedTables)},
	}
}

func (ts *TestSuite) testProjectCRUD(ctx context.Context) TestResult {
	start := time.Now()
	
	// Create project
	projectID := uuid.New().String()
	params := db.CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Integration test project", Valid: true},
	}
	
	project, err := ts.queries.CreateProject(ctx, params)
	if err != nil {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create project: %v", err),
		}
	}
	
	// Read project
	retrieved, err := ts.queries.GetProject(ctx, projectID)
	if err != nil {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to get project: %v", err),
		}
	}
	
	if retrieved.Name != project.Name {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    "Retrieved project name doesn't match created project",
		}
	}
	
	// Update project
	updateParams := db.UpdateProjectParams{
		ID:          projectID,
		Name:        "Updated Test Project",
		Theme:       sql.NullString{String: "Mystery", Valid: true},
		Genre:       sql.NullString{String: "Thriller", Valid: true},
		Description: sql.NullString{String: "Updated description", Valid: true},
	}
	
	updated, err := ts.queries.UpdateProject(ctx, updateParams)
	if err != nil {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to update project: %v", err),
		}
	}
	
	if updated.Name != "Updated Test Project" {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    "Project name not updated correctly",
		}
	}
	
	// Delete project
	err = ts.queries.DeleteProject(ctx, projectID)
	if err != nil {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to delete project: %v", err),
		}
	}
	
	// Verify deletion
	_, err = ts.queries.GetProject(ctx, projectID)
	if err == nil {
		return TestResult{
			Name:     "Project CRUD Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    "Project still exists after deletion",
		}
	}
	
	return TestResult{
		Name:     "Project CRUD Operations",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"operations": []string{"create", "read", "update", "delete"}},
	}
}

func (ts *TestSuite) testGraphVersions(ctx context.Context) TestResult {
	start := time.Now()
	
	// Create project first
	projectID := uuid.New().String()
	projectParams := db.CreateProjectParams{
		ID:          projectID,
		Name:        "Version Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Test project for versions", Valid: true},
	}
	
	_, err := ts.queries.CreateProject(ctx, projectParams)
	if err != nil {
		return TestResult{
			Name:     "Graph Version Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create project: %v", err),
		}
	}
	
	// Create initial version
	version1ID := uuid.New().String()
	version1Params := db.CreateGraphVersionParams{
		ID:            version1ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Initial Version", Valid: true},
		Description:   sql.NullString{String: "First version", Valid: true},
		IsWorkingSet:  true,
	}
	
	_, err = ts.queries.CreateGraphVersion(ctx, version1Params)
	if err != nil {
		return TestResult{
			Name:     "Graph Version Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create initial version: %v", err),
		}
	}
	
	// Create child version
	version2ID := uuid.New().String()
	version2Params := db.CreateGraphVersionParams{
		ID:            version2ID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{String: version1ID, Valid: true},
		Name:          sql.NullString{String: "Second Version", Valid: true},
		Description:   sql.NullString{String: "Child version", Valid: true},
		IsWorkingSet:  false,
	}
	
	_, err = ts.queries.CreateGraphVersion(ctx, version2Params)
	if err != nil {
		return TestResult{
			Name:     "Graph Version Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create child version: %v", err),
		}
	}
	
	// Test working set switching
	setWorkingSetParams := db.SetWorkingSetParams{
		ID:        version2ID,
		ProjectID: projectID,
	}
	
	err = ts.queries.SetWorkingSet(ctx, setWorkingSetParams)
	if err != nil {
		return TestResult{
			Name:     "Graph Version Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to set working set: %v", err),
		}
	}
	
	// Verify working set
	workingSet, err := ts.queries.GetWorkingSetVersion(ctx, projectID)
	if err != nil {
		return TestResult{
			Name:     "Graph Version Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to get working set: %v", err),
		}
	}
	
	if workingSet.ID != version2ID {
		return TestResult{
			Name:     "Graph Version Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    "Working set not switched correctly",
		}
	}
	
	return TestResult{
		Name:     "Graph Version Management",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"versions_created": 2, "working_set_switched": true},
	}
}

func (ts *TestSuite) testEntityManagement(ctx context.Context) TestResult {
	start := time.Now()
	
	// Setup project and version
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Entity Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Test different entity types
	entityTypes := []struct {
		entityType string
		data       interface{}
	}{
		{
			string(types.EntityTypeScene),
			&types.SceneData{
				Title:         "Test Scene",
				Summary:       "A test scene",
				Content:       "Scene content",
				Act:           "Act1",
				Sequence:      1,
				EmotionalTone: "neutral",
				Pacing:        "medium",
			},
		},
		{
			string(types.EntityTypeCharacter),
			&types.CharacterData{
				Name:        "Test Character",
				Role:        "protagonist",
				Description: "A test character",
				PersonalityTraits: []string{"brave", "kind"},
			},
		},
		{
			string(types.EntityTypeLocation),
			&types.LocationData{
				Name:        "Test Location",
				Description: "A test location",
				Atmosphere:  "peaceful",
			},
		},
	}
	
	createdEntities := make([]string, 0, len(entityTypes))
	
	for _, et := range entityTypes {
		entityID := uuid.New().String()
		data, err := types.MarshalEntityData(et.data)
		if err != nil {
			return TestResult{
				Name:     "Entity Management",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to marshal %s data: %v", et.entityType, err),
			}
		}
		
		params := db.CreateEntityParams{
			ID:         entityID,
			VersionID:  versionID,
			EntityType: et.entityType,
			Name:       fmt.Sprintf("Test %s", et.entityType),
			Data:       data,
		}
		
		_, err = ts.queries.CreateEntity(ctx, params)
		if err != nil {
			return TestResult{
				Name:     "Entity Management",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create %s entity: %v", et.entityType, err),
			}
		}
		
		createdEntities = append(createdEntities, entityID)
	}
	
	// Test entity listing
	entities, err := ts.queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		return TestResult{
			Name:     "Entity Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list entities: %v", err),
		}
	}
	
	if len(entities) != len(entityTypes) {
		return TestResult{
			Name:     "Entity Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected %d entities, got %d", len(entityTypes), len(entities)),
		}
	}
	
	// Test entity counting by type
	for _, et := range entityTypes {
		countParams := db.CountEntitiesByTypeParams{
			VersionID:  versionID,
			EntityType: et.entityType,
		}
		count, err := ts.queries.CountEntitiesByType(ctx, countParams)
		if err != nil {
			return TestResult{
				Name:     "Entity Management",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to count %s entities: %v", et.entityType, err),
			}
		}
		if count != 1 {
			return TestResult{
				Name:     "Entity Management",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Expected 1 %s entity, got %d", et.entityType, count),
			}
		}
	}
	
	return TestResult{
		Name:     "Entity Management",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"entities_created": len(createdEntities), "types_tested": len(entityTypes)},
	}
}

func (ts *TestSuite) testRelationshipManagement(ctx context.Context) TestResult {
	start := time.Now()
	
	// Setup project, version, and entities
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Relationship Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Create test entities
	entity1ID := uuid.New().String()
	entity2ID := uuid.New().String()
	
	sceneData, _ := types.MarshalEntityData(&types.SceneData{
		Title:   "Test Scene",
		Summary: "A test scene",
		Content: "Scene content",
	})
	
	characterData, _ := types.MarshalEntityData(&types.CharacterData{
		Name: "Test Character",
		Role: "protagonist",
	})
	
	entities := []db.CreateEntityParams{
		{
			ID:         entity1ID,
			VersionID:  versionID,
			EntityType: string(types.EntityTypeScene),
			Name:       "Test Scene",
			Data:       sceneData,
		},
		{
			ID:         entity2ID,
			VersionID:  versionID,
			EntityType: string(types.EntityTypeCharacter),
			Name:       "Test Character",
			Data:       characterData,
		},
	}
	
	for _, entity := range entities {
		_, err := ts.queries.CreateEntity(ctx, entity)
		if err != nil {
			return TestResult{
				Name:     "Relationship Management",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create entity: %v", err),
			}
		}
	}
	
	// Create relationship
	relationshipID := uuid.New().String()
	properties := map[string]interface{}{"role": "protagonist", "importance": "high"}
	propertiesJSON, _ := json.Marshal(properties)
	
	relParams := db.CreateRelationshipParams{
		ID:               relationshipID,
		VersionID:        versionID,
		FromEntityID:     entity1ID,
		ToEntityID:       entity2ID,
		RelationshipType: string(types.RelationshipFeatures),
		Properties:       propertiesJSON,
	}
	
	_, err = ts.queries.CreateRelationship(ctx, relParams)
	if err != nil {
		return TestResult{
			Name:     "Relationship Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create relationship: %v", err),
		}
	}
	
	// Test relationship queries
	relationships, err := ts.queries.ListRelationshipsByVersion(ctx, versionID)
	if err != nil {
		return TestResult{
			Name:     "Relationship Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list relationships: %v", err),
		}
	}
	
	if len(relationships) != 1 {
		return TestResult{
			Name:     "Relationship Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected 1 relationship, got %d", len(relationships)),
		}
	}
	
	// Test relationship by entity
	entityRelParams := db.ListRelationshipsByEntityParams{
		FromEntityID: entity1ID,
		ToEntityID:   entity1ID,
	}
	entityRels, err := ts.queries.ListRelationshipsByEntity(ctx, entityRelParams)
	if err != nil {
		return TestResult{
			Name:     "Relationship Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list relationships by entity: %v", err),
		}
	}
	
	if len(entityRels) != 1 {
		return TestResult{
			Name:     "Relationship Management",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected 1 entity relationship, got %d", len(entityRels)),
		}
	}
	
	return TestResult{
		Name:     "Relationship Management",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"relationships_created": 1, "queries_tested": 2},
	}
}

func (ts *TestSuite) testAnnotationSystem(ctx context.Context) TestResult {
	start := time.Now()
	
	// Setup project, version, and entity
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Annotation System",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Create test entity
	entityID := uuid.New().String()
	sceneData, _ := types.MarshalEntityData(&types.SceneData{
		Title:   "Test Scene",
		Summary: "A test scene",
		Content: "Scene content",
	})
	
	entityParams := db.CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: string(types.EntityTypeScene),
		Name:       "Test Scene",
		Data:       sceneData,
	}
	
	_, err = ts.queries.CreateEntity(ctx, entityParams)
	if err != nil {
		return TestResult{
			Name:     "Annotation System",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create entity: %v", err),
		}
	}
	
	// Create different types of annotations
	annotationTypes := []struct {
		annotationType string
		data           interface{}
		content        string
		agent          string
	}{
		{
			string(types.AnnotationEmotionalAnalysis),
			&types.EmotionalAnalysisData{
				Sentiment:    0.7,
				Emotions:     map[string]float64{"joy": 0.8, "excitement": 0.6},
				EmotionalArc: "rising",
				ImpactScore:  0.75,
				AnalyzedAt:   time.Now(),
			},
			"Positive emotional tone with rising arc",
			"empath_agent",
		},
		{
			string(types.AnnotationThematicScore),
			&types.ThematicScoreData{
				RelevanceScore: 0.85,
				ThemeAlignment: map[string]float64{"courage": 0.9, "friendship": 0.8},
				AnalyzedAt:     time.Now(),
			},
			"Strong thematic alignment with core themes",
			"thematic_steward",
		},
	}
	
	createdAnnotations := make([]string, 0, len(annotationTypes))
	
	for _, at := range annotationTypes {
		annotationID := uuid.New().String()
		metadata, err := json.Marshal(at.data)
		if err != nil {
			return TestResult{
				Name:     "Annotation System",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to marshal %s data: %v", at.annotationType, err),
			}
		}
		
		params := db.CreateAnnotationParams{
			ID:             annotationID,
			EntityID:       entityID,
			AnnotationType: at.annotationType,
			Content:        at.content,
			Metadata:       metadata,
			AgentName:      sql.NullString{String: at.agent, Valid: true},
		}
		
		_, err = ts.queries.CreateAnnotation(ctx, params)
		if err != nil {
			return TestResult{
				Name:     "Annotation System",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create %s annotation: %v", at.annotationType, err),
			}
		}
		
		createdAnnotations = append(createdAnnotations, annotationID)
	}
	
	// Test annotation queries
	annotations, err := ts.queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		return TestResult{
			Name:     "Annotation System",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list annotations: %v", err),
		}
	}
	
	if len(annotations) != len(annotationTypes) {
		return TestResult{
			Name:     "Annotation System",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected %d annotations, got %d", len(annotationTypes), len(annotations)),
		}
	}
	
	// Test annotation by type
	typeParams := db.ListAnnotationsByTypeParams{
		EntityID:       entityID,
		AnnotationType: string(types.AnnotationEmotionalAnalysis),
	}
	typeAnnotations, err := ts.queries.ListAnnotationsByType(ctx, typeParams)
	if err != nil {
		return TestResult{
			Name:     "Annotation System",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list annotations by type: %v", err),
		}
	}
	
	if len(typeAnnotations) != 1 {
		return TestResult{
			Name:     "Annotation System",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected 1 emotional analysis annotation, got %d", len(typeAnnotations)),
		}
	}
	
	return TestResult{
		Name:     "Annotation System",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"annotations_created": len(createdAnnotations), "types_tested": len(annotationTypes)},
	}
}

func (ts *TestSuite) testDataModelValidation(ctx context.Context) TestResult {
	start := time.Now()
	
	// Test all entity data types for marshal/unmarshal
	testCases := []struct {
		name string
		data interface{}
		test func(json.RawMessage) error
	}{
		{
			"SceneData",
			&types.SceneData{
				Title:         "Test Scene",
				Summary:       "Test summary",
				Content:       "Test content",
				Act:           "Act1",
				Sequence:      1,
				EmotionalTone: "neutral",
				Pacing:        "medium",
				Characters:    []string{"char1", "char2"},
				Location:      "loc1",
				Themes:        []string{"theme1"},
				Metadata:      map[string]any{"test": "value"},
			},
			func(raw json.RawMessage) error {
				_, err := types.UnmarshalSceneData(raw)
				return err
			},
		},
		{
			"CharacterData",
			&types.CharacterData{
				Name:        "Test Character",
				Role:        "protagonist",
				Description: "Test description",
				PersonalityTraits: []string{"brave", "kind"},
				Background:        "Test background",
				VoiceCharacteristics: types.VoiceCharacteristics{
					Tone:           "confident",
					Vocabulary:     "formal",
					SpeechPatterns: []string{"uses metaphors"},
				},
				CharacterArc: types.CharacterArc{
					StartingState: "naive",
					CurrentState:  "learning",
					TargetState:   "wise",
				},
			},
			func(raw json.RawMessage) error {
				_, err := types.UnmarshalCharacterData(raw)
				return err
			},
		},
		{
			"EmotionalAnalysisData",
			&types.EmotionalAnalysisData{
				Sentiment:    0.7,
				Emotions:     map[string]float64{"joy": 0.8, "excitement": 0.6},
				EmotionalArc: "rising",
				ImpactScore:  0.75,
				Suggestions:  []string{"Add more tension"},
				AnalyzedAt:   time.Now(),
			},
			func(raw json.RawMessage) error {
				_, err := types.UnmarshalEmotionalAnalysisData(raw)
				return err
			},
		},
	}
	
	for _, tc := range testCases {
		// Marshal
		marshaled, err := types.MarshalEntityData(tc.data)
		if err != nil {
			return TestResult{
				Name:     "Data Model Validation",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to marshal %s: %v", tc.name, err),
			}
		}
		
		// Unmarshal
		err = tc.test(marshaled)
		if err != nil {
			return TestResult{
				Name:     "Data Model Validation",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to unmarshal %s: %v", tc.name, err),
			}
		}
	}
	
	return TestResult{
		Name:     "Data Model Validation",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"data_types_tested": len(testCases)},
	}
}

func (ts *TestSuite) testComplexNarrativeGraph(ctx context.Context) TestResult {
	start := time.Now()
	
	// Create a complex narrative graph similar to the fantasy story
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Complex Narrative Graph",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Create multiple entities of different types
	entityCount := 0
	relationshipCount := 0
	annotationCount := 0
	
	// Create scenes
	scenes := []string{"Opening", "Conflict", "Resolution"}
	sceneIDs := make([]string, len(scenes))
	
	for i, sceneName := range scenes {
		sceneID := uuid.New().String()
		sceneIDs[i] = sceneID
		
		sceneData, _ := types.MarshalEntityData(&types.SceneData{
			Title:         sceneName,
			Summary:       fmt.Sprintf("%s scene", sceneName),
			Content:       fmt.Sprintf("Content for %s", sceneName),
			Act:           fmt.Sprintf("Act%d", i+1),
			Sequence:      i + 1,
			EmotionalTone: "dramatic",
			Pacing:        "medium",
		})
		
		params := db.CreateEntityParams{
			ID:         sceneID,
			VersionID:  versionID,
			EntityType: string(types.EntityTypeScene),
			Name:       sceneName,
			Data:       sceneData,
		}
		
		_, err := ts.queries.CreateEntity(ctx, params)
		if err != nil {
			return TestResult{
				Name:     "Complex Narrative Graph",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create scene %s: %v", sceneName, err),
			}
		}
		entityCount++
	}
	
	// Create characters
	characters := []string{"Hero", "Villain", "Mentor"}
	characterIDs := make([]string, len(characters))
	
	for i, charName := range characters {
		charID := uuid.New().String()
		characterIDs[i] = charID
		
		charData, _ := types.MarshalEntityData(&types.CharacterData{
			Name:        charName,
			Role:        []string{"protagonist", "antagonist", "mentor"}[i],
			Description: fmt.Sprintf("%s character", charName),
			PersonalityTraits: []string{"brave", "cunning", "wise"}[i:i+1],
		})
		
		params := db.CreateEntityParams{
			ID:         charID,
			VersionID:  versionID,
			EntityType: string(types.EntityTypeCharacter),
			Name:       charName,
			Data:       charData,
		}
		
		_, err := ts.queries.CreateEntity(ctx, params)
		if err != nil {
			return TestResult{
				Name:     "Complex Narrative Graph",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create character %s: %v", charName, err),
			}
		}
		entityCount++
	}
	
	// Create relationships between scenes and characters
	for i, sceneID := range sceneIDs {
		for j, charID := range characterIDs {
			relID := uuid.New().String()
			properties := map[string]interface{}{
				"importance": []string{"primary", "secondary", "tertiary"}[j],
				"scene_role": fmt.Sprintf("role_in_scene_%d", i+1),
			}
			propertiesJSON, _ := json.Marshal(properties)
			
			relParams := db.CreateRelationshipParams{
				ID:               relID,
				VersionID:        versionID,
				FromEntityID:     sceneID,
				ToEntityID:       charID,
				RelationshipType: string(types.RelationshipFeatures),
				Properties:       propertiesJSON,
			}
			
			_, err := ts.queries.CreateRelationship(ctx, relParams)
			if err != nil {
				return TestResult{
					Name:     "Complex Narrative Graph",
					Passed:   false,
					Duration: time.Since(start),
					Error:    fmt.Sprintf("Failed to create relationship: %v", err),
				}
			}
			relationshipCount++
		}
	}
	
	// Create annotations for scenes
	for _, sceneID := range sceneIDs {
		// Emotional analysis
		emotionalData, _ := json.Marshal(&types.EmotionalAnalysisData{
			Sentiment:    0.6,
			Emotions:     map[string]float64{"tension": 0.7, "excitement": 0.8},
			EmotionalArc: "rising",
			ImpactScore:  0.75,
			AnalyzedAt:   time.Now(),
		})
		
		annotationID := uuid.New().String()
		annotationParams := db.CreateAnnotationParams{
			ID:             annotationID,
			EntityID:       sceneID,
			AnnotationType: string(types.AnnotationEmotionalAnalysis),
			Content:        "Emotional analysis of scene",
			Metadata:       emotionalData,
			AgentName:      sql.NullString{String: "empath_agent", Valid: true},
		}
		
		_, err := ts.queries.CreateAnnotation(ctx, annotationParams)
		if err != nil {
			return TestResult{
				Name:     "Complex Narrative Graph",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create annotation: %v", err),
			}
		}
		annotationCount++
	}
	
	// Verify the complete graph
	entities, err := ts.queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		return TestResult{
			Name:     "Complex Narrative Graph",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list entities: %v", err),
		}
	}
	
	relationships, err := ts.queries.ListRelationshipsByVersion(ctx, versionID)
	if err != nil {
		return TestResult{
			Name:     "Complex Narrative Graph",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list relationships: %v", err),
		}
	}
	
	if len(entities) != entityCount {
		return TestResult{
			Name:     "Complex Narrative Graph",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected %d entities, got %d", entityCount, len(entities)),
		}
	}
	
	if len(relationships) != relationshipCount {
		return TestResult{
			Name:     "Complex Narrative Graph",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected %d relationships, got %d", relationshipCount, len(relationships)),
		}
	}
	
	return TestResult{
		Name:     "Complex Narrative Graph",
		Passed:   true,
		Duration: time.Since(start),
		Details: map[string]interface{}{
			"entities":      entityCount,
			"relationships": relationshipCount,
			"annotations":   annotationCount,
		},
	}
}

func (ts *TestSuite) testPerformance(ctx context.Context) TestResult {
	start := time.Now()
	
	// Performance benchmarks
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Performance Benchmarks",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Benchmark entity creation
	entityCreationStart := time.Now()
	entityCount := 100
	
	for i := 0; i < entityCount; i++ {
		entityID := uuid.New().String()
		sceneData, _ := types.MarshalEntityData(&types.SceneData{
			Title:   fmt.Sprintf("Scene %d", i),
			Summary: fmt.Sprintf("Summary %d", i),
			Content: fmt.Sprintf("Content %d", i),
		})
		
		params := db.CreateEntityParams{
			ID:         entityID,
			VersionID:  versionID,
			EntityType: string(types.EntityTypeScene),
			Name:       fmt.Sprintf("Scene %d", i),
			Data:       sceneData,
		}
		
		_, err := ts.queries.CreateEntity(ctx, params)
		if err != nil {
			return TestResult{
				Name:     "Performance Benchmarks",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Failed to create entity %d: %v", i, err),
			}
		}
	}
	
	entityCreationDuration := time.Since(entityCreationStart)
	
	// Benchmark entity listing
	listingStart := time.Now()
	entities, err := ts.queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		return TestResult{
			Name:     "Performance Benchmarks",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list entities: %v", err),
		}
	}
	listingDuration := time.Since(listingStart)
	
	if len(entities) != entityCount {
		return TestResult{
			Name:     "Performance Benchmarks",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected %d entities, got %d", entityCount, len(entities)),
		}
	}
	
	// Performance thresholds (adjust based on requirements)
	maxCreationTime := time.Duration(entityCount) * 10 * time.Millisecond // 10ms per entity
	maxListingTime := 100 * time.Millisecond                              // 100ms for listing
	
	if entityCreationDuration > maxCreationTime {
		return TestResult{
			Name:     "Performance Benchmarks",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Entity creation too slow: %v > %v", entityCreationDuration, maxCreationTime),
		}
	}
	
	if listingDuration > maxListingTime {
		return TestResult{
			Name:     "Performance Benchmarks",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Entity listing too slow: %v > %v", listingDuration, maxListingTime),
		}
	}
	
	return TestResult{
		Name:     "Performance Benchmarks",
		Passed:   true,
		Duration: time.Since(start),
		Details: map[string]interface{}{
			"entities_created":        entityCount,
			"creation_duration_ms":    float64(entityCreationDuration.Nanoseconds()) / 1e6,
			"listing_duration_ms":     float64(listingDuration.Nanoseconds()) / 1e6,
			"creation_per_entity_ms":  float64(entityCreationDuration.Nanoseconds()) / float64(entityCount) / 1e6,
		},
	}
}

func (ts *TestSuite) testDataIntegrity(ctx context.Context) TestResult {
	start := time.Now()
	
	// Test foreign key constraints and cascade deletes
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Create entity
	entityID := uuid.New().String()
	sceneData, _ := types.MarshalEntityData(&types.SceneData{
		Title:   "Test Scene",
		Summary: "Test summary",
		Content: "Test content",
	})
	
	entityParams := db.CreateEntityParams{
		ID:         entityID,
		VersionID:  versionID,
		EntityType: string(types.EntityTypeScene),
		Name:       "Test Scene",
		Data:       sceneData,
	}
	
	_, err = ts.queries.CreateEntity(ctx, entityParams)
	if err != nil {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create entity: %v", err),
		}
	}
	
	// Create annotation for entity
	annotationID := uuid.New().String()
	annotationParams := db.CreateAnnotationParams{
		ID:             annotationID,
		EntityID:       entityID,
		AnnotationType: string(types.AnnotationEmotionalAnalysis),
		Content:        "Test annotation",
		Metadata:       json.RawMessage(`{"test": "data"}`),
		AgentName:      sql.NullString{String: "test_agent", Valid: true},
	}
	
	_, err = ts.queries.CreateAnnotation(ctx, annotationParams)
	if err != nil {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to create annotation: %v", err),
		}
	}
	
	// Verify annotation exists
	annotations, err := ts.queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list annotations: %v", err),
		}
	}
	
	if len(annotations) != 1 {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected 1 annotation, got %d", len(annotations)),
		}
	}
	
	// Delete entity (should cascade delete annotation)
	err = ts.queries.DeleteEntity(ctx, entityID)
	if err != nil {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to delete entity: %v", err),
		}
	}
	
	// Verify annotation was cascade deleted
	annotationsAfterDelete, err := ts.queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list annotations after delete: %v", err),
		}
	}
	
	if len(annotationsAfterDelete) != 0 {
		return TestResult{
			Name:     "Data Integrity",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected 0 annotations after cascade delete, got %d", len(annotationsAfterDelete)),
		}
	}
	
	return TestResult{
		Name:     "Data Integrity",
		Passed:   true,
		Duration: time.Since(start),
		Details:  map[string]interface{}{"cascade_delete_verified": true},
	}
}

func (ts *TestSuite) testConcurrentOperations(ctx context.Context) TestResult {
	start := time.Now()
	
	// Test sequential entity creation (SQLite in-memory doesn't support true concurrency)
	_, versionID, err := ts.setupProjectAndVersion(ctx)
	if err != nil {
		return TestResult{
			Name:     "Concurrent Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to setup project: %v", err),
		}
	}
	
	// Create entities sequentially to simulate concurrent load
	totalEntities := 50
	
	for i := 0; i < totalEntities; i++ {
		entityID := uuid.New().String()
		sceneData, _ := types.MarshalEntityData(&types.SceneData{
			Title:   fmt.Sprintf("Load Test Scene %d", i),
			Summary: fmt.Sprintf("Summary %d", i),
			Content: fmt.Sprintf("Content %d", i),
		})
		
		params := db.CreateEntityParams{
			ID:         entityID,
			VersionID:  versionID,
			EntityType: string(types.EntityTypeScene),
			Name:       fmt.Sprintf("Load Test Scene %d", i),
			Data:       sceneData,
		}
		
		_, err := ts.queries.CreateEntity(ctx, params)
		if err != nil {
			return TestResult{
				Name:     "Concurrent Operations",
				Passed:   false,
				Duration: time.Since(start),
				Error:    fmt.Sprintf("Entity creation failed at %d: %v", i, err),
			}
		}
	}
	
	// Verify all entities exist
	entities, err := ts.queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		return TestResult{
			Name:     "Concurrent Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Failed to list entities: %v", err),
		}
	}
	
	if len(entities) != totalEntities {
		return TestResult{
			Name:     "Concurrent Operations",
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Expected %d entities in database, got %d", totalEntities, len(entities)),
		}
	}
	
	return TestResult{
		Name:     "Concurrent Operations",
		Passed:   true,
		Duration: time.Since(start),
		Details: map[string]interface{}{
			"total_entities_created": totalEntities,
			"test_type":             "sequential_load_test",
		},
	}
}

// Helper methods

func (ts *TestSuite) setupProjectAndVersion(ctx context.Context) (string, string, error) {
	// Create project
	projectID := uuid.New().String()
	projectParams := db.CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Integration test project", Valid: true},
	}
	
	_, err := ts.queries.CreateProject(ctx, projectParams)
	if err != nil {
		return "", "", err
	}
	
	// Create version
	versionID := uuid.New().String()
	versionParams := db.CreateGraphVersionParams{
		ID:            versionID,
		ProjectID:     projectID,
		ParentVersionID: sql.NullString{},
		Name:          sql.NullString{String: "Test Version", Valid: true},
		Description:   sql.NullString{String: "Test version", Valid: true},
		IsWorkingSet:  true,
	}
	
	_, err = ts.queries.CreateGraphVersion(ctx, versionParams)
	if err != nil {
		return "", "", err
	}
	
	return projectID, versionID, nil
}

func (ts *TestSuite) SaveReport(report *TestReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}

func (ts *TestSuite) PrintSummary(report *TestReport) {
	fmt.Printf("\n=== TEST SUMMARY ===\n")
	fmt.Printf("Total Tests: %d\n", report.TotalTests)
	fmt.Printf("Passed: %d\n", report.PassedTests)
	fmt.Printf("Failed: %d\n", report.FailedTests)
	fmt.Printf("Total Time: %.2fms\n", float64(report.TotalTime.Nanoseconds())/1e6)
	fmt.Printf("Success Rate: %.1f%%\n", float64(report.PassedTests)/float64(report.TotalTests)*100)
	
	if report.FailedTests > 0 {
		fmt.Printf("\n=== FAILED TESTS ===\n")
		for _, result := range report.Results {
			if !result.Passed {
				fmt.Printf("✗ %s: %s\n", result.Name, result.Error)
			}
		}
	}
	
	fmt.Printf("\n=== DETAILED RESULTS ===\n")
	for _, result := range report.Results {
		status := "✓"
		if !result.Passed {
			status = "✗"
		}
		fmt.Printf("%s %s (%.2fms)\n", status, result.Name, float64(result.Duration.Nanoseconds())/1e6)
		if result.Details != nil {
			detailsJSON, _ := json.MarshalIndent(result.Details, "  ", "  ")
			fmt.Printf("  Details: %s\n", string(detailsJSON))
		}
	}
}

func applyMigrations(database *sql.DB) error {
	// Enable foreign key constraints
	if _, err := database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %v", err)
	}
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
		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_graph_versions_project_id ON graph_versions(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_graph_versions_working_set ON graph_versions(project_id, is_working_set) WHERE is_working_set = TRUE`,
		`CREATE INDEX IF NOT EXISTS idx_entities_version_id ON entities(version_id)`,
		`CREATE INDEX IF NOT EXISTS idx_entities_type ON entities(version_id, entity_type)`,
		`CREATE INDEX IF NOT EXISTS idx_relationships_version_id ON relationships(version_id)`,
		`CREATE INDEX IF NOT EXISTS idx_relationships_from_entity ON relationships(from_entity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_relationships_to_entity ON relationships(to_entity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_relationships_type ON relationships(version_id, relationship_type)`,
		`CREATE INDEX IF NOT EXISTS idx_annotations_entity_id ON annotations(entity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_annotations_type ON annotations(entity_id, annotation_type)`,
		// Triggers
		`CREATE TRIGGER IF NOT EXISTS update_projects_updated_at 
			AFTER UPDATE ON projects
			FOR EACH ROW
		BEGIN
			UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END`,
		`CREATE TRIGGER IF NOT EXISTS update_entities_updated_at 
			AFTER UPDATE ON entities
			FOR EACH ROW
		BEGIN
			UPDATE entities SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END`,
		// Unique constraint for working set
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_working_set_per_project 
		ON graph_versions(project_id) 
		WHERE is_working_set = TRUE`,
	}

	for _, migration := range migrations {
		if _, err := database.Exec(migration); err != nil {
			return fmt.Errorf("failed to apply migration: %v", err)
		}
	}

	return nil
}