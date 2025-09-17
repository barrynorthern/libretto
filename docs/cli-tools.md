# CLI Tools Documentation

The Libretto narrative engine includes several command-line tools for database management, testing, and monitoring. This document provides comprehensive usage information for each tool.

## Overview

| Tool | Purpose | Make Target |
|------|---------|-------------|
| `dbinspect` | Database inspection and analysis | `make db-inspect` |
| `dbseed` | Database seeding with test data | `make db-seed` |
| `dashboard` | Web-based monitoring interface | `make dashboard` |
| `integration-test` | Comprehensive test suite | `make test-integration` |

## Database Inspector (`dbinspect`)

Interactive CLI tool for examining database state and structure.

### Usage

```bash
go run cmd/dbinspect/main.go [options]

# Or using Make targets
make db-inspect-projects    # List all projects
make db-inspect-schema      # Show database schema
make db-inspect-stats       # Show statistics
make db-inspect-graph       # Visualize narrative graph
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-db` | Path to SQLite database | `libretto.db` |
| `-cmd` | Command to execute | `schema` |
| `-project` | Project ID for filtering | - |
| `-version` | Version ID for filtering | - |
| `-entity` | Entity ID for filtering | - |
| `-v` | Verbose output | `false` |

### Commands

#### `schema` - Database Schema
Shows complete database schema with table structures, indexes, and row counts.

```bash
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd schema
```

**Output:**
```
=== DATABASE SCHEMA ===

--- Table: projects ---
Column          Type        Not Null    Default             PK
id              TEXT        true        NULL                true
name            TEXT        true        NULL                false
theme           TEXT        false       NULL                false
genre           TEXT        false       NULL                false
description     TEXT        false       ''                  false
created_at      DATETIME    true        CURRENT_TIMESTAMP   false
updated_at      DATETIME    true        CURRENT_TIMESTAMP   false
Rows: 1

--- Table: entities ---
Column          Type        Not Null    Default             PK
id              TEXT        true        NULL                true
version_id      TEXT        true        NULL                false
entity_type     TEXT        true        NULL                false
name            TEXT        true        ''                  false
data            JSON        true        NULL                false
created_at      DATETIME    true        CURRENT_TIMESTAMP   false
updated_at      DATETIME    true        CURRENT_TIMESTAMP   false
Rows: 9
```

#### `projects` - Project Listing
Lists all projects with metadata.

```bash
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd projects -v
```

**Output:**
```
=== PROJECTS ===
ID                                    Name                          Theme         Genre         Description                    Created
cd79795f-55bb-4f54-bdb2-3256151524e4  The Crystal of Eternal Light  Good vs Evil  Epic Fantasy  A tale of heroes seeking...    2024-01-15 10:30
```

#### `entities` - Entity Listing
Lists entities by version or project, organized by type.

```bash
# List entities for a project (uses working set version)
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd entities -project <project-id> -v

# List entities for specific version
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd entities -version <version-id>
```

**Output:**
```
=== ENTITIES ===
Using working set version: db7e7042-c482-4c55-84b7-d4f0a238ba88

--- Scene Entities (3) ---
ID                                    Name                    Data Preview           Created
5054dfb2-0934-4693-8793-95bd3643ac1a  The Call to Adventure   Act: Act1, Seq: 1     2024-01-15 10:30
0a2b2700-ee29-44f0-8020-1beea8e16585  The Dark Forest         Act: Act2, Seq: 8     2024-01-15 10:31
09e0afde-dc6d-4732-a93e-6d35389ecbc4  The Final Battle        Act: Act3, Seq: 25    2024-01-15 10:32

--- Character Entities (2) ---
ID                                    Name                    Data Preview           Created
fb074932-cb88-4d1c-9e5d-832fb9352d36  Elara the Brave         Role: protagonist      2024-01-15 10:33
9656270a-fd76-4c88-afc4-6eddb8fcfb04  Shadow Lord Malachar    Role: antagonist       2024-01-15 10:34
```

#### `relationships` - Relationship Analysis
Shows typed connections between entities.

```bash
# List relationships for a version
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd relationships -version <version-id> -v

# List relationships for specific entity
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd relationships -entity <entity-id>
```

**Output:**
```
=== RELATIONSHIPS ===
Relationships for version: db7e7042-c482-4c55-84b7-d4f0a238ba88

--- features Relationships (6) ---
From Entity                           To Entity                             Properties                     Created
5054dfb2-0934-4693-8793-95bd3643ac1a  fb074932-cb88-4d1c-9e5d-832fb9352d36  {"role": "protagonist"}        2024-01-15 10:35
0a2b2700-ee29-44f0-8020-1beea8e16585  fb074932-cb88-4d1c-9e5d-832fb9352d36  {"role": "protagonist"}        2024-01-15 10:36

--- occurs_at Relationships (2) ---
From Entity                           To Entity                             Properties                     Created
0a2b2700-ee29-44f0-8020-1beea8e16585  5c9e7f6e-6c7e-475a-a6c0-7a38a11a14e8  {}                            2024-01-15 10:37
09e0afde-dc6d-4732-a93e-6d35389ecbc4  f0e66964-8698-4291-820f-7a5f820bb3cf  {}                            2024-01-15 10:38
```

