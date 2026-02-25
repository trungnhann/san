CREATE TABLE storage_blobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR NOT NULL,
    filename VARCHAR NOT NULL,
    content_type VARCHAR,
    metadata JSONB,
    byte_size BIGINT NOT NULL,
    checksum VARCHAR,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE storage_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    record_type VARCHAR NOT NULL,
    record_id VARCHAR NOT NULL,
    blob_id UUID NOT NULL REFERENCES storage_blobs(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (record_type, record_id, name)
);

-- Index for polymorphic queries
CREATE INDEX idx_storage_attachments_record ON storage_attachments(record_type, record_id, name);

-- Remove image column from users (data migration might be needed if you have existing data)
ALTER TABLE users DROP COLUMN image;
