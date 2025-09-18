package graphwrite_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/graphwrite"
	"github.com/google/uuid"
)

// Example demonstrates how to use the GraphWrite service to create and manage
// a narrative graph with versioning.
func Example() {
	// Create temporary database for this example
	tmpFile, err := os.CreateTemp("", "libretto_example_*.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Initialize database and run migrations
	database, err := db.NewDatabase(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	// Create a project
	projectID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "Example Story",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "An example fantasy adventure", Valid: true},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create initial graph version
	initialVersionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           initialVersionID,
		ProjectID:    projectID,
		Name:         sql.NullString{String: "Initial Version", Valid: true},
		Description:  sql.NullString{String: "Starting point for our story", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize GraphWrite service
	service := graphwrite.NewService(database)

	// Create the first version with a scene and character
	sceneID := uuid.New().String()
	characterID := uuid.New().String()

	response1, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: initialVersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   sceneID,
				Fields: map[string]any{
					"name":    "Opening Scene",
					"title":   "The Tavern",
					"summary": "Our hero enters a mysterious tavern",
					"content": "The wooden door creaked as Elena pushed it open...",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   characterID,
				Fields: map[string]any{
					"name":        "Elena",
					"role":        "protagonist",
					"description": "A brave archaeologist seeking ancient artifacts",
				},
				Relationships: []*graphwrite.RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     sceneID,
						ToEntityID:       characterID,
						RelationshipType: "features",
						Properties: map[string]any{
							"importance": "primary",
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created version 1 with %d entities\n", response1.Applied)

	// Create a second version that updates the scene
	response2, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: response1.GraphVersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "update",
				EntityType: "Scene",
				EntityID:   sceneID, // Reference original ID, service will map to new version
				Fields: map[string]any{
					"name":    "Opening Scene - Revised",
					"title":   "The Mysterious Tavern",
					"summary": "Our hero enters a tavern filled with strange patrons",
					"content": "The wooden door creaked ominously as Elena pushed it open, revealing a dimly lit tavern filled with hooded figures...",
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created version 2 with %d updates\n", response2.Applied)

	// List entities in the latest version
	entities, err := service.ListEntities(ctx, response2.GraphVersionID, graphwrite.EntityFilter{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Latest version contains %d entities:\n", len(entities))
	for _, entity := range entities {
		fmt.Printf("- %s: %s\n", entity.EntityType, entity.Name)
	}

	// Get neighbors of the scene
	neighbors, err := service.GetNeighborsInVersion(ctx, response2.GraphVersionID, sceneID, "features")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Scene features %d characters:\n", len(neighbors))
	for _, neighbor := range neighbors {
		fmt.Printf("- %s\n", neighbor.Name)
	}

	// Get version information
	version, err := service.GetVersion(ctx, response2.GraphVersionID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version has parent: %t\n", version.ParentVersionID != nil)

	// Output:
	// Created version 1 with 2 entities
	// Created version 2 with 1 updates
	// Latest version contains 2 entities:
	// - Character: Elena
	// - Scene: Opening Scene - Revised
	// Scene features 1 characters:
	// - Elena
	// Version has parent: true
}