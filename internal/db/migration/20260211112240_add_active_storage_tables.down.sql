DROP TABLE IF EXISTS storage_attachments;
DROP TABLE IF EXISTS storage_blobs;

-- Restore image column to users
ALTER TABLE users ADD COLUMN image VARCHAR;