#### `annotations` - Annotation Analysis
Shows AI agent annotations for entities.

```bash
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd annotations -entity <entity-id> -v
```

**Output:**
```
=== ANNOTATIONS ===

--- emotional_analysis Annotations (1) ---
Agent         Content                                   Metadata Preview               Created
empath_agent  Strong opening with good emotional...     {"sentiment": 0.7, "emot...   2024-01-15 10:40

--- thematic_score Annotations (1) ---
Agent              Content                              Metadata Preview               Created
thematic_steward   Perfect thematic culmination...      {"relevance_score": 0.98...   2024-01-15 10:41
```

#### `graph` - Narrative Graph Visualization
Shows the complete narrative graph structure in text format.

```bash
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd graph -project <project-id>
```

**Output:**
```
=== NARRATIVE GRAPH ===
Graph for version: db7e7042-c482-4c55-84b7-d4f0a238ba88
Entities: 9, Relationships: 7

The Call to Adventure (Scene) --features--> Elara the Brave (Character)
The Dark Forest (Scene) --features--> Elara the Brave (Character)
The Dark Forest (Scene) --occurs_at--> Shadowwood Forest (Location)
The Final Battle (Scene) --features--> Elara the Brave (Character)
The Final Battle (Scene) --features--> Shadow Lord Malachar (Character)
The Final Battle (Scene) --occurs_at--> Crystal Caverns (Location)
Elara the Brave (Character) --conflicts--> Shadow Lord Malachar (Character)
```

#### `stats` - Database Statistics
Shows comprehensive statistics about the narrative graph.

```bash
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd stats -project <project-id>
```

**Output:**
```
=== STATISTICS ===

Entity Counts:
Type          Count
Scene         3
Character     2
Location      2
Theme         2
PlotPoint     0
Arc           0
TOTAL         9

Relationship Counts:
Type          Count
features      6
occurs_at     2
conflicts     1
TOTAL         9
```

## Database Seeder (`dbseed`)

Creates realistic test data for development and testing.

### Usage

```bash
go run cmd/dbseed/main.go [options]

# Or using Make targets
make db-seed                    # Seed with default settings
DB_PRESET=fantasy make db-seed  # Specify preset
DB_FILE=custom.db make db-seed  # Custom database file
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-db` | Path to SQLite database | `libretto.db` |
| `-preset` | Data preset to load | `fantasy` |
| `-clean` | Clean database before seeding | `false` |

### Presets

#### `fantasy` - Epic Fantasy Story
Creates a complete epic fantasy narrative with:

**Project:**
- Name: "The Crystal of Eternal Light"
- Theme: "Good vs Evil"
- Genre: "Epic Fantasy"
- Description: "A tale of heroes seeking an ancient crystal to save their realm"

**Entities:**
- **3 Scenes**: The Call to Adventure, The Dark Forest, The Final Battle
- **2 Characters**: Elara the Brave (protagonist), Shadow Lord Malachar (antagonist)
- **2 Locations**: Shadowwood Forest, Crystal Caverns
- **2 Themes**: Good vs Evil, Courage and Sacrifice

**Relationships:**
- Scene-Character relationships (features)
- Scene-Location relationships (occurs_at)
- Character conflicts

**Annotations:**
- Emotional analysis for scenes
- Thematic scoring
- AI agent metadata

#### Future Presets
- `scifi` - Science fiction story (planned)
- `mystery` - Mystery/thriller story (planned)

### Example Usage

```bash
# Create fresh fantasy database
go run cmd/dbseed/main.go -db fantasy-story.db -preset fantasy -clean

# Add fantasy data to existing database
go run cmd/dbseed/main.go -db existing.db -preset fantasy

# Verify seeded data
go run cmd/dbinspect/main.go -db fantasy-story.db -cmd projects
```

## Web Dashboard (`dashboard`)

Web-based interface for monitoring and visualizing narrative graphs.

### Usage

```bash
go run cmd/dashboard/main.go [options]

# Or using Make targets
make dashboard                      # Start with default settings
DASHBOARD_PORT=3000 make dashboard  # Custom port
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-db` | Path to SQLite database | `libretto.db` |
| `-port` | Port to serve on | `8080` |

### Features

#### Home Page (`/`)
- Project overview cards with statistics
- Entity/relationship/annotation counts
- Quick navigation to detailed views

#### Project Details (`/project/<project-id>`)
- Complete entity listings organized by type
- Relationship visualization
- Version management information
- Real-time statistics

#### Graph Visualization (`/graph/<project-id>`)
- Interactive D3.js force-directed graph
- Color-coded entity types:
  - **Red**: Scene
  - **Blue**: Character
  - **Green**: Location
  - **Orange**: Theme
  - **Purple**: PlotPoint
  - **Teal**: Arc
