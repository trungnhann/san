-- name: CreateBlob :one
INSERT INTO storage_blobs (id, key, filename, content_type, byte_size, checksum, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CreateAttachment :one
INSERT INTO storage_attachments (name, record_type, record_id, blob_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (record_type, record_id, name) 
DO UPDATE SET blob_id = EXCLUDED.blob_id, created_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetAttachment :one
SELECT 
    a.*,
    b.key,
    b.filename,
    b.content_type,
    b.byte_size
FROM storage_attachments a
JOIN storage_blobs b ON a.blob_id = b.id
WHERE a.record_type = $1 AND a.record_id = $2 AND a.name = $3
LIMIT 1;

-- name: GetAttachmentByRecord :one
SELECT id, blob_id, record_type, record_id, name
FROM storage_attachments
WHERE record_type = $1 AND record_id = $2 AND name = $3
LIMIT 1;

-- name: GetBlob :one
SELECT * FROM storage_blobs
WHERE id = $1;

-- name: DeleteAttachment :exec
DELETE FROM storage_attachments
WHERE id = $1;

-- name: DeleteBlob :exec
DELETE FROM storage_blobs
WHERE id = $1;
