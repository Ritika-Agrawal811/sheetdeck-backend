-- name: CreateCheatsheet :exec
INSERT INTO cheatsheets (slug, title, category, subcategory, image_url)
VALUES ($1, $2, $3, $4, $5);


-- name: GetCheatsheetByID :one
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at
FROM cheatsheets
WHERE id = $1;

-- name: GetCheatsheetBySlug :one
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at
FROM cheatsheets
WHERE slug = $1;    

-- name: ListCheatsheets :many
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at
FROM cheatsheets
WHERE (sqlc.narg(category)::category IS NULL OR category = sqlc.narg(category))
  AND (sqlc.narg(subcategory)::subcategory IS NULL OR subcategory = sqlc.narg(subcategory))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
