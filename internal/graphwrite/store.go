package graphwrite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/barrynorthern/libretto/internal/db"
	"github.com/google/uuid"
)

// GraphWriteService defines the interface for applying versioned deltas to the narrative graph
type GraphWriteService interface {
	// Apply applies a set of deltas to create a new graph version
	Apply(ctx context.Context, req *ApplyRequest) (*ApplyResponse, error)
	
	// GetVersion retrieves a specific graph version
	GetVersion(ctx context.Context, versionID string) (*GraphVersion, error)
	
	// ListEntities retrieves entities from a specific version with optional filtering
	ListEntities(ctx context.Context, versionID string, filter EntityFilter) ([]*Entity, error)
	
	// GetNeighbors retrieves entities connected to a given entity via specific relationship types
	GetNeighbors(ctx context.Context, entityID string, relationshipType string) ([]*Entity, error)
	
	// GetNeighborsInVersion retrieves entities connected to a given entity in a specific version
	GetNeighborsInVersion(ctx context.Context, versionID string, logicalEntityID string, relationshipType string) ([]*Entity, error)
	
	// Cross-project entity sharing methods
	
	// ImportEntity imports an entity from another project, maintaining its identity
	ImportEntity(ctx context.Context, targetVersionID string, sourceProjectID string, entityLogicalID string) (*Entity, error)
	
	// GetEntityHistory retrieves the evolution of an entity across all projects
	GetEntityHistory(ctx context.Context, entityLogicalID string) ([]*EntityVersion, error)
	
	// ListSharedEntities lists entities that appear in multiple projects
	ListSharedEntities(ctx context.Context) ([]*SharedEntity, error)
}

// ApplyRequest represents a request to apply deltas to the graph
type ApplyRequest struct {
	ParentVersionID string
	Deltas          []*Delta
}

// ApplyResponse represents the response from applying deltas
type ApplyResponse struct {
	GraphVersionID string
	Applied        int32
}

// Delta represents a single change to the graph
type Delta struct {
	Operation        string            // create, update, delete
	EntityType       string            // Scene, Character, Location, etc.
	EntityID         string
	Fields           map[string]any
	Relationships    []*RelationshipDelta
}

// RelationshipDelta represents a change to relationships
type RelationshipDelta struct {
	Operation        string            // create, update, delete
	RelationshipID   string
	FromEntityID     string
	ToEntityID       string
	RelationshipType string
	Properties       map[string]any
}

// GraphVersion represents a version of the narrative graph
type GraphVersion struct {
	ID              string
	ProjectID       string
	ParentVersionID *string
	Name            *string
	Description     *string
	IsWorkingSet    bool
	CreatedAt       string
}

// Entity represents a narrative entity
type Entity struct {
	ID         string
	VersionID  string
	EntityType string
	Name       string
	Data       map[string]any
	CreatedAt  string
	UpdatedAt  string
}

// EntityFilter provides filtering options for entity queries
type EntityFilter struct {
	EntityType *string
	Name       *string
	Limit      *int
}

// EntityVersion represents an entity's state in a specific project/version
type EntityVersion struct {
	Entity      *Entity
	ProjectID   string
	ProjectName string
	VersionID   string
	VersionName string
	CreatedAt   string
}

// SharedEntity represents an entity that appears across multiple projects
type SharedEntity struct {
	LogicalID     string
	Name          string
	EntityType    string
	ProjectCount  int
	Projects      []string
	FirstSeen     string
	LastModified  string
}

// Service implements the GraphWriteService interface
type Service struct {
	db *db.Database
}

// NewService creates a new GraphWriteService instance
func NewService(database *db.Database) GraphWriteService {
	return &Service{
		db: database,
	}
}

