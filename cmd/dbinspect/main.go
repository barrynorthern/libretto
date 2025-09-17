package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dbPath    = flag.String("db", "libretto.db", "Path to SQLite database")
		command   = flag.String("cmd", "schema", "Command: schema, projects, entities, relationships, annotations, graph")
		projectID = flag.String("project", "", "Project ID for filtering")
		versionID = flag.String("version", "", "Version ID for filtering")
		entityID  = flag.String("entity", "", "Entity ID for filtering")
		verbose   = flag.Bool("v", false, "Verbose output")
	)
	flag.Parse()

	database, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	queries := db.New(database)
	ctx := context.Background()

	switch *command {
	case "schema":
		showSchema(database)
	case "projects":
		showProjects(ctx, queries, *verbose)
	case "entities":
		showEntities(ctx, queries, *projectID, *versionID, *verbose)
	case "relationships":
		showRelationships(ctx, queries, *versionID, *entityID, *verbose)
	case "annotations":
		showAnnotations(ctx, queries, *entityID, *verbose)
	case "graph":
		showGraph(ctx, queries, *projectID, *versionID)
	case "stats":
		showStats(ctx, queries, *projectID, *versionID)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: schema, projects, entities, relationships, annotations, graph, stats")
	}
}

func showSchema(db *sql.DB) {
	fmt.Println("=== DATABASE SCHEMA ===")
	
	// Get all tables
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		log.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	for _, table := range tables {
		fmt.Printf("\n--- Table: %s ---\n", table)
		
		// Get table schema
		schemaRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			log.Printf("Failed to get schema for %s: %v", table, err)
			continue
		}
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Column\tType\tNot Null\tDefault\tPK")
		
		for schemaRows.Next() {
			var cid int
			var name, dataType string
			var notNull, pk int
			var defaultValue sql.NullString
			
			schemaRows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
			
			defaultStr := "NULL"
			if defaultValue.Valid {
				defaultStr = defaultValue.String
			}
			
			fmt.Fprintf(w, "%s\t%s\t%t\t%s\t%t\n", 
				name, dataType, notNull == 1, defaultStr, pk == 1)
		}
		w.Flush()
		schemaRows.Close()
		
		// Get row count
		var count int
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		fmt.Printf("Rows: %d\n", count)
	}
}

func showProjects(ctx context.Context, queries *db.Queries, verbose bool) {
	fmt.Println("=== PROJECTS ===")
	
	projects, err := queries.ListProjects(ctx)
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if verbose {
		fmt.Fprintln(w, "ID\tName\tTheme\tGenre\tDescription\tCreated")
		for _, p := range projects {
			theme := "N/A"
			if p.Theme.Valid {
				theme = p.Theme.String
			}
			genre := "N/A"
			if p.Genre.Valid {
				genre = p.Genre.String
			}
			desc := "N/A"
			if p.Description.Valid {
				desc = truncate(p.Description.String, 30)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", 
				p.ID, p.Name, theme, genre, desc, p.CreatedAt.Format("2006-01-02 15:04"))
		}
	} else {
		fmt.Fprintln(w, "ID\tName\tTheme\tGenre")
		for _, p := range projects {
			theme := "N/A"
			if p.Theme.Valid {
				theme = p.Theme.String
			}
			genre := "N/A"
			if p.Genre.Valid {
				genre = p.Genre.String
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.ID, p.Name, theme, genre)
		}
	}
	w.Flush()
}

func showEntities(ctx context.Context, queries *db.Queries, projectID, versionID string, verbose bool) {
	fmt.Println("=== ENTITIES ===")
	
	if versionID == "" && projectID != "" {
		// Get working set version for project
		workingSet, err := queries.GetWorkingSetVersion(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to get working set for project %s: %v", projectID, err)
		}
		versionID = workingSet.ID
		fmt.Printf("Using working set version: %s\n", versionID)
	}
	
	if versionID == "" {
		fmt.Println("Please specify either -project or -version")
		return
	}

	entities, err := queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		log.Fatalf("Failed to list entities: %v", err)
	}

	if len(entities) == 0 {
		fmt.Println("No entities found.")
		return
	}

	// Group by type
	entityTypes := make(map[string][]db.Entity)
	for _, entity := range entities {
		entityTypes[entity.EntityType] = append(entityTypes[entity.EntityType], entity)
	}

	for entityType, typeEntities := range entityTypes {
		fmt.Printf("\n--- %s Entities (%d) ---\n", entityType, len(typeEntities))
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if verbose {
			fmt.Fprintln(w, "ID\tName\tData Preview\tCreated")
			for _, e := range typeEntities {
				dataPreview := getDataPreview(e.Data, entityType)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
					e.ID, e.Name, dataPreview, e.CreatedAt.Format("2006-01-02 15:04"))
			}
		} else {
			fmt.Fprintln(w, "ID\tName")
			for _, e := range typeEntities {
				fmt.Fprintf(w, "%s\t%s\n", e.ID, e.Name)
			}
		}
		w.Flush()
	}
}

