CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS feeds (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    name           TEXT NOT NULL UNIQUE,
    url            TEXT NOT NULL,
    last_polled_at TIMESTAMP,  
    last_changed_at TIMESTAMP,   
    etag           TEXT,           
    last_modified  TEXT            
);

CREATE INDEX IF NOT EXISTS idx_feeds_last_polled ON feeds (last_polled_at NULLS FIRST);