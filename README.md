# Libretto Narrative Engine

A revolutionary multi-agent narrative orchestration engine where users act as "Conductors" directing specialized AI agents to build interconnected Living Narrative graphs.

## 🌟 Key Innovation: Cross-Project Entity Continuity

**"Elena must always be Elena"** - Characters maintain their identity and evolution across multiple related projects, enabling:

- **Epic Multi-Book Series**: Elena from Book 1 is the same Elena in Book 3, carrying forward her complete character development
- **Shared Universes**: Characters appear across spin-offs, crossovers, and collaborative stories
- **Narrative Continuity**: Relationships, skills, and character growth persist across the entire story universe
- **Collaborative World-Building**: Multiple creators can share characters while maintaining consistency

[Learn more about Cross-Project Entity Continuity →](docs/cross-project-continuity.md)

## Quick Start

### Prerequisites

- Go 1.22.5+
- SQLite3
- Make

### Setup Development Environment

```bash
# Clone and setup
git clone <repository-url>
cd libretto

# Setup development environment (creates database, builds tools)
make dev-setup

# Launch web dashboard
make dashboard
# Visit http://localhost:8080
```

### Basic Usage

```bash
# Run all tests
make test

# Inspect database
make db-inspect-projects
make db-inspect-graph

# Clean and reseed database
make db-clean
```

## Architecture

Libretto is a **local-first**, **event-driven** narrative engine that generates prose from a structured graph model rather than traditional document editing.

### Core Concepts

- **The Conductor, Not the Typist**: Users provide high-level creative directives rather than typing prose directly
- **Living Narrative**: A versioned, validated graph of scenes, arcs, characters, settings, and relationships
- **Multi-agent AI**: Event-driven agents (Plot Weaver, Thematic Steward, etc.) collaborate over a Narrative Event Bus
- **Generated Prose**: Prose is a read-only view generated from the model; refinements happen via Tuners or graph edits

### Technology Stack

- **Backend**: Go 1.22.5 - Performance & concurrency for multi-agent architecture
- **Desktop Framework**: Wails - Go-to-UI bridge, more efficient than Electron
- **Frontend**: React + TypeScript - Mature ecosystem with developer familiarity
- **UI Components**: shadcn/ui + Tailwind CSS - Maximum development velocity
- **Database**: SQLite - Zero-config local-first storage
- **Database Layer**: sqlc - Type-safe SQL code generation from raw queries
- **Vector DB**: sqlite-vec - Local embeddings and similarity search for RAG
- **Build System**: Bazel (Go) + pnpm (frontend) + buf (protobuf codegen)

## Development

### Make Targets

#### Build & Development
```bash
make build            # Build entire project
make proto            # Generate protobuf code
make sqlc             # Generate type-safe database code
make wails-dev        # Desktop app development with live reload
```

#### Testing
```bash
make test             # Run all tests (unit + integration)
make test-unit        # Run unit tests only
make test-integration # Run comprehensive integration tests
make test-coverage    # Generate coverage report
make test-watch       # Run tests in watch mode (requires entr)
```

#### Database & Tools
```bash
make db-seed          # Seed database with fantasy story
make db-inspect       # Interactive database inspection
make db-clean         # Clean and reseed database
make dashboard        # Launch web dashboard
make tools-build      # Build CLI tools
```

#### Monitoring & Documentation
```bash
make monitoring-start # Start monitoring dashboard
make docs-serve       # Serve documentation locally
```

### Project Structure

```
/
├── apps/desktop/           # Wails desktop application
├── cmd/                   # Command-line tools and binaries
│   ├── dbinspect/         # Database inspection CLI
│   ├── dbseed/            # Database seeding tool
│   ├── dashboard/         # Web monitoring dashboard
│   └── integration-test/  # Integration test suite
├── internal/              # Private application code
│   ├── db/               # Database models and queries (sqlc generated)
│   ├── types/            # Domain types and data models
│   ├── monitoring/       # Logging and metrics
│   └── repository/       # Data access layer
├── docs/                 # Documentation
└── scripts/              # Development and CI scripts
```

## Testing & Quality Assurance

### Test Coverage

- **Unit Tests**: 32 tests covering all database operations and data models
- **Integration Tests**: 11 comprehensive end-to-end test scenarios
- **Performance Tests**: Benchmarks for entity creation and query performance
- **Data Integrity Tests**: Foreign key constraints and cascade operations

### Test Categories

1. **Database Schema Validation** - Ensures all tables and constraints exist
2. **CRUD Operations** - Tests all create, read, update, delete operations
3. **Graph Version Management** - Tests versioning and working set switching
4. **Entity Management** - Tests all entity types (Scene, Character, Location, etc.)
5. **Relationship Management** - Tests typed connections between entities
6. **Annotation System** - Tests AI agent annotations and metadata
7. **Data Model Validation** - Tests JSON serialization of complex types
8. **Complex Narrative Graph** - Tests complete story creation workflow
9. **Performance Benchmarks** - Validates response times and throughput
10. **Data Integrity** - Tests foreign key constraints and cascade deletes
11. **Concurrent Operations** - Tests system behavior under load

### Running Tests