- Draggable nodes
- Relationship links with labels
- Node information panel
- Legend and controls

#### API Endpoints
- `/api/graph/<project-id>` - JSON graph data for visualization

### Example Usage

```bash
# Start dashboard
go run cmd/dashboard/main.go -db libretto-dev.db -port 8080

# Visit in browser
open http://localhost:8080

# View specific project
open http://localhost:8080/project/<project-id>

# Interactive graph
open http://localhost:8080/graph/<project-id>
```

## Integration Test Suite (`integration-test`)

Comprehensive end-to-end testing framework.

### Usage

```bash
go run cmd/integration-test/main.go [options]

# Or using Make targets
make test-integration               # Run with default settings
make test-integration -v           # Verbose output
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-db` | Database path | `:memory:` |
| `-output` | JSON output file | - |
| `-v` | Verbose output | `false` |

### Test Categories

1. **Database Schema Validation** - Verifies all tables and constraints exist
2. **Project CRUD Operations** - Tests create, read, update, delete for projects
3. **Graph Version Management** - Tests versioning and working set switching
4. **Entity Management** - Tests all entity types and operations
5. **Relationship Management** - Tests typed connections between entities
6. **Annotation System** - Tests AI agent annotations and metadata
7. **Data Model Validation** - Tests JSON serialization of complex types
8. **Complex Narrative Graph** - Tests complete story creation workflow
9. **Performance Benchmarks** - Validates response times and throughput
10. **Data Integrity** - Tests foreign key constraints and cascade deletes
11. **Concurrent Operations** - Tests system behavior under load

### Output Format

#### Console Output
```
Running test: Database Schema Validation...
✓ Database Schema Validation (0.13ms)
Running test: Project CRUD Operations...
✓ Project CRUD Operations (0.20ms)
...

=== TEST SUMMARY ===
Total Tests: 11
Passed: 11
Failed: 0
Total Time: 5.45ms
Success Rate: 100.0%
```

#### JSON Output (`-output results.json`)
```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "total_tests": 11,
  "passed_tests": 11,
  "failed_tests": 0,
  "total_time": "5.45ms",
  "results": [
    {
      "name": "Database Schema Validation",
      "passed": true,
      "duration": "0.13ms",
      "details": {
        "tables_validated": 6
      }
    },
    {
      "name": "Performance Benchmarks",
      "passed": true,
      "duration": "2.23ms",
      "details": {
        "entities_created": 100,
        "creation_duration_ms": 1.774625,
        "creation_per_entity_ms": 0.01774625,
        "listing_duration_ms": 0.250458
      }
    }
  ]
}
```

### Example Usage

```bash
# Basic test run
go run cmd/integration-test/main.go

# Verbose output with JSON report
go run cmd/integration-test/main.go -v -output test-results.json

# Use persistent database for debugging
go run cmd/integration-test/main.go -db integration-test.db -v

# CI/CD integration
go run cmd/integration-test/main.go -output ci-results.json
if [ $? -eq 0 ]; then echo "Tests passed"; else echo "Tests failed"; fi
```

## Building and Installing Tools

### Build All Tools

```bash
# Build to ./bin/ directory
make tools-build

# Install to Go bin directory
make tools-install
```

### Individual Tool Building

```bash
# Build specific tools
go build -o bin/dbinspect cmd/dbinspect/main.go
go build -o bin/dbseed cmd/dbseed/main.go
go build -o bin/dashboard cmd/dashboard/main.go
go build -o bin/integration-test cmd/integration-test/main.go
```

### Using Installed Tools

After `make tools-install`, you can use the tools directly:

```bash
# Database inspection
dbinspect -db libretto-dev.db -cmd projects

# Database seeding
dbseed -db test.db -preset fantasy -clean

# Web dashboard
dashboard -db libretto-dev.db -port 8080

# Integration testing
integration-test -v -output results.json
```

## Troubleshooting

### Common Issues

**Database not found:**
```bash
# Create database first
make db-seed
```

**Permission denied:**
```bash
# Check file permissions
ls -la *.db
chmod 644 libretto-dev.db
```

**Port already in use:**
```bash
# Use different port
go run cmd/dashboard/main.go -db libretto-dev.db -port 3000
```

**SQLite locked:**
```bash
# Check for running processes
lsof libretto-dev.db
pkill -f dashboard
```

### Debug Mode

Enable verbose logging for all tools:

```bash
LOG_LEVEL=DEBUG go run cmd/dbinspect/main.go -db libretto-dev.db -cmd projects -v
LOG_LEVEL=DEBUG go run cmd/dashboard/main.go -db libretto-dev.db
```

### Performance Issues

```bash
# Check database size
ls -lh *.db

# Vacuum database
sqlite3 libretto-dev.db "VACUUM;"

# Monitor memory usage
go run cmd/integration-test/main.go -v | grep -i memory
```

This comprehensive CLI tools documentation provides everything needed to effectively use the Libretto development and monitoring tools.