func showRelationships(ctx context.Context, queries *db.Queries, versionID, entityID string, verbose bool) {
	fmt.Println("=== RELATIONSHIPS ===")
	
	var relationships []db.Relationship
	var err error
	
	if entityID != "" {
		params := db.ListRelationshipsByEntityParams{
			FromEntityID: entityID,
			ToEntityID:   entityID,
		}
		relationships, err = queries.ListRelationshipsByEntity(ctx, params)
		if err != nil {
			log.Fatalf("Failed to list relationships for entity %s: %v", entityID, err)
		}
		fmt.Printf("Relationships for entity: %s\n", entityID)
	} else if versionID != "" {
		relationships, err = queries.ListRelationshipsByVersion(ctx, versionID)
		if err != nil {
			log.Fatalf("Failed to list relationships for version %s: %v", versionID, err)
		}
		fmt.Printf("Relationships for version: %s\n", versionID)
	} else {
		fmt.Println("Please specify either -version or -entity")
		return
	}

	if len(relationships) == 0 {
		fmt.Println("No relationships found.")
		return
	}

	// Group by type
	relTypes := make(map[string][]db.Relationship)
	for _, rel := range relationships {
		relTypes[rel.RelationshipType] = append(relTypes[rel.RelationshipType], rel)
	}

	for relType, typeRels := range relTypes {
		fmt.Printf("\n--- %s Relationships (%d) ---\n", relType, len(typeRels))
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if verbose {
			fmt.Fprintln(w, "From Entity\tTo Entity\tProperties\tCreated")
			for _, r := range typeRels {
				props := "N/A"
				if len(r.Properties) > 0 {
					props = truncate(string(r.Properties), 30)
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
					r.FromEntityID, r.ToEntityID, props, r.CreatedAt.Format("2006-01-02 15:04"))
			}
		} else {
			fmt.Fprintln(w, "From Entity\tTo Entity")
			for _, r := range typeRels {
				fmt.Fprintf(w, "%s\t%s\n", r.FromEntityID, r.ToEntityID)
			}
		}
		w.Flush()
	}
}

func showAnnotations(ctx context.Context, queries *db.Queries, entityID string, verbose bool) {
	fmt.Println("=== ANNOTATIONS ===")
	
	if entityID == "" {
		fmt.Println("Please specify -entity")
		return
	}

	annotations, err := queries.ListAnnotationsByEntity(ctx, entityID)
	if err != nil {
		log.Fatalf("Failed to list annotations for entity %s: %v", entityID, err)
	}

	if len(annotations) == 0 {
		fmt.Println("No annotations found.")
		return
	}

	// Group by type
	annotationTypes := make(map[string][]db.Annotation)
	for _, annotation := range annotations {
		annotationTypes[annotation.AnnotationType] = append(annotationTypes[annotation.AnnotationType], annotation)
	}

	for annotationType, typeAnnotations := range annotationTypes {
		fmt.Printf("\n--- %s Annotations (%d) ---\n", annotationType, len(typeAnnotations))
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if verbose {
			fmt.Fprintln(w, "Agent\tContent\tMetadata Preview\tCreated")
			for _, a := range typeAnnotations {
				agent := "N/A"
				if a.AgentName.Valid {
					agent = a.AgentName.String
				}
				content := truncate(a.Content, 40)
				metadata := truncate(string(a.Metadata), 30)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
					agent, content, metadata, a.CreatedAt.Format("2006-01-02 15:04"))
			}
		} else {
			fmt.Fprintln(w, "Agent\tContent")
			for _, a := range typeAnnotations {
				agent := "N/A"
				if a.AgentName.Valid {
					agent = a.AgentName.String
				}
				content := truncate(a.Content, 60)
				fmt.Fprintf(w, "%s\t%s\n", agent, content)
			}
		}
		w.Flush()
	}
}

