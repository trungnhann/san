-- name: CreatePost :one
INSERT INTO posts (
    id, user_id, title, slug, abstract, body, published, publish_date, location, lat, lon, locale, tags
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: GetPostBySlug :one
SELECT * FROM posts
WHERE slug = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
SET
    title = COALESCE(sqlc.narg('title'), title),
    slug = COALESCE(sqlc.narg('slug'), slug),
    abstract = COALESCE(sqlc.narg('abstract'), abstract),
    body = COALESCE(sqlc.narg('body'), body),
    published = COALESCE(sqlc.narg('published'), published),
    publish_date = COALESCE(sqlc.narg('publish_date'), publish_date),
    location = COALESCE(sqlc.narg('location'), location),
    lat = COALESCE(sqlc.narg('lat'), lat),
    lon = COALESCE(sqlc.narg('lon'), lon),
    locale = COALESCE(sqlc.narg('locale'), locale),
    tags = COALESCE(sqlc.narg('tags'), tags),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: ListPostsByUserID :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
