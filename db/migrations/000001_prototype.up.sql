-- install extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto CASCADE;
CREATE EXTENSION IF NOT EXISTS ai CASCADE;
CREATE EXTENSION IF NOT EXISTS vectorscale CASCADE;

CREATE TYPE source_status AS ENUM ('inactive', 'active', 'pending', 'error');

-- create tables
CREATE TABLE sources (
    id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    kind VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    credentials BYTEA NOT NULL,
    status source_status DEFAULT 'inactive',
    sync_frequency INTERVAL,
    synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT sources_pk UNIQUE (id, created_at)
);

CREATE TABLE content (
    id UUID NOT NULL,
    source_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    raw_data TEXT NOT NULL,
    metadata JSONB,
    content_type VARCHAR(50),
    version INT DEFAULT 1,
    checksum VARCHAR(64),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT content_pk UNIQUE (id, created_at)
);

CREATE TABLE embeddings (
    id UUID NOT NULL,
    content_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    chunk_text TEXT NOT NULL,
    embedding VECTOR(768),
    metadata JSONB,
    entity_type VARCHAR(50),
    version INT DEFAULT 1,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT embeddings_pk UNIQUE (id, created_at)
);

CREATE TABLE relationships (
    id UUID NOT NULL,
    source_embedding_id UUID NOT NULL,
    target_embedding_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    relationship_type VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT relationships_pk UNIQUE (id, created_at)
);

-- Hypertables
SELECT create_hypertable('sources', 'created_at');
SELECT create_hypertable('content', 'created_at');
SELECT create_hypertable('embeddings', 'created_at');
SELECT create_hypertable('relationships', 'created_at');

-- Indexes
CREATE INDEX idx_sources_id ON sources (id);
CREATE INDEX idx_sources_status ON sources (status);
CREATE INDEX idx_content_id ON content (id);
CREATE INDEX idx_content_source_id ON content (source_id);
CREATE INDEX idx_content_checksum ON content (checksum);
CREATE INDEX idx_embeddings_id ON embeddings (id);
CREATE INDEX idx_embeddings_content_id ON embeddings (content_id);
CREATE INDEX idx_embeddings_embedding ON embeddings USING ivfflat (embedding vector_l2_ops) WITH (lists = 100);
CREATE INDEX idx_relationships_id ON relationships (id);
CREATE INDEX idx_relationships_source ON relationships (source_embedding_id);
CREATE INDEX idx_relationships_target ON relationships (target_embedding_id);
