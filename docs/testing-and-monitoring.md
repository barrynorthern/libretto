# Testing and Monitoring Guide

This guide covers the comprehensive testing, visualization, and monitoring tools available for the Libretto narrative engine.

## Overview

The Libretto system includes multiple layers of testing and observability:

1. **Unit Tests** - Comprehensive test coverage for all database operations and data models
2. **Integration Tests** - End-to-end testing of complex workflows
3. **Database Inspection** - CLI tools for examining database state
4. **Web Dashboard** - Visual interface for monitoring narrative graphs
5. **Logging & Metrics** - Structured logging and performance monitoring
6. **Performance Testing** - Benchmarks and load testing

## Quick Start

### 1. Run Unit Tests

```bash
# Run all unit tests
go test ./internal/... -v

# Run specific package tests
go test ./internal/db -v
go test ./internal/types -v

# Run tests with coverage
go test ./internal/... -cover
```

### 2. Set Up Test Database

```bash
# Create and seed a test database
go run cmd/dbseed/main.go -db test.db -preset fantasy -clean

# Inspect the database
go run cmd/dbinspect/main.go -db test.db -cmd projects
go run cmd/dbinspect/main.go -db test.db -cmd graph -project <project-id>
```

### 3. Launch Web Dashboard

```bash
# Start the web dashboard
go run cmd/dashboard/main.go -db test.db -port 8080

# Open browser to http://localhost:8080
```

### 4. Run Integration Tests

```bash
# Run comprehensive integration tests
go run cmd/integration-test/main.go -v -output test-results.json

# Run with persistent database
go run cmd/integration-test/main.go -db integration-test.db -v
```

## Testing Tools

### Unit Tests

Located in `internal/*/` directories with `*_test.go` files:

- **Database Tests**: `internal/db/*_test.go`
  - CRUD operations for all entities
  - Constraint validation
  - Foreign key relationships
  - Cascade deletes
  - Unique constraints

- **Type Tests**: `internal/types/entities_test.go`
  - Data model marshaling/unmarshaling
  - Type constant validation
  - JSON serialization integrity

- **Integration Test**: `internal/db/integration_test.go`
  - Complete narrative graph creation
  - Multi-entity workflows
  - Cross-table relationships

### Integration Test Suite

The `cmd/integration-test` tool provides comprehensive end-to-end testing:

**Test Categories:**
- Database Schema Validation
- Project CRUD Operations
- Graph Version Management
- Entity Management
- Relationship Management
- Annotation System
- Data Model Validation
- Complex Narrative Graph
- Performance Benchmarks
- Data Integrity
- Concurrent Operations

**Usage:**
```bash
# Basic run
go run cmd/integration-test/main.go

# Verbose output with JSON report
go run cmd/integration-test/main.go -v -output results.json

# Use persistent database
go run cmd/integration-test/main.go -db persistent.db -v
```

**Sample Output:**
```
Running test: Database Schema Validation...
✓ Database Schema Validation (2.34ms)
Running test: Project CRUD Operations...
✓ Project CRUD Operations (5.67ms)
...

=== TEST SUMMARY ===
Total Tests: 11
Passed: 11
Failed: 0
Total Time: 45.23ms
Success Rate: 100.0%
```

## Database Inspection Tools

### dbinspect CLI

The `cmd/dbinspect` tool provides detailed database inspection:

**Commands:**
```bash
# Show database schema
go run cmd/dbinspect/main.go -db test.db -cmd schema

# List all projects
go run cmd/dbinspect/main.go -db test.db -cmd projects -v

# Show entities in a project
go run cmd/dbinspect/main.go -db test.db -cmd entities -project <project-id> -v

# Show relationships
go run cmd/dbinspect/main.go -db test.db -cmd relationships -version <version-id> -v

# Show annotations for an entity
go run cmd/dbinspect/main.go -db test.db -cmd annotations -entity <entity-id> -v

# Visualize narrative graph
go run cmd/dbinspect/main.go -db test.db -cmd graph -project <project-id>

# Show statistics
go run cmd/dbinspect/main.go -db test.db -cmd stats -project <project-id>
```

