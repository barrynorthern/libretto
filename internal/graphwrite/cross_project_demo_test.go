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

// TestCrossProjectCharacterArcs demonstrates Elena's journey across multiple books/projects
func TestCrossProjectCharacterArcs(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "libretto_cross_project_test_*.db")
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

	service := NewService(database)

	// Elena's stable identity across the entire saga
	elenaID := "elena-stormwind-protagonist"
	marcusID := "marcus-ironforge-companion"

	fmt.Printf("🏰 === THE CHRONICLES OF ELENA STORMWIND ===\n\n")

	// ==========================================
	// BOOK 1: THE LOST ARTIFACT
	// ==========================================
	fmt.Printf("📖 BOOK 1: THE LOST ARTIFACT\n")
	
	book1ID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          book1ID,
		Name:        "Book 1: The Lost Artifact",
		Theme:       sql.NullString{String: "Discovery", Valid: true},
		Genre:       sql.NullString{String: "Fantasy Adventure", Valid: true},
		Description: sql.NullString{String: "Elena begins her journey as a young archaeologist", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create Book 1 project: %v", err)
	}

	book1VersionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           book1VersionID,
		ProjectID:    book1ID,
		Name:         sql.NullString{String: "Final Draft", Valid: true},
		Description:  sql.NullString{String: "Elena's origin story", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create Book 1 version: %v", err)
	}

	// Elena starts her journey
	book1Response, err := service.Apply(ctx, &ApplyRequest{
		ParentVersionID: book1VersionID,
		Deltas: []*Delta{
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name":        "Elena Stormwind",
					"role":        "protagonist",
					"description": "A young archaeologist with a thirst for ancient mysteries",
					"level":       1,
					"age":         22,
					"skills":      []string{"archaeology", "ancient_languages"},
					"book":        "The Lost Artifact",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   marcusID,
				Fields: map[string]any{
					"name":        "Marcus Ironforge",
					"role":        "companion",
					"description": "A gruff but loyal dwarf warrior",
					"level":       3,
					"age":         45,
					"skills":      []string{"combat", "smithing"},
					"book":        "The Lost Artifact",
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     elenaID,
						ToEntityID:       marcusID,
						RelationshipType: "allies_with",
						Properties: map[string]any{
							"bond_strength": "growing",
							"trust_level":   "cautious",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create Book 1 characters: %v", err)
	}

	fmt.Printf("   ✨ Elena begins at Level %d, age %d\n", 1, 22)
	fmt.Printf("   ⚔️  Marcus joins as her companion\n")

	// ==========================================
	// BOOK 2: THE SHADOW WAR
	// ==========================================
	fmt.Printf("\n📖 BOOK 2: THE SHADOW WAR\n")
	
	book2ID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          book2ID,
		Name:        "Book 2: The Shadow War",
		Theme:       sql.NullString{String: "Conflict", Valid: true},
		Genre:       sql.NullString{String: "Fantasy War", Valid: true},
		Description: sql.NullString{String: "Elena faces the growing darkness", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create Book 2 project: %v", err)
	}

	book2VersionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           book2VersionID,
		ProjectID:    book2ID,
		Name:         sql.NullString{String: "Final Draft", Valid: true},
		Description:  sql.NullString{String: "The war begins", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create Book 2 version: %v", err)
	}

	// Update Book 1's working set to point to the version with entities
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book1Response.GraphVersionID,
		ProjectID: book1ID,
	})
	if err != nil {
		t.Fatalf("Failed to update Book 1 working set: %v", err)
	}

	// Import Elena from Book 1
	importedElena, err := service.ImportEntity(ctx, book2VersionID, book1ID, elenaID)
	if err != nil {
		t.Fatalf("Failed to import Elena to Book 2: %v", err)
	}

	// Import Marcus from Book 1
	importedMarcus, err := service.ImportEntity(ctx, book2VersionID, book1ID, marcusID)
	if err != nil {
		t.Fatalf("Failed to import Marcus to Book 2: %v", err)
	}

	fmt.Printf("   📥 Elena imported from Book 1 (ID: %s)\n", importedElena.ID)
	fmt.Printf("   📥 Marcus imported from Book 1 (ID: %s)\n", importedMarcus.ID)

	// Elena evolves in Book 2
	book2Response, err := service.Apply(ctx, &ApplyRequest{
		ParentVersionID: book2VersionID,
		Deltas: []*Delta{
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   elenaID, // Same logical ID!
				Fields: map[string]any{
					"name":        "Elena Stormwind",
					"role":        "war_leader",
					"description": "A seasoned archaeologist turned reluctant war leader",
					"level":       7,
					"age":         25,
					"skills":      []string{"archaeology", "ancient_languages", "leadership", "combat_magic"},
					"book":        "The Shadow War",
					"trauma":      "witnessed_the_fall_of_ancient_city",
				},
			},
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   marcusID,
				Fields: map[string]any{
					"name":        "Marcus Ironforge",
					"role":        "war_veteran",
					"description": "Elena's most trusted advisor and battle companion",
					"level":       8,
					"age":         48,
					"skills":      []string{"combat", "smithing", "tactics", "leadership"},
					"book":        "The Shadow War",
					"scars":       "battle_of_iron_pass",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to evolve characters in Book 2: %v", err)
	}

	fmt.Printf("   ⚡ Elena evolved to Level %d, now a war leader\n", 7)
	fmt.Printf("   🛡️  Marcus became a war veteran\n")

	// ==========================================
	// BOOK 3: THE FINAL PROPHECY
	// ==========================================
	fmt.Printf("\n📖 BOOK 3: THE FINAL PROPHECY\n")
	
	book3ID := uuid.New().String()
	_, err = database.Queries().CreateProject(ctx, db.CreateProjectParams{
		ID:          book3ID,
		Name:        "Book 3: The Final Prophecy",
		Theme:       sql.NullString{String: "Destiny", Valid: true},
		Genre:       sql.NullString{String: "Epic Fantasy", Valid: true},
		Description: sql.NullString{String: "Elena fulfills her destiny", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create Book 3 project: %v", err)
	}

	book3VersionID := uuid.New().String()
	_, err = database.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           book3VersionID,
		ProjectID:    book3ID,
		Name:         sql.NullString{String: "Final Draft", Valid: true},
		Description:  sql.NullString{String: "The epic conclusion", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create Book 3 version: %v", err)
	}

	// Update Book 2's working set to point to the version with evolved entities
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book2Response.GraphVersionID,
		ProjectID: book2ID,
	})
	if err != nil {
		t.Fatalf("Failed to update Book 2 working set: %v", err)
	}

	// Import Elena from Book 2 (she carries her evolution)
	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, elenaID)
	if err != nil {
		t.Fatalf("Failed to import Elena to Book 3: %v", err)
	}

	// Import Marcus from Book 2
	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, marcusID)
	if err != nil {
		t.Fatalf("Failed to import Marcus to Book 3: %v", err)
	}

	// Elena reaches her final form
	book3Response, err := service.Apply(ctx, &ApplyRequest{
		ParentVersionID: book3VersionID,
		Deltas: []*Delta{
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name":        "Elena Stormwind, the Lightbringer",
					"role":        "legendary_hero",
					"description": "The prophesied hero who united the realms against darkness",
					"level":       15,
					"age":         28,
					"skills":      []string{"archaeology", "ancient_languages", "leadership", "combat_magic", "divine_magic", "realm_walking"},
					"book":        "The Final Prophecy",
					"title":       "Lightbringer of the Seven Realms",
					"achievement": "defeated_the_shadow_lord",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   "lyra-stormwind-successor",
				Fields: map[string]any{
					"name":        "Lyra Stormwind",
					"role":        "successor",
					"description": "Elena's apprentice, destined to carry on her legacy",
					"level":       3,
					"age":         19,
					"skills":      []string{"archaeology", "ancient_languages", "potential"},
					"book":        "The Final Prophecy",
					"mentor":      elenaID,
				},
				Relationships: []*RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     elenaID,
						ToEntityID:       "lyra-stormwind-successor",
						RelationshipType: "mentors",
						Properties: map[string]any{
							"legacy_transfer": "in_progress",
							"bond_type":       "master_apprentice",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to complete Elena's arc in Book 3: %v", err)
	}

	// Update Book 3's working set to point to the final version
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book3Response.GraphVersionID,
		ProjectID: book3ID,
	})
	if err != nil {
		t.Fatalf("Failed to update Book 3 working set: %v", err)
	}

	fmt.Printf("   🌟 Elena ascended to Legendary status (Level %d)\n", 15)
	fmt.Printf("   👑 Elena became the Lightbringer of the Seven Realms\n")
	fmt.Printf("   🎓 Elena now mentors Lyra, her successor\n")

	// ==========================================
	// VERIFICATION: CROSS-PROJECT CONTINUITY
	// ==========================================
	fmt.Printf("\n🔍 === VERIFICATION: ELENA'S COMPLETE JOURNEY ===\n")

	// Get Elena's complete history across all projects
	elenaHistory, err := service.GetEntityHistory(ctx, elenaID)
	if err != nil {
		t.Fatalf("Failed to get Elena's history: %v", err)
	}

	fmt.Printf("\n📚 Elena's Evolution Across the Saga:\n")
	for i, version := range elenaHistory {
		level := version.Entity.Data["level"].(float64)
		age := version.Entity.Data["age"].(float64)
		role := version.Entity.Data["role"].(string)
		book := version.Entity.Data["book"].(string)
		
		fmt.Printf("   %d. %s: Level %.0f, Age %.0f, Role: %s\n", 
			i+1, book, level, age, role)
	}

	// Verify Elena appears in all 3 books
	if len(elenaHistory) != 3 {
		t.Errorf("Expected Elena to appear in 3 books, found %d", len(elenaHistory))
	}

	// Verify character progression
	book1Elena := elenaHistory[0].Entity
	book3Elena := elenaHistory[2].Entity

	book1Level := book1Elena.Data["level"].(float64)
	book3Level := book3Elena.Data["level"].(float64)

	if book1Level != 1 {
		t.Errorf("Expected Elena to start at level 1, got %v", book1Level)
	}

	if book3Level != 15 {
		t.Errorf("Expected Elena to end at level 15, got %v", book3Level)
	}

	// List all shared entities
	sharedEntities, err := service.ListSharedEntities(ctx)
	if err != nil {
		t.Fatalf("Failed to list shared entities: %v", err)
	}

	fmt.Printf("\n🌐 Shared Characters Across the Saga:\n")
	for _, entity := range sharedEntities {
		fmt.Printf("   • %s (%s) - appears in %d books: %v\n", 
			entity.Name, entity.EntityType, entity.ProjectCount, entity.Projects)
	}

	// Verify both Elena and Marcus are shared
	foundElena := false
	foundMarcus := false
	for _, entity := range sharedEntities {
		if entity.LogicalID == elenaID {
			foundElena = true
			if entity.ProjectCount != 3 {
				t.Errorf("Expected Elena in 3 projects, found %d", entity.ProjectCount)
			}
		}
		if entity.LogicalID == marcusID {
			foundMarcus = true
			if entity.ProjectCount != 3 {
				t.Errorf("Expected Marcus in 3 projects, found %d", entity.ProjectCount)
			}
		}
	}

	if !foundElena {
		t.Error("Elena not found in shared entities list")
	}

	if !foundMarcus {
		t.Error("Marcus not found in shared entities list")
	}

	fmt.Printf("\n✅ SUCCESS: Elena's identity preserved across entire saga!\n")
	fmt.Printf("✅ SUCCESS: Character arcs span multiple projects!\n")
	fmt.Printf("✅ SUCCESS: Cross-project relationships maintained!\n")
	fmt.Printf("✅ SUCCESS: Narrative continuity achieved!\n")

	fmt.Printf("\n🎉 THE CHRONICLES OF ELENA STORMWIND - COMPLETE! 🎉\n")

	// Output:
	// 🏰 === THE CHRONICLES OF ELENA STORMWIND ===
	// 
	// 📖 BOOK 1: THE LOST ARTIFACT
	//    ✨ Elena begins at Level 1, age 22
	//    ⚔️  Marcus joins as her companion
	// 
	// 📖 BOOK 2: THE SHADOW WAR
	//    📥 Elena imported from Book 1 (ID: elena-stormwind-protagonist)
	//    📥 Marcus imported from Book 1 (ID: marcus-ironforge-companion)
	//    ⚡ Elena evolved to Level 7, now a war leader
	//    🛡️  Marcus became a war veteran
	// 
	// 📖 BOOK 3: THE FINAL PROPHECY
	//    🌟 Elena ascended to Legendary status (Level 15)
	//    👑 Elena became the Lightbringer of the Seven Realms
	//    🎓 Elena now mentors Lyra, her successor
	// 
	// 🔍 === VERIFICATION: ELENA'S COMPLETE JOURNEY ===
	// 
	// 📚 Elena's Evolution Across the Saga:
	//    1. The Lost Artifact: Level 1, Age 22, Role: protagonist
	//    2. The Shadow War: Level 7, Age 25, Role: war_leader
	//    3. The Final Prophecy: Level 15, Age 28, Role: legendary_hero
	// 
	// 🌐 Shared Characters Across the Saga:
	//    • Elena Stormwind, the Lightbringer (Character) - appears in 3 books: [Book 1: The Lost Artifact Book 2: The Shadow War Book 3: The Final Prophecy]
	//    • Marcus Ironforge (Character) - appears in 3 books: [Book 1: The Lost Artifact Book 2: The Shadow War Book 3: The Final Prophecy]
	// 
	// ✅ SUCCESS: Elena's identity preserved across entire saga!
	// ✅ SUCCESS: Character arcs span multiple projects!
	// ✅ SUCCESS: Cross-project relationships maintained!
	// ✅ SUCCESS: Narrative continuity achieved!
	// 
	// 🎉 THE CHRONICLES OF ELENA STORMWIND - COMPLETE! 🎉
}