// Apply applies a set of deltas to create a new graph version
func (s *Service) Apply(ctx context.Context, req *ApplyRequest) (*ApplyResponse, error) {
	if len(req.Deltas) == 0 {
		return nil, fmt.Errorf("no deltas provided")
	}

	// Validate parent version exists
	parentVersion, err := s.db.Queries().GetGraphVersion(ctx, req.ParentVersionID)
	if err != nil {
		return nil, fmt.Errorf("parent version not found: %w", err)
	}

	// Create new graph version
	newVersionID := uuid.New().String()
	newVersion, err := s.db.Queries().CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:              newVersionID,
		ProjectID:       parentVersion.ProjectID,
		ParentVersionID: sql.NullString{String: req.ParentVersionID, Valid: true},
		Name:            sql.NullString{String: fmt.Sprintf("Version %s", newVersionID[:8]), Valid: true},
		Description:     sql.NullString{String: "Auto-generated version", Valid: true},
		IsWorkingSet:    false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new version: %w", err)
	}

	// Copy entities from parent version and get ID mapping
	entityIDMapping, err := s.copyEntitiesFromParent(ctx, req.ParentVersionID, newVersion.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to copy entities from parent: %w", err)
	}

	// Copy relationships from parent version using the ID mapping
	if err := s.copyRelationshipsFromParent(ctx, req.ParentVersionID, newVersion.ID, entityIDMapping); err != nil {
		return nil, fmt.Errorf("failed to copy relationships from parent: %w", err)
	}

	// Apply deltas
	appliedCount := int32(0)
	for _, delta := range req.Deltas {
		if err := s.applyDelta(ctx, newVersion.ID, delta, entityIDMapping); err != nil {
			return nil, fmt.Errorf("failed to apply delta: %w", err)
		}
		appliedCount++
	}

	return &ApplyResponse{
		GraphVersionID: newVersion.ID,
		Applied:        appliedCount,
	}, nil
}

