-- Living Narrative data model schema
-- This extends the basic scenes table to support the full narrative graph model
-- with entities, relationships, versioning, and annotations

-- Projects table - top-level container for narratives
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    theme TEXT,
    genre TEXT,
    description TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Graph versions table - supports versioning and branching of narrative graphs
CREATE TABLE graph_versions (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    parent_version_id TEXT,
    name TEXT DEFAULT '',
    description TEXT DEFAULT '',
    is_working_set BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_version_id) REFERENCES graph_versions(id)
);

-- Entities table - core narrative elements (scenes, characters, locations, themes, etc.)
CREATE TABLE entities (
    id TEXT PRIMARY KEY,
    version_id TEXT NOT NULL,
    entity_type TEXT NOT NULL, -- Scene, Character, Location, Theme, PlotPoint, etc.
    name TEXT NOT NULL DEFAULT '',
    data JSON NOT NULL, -- Flexible storage for entity-specific data
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (version_id) REFERENCES graph_versions(id) ON DELETE CASCADE
);

-- Relationships table - typed connections between entities
CREATE TABLE relationships (
    id TEXT PRIMARY KEY,
    version_id TEXT NOT NULL,
    from_entity_id TEXT NOT NULL,
    to_entity_id TEXT NOT NULL,
    relationship_type TEXT NOT NULL, -- contains, advances, features, occurs_at, influences, etc.
    properties JSON, -- Optional relationship metadata
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (version_id) REFERENCES graph_versions(id) ON DELETE CASCADE,
    FOREIGN KEY (from_entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    FOREIGN KEY (to_entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    UNIQUE(version_id, from_entity_id, to_entity_id, relationship_type)
);

-- Annotations table - agent-generated metadata and user comments
CREATE TABLE annotations (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    annotation_type TEXT NOT NULL, -- emotional_analysis, thematic_score, continuity_check, etc.
    content TEXT NOT NULL,
    metadata JSON, -- Additional structured data
    agent_name TEXT, -- Which agent created this annotation
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE
);

-- Indexes for efficient queries
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);
CREATE INDEX idx_graph_versions_project_id ON graph_versions(project_id);
CREATE INDEX idx_graph_versions_working_set ON graph_versions(project_id, is_working_set) WHERE is_working_set = TRUE;
CREATE INDEX idx_entities_version_id ON entities(version_id);
CREATE INDEX idx_entities_type ON entities(version_id, entity_type);
CREATE INDEX idx_relationships_version_id ON relationships(version_id);
CREATE INDEX idx_relationships_from_entity ON relationships(from_entity_id);
CREATE INDEX idx_relationships_to_entity ON relationships(to_entity_id);
CREATE INDEX idx_relationships_type ON relationships(version_id, relationship_type);
CREATE INDEX idx_annotations_entity_id ON annotations(entity_id);
CREATE INDEX idx_annotations_type ON annotations(entity_id, annotation_type);

-- Triggers to update updated_at timestamps
CREATE TRIGGER update_projects_updated_at 
    AFTER UPDATE ON projects
    FOR EACH ROW
BEGIN
    UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_entities_updated_at 
    AFTER UPDATE ON entities
    FOR EACH ROW
BEGIN
    UPDATE entities SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Constraint to ensure only one working set per project
CREATE UNIQUE INDEX idx_unique_working_set_per_project 
ON graph_versions(project_id) 
WHERE is_working_set = TRUE;