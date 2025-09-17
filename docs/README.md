# Libretto Documentation

Welcome to the Libretto Narrative Engine documentation. This directory contains comprehensive guides for developers, users, and contributors.

## Quick Start

New to Libretto? Start here:

1. **[README.md](../README.md)** - Project overview and quick setup
2. **[Quick Reference](quick-reference.md)** - Essential commands and workflows
3. **[Development Guide](development-guide.md)** - Comprehensive development documentation

## Documentation Structure

### Core Documentation

| Document | Purpose | Audience |
|----------|---------|----------|
| **[README.md](../README.md)** | Project overview, architecture, quick start | All users |
| **[Quick Reference](quick-reference.md)** | Essential commands, common workflows | Developers |
| **[Development Guide](development-guide.md)** | Comprehensive development documentation | Contributors |

### Specialized Guides

| Document | Purpose | Audience |
|----------|---------|----------|
| **[Testing & Monitoring](testing-and-monitoring.md)** | Complete testing and observability guide | QA, DevOps |
| **[CLI Tools](cli-tools.md)** | Detailed CLI tool documentation | Developers, Ops |

### Reference Materials

| Document | Purpose | Audience |
|----------|---------|----------|
| **[Architecture Decision Records](adr/)** | Design decisions and rationale | Architects, Senior Devs |
| **[API Documentation](api/)** | Generated API documentation | Integrators |

## Getting Started Workflows

### For New Developers

```bash
# 1. Clone and setup
git clone <repository-url>
cd libretto
make dev-setup

# 2. Explore the system
make dashboard
# Visit http://localhost:8080

# 3. Run tests
make test

# 4. Read documentation
make docs-serve
# Visit http://localhost:8000
```

### For Contributors

```bash
# 1. Setup development environment
make dev-setup

# 2. Make changes
# ... edit code ...

# 3. Test changes
make test-unit
make test-integration

# 4. Verify visually
make dashboard

# 5. Build tools
make tools-build
```

### For QA/Testing

```bash
# 1. Run comprehensive tests
make test

# 2. Generate coverage report
make test-coverage
open coverage.html

# 3. Run integration tests with output
make test-integration

# 4. Inspect database state
make db-inspect-stats
make db-inspect-graph
```

### For Operations/DevOps

```bash
# 1. Monitor system health
make dashboard

# 2. Inspect database
make db-inspect-schema
make db-inspect-stats

# 3. Performance testing
make test-integration

# 4. Log analysis
LOG_LEVEL=DEBUG make test-integration
```

## Key Concepts

### Architecture Overview

Libretto is a **local-first**, **event-driven** narrative engine:

- **Local-First**: SQLite database, no cloud dependencies
- **Event-Driven**: Multi-agent AI system with publish-subscribe patterns
- **Graph-Based**: Narrative stored as entities and relationships
- **Generated Prose**: Content generated from structured data

### Core Components

1. **Living Narrative Graph**: Versioned graph of scenes, characters, locations, themes
2. **Multi-Agent System**: AI agents that analyze and enhance the narrative
3. **Desktop Application**: Wails-based UI for conductors
4. **CLI Tools**: Database inspection, seeding, monitoring
5. **Web Dashboard**: Visual monitoring and graph exploration

### Data Model

- **Projects**: Top-level narrative containers
- **Graph Versions**: Versioned snapshots with working set management
- **Entities**: Core narrative elements (Scene, Character, Location, Theme, PlotPoint, Arc)
- **Relationships**: Typed connections between entities
- **Annotations**: AI agent analysis and metadata

## Development Workflows

### Daily Development

1. **Setup**: `make dev-setup`
2. **Code**: Make changes
3. **Test**: `make test-unit`
4. **Verify**: `make dashboard`
5. **Integration**: `make test-integration`

### Feature Development

1. **Plan**: Define entity types, relationships, or agents
2. **Implement**: Add code with tests
3. **Test**: Comprehensive testing
4. **Document**: Update documentation
5. **Review**: Visual verification with dashboard

### Debugging

1. **Inspect**: `make db-inspect-*` commands
2. **Monitor**: `make dashboard` for visual analysis
3. **Test**: `make test-integration` for comprehensive validation
4. **Logs**: `LOG_LEVEL=DEBUG` for detailed logging

## Testing Strategy

### Test Levels

1. **Unit Tests** (`make test-unit`): Individual component testing
2. **Integration Tests** (`make test-integration`): End-to-end workflows
3. **Performance Tests**: Benchmarks and load testing
4. **Manual Testing**: Dashboard and CLI tool exploration

