package graphwrite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/google/uuid"
)

// TestStableEntityIDs demonstrates that Elena remains Elena across versions
func TestStableEntityIDs(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_stable_id_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Initialize database
	database, err := db.NewDatabase(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Create project and initial version
	projectID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "Elena's Story",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "Following Elena's journey", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	initialVersionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           initialVersionID,
		ProjectID:    projectID,
		Name:         sql.NullString{String: "Chapter 1", Valid: true},
		Description:  sql.NullString{String: "Elena's introduction", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create initial version: %v", err)
	}

	service := NewService(database)

	// Elena's logical ID - this should remain constant
	elenaID := "elena-protagonist-001"
	sceneID := "opening-scene-001"

	fmt.Printf("=== Creating Elena in Version 1 ===\n")
	
	// Version 1: Create Elena
	response1, err := service.Apply(ctx, &ApplyRequest{
		ParentVersionID: initialVersionID,
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   elenaID, // Elena's stable logical ID
				Fields: map[string]any{
					"name":        "Elena",
					"role":        "protagonist",
					"description": "A brave archaeologist",
					"level":       1,
				},
			},
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   sceneID,
				Fields: map[string]any{
					"name":  "Opening Scene",
					"title": "The Journey Begins",
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     sceneID,
						ToEntityID:       elenaID,
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
		t.Fatalf("Failed to create version 1: %v", err)
	}

	// Verify Elena exists in version 1
	entities1, err := service.ListEntities(ctx, response1.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("Failed to list entities in version 1: %v", err)
	}

	var elena1 *Entity
	for _, entity := range entities1 {
		if entity.ID == elenaID {
			elena1 = entity
			break
		}
	}

	if elena1 == nil {
		t.Fatalf("Elena not found in version 1")
	}

	fmt.Printf("Version 1 - Elena ID: %s, Level: %v\n", elena1.ID, elena1.Data["level"])

	fmt.Printf("\n=== Updating Elena in Version 2 ===\n")

	// Version 2: Elena levels up
	response2, err := service.Apply(ctx, &ApplyRequest{
		ParentVersionID: response1.GraphVersionID,
		Deltas: []*Delta{
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   elenaID, // Same logical ID
				Fields: map[string]any{
					"name":        "Elena",
					"role":        "protagonist",
					"description": "A brave archaeologist with growing confidence",
					"level":       2, // She leveled up!
					"skills":      []string{"archaeology", "combat"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create version 2: %v", err)
	}

	// Verify Elena exists in version 2 with same ID
	entities2, err := service.ListEntities(ctx, response2.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("Failed to list entities in version 2: %v", err)
	}

	var elena2 *Entity
	for _, entity := range entities2 {
		if entity.ID == elenaID {
			elena2 = entity
			break
		}
	}

	if elena2 == nil {
		t.Fatalf("Elena not found in version 2")
	}

	fmt.Printf("Version 2 - Elena ID: %s, Level: %v\n", elena2.ID, elena2.Data["level"])

	fmt.Printf("\n=== Adding Elena's companion in Version 3 ===\n")

	// Version 3: Elena gets a companion
	companionID := "marcus-companion-001"
	response3, err := service.Apply(ctx, &ApplyRequest{
		ParentVersionID: response2.GraphVersionID,
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   companionID,
				Fields: map[string]any{
					"name":        "Marcus",
					"role":        "companion",
					"description": "Elena's trusted ally",
					"level":       1,
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     elenaID, // Elena's stable ID
						ToEntityID:       companionID,
						RelationshipType: "allies_with",
						Properties: map[string]any{
							"bond_strength": "strong",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create version 3: %v", err)
	}

	// Verify Elena still exists with same ID in version 3
	entities3, err := service.ListEntities(ctx, response3.GraphVersionID, EntityFilter{})
	if err != nil {
		t.Fatalf("Failed to list entities in version 3: %v", err)
	}

	var elena3 *Entity
	for _, entity := range entities3 {
		if entity.ID == elenaID {
			elena3 = entity
			break
		}
	}

	if elena3 == nil {
		t.Fatalf("Elena not found in version 3")
	}

	fmt.Printf("Version 3 - Elena ID: %s, Level: %v\n", elena3.ID, elena3.Data["level"])

	// Verify Elena's relationships
	companions, err := service.GetNeighborsInVersion(ctx, response3.GraphVersionID, elenaID, "allies_with")
	if err != nil {
		t.Fatalf("Failed to get Elena's companions: %v", err)
	}

	if len(companions) != 1 {
		t.Errorf("Expected Elena to have 1 companion, got %d", len(companions))
	}

	if len(companions) > 0 {
		fmt.Printf("Elena's companion: %s (ID: %s)\n", companions[0].Name, companions[0].ID)
	}

	fmt.Printf("\n=== VERIFICATION: Elena's Identity Across Versions ===\n")

	// The key test: Elena's ID should be the same across all versions
	if elena1.ID != elenaID {
		t.Errorf("Elena's ID changed in version 1: expected %s, got %s", elenaID, elena1.ID)
	}

	if elena2.ID != elenaID {
		t.Errorf("Elena's ID changed in version 2: expected %s, got %s", elenaID, elena2.ID)
	}

	if elena3.ID != elenaID {
		t.Errorf("Elena's ID changed in version 3: expected %s, got %s", elenaID, elena3.ID)
	}

	// Elena should have evolved across versions
	level1 := elena1.Data["level"].(float64)
	level2 := elena2.Data["level"].(float64)
	level3 := elena3.Data["level"].(float64)

	if level1 != 1 {
		t.Errorf("Expected Elena to be level 1 in version 1, got %v", level1)
	}

	if level2 != 2 {
		t.Errorf("Expected Elena to be level 2 in version 2, got %v", level2)
	}

	if level3 != 2 {
		t.Errorf("Expected Elena to be level 2 in version 3, got %v", level3)
	}

	fmt.Printf("✅ SUCCESS: Elena maintains her identity (%s) across all versions!\n", elenaID)
	fmt.Printf("✅ SUCCESS: Elena's character development is preserved!\n")
	fmt.Printf("✅ SUCCESS: Elena's relationships are maintained!\n")

	// Output:
	// === Creating Elena in Version 1 ===
	// Version 1 - Elena ID: elena-protagonist-001, Level: 1
	// 
	// === Updating Elena in Version 2 ===
	// Version 2 - Elena ID: elena-protagonist-001, Level: 2
	// 
	// === Adding Elena's companion in Version 3 ===
	// Version 3 - Elena ID: elena-protagonist-001, Level: 2
	// Elena's companion: Marcus (ID: marcus-companion-001)
	// 
	// === VERIFICATION: Elena's Identity Across Versions ===
	// ✅ SUCCESS: Elena maintains her identity (elena-protagonist-001) across all versions!
	// ✅ SUCCESS: Elena's character development is preserved!
	// ✅ SUCCESS: Elena's relationships are maintained!
}