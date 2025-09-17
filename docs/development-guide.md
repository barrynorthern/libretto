# Libretto Development Guide

## Getting Started

### Prerequisites

Ensure you have the following installed:

- **Go 1.22.5+**: `go version`
- **SQLite3**: `sqlite3 --version`
- **Make**: `make --version`
- **Git**: `git --version`

Optional but recommended:
- **entr**: For watch mode testing (`brew install entr` on macOS)
- **jq**: For JSON processing (`brew install jq` on macOS)

### Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd libretto

# Complete development setup
make dev-setup
```

This command will:
1. Clean any existing database
2. Seed database with fantasy story
3. Build all CLI tools
4. Verify setup with tests

### Verify Installation

```bash
# Run all tests
make test

# Launch dashboard
make dashboard
# Visit http://localhost:8080

# Inspect database
make db-inspect-projects
```

## Development Workflow

### Daily Development Cycle

```bash
# 1. Start with clean environment
make db-clean

# 2. Run tests to ensure baseline
make test

# 3. Start dashboard for visual feedback
make dashboard &

# 4. Make code changes
# ... edit files ...

# 5. Test changes continuously
make test-unit

# 6. Run integration tests
make test-integration

# 7. Verify visually in dashboard
# Visit http://localhost:8080
```

### Code Organization

#### Adding New Entity Types

1. **Define Type Constants** (`internal/types/entities.go`):
```go
const (
    EntityTypeNewType EntityType = "NewType"
)
```

2. **Create Data Structure**:
```go
type NewTypeData struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    // ... other fields
}
```

3. **Add Marshal/Unmarshal Functions**:
```go
func UnmarshalNewTypeData(raw json.RawMessage) (*NewTypeData, error) {
    var data NewTypeData
    err := json.Unmarshal(raw, &data)
    return &data, err
}
```

4. **Add Unit Tests** (`internal/types/entities_test.go`):
```go
func TestNewTypeDataMarshalUnmarshal(t *testing.T) {
    original := &NewTypeData{
        Name:        "Test",
        Description: "Test description",
    }
    
    data, err := types.MarshalEntityData(original)
    if err != nil {
        t.Fatalf("Failed to marshal: %v", err)
    }
    
    unmarshaled, err := types.UnmarshalNewTypeData(data)
    if err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }
    
    if unmarshaled.Name != original.Name {
        t.Errorf("Expected name %s, got %s", original.Name, unmarshaled.Name)
    }
}
```

5. **Update Database Seeding** (`cmd/dbseed/main.go`):
```go
// Add to createFantasyEntities function
newTypeData := &types.NewTypeData{
    Name:        "Example New Type",
    Description: "Example description",
}

data, _ := types.MarshalEntityData(newTypeData)
entities = append(entities, db.CreateEntityParams{
    ID:         uuid.New().String(),
    VersionID:  versionID,
    EntityType: string(types.EntityTypeNewType),
    Name:       newTypeData.Name,
    Data:       data,
})
```

6. **Update Dashboard Visualization** (`cmd/dashboard/main.go`):
```go
// Add to color mapping
const colors = {
    'NewType': '#color-hex',
    // ... existing colors
};

// Add to type groups
typeGroups := map[string]int{
    "NewType": 7,
    // ... existing types
}
```

#### Adding New Relationship Types

1. **Define Relationship Type** (`internal/types/entities.go`):
```go
const (
    RelationshipNewType RelationshipType = "new_type"
)
```

2. **Update Database Seeding**:
```go
relationships = append(relationships, db.CreateRelationshipParams{
    ID:               uuid.New().String(),
    VersionID:        versionID,
    FromEntityID:     fromID,
    ToEntityID:       toID,
    RelationshipType: string(types.RelationshipNewType),
    Properties:       json.RawMessage(`{"property": "value"}`),
})
```

3. **Add Tests**:
```go
func TestNewRelationshipType(t *testing.T) {
    // Test relationship creation and queries
}
```

#### Adding New Annotation Types

1. **Define Annotation Type** (`internal/types/entities.go`):
```go
const (
    AnnotationNewAnalysis AnnotationType = "new_analysis"
)