**Sample Output:**
```
=== PROJECTS ===
ID                                    Name                    Theme           Genre
550e8400-e29b-41d4-a716-446655440000  The Crystal of Light    Good vs Evil    Epic Fantasy

=== ENTITIES ===
--- Scene Entities (3) ---
ID                                    Name                    Data Preview           Created
550e8400-e29b-41d4-a716-446655440001  The Call to Adventure   Act: Act1, Seq: 1     2024-01-15 10:30
550e8400-e29b-41d4-a716-446655440002  The Dark Forest         Act: Act2, Seq: 8     2024-01-15 10:31
550e8400-e29b-41d4-a716-446655440003  The Final Battle        Act: Act3, Seq: 25    2024-01-15 10:32

--- Character Entities (2) ---
ID                                    Name                    Data Preview           Created
550e8400-e29b-41d4-a716-446655440004  Elara the Brave         Role: protagonist      2024-01-15 10:33
550e8400-e29b-41d4-a716-446655440005  Shadow Lord Malachar    Role: antagonist       2024-01-15 10:34
```

### dbseed Tool

The `cmd/dbseed` tool creates realistic test data:

**Presets:**
- `fantasy` - Epic fantasy story with heroes, villains, and magic
- `scifi` - Science fiction story (planned)
- `mystery` - Mystery/thriller story (planned)

**Usage:**
```bash
# Create fresh fantasy database
go run cmd/dbseed/main.go -db fantasy.db -preset fantasy -clean

# Add to existing database
go run cmd/dbseed/main.go -db existing.db -preset fantasy
```

**Generated Content:**
- Complete project with theme and genre
- Multiple scenes across story acts
- Characters with detailed profiles
- Locations with atmospheric descriptions
- Themes with symbolic elements
- Relationships between all entities
- AI agent annotations with realistic analysis

## Web Dashboard

### Features

The `cmd/dashboard` tool provides a web-based interface:

**Home Page:**
- Project overview cards
- Entity/relationship/annotation counts
- Quick navigation to detailed views

**Project Detail Page:**
- Complete entity listings by type
- Relationship visualization
- Version management
- Statistics dashboard

**Graph Visualization:**
- Interactive D3.js force-directed graph
- Color-coded entity types
- Draggable nodes
- Relationship links
- Node information panel

### Usage

```bash
# Start dashboard on default port 8080
go run cmd/dashboard/main.go -db test.db

# Use custom port
go run cmd/dashboard/main.go -db test.db -port 3000

# Open browser to http://localhost:8080
```

### Screenshots

The dashboard provides:
- **Project Cards**: Overview of all projects with key metrics
- **Entity Grids**: Organized display of scenes, characters, locations, etc.
- **Relationship Lists**: Typed connections between entities
- **Interactive Graph**: Visual representation of narrative structure
- **Statistics**: Real-time counts and metrics

## Monitoring and Logging

### Structured Logging

The `internal/monitoring` package provides comprehensive logging:

**Features:**
- Structured JSON logging
- Context-aware log entries
- Operation timing
- Database metrics
- Error tracking

**Usage:**
```go
logger := monitoring.NewLogger("narrative-engine")

// Basic logging
logger.Info(ctx, "Processing scene", 
    monitoring.String("scene_id", sceneID),
    monitoring.String("operation", "create"))

// Operation timing
timer := logger.StartOperation(ctx, "create_narrative_graph")
// ... do work ...
timer.Complete(ctx, "Successfully created graph")

// Error logging
logger.Error(ctx, "Failed to create entity", err,
    monitoring.String("entity_type", "Scene"),
    monitoring.String("entity_id", entityID))
```

**Log Output:**
```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "message": "Processing scene",
  "component": "narrative-engine",
  "scene_id": "550e8400-e29b-41d4-a716-446655440001",
  "operation": "create"
}
```

### Database Metrics

Track database performance and usage:

