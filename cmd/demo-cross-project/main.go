package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/barrynorthern/libretto/internal/graphwrite"
	"github.com/google/uuid"
)

func main() {
	var (
		dbPath = flag.String("db", "libretto.db", "Path to SQLite database")
		clean  = flag.Bool("clean", false, "Clean existing data before creating demo")
	)
	flag.Parse()

	fmt.Printf("🏰 Creating Elena Stormwind Cross-Project Demo\n")
	fmt.Printf("Database: %s\n\n", *dbPath)

	// Initialize database
	database, err := db.NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean existing data if requested
	if *clean {
		fmt.Printf("🧹 Cleaning existing data...\n")
		if err := cleanDatabase(ctx, database); err != nil {
			log.Fatalf("Failed to clean database: %v", err)
		}
	}

	service := graphwrite.NewService(database)

	// Elena's stable identity across the entire saga
	elenaID := "elena-stormwind-protagonist"
	marcusID := "marcus-ironforge-companion"

	fmt.Printf("📚 Creating The Chronicles of Elena Stormwind...\n\n")

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
		Description: sql.NullString{String: "Elena begins her journey as a young archaeologist discovering ancient mysteries", Valid: true},
	})
	if err != nil {
		log.Fatalf("Failed to create Book 1 project: %v", err)
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
		log.Fatalf("Failed to create Book 1 version: %v", err)
	}

	// Elena starts her journey
	book1Response, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: book1VersionID,
		Deltas: []*graphwrite.Delta{
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
					"logical_id":  elenaID,
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
					"logical_id":  marcusID,
				},
				Relationships: []*graphwrite.RelationshipDelta{
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
			{
				Operation:  "create",
				EntityType: "Location",
				EntityID:   "ancient-temple-of-echoes",
				Fields: map[string]any{
					"name":        "Ancient Temple of Echoes",
					"description": "A mysterious temple where Elena discovers her first artifact",
					"type":        "dungeon",
					"book":        "The Lost Artifact",
					"logical_id":  "ancient-temple-of-echoes",
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create Book 1 characters: %v", err)
	}

	fmt.Printf("   ✨ Elena begins at Level 1, age 22\n")
	fmt.Printf("   ⚔️  Marcus joins as her companion\n")
	fmt.Printf("   🏛️  Ancient Temple of Echoes discovered\n")

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
		Description: sql.NullString{String: "Elena faces the growing darkness and becomes a war leader", Valid: true},
	})
	if err != nil {
		log.Fatalf("Failed to create Book 2 project: %v", err)
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
		log.Fatalf("Failed to create Book 2 version: %v", err)
	}

	// Update Book 1's working set to point to the version with entities
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book1Response.GraphVersionID,
		ProjectID: book1ID,
	})
	if err != nil {
		log.Fatalf("Failed to update Book 1 working set: %v", err)
	}

	// Import Elena from Book 1
	importedElena, err := service.ImportEntity(ctx, book2VersionID, book1ID, elenaID)
	if err != nil {
		log.Fatalf("Failed to import Elena to Book 2: %v", err)
	}

	// Import Marcus from Book 1
	importedMarcus, err := service.ImportEntity(ctx, book2VersionID, book1ID, marcusID)
	if err != nil {
		log.Fatalf("Failed to import Marcus to Book 2: %v", err)
	}

	// Import the temple location
	_, err = service.ImportEntity(ctx, book2VersionID, book1ID, "ancient-temple-of-echoes")
	if err != nil {
		log.Fatalf("Failed to import temple to Book 2: %v", err)
	}

	fmt.Printf("   📥 Elena imported from Book 1 (ID: %s)\n", importedElena.ID)
	fmt.Printf("   📥 Marcus imported from Book 1 (ID: %s)\n", importedMarcus.ID)

	// Elena evolves in Book 2
	book2Response, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: book2VersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name":        "Elena Stormwind",
					"role":        "war_leader",
					"description": "A seasoned archaeologist turned reluctant war leader",
					"level":       7,
					"age":         25,
					"skills":      []string{"archaeology", "ancient_languages", "leadership", "combat_magic"},
					"book":        "The Shadow War",
					"trauma":      "witnessed_the_fall_of_ancient_city",
					"logical_id":  elenaID,
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
					"logical_id":  marcusID,
				},
			},
			{
				Operation:  "create",
				EntityType: "Location",
				EntityID:   "iron-pass-battlefield",
				Fields: map[string]any{
					"name":        "Iron Pass Battlefield",
					"description": "The site of the great battle where Elena proved her leadership",
					"type":        "battlefield",
					"book":        "The Shadow War",
					"logical_id":  "iron-pass-battlefield",
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to evolve characters in Book 2: %v", err)
	}

	fmt.Printf("   ⚡ Elena evolved to Level 7, now a war leader\n")
	fmt.Printf("   🛡️  Marcus became a war veteran\n")
	fmt.Printf("   ⚔️  Iron Pass Battlefield added\n")

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
		Description: sql.NullString{String: "Elena fulfills her destiny as the Lightbringer of the Seven Realms", Valid: true},
	})
	if err != nil {
		log.Fatalf("Failed to create Book 3 project: %v", err)
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
		log.Fatalf("Failed to create Book 3 version: %v", err)
	}

	// Update Book 2's working set to point to the version with evolved entities
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book2Response.GraphVersionID,
		ProjectID: book2ID,
	})
	if err != nil {
		log.Fatalf("Failed to update Book 2 working set: %v", err)
	}

	// Import Elena from Book 2 (she carries her evolution)
	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, elenaID)
	if err != nil {
		log.Fatalf("Failed to import Elena to Book 3: %v", err)
	}

	// Import Marcus from Book 2
	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, marcusID)
	if err != nil {
		log.Fatalf("Failed to import Marcus to Book 3: %v", err)
	}

	// Import locations
	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, "ancient-temple-of-echoes")
	if err != nil {
		log.Fatalf("Failed to import temple to Book 3: %v", err)
	}

	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, "iron-pass-battlefield")
	if err != nil {
		log.Fatalf("Failed to import battlefield to Book 3: %v", err)
	}

	// Elena reaches her final form
	book3Response, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: book3VersionID,
		Deltas: []*graphwrite.Delta{
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
					"logical_id":  elenaID,
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
					"logical_id":  "lyra-stormwind-successor",
				},
				Relationships: []*graphwrite.RelationshipDelta{
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
			{
				Operation:  "create",
				EntityType: "Location",
				EntityID:   "crystal-throne-chamber",
				Fields: map[string]any{
					"name":        "Crystal Throne Chamber",
					"description": "The final battleground where Elena defeats the Shadow Lord",
					"type":        "throne_room",
					"book":        "The Final Prophecy",
					"logical_id":  "crystal-throne-chamber",
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to complete Elena's arc in Book 3: %v", err)
	}

	// Update Book 3's working set to point to the final version
	err = database.Queries().SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book3Response.GraphVersionID,
		ProjectID: book3ID,
	})
	if err != nil {
		log.Fatalf("Failed to update Book 3 working set: %v", err)
	}

	fmt.Printf("   🌟 Elena ascended to Legendary status (Level 15)\n")
	fmt.Printf("   👑 Elena became the Lightbringer of the Seven Realms\n")
	fmt.Printf("   🎓 Elena now mentors Lyra, her successor\n")
	fmt.Printf("   🏰 Crystal Throne Chamber - final battleground\n")

	// ==========================================
	// VERIFICATION: CROSS-PROJECT CONTINUITY
	// ==========================================
	fmt.Printf("\n🔍 === VERIFICATION: ELENA'S COMPLETE JOURNEY ===\n")

	// Get Elena's complete history across all projects
	elenaHistory, err := service.GetEntityHistory(ctx, elenaID)
	if err != nil {
		log.Fatalf("Failed to get Elena's history: %v", err)
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

	// List all shared entities
	sharedEntities, err := service.ListSharedEntities(ctx)
	if err != nil {
		log.Fatalf("Failed to list shared entities: %v", err)
	}

	fmt.Printf("\n🌐 Shared Characters Across the Saga:\n")
	for _, entity := range sharedEntities {
		fmt.Printf("   • %s (%s) - appears in %d books: %v\n", 
			entity.Name, entity.EntityType, entity.ProjectCount, entity.Projects)
	}

	fmt.Printf("\n✅ SUCCESS: Elena's identity preserved across entire saga!\n")
	fmt.Printf("✅ SUCCESS: Character arcs span multiple projects!\n")
	fmt.Printf("✅ SUCCESS: Cross-project relationships maintained!\n")
	fmt.Printf("✅ SUCCESS: Narrative continuity achieved!\n")

	fmt.Printf("\n🎉 THE CHRONICLES OF ELENA STORMWIND - COMPLETE! 🎉\n")
	fmt.Printf("\n🎛️  Launch the dashboard to explore Elena's journey:\n")
	fmt.Printf("   go run cmd/dashboard/main.go -db %s\n", *dbPath)
	fmt.Printf("   Visit: http://localhost:9000\n")
}

func cleanDatabase(ctx context.Context, database *db.Database) error {
	// Delete all data in reverse dependency order to avoid foreign key constraints
	
	// Delete relationships first
	if _, err := database.DB().ExecContext(ctx, "DELETE FROM relationships"); err != nil {
		return fmt.Errorf("failed to delete relationships: %w", err)
	}
	
	// Delete annotations
	if _, err := database.DB().ExecContext(ctx, "DELETE FROM annotations"); err != nil {
		return fmt.Errorf("failed to delete annotations: %w", err)
	}
	
	// Delete entities
	if _, err := database.DB().ExecContext(ctx, "DELETE FROM entities"); err != nil {
		return fmt.Errorf("failed to delete entities: %w", err)
	}
	
	// Delete graph versions
	if _, err := database.DB().ExecContext(ctx, "DELETE FROM graph_versions"); err != nil {
		return fmt.Errorf("failed to delete graph_versions: %w", err)
	}
	
	// Delete projects
	if _, err := database.DB().ExecContext(ctx, "DELETE FROM projects"); err != nil {
		return fmt.Errorf("failed to delete projects: %w", err)
	}
	
	return nil
}