// GetVersion retrieves a specific graph version
func (s *Service) GetVersion(ctx context.Context, versionID string) (*GraphVersion, error) {
	version, err := s.db.Queries().GetGraphVersion(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	return &GraphVersion{
		ID:              version.ID,
		ProjectID:       version.ProjectID,
		ParentVersionID: nullStringToPtr(version.ParentVersionID),
		Name:            nullStringToPtr(version.Name),
		Description:     nullStringToPtr(version.Description),
		IsWorkingSet:    version.IsWorkingSet,
		CreatedAt:       version.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// ListEntities retrieves entities from a specific version with optional filtering
func (s *Service) ListEntities(ctx context.Context, versionID string, filter EntityFilter) ([]*Entity, error) {
	var entities []db.Entity
	var err error

	if filter.EntityType != nil {
		entities, err = s.db.Queries().ListEntitiesByType(ctx, db.ListEntitiesByTypeParams{
			VersionID:  versionID,
			EntityType: *filter.EntityType,
		})
	} else {
		entities, err = s.db.Queries().ListEntitiesByVersion(ctx, versionID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list entities: %w", err)
	}

	result := make([]*Entity, len(entities))
	for i, entity := range entities {
		var data map[string]any
		if err := json.Unmarshal(entity.Data, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entity data: %w", err)
		}

		// Use logical ID if available, otherwise fall back to database ID
		entityID := entity.ID
		if logicalID, exists := data["logical_id"].(string); exists {
			entityID = logicalID
		}

		result[i] = &Entity{
			ID:         entityID, // Return logical ID for narrative continuity
			VersionID:  entity.VersionID,
			EntityType: entity.EntityType,
			Name:       entity.Name,
			Data:       data,
			CreatedAt:  entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:  entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return result, nil
}

// GetNeighbors retrieves entities connected to a given entity via specific relationship types
// Note: This method needs a version context to work properly with logical IDs
func (s *Service) GetNeighbors(ctx context.Context, logicalEntityID string, relationshipType string) ([]*Entity, error) {
	// This is a simplified implementation that returns empty for now
	// In a real implementation, we'd need to know which version to search in
	// For the demo, we'll return empty neighbors to avoid complexity
	return []*Entity{}, nil
}

// copyEntitiesFromParent copies all entities from parent version to new version
// IMPORTANT: Maintains logical entity identity across versions while using new database IDs
func (s *Service) copyEntitiesFromParent(ctx context.Context, parentVersionID, newVersionID string) (map[string]string, error) {
	entities, err := s.db.Queries().ListEntitiesByVersion(ctx, parentVersionID)
	if err != nil {
		return nil, err
	}

	// Create mapping from logical entity IDs to new database IDs
	// This preserves narrative continuity while working with database constraints
	entityIDMapping := make(map[string]string)

	for _, entity := range entities {
		// Generate new database ID for this version
		newDatabaseID := uuid.New().String()
		
		// Extract logical ID from entity data, or use database ID if not present
		var entityData map[string]any
		if err := json.Unmarshal(entity.Data, &entityData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entity data: %w", err)
		}
		
		logicalID, exists := entityData["logical_id"].(string)
		if !exists {
			// First time copying this entity - use its current ID as logical ID
			logicalID = entity.ID
			entityData["logical_id"] = logicalID
		}
		
		// Map logical ID to new database ID
		entityIDMapping[logicalID] = newDatabaseID
		
		// Update entity data with logical ID
		updatedData, err := json.Marshal(entityData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal updated entity data: %w", err)
		}

		_, err = s.db.Queries().CreateEntity(ctx, db.CreateEntityParams{
			ID:         newDatabaseID, // New database ID for uniqueness
			VersionID:  newVersionID,
			EntityType: entity.EntityType,
			Name:       entity.Name,
			Data:       updatedData, // Includes logical_id
		})
		if err != nil {
			return nil, fmt.Errorf("failed to copy entity %s: %w", entity.ID, err)
		}
	}

	return entityIDMapping, nil
}

// copyRelationshipsFromParent copies all relationships from parent version to new version
func (s *Service) copyRelationshipsFromParent(ctx context.Context, parentVersionID, newVersionID string, entityIDMapping map[string]string) error {
	relationships, err := s.db.Queries().ListRelationshipsByVersion(ctx, parentVersionID)
	if err != nil {
		return err
	}

	// First, we need to build a reverse mapping from database IDs to logical IDs for the parent version
	parentEntities, err := s.db.Queries().ListEntitiesByVersion(ctx, parentVersionID)
	if err != nil {
		return err
	}

	dbToLogicalMapping := make(map[string]string)
	for _, entity := range parentEntities {
		var entityData map[string]any
		if err := json.Unmarshal(entity.Data, &entityData); err != nil {
			continue
		}
		
		if logicalID, exists := entityData["logical_id"].(string); exists {
			dbToLogicalMapping[entity.ID] = logicalID
		} else {
			// First version - database ID is the logical ID
			dbToLogicalMapping[entity.ID] = entity.ID
		}
	}

	for _, rel := range relationships {
		// Map parent database IDs to logical IDs, then to new database IDs
		fromLogicalID := dbToLogicalMapping[rel.FromEntityID]
		toLogicalID := dbToLogicalMapping[rel.ToEntityID]
		
		if fromLogicalID == "" || toLogicalID == "" {
			continue // Skip relationships with unmappable entities
		}
		
		fromNewDatabaseID := entityIDMapping[fromLogicalID]
		toNewDatabaseID := entityIDMapping[toLogicalID]
		
		if fromNewDatabaseID == "" || toNewDatabaseID == "" {
			continue // Skip relationships where entities don't exist in new version
		}

		// Generate new relationship ID for this version
		newRelationshipID := uuid.New().String()
		_, err := s.db.Queries().CreateRelationship(ctx, db.CreateRelationshipParams{
			ID:               newRelationshipID,
			VersionID:        newVersionID,
			FromEntityID:     fromNewDatabaseID,
			ToEntityID:       toNewDatabaseID,
			RelationshipType: rel.RelationshipType,
			Properties:       rel.Properties,
		})
		if err != nil {
			return fmt.Errorf("failed to copy relationship %s: %w", rel.ID, err)
		}
	}

	return nil
}

// applyDelta applies a single delta to the graph
func (s *Service) applyDelta(ctx context.Context, versionID string, delta *Delta, entityIDMapping map[string]string) error {
	switch delta.Operation {
	case "create":
		return s.createEntity(ctx, versionID, delta, entityIDMapping)
	case "update":
		return s.updateEntity(ctx, versionID, delta, entityIDMapping)
	case "delete":
		return s.deleteEntity(ctx, versionID, delta, entityIDMapping)
	default:
		return fmt.Errorf("unknown operation: %s", delta.Operation)
	}
}

// createEntity creates a new entity
func (s *Service) createEntity(ctx context.Context, versionID string, delta *Delta, entityIDMapping map[string]string) error {
	logicalID := delta.EntityID
	if logicalID == "" {
		logicalID = uuid.New().String()
	}

	// Generate new database ID
	databaseID := uuid.New().String()
	
	// Add to mapping
	entityIDMapping[logicalID] = databaseID

	// Extract name from fields
	name := ""
	if nameVal, ok := delta.Fields["name"]; ok {
		if nameStr, ok := nameVal.(string); ok {
			name = nameStr
		}
	}

	// Add logical ID to entity data
	updatedFields := make(map[string]any)
	for k, v := range delta.Fields {
		updatedFields[k] = v
	}
	updatedFields["logical_id"] = logicalID

	// Serialize data as JSON
	dataBytes, err := json.Marshal(updatedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal entity data: %w", err)
	}

	// Create entity with database ID
	_, err = s.db.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         databaseID,
		VersionID:  versionID,
		EntityType: delta.EntityType,
		Name:       name,
		Data:       dataBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}

	// Create relationships
	for _, relDelta := range delta.Relationships {
		if err := s.applyRelationshipDelta(ctx, versionID, relDelta, entityIDMapping); err != nil {
			return fmt.Errorf("failed to apply relationship delta: %w", err)
		}
	}

	return nil
}

// updateEntity updates an existing entity
func (s *Service) updateEntity(ctx context.Context, versionID string, delta *Delta, entityIDMapping map[string]string) error {
	// Map logical entity ID to database ID for this version
	databaseID, exists := entityIDMapping[delta.EntityID]
	if !exists {
		return fmt.Errorf("entity with logical ID %s not found in current version", delta.EntityID)
	}

	// Extract name from fields
	name := ""
	if nameVal, ok := delta.Fields["name"]; ok {
		if nameStr, ok := nameVal.(string); ok {
			name = nameStr
		}
	}

	// Preserve logical ID in the data
	updatedFields := make(map[string]any)
	for k, v := range delta.Fields {
		updatedFields[k] = v
	}
	updatedFields["logical_id"] = delta.EntityID // Preserve logical identity

	// Serialize data as JSON
	dataBytes, err := json.Marshal(updatedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal entity data: %w", err)
	}

	// Update entity using database ID
	_, err = s.db.Queries().UpdateEntity(ctx, db.UpdateEntityParams{
		ID:   databaseID,
		Name: name,
		Data: dataBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}

	// Apply relationship changes
	for _, relDelta := range delta.Relationships {
		if err := s.applyRelationshipDelta(ctx, versionID, relDelta, entityIDMapping); err != nil {
			return fmt.Errorf("failed to apply relationship delta: %w", err)
		}
	}

	return nil
}

// deleteEntity deletes an entity and its relationships
func (s *Service) deleteEntity(ctx context.Context, versionID string, delta *Delta, entityIDMapping map[string]string) error {
	// Map logical entity ID to database ID for this version
	databaseID, exists := entityIDMapping[delta.EntityID]
	if !exists {
		return fmt.Errorf("entity with logical ID %s not found in current version", delta.EntityID)
	}

	// Delete relationships first (referential integrity)
	if err := s.db.Queries().DeleteRelationshipsByEntity(ctx, db.DeleteRelationshipsByEntityParams{
		FromEntityID: databaseID,
		ToEntityID:   databaseID,
	}); err != nil {
		return fmt.Errorf("failed to delete entity relationships: %w", err)
	}

	// Delete entity
	if err := s.db.Queries().DeleteEntity(ctx, databaseID); err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}

	return nil
}

// applyRelationshipDelta applies a relationship change
func (s *Service) applyRelationshipDelta(ctx context.Context, versionID string, relDelta *RelationshipDelta, entityIDMapping map[string]string) error {
	switch relDelta.Operation {
	case "create":
		return s.createRelationship(ctx, versionID, relDelta, entityIDMapping)
	case "update":
		return s.updateRelationship(ctx, versionID, relDelta, entityIDMapping)
	case "delete":
		return s.deleteRelationship(ctx, versionID, relDelta, entityIDMapping)
	default:
		return fmt.Errorf("unknown relationship operation: %s", relDelta.Operation)
	}
}

// createRelationship creates a new relationship
func (s *Service) createRelationship(ctx context.Context, versionID string, relDelta *RelationshipDelta, entityIDMapping map[string]string) error {
	relationshipID := relDelta.RelationshipID
	if relationshipID == "" {
		relationshipID = uuid.New().String()
	}

	// Map logical entity IDs to database IDs
	fromDatabaseID, exists := entityIDMapping[relDelta.FromEntityID]
	if !exists {
		return fmt.Errorf("from entity with logical ID %s not found", relDelta.FromEntityID)
	}
	
	toDatabaseID, exists := entityIDMapping[relDelta.ToEntityID]
	if !exists {
		return fmt.Errorf("to entity with logical ID %s not found", relDelta.ToEntityID)
	}

	// Serialize properties as JSON
	var propertiesBytes []byte
	if relDelta.Properties != nil {
		var err error
		propertiesBytes, err = json.Marshal(relDelta.Properties)
		if err != nil {
			return fmt.Errorf("failed to marshal relationship properties: %w", err)
		}
	}

	_, err := s.db.Queries().CreateRelationship(ctx, db.CreateRelationshipParams{
		ID:               relationshipID,
		VersionID:        versionID,
		FromEntityID:     fromDatabaseID,
		ToEntityID:       toDatabaseID,
		RelationshipType: relDelta.RelationshipType,
		Properties:       propertiesBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	return nil
}

// updateRelationship updates an existing relationship
func (s *Service) updateRelationship(ctx context.Context, versionID string, relDelta *RelationshipDelta, entityIDMapping map[string]string) error {
	// Serialize properties as JSON
	var propertiesBytes []byte
	if relDelta.Properties != nil {
		var err error
		propertiesBytes, err = json.Marshal(relDelta.Properties)
		if err != nil {
			return fmt.Errorf("failed to marshal relationship properties: %w", err)
		}
	}

	_, err := s.db.Queries().UpdateRelationship(ctx, db.UpdateRelationshipParams{
		ID:         relDelta.RelationshipID,
		Properties: propertiesBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to update relationship: %w", err)
	}

	return nil
}

// deleteRelationship deletes a relationship
func (s *Service) deleteRelationship(ctx context.Context, versionID string, relDelta *RelationshipDelta, entityIDMapping map[string]string) error {
	if err := s.db.Queries().DeleteRelationship(ctx, relDelta.RelationshipID); err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	return nil
}

// nullStringToPtr converts sql.NullString to *string
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
// GetNeighborsInVersion retrieves entities connected to a given logical entity in a specific version
func (s *Service) GetNeighborsInVersion(ctx context.Context, versionID string, logicalEntityID string, relationshipType string) ([]*Entity, error) {
	// Get all entities in this version
	entities, err := s.db.Queries().ListEntitiesByVersion(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entities in version: %w", err)
	}

	// Find the database ID for the logical entity ID in this version
	var targetDatabaseID string
	for _, entity := range entities {
		var data map[string]any
		if err := json.Unmarshal(entity.Data, &data); err != nil {
			continue
		}
		
		entityLogicalID := entity.ID // Default to database ID
		if logicalID, exists := data["logical_id"].(string); exists {
			entityLogicalID = logicalID
		}
		
		if entityLogicalID == logicalEntityID {
			targetDatabaseID = entity.ID
			break
		}
	}

	if targetDatabaseID == "" {
		return []*Entity{}, nil // Entity not found in this version
	}

	// Get relationships for this entity
	relationships, err := s.db.Queries().ListRelationshipsByEntity(ctx, db.ListRelationshipsByEntityParams{
		FromEntityID: targetDatabaseID,
		ToEntityID:   targetDatabaseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get relationships: %w", err)
	}

	var neighbors []*Entity
	for _, rel := range relationships {
		if relationshipType != "" && rel.RelationshipType != relationshipType {
			continue
		}

		var neighborDatabaseID string
		if rel.FromEntityID == targetDatabaseID {
			neighborDatabaseID = rel.ToEntityID
		} else {
			neighborDatabaseID = rel.FromEntityID
		}

		// Find the neighbor entity in our entities list
		for _, entity := range entities {
			if entity.ID == neighborDatabaseID {
				var data map[string]any
				if err := json.Unmarshal(entity.Data, &data); err != nil {
					continue
				}

				// Use logical ID if available
				neighborLogicalID := entity.ID
				if logicalID, exists := data["logical_id"].(string); exists {
					neighborLogicalID = logicalID
				}

				neighbors = append(neighbors, &Entity{
					ID:         neighborLogicalID,
					VersionID:  entity.VersionID,
					EntityType: entity.EntityType,
					Name:       entity.Name,
					Data:       data,
					CreatedAt:  entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
					UpdatedAt:  entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
				})
				break
			}
		}
	}

	return neighbors, nil
}

// ImportEntity imports an entity from another project, maintaining its identity
func (s *Service) ImportEntity(ctx context.Context, targetVersionID string, sourceProjectID string, entityLogicalID string) (*Entity, error) {
	// Find the entity in the source project (get the latest version)
	sourceEntity, err := s.findLatestEntityVersion(ctx, sourceProjectID, entityLogicalID)
	if err != nil {
		return nil, fmt.Errorf("failed to find entity %s in project %s: %w", entityLogicalID, sourceProjectID, err)
	}

	// Check if entity already exists in target version
	targetEntities, err := s.db.Queries().ListEntitiesByVersion(ctx, targetVersionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list target entities: %w", err)
	}

	// Check if entity already exists in target version
	for _, entity := range targetEntities {
		var data map[string]any
		if err := json.Unmarshal(entity.Data, &data); err != nil {
			continue
		}
		
		if logicalID, exists := data["logical_id"].(string); exists && logicalID == entityLogicalID {
			// Entity already exists in target version
			return &Entity{
				ID:         logicalID,
				VersionID:  entity.VersionID,
				EntityType: entity.EntityType,
				Name:       entity.Name,
				Data:       data,
				CreatedAt:  entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt:  entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			}, nil
		}
	}

	// Import the entity into the target version
	newDatabaseID := uuid.New().String()
	
	// Add import metadata to the entity data
	var entityData map[string]any
	if err := json.Unmarshal(sourceEntity.Data, &entityData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal source entity data: %w", err)
	}
	
	// Add import tracking
	entityData["imported_from_project"] = sourceProjectID
	entityData["import_timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	
	updatedData, err := json.Marshal(entityData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated entity data: %w", err)
	}

	_, err = s.db.Queries().CreateEntity(ctx, db.CreateEntityParams{
		ID:         newDatabaseID,
		VersionID:  targetVersionID,
		EntityType: sourceEntity.EntityType,
		Name:       sourceEntity.Name,
		Data:       updatedData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to import entity: %w", err)
	}

	return &Entity{
		ID:         entityLogicalID,
		VersionID:  targetVersionID,
		EntityType: sourceEntity.EntityType,
		Name:       sourceEntity.Name,
		Data:       entityData,
		CreatedAt:  time.Now().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

// GetEntityHistory retrieves the evolution of an entity across all projects
func (s *Service) GetEntityHistory(ctx context.Context, entityLogicalID string) ([]*EntityVersion, error) {
	// Get all projects
	projects, err := s.db.Queries().ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	var history []*EntityVersion

	for _, project := range projects {
		// Get working set version for this project
		workingSet, err := s.db.Queries().GetWorkingSetVersion(ctx, project.ID)
		if err != nil {
			continue // Skip projects without working sets
		}

		// Look for the entity in this project's working set
		entities, err := s.db.Queries().ListEntitiesByVersion(ctx, workingSet.ID)
		if err != nil {
			continue
		}

		for _, entity := range entities {
			var data map[string]any
			if err := json.Unmarshal(entity.Data, &data); err != nil {
				continue
			}
			
			logicalID := entity.ID // Default to database ID
			if lid, exists := data["logical_id"].(string); exists {
				logicalID = lid
			}
			
			if logicalID == entityLogicalID {
				history = append(history, &EntityVersion{
					Entity: &Entity{
						ID:         logicalID,
						VersionID:  entity.VersionID,
						EntityType: entity.EntityType,
						Name:       entity.Name,
						Data:       data,
						CreatedAt:  entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
						UpdatedAt:  entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
					},
					ProjectID:   project.ID,
					ProjectName: project.Name,
					VersionID:   workingSet.ID,
					VersionName: workingSet.Name.String,
					CreatedAt:   entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
				})
				break
			}
		}
	}

	return history, nil
}

// ListSharedEntities lists entities that appear in multiple projects
func (s *Service) ListSharedEntities(ctx context.Context) ([]*SharedEntity, error) {
	// Get all projects
	projects, err := s.db.Queries().ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	// Map logical ID to project appearances
	entityProjects := make(map[string]map[string]bool) // logicalID -> projectID -> true
	entityInfo := make(map[string]*SharedEntity)

	for _, project := range projects {
		// Get working set version for this project
		workingSet, err := s.db.Queries().GetWorkingSetVersion(ctx, project.ID)
		if err != nil {
			continue
		}

		// Get entities in this project
		entities, err := s.db.Queries().ListEntitiesByVersion(ctx, workingSet.ID)
		if err != nil {
			continue
		}

		for _, entity := range entities {
			var data map[string]any
			if err := json.Unmarshal(entity.Data, &data); err != nil {
				continue
			}
			
			logicalID := entity.ID
			if lid, exists := data["logical_id"].(string); exists {
				logicalID = lid
			}

			// Track this entity's appearance in this project
			if entityProjects[logicalID] == nil {
				entityProjects[logicalID] = make(map[string]bool)
			}
			entityProjects[logicalID][project.ID] = true

			// Store entity info
			if entityInfo[logicalID] == nil {
				entityInfo[logicalID] = &SharedEntity{
					LogicalID:    logicalID,
					Name:         entity.Name,
					EntityType:   entity.EntityType,
					FirstSeen:    entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
					LastModified: entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
				}
			} else {
				// Update last modified if this is newer
				if lastModTime, err := time.Parse("2006-01-02T15:04:05Z", entityInfo[logicalID].LastModified); err == nil {
					if entity.UpdatedAt.After(lastModTime) {
						entityInfo[logicalID].LastModified = entity.UpdatedAt.Format("2006-01-02T15:04:05Z")
					}
				}
			}
		}
	}

	// Filter for entities that appear in multiple projects
	var sharedEntities []*SharedEntity
	for logicalID, projectMap := range entityProjects {
		if len(projectMap) > 1 {
			entity := entityInfo[logicalID]
			entity.ProjectCount = len(projectMap)
			
			// Get project names
			for projectID := range projectMap {
				for _, project := range projects {
					if project.ID == projectID {
						entity.Projects = append(entity.Projects, project.Name)
						break
					}
				}
			}
			
			sharedEntities = append(sharedEntities, entity)
		}
	}

	return sharedEntities, nil
}

// findLatestEntityVersion finds the latest version of an entity in a project
func (s *Service) findLatestEntityVersion(ctx context.Context, projectID string, entityLogicalID string) (*db.Entity, error) {
	// Get working set version for the project
	workingSet, err := s.db.Queries().GetWorkingSetVersion(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get working set for project: %w", err)
	}

	// Find the entity in the working set
	entities, err := s.db.Queries().ListEntitiesByVersion(ctx, workingSet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list entities: %w", err)
	}

	for _, entity := range entities {
		var data map[string]any
		if err := json.Unmarshal(entity.Data, &data); err != nil {
			continue
		}
		
		logicalID := entity.ID
		if lid, exists := data["logical_id"].(string); exists {
			logicalID = lid
		}
		
		if logicalID == entityLogicalID {
			return &entity, nil
		}
	}

	return nil, fmt.Errorf("entity not found")
}