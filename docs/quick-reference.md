# Libretto Quick Reference

## Essential Commands

### Setup & Development
```bash
make dev-setup          # Complete development environment setup
make test              # Run all tests (unit + integration)
make dashboard         # Launch web dashboard (http://localhost:8080)
```

### Database Operations
```bash
make db-seed           # Create fantasy story database
make db-clean          # Clean and reseed database
make db-inspect-projects    # List all projects
make db-inspect-graph      # Visualize narrative graph
make db-inspect-stats      # Show database statistics
```

### Testing
```bash
make test-unit         # Unit tests only
make test-integration  # Integration tests only
make test-coverage     # Generate coverage report
```

## CLI Tools

### Database Inspector (`dbinspect`)
```bash
# Basic usage
go run cmd/dbinspect/main.go -db <database> -cmd <command>

# Commands
-cmd schema           # Show database schema
-cmd projects         # List projects
-cmd entities         # List entities (requires -project or -version)
-cmd relationships    # List relationships (requires -version or -entity)
-cmd annotations      # List annotations (requires -entity)
-cmd graph           # Show narrative graph (requires -project or -version)
-cmd stats           # Show statistics (requires -project or -version)

# Options
-v                   # Verbose output
-project <id>        # Filter by project ID
-version <id>        # Filter by version ID
-entity <id>         # Filter by entity ID
```

### Database Seeder (`dbseed`)
```bash
# Basic usage
go run cmd/dbseed/main.go -db <database> -preset <preset> [options]

# Presets
-preset fantasy      # Epic fantasy story
-preset scifi        # Science fiction (planned)
-preset mystery      # Mystery/thriller (planned)

# Options
-clean              # Clean database before seeding
```

### Web Dashboard (`dashboard`)
```bash
# Basic usage
go run cmd/dashboard/main.go -db <database> [options]

# Options
-port <port>        # Port to serve on (default: 8080)
```

### Integration Tests (`integration-test`)
```bash
# Basic usage
go run cmd/integration-test/main.go [options]

# Options
-db <path>          # Database path (default: :memory:)
-output <file>      # JSON output file
-v                  # Verbose output
```

## Database Schema

### Core Tables
- `projects` - Top-level narrative containers
- `graph_versions` - Versioned snapshots with working set management
- `entities` - Core narrative elements
- `relationships` - Typed connections between entities
- `annotations` - AI agent analysis and metadata

### Entity Types
- `Scene` - Narrative scenes with content and pacing
- `Character` - Characters with personality and arcs
- `Location` - Settings with atmosphere
- `Theme` - Thematic elements
- `PlotPoint` - Key story moments
- `Arc` - Character and plot arcs

### Relationship Types
- `features` - Scene features Character
- `occurs_at` - Scene occurs at Location
- `influences` - Entity influences Theme
- `conflicts` - Character conflicts with Character
- `supports` - Entity supports Entity
- `precedes`/`follows` - Sequential relationships

## Environment Variables

```bash
# Make targets
DB_FILE=custom.db make db-seed     # Use custom database file
DB_PRESET=fantasy make db-seed     # Use specific preset
DASHBOARD_PORT=3000 make dashboard # Use custom port

# Go applications
LOG_LEVEL=DEBUG                    # Enable debug logging
```

## File Locations

### Databases
- `libretto-dev.db` - Default development database
- `test-*.db` - Test databases
- `:memory:` - In-memory database for tests

### Generated Files
- `coverage.out` - Test coverage data
- `coverage.html` - Coverage report
- `test-results.json` - Integration test results
- `bin/` - Built CLI tools

### Configuration
- `sqlc.yaml` - Database code generation
- `Makefile` - Build and development targets
- `.kiro/` - Kiro IDE configuration

## Common Workflows

### New Feature Development
```bash
# 1. Setup clean environment
make dev-setup

# 2. Run tests to ensure baseline
make test

# 3. Make code changes
# ... edit files ...

# 4. Test changes
make test-unit
make test-integration

# 5. Verify visually
make dashboard
# Visit http://localhost:8080

# 6. Check database state
make db-inspect-graph
```

### Debugging Database Issues
```bash
# 1. Inspect schema
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd schema

# 2. Check data integrity
make test-integration

# 3. View specific entities
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd entities -project <id> -v

# 4. Examine relationships
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd relationships -version <id> -v

# 5. Reset if needed
make db-clean
```

### Performance Analysis
```bash
# 1. Run performance tests
make test-integration

# 2. Check test results
cat test-results.json | jq '.results[] | select(.name == "Performance Benchmarks")'

# 3. Generate coverage report
make test-coverage
open coverage.html

# 4. Monitor with dashboard
make dashboard
```

### CI/CD Integration
```bash
# Continuous Integration
make ci-test          # Run all tests for CI
make ci-build         # Build all components

# Test output for CI systems
go run cmd/integration-test/main.go -output ci-results.json
```

## Troubleshooting

### Common Issues

**Database locked errors:**
```bash
# Check for running processes
lsof libretto-dev.db
pkill -f dashboard
```

**Missing database:**
```bash
make db-seed
```

**Test failures:**
```bash
# Clean environment
make db-clean
make test
```

**Port already in use:**
```bash
DASHBOARD_PORT=3000 make dashboard
```

### Debug Mode
```bash
# Enable verbose logging
LOG_LEVEL=DEBUG make test-integration
LOG_LEVEL=DEBUG make dashboard
```

### Performance Issues
```bash
# Check database size
ls -lh *.db

# Vacuum database
sqlite3 libretto-dev.db "VACUUM;"

# Monitor memory usage
go test ./internal/db -memprofile=mem.prof
go tool pprof mem.prof
```

## Integration with IDEs

### VS Code
- Install Go extension
- Use integrated terminal for make commands
- Configure test runner for `go test ./internal/...`

### Kiro IDE
- Use `#File` and `#Folder` context
- Steering files in `.kiro/steering/`
- MCP configuration in `.kiro/settings/mcp.json`

## API Reference

### Database Queries (sqlc generated)
```go
// Projects
queries.CreateProject(ctx, params)
queries.GetProject(ctx, id)
queries.ListProjects(ctx)
queries.UpdateProject(ctx, params)
queries.DeleteProject(ctx, id)

// Entities
queries.CreateEntity(ctx, params)
queries.GetEntity(ctx, id)
queries.ListEntitiesByVersion(ctx, versionID)
queries.ListEntitiesByType(ctx, params)
queries.UpdateEntity(ctx, params)
queries.DeleteEntity(ctx, id)

// Relationships
queries.CreateRelationship(ctx, params)
queries.ListRelationshipsByVersion(ctx, versionID)
queries.ListRelationshipsByEntity(ctx, params)
queries.GetRelationshipsBetweenEntities(ctx, params)
```

### Data Types
```go
// Entity data structures
types.SceneData
types.CharacterData
types.LocationData
types.ThemeData
types.PlotPointData
types.ArcData

// Annotation data structures
types.EmotionalAnalysisData
types.ThematicScoreData
types.ContinuityCheckData

// Helper functions
types.MarshalEntityData(data)
types.UnmarshalSceneData(raw)
// ... etc for each type
```

### Monitoring
```go
// Structured logging
logger := monitoring.NewLogger("component")
logger.Info(ctx, "message", monitoring.String("key", "value"))

// Operation timing
timer := logger.StartOperation(ctx, "operation")
timer.Complete(ctx, "success message")

// Database metrics
dbMetrics := monitoring.NewDatabaseMetrics(logger)
dbMetrics.RecordQuery(ctx, "operation", duration, err)
```

This quick reference covers the most common operations and workflows for developing with the Libretto narrative engine.