type NewAnalysisData struct {
    Score      float64   `json:"score"`
    Analysis   string    `json:"analysis"`
    AnalyzedAt time.Time `json:"analyzed_at"`
}
```

2. **Add Marshal/Unmarshal Functions**:
```go
func UnmarshalNewAnalysisData(raw json.RawMessage) (*NewAnalysisData, error) {
    var data NewAnalysisData
    err := json.Unmarshal(raw, &data)
    return &data, err
}
```

3. **Update Database Seeding**:
```go
newAnalysisData := &types.NewAnalysisData{
    Score:      0.85,
    Analysis:   "Analysis result",
    AnalyzedAt: time.Now(),
}

data, _ := json.Marshal(newAnalysisData)
annotations = append(annotations, db.CreateAnnotationParams{
    ID:             uuid.New().String(),
    EntityID:       entityID,
    AnnotationType: string(types.AnnotationNewAnalysis),
    Content:        "New analysis annotation",
    Metadata:       data,
    AgentName:      sql.NullString{String: "new_agent", Valid: true},
})
```

### Database Development

#### Schema Changes

1. **Create Migration** (`internal/db/migrations/`):
```sql
-- 003_add_new_feature.sql
ALTER TABLE entities ADD COLUMN new_field TEXT;
CREATE INDEX idx_entities_new_field ON entities(new_field);
```

2. **Update Queries** (`internal/db/queries/`):
```sql
-- name: GetEntitiesByNewField :many
SELECT * FROM entities
WHERE new_field = ?
ORDER BY created_at DESC;
```

3. **Regenerate Code**:
```bash
make sqlc
```

4. **Update Tests**:
```go
func TestNewFieldQuery(t *testing.T) {
    // Test new query functionality
}
```

#### Query Optimization

1. **Add Indexes** for frequently queried columns:
```sql
CREATE INDEX idx_table_column ON table(column);
CREATE INDEX idx_table_composite ON table(col1, col2);
```

2. **Analyze Query Performance**:
```bash
# Enable query logging in SQLite
sqlite3 libretto-dev.db "PRAGMA query_only = ON;"

# Use EXPLAIN QUERY PLAN
sqlite3 libretto-dev.db "EXPLAIN QUERY PLAN SELECT * FROM entities WHERE entity_type = 'Scene';"
```

3. **Monitor with Integration Tests**:
```bash
make test-integration
# Check performance benchmarks in output
```

### Testing Strategy

#### Unit Tests

Focus on individual components:

```go
func TestEntityCreation(t *testing.T) {
    // Test single entity creation
    queries := setupTestDB(t)
    ctx := context.Background()
    
    entity, err := queries.CreateEntity(ctx, params)
    if err != nil {
        t.Fatalf("Failed to create entity: %v", err)
    }
    
    // Verify entity properties
    if entity.Name != expectedName {
        t.Errorf("Expected name %s, got %s", expectedName, entity.Name)
    }
}
```

#### Integration Tests

Test complete workflows:

```go
func TestCompleteNarrativeWorkflow(t *testing.T) {
    // Test project -> version -> entities -> relationships -> annotations
    queries := setupTestDB(t)
    ctx := context.Background()
    
    // Create project
    project, err := queries.CreateProject(ctx, projectParams)
    // ... continue with complete workflow
    
    // Verify end-to-end functionality
    entities, err := queries.ListEntitiesByVersion(ctx, versionID)
    if len(entities) != expectedCount {
        t.Errorf("Expected %d entities, got %d", expectedCount, len(entities))
    }
}
```

#### Performance Tests

Establish benchmarks:

```go
func BenchmarkEntityCreation(b *testing.B) {
    queries := setupTestDB(b)
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        entityID := uuid.New().String()
        params := CreateEntityParams{
            ID:         entityID,
            VersionID:  versionID,
            EntityType: "Scene",
            Name:       fmt.Sprintf("Scene %d", i),
            Data:       sceneData,
        }
        
        _, err := queries.CreateEntity(ctx, params)
        if err != nil {
            b.Fatalf("Failed to create entity: %v", err)
        }
    }
}
```

### Monitoring and Debugging

#### Structured Logging

Use consistent logging patterns:

```go
logger := monitoring.NewLogger("component-name")