```go
dbMetrics := monitoring.NewDatabaseMetrics(logger)

// Record query performance
dbMetrics.RecordQuery(ctx, "CreateEntity", duration, err)

// Record entity operations
dbMetrics.RecordEntityOperation(ctx, "Scene", "create", duration, err)

// Log current metrics
dbMetrics.LogDatabaseMetrics(ctx)
```

## Performance Testing

### Benchmarks

The integration test suite includes performance benchmarks:

**Entity Creation Benchmark:**
- Creates 100 entities
- Measures creation time per entity
- Validates performance thresholds
- Reports detailed timing metrics

**Query Performance:**
- Tests entity listing performance
- Measures relationship queries
- Validates response times

**Concurrent Operations:**
- Tests concurrent entity creation
- Validates data consistency
- Measures throughput under load

**Sample Results:**
```json
{
  "name": "Performance Benchmarks",
  "passed": true,
  "duration": "125.45ms",
  "details": {
    "entities_created": 100,
    "creation_duration_ms": 89.23,
    "listing_duration_ms": 12.45,
    "creation_per_entity_ms": 0.89
  }
}
```

### Load Testing

For production load testing:

1. **Use Integration Tests**: Scale up entity counts in performance tests
2. **Database Profiling**: Enable SQLite query logging
3. **Memory Monitoring**: Track Go memory usage during operations
4. **Concurrent Users**: Simulate multiple narrative graph operations

## Continuous Integration

### Test Automation

Recommended CI pipeline:

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
      
      # Unit tests
      - name: Run unit tests
        run: go test ./internal/... -v -cover
      
      # Integration tests
      - name: Run integration tests
        run: go run cmd/integration-test/main.go -output results.json
      
      # Upload results
      - name: Upload test results
        uses: actions/upload-artifact@v3
        with:
          name: test-results
          path: results.json
```

### Quality Gates

Set up quality thresholds:
- **Test Coverage**: Minimum 80% coverage
- **Performance**: Entity creation < 10ms per entity
- **Integration**: All integration tests must pass
- **Data Integrity**: No foreign key violations

## Troubleshooting

### Common Issues

**Database Lock Errors:**
```bash
# Check for concurrent access
lsof test.db

# Use WAL mode for better concurrency
PRAGMA journal_mode=WAL;
```

**Performance Issues:**
```bash
# Analyze query performance
go run cmd/dbinspect/main.go -db test.db -cmd stats

# Check database size
ls -lh test.db

# Vacuum database
sqlite3 test.db "VACUUM;"
```

**Memory Usage:**
```bash
# Monitor during tests
go test ./internal/db -memprofile=mem.prof
go tool pprof mem.prof
```

### Debug Mode

Enable verbose logging:
```bash
export LOG_LEVEL=DEBUG
go run cmd/integration-test/main.go -v
```

### Database Recovery

If database becomes corrupted:
```bash
# Check integrity
sqlite3 test.db "PRAGMA integrity_check;"

# Recreate from seed
go run cmd/dbseed/main.go -db recovered.db -preset fantasy -clean
```

## Best Practices

### Testing Strategy

1. **Unit Tests First**: Test individual components in isolation
2. **Integration Tests**: Validate complete workflows
3. **Performance Tests**: Establish baseline performance
4. **Manual Testing**: Use dashboard for exploratory testing
5. **Continuous Monitoring**: Track metrics in production

### Database Management

1. **Regular Backups**: Copy database files before major changes
2. **Migration Testing**: Test schema changes with existing data
3. **Performance Monitoring**: Track query performance over time
4. **Data Validation**: Use constraints and foreign keys
5. **Cleanup**: Remove test data regularly

### Monitoring

1. **Structured Logging**: Use consistent log formats
2. **Error Tracking**: Monitor error rates and patterns
3. **Performance Metrics**: Track operation timing
4. **Resource Usage**: Monitor memory and disk usage
5. **Alerting**: Set up alerts for critical failures

This comprehensive testing and monitoring setup ensures the Libretto narrative engine is reliable, performant, and observable in all environments.