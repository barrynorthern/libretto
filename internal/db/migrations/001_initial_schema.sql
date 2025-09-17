-- Initial schema for Libretto scenes
-- This creates the basic scenes table for the MVP

CREATE TABLE scenes (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for efficient listing by creation time
CREATE INDEX idx_scenes_created_at ON scenes(created_at DESC);

-- Trigger to update updated_at on modifications
CREATE TRIGGER update_scenes_updated_at 
    AFTER UPDATE ON scenes
    FOR EACH ROW
BEGIN
    UPDATE scenes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