```bash
# Quick test run
make test

# Detailed test output
make test-unit
make test-integration

# Coverage analysis
make test-coverage
open coverage.html
```

## Database Tools

### Database Inspection (`dbinspect`)

Interactive CLI for examining database state:

```bash
# List all projects
make db-inspect-projects

# Show database schema
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd schema

# View narrative graph
make db-inspect-graph

# Show statistics
make db-inspect-stats

# Inspect specific entities
go run cmd/dbinspect/main.go -db libretto-dev.db -cmd entities -project <project-id> -v
```

### Database Seeding (`dbseed`)

Creates realistic test data:

```bash
# Seed with fantasy story
make db-seed

# Custom database and preset
go run cmd/dbseed/main.go -db custom.db -preset fantasy -clean
```

**Generated Content:**
- Complete project with theme and genre
- Multiple scenes across story acts
- Characters with detailed profiles and arcs
- Locations with atmospheric descriptions
- Themes with symbolic elements
- Relationships between all entities
- AI agent annotations with realistic analysis

## Web Dashboard

Interactive web interface for monitoring and visualization:

```bash
# Launch dashboard
make dashboard
# Visit http://localhost:8080
```

**Features:**
- Project overview with statistics
- Interactive entity listings
- D3.js graph visualization with:
  - Color-coded entity types
  - Draggable nodes
  - Relationship links
  - Node information panel
- Real-time metrics and statistics

## Monitoring & Logging

### Structured Logging

The system includes comprehensive structured logging:

```go
logger := monitoring.NewLogger("narrative-engine")

// Context-aware logging
logger.Info(ctx, "Processing scene", 
    monitoring.String("scene_id", sceneID),
    monitoring.String("operation", "create"))

// Operation timing
timer := logger.StartOperation(ctx, "create_narrative_graph")
// ... do work ...
timer.Complete(ctx, "Successfully created graph")
```

### Database Metrics

Track database performance:

```go
dbMetrics := monitoring.NewDatabaseMetrics(logger)
dbMetrics.RecordQuery(ctx, "CreateEntity", duration, err)
dbMetrics.RecordEntityOperation(ctx, "Scene", "create", duration, err)
```

## Data Model

### Core Entities

- **Projects**: Top-level narrative containers
- **Graph Versions**: Versioned snapshots with working set management
- **Entities**: Core narrative elements (Scene, Character, Location, Theme, PlotPoint, Arc)
- **Relationships**: Typed connections between entities
- **Annotations**: AI agent analysis and metadata

### Entity Types

1. **Scene**: Narrative scenes with content, emotional tone, pacing
2. **Character**: Characters with personality, voice, and character arcs
3. **Location**: Settings with atmosphere and physical details
4. **Theme**: Thematic elements with questions and symbols
5. **PlotPoint**: Key story moments and turning points
6. **Arc**: Character and plot arcs spanning multiple scenes

### Relationship Types

- `features`: Scene features Character
- `occurs_at`: Scene occurs at Location
- `influences`: Entity influences Theme
- `conflicts`: Character conflicts with Character
- `supports`: Entity supports Entity
- `precedes`/`follows`: Sequential relationships
- `contains`: Hierarchical relationships

## Performance

### Benchmarks

Current performance metrics (from integration tests):

- **Entity Creation**: ~0.017ms per entity
- **Database Queries**: Sub-millisecond response times
- **Graph Traversal**: Efficient relationship queries
- **Load Testing**: Successfully handles 100+ entities

### Optimization

- **Indexed Queries**: Strategic database indexing
- **Type-Safe Operations**: sqlc-generated queries
- **Connection Pooling**: Efficient database connections
- **JSON Storage**: Flexible entity data with fast serialization

## Contributing

### Development Workflow

1. **Setup**: `make dev-setup`
2. **Code**: Make changes to Go code
3. **Test**: `make test` (runs unit + integration tests)
4. **Verify**: `make dashboard` to visually inspect changes
5. **Commit**: Ensure all tests pass

### Code Standards

- **Go**: Follow standard Go conventions
- **SQL**: Use sqlc for type-safe database operations
- **Testing**: Maintain >90% test coverage
- **Documentation**: Update docs for new features

### Adding New Entity Types

1. Add type constants to `internal/types/entities.go`
2. Create data structure and marshal/unmarshal functions
3. Add unit tests in `internal/types/entities_test.go`
4. Update database seeding in `cmd/dbseed/main.go`
5. Add visualization support in dashboard

## Deployment

### Local Development

```bash
make dev-setup    # Complete development setup
make dashboard    # Launch monitoring interface
```

### Production Considerations

- **Database**: SQLite with WAL mode for better concurrency
- **Monitoring**: Structured JSON logs for log aggregation
- **Backup**: Regular database file backups
- **Performance**: Monitor query performance and optimize indexes

## Documentation

- **[Testing & Monitoring Guide](docs/testing-and-monitoring.md)**: Comprehensive testing documentation
- **[Architecture Decision Records](docs/adr/)**: Design decisions and rationale
- **[API Documentation](docs/api/)**: Generated API documentation

## License

[License information]

## Support

For questions and support:
- Check the documentation in `docs/`
- Run `make help` for available commands
- Use `make db-inspect` for database exploration
- Launch `make dashboard` for visual monitoring