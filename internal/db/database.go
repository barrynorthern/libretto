package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Database wraps the sqlc generated Queries with connection management
type Database struct {
	db      *sql.DB
	queries *Queries
}

// NewDatabase creates a new Database instance with SQLite connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{
		db:      db,
		queries: New(db),
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Queries returns the sqlc generated queries interface
func (d *Database) Queries() *Queries {
	return d.queries
}

// Migrate runs all pending migrations
func (d *Database) Migrate(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := d.db.ExecContext(ctx, createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := d.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all migration files
	migrationFiles, err := d.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrationFiles {
		if appliedMigrations[migration] {
			continue // Already applied
		}

		if err := d.applyMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration, err)
		}
	}

	return nil
}

func (d *Database) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func (d *Database) getMigrationFiles() ([]string, error) {
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		migrations = append(migrations, entry.Name())
	}

	sort.Strings(migrations)
	return migrations, nil
}

func (d *Database) applyMigration(ctx context.Context, filename string) error {
	content, err := migrationFiles.ReadFile(filepath.Join("migrations", filename))
	if err != nil {
		return err
	}

	// Execute migration in a transaction
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return err
	}

	// Record migration as applied
	if _, err := tx.ExecContext(ctx, 
		"INSERT INTO schema_migrations (version) VALUES (?)", 
		filename); err != nil {
		return err
	}

	return tx.Commit()
}