func showGraph(ctx context.Context, queries *db.Queries, projectID, versionID string) {
	fmt.Println("=== NARRATIVE GRAPH ===")
	
	if versionID == "" && projectID != "" {
		workingSet, err := queries.GetWorkingSetVersion(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to get working set for project %s: %v", projectID, err)
		}
		versionID = workingSet.ID
	}
	
	if versionID == "" {
		fmt.Println("Please specify either -project or -version")
		return
	}

	// Get entities and relationships
	entities, err := queries.ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		log.Fatalf("Failed to list entities: %v", err)
	}

	relationships, err := queries.ListRelationshipsByVersion(ctx, versionID)
	if err != nil {
		log.Fatalf("Failed to list relationships: %v", err)
	}

	// Create entity lookup
	entityLookup := make(map[string]db.Entity)
	for _, entity := range entities {
		entityLookup[entity.ID] = entity
	}

	fmt.Printf("Graph for version: %s\n", versionID)
	fmt.Printf("Entities: %d, Relationships: %d\n\n", len(entities), len(relationships))

	// Show graph structure
	for _, rel := range relationships {
		fromEntity := entityLookup[rel.FromEntityID]
		toEntity := entityLookup[rel.ToEntityID]
		
		fmt.Printf("%s (%s) --%s--> %s (%s)\n", 
			fromEntity.Name, fromEntity.EntityType,
			rel.RelationshipType,
			toEntity.Name, toEntity.EntityType)
	}
}

func showStats(ctx context.Context, queries *db.Queries, projectID, versionID string) {
	fmt.Println("=== STATISTICS ===")
	
	if versionID == "" && projectID != "" {
		workingSet, err := queries.GetWorkingSetVersion(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to get working set for project %s: %v", projectID, err)
		}
		versionID = workingSet.ID
	}
	
	if versionID == "" {
		fmt.Println("Please specify either -project or -version")
		return
	}

	// Entity counts by type
	entityTypes := []string{"Scene", "Character", "Location", "Theme", "PlotPoint", "Arc"}
	
	fmt.Println("Entity Counts:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Type\tCount")
	
	totalEntities := int64(0)
	for _, entityType := range entityTypes {
		params := db.CountEntitiesByTypeParams{
			VersionID:  versionID,
			EntityType: entityType,
		}
		count, err := queries.CountEntitiesByType(ctx, params)
		if err != nil {
			log.Printf("Failed to count %s entities: %v", entityType, err)
			continue
		}
		fmt.Fprintf(w, "%s\t%d\n", entityType, count)
		totalEntities += count
	}
	fmt.Fprintf(w, "TOTAL\t%d\n", totalEntities)
	w.Flush()

	// Relationship counts by type
	relationships, err := queries.ListRelationshipsByVersion(ctx, versionID)
	if err != nil {
		log.Printf("Failed to list relationships: %v", err)
		return
	}

	relTypeCounts := make(map[string]int)
	for _, rel := range relationships {
		relTypeCounts[rel.RelationshipType]++
	}

	fmt.Println("\nRelationship Counts:")
	w2 := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w2, "Type\tCount")
	
	totalRels := 0
	for relType, count := range relTypeCounts {
		fmt.Fprintf(w2, "%s\t%d\n", relType, count)
		totalRels += count
	}
	fmt.Fprintf(w2, "TOTAL\t%d\n", totalRels)
	w2.Flush()
}

func getDataPreview(data json.RawMessage, entityType string) string {
	switch entityType {
	case "Scene":
		if sceneData, err := types.UnmarshalSceneData(data); err == nil {
			return fmt.Sprintf("Act: %s, Seq: %d", sceneData.Act, sceneData.Sequence)
		}
	case "Character":
		if charData, err := types.UnmarshalCharacterData(data); err == nil {
			return fmt.Sprintf("Role: %s", charData.Role)
		}
	case "Location":
		if locData, err := types.UnmarshalLocationData(data); err == nil {
			return fmt.Sprintf("Atmosphere: %s", locData.Atmosphere)
		}
	case "Theme":
		if themeData, err := types.UnmarshalThemeData(data); err == nil {
			return fmt.Sprintf("Relevance: %.2f", themeData.Relevance)
		}
	}
	return truncate(string(data), 30)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}