### Test Coverage

- **32 Unit Tests**: Database operations and data models
- **11 Integration Tests**: Complete system workflows
- **Performance Benchmarks**: Entity creation, query performance
- **Data Integrity**: Foreign key constraints, cascade operations

### Quality Gates

- **Test Coverage**: >90% for critical components
- **Performance**: <20ms per entity operation
- **Integration**: All 11 integration tests must pass
- **Data Integrity**: No foreign key violations

## Tools and Utilities

### Make Targets

| Category | Targets | Purpose |
|----------|---------|---------|
| **Testing** | `test`, `test-unit`, `test-integration`, `test-coverage` | Run various test suites |
| **Database** | `db-seed`, `db-inspect-*`, `db-clean` | Database management |
| **Tools** | `tools-build`, `tools-install`, `dashboard` | CLI tools and monitoring |
| **Development** | `dev-setup`, `build`, `sqlc`, `proto` | Development workflow |

### CLI Tools

| Tool | Purpose | Documentation |
|------|---------|---------------|
| **dbinspect** | Database inspection and analysis | [CLI Tools](cli-tools.md#database-inspector-dbinspect) |
| **dbseed** | Database seeding with test data | [CLI Tools](cli-tools.md#database-seeder-dbseed) |
| **dashboard** | Web-based monitoring interface | [CLI Tools](cli-tools.md#web-dashboard-dashboard) |
| **integration-test** | Comprehensive test suite | [CLI Tools](cli-tools.md#integration-test-suite-integration-test) |

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `DB_FILE` | Database file path | `libretto-dev.db` |
| `DB_PRESET` | Database preset for seeding | `fantasy` |
| `DASHBOARD_PORT` | Web dashboard port | `8080` |
| `LOG_LEVEL` | Logging verbosity | `INFO` |

## Performance and Monitoring

### Performance Metrics

Current benchmarks from integration tests:
- **Entity Creation**: ~0.026ms per entity
- **Database Queries**: Sub-millisecond response times
- **Graph Traversal**: Efficient relationship queries
- **Load Testing**: Successfully handles 100+ entities

### Monitoring Tools

1. **Web Dashboard**: Real-time visualization and statistics
2. **Structured Logging**: JSON logs with context and timing
3. **Database Metrics**: Query performance and operation tracking
4. **Integration Tests**: Automated performance validation

### Observability

- **Logs**: Structured JSON logging with context
- **Metrics**: Database operation timing and counts
- **Visualization**: Interactive graph exploration
- **Health Checks**: Comprehensive integration test suite

## Contributing

### Code Standards

- **Go**: Follow standard Go conventions and gofmt
- **SQL**: Use sqlc for type-safe database operations
- **Testing**: Maintain >90% test coverage for critical components
- **Documentation**: Update docs for new features

### Pull Request Process

1. **Setup**: `make dev-setup`
2. **Develop**: Make changes with tests
3. **Test**: `make test` (all tests must pass)
4. **Document**: Update relevant documentation
5. **Review**: Submit PR with clear description

### Adding Features

See the [Development Guide](development-guide.md) for detailed instructions on:
- Adding new entity types
- Creating new relationship types
- Implementing new annotation types
- Extending the database schema

## Troubleshooting

### Common Issues

| Issue | Solution | Documentation |
|-------|----------|---------------|
| Database not found | `make db-seed` | [Quick Reference](quick-reference.md#troubleshooting) |
| Tests failing | `make db-clean && make test` | [Testing Guide](testing-and-monitoring.md#troubleshooting) |
| Port in use | `DASHBOARD_PORT=3000 make dashboard` | [CLI Tools](cli-tools.md#troubleshooting) |
| Performance issues | Check database size, run `VACUUM` | [Development Guide](development-guide.md#debugging-common-issues) |

### Debug Mode

Enable verbose logging for all tools:
```bash
LOG_LEVEL=DEBUG make test-integration
LOG_LEVEL=DEBUG make dashboard
```

### Getting Help

1. **Documentation**: Check relevant guides in this directory
2. **Commands**: Run `make help` for available targets
3. **Database**: Use `make db-inspect-*` for database exploration
4. **Monitoring**: Launch `make dashboard` for visual inspection
5. **Testing**: Run `make test-integration -v` for detailed test output

## License and Support

- **License**: [License information]
- **Issues**: Report issues via GitHub Issues
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: All documentation is in this `docs/` directory

---

This documentation is maintained alongside the codebase. When making changes, please update the relevant documentation files to keep them current and accurate.