// Context-aware logging
logger.Info(ctx, "Operation started", 
    monitoring.String("operation", "create_entity"),
    monitoring.String("entity_type", "Scene"),
    monitoring.String("entity_id", entityID))

// Error logging with context
logger.Error(ctx, "Operation failed", err,
    monitoring.String("operation", "create_entity"),
    monitoring.Duration("duration", time.Since(start)))

// Operation timing
timer := logger.StartOperation(ctx, "complex_operation")
// ... do work ...
timer.Complete(ctx, "Operation completed successfully")
```

#### Database Metrics

Track database performance:

```go
dbMetrics := monitoring.NewDatabaseMetrics(logger)

// Record query performance
start := time.Now()
result, err := queries.CreateEntity(ctx, params)
duration := time.Since(start)

dbMetrics.RecordQuery(ctx, "CreateEntity", duration, err)
dbMetrics.RecordEntityOperation(ctx, "Scene", "create", duration, err)
```

#### Visual Debugging

Use the dashboard for visual inspection:

```bash
# Start dashboard
make dashboard

# Navigate to:
# - Project overview: http://localhost:8080
# - Project details: http://localhost:8080/project/<project-id>
# - Graph visualization: http://localhost:8080/graph/<project-id>
```

### Code Quality

#### Linting and Formatting

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

#### Test Coverage

```bash
# Generate coverage report
make test-coverage

# View coverage
open coverage.html

# Check coverage percentage
go tool cover -func=coverage.out | grep total
```

#### Documentation

1. **Code Comments**: Document public functions and complex logic
2. **README Updates**: Keep README.md current with new features
3. **API Documentation**: Document database queries and types
4. **Architecture Decisions**: Record significant decisions in `docs/adr/`

### Debugging Common Issues

#### Database Issues

```bash
# Check database integrity
sqlite3 libretto-dev.db "PRAGMA integrity_check;"

# View schema
make db-inspect-schema

# Check foreign key constraints
sqlite3 libretto-dev.db "PRAGMA foreign_key_check;"

# Reset database
make db-clean
```

#### Test Failures

```bash
# Run specific test
go test ./internal/db -run TestSpecificFunction -v

# Run with race detection
go test ./internal/db -race -v

# Debug with verbose output
go test ./internal/db -v -count=1
```

#### Performance Issues

```bash
# Profile memory usage
go test ./internal/db -memprofile=mem.prof
go tool pprof mem.prof

# Profile CPU usage
go test ./internal/db -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Check database size
ls -lh *.db

# Vacuum database
sqlite3 libretto-dev.db "VACUUM;"
```

### Continuous Integration

#### Local CI Simulation

```bash
# Run CI test suite
make ci-test

# Build all components
make ci-build

# Generate test reports
go run cmd/integration-test/main.go -output ci-results.json
```

#### GitHub Actions Integration

```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.22'
      
      - name: Run tests
        run: make ci-test
      
      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: test-results
          path: test-results.json
```

### Release Process

#### Version Management

1. **Update Version**: Update version constants in code
2. **Run Full Test Suite**: `make test`
3. **Generate Documentation**: Update README and docs
4. **Build Release**: `make ci-build`
5. **Tag Release**: `git tag v1.0.0`

#### Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] Performance benchmarks acceptable
- [ ] Database migrations tested
- [ ] CLI tools working
- [ ] Dashboard functional
- [ ] Integration tests passing

This development guide provides a comprehensive foundation for contributing to the Libretto narrative engine. Follow these patterns and practices to maintain code quality and system reliability.