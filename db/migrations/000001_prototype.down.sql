-- Drop Indexes
DROP INDEX idx_relationships_target;
DROP INDEX idx_relationships_source;
DROP INDEX idx_embeddings_embedding;
DROP INDEX idx_embeddings_content_id;
DROP INDEX idx_content_checksum;
DROP INDEX idx_content_source_id;
DROP INDEX idx_sources_status;

-- Drop Hypertables and Tables
-- Drop relationships table
DROP TABLE IF EXISTS relationships CASCADE;

-- Drop embeddings table
DROP TABLE IF EXISTS embeddings CASCADE;

-- Drop content table
DROP TABLE IF EXISTS content CASCADE;

-- Drop sources table
DROP TABLE IF EXISTS sources